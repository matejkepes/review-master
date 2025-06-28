package google_my_business_api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"google_my_business/ai_service/review"
	"google_my_business/database"
	"google_my_business/email_service"
	"google_my_business/shared"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

// Function variable aliases to make testing easier
var (
	// GetAccountsFunc is the function used to fetch Google accounts
	GetAccountsFunc = GetAccounts

	// GetLocationsFunc is the function used to fetch locations
	GetLocationsFunc = GetLocations

	// FetchReviewsForMonthFunc is the function used to fetch reviews for a specific month
	FetchReviewsForMonthFunc = FetchReviewsForMonth

	// ConvertHTMLToPDFFunc is the function used to convert HTML to PDF
	ConvertHTMLToPDFFunc = ConvertHTMLToPDF
)

// DBInterface defines the database operations needed for monthly analysis
type DBInterface interface {
	// Client returns the underlying database client
	Client() *sql.DB

	// GetClientsWithMonthlyReviewAnalysisEnabled returns clients with monthly analysis enabled
	GetClientsWithMonthlyReviewAnalysisEnabled() ([]database.ClientWithMonthlyReviewAnalysis, error)

	// GetClientReportByClientAndPeriod checks if a report exists for a time period
	GetClientReportByClientAndPeriod(clientID int, periodStart, periodEnd time.Time) (*shared.ClientReportData, error)

	// DeleteClientReport deletes an existing report
	DeleteClientReport(reportID int64) error

	// SaveClientReport saves a new report
	SaveClientReport(clientID int, periodStart, periodEnd time.Time, locationResults []byte) (int64, error)
}

// ReviewAnalyzerInterface defines the operations needed from the analyzer
type ReviewAnalyzerInterface interface {
	// Analyze processes a batch of reviews
	Analyze(batch review.ReviewBatch) (*shared.AnalysisResult, error)
}

// FailedClientInfo contains information about a client that failed processing
type FailedClientInfo struct {
	ClientID   int
	ClientName string
	Error      string
}

// ProcessingSummary contains statistics about the processing run
type ProcessingSummary struct {
	TargetMonth          time.Time
	ClientsProcessed     int
	ClientsSucceeded     int
	ClientsFailed        int
	ClientsSkipped       int   // Clients already having reports
	SkippedClientIDs     []int // IDs of skipped clients
	TotalLocations       int
	TotalReviews         int
	TotalReviewsAnalyzed int
	PDFsGenerated        int // Track how many PDFs were generated
	EmailsSent           int // Track how many emails were sent
	ElapsedTime          time.Duration
	FailedClients        []FailedClientInfo
}

// LocationReport is a structure to hold analysis results for a location
type LocationReport struct {
	LocationID      string `json:"locationID"`
	LocationName    string `json:"locationName"`
	LocationAddress string `json:"locationAddress"`
	ReviewCount     int    `json:"reviewCount"`
	// Embed the AnalysisResult fields directly to avoid double nesting
	Analysis shared.Analysis          `json:"analysis"`
	Metadata shared.AnalysisMetadata  `json:"metadata"`
	Reviews  []map[string]interface{} `json:"reviews,omitempty"`
}

