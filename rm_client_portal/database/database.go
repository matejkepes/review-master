package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"shared_templates"
	"time"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// validateLocationResults safely unmarshals location results JSON with NULL and error handling
func validateLocationResults(data []byte) ([]shared_templates.AnalysisResult, error) {
	// Handle NULL or empty data
	if len(data) == 0 {
		log.Printf("Warning: NULL or empty location results data, returning empty slice")
		return []shared_templates.AnalysisResult{}, nil
	}

	var results []shared_templates.AnalysisResult
	if err := json.Unmarshal(data, &results); err != nil {
		log.Printf("Warning: Failed to unmarshal location results JSON: %v", err)
		return []shared_templates.AnalysisResult{}, nil // Return empty slice instead of error to prevent crashes
	}

	return results, nil
}

var Db *sql.DB

// GoogleReviewsConfigAndGoogleMyBusinessLocation - represents the google reviews config and the google my business location for reports
type GoogleReviewsConfigAndGoogleMyBusinessLocation struct {
	GoogleMyBusinessLocationName string `json:"location_name"` // Google My Business google my business location name
	GoogleMyBusinessPostalCode   string `json:"postal_code"`   // Google My Business google my business postal code address
	// GoogleMyBusinessReportEnabled bool   // Google My Business report enabled
	// EmailAddress                  string // Email address used for reporting
	GoogleMyBusinessLocationPath    string                                  `json:"-"`                // added later from the Google My Business Location API
	GoogleMyBusinessLocationAddress string                                  `json:"location_address"` // added later from the Google My Business Location API
	GoogleMyBusinessLocality        string                                  `json:"locality"`         // added later from the Google My Business Location API
	TimeZone                        string                                  `json:"time_zone"`        // this is for reporting the correct time period
	GoogleReviewRatings             GoogleReviewRatingsFromGoogleMyBusiness `json:"review_ratings"`   // added later from the Google My Business Location API
	GoogleInsights                  GoogleInsightsFromGoogleMyBusiness      `json:"insights"`         // added later from the Google My Business Location API
	ClientID                        uint64                                  `json:"client_id"`        // client id
}

// GoogleReviewRatingsFromGoogleMyBusiness - represents the google reviews star ratings from google my business
type GoogleReviewRatingsFromGoogleMyBusiness struct {
	StarRatingUnspecified int `json:"unspecified"`
	StarRatingOne         int `json:"one"`
	StarRatingTwo         int `json:"two"`
	StarRatingThree       int `json:"three"`
	StarRatingFour        int `json:"four"`
	StarRatingFive        int `json:"five"`
}

// GoogleInsightsFromGoogleMyBusiness - represents the insights from google business
type GoogleInsightsFromGoogleMyBusiness struct {
	NumberOfBusinessProfileCallButtonClicked int `json:"number_of_business_profile_call_button_clicked"`
	NumberOfBusinessProfileWebsiteClicked    int `json:"number_of_business_profile_website_clicked"`
}

// OpenDB - open database connection
func OpenDB(database string, host string, port string, username string, password string) {
	// NOTE: ?parseTime=true which allows DATE and DATETIME database types to be parsed into golang time.Time
	db, err := sql.Open("mysql", username+":"+password+"@tcp("+host+":"+port+")/"+database+"?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	// database connection is pooled and used by many connections, even using defer db.Close() causes issues with these
	// where it indicates the connection is closed.
	// defer db.Close()

	// FIX: issues with this program using all the MySQL connections (max_connections - default 151)
	// this is because there are so many concurrent calls to this program and by default there is
	// no limit on the number of open connections (in-use + idle) at the same time
	//
	// Set the maximum number of concurrently open connections (in-use + idle)
	// to 40. Setting this to less than or equal to 0 will mean there is no
	// maximum limit (which is also the default setting).
	db.SetMaxOpenConns(40)

	// Test connection to the database (this will result in the program ending if unable to establish a connection)
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}
	Db = db
}

