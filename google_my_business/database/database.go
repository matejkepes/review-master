package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"shared-templates"
	"log"
	"time"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode - represents the google reviews config fields required to reply to a review from the google my business location name and postal address
type GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode struct {
	ClientID                                      uint64
	GoogleMyBusinessLocationName                  string
	GoogleMyBusinessLocationPath                  string
	GoogleMyBusinessPostalCode                    string
	GoogleMyBusinessReplyToUnspecfifiedStarRating bool
	GoogleMyBusinessUnspecfifiedStarRatingReply   string
	GoogleMyBusinessReplyToOneStarRating          bool
	GoogleMyBusinessOneStarRatingReply            string
	GoogleMyBusinessReplyToTwoStarRating          bool
	GoogleMyBusinessTwoStarRatingReply            string
	GoogleMyBusinessReplyToThreeStarRating        bool
	GoogleMyBusinessThreeStarRatingReply          string
	GoogleMyBusinessReplyToFourStarRating         bool
	GoogleMyBusinessFourStarRatingReply           string
	GoogleMyBusinessReplyToFiveStarRating         bool
	GoogleMyBusinessFiveStarRatingReply           string
	GoogleMyBusinessReportEnabled                 bool
	TimeZone                                      string
	MultiMessageSeparator                         string
	EmailAddress                                  string
	ContactMethod                                 *string // Nullable field
	AIResponsesEnabled                            bool
}

// Lookup modes for ConfigFromGoogleMyBusinessLocationNameAndPostalCode
const (
	LookupModeReply    = 0 // For review replies
	LookupModeReport   = 1 // For regular reports
	LookupModeAnalysis = 2 // For monthly analysis
)

// OpenDB - open database connection
func OpenDB(database string, host string, port string, username string, password string) *sql.DB {
	// NOTE: ?parseTime=true which allows DATE and DATETIME database types to be parsed into golang time.Time
	db, err := sql.Open("mysql", username+":"+password+"@tcp("+host+":"+port+")/"+database+"?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	// database connection is pooled and used by many connections, even using defer db.Close() causes issues with these
	// where it indicates the connection is closed.
	// defer db.Close()
	// Test connection to the database (this will result in the program ending if unable to establish a connection)
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}
	return db
}

// GoogleMyBusinessLocationName - get the config required fields for a reviews reply from the google my business location name and postal code address
// This is also used for reporting where it is used to fetch info from google my business.
// mode parameter controls which condition to use:
// - LookupModeReply (0): For review replies, checks google_my_business_review_reply_enabled
// - LookupModeReport (1): For regular reports, checks google_my_business_report_enabled
// - LookupModeAnalysis (2): For monthly analysis, checks monthly_review_analysis_enabled
func ConfigFromGoogleMyBusinessLocationNameAndPostalCode(db *sql.DB, googleMyBusinessLocationName string, googleMyBusinessPostalCode string, mode int) GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode {
	qry := "SELECT config.google_my_business_location_name," +
		" config.google_my_business_postal_code," +
		" config.google_my_business_reply_to_unspecfified_star_rating," +
		" config.google_my_business_unspecfified_star_rating_reply," +
		" config.google_my_business_reply_to_one_star_rating," +
		" config.google_my_business_one_star_rating_reply," +
		" config.google_my_business_reply_to_two_star_rating," +
		" config.google_my_business_two_star_rating_reply," +
		" config.google_my_business_reply_to_three_star_rating," +
		" config.google_my_business_three_star_rating_reply," +
		" config.google_my_business_reply_to_four_star_rating," +
		" config.google_my_business_four_star_rating_reply," +
		" config.google_my_business_reply_to_five_star_rating," +
		" config.google_my_business_five_star_rating_reply," +
		" config.google_my_business_report_enabled," +
		" config.email_address," +
		" config.multi_message_separator," +
		" config.time_zone," +
		" config.client_id," +
		" config.contact_method," +
		" config.ai_responses_enabled" +
		" FROM google_reviews_configs AS config" +
		" JOIN clients AS client ON client.id = config.client_id" +
		" WHERE config.enabled = 1" +
		// " AND config.google_my_business_review_reply_enabled = 1" +
		" AND TRIM(config.google_my_business_location_name) = TRIM(?)" +
		" AND REPLACE(config.google_my_business_postal_code, ' ', '') = REPLACE(?, ' ', '')" +
		" AND client.enabled = 1"

	// Add the appropriate condition based on the mode
	switch mode {
	case LookupModeReport:
		qry += " AND config.google_my_business_report_enabled = 1"
	case LookupModeAnalysis:
		qry += " AND config.monthly_review_analysis_enabled = 1"
	case LookupModeReply:
		fallthrough
	default:
		// Default to reply mode for backward compatibility
		qry += " AND config.google_my_business_review_reply_enabled = 1"
	}

	rows, err := db.Query(qry, googleMyBusinessLocationName, googleMyBusinessPostalCode)
	var grcfgmbln GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode
	if err != nil {
		log.Println("Error retrieving google my business location name", googleMyBusinessLocationName, "from database. Error: ", err)
		return grcfgmbln
	}
	defer rows.Close()
	for rows.Next() {
		if err1 := rows.Scan(&grcfgmbln.GoogleMyBusinessLocationName,
			&grcfgmbln.GoogleMyBusinessPostalCode,
			&grcfgmbln.GoogleMyBusinessReplyToUnspecfifiedStarRating,
			&grcfgmbln.GoogleMyBusinessUnspecfifiedStarRatingReply,
			&grcfgmbln.GoogleMyBusinessReplyToOneStarRating,
			&grcfgmbln.GoogleMyBusinessOneStarRatingReply,
			&grcfgmbln.GoogleMyBusinessReplyToTwoStarRating,
			&grcfgmbln.GoogleMyBusinessTwoStarRatingReply,
			&grcfgmbln.GoogleMyBusinessReplyToThreeStarRating,
			&grcfgmbln.GoogleMyBusinessThreeStarRatingReply,
			&grcfgmbln.GoogleMyBusinessReplyToFourStarRating,
			&grcfgmbln.GoogleMyBusinessFourStarRatingReply,
			&grcfgmbln.GoogleMyBusinessReplyToFiveStarRating,
			&grcfgmbln.GoogleMyBusinessFiveStarRatingReply,
			&grcfgmbln.GoogleMyBusinessReportEnabled,
			&grcfgmbln.EmailAddress,
			&grcfgmbln.MultiMessageSeparator,
			&grcfgmbln.TimeZone,
			&grcfgmbln.ClientID,
			&grcfgmbln.ContactMethod,
			&grcfgmbln.AIResponsesEnabled); err1 != nil {
			log.Println("Error retrieving google my business location name", googleMyBusinessLocationName, "from database whilst reading returned results. Error: ", err1)
			return grcfgmbln
		}
	}
	return grcfgmbln
}

