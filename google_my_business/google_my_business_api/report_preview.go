package google_my_business_api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"google_my_business/database"
	"shared-templates"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// Client represents a simple client for the dropdown
type Client struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// StartReportPreviewServer starts a simple HTTP server for previewing report templates
func StartReportPreviewServer(port int, db *sql.DB) {
	// Register handlers
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleClientSelection(w, r, db)
	})
	http.HandleFunc("/preview", func(w http.ResponseWriter, r *http.Request) {
		handleReportPreview(w, r, db)
	})
	http.HandleFunc("/download-html", func(w http.ResponseWriter, r *http.Request) {
		handleDownloadHTML(w, r, db)
	})
	http.HandleFunc("/generate-pdf", func(w http.ResponseWriter, r *http.Request) {
		handleGeneratePDF(w, r, db)
	})
	http.HandleFunc("/save-to-database", func(w http.ResponseWriter, r *http.Request) {
		handleSaveToDatabase(w, r, db)
	})

	// Start the server
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Starting template preview server on http://localhost%s\n", addr)
	fmt.Printf("Visit http://localhost%s to select a client and generate reports\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// getAllClients fetches all clients from the database
func getAllClients(db *sql.DB) ([]Client, error) {
	query := `
		SELECT DISTINCT c.id, c.name 
		FROM clients c 
		WHERE c.enabled = 1 
		ORDER BY c.name
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query clients: %w", err)
	}
	defer rows.Close()

	var clients []Client
	for rows.Next() {
		var client Client
		if err := rows.Scan(&client.ID, &client.Name); err != nil {
			return nil, fmt.Errorf("failed to scan client: %w", err)
		}
		clients = append(clients, client)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating clients: %w", err)
	}

	return clients, nil
}

// handleClientSelection shows the client selection form with loaded clients
func handleClientSelection(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == "POST" {
		// Handle form submission - redirect to preview with parameters
		clientID := r.FormValue("client_id")
		month := r.FormValue("month")

		if clientID == "" || month == "" {
			http.Error(w, "Client ID and month are required", http.StatusBadRequest)
			return
		}

		// Redirect to preview with parameters
		http.Redirect(w, r, fmt.Sprintf("/preview?client_id=%s&month=%s", clientID, month), http.StatusSeeOther)
		return
	}

	// Fetch clients from database
	clients, err := getAllClients(db)
	if err != nil {
		log.Printf("Error fetching clients: %v", err)
		http.Error(w, "Error loading clients", http.StatusInternalServerError)
		return
	}

	// Generate client options HTML
	clientOptions := ""
	for _, client := range clients {
		clientOptions += fmt.Sprintf(`<option value="%d">%s (ID: %d)</option>`, client.ID, client.Name, client.ID)
	}

	// Show the client selection form
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Report Preview - Client Selection</title>
			<style>
				body { 
					font-family: Arial, sans-serif; 
					margin: 40px auto; 
					max-width: 600px; 
					line-height: 1.6; 
					background-color: #f5f5f5;
				}
				.container { 
					background: white; 
					padding: 30px; 
					border-radius: 8px; 
					box-shadow: 0 2px 10px rgba(0,0,0,0.1);
				}
				h1 { 
					color: #333; 
					text-align: center; 
					margin-bottom: 30px;
				}
				.form-group { 
					margin-bottom: 20px; 
				}
				label { 
					display: block; 
					margin-bottom: 5px; 
					font-weight: bold; 
					color: #555;
				}
				select, input[type="text"] { 
					width: 100%%; 
					padding: 10px; 
					border: 1px solid #ddd; 
					border-radius: 4px; 
					font-size: 16px;
					box-sizing: border-box;
				}
				button { 
					width: 100%%; 
					padding: 12px; 
					background-color: #0066cc; 
					color: white; 
					border: none; 
					border-radius: 4px; 
					font-size: 16px; 
					font-weight: bold; 
					cursor: pointer;
				}
				button:hover { 
					background-color: #0052a3; 
				}
				.help-text {
					font-size: 14px;
					color: #666;
					margin-top: 5px;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h1>Monthly Report Preview</h1>
				<p>Select a client and month to generate a preview report with dummy data.</p>
				
				<form method="POST">
					<div class="form-group">
						<label for="client_id">Client:</label>
						<select name="client_id" id="client_id" required>
							<option value="">Select a client...</option>
							%s
						</select>
					</div>
					
					<div class="form-group">
						<label for="month">Month (YYYY-MM):</label>
						<input type="text" name="month" id="month" placeholder="2024-01" pattern="[0-9]{4}-[0-9]{2}" required>
						<div class="help-text">Enter the month in YYYY-MM format (e.g., 2024-01 for January 2024)</div>
					</div>
					
					<button type="submit">Generate Report Preview</button>
				</form>
			</div>

			<script>
				// Auto-populate current month - 1
				document.addEventListener('DOMContentLoaded', function() {
					const monthInput = document.getElementById('month');
					const now = new Date();
					const lastMonth = new Date(now.getFullYear(), now.getMonth() - 1, 1);
					const year = lastMonth.getFullYear();
					const month = String(lastMonth.getMonth() + 1).padStart(2, '0');
					monthInput.value = year + '-' + month;
				});
			</script>
		</body>
		</html>
	`, clientOptions)
}

