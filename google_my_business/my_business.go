// Google My Business used to respond to Google reviews and get review stats.
//
// The credentials have been constructed using information from google sheets API: https://developers.google.com/sheets/api/quickstart/go?authuser=2
// which I used in another project and modifying the downloaded credentials.json file using the
// documents online about google my business: https://developers.google.com/my-business/content/implement-oauth
// (also see https://developers.google.com/my-business/content/basic-setup)
//
// This involved modifying:
// "client_id":"YOUR_CLIENT_ID.apps.googleusercontent.com"
// "project_id":"YOUR_PROJECT_ID"
// "token_uri":"https://www.googleapis.com/oauth2/v4/token"
// "client_secret":"YOUR_CLIENT_SECRET"
//
// These are acquired from logging into the Google API Console: https://console.developers.google.com/apis/credentials
// First select the correct project near the top of the page (the project_id can be got from here).
// The create a OAuth 2.0 Client ID using the CREATE CREDENTIATIALS option, selecting OAauth client ID and
// selecting Desktop app (IMPORTANT) as the Application type.
//
// NOTE: The credentials.json file can be downloaded from the Google API Console under credentials menu.
//
// To allow sending emails from apps it is necessary to turn on less secure app
// access on Google see: https://support.google.com/accounts/answer/6010255
// NOTE: Google will automatically turn this setting off if it's not being used.
// This is quite useful also: https://gist.github.com/jpillora/cb46d183eca0710d909a
//
// This application will run on a Linux 64 bit system so will need compiling with the following:
//
//	$ env GOOS=linux GOARCH=amd64 go build my_business.go
//
// Use this in a cron job (use via a shell script for changing into the desired directory)
// for responding to reviews use no flags
// for running report (done each month) use the -reportmonthback <integer of how many months back to run monthly report for>
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gopkg.in/natefinch/lumberjack.v2"

	"google_my_business/ai_service"
	"google_my_business/ai_service/review"
	"google_my_business/config"
	"google_my_business/database"
	"google_my_business/email_service"
	"google_my_business/google_my_business_api"
	"shared_templates"

	_ "github.com/go-sql-driver/mysql"
)

// token from Google OAuth
var googleToken *oauth2.Token

// RetryTracking represents the structure for tracking client retry attempts
type RetryTracking map[int]map[string]ClientRetryInfo

// ClientRetryInfo tracks retry information for a specific client and month
type ClientRetryInfo struct {
	Month     string `json:"month"`
	Attempts  int    `json:"attempts"`
	LastError string `json:"last_error"`
}

// Constants for monthly analysis
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

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		getTokenFromWeb(config)
		if googleToken != nil {
			saveToken(tokFile, googleToken)
		}
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, the token is set in the global variable googleToken.
func getTokenFromWeb(config *oauth2.Config) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	log.Printf("getTokenFromWeb: starting HTTP server\n")
	httpServerExitDone := &sync.WaitGroup{}
	httpServerExitDone.Add(1)
	startHttpServerForTokenUri(httpServerExitDone, config)
	// wait for goroutine started in startHttpServerForTokenUri() to stop
	httpServerExitDone.Wait()
	log.Printf("getTokenFromWeb: done. exiting\n")
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

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// start HTTP server for receiving the OAuth token from Google.
// The credentials.json file contain the "redirect_uris":["http://localhost"] to reference this server.
// request example: http://localhost/?state=state-token&code=4/0AX4XfWjODE_aomPgdCJRX5AWAzfzIWUbRxrZJ_M4Z2NY8VvdyeH6XbsM1Sm3sYLcIwgFag&scope=https://www.googleapis.com/auth/business.manage
func startHttpServerForTokenUri(wg *sync.WaitGroup, config *oauth2.Config) {
	srv := &http.Server{
		Addr: ":80",
	}
	idleConnsClosed := make(chan struct{})

	// sends an HTTP response to the browser
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("startHttpServerForTokenUri: request raw query: %s\n", r.URL.RawQuery)
		authCode := r.URL.Query().Get("code")
		if authCode != "" {
			log.Println("startHttpServerForTokenUri: Found authorization code:", authCode)
			var err error
			googleToken, err = config.Exchange(context.Background(), authCode)
			if err != nil {
				log.Fatalf("startHttpServerForTokenUri: Unable to retrieve token from web: %v", err)
			}
			log.Println("startHttpServerForTokenUri: Shutting down ...")
			// graceful-shutdown
			err = srv.Shutdown(context.Background())
			if err != nil {
				log.Printf("startHttpServerForTokenUri: server.Shutdown: %v\n", err)
			}
			close(idleConnsClosed)
		} else {
			fmt.Fprintln(w, "No token found") // Just display a default page if the query string does not have a token (code) in it return to server at /
		}
	})

	go func() {
		defer wg.Done() // let calling function know done cleaning up
		// always return error. ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Printf("startHttpServerForTokenUri: Error ListenAndServer(): %v\n", err)
		}
		log.Println("startHttpServerForTokenUri: Bye.")
	}()
}