// AllConfigsWithReplyToGoogleMyBusiness - get all the google my business enabled configs required fields for a reviews reply
// This is used to find all the google my businesses that should be replied to.
// The locations found by name and postal code are removed from this list to find the names and / or postal codes that may have changed.
func AllConfigsWithReplyToGoogleMyBusiness(db *sql.DB) []GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode {
	qry := "SELECT config.google_my_business_location_name," +
		" config.google_my_business_postal_code," +
		" config.google_my_business_reply_to_unspecfified_star_rating," +
		" config.google_my_business_unspecfified_star_rating_reply," +
		" config.google_my_business_reply_to_one_star_rating," +
		" config.google_my_business_one_star_rating_reply," +
		" config.google_my_business_reply_to_two_star_rating," +
		" config.google_my_business_two_star_rating_reply," +
		" config.google_my_business_reply_to_three_star_rating," +
		" config.google_my_business_three_star_rating_reply," +
		" config.google_my_business_reply_to_four_star_rating," +
		" config.google_my_business_four_star_rating_reply," +
		" config.google_my_business_reply_to_five_star_rating," +
		" config.google_my_business_five_star_rating_reply," +
		" config.google_my_business_report_enabled," +
		" config.email_address," +
		" config.multi_message_separator," +
		" config.time_zone," +
		" config.client_id," +
		" config.contact_method," +
		" config.ai_responses_enabled" +
		" FROM google_reviews_configs AS config" +
		" JOIN clients AS client ON client.id = config.client_id" +
		" WHERE config.enabled = 1" +
		" AND config.google_my_business_review_reply_enabled = 1" +
		" AND (config.google_my_business_reply_to_unspecfified_star_rating = 1" +
		"   OR config.google_my_business_reply_to_one_star_rating = 1" +
		"   OR config.google_my_business_reply_to_two_star_rating = 1" +
		"   OR config.google_my_business_reply_to_three_star_rating = 1" +
		"   OR config.google_my_business_reply_to_four_star_rating = 1" +
		"   OR config.google_my_business_reply_to_five_star_rating = 1)" +
		" AND client.enabled = 1"
	rows, err := db.Query(qry)
	var grcfgmblns []GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode
	if err != nil {
		log.Println("Error retrieving all with reply to google my business from database. Error: ", err)
		return grcfgmblns
	}
	defer rows.Close()
	var grcfgmbln GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode
	for rows.Next() {
		if err1 := rows.Scan(&grcfgmbln.GoogleMyBusinessLocationName,
			&grcfgmbln.GoogleMyBusinessPostalCode,
			&grcfgmbln.GoogleMyBusinessReplyToUnspecfifiedStarRating,
			&grcfgmbln.GoogleMyBusinessUnspecfifiedStarRatingReply,
			&grcfgmbln.GoogleMyBusinessReplyToOneStarRating,
			&grcfgmbln.GoogleMyBusinessOneStarRatingReply,
			&grcfgmbln.GoogleMyBusinessReplyToTwoStarRating,
			&grcfgmbln.GoogleMyBusinessTwoStarRatingReply,
			&grcfgmbln.GoogleMyBusinessReplyToThreeStarRating,
			&grcfgmbln.GoogleMyBusinessThreeStarRatingReply,
			&grcfgmbln.GoogleMyBusinessReplyToFourStarRating,
			&grcfgmbln.GoogleMyBusinessFourStarRatingReply,
			&grcfgmbln.GoogleMyBusinessReplyToFiveStarRating,
			&grcfgmbln.GoogleMyBusinessFiveStarRatingReply,
			&grcfgmbln.GoogleMyBusinessReportEnabled,
			&grcfgmbln.EmailAddress,
			&grcfgmbln.MultiMessageSeparator,
			&grcfgmbln.TimeZone,
			&grcfgmbln.ClientID,
			&grcfgmbln.ContactMethod,
			&grcfgmbln.AIResponsesEnabled); err1 != nil {
			log.Println("Error retrieving all with reply to google my business from database whilst reading returned results. Error: ", err1)
		}
		grcfgmblns = append(grcfgmblns, grcfgmbln)
	}
	return grcfgmblns
}