// handleReportPreview renders the monthly report template with dummy data
func handleReportPreview(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Get parameters from URL
	clientIDStr := r.URL.Query().Get("client_id")
	monthStr := r.URL.Query().Get("month")

	if clientIDStr == "" || monthStr == "" {
		http.Error(w, "Missing client_id or month parameter", http.StatusBadRequest)
		return
	}

	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		http.Error(w, "Invalid client_id", http.StatusBadRequest)
		return
	}

	// Parse month
	var targetMonth time.Time
	targetMonth, err = time.Parse("2006-01", monthStr)
	if err != nil {
		http.Error(w, "Invalid month format. Use YYYY-MM", http.StatusBadRequest)
		return
	}

	// Generate dummy data for the specific client and month
	dummyData, err := generateDummyClientReportForClient(db, clientID, targetMonth)
	if err != nil {
		log.Printf("Error generating dummy data: %v", err)
		http.Error(w, "Error generating report data", http.StatusInternalServerError)
		return
	}

	// Add a navigation header with download buttons
	w.Write([]byte(fmt.Sprintf(`
		<div style="position: sticky; top: 0; background-color: #f8f9fa; padding: 10px; text-align: center; margin-bottom: 20px; z-index: 100; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
			<a href="/" style="display: inline-block; padding: 8px 16px; background-color: #6c757d; color: white; text-decoration: none; border-radius: 4px; font-weight: bold; margin-right: 10px;">← Back to Client Selection</a>
			<a href="/download-html?client_id=%s&month=%s" style="display: inline-block; padding: 8px 16px; background-color: #0066cc; color: white; text-decoration: none; border-radius: 4px; font-weight: bold; margin-right: 10px;">Download HTML</a>
			<a href="/generate-pdf?client_id=%s&month=%s" style="display: inline-block; padding: 8px 16px; background-color: #28a745; color: white; text-decoration: none; border-radius: 4px; font-weight: bold; margin-right: 10px;">Generate PDF</a>
			<button onclick="saveToDatabase('%s', '%s')" style="display: inline-block; padding: 8px 16px; background-color: #dc3545; color: white; border: none; border-radius: 4px; font-weight: bold; cursor: pointer;">Save to Database</button>
		</div>
		
		<script>
		function saveToDatabase(clientId, month) {
			if (confirm('Save this report data to the database for client ID ' + clientId + ' and month ' + month + '? This will overwrite any existing report for this period.')) {
				fetch('/save-to-database?client_id=' + clientId + '&month=' + month, { method: 'POST' })
					.then(response => response.json())
					.then(data => {
						if (data.success) {
							alert('Report saved successfully! Report ID: ' + data.report_id);
						} else {
							alert('Error saving report: ' + data.error);
						}
					})
					.catch(error => {
						alert('Error saving report: ' + error);
					});
			}
		}
		</script>
	`, clientIDStr, monthStr, clientIDStr, monthStr, clientIDStr, monthStr)))

	// Create and parse template
	tmpl, err := template.New("report").Parse(shared_templates.MonthlyReportTemplate)

	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	// Render template
	err = tmpl.Execute(w, dummyData)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// handleDownloadHTML generates the HTML content for download