// getOAuthClient creates an authenticated Google OAuth client
func getOAuthClient() (*http.Client, error) {
	// Read OAuth credentials file
	credFile := "credentials.json"
	b, err := os.ReadFile(credFile)
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

// send - send email
// returns true if successfully sent
func send(smtpServer string, smtpServerPort string, password string, from string, to string, subject string, msg string) bool {
	// Set up authentication information.
	auth := smtp.PlainAuth("", from, password, smtpServer)

	// Connect to the server, authenticate, set the sender and recipient
	// and send the email all in one step.

	// check recipient
	// sendTo := strings.Split(to, ",")
	toSplit := regexp.MustCompile(` *, *`)
	sendTo := toSplit.Split(to, -1)
	// check that sendTo is not empty (could be array of empty strings)
	sendToEmpty := true
	for i := 0; i < len(sendTo); i++ {
		if len(sendTo[i]) > 0 {
			sendToEmpty = false
			break
		}
	}
	if sendToEmpty {
		return true
	}

	message := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		msg + "\r\n")

	err := smtp.SendMail(smtpServer+":"+smtpServerPort, auth, from, sendTo, message)
	if err != nil {
		log.Printf("error sending email: %s\n", err)
		return false
	}
	return true
}

// Process - process - this is public to help with testing (command line flags e.g. -reportmonthback 2)
func Process(reportMonthBack int, reportNameOrPostalCodeNotFound bool, csvReportOnly bool) {
	// ---------------
	// Properties file
	// ---------------

	props := config.ReadProperties()
	// fmt.Printf("ReviewsNotBeforeDays = %d\n", props.ReviewsNotBeforeDays)

	// --------
	// Database
	// --------

	db := database.OpenDB(props.DbName, props.DbAddress, props.DbPort, props.DbUsername, props.DbPassword)
	if db == nil {
		log.Fatal("Error opening database")
	}

	// ------------------
	// Google my business
	// ------------------

	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/business.manage")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	// see: https://github.com/googleapis/google-api-go-client
	// This is the context that the client has been already set up with (only required if want to reference for other operations in the same context)
	// ctx := context.Background()

	accounts := google_my_business_api.GetAccounts(client)
	// fmt.Println(accounts)

	// CSV report
	csvReport := "client ID,location name,postal code,report month,report year,unspecified star rating," +
		"one star rating,two star rating,three star rating,four star rating," +
		"five star rating,number of times the business profile call button was clicked," +
		"number of times the business profile website was clicked\n"
	var (
		csvReportMonth time.Month
		csvReportYear  int
	)

	// location name and / or postal code not found report
	var nopgrcfgmblns []database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode
	if reportNameOrPostalCodeNotFound {
		// set the name or postal code not found list initially to all the configs with reply to google my business enabled
		nopgrcfgmblns = database.AllConfigsWithReplyToGoogleMyBusiness(db)
	}

	for _, a := range accounts {
		var grcfgmblns []database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode
		if reportMonthBack > -1 {
			grcfgmblns = google_my_business_api.GetLocations(client, a, db, database.LookupModeReport)
		} else {
			grcfgmblns = google_my_business_api.GetLocations(client, a, db, database.LookupModeReply)
		}
		// fmt.Printf("locations: %+v\n", grcfgmblns)
		for _, g := range grcfgmblns {
			if reportNameOrPostalCodeNotFound {
				// remove from name or postal code not found list
				var tmp []database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode
				for _, x := range nopgrcfgmblns {
					if x.GoogleMyBusinessLocationName == g.GoogleMyBusinessLocationName &&
						x.GoogleMyBusinessPostalCode == g.GoogleMyBusinessPostalCode {
						continue
					}
					tmp = append(tmp, x)
				}
				nopgrcfgmblns = tmp
			} else if reportMonthBack > -1 {
				csvReport += fmt.Sprintf("%d,\"%s\",\"%s\",", g.ClientID, g.GoogleMyBusinessLocationName, g.GoogleMyBusinessPostalCode)
				if g.GoogleMyBusinessReportEnabled && len(strings.Trim(g.EmailAddress, " ")) > 0 {
					review_ratings, reportYear, reportMonth := google_my_business_api.ReportOnReviews(client, g, reportMonthBack)
					csvReportMonth = reportMonth
					csvReportYear = reportYear
					csvReport += fmt.Sprintf("\"%s\",%d,", csvReportMonth, csvReportYear)
					// fmt.Printf("location: %s, %s\nreview_ratings: %+v\n", g.GoogleMyBusinessLocationName, g.GoogleMyBusinessPostalCode, review_ratings)
					// format message:
					starRatings := ""
					if v, ok := review_ratings[google_my_business_api.StarRatingUnspecified]; ok {
						starRatings = fmt.Sprintf("unspecified rating: %d\n", v)
						csvReport += fmt.Sprintf("%d,", v)
					} else {
						csvReport += "0,"
					}
					if v, ok := review_ratings[google_my_business_api.StarRatingOne]; ok {
						starRatings += fmt.Sprintf("one star rating: %d\n", v)
						csvReport += fmt.Sprintf("%d,", v)
					} else {
						csvReport += "0,"
					}
					if v, ok := review_ratings[google_my_business_api.StarRatingTwo]; ok {
						starRatings += fmt.Sprintf("two star rating: %d\n", v)
						csvReport += fmt.Sprintf("%d,", v)
					} else {
						csvReport += "0,"
					}
					if v, ok := review_ratings[google_my_business_api.StarRatingThree]; ok {
						starRatings += fmt.Sprintf("three star rating: %d\n", v)
						csvReport += fmt.Sprintf("%d,", v)
					} else {
						csvReport += "0,"
					}
					if v, ok := review_ratings[google_my_business_api.StarRatingFour]; ok {
						starRatings += fmt.Sprintf("four star rating: %d\n", v)
						csvReport += fmt.Sprintf("%d,", v)
					} else {
						csvReport += "0,"
					}
					if v, ok := review_ratings[google_my_business_api.StarRatingFive]; ok {
						starRatings += fmt.Sprintf("five star rating: %d\n", v)
						csvReport += fmt.Sprintf("%d,", v)
					} else {
						csvReport += "0,"
					}
					msg := fmt.Sprintf("%s, %s (%s %d)\n\n%s\n", g.GoogleMyBusinessLocationName, g.GoogleMyBusinessPostalCode, reportMonth, reportYear, starRatings)
					// fmt.Println(msg)

					// insights
					insights := google_my_business_api.ReportOnInsights(client, g, reportMonthBack)
					msg += fmt.Sprintf("\nInsights:\n%s\n", insights)

					// CSV report minor bodge with string maniplutaion rather than getting numbers in functions
					// each insight is separated by a line and of the form: <text>: <number>
					a1 := strings.Split(insights, "\n")
					for i, s := range a1 {
						// ignore last array item it is empty since insights string has newline at end
						if i+1 < len(a1) {
							mt := strings.Split(s, ":")
							if len(mt) == 2 {
								m := strings.TrimSpace(mt[1])
								csvReport += m
							} else {
								csvReport += "0"
							}
							if i+2 < len(a1) {
								csvReport += ","
							}
						}
					}
					csvReport += "\n"

					// send email
					// Initially removed because it causes Google to not send email with error: Too many login attempts, when there are a lot of companies, this I guess is to stop spam.
					// This has been restricted by only sending to those email addresses that are not the same as the CSV report.
					// This can be disabled by setting the csvReportOnly flag to true
					if !csvReportOnly && (strings.TrimSpace(g.EmailAddress) != strings.TrimSpace(props.EmailReportTo)) {
						emailSubject := fmt.Sprintf("%s (%s %d)", props.EmailSubject, reportMonth, reportYear)
						send(props.SMTPServer, props.SMTPServerPort, props.EmailPassword, props.EmailFrom, g.EmailAddress, emailSubject, msg)
					}
				} else {
					csvReport += fmt.Sprintf("\"%s\",%d,0,0,0,0,0,0,0,0\n", csvReportMonth, csvReportYear)
				}
			} else {
				// Get multi message (reply) separator
				sep := strings.Trim(g.MultiMessageSeparator, " ")
				google_my_business_api.ProcessReviews(client, g, props.ReviewsNotBeforeDays, sep)
			}
		}
	}
	if reportNameOrPostalCodeNotFound {
		// CSV report for location name and / or postal code not found
		nameOrPostalCodeNotFoundReport := "All location names (and postal codes) have been found\n"
		if len(nopgrcfgmblns) > 0 {
			nameOrPostalCodeNotFoundReport = "client ID,location name,postal code\n"
			for _, n := range nopgrcfgmblns {
				nameOrPostalCodeNotFoundReport += fmt.Sprintf("%d,%s,%s\n", n.ClientID, n.GoogleMyBusinessLocationName, n.GoogleMyBusinessPostalCode)
			}
		}
		// send name or postal code not found report email
		emailSubject := fmt.Sprintf("%s (%d %s %d)", props.EmailNameOrPostalCodeNotFoundReportSubject, time.Now().Day(), time.Now().Month(), time.Now().Year())
		send(props.SMTPServer, props.SMTPServerPort, props.EmailPassword, props.EmailFrom, props.EmailReportTo, emailSubject, nameOrPostalCodeNotFoundReport)
	} else if reportMonthBack > -1 {
		// send CSV report email
		emailSubject := fmt.Sprintf("%s (%s %d)", props.EmailCsvReportSubject, csvReportMonth, csvReportYear)
		send(props.SMTPServer, props.SMTPServerPort, props.EmailPassword, props.EmailFrom, props.EmailReportTo, emailSubject, csvReport)
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
	data, err := os.ReadFile(RetryFilePath)
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
	err = os.WriteFile(RetryFilePath, data, 0644)
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

// runMonthlyAnalysis runs the monthly analysis process
func runMonthlyAnalysis(targetMonthStr string, forceReprocess bool, retryOnly bool, emailSummary string, noEmail bool, noSave bool, clientsFilter string, debugMode bool) {
	// Parse and validate client filter if provided
	var clientIDs []int
	if clientsFilter != "" {
		// Split by comma and trim whitespace
		clientStrs := strings.Split(clientsFilter, ",")
		for _, clientStr := range clientStrs {
			clientStr = strings.TrimSpace(clientStr)
			if clientStr == "" {
				continue // Skip empty strings
			}
			
			// Parse to integer
			clientID, err := strconv.Atoi(clientStr)
			if err != nil {
				log.Fatalf("Invalid client ID '%s': must be a valid integer", clientStr)
			}
			
			if clientID <= 0 {
				log.Fatalf("Invalid client ID '%d': must be a positive integer", clientID)
			}
			
			clientIDs = append(clientIDs, clientID)
		}
		
		if len(clientIDs) == 0 {
			log.Fatalf("No valid client IDs found in filter: %s", clientsFilter)
		}
		
		fmt.Printf("Processing specific clients: %v\n", clientIDs)
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
	if targetMonthStr != "" {
		if tm, err := time.Parse("2006-01", targetMonthStr); err == nil {
			targetMonth = tm
		} else {
			log.Printf("Warning: Couldn't parse month '%s', using previous month. Expected format: YYYY-MM", targetMonthStr)
		}
	}

	targetMonthFormatted := targetMonth.Format("2006-01")
	fmt.Printf("Running monthly analysis for %s\n", targetMonthFormatted)

	// Load retry tracking if in retry mode or to update after processing
	retryTracking := loadRetryTracking()

	// If in retry-only mode, determine which clients to process
	var clientsToProcess []int
	if retryOnly {
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
	summary, err := google_my_business_api.AnalyzeClientReviews(dbAdapter, httpClient, analyzer, &targetMonth, forceReprocess, emailSvc, debugMode, noSave, noEmail, clientIDs)
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
	if emailSummary != "" && emailSvc != nil {
		sendSummaryEmail(emailSvc, emailSummary, summary, targetMonthFormatted)
	}
}

func main() {

	// -----------------------------------------------
	// Get command line flags
	// -----------------------------------------------

	// logging flag useful for output during testing with go run ...
	var testLog bool
	flag.BoolVar(&testLog, "testlog", false, "set for logging output shown on console when running with go run ...")

	// report how many months back is passed as a flag on command line (e.g. my_business -reportmonthback 2)
	// default set to -1 which means do not do report.
	// NOTE: flags appear to have to come before other command line arguments.
	var reportMonthBack int
	flag.IntVar(&reportMonthBack, "reportmonthback", -1, "Report on how many months back from this month, defaults -1 i.e. not set")

	// report on name or postal code not found in google my business locations
	// default is false i.e. not set
	var reportNameOrPostalCodeNotFound bool
	flag.BoolVar(&reportNameOrPostalCodeNotFound, "reportnameorpostalcodenotfound", false, "Report on name or postal code not found in google my business locations, default false i.e. not enabled")

	// CSV report only, this prevents sending the individual reports
	// default is false i.e. not set
	var csvReportOnly bool
	flag.BoolVar(&csvReportOnly, "csvreportonly", false, "Send CSV report only, do not send the individual reports, default false i.e. not enabled, send all reports")

	// Monthly analysis command
	var runMonthlyAnalysisCmd bool
	flag.BoolVar(&runMonthlyAnalysisCmd, "run-monthly-analysis", false, "Run the monthly analysis process")

	// Monthly analysis options
	var monthlyAnalysisMonth string
	flag.StringVar(&monthlyAnalysisMonth, "month", "", "Target month for monthly analysis in YYYY-MM format (defaults to previous month)")

	var forceReprocess bool
	flag.BoolVar(&forceReprocess, "force-reprocess", false, "Force reprocessing of existing reports")

	var retryOnly bool
	flag.BoolVar(&retryOnly, "retry-only", false, "Only retry previously failed clients")

	var emailSummary string
	flag.StringVar(&emailSummary, "email-summary", "", "Email address to send processing summary")

	// New debugging and control flags for monthly analysis
	var noEmail bool
	flag.BoolVar(&noEmail, "no-email", false, "Save reports but skip email sending")

	var noSave bool
	flag.BoolVar(&noSave, "no-save", false, "Skip database save operations")

	var clientsFilter string
	flag.StringVar(&clientsFilter, "clients", "", "Process only specified client IDs (comma-separated)")

	var debugMode bool
	flag.BoolVar(&debugMode, "debug", false, "Print detailed debug information to stdout")

	flag.Parse()

	// Process positional arguments for monthly analysis (for backward compatibility)
	// If first argument is not a flag and looks like a date, use it as target month
	args := flag.Args()
	if runMonthlyAnalysisCmd && len(args) > 0 && len(args[0]) == 7 && args[0][4] == '-' {
		monthlyAnalysisMonth = args[0]
	}

	// Handle monthly analysis command
	if runMonthlyAnalysisCmd {
		runMonthlyAnalysis(monthlyAnalysisMonth, forceReprocess, retryOnly, emailSummary, noEmail, noSave, clientsFilter, debugMode)
		return
	}

	// -----------------------------------------------
	// Create log file in same directory as executable
	// -----------------------------------------------

	logFilename := ""
	if !testLog {
		// create log file in same directory as executable
		ex, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(filepath.Dir(ex))
		logFilename = filepath.Dir(ex) + "/my_business.log"
		// f, err := os.OpenFile(filepath.Dir(ex)+"/google_reviews.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		f, err := os.OpenFile(logFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		// log.SetOutput(f)
		log.SetOutput(&lumberjack.Logger{
			Filename:   logFilename,
			MaxSize:    20, // megabytes
			MaxBackups: 3,
			MaxAge:     28, //days
		})
	}

	// NOTE: if reportNameOrPostalCodeNotFound is set to true then should also set the reportMonthBack to greater than -1 e.g. 0 to get required info from locations
	if reportNameOrPostalCodeNotFound {
		reportMonthBack = 0
	}

	// log.Printf("reportMonthBack: %d, reportNameOrPostalCodeNotFound: %t, csvReportOnly: %t\n", reportMonthBack, reportNameOrPostalCodeNotFound, csvReportOnly)

	Process(reportMonthBack, reportNameOrPostalCodeNotFound, csvReportOnly)
}