// SaveClientReport stores a complete client report with multiple location analyses
func SaveClientReport(db *sql.DB, clientID int, periodStart, periodEnd time.Time, locationResults []byte) (int64, error) {
	// Prepare the insert statement
	query := `
		INSERT INTO client_reports 
		(client_id, report_period_start, report_period_end, generated_at, locations) 
		VALUES (?, ?, ?, NOW(), ?)
	`

	// Execute the insert
	result, err := db.Exec(query, clientID, periodStart, periodEnd, locationResults)
	if err != nil {
		return 0, fmt.Errorf("failed to insert client report: %w", err)
	}

	// Get the inserted ID
	reportID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get inserted report ID: %w", err)
	}

	return reportID, nil
}

// GetClientReportByID retrieves a specific report by its ID
func GetClientReportByID(db *sql.DB, reportID int64) (*shared_templates.ClientReportData, error) {
	query := `
		SELECT report_id, client_id, report_period_start, report_period_end, 
			   generated_at, locations
		FROM client_reports
		WHERE report_id = ?
	`

	var reportID2 int64
	var clientID int
	var periodStart, periodEnd, generatedAt time.Time
	var locationResultsJSON []byte

	err := db.QueryRow(query, reportID).Scan(
		&reportID2,
		&clientID,
		&periodStart,
		&periodEnd,
		&generatedAt,
		&locationResultsJSON,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve client report: %w", err)
	}

	// Get client name (assuming a clients table with name column)
	var clientName string
	err = db.QueryRow("SELECT name FROM clients WHERE id = ?", clientID).Scan(&clientName)
	if err != nil {
		clientName = fmt.Sprintf("Client %d", clientID) // Fallback if name not found
	}

	// Unmarshal the location results JSON
	var locationResults []shared_templates.AnalysisResult
	if err := json.Unmarshal(locationResultsJSON, &locationResults); err != nil {
		return nil, fmt.Errorf("failed to unmarshal location results: %w", err)
	}

	// Create and return the client report data
	return &shared_templates.ClientReportData{
		ReportID:        reportID,
		ClientID:        clientID,
		ClientName:      clientName,
		PeriodStart:     periodStart,
		PeriodEnd:       periodEnd,
		GeneratedAt:     generatedAt,
		LocationResults: locationResults,
	}, nil
}