func handleDownloadHTML(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Get parameters from URL
	clientIDStr := r.URL.Query().Get("client_id")
	monthStr := r.URL.Query().Get("month")

	if clientIDStr == "" || monthStr == "" {
		http.Error(w, "Missing client_id or month parameter", http.StatusBadRequest)
		return
	}

	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		http.Error(w, "Invalid client_id", http.StatusBadRequest)
		return
	}

	// Parse month
	var targetMonth time.Time
	targetMonth, err = time.Parse("2006-01", monthStr)
	if err != nil {
		http.Error(w, "Invalid month format. Use YYYY-MM", http.StatusBadRequest)
		return
	}

	// Generate dummy data for the specific client and month
	dummyData, err := generateDummyClientReportForClient(db, clientID, targetMonth)
	if err != nil {
		log.Printf("Error generating dummy data: %v", err)
		http.Error(w, "Error generating report data", http.StatusInternalServerError)
		return
	}

	// Create and parse template
	tmpl, err := template.New("report").Parse(shared_templates.MonthlyReportTemplate)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	// Set headers for download
	filename := fmt.Sprintf("report_%s_%s.html", dummyData.ClientName, monthStr)
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	// Render template directly to response writer
	err = tmpl.Execute(w, dummyData)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// handleGeneratePDF generates a PDF from the HTML template using the external API
func handleGeneratePDF(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Get parameters from URL
	clientIDStr := r.URL.Query().Get("client_id")
	monthStr := r.URL.Query().Get("month")

	if clientIDStr == "" || monthStr == "" {
		http.Error(w, "Missing client_id or month parameter", http.StatusBadRequest)
		return
	}

	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		http.Error(w, "Invalid client_id", http.StatusBadRequest)
		return
	}

	// Parse month
	var targetMonth time.Time
	targetMonth, err = time.Parse("2006-01", monthStr)
	if err != nil {
		http.Error(w, "Invalid month format. Use YYYY-MM", http.StatusBadRequest)
		return
	}

	// Generate dummy data for the specific client and month
	dummyData, err := generateDummyClientReportForClient(db, clientID, targetMonth)
	if err != nil {
		log.Printf("Error generating dummy data: %v", err)
		http.Error(w, "Error generating report data", http.StatusInternalServerError)
		return
	}

	// Generate HTML content
	var htmlBuffer bytes.Buffer

	// Create and parse template
	tmpl, err := template.New("report").Parse(shared_templates.MonthlyReportTemplate)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	// Execute template to buffer
	err = tmpl.Execute(&htmlBuffer, dummyData)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

	// Use the PDF service to convert HTML to PDF
	pdfData, err := ConvertHTMLToPDF(htmlBuffer.String())
	if err != nil {
		log.Printf("Error generating PDF: %v", err)
		http.Error(w, "Error generating PDF", http.StatusInternalServerError)
		return
	}

	// Set response headers for PDF download
	filename := fmt.Sprintf("report_%s_%s.pdf", dummyData.ClientName, monthStr)
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pdfData)))

	// Write PDF data to response
	_, err = w.Write(pdfData)
	if err != nil {
		log.Printf("Error writing PDF to response: %v", err)
		return
	}
}

// handleSaveToDatabase saves the dummy report data to the database
func handleSaveToDatabase(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get parameters from URL
	clientIDStr := r.URL.Query().Get("client_id")
	monthStr := r.URL.Query().Get("month")

	if clientIDStr == "" || monthStr == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Missing client_id or month parameter",
		})
		return
	}

	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid client_id",
		})
		return
	}

	// Parse month
	var targetMonth time.Time
	targetMonth, err = time.Parse("2006-01", monthStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid month format. Use YYYY-MM",
		})
		return
	}

	// Calculate period start and end
	periodStart := time.Date(targetMonth.Year(), targetMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
	periodEnd := periodStart.AddDate(0, 1, 0)

	// Check if a report already exists for this client and period (to overwrite)
	existingReport, err := database.GetClientReportByClientAndPeriod(db, clientID, periodStart, periodEnd)
	if err != nil && err.Error() != "failed to retrieve client report: sql: no rows in result set" {
		log.Printf("Error checking for existing report: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Database error checking existing report",
		})
		return
	}

	// If report exists, delete it (overwrite functionality)
	if existingReport != nil {
		log.Printf("Deleting existing report ID %d for client %d and period %s", existingReport.ReportID, clientID, monthStr)
		_, err = db.Exec("DELETE FROM client_reports WHERE report_id = ?", existingReport.ReportID)
		if err != nil {
			log.Printf("Error deleting existing report: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Error deleting existing report",
			})
			return
		}
	}

	// Generate dummy data for the specific client and month
	dummyData, err := generateDummyClientReportForClient(db, clientID, targetMonth)
	if err != nil {
		log.Printf("Error generating dummy data: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Error generating report data",
		})
		return
	}

	// Convert location results to JSON for database storage
	locationResultsJSON, err := json.Marshal(dummyData.LocationResults)
	if err != nil {
		log.Printf("Error marshaling location results: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Error processing report data",
		})
		return
	}

	// Save the report to database using the existing function
	reportID, err := database.SaveClientReport(db, clientID, periodStart, periodEnd, locationResultsJSON)
	if err != nil {
		log.Printf("Error saving client report: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Error saving report to database",
		})
		return
	}

	log.Printf("Successfully saved report with ID %d for client %d and period %s", reportID, clientID, monthStr)

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"report_id": reportID,
		"message":   fmt.Sprintf("Report saved successfully for %s (ID: %d) for period %s", dummyData.ClientName, clientID, monthStr),
	})
}

