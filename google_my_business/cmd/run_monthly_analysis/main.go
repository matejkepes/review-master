package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"google_my_business/ai_service"
	"google_my_business/ai_service/review"
	"google_my_business/config"
	"google_my_business/database"
	"google_my_business/email_service"
	"google_my_business/google_my_business_api"
	"shared_templates"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// RetryTracking represents the structure for tracking client retry attempts
type RetryTracking map[int]map[string]ClientRetryInfo

// ClientRetryInfo tracks retry information for a specific client and month
type ClientRetryInfo struct {
	Month     string `json:"month"`
	Attempts  int    `json:"attempts"`
	LastError string `json:"last_error"`
}

// Constants
const (
	MaxRetryAttempts = 3
	RetryFilePath    = "retry_tracking.json"
)

// DBAdapter implements the google_my_business_api.DBInterface
type DBAdapter struct {
	db *sql.DB
}

// Client returns the underlying database client
func (a *DBAdapter) Client() *sql.DB {
	return a.db
}

// GetClientsWithMonthlyReviewAnalysisEnabled returns clients with monthly analysis enabled
func (a *DBAdapter) GetClientsWithMonthlyReviewAnalysisEnabled() ([]database.ClientWithMonthlyReviewAnalysis, error) {
	clients, err := database.GetClientsWithMonthlyReviewAnalysisEnabled(a.db)
	return clients, err
}

// GetClientReportByClientAndPeriod checks if a report exists for a time period
func (a *DBAdapter) GetClientReportByClientAndPeriod(clientID int, periodStart, periodEnd time.Time) (*shared_templates.ClientReportData, error) {
	report, err := database.GetClientReportByClientAndPeriod(a.db, clientID, periodStart, periodEnd)
	if err != nil {
		// Check if it's a "not found" error
		if err.Error() == "failed to retrieve client report: sql: no rows in result set" {
			return nil, nil // Return nil, nil to indicate no report found
		}
		// Otherwise, return the error
		return nil, err
	}
	return report, nil
}

// DeleteClientReport deletes an existing report
func (a *DBAdapter) DeleteClientReport(reportID int64) error {
	_, err := a.db.Exec("DELETE FROM client_reports WHERE report_id = ?", reportID)
	return err
}

// SaveClientReport saves a new report
func (a *DBAdapter) SaveClientReport(clientID int, periodStart, periodEnd time.Time, locationResults []byte) (int64, error) {
	return database.SaveClientReport(a.db, clientID, periodStart, periodEnd, locationResults)
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Creates an authenticated Google OAuth client
func getOAuthClient() (*http.Client, error) {
	// Read OAuth credentials file
	credFile := "credentials.json"
	b, err := ioutil.ReadFile(credFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %v", err)
	}

	// Parse credentials
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/business.manage")
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

	// Get token from file
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		return nil, fmt.Errorf("token file error: %v", err)
	}

	// Create client with automatic token refresh
	return config.Client(context.Background(), tok), nil
}