// AnalyzeClientReviews processes reviews for all eligible clients
func AnalyzeClientReviews(db DBInterface, httpClient *http.Client, analyzer ReviewAnalyzerInterface,
	targetMonth *time.Time, forceReprocess bool, emailSvc email_service.EmailService) (ProcessingSummary, error) {
	startTime := time.Now()

	// Initialize processing summary
	summary := ProcessingSummary{
		SkippedClientIDs: []int{},
		FailedClients:    []FailedClientInfo{},
	}

	// If targetMonth is nil, use previous month
	if targetMonth == nil {
		now := time.Now()
		lastMonth := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, time.UTC)
		targetMonth = &lastMonth
	}

	summary.TargetMonth = *targetMonth

	// Calculate period start and end dates (first day of month to first day of next month)
	periodStart := time.Date(targetMonth.Year(), targetMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(targetMonth.Year(), targetMonth.Month()+1, 1, 0, 0, 0, 0, time.UTC)

	log.Printf("Monthly analysis started for %s", targetMonth.Format("2006-01"))

	// Step 1: Get clients with monthly review analysis enabled
	clients, err := db.GetClientsWithMonthlyReviewAnalysisEnabled()
	if err != nil {
		return summary, err
	}

	log.Printf("Found %d clients with monthly review analysis enabled", len(clients))

	// Early exit if no clients
	if len(clients) == 0 {
		summary.ElapsedTime = time.Since(startTime)
		return summary, nil
	}

	// Step 2: Process each client
	for _, clientInfo := range clients {
		summary.ClientsProcessed++

		log.Printf("Processing client: %s (ID: %d)", clientInfo.ClientName, clientInfo.ClientID)

		// Check if client already has a report for this month
		existingReport, err := db.GetClientReportByClientAndPeriod(clientInfo.ClientID, periodStart, periodEnd)
		if err != nil {
			// Database error
			log.Printf("Error checking for existing report for client %d: %v", clientInfo.ClientID, err)
			summary.ClientsFailed++
			summary.FailedClients = append(summary.FailedClients, FailedClientInfo{
				ClientID:   clientInfo.ClientID,
				ClientName: clientInfo.ClientName,
				Error:      err.Error(),
			})
			continue
		}

		// If report exists and we're not forcing reprocessing, skip this client
		if existingReport != nil && !forceReprocess {
			log.Printf("Skipping client %s (ID: %d) - report already exists for %s",
				clientInfo.ClientName, clientInfo.ClientID, targetMonth.Format("2006-01"))
			summary.ClientsSkipped++
			summary.SkippedClientIDs = append(summary.SkippedClientIDs, clientInfo.ClientID)
			continue
		}

		// If we're reprocessing, delete existing report first
		if existingReport != nil && forceReprocess {
			log.Printf("Deleting existing report for client %s (ID: %d) for %s",
				clientInfo.ClientName, clientInfo.ClientID, targetMonth.Format("2006-01"))

			err = db.DeleteClientReport(existingReport.ReportID)
			if err != nil {
				log.Printf("Error deleting existing report for client %d: %v", clientInfo.ClientID, err)
				summary.ClientsFailed++
				summary.FailedClients = append(summary.FailedClients, FailedClientInfo{
					ClientID:   clientInfo.ClientID,
					ClientName: clientInfo.ClientName,
					Error:      err.Error(),
				})
				continue
			}
		}

		// Step 3: Get accounts and locations for this client
		// We'll use the direct functions rather than an interface for these
		// since they're tightly coupled to the Google API
		sqlDb := db.Client()

		// Get Google accounts
		accounts := GetAccountsFunc(httpClient)
		if len(accounts) == 0 {
			log.Printf("No Google accounts found for client %d", clientInfo.ClientID)
			summary.ClientsFailed++
			summary.FailedClients = append(summary.FailedClients, FailedClientInfo{
				ClientID:   clientInfo.ClientID,
				ClientName: clientInfo.ClientName,
				Error:      "No Google accounts found",
			})
			continue
		}

		log.Printf("Found %d Google accounts for client %d", len(accounts), clientInfo.ClientID)

		// Collection of location results for this client
		var locationResults []shared.AnalysisResult
		var locationReports []LocationReport

		// Step 4: Process each account
		for _, account := range accounts {
			// Get locations for this account
			// Set report=true to get only locations with report enabled
			allLocations := GetLocationsFunc(httpClient, account, sqlDb, database.LookupModeReport)
			if len(allLocations) == 0 {
				log.Printf("No locations found for account %s", account)
				continue
			}

			// Filter locations to only include those belonging to the current client
			var clientLocations []database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode
			for _, loc := range allLocations {
				// If the location already has the correct client ID, use it directly
				// This is especially useful for testing where we can pre-populate the client ID
				if loc.ClientID == uint64(clientInfo.ClientID) {
					clientLocations = append(clientLocations, loc)
					continue
				}

				// Otherwise, look up each location by name and postal code to verify if it belongs to this client
				dbLocation := database.ConfigFromGoogleMyBusinessLocationNameAndPostalCode(
					sqlDb,
					loc.GoogleMyBusinessLocationName,
					loc.GoogleMyBusinessPostalCode,
					database.LookupModeAnalysis)

				// If the location belongs to this client, add it to our filtered list
				if dbLocation.ClientID == uint64(clientInfo.ClientID) {
					// Use the location from the API but ensure it has the correct client ID and settings
					loc.ClientID = dbLocation.ClientID
					clientLocations = append(clientLocations, loc)
				}
			}

			if len(clientLocations) == 0 {
				log.Printf("No locations found for client %s (ID: %d) in account %s",
					clientInfo.ClientName, clientInfo.ClientID, account)
				continue
			}

			log.Printf("Found %d locations for client %s (ID: %d) in account %s",
				len(clientLocations), clientInfo.ClientName, clientInfo.ClientID, account)
			summary.TotalLocations += len(clientLocations)

			// Step 5: Process all client locations
			for _, location := range clientLocations {
				log.Printf("Processing location: %s", location.GoogleMyBusinessLocationName)

				// Step 5.1: Fetch reviews for this location for the target month
				reviews, _, _, err := FetchReviewsForMonthFunc(httpClient, location, *targetMonth)
				if err != nil {
					log.Printf("Error fetching reviews for location %s: %v", location.GoogleMyBusinessLocationName, err)
					continue
				}

				reviewCount := len(reviews)
				log.Printf("Found %d reviews for location %s", reviewCount, location.GoogleMyBusinessLocationName)
				summary.TotalReviews += reviewCount

				// If no reviews found, add empty report and continue
				if reviewCount == 0 {
					locationReports = append(locationReports, LocationReport{
						LocationID:      location.GoogleMyBusinessLocationPath,
						LocationName:    location.GoogleMyBusinessLocationName,
						LocationAddress: location.GoogleMyBusinessPostalCode,
						ReviewCount:     0,
						Analysis:        shared.Analysis{}, // Empty analysis struct
						Metadata: shared.AnalysisMetadata{
							LocationName: location.GoogleMyBusinessLocationName,
							LocationID:   location.GoogleMyBusinessLocationPath,
							BusinessName: clientInfo.ClientName,
							ClientID:     clientInfo.ClientID,
						},
					})
					continue
				}

				// Step 5.2: Convert reviews to the format expected by the analyzer
				reviewBatch := prepareReviewBatch(reviews, location, clientInfo, periodStart, periodEnd)

				// Step 5.3: Analyze the reviews
				// Check if the batch has any reviews before attempting analysis
				if len(reviewBatch.Reviews) == 0 {
					log.Printf("No valid reviews to analyze for location %s", location.GoogleMyBusinessLocationName)

					// Add an empty report for this location
					locationReports = append(locationReports, LocationReport{
						LocationID:      location.GoogleMyBusinessLocationPath,
						LocationName:    location.GoogleMyBusinessLocationName,
						LocationAddress: location.GoogleMyBusinessPostalCode,
						ReviewCount:     reviewCount,
						Analysis:        shared.Analysis{}, // Empty analysis struct
						Metadata: shared.AnalysisMetadata{
							LocationName: location.GoogleMyBusinessLocationName,
							LocationID:   location.GoogleMyBusinessLocationPath,
							BusinessName: clientInfo.ClientName,
							ClientID:     clientInfo.ClientID,
						},
					})
					continue
				}

				analysisResult, err := analyzer.Analyze(reviewBatch)
				if err != nil {
					log.Printf("Error analyzing reviews for location %s: %v", location.GoogleMyBusinessLocationName, err)
					continue
				}

				// If location name is not properly set in metadata, set it explicitly
				if analysisResult.Metadata.LocationName == "" {
					analysisResult.Metadata.LocationName = location.GoogleMyBusinessLocationName
				}

				// For LocationReport structure, make sure to also set the name there
				locationReport := LocationReport{
					LocationID:      location.GoogleMyBusinessLocationPath,
					LocationName:    location.GoogleMyBusinessLocationName,
					LocationAddress: location.GoogleMyBusinessPostalCode,
					ReviewCount:     len(reviews),
					Analysis:        analysisResult.Analysis, // Extract the Analysis field
					Metadata:        analysisResult.Metadata, // Extract the Metadata field
				}

				summary.TotalReviewsAnalyzed += len(reviews)

				// Add to location results for PDF generation
				locationResults = append(locationResults, *analysisResult)

				// Add the location report to the list
				locationReports = append(locationReports, locationReport)
			}
		}

		// Skip saving if no location reports were generated
		if len(locationReports) == 0 {
			log.Printf("No location reports generated for client %s (ID: %d)", clientInfo.ClientName, clientInfo.ClientID)
			summary.ClientsFailed++
			summary.FailedClients = append(summary.FailedClients, FailedClientInfo{
				ClientID:   clientInfo.ClientID,
				ClientName: clientInfo.ClientName,
				Error:      "No reviews found for any location",
			})
			continue
		}

		// Convert location reports to JSON for storage
		reportJSON, err := json.Marshal(locationReports)
		if err != nil {
			log.Printf("Error marshaling location reports for client %d: %v", clientInfo.ClientID, err)
			summary.ClientsFailed++
			summary.FailedClients = append(summary.FailedClients, FailedClientInfo{
				ClientID:   clientInfo.ClientID,
				ClientName: clientInfo.ClientName,
				Error:      err.Error(),
			})
			continue
		}

		// Save the client report
		reportID, err := db.SaveClientReport(clientInfo.ClientID, periodStart, periodEnd, reportJSON)
		if err != nil {
			log.Printf("Error saving report for client %d: %v", clientInfo.ClientID, err)
			summary.ClientsFailed++
			summary.FailedClients = append(summary.FailedClients, FailedClientInfo{
				ClientID:   clientInfo.ClientID,
				ClientName: clientInfo.ClientName,
				Error:      err.Error(),
			})
			continue
		}

		// Step 6: Generate PDF report
		var pdfContent []byte
		if len(locationResults) > 0 {
			pdfContent, err = generatePDFReport(clientInfo, periodStart, periodEnd, locationResults, reportID)
			if err != nil {
				log.Printf("Error generating PDF for client %d: %v", clientInfo.ClientID, err)
				// Don't fail the entire client process just because PDF generation failed
				// But log the error for investigation
			} else {
				summary.PDFsGenerated++
				log.Printf("PDF report generated for client %d", clientInfo.ClientID)

				// Step 7: Send email with PDF report if we have an email service and PDF content
				if emailSvc != nil && len(pdfContent) > 0 {
					err = sendReportEmail(emailSvc, clientInfo, periodStart, pdfContent)
					if err != nil {
						log.Printf("Error sending report email for client %d: %v", clientInfo.ClientID, err)
						// Don't fail the client processing just because email sending failed
					} else {
						summary.EmailsSent++
						log.Printf("Report email sent for client %d", clientInfo.ClientID)
					}
				}
			}
		}

		// Mark as succeeded
		summary.ClientsSucceeded++
	}

	// Set final elapsed time
	summary.ElapsedTime = time.Since(startTime)

	return summary, nil
}