// getClientByID fetches a client by ID from database
func getClientByID(db *sql.DB, clientID int) (*Client, error) {
	query := "SELECT id, name FROM clients WHERE id = ? AND enabled = 1"

	var client Client
	err := db.QueryRow(query, clientID).Scan(&client.ID, &client.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	return &client, nil
}

// generateDummyClientReportForClient creates sample data for a specific client and month
func generateDummyClientReportForClient(db *sql.DB, clientID int, targetMonth time.Time) (*shared_templates.ClientReportData, error) {
	// Get client info from database
	client, err := getClientByID(db, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get client info: %w", err)
	}

	// Seed random number generator with client ID and month for consistency
	// Use a more complex seed to avoid edge cases
	seed := int64(clientID)*1000 + targetMonth.Unix()/86400 // Use days since epoch
	log.Printf("DEBUG: Using random seed %d for client %d and month %s", seed, clientID, targetMonth.Format("2006-01"))
	rand.Seed(seed)

	// Calculate period start and end
	periodStart := time.Date(targetMonth.Year(), targetMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
	periodEnd := periodStart.AddDate(0, 1, 0)

	// Create dummy locations (2-4)
	locationCount := 2 + rand.Intn(3)
	locationResults := make([]shared_templates.AnalysisResult, locationCount)

	locationNames := []string{
		fmt.Sprintf("%s - Main Location", client.Name),
		fmt.Sprintf("%s - Branch Office", client.Name),
		fmt.Sprintf("%s - Service Center", client.Name),
		fmt.Sprintf("%s - Express Station", client.Name),
	}

	// Generate data for each location
	for i := 0; i < locationCount; i++ {
		locationName := locationNames[i]

		// Create a realistic number of reviews (10-40)
		reviewCount := 10 + rand.Intn(31)
		log.Printf("DEBUG: Location %d (%s) - Generated reviewCount: %d", i, locationName, reviewCount)

		// Calculate sentiment distribution
		positiveCount := int(float64(reviewCount) * (0.6 + rand.Float64()*0.2))
		neutralCount := int(float64(reviewCount-positiveCount) * (0.3 + rand.Float64()*0.4))
		negativeCount := reviewCount - positiveCount - neutralCount
		if negativeCount < 0 {
			negativeCount = 0
		}

		log.Printf("DEBUG: Location %d - Sentiment counts: Positive=%d, Neutral=%d, Negative=%d, Total=%d",
			i, positiveCount, neutralCount, negativeCount, positiveCount+neutralCount+negativeCount)

		// Calculate percentages
		positivePercentage := float64(positiveCount) / float64(reviewCount) * 100
		neutralPercentage := float64(neutralCount) / float64(reviewCount) * 100
		negativePercentage := float64(negativeCount) / float64(reviewCount) * 100

		// Create location-specific analysis
		positiveThemes := []string{"Professional service", "Friendly staff", "Clean environment"}
		negativeThemes := []string{"Wait times", "Parking issues"}

		strengths := []shared_templates.Insight{
			{
				Category:    "Service Quality",
				Description: "Customers consistently praise the professional and efficient service delivery.",
				Example:     "Excellent service from start to finish. Very professional team.",
			},
			{
				Category:    "Staff Friendliness",
				Description: "Multiple reviews highlight the friendly and helpful nature of the staff.",
				Example:     "Staff were incredibly friendly and went out of their way to help us.",
			},
		}

		improvements := []shared_templates.Insight{
			{
				Category:    "Wait Times",
				Description: "Some customers mention longer than expected wait times during busy periods.",
				Example:     "Service was great but had to wait about 25 minutes longer than scheduled.",
			},
		}

		// Create negative review categories
		negativeCategories := []shared_templates.ReviewCategory{}
		if negativeCount > 0 {
			negativeCategories = []shared_templates.ReviewCategory{
				{
					Name:       "Wait Time",
					Count:      negativeCount / 2,
					Percentage: 50.0,
				},
				{
					Name:       "Service Issues",
					Count:      negativeCount - (negativeCount / 2),
					Percentage: 50.0,
				},
			}
		}

		recommendations := []string{
			"Implement appointment scheduling system to reduce wait times",
			"Enhance staff training on customer communication during delays",
			"Consider expanding service hours during peak demand periods",
		}

		operatorTraining := []string{
			"Customer communication best practices for managing expectations",
			"Efficient workflow management during high-volume periods",
			"Problem resolution techniques for common customer concerns",
		}

		driverTraining := []string{
			"Professional customer interaction protocols",
			"Time management strategies for punctual service delivery",
			"Effective communication during service delays or changes",
		}

		// Create the analysis result
		locationResults[i] = shared_templates.AnalysisResult{
			Analysis: shared_templates.Analysis{
				OverallSummary: shared_templates.OverallSummary{
					SummaryText:       fmt.Sprintf("Analysis of %s for %s shows positive customer sentiment with opportunities for improvement in operational efficiency.", locationName, targetMonth.Format("January 2006")),
					PositiveThemes:    positiveThemes,
					NegativeThemes:    negativeThemes,
					OverallPerception: "Customers generally view the business positively with high satisfaction rates.",
					AverageRating:     3.8 + rand.Float64()*1.0,
				},
				SentimentAnalysis: shared_templates.SentimentAnalysis{
					PositiveCount:      positiveCount,
					PositivePercentage: positivePercentage,
					NeutralCount:       neutralCount,
					NeutralPercentage:  neutralPercentage,
					NegativeCount:      negativeCount,
					NegativePercentage: negativePercentage,
					TotalReviews:       reviewCount,
					SentimentTrend:     "stable",
				},
				KeyTakeaways: shared_templates.KeyTakeaways{
					Strengths:           strengths,
					AreasForImprovement: improvements,
				},
				NegativeReviewBreakdown: shared_templates.NegativeReviewBreakdown{
					Categories:                 negativeCategories,
					ImprovementRecommendations: recommendations,
				},
				TrainingRecommendations: shared_templates.TrainingRecommendations{
					ForOperators: operatorTraining,
					ForDrivers:   driverTraining,
				},
			},
			Metadata: shared_templates.AnalysisMetadata{
				GeneratedAt:   time.Now(),
				ReviewCount:   reviewCount,
				LocationID:    fmt.Sprintf("accounts/123/locations/%d", 100+i),
				LocationName:  locationName,
				BusinessName:  client.Name,
				AnalyzerID:    "gpt-4",
				AnalyzerName:  "AI Review Analyzer",
				AnalyzerModel: "gpt-4-0613",
				ReportPeriod: struct {
					StartDate string `json:"start_date"`
					EndDate   string `json:"end_date"`
				}{
					StartDate: periodStart.Format("2006-01-02"),
					EndDate:   periodEnd.AddDate(0, 0, -1).Format("2006-01-02"),
				},
			},
		}
	}

	// Create the complete client report
	return &shared_templates.ClientReportData{
		ReportID:        int64(rand.Intn(10000) + 1000),
		ClientID:        clientID,
		ClientName:      client.Name,
		PeriodStart:     periodStart,
		PeriodEnd:       periodEnd,
		GeneratedAt:     time.Now(),
		LocationResults: locationResults,
	}, nil
}

// generateDummyClientReport creates sample data for the client report preview (legacy function for compatibility)
func generateDummyClientReport() *shared_templates.ClientReportData {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate time period (last month)
	now := time.Now()
	lastMonth := now.AddDate(0, -1, 0)
	periodStart := time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	// Create dummy locations (3-5)
	locationCount := 3 + rand.Intn(3)
	locationResults := make([]shared_templates.AnalysisResult, locationCount)

	locationNames := []string{
		"Downtown Office",
		"Airport Terminal",
		"Shopping Center",
		"Hotel Branch",
		"Mall Kiosk",
		"Business District",
	}

	// Generate data for each location
	for i := 0; i < locationCount; i++ {
		locationName := locationNames[i]

		// Create a realistic number of reviews (15-45)
		reviewCount := 15 + rand.Intn(31)

		// Define location-specific negative issues for consistency
		var negativeThemes []string
		var locationImprovements []shared_templates.Insight
		var negativeCategories []string

		// Generate consistent negative themes based on location
		if i == 0 { // Downtown Office
			negativeThemes = []string{"Wait times", "Parking issues"}
			locationImprovements = []shared_templates.Insight{
				{
					Category:    "Wait Times",
					Description: "Some customers mention long wait times during peak hours.",
					Example:     "Had to wait 20 minutes to be served on a Tuesday afternoon.",
				},
				{
					Category:    "Parking Availability",
					Description: "Limited parking options cause frustration for customers.",
					Example:     "Spent 15 minutes looking for parking before my appointment.",
				},
			}
			negativeCategories = []string{"Wait Time", "Parking Problems", "Staff Issues"}
		} else if i == 1 { // Airport Terminal
			negativeThemes = []string{"Finding location", "Communication issues"}
			locationImprovements = []shared_templates.Insight{
				{
					Category:    "Signage Issues",
					Description: "Several customers reported difficulty finding the location due to poor signage.",
					Example:     "Walked around for 10 minutes trying to find your office. Need better signs.",
				},
				{
					Category:    "Communication Problems",
					Description: "Customers report unclear information about location and procedures.",
					Example:     "Instructions were confusing and staff didn't explain the process clearly.",
				},
			}
			negativeCategories = []string{"Location Problems", "Communication Issues", "Service Quality"}
		} else { // Shopping Center, Hotel Branch, etc.
			negativeThemes = []string{"Service speed", "Staff availability"}
			locationImprovements = []shared_templates.Insight{
				{
					Category:    "Service Speed",
					Description: "Customers mention slower than expected service delivery.",
					Example:     "Service took much longer than promised, affecting my schedule.",
				},
				{
					Category:    "Staff Availability",
					Description: "Limited staff during busy periods causes delays.",
					Example:     "Only one person working during lunch rush, caused long delays.",
				},
			}
			negativeCategories = []string{"Service Quality", "Staff Issues", "Wait Time"}
		}

		// Create consistent overall summary
		positiveThemes := []string{"Helpful staff", "Professional service", "Clean environment"}

		// Calculate sentiment distribution
		positiveCount := int(float64(reviewCount) * (0.5 + rand.Float64()*0.3))
		neutralCount := int(float64(reviewCount-positiveCount) * (0.4 + rand.Float64()*0.3))
		negativeCount := reviewCount - positiveCount - neutralCount
		if negativeCount < 0 {
			negativeCount = 0
		}

		// Calculate percentages
		positivePercentage := float64(positiveCount) / float64(reviewCount) * 100
		neutralPercentage := float64(neutralCount) / float64(reviewCount) * 100
		negativePercentage := float64(negativeCount) / float64(reviewCount) * 100

		// Create strengths (2-3)
		strengths := []shared_templates.Insight{
			{
				Category:    "Staff Friendliness",
				Description: "Customers frequently praise the friendly and helpful staff.",
				Example:     "Every employee I've interacted with has been incredibly friendly and accommodating.",
			},
		}

		// Add location-specific strength
		if i == 0 {
			strengths = append(strengths, shared_templates.Insight{
				Category:    "Central Location",
				Description: "Many clients appreciate the convenient downtown location.",
				Example:     "Located right in the center of town, very easy to find and get to.",
			})
		} else if i == 1 {
			strengths = append(strengths, shared_templates.Insight{
				Category:    "Transportation Access",
				Description: "Customers value the easy access to transportation options.",
				Example:     "Right next to the train station, making it very convenient for travelers.",
			})
		} else {
			strengths = append(strengths, shared_templates.Insight{
				Category:    "Modern Facilities",
				Description: "Clients appreciate the clean and up-to-date facilities.",
				Example:     "Everything is new and well-maintained, making for a pleasant experience.",
			})
		}

		// Create improvement areas (1-2)
		improvements := []shared_templates.Insight{
			{
				Category:    "Wait Times",
				Description: "Some customers mention long wait times during peak hours.",
				Example:     "Had to wait 20 minutes to be served on a Tuesday afternoon.",
			},
		}

		// Add location-specific improvement area
		if i == 1 {
			improvements = append(improvements, shared_templates.Insight{
				Category:    "Signage Issues",
				Description: "Several customers reported difficulty finding the location due to poor signage.",
				Example:     "Walked around for 10 minutes trying to find your office. Need better signs.",
			})
		} else if i == 2 {
			improvements = append(improvements, shared_templates.Insight{
				Category:    "Parking Concerns",
				Description: "Customers mention limited parking options near the location.",
				Example:     "Great service but the parking situation is terrible. Had to park three blocks away.",
			})
		}

		// Create location-specific recommendations that match the issues
		var recommendations []string
		if i == 0 { // Downtown Office - Wait times & Parking
			recommendations = []string{
				"Implement a digital queue management system to reduce perceived wait times",
				"Explore partnerships with nearby parking facilities to offer customer validation",
				"Consider extending hours to spread customer traffic more evenly throughout the day",
			}
		} else if i == 1 { // Airport Terminal - Finding location & Communication
			recommendations = []string{
				"Install clearer directional signage and digital displays throughout the terminal",
				"Provide detailed location instructions in confirmation emails and messages",
				"Train staff on clear communication protocols for explaining procedures to customers",
			}
		} else { // Other locations - Service speed & Staff availability
			recommendations = []string{
				"Implement process improvements to reduce service delivery times",
				"Increase staffing during historically busy periods based on traffic analysis",
				"Cross-train team members to provide backup support during peak hours",
			}
		}

		// Create location-specific training recommendations that match the issues
		var operatorTraining, driverTraining []string
		if i == 0 { // Downtown Office - Wait times & Parking
			operatorTraining = []string{
				"Queue management techniques to optimize customer flow during peak hours",
				"Proactive communication strategies for managing customer expectations during delays",
				"Alternative parking guidance and validation partnership protocols",
			}
			driverTraining = []string{
				"Efficient route planning to minimize customer travel time to location",
				"Customer notification procedures for pickup/delivery timing adjustments",
				"Professional handling of parking-related customer concerns",
			}
		} else if i == 1 { // Airport Terminal - Finding location & Communication
			operatorTraining = []string{
				"Clear directional assistance and location guidance for confused customers",
				"Multi-language communication basics for diverse airport clientele",
				"Terminal navigation support and wayfinding assistance techniques",
			}
			driverTraining = []string{
				"Airport terminal pickup/dropoff procedures and designated zones",
				"Customer contact protocols for location coordination",
				"Professional communication during complex airport logistics",
			}
		} else { // West Side Branch - Service speed & Staff availability
			operatorTraining = []string{
				"Workflow optimization techniques to improve service delivery speed",
				"Cross-training protocols for better staff coverage during busy periods",
				"Time management strategies for handling multiple customer requests efficiently",
			}
			driverTraining = []string{
				"Expedited service delivery methods for time-sensitive customers",
				"Professional customer interaction during rush periods",
				"Efficient vehicle preparation and maintenance for faster turnaround",
			}
		}

		// Create the analysis result
		locationResults[i] = shared_templates.AnalysisResult{
			Analysis: shared_templates.Analysis{
				OverallSummary: shared_templates.OverallSummary{
					SummaryText:       fmt.Sprintf("Reviews for %s show generally positive customer experiences with some areas for improvement. Customers appreciate the staff friendliness but note concerns about wait times during busy periods.", locationName),
					PositiveThemes:    positiveThemes,
					NegativeThemes:    negativeThemes,
					OverallPerception: "Customers generally view the business positively with most expressing satisfaction.",
					AverageRating:     4.0 + rand.Float64(),
				},
				SentimentAnalysis: shared_templates.SentimentAnalysis{
					PositiveCount:      positiveCount,
					PositivePercentage: positivePercentage,
					NeutralCount:       neutralCount,
					NeutralPercentage:  neutralPercentage,
					NegativeCount:      negativeCount,
					NegativePercentage: negativePercentage,
					TotalReviews:       reviewCount,
					SentimentTrend:     "improving",
				},
				KeyTakeaways: shared_templates.KeyTakeaways{
					Strengths:           strengths,
					AreasForImprovement: locationImprovements,
				},
				NegativeReviewBreakdown: shared_templates.NegativeReviewBreakdown{
					Categories:                 generateCoherentNegativeCategories(negativeCount, negativeCategories),
					ImprovementRecommendations: recommendations,
				},
				TrainingRecommendations: shared_templates.TrainingRecommendations{
					ForOperators: operatorTraining,
					ForDrivers:   driverTraining,
				},
			},
			Metadata: shared_templates.AnalysisMetadata{
				GeneratedAt:   time.Now(),
				ReviewCount:   reviewCount,
				LocationID:    fmt.Sprintf("accounts/123/locations/%d", 100+i),
				LocationName:  locationName,
				BusinessName:  "ABC Business Services",
				AnalyzerID:    "gpt-4",
				AnalyzerName:  "AI Review Analyzer",
				AnalyzerModel: "gpt-4-0613",
			},
		}

		// Set report period
		locationResults[i].Metadata.ReportPeriod.StartDate = periodStart.Format("2006-01-02")
		locationResults[i].Metadata.ReportPeriod.EndDate = periodEnd.AddDate(0, 0, -1).Format("2006-01-02")
	}

	// Create the complete client report
	return &shared_templates.ClientReportData{
		ReportID:        123,
		ClientID:        456,
		ClientName:      "ABC Business Services",
		PeriodStart:     periodStart,
		PeriodEnd:       periodEnd,
		GeneratedAt:     time.Now(),
		LocationResults: locationResults,
	}
}

// Helper function to format dates
func formatDate(t time.Time) string {
	return t.Format("January 2, 2006")
}

// Helper function to format percentages
func formatPercentage(val float64) string {
	return fmt.Sprintf("%.1f%%", val)
}

// generateRandomNegativeCategories generates random negative categories for the NegativeReviewBreakdown
func generateRandomNegativeCategories(negativeCount int, locationIndex int) []shared_templates.ReviewCategory {
	if negativeCount == 0 {
		return []shared_templates.ReviewCategory{}
	}

	// Define possible negative categories with realistic names
	possibleCategories := []string{
		"Wait Time",
		"Staff Issues",
		"Location Problems",
		"Service Quality",
		"Cleanliness Concerns",
		"Communication Issues",
		"Pricing Complaints",
		"Technical Problems",
		"Accessibility Issues",
		"Parking Problems",
	}

	// Select 2-4 categories randomly
	numCategories := 2 + rand.Intn(3) // 2-4 categories
	if numCategories > negativeCount {
		numCategories = negativeCount
	}
	if numCategories > len(possibleCategories) {
		numCategories = len(possibleCategories)
	}

	// Shuffle and select categories
	selectedCategories := make([]string, numCategories)
	shuffled := make([]string, len(possibleCategories))
	copy(shuffled, possibleCategories)

	// Fisher-Yates shuffle
	for i := len(shuffled) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	for i := 0; i < numCategories; i++ {
		selectedCategories[i] = shuffled[i]
	}

	// Generate random distribution that adds up to 100%
	remainingCount := negativeCount

	categories := make([]shared_templates.ReviewCategory, numCategories)

	for i := 0; i < numCategories-1; i++ {
		// Random count between 1 and remaining count
		maxCount := remainingCount - (numCategories - i - 1)
		if maxCount < 1 {
			maxCount = 1
		}
		count := 1 + rand.Intn(maxCount)

		categories[i] = shared_templates.ReviewCategory{
			Name:       selectedCategories[i],
			Count:      count,
			Percentage: float64(count) / float64(negativeCount) * 100,
		}

		remainingCount -= count
	}

	// Last category gets remaining count
	categories[numCategories-1] = shared_templates.ReviewCategory{
		Name:       selectedCategories[numCategories-1],
		Count:      remainingCount,
		Percentage: float64(remainingCount) / float64(negativeCount) * 100,
	}

	return categories
}

// generateCoherentNegativeCategories generates negative categories for the NegativeReviewBreakdown
func generateCoherentNegativeCategories(negativeCount int, negativeCategories []string) []shared_templates.ReviewCategory {
	if negativeCount == 0 {
		return []shared_templates.ReviewCategory{}
	}

	// Define possible negative categories with realistic names
	possibleCategories := negativeCategories

	// Select 2-4 categories randomly
	numCategories := 2 + rand.Intn(3) // 2-4 categories
	if numCategories > negativeCount {
		numCategories = negativeCount
	}
	if numCategories > len(possibleCategories) {
		numCategories = len(possibleCategories)
	}

	// Shuffle and select categories
	selectedCategories := make([]string, numCategories)
	shuffled := make([]string, len(possibleCategories))
	copy(shuffled, possibleCategories)

	// Fisher-Yates shuffle
	for i := len(shuffled) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	for i := 0; i < numCategories; i++ {
		selectedCategories[i] = shuffled[i]
	}

	// Generate random distribution that adds up to 100%
	remainingCount := negativeCount

	categories := make([]shared_templates.ReviewCategory, numCategories)

	for i := 0; i < numCategories-1; i++ {
		// Random count between 1 and remaining count
		maxCount := remainingCount - (numCategories - i - 1)
		if maxCount < 1 {
			maxCount = 1
		}
		count := 1 + rand.Intn(maxCount)

		categories[i] = shared_templates.ReviewCategory{
			Name:       selectedCategories[i],
			Count:      count,
			Percentage: float64(count) / float64(negativeCount) * 100,
		}

		remainingCount -= count
	}

	// Last category gets remaining count
	categories[numCategories-1] = shared_templates.ReviewCategory{
		Name:       selectedCategories[numCategories-1],
		Count:      remainingCount,
		Percentage: float64(remainingCount) / float64(negativeCount) * 100,
	}

	return categories
}