func main() {
	// Define command-line flags
	targetMonthStr := flag.String("month", "", "Target month in YYYY-MM format (defaults to previous month)")
	forceReprocess := flag.Bool("force-reprocess", false, "Force reprocessing of existing reports")
	retryOnly := flag.Bool("retry-only", false, "Only retry previously failed clients")
	emailSummary := flag.String("email-summary", "", "Email address to send processing summary")
	flag.Parse()

	// Process positional arguments (for backward compatibility)
	// If first argument is not a flag and looks like a date, use it as target month
	args := flag.Args()
	if len(args) > 0 && len(args[0]) == 7 && args[0][4] == '-' {
		*targetMonthStr = args[0]
	}

	// Load system configuration
	cfg := config.ReadProperties()

	// Initialize database connection
	dbConnectionString := cfg.DbUsername + ":" + cfg.DbPassword + "@tcp(" + cfg.DbAddress + ":" + cfg.DbPort + ")/" + cfg.DbName + "?parseTime=true"

	db, err := sql.Open("mysql", dbConnectionString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Create the database adapter
	dbAdapter := &DBAdapter{db: db}

	// Initialize Google OAuth HTTP client for Google My Business API
	httpClient, err := getOAuthClient()
	if err != nil {
		log.Fatalf("Failed to initialize OAuth client: %v", err)
	}

	// Initialize LLM provider for analysis
	llmProvider, err := ai_service.NewOpenAIProvider(ai_service.GetConfigForUseCase(ai_service.MonthlyAnalysis, cfg.OpenAIAPIKey))
	if err != nil {
		log.Fatalf("Failed to initialize AI provider: %v", err)
	}

	// Initialize review analyzer
	analyzer := review.NewAnalyzer(llmProvider, review.AnalyzerConfig{
		SystemPrompt: "You are a business analytics expert analyzing customer reviews. Provide insights in JSON format.",
		MaxTokens:    2000,
		ModelName:    "gpt-4",
		AnalyzerID:   "monthly-report-analyzer",
		AnalyzerName: "Monthly Report Analyzer",
	})

	// Initialize email service for sending reports
	emailSvc := email_service.NewSendGridEmailService(
		cfg.SendGridAPIKey,
		"admin@review-assistant.com",
		"Monthly Review Report", // From name
		cfg.SendGridTemplateID,
	)

	// Determine target month (default to previous month)
	now := time.Now()
	targetMonth := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, time.UTC)

	// Check if a specific month was provided
	if *targetMonthStr != "" {
		if tm, err := time.Parse("2006-01", *targetMonthStr); err == nil {
			targetMonth = tm
		} else {
			log.Printf("Warning: Couldn't parse month '%s', using previous month. Expected format: YYYY-MM", *targetMonthStr)
		}
	}

	targetMonthFormatted := targetMonth.Format("2006-01")
	fmt.Printf("Running monthly analysis for %s\n", targetMonthFormatted)

	// Load retry tracking if in retry mode or to update after processing
	retryTracking := loadRetryTracking()

	// If in retry-only mode, determine which clients to process
	var clientsToProcess []int
	if *retryOnly {
		// Extract clients needing retry for this month
		for clientID, retryInfo := range retryTracking {
			if monthInfo, exists := retryInfo[targetMonthFormatted]; exists {
				if monthInfo.Attempts < MaxRetryAttempts {
					clientsToProcess = append(clientsToProcess, clientID)
					fmt.Printf("Will retry client %d (previous attempts: %d)\n", clientID, monthInfo.Attempts)
				} else {
					fmt.Printf("Skipping client %d - max retries (%d) exceeded\n", clientID, MaxRetryAttempts)
				}
			}
		}

		if len(clientsToProcess) == 0 {
			fmt.Println("No clients to retry for this month")
			return
		}
	}

	// Run the analysis process
	summary, err := google_my_business_api.AnalyzeClientReviews(dbAdapter, httpClient, analyzer, &targetMonth, *forceReprocess, emailSvc)
	if err != nil {
		log.Fatalf("Analysis failed: %v", err)
	}

	// Update retry tracking
	if len(summary.FailedClients) > 0 {
		for _, fc := range summary.FailedClients {
			clientID := fc.ClientID

			// Initialize client in tracking if needed
			if _, exists := retryTracking[clientID]; !exists {
				retryTracking[clientID] = make(map[string]ClientRetryInfo)
			}

			// Get current retry info or initialize
			retryInfo, exists := retryTracking[clientID][targetMonthFormatted]
			if !exists {
				retryInfo = ClientRetryInfo{
					Month:    targetMonthFormatted,
					Attempts: 0,
				}
			}

			// Update retry info
			retryInfo.Attempts++
			retryInfo.LastError = fc.Error
			retryTracking[clientID][targetMonthFormatted] = retryInfo
		}

		// Save updated retry tracking
		saveRetryTracking(retryTracking)
	}

	// Print summary
	fmt.Printf("\nMonthly Analysis Complete\n")
	fmt.Printf("------------------------\n")
	fmt.Printf("Analyzed month: %s\n", summary.TargetMonth.Format("2006-01"))
	fmt.Printf("Clients processed: %d\n", summary.ClientsProcessed)
	fmt.Printf("Clients succeeded: %d\n", summary.ClientsSucceeded)
	fmt.Printf("Clients failed: %d\n", summary.ClientsFailed)
	fmt.Printf("Clients skipped: %d\n", summary.ClientsSkipped)
	fmt.Printf("Total locations: %d\n", summary.TotalLocations)
	fmt.Printf("Total reviews: %d\n", summary.TotalReviews)
	fmt.Printf("Total reviews analyzed: %d\n", summary.TotalReviewsAnalyzed)
	fmt.Printf("PDFs generated: %d\n", summary.PDFsGenerated)
	fmt.Printf("Emails sent: %d\n", summary.EmailsSent)
	fmt.Printf("Total processing time: %v\n", summary.ElapsedTime)

	if len(summary.FailedClients) > 0 {
		fmt.Printf("\nFailed Clients:\n")
		for _, fc := range summary.FailedClients {
			retryInfo := retryTracking[fc.ClientID][targetMonthFormatted]
			fmt.Printf("- Client %d (%s): %s (Attempt %d/%d)\n",
				fc.ClientID, fc.ClientName, fc.Error, retryInfo.Attempts, MaxRetryAttempts)
		}
	}

	// Send summary email if requested
	if *emailSummary != "" && emailSvc != nil {
		sendSummaryEmail(emailSvc, *emailSummary, summary, targetMonthFormatted)
	}
}