// generatePDFReport creates a PDF report for a client
func generatePDFReport(clientInfo database.ClientWithMonthlyReviewAnalysis,
	periodStart time.Time, periodEnd time.Time,
	locationResults []shared.AnalysisResult, reportID int64) ([]byte, error) {

	log.Printf("Generating PDF report with %d location results", len(locationResults))

	// Ensure location names are set - this is a failsafe
	for i, loc := range locationResults {
		// If location name is not set, use a placeholder
		if loc.Metadata.LocationName == "" {
			log.Printf("WARNING: Location name is empty in PDF generation! Setting a placeholder for location %d", i+1)
			locationResults[i].Metadata.LocationName = fmt.Sprintf("Location %d", i+1)
		}
	}

	// Create client report data structure
	clientReport := &shared.ClientReportData{
		ReportID:        reportID,
		ClientID:        clientInfo.ClientID,
		ClientName:      clientInfo.ClientName,
		PeriodStart:     periodStart,
		PeriodEnd:       periodEnd,
		GeneratedAt:     time.Now(),
		LocationResults: locationResults,
	}

	// Create and parse template
	tmpl, err := template.New("report").Parse(shared.MonthlyReportTemplate)
	if err != nil {
		return nil, fmt.Errorf("error parsing template: %v", err)
	}

	// Execute template to buffer
	var htmlBuffer bytes.Buffer
	err = tmpl.Execute(&htmlBuffer, clientReport)
	if err != nil {
		return nil, fmt.Errorf("error rendering template: %v", err)
	}

	// Convert HTML to PDF
	pdfContent, err := ConvertHTMLToPDFFunc(htmlBuffer.String())
	if err != nil {
		return nil, fmt.Errorf("error converting to PDF: %v", err)
	}

	return pdfContent, nil
}

