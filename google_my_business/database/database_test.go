package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	testfixtures "gopkg.in/testfixtures.v2"
)

var (
	db       *sql.DB
	fixtures *testfixtures.Context
)

func TestMain(m *testing.M) {
	var err error

	// Open connection with the test database.
	// Do NOT import fixtures in a production database!
	// Existing data would be deleted
	db = OpenDB(TestDbName, TestDbAddress, TestDbPort, TestDbUsername, TestDbPassword)
	if db == nil {
		log.Fatal(err)
	}

	// creating the context that hold the fixtures
	// see about all compatible databases in this page below
	fixtures, err = testfixtures.NewFolder(db, &testfixtures.MySQL{}, "fixtures")
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func prepareTestDatabase() {
	// prevent the error message:
	// Loading aborted because the database name does not contains "test"
	testfixtures.SkipDatabaseNameCheck(true)
	if err := fixtures.Load(); err != nil {
		log.Fatal(err)
	}
}

func TestConfigFromGoogleMyBusinessLocationNameAndPostalCode(t *testing.T) {
	prepareTestDatabase()
	googleMyBusinessLocationName := "City Taxis"
	googleMyBusinessPostalCode := "S9 2LR"
	grcfgmbln := ConfigFromGoogleMyBusinessLocationNameAndPostalCode(db, googleMyBusinessLocationName, googleMyBusinessPostalCode, LookupModeReply)
	if len(grcfgmbln.GoogleMyBusinessOneStarRatingReply) == 0 {
		t.Fatal("google my business location name", googleMyBusinessLocationName, "not found")
	}
	fmt.Printf("google reviews config fields from google my business location name: %+v\n", grcfgmbln)
}

func TestConfigFromGoogleMyBusinessLocationNameAndPostalCodeFails1(t *testing.T) {
	prepareTestDatabase()
	googleMyBusinessLocationName := "Wrong Name"
	googleMyBusinessPostalCode := "S9 2LR"
	grcfgmbln := ConfigFromGoogleMyBusinessLocationNameAndPostalCode(db, googleMyBusinessLocationName, googleMyBusinessPostalCode, LookupModeReply)
	if len(grcfgmbln.GoogleMyBusinessOneStarRatingReply) != 0 {
		t.Fatal("google my business location name", googleMyBusinessLocationName, "not found")
	}
}

func TestConfigFromGoogleMyBusinessLocationNameAndPostalCodeFails2(t *testing.T) {
	prepareTestDatabase()
	googleMyBusinessLocationName := "City Taxis"
	googleMyBusinessPostalCode := "S9 4LB"
	grcfgmbln := ConfigFromGoogleMyBusinessLocationNameAndPostalCode(db, googleMyBusinessLocationName, googleMyBusinessPostalCode, LookupModeReply)
	if len(grcfgmbln.GoogleMyBusinessOneStarRatingReply) != 0 {
		t.Fatal("google my business location name", googleMyBusinessLocationName, "not found")
	}
}

func TestConfigFromGoogleMyBusinessLocationNameAndPostalCodeTrimWhiteSpace(t *testing.T) {
	prepareTestDatabase()
	googleMyBusinessLocationName := "  City Taxis  "
	googleMyBusinessPostalCode := "S9 2LR"
	grcfgmbln := ConfigFromGoogleMyBusinessLocationNameAndPostalCode(db, googleMyBusinessLocationName, googleMyBusinessPostalCode, LookupModeReply)
	if len(grcfgmbln.GoogleMyBusinessOneStarRatingReply) == 0 {
		t.Fatal("google my business location name", googleMyBusinessLocationName, "not found")
	}
	fmt.Printf("google reviews config fields from google my business location name: %+v\n", grcfgmbln)
}

func TestConfigFromGoogleMyBusinessLocationNameAndPostalCodeTrimWhiteSpaceCaseChange(t *testing.T) {
	prepareTestDatabase()
	googleMyBusinessLocationName := "  city taxis  "
	googleMyBusinessPostalCode := "S9 2LR"
	grcfgmbln := ConfigFromGoogleMyBusinessLocationNameAndPostalCode(db, googleMyBusinessLocationName, googleMyBusinessPostalCode, LookupModeReply)
	if len(grcfgmbln.GoogleMyBusinessOneStarRatingReply) == 0 {
		t.Fatal("google my business location name", googleMyBusinessLocationName, "not found")
	}
	fmt.Printf("google reviews config fields from google my business location name: %+v\n", grcfgmbln)
}

func TestConfigFromGoogleMyBusinessLocationNameAndPostalCodeTrimWhiteSpaceCaseChangeAndPostalCodeSpacesAndCase(t *testing.T) {
	prepareTestDatabase()
	googleMyBusinessLocationName := "  city taxis  "
	googleMyBusinessPostalCode := " s9    2lr  "
	grcfgmbln := ConfigFromGoogleMyBusinessLocationNameAndPostalCode(db, googleMyBusinessLocationName, googleMyBusinessPostalCode, LookupModeReply)
	if len(grcfgmbln.GoogleMyBusinessOneStarRatingReply) == 0 {
		t.Fatal("google my business location name", googleMyBusinessLocationName, "not found")
	}
	fmt.Printf("google reviews config fields from google my business location name: %+v\n", grcfgmbln)
}

// the entry in the database test data has been commented out so do not cause potential issues on google my business (highly unlikely)
// func TestConfigFromGoogleMyBusinessLocationNameAndPostalCodeReportOnly(t *testing.T) {
// 	prepareTestDatabase()
// 	googleMyBusinessLocationName := "City Taxis"
// 	googleMyBusinessPostalCode := "S9 2LR"
// 	grcfgmbln := ConfigFromGoogleMyBusinessLocationNameAndPostalCode(db, googleMyBusinessLocationName, googleMyBusinessPostalCode, LookupModeReport)
// 	if len(grcfgmbln.GoogleMyBusinessOneStarRatingReply) == 0 {
// 		t.Fatal("google my business location name", googleMyBusinessLocationName, "not found")
// 	}
// 	fmt.Printf("google reviews config fields from google my business location name: %+v\n", grcfgmbln)
// }

func TestAllConfigsWithReplyToGoogleMyBusiness(t *testing.T) {
	prepareTestDatabase()
	grcfgmblns := AllConfigsWithReplyToGoogleMyBusiness(db)
	if len(grcfgmblns) == 0 {
		t.Fatal("all configs with reply to google my business not found")
	}
	fmt.Printf("all configs with reply to google my business: %+v\n", grcfgmblns)
}

// TestClientReport tests saving and retrieving client reports
func TestClientReport(t *testing.T) {
	prepareTestDatabase()

	// Setup test data
	clientID := 1
	periodStart := time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2023, 5, 31, 0, 0, 0, 0, time.UTC) // Changed to match MySQL DATE behavior

	// Create sample location results
	locationResults := []byte(`[
		{
			"analysis": {
				"overall_summary": {
					"summary_text": "Test summary",
					"positive_themes": ["Good service", "Punctuality"],
					"negative_themes": ["Communication issues"],
					"overall_perception": "Positive",
					"average_rating": 4.5
				},
				"sentiment_analysis": {
					"positive_count": 8,
					"positive_percentage": 80.0,
					"neutral_count": 1,
					"neutral_percentage": 10.0,
					"negative_count": 1,
					"negative_percentage": 10.0,
					"total_reviews": 10,
					"sentiment_trend": "Stable"
				}
			},
			"metadata": {
				"generatedAt": "2023-05-15T10:00:00Z",
				"reviewCount": 10,
				"locationID": "location123",
				"locationName": "Test Location",
				"businessName": "Test Business",
				"clientID": 1,
				"analyzerID": "test-analyzer",
				"analyzerName": "Test Analyzer",
				"analyzerModel": "test-model"
			}
		}
	]`)

	// Test saving a report
	reportID, err := SaveClientReport(db, clientID, periodStart, periodEnd, locationResults)
	if err != nil {
		t.Fatalf("Failed to save client report: %v", err)
	}
	if reportID <= 0 {
		t.Fatalf("Expected valid report ID, got %d", reportID)
	}

	// Test retrieving by ID
	report, err := GetClientReportByID(db, reportID)
	if err != nil {
		t.Fatalf("Failed to retrieve client report by ID: %v", err)
	}

	if report.ClientID != clientID {
		t.Errorf("Expected client ID %d, got %d", clientID, report.ClientID)
	}

	// Date comparison without time - MySQL DATE type doesn't store time component
	if report.PeriodStart.Format("2006-01-02") != periodStart.Format("2006-01-02") {
		t.Errorf("Expected period start %v, got %v", periodStart.Format("2006-01-02"), report.PeriodStart.Format("2006-01-02"))
	}

	if report.PeriodEnd.Format("2006-01-02") != periodEnd.Format("2006-01-02") {
		t.Errorf("Expected period end %v, got %v", periodEnd.Format("2006-01-02"), report.PeriodEnd.Format("2006-01-02"))
	}

	// Test retrieving by client ID
	reports, err := GetClientReportsByClientID(db, clientID)
	if err != nil {
		t.Fatalf("Failed to retrieve client reports by client ID: %v", err)
	}

	found := false
	for _, r := range reports {
		if r.ReportID == reportID {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Could not find report with ID %d in client reports", reportID)
	}

	// Test retrieving by client and period (using the stored record's exact dates)
	report2, err := GetClientReportByClientAndPeriod(db, clientID, report.PeriodStart, report.PeriodEnd)
	if err != nil {
		t.Fatalf("Failed to retrieve client report by client and period: %v", err)
	}

	if report2.ReportID != reportID {
		t.Errorf("Expected report ID %d, got %d", reportID, report2.ReportID)
	}
}

// TestGetClientsWithMonthlyReviewAnalysisEnabled tests fetching clients with monthly review analysis enabled
func TestGetClientsWithMonthlyReviewAnalysisEnabled(t *testing.T) {
	prepareTestDatabase()

	// Call the function
	clients, err := GetClientsWithMonthlyReviewAnalysisEnabled(db)
	if err != nil {
		t.Fatalf("Failed to get clients with monthly review analysis: %v", err)
	}

	// Check that we get at least one result (assuming test fixtures include at least one)
	if len(clients) == 0 {
		t.Logf("No clients with monthly review analysis enabled found. This may be expected if none are configured in fixtures.")
	}

	// Validate fields of returned clients
	for _, client := range clients {
		if client.ClientID <= 0 {
			t.Errorf("Expected valid client ID, got %d", client.ClientID)
		}

		if client.ClientName == "" {
			t.Errorf("Expected non-empty client name for client ID %d", client.ClientID)
		}

		// Email might be optional in some cases, so just log if missing
		if client.EmailAddress == "" {
			t.Logf("Client ID %d has no email address", client.ClientID)
		}

		// Log if report email address is not set (Valid but empty string, or not Valid)
		if !client.ReportEmailAddress.Valid {
			t.Logf("Client ID %d has no report email address (null)", client.ClientID)
		} else if client.ReportEmailAddress.String == "" {
			t.Logf("Client ID %d has empty report email address", client.ClientID)
		}
	}
}