// type UserClient struct {
// 	UserID    uint
// 	Email     string
// 	ClientsID uint
// 	Name      string
// }

// type UserAndClients struct {
// 	Email   string        `json:"email"`
// 	Clients []UsersClient `json:"clients"`
// }

// UsersClient - represents a client for a user (identified by email)
type UsersClient struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// StatsRequest - represents a statistic request.
type StatsRequest struct {
	StartDay     string `json:"start_day"`     // start day
	EndDay       string `json:"end_day"`       // end day
	TimeGrouping string `json:"time_grouping"` // time grouping
}

// StatsNewResult - represents a statistic (new) result.
type StatsNewResult struct {
	ClientID    uint64 `json:"client_id"`    // client id
	ClientName  string `json:"client_name"`  // client name
	Sent        uint64 `json:"sent"`         // sent
	Requested   uint64 `json:"requested"`    // requested
	GroupPeriod string `json:"group_period"` // group period
}

func GetUser(email, password string) string {
	qry := "SELECT u.email" +
		" FROM users AS u" +
		" WHERE u.email = ?" +
		" AND u.password = SHA2(?, 224)"
	var e string
	err := Db.QueryRow(qry, email, password).Scan(&e)
	switch {
	case err == sql.ErrNoRows:
		return e
	case err != nil:
		log.Printf("query error: %v\n", err)
		return e
	default:
		// nothing
	}
	return e
}