// prepareReviewBatch converts raw review data from the API to the format expected by the analyzer
func prepareReviewBatch(reviews []map[string]interface{},
	location database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode,
	clientInfo database.ClientWithMonthlyReviewAnalysis,
	startDate time.Time, endDate time.Time) review.ReviewBatch {

	// Create the batch with location information
	batch := review.ReviewBatch{
		LocationName: location.GoogleMyBusinessLocationName,
		LocationID:   location.GoogleMyBusinessLocationPath,
		BusinessName: clientInfo.ClientName,
		ClientID:     clientInfo.ClientID,
		PostalCode:   location.GoogleMyBusinessPostalCode,
	}

	// Validate the batch has a location name
	if batch.LocationName == "" {
		log.Printf("WARNING: Location name is empty after creating batch! Setting a fallback name.")
		// Extract a name from the location ID as fallback
		pathParts := strings.Split(batch.LocationID, "/")
		if len(pathParts) > 0 {
			batch.LocationName = fmt.Sprintf("Location %s", pathParts[len(pathParts)-1])
		} else {
			batch.LocationName = "Unknown Location"
		}
	}

	// Set the report period
	batch.ReportPeriod.StartDate = startDate.Format("2006-01-02")
	batch.ReportPeriod.EndDate = endDate.AddDate(0, 0, -1).Format("2006-01-02") // Adjust end date to last day of month

	// Tracking stats for diagnostics
	totalReviews := len(reviews)
	reviewsWithComments := 0
	reviewsWithText := 0

	// Convert each review to the Review type
	for _, rawReview := range reviews {
		// Extract review text - but make it optional
		reviewText := ""
		comment, hasComment := rawReview["comment"].(map[string]interface{})
		if hasComment {
			reviewsWithComments++
			if text, hasText := comment["text"].(string); hasText {
				reviewText = text
				reviewsWithText++
			}
		} else {
			// Try extracting text directly if comment is not a map
			if commentStr, ok := rawReview["comment"].(string); ok {
				reviewText = commentStr
				reviewsWithComments++
				reviewsWithText++
			}
		}

		// Extract star rating
		starRating := 0
		ratingObj, hasRating := rawReview["starRating"].(map[string]interface{})
		if hasRating {
			ratingVal, hasCount := ratingObj["count"].(float64)
			if hasCount {
				starRating = int(ratingVal)
			}
		}

		// If the above approach didn't work, try the string-based rating format
		if starRating == 0 {
			if strRating, ok := rawReview["starRating"].(string); ok {
				switch strRating {
				case "ONE":
					starRating = 1
				case "TWO":
					starRating = 2
				case "THREE":
					starRating = 3
				case "FOUR":
					starRating = 4
				case "FIVE":
					starRating = 5
				}
			}
		}

		// Third approach - try to extract rating from a nested object with "name" field
		if starRating == 0 && hasRating {
			if nameField, ok := ratingObj["name"].(string); ok {
				// The name field might contain values like "ONE", "TWO", etc. at the end
				nameParts := strings.Split(nameField, "/")
				if len(nameParts) > 0 {
					lastPart := nameParts[len(nameParts)-1]
					switch lastPart {
					case "ONE":
						starRating = 1
					case "TWO":
						starRating = 2
					case "THREE":
						starRating = 3
					case "FOUR":
						starRating = 4
					case "FIVE":
						starRating = 5
					}
				}
			}
		}

		// Fourth approach - handle case where "starRating" is a map with a different schema
		if starRating == 0 && hasRating {
			// Check if there's a direct numerical value
			if val, ok := ratingObj["value"].(float64); ok {
				starRating = int(val)
			}

			// Check if there's an enum value
			if enum, ok := ratingObj["enum"].(string); ok {
				switch enum {
				case "ONE", "STAR_RATING_ONE":
					starRating = 1
				case "TWO", "STAR_RATING_TWO":
					starRating = 2
				case "THREE", "STAR_RATING_THREE":
					starRating = 3
				case "FOUR", "STAR_RATING_FOUR":
					starRating = 4
				case "FIVE", "STAR_RATING_FIVE":
					starRating = 5
				}
			}
		}

		// Extract review date
		createTime := time.Time{}
		createTimeStr, hasCreateTime := rawReview["createTime"].(string)
		if hasCreateTime {
			parsedTime, err := time.Parse(time.RFC3339, createTimeStr)
			if err == nil {
				createTime = parsedTime
			}
		}

		// Extract review ID
		reviewID := fmt.Sprintf("%v", rawReview["name"])

		// Add to batch - include all reviews regardless of whether they have text
		batch.Reviews = append(batch.Reviews, review.Review{
			ID:     reviewID,
			Text:   reviewText,
			Rating: starRating,
			Date:   createTime,
		})
	}

	log.Printf("Review batch summary for %s: total=%d, with_comments=%d, with_text=%d, final_batch_size=%d",
		location.GoogleMyBusinessLocationName, totalReviews, reviewsWithComments, reviewsWithText, len(batch.Reviews))

	return batch
}