// loadRetryTracking loads the retry tracking data from file
func loadRetryTracking() RetryTracking {
	retryTracking := make(RetryTracking)

	// Check if file exists
	if _, err := os.Stat(RetryFilePath); os.IsNotExist(err) {
		return retryTracking
	}

	// Read file
	data, err := ioutil.ReadFile(RetryFilePath)
	if err != nil {
		log.Printf("Warning: Could not read retry tracking file: %v", err)
		return retryTracking
	}

	// Parse JSON
	err = json.Unmarshal(data, &retryTracking)
	if err != nil {
		log.Printf("Warning: Could not parse retry tracking data: %v", err)
		return make(RetryTracking)
	}

	return retryTracking
}

// saveRetryTracking saves the retry tracking data to file
func saveRetryTracking(retryTracking RetryTracking) {
	// Convert to JSON
	data, err := json.MarshalIndent(retryTracking, "", "  ")
	if err != nil {
		log.Printf("Warning: Could not marshal retry tracking data: %v", err)
		return
	}

	// Write to file
	err = ioutil.WriteFile(RetryFilePath, data, 0644)
	if err != nil {
		log.Printf("Warning: Could not write retry tracking file: %v", err)
	}
}

// sendSummaryEmail sends an email with the processing summary
func sendSummaryEmail(emailSvc email_service.EmailService, recipient string, summary google_my_business_api.ProcessingSummary, month string) {
	// Create email subject
	subject := fmt.Sprintf("Monthly Review Analysis Summary - %s", month)

	// Create email body
	body := fmt.Sprintf(`
Monthly Review Analysis Summary for %s

Processing Results:
- Clients processed: %d
- Clients succeeded: %d
- Clients failed: %d
- Clients skipped: %d
- Total locations: %d
- Total reviews: %d
- Total reviews analyzed: %d
- PDFs generated: %d
- Emails sent: %d
- Processing time: %v

`,
		month,
		summary.ClientsProcessed,
		summary.ClientsSucceeded,
		summary.ClientsFailed,
		summary.ClientsSkipped,
		summary.TotalLocations,
		summary.TotalReviews,
		summary.TotalReviewsAnalyzed,
		summary.PDFsGenerated,
		summary.EmailsSent,
		summary.ElapsedTime,
	)

	// Add failed clients if any
	if len(summary.FailedClients) > 0 {
		body += "\nFailed Clients:\n"
		for _, fc := range summary.FailedClients {
			body += fmt.Sprintf("- Client %d (%s): %s\n", fc.ClientID, fc.ClientName, fc.Error)
		}
	}

	// Send email using the new plain text email method
	err := emailSvc.SendPlainTextEmail(
		subject,
		recipient,
		"Admin", // Recipient name
		body,
	)

	if err != nil {
		log.Printf("Warning: Could not send summary email: %v", err)
	} else {
		log.Printf("Summary email sent to %s", recipient)
	}
}

// Helper function to get minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