// func GetUserAndClients(email, password string) UserAndClients {
// 	qry := "SELECT u.id, u.email," +
// 		" c.id, c.name" +
// 		" FROM users_clients AS u_c" +
// 		" JOIN users AS u ON u.id = u_c.user_id" +
// 		" JOIN clients AS c ON c.id = u_c.client_id" +
// 		" WHERE u.email = ?" +
// 		" AND u.password = SHA2(?, 224)" +
// 		" AND c.enabled = 1"
// 	rows, err := Db.Query(qry, email, password)
// 	var ucs UserAndClients
// 	var uc UserClient
// 	if err != nil {
// 		log.Println("Error retrieving user client for email:", email, "from database. Error:", err)
// 		return ucs
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		if err1 := rows.Scan(&uc.UserID, &uc.Email, &uc.ClientsID, &uc.Name); err1 != nil {
// 			log.Println("Error retrieving user client for email", email, "from database whilst reading returned results. Error: ", err1)
// 			return ucs
// 		}
// 		if uc.UserID == 0 {
// 			log.Printf("email %s not found\n", email)
// 			continue
// 		}
// 		ucs.Email = uc.Email
// 		var c UsersClient
// 		c.ID = uc.ClientsID
// 		c.Name = uc.Name
// 		ucs.Clients = append(ucs.Clients, c)
// 	}

// 	// marshall into json
// 	// r, _ := json.Marshal(ucs)
// 	// fmt.Sprintf("%s", r)

// 	return ucs
// }

// GetClientsForUserEmail - client list for a user (identified by email)
func GetClientsForUserEmail(email string) []UsersClient {
	qry := "SELECT c.id, c.name" +
		" FROM users_clients AS u_c" +
		" JOIN users AS u ON u.id = u_c.user_id" +
		" JOIN clients AS c ON c.id = u_c.client_id" +
		" WHERE u.email = ?" +
		" AND c.enabled = 1"
	rows, err := Db.Query(qry, email)
	var ucs []UsersClient
	if err != nil {
		log.Println("Error retrieving user clients for email:", email, "from database. Error:", err)
		return ucs
	}
	defer rows.Close()
	for rows.Next() {
		var uc UsersClient
		if err1 := rows.Scan(&uc.ID, &uc.Name); err1 != nil {
			log.Println("Error retrieving user client for email", email, "from database whilst reading returned results. Error: ", err1)
			return ucs
		}
		if uc.ID == 0 {
			log.Printf("client not found for email %s\n", email)
			continue
		}
		ucs = append(ucs, uc)
	}
	return ucs
}

// GetClientCheckUserEmail - client for a user (identified by email)
func GetClientCheckUserEmail(clientID uint, email string) UsersClient {
	qry := "SELECT c.id, c.name" +
		" FROM users_clients AS u_c" +
		" JOIN users AS u ON u.id = u_c.user_id" +
		" JOIN clients AS c ON c.id = u_c.client_id" +
		" WHERE c.id = ?" +
		" AND u.email = ?" +
		" AND c.enabled = 1"
	var uc UsersClient
	err := Db.QueryRow(qry, clientID, email).Scan(&uc.ID, &uc.Name)
	switch {
	case err == sql.ErrNoRows:
		return uc
	case err != nil:
		log.Printf("query error: %v\n", err)
		return uc
	default:
		// nothing
	}
	return uc
}

// GetClientsForUserEmail - client list for a user (identified by email)
func GetClientIDsForUserEmail(email string) []uint64 {
	qry := "SELECT c.id" +
		" FROM users_clients AS u_c" +
		" JOIN users AS u ON u.id = u_c.user_id" +
		" JOIN clients AS c ON c.id = u_c.client_id" +
		" WHERE u.email = ?" +
		" AND c.enabled = 1"
	rows, err := Db.Query(qry, email)
	var ids []uint64
	if err != nil {
		log.Println("Error retrieving client ids for email:", email, "from database. Error:", err)
		return ids
	}
	defer rows.Close()
	for rows.Next() {
		var id uint64
		if err1 := rows.Scan(&id); err1 != nil {
			log.Println("Error retrieving client ids for email", email, "from database whilst reading returned results. Error: ", err1)
			return ids
		}
		if id == 0 {
			log.Printf("client id not found for email %s\n", email)
			continue
		}
		ids = append(ids, id)
	}
	return ids
}

// ClientStats - get some stats using the stats table
func ClientStats(statsRequest StatsRequest, clientID int) ([]StatsNewResult, error) {
	var outerQryStart = "SELECT a.client_id, a.client_name, a.sent, a.requested, a.group_period FROM ("
	var outerQryEnd = ") AS a WHERE a.client_id = ?"
	var qry = "SELECT g.client_id AS client_id, c.name AS client_name," +
		" sum(g.sent_count) AS sent,  sum(g.requested_count) AS requested, "
	const qryEnd = " AS group_period" +
		" FROM google_reviews_stats AS g" +
		" LEFT JOIN clients AS c ON g.client_id = c.id" +
		" WHERE g.stats_date BETWEEN ? AND ?" +
		" GROUP BY client_id, group_period"
	timeGrouping := "DATE(g.stats_date)"
	switch statsRequest.TimeGrouping {
	case "Week":
		// timeGrouping = "WEEK(g.stats_date)"
		timeGrouping = "DATE_FORMAT(g.stats_date, '%Y-%u-01T00:00:00Z')"
	case "Month":
		timeGrouping = "DATE_FORMAT(g.stats_date, '%Y-%m-01T00:00:00Z')"
	case "Year":
		timeGrouping = "DATE_FORMAT(g.stats_date, '%Y-01-01T00:00:00Z')"
	}
	rows, err := Db.Query(outerQryStart+qry+timeGrouping+qryEnd+outerQryEnd, statsRequest.StartDay, statsRequest.EndDay, clientID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	var statsResult []StatsNewResult
	for rows.Next() {
		var s StatsNewResult
		if err := rows.Scan(&s.ClientID, &s.ClientName, &s.Sent, &s.Requested, &s.GroupPeriod); err != nil {
			log.Printf("Error getting stats results: %v\n", err)
		}
		statsResult = append(statsResult, s)
	}
	return statsResult, nil
}

// GoogleMyBusinessLocationName - get the config required fields for a reviews reporting from the google my business location name and postal code address
func ConfigFromGoogleMyBusinessLocationNameAndPostalCode(googleMyBusinessLocationName string, googleMyBusinessPostalCode string) GoogleReviewsConfigAndGoogleMyBusinessLocation {
	qry := "SELECT config.google_my_business_location_name," +
		" config.google_my_business_postal_code," +
		// " config.google_my_business_report_enabled," +
		// " config.email_address," +
		" config.time_zone," +
		" config.client_id" +
		" FROM google_reviews_configs AS config" +
		" JOIN clients AS client ON client.id = config.client_id" +
		" WHERE config.enabled = 1" +
		" AND TRIM(config.google_my_business_location_name) = TRIM(?)" +
		" AND REPLACE(config.google_my_business_postal_code, ' ', '') = REPLACE(?, ' ', '')" +
		" AND client.enabled = 1" +
		// TODO: is this necessary
		" AND config.google_my_business_report_enabled = 1"
	rows, err := Db.Query(qry, googleMyBusinessLocationName, googleMyBusinessPostalCode)
	var grcfgmbln GoogleReviewsConfigAndGoogleMyBusinessLocation
	if err != nil {
		log.Println("Error retrieving google my business location name", googleMyBusinessLocationName, "from database. Error: ", err)
		return grcfgmbln
	}
	defer rows.Close()
	for rows.Next() {
		if err1 := rows.Scan(&grcfgmbln.GoogleMyBusinessLocationName,
			&grcfgmbln.GoogleMyBusinessPostalCode,
			// &grcfgmbln.GoogleMyBusinessReportEnabled,
			// &grcfgmbln.EmailAddress,
			&grcfgmbln.TimeZone,
			&grcfgmbln.ClientID); err1 != nil {
			log.Println("Error retrieving google my business location name", googleMyBusinessLocationName, "from database whilst reading returned results. Error: ", err1)
			return grcfgmbln
		}
	}
	return grcfgmbln
}

// GetClientReportsByClientID retrieves all reports for a specific client
func GetClientReportsByClientID(clientID int) ([]*shared_templates.ClientReportData, error) {
	query := `
		SELECT report_id, client_id, report_period_start, report_period_end, 
			   generated_at, locations
		FROM client_reports
		WHERE client_id = ?
		ORDER BY report_period_start DESC
	`

	rows, err := Db.Query(query, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to query client reports: %w", err)
	}
	defer rows.Close()

	// Get client name once (assuming a clients table with name column)
	var clientName string
	err = Db.QueryRow("SELECT name FROM clients WHERE id = ?", clientID).Scan(&clientName)
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

		// Safely unmarshal the location results JSON with NULL handling
		locationResults, err := validateLocationResults(locationResultsJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to validate location results: %w", err)
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

// GetClientReportByIDForClient retrieves a specific report by its ID with client permission check
func GetClientReportByIDForClient(reportID int64, clientID int) (*shared_templates.ClientReportData, error) {
	query := `
		SELECT report_id, client_id, report_period_start, report_period_end, 
			   generated_at, locations
		FROM client_reports
		WHERE report_id = ? AND client_id = ?
	`

	var reportID2 int64
	var clientID2 int
	var periodStart, periodEnd, generatedAt time.Time
	var locationResultsJSON []byte

	err := Db.QueryRow(query, reportID, clientID).Scan(
		&reportID2,
		&clientID2,
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
	err = Db.QueryRow("SELECT name FROM clients WHERE id = ?", clientID).Scan(&clientName)
	if err != nil {
		clientName = fmt.Sprintf("Client %d", clientID) // Fallback if name not found
	}

	// Safely unmarshal the location results JSON with NULL handling
	locationResults, err := validateLocationResults(locationResultsJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to validate location results: %w", err)
	}

	// Create and return the client report data
	return &shared_templates.ClientReportData{
		ReportID:        reportID2,
		ClientID:        clientID2,
		ClientName:      clientName,
		PeriodStart:     periodStart,
		PeriodEnd:       periodEnd,
		GeneratedAt:     generatedAt,
		LocationResults: locationResults,
	}, nil
}