// GetClientReportsByClientID retrieves all reports for a specific client
func GetClientReportsByClientID(db *sql.DB, clientID int) ([]*shared_templates.ClientReportData, error) {
	query := `
		SELECT report_id, client_id, report_period_start, report_period_end, 
			   generated_at, locations
		FROM client_reports
		WHERE client_id = ?
		ORDER BY report_period_start DESC
	`

	rows, err := db.Query(query, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to query client reports: %w", err)
	}
	defer rows.Close()

	// Get client name once (assuming a clients table with name column)
	var clientName string
	err = db.QueryRow("SELECT name FROM clients WHERE id = ?", clientID).Scan(&clientName)
	if err != nil {
		clientName = fmt.Sprintf("Client %d", clientID) // Fallback if name not found
	}

	var reports []*shared_templates.ClientReportData

	for rows.Next() {
		var reportID int64
		var clientID2 int
		var periodStart, periodEnd, generatedAt time.Time
		var locationResultsJSON []byte

		err := rows.Scan(
			&reportID,
			&clientID2,
			&periodStart,
			&periodEnd,
			&generatedAt,
			&locationResultsJSON,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan client report row: %w", err)
		}

		// Unmarshal the location results JSON
		var locationResults []shared_templates.AnalysisResult
		if err := json.Unmarshal(locationResultsJSON, &locationResults); err != nil {
			return nil, fmt.Errorf("failed to unmarshal location results: %w", err)
		}

		// Create the client report data
		report := &shared_templates.ClientReportData{
			ReportID:        reportID,
			ClientID:        clientID,
			ClientName:      clientName,
			PeriodStart:     periodStart,
			PeriodEnd:       periodEnd,
			GeneratedAt:     generatedAt,
			LocationResults: locationResults,
		}

		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating client report rows: %w", err)
	}

	return reports, nil
}

// GetClientReportByClientAndPeriod retrieves a report for a specific client and time period
func GetClientReportByClientAndPeriod(db *sql.DB, clientID int, periodStart, periodEnd time.Time) (*shared_templates.ClientReportData, error) {
	query := `
		SELECT report_id, client_id, report_period_start, report_period_end, 
			   generated_at, locations
		FROM client_reports
		WHERE client_id = ? 
		AND report_period_start = ? 
		AND report_period_end = ?
	`

	var reportID int64
	var clientID2 int
	var periodStart2, periodEnd2, generatedAt time.Time
	var locationResultsJSON []byte

	err := db.QueryRow(query, clientID, periodStart, periodEnd).Scan(
		&reportID,
		&clientID2,
		&periodStart2,
		&periodEnd2,
		&generatedAt,
		&locationResultsJSON,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve client report: %w", err)
	}

	// Get client name (assuming a clients table with name column)
	var clientName string
	err = db.QueryRow("SELECT name FROM clients WHERE id = ?", clientID).Scan(&clientName)
	if err != nil {
		clientName = fmt.Sprintf("Client %d", clientID) // Fallback if name not found
	}

	// Unmarshal the location results JSON
	var locationResults []shared_templates.AnalysisResult
	if err := json.Unmarshal(locationResultsJSON, &locationResults); err != nil {
		return nil, fmt.Errorf("failed to unmarshal location results: %w", err)
	}

	// Create and return the client report data
	return &shared_templates.ClientReportData{
		ReportID:        reportID,
		ClientID:        clientID,
		ClientName:      clientName,
		PeriodStart:     periodStart,
		PeriodEnd:       periodEnd,
		GeneratedAt:     generatedAt,
		LocationResults: locationResults,
	}, nil
}

// ClientWithMonthlyReviewAnalysis represents a client with monthly review analysis enabled
type ClientWithMonthlyReviewAnalysis struct {
	ClientID           int
	ClientName         string
	EmailAddress       string
	ReportEmailAddress sql.NullString // New field for storing the report_email_address from clients table
}

// GetClientsWithMonthlyReviewAnalysisEnabled returns all active clients with monthly review analysis enabled
func GetClientsWithMonthlyReviewAnalysisEnabled(db *sql.DB) ([]ClientWithMonthlyReviewAnalysis, error) {
	query := `
		SELECT 
			c.id AS client_id, 
			c.name AS client_name, 
			g.email_address,
			c.report_email_address
		FROM 
			clients AS c
		JOIN 
			google_reviews_configs AS g ON c.id = g.client_id
		WHERE 
			c.enabled = 1
			AND g.enabled = 1
			AND g.monthly_review_analysis_enabled = TRUE
		ORDER BY 
			c.name
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query clients with monthly review analysis: %w", err)
	}
	defer rows.Close()

	var clients []ClientWithMonthlyReviewAnalysis

	for rows.Next() {
		var client ClientWithMonthlyReviewAnalysis

		if err := rows.Scan(&client.ClientID, &client.ClientName, &client.EmailAddress, &client.ReportEmailAddress); err != nil {
			return nil, fmt.Errorf("failed to scan client row: %w", err)
		}

		clients = append(clients, client)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating client rows: %w", err)
	}

	return clients, nil
}