// sendReportEmail sends the monthly report PDF to the client via email
func sendReportEmail(emailSvc email_service.EmailService, clientInfo database.ClientWithMonthlyReviewAnalysis,
	periodStart time.Time, pdfContent []byte) error {

	if len(pdfContent) == 0 {
		return fmt.Errorf("cannot send email: PDF content is empty")
	}

	// Check if report_email_address is set (not null and not empty)
	var recipientEmail string
	if clientInfo.ReportEmailAddress.Valid && clientInfo.ReportEmailAddress.String != "" {
		recipientEmail = clientInfo.ReportEmailAddress.String
	} else {
		// We don't send the email if report_email_address is not set
		log.Printf("Skipping report email for client %s (ID: %d) - report_email_address is not set",
			clientInfo.ClientName, clientInfo.ClientID)
		return fmt.Errorf("client report_email_address is not set")
	}

	// Format the month for display in the email (e.g., "March 2025")
	monthStr := periodStart.Format("January 2006")

	// Send the email with the PDF attachment
	err := emailSvc.SendMonthlyReport(
		clientInfo.ClientName,
		recipientEmail,
		monthStr,
		pdfContent,
	)

	if err != nil {
		return fmt.Errorf("failed to send email to %s: %w", recipientEmail, err)
	}

	log.Printf("Monthly report email for %s sent successfully to %s",
		monthStr, recipientEmail)

	return nil
}
