package database

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"gopkg.in/testfixtures.v2"
)

var fixtures *testfixtures.Context

func TestMain(m *testing.M) {
	var err error

	// Open connection with the test database.
	// Do NOT import fixtures in a production database!
	// Existing data would be deleted
	OpenDB(TestDbName, TestDbAddress, TestDbPort, TestDbUsername, TestDbPassword)

	// creating the context that hold the fixtures
	// see about all compatible databases in this page below
	fixtures, err = testfixtures.NewFolder(Db, &testfixtures.MySQL{}, "fixtures")
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

func TestGetUser(t *testing.T) {
	prepareTestDatabase()
	e := GetUser("test@testing.com", "test")
	fmt.Printf("%+v\n", e)
	if e == "" {
		t.Error("Failed to find user")
	}
}

func TestGetUserIncorrectPassword(t *testing.T) {
	prepareTestDatabase()
	e := GetUser("test@testing.com", "notright")
	fmt.Printf("%+v\n", e)
	if e != "" {
		t.Error("Should Not have found user")
	}
}

// func TestGetUserAndClients(t *testing.T) {
// 	prepareTestDatabase()
// 	uc := GetUserAndClients("test@testing.com", "test")
// 	fmt.Printf("%+v\n", uc)
// 	if uc.Email == "" {
// 		t.Error("Failed to find user")
// 	}
// }

// func TestGetUserAndClientsIncorrect(t *testing.T) {
// 	prepareTestDatabase()
// 	uc := GetUserAndClients("rubbish@testing.com", "wrong")
// 	fmt.Printf("%+v\n", uc)
// 	if uc.Email != "" {
// 		t.Error("Failed to find user")
// 	}
// }

func TestGetClientsForUserEmail(t *testing.T) {
	prepareTestDatabase()
	ucs := GetClientsForUserEmail("test@testing.com")
	fmt.Printf("%+v\n", ucs)
	if len(ucs) == 0 {
		t.Error("Failed to find user clients")
	}
}

func TestGetClientsForUserEmailIncorrect(t *testing.T) {
	prepareTestDatabase()
	ucs := GetClientsForUserEmail("rubbish@testing.com")
	fmt.Printf("%+v\n", ucs)
	if len(ucs) > 0 {
		t.Error("Failed to find user clients")
	}
}

func TestGetClientCheckUserEmail(t *testing.T) {
	prepareTestDatabase()
	uc := GetClientCheckUserEmail(1, "test@testing.com")
	fmt.Printf("%+v\n", uc)
	if uc.ID != 1 {
		t.Error("Failed to retrieve client")
	}
}

func TestGetClientCheckUserEmailWrongEmail(t *testing.T) {
	prepareTestDatabase()
	uc := GetClientCheckUserEmail(1, "test2@testing.com")
	fmt.Printf("%+v\n", uc)
	if uc.ID != 0 {
		t.Error("Should have failed to retrieve client")
	}
}

func TestGetClientCheckUserID(t *testing.T) {
	prepareTestDatabase()
	uc := GetClientCheckUserEmail(3, "test@testing.com")
	fmt.Printf("%+v\n", uc)
	if uc.ID != 0 {
		t.Error("Should have failed to retrieve client")
	}
}

func TestGetClientIDsForUserEmail(t *testing.T) {
	prepareTestDatabase()
	ids := GetClientIDsForUserEmail("test@testing.com")
	fmt.Printf("%+v\n", ids)
	if len(ids) == 0 {
		t.Error("Failed to find user clients")
	}
}

func TestGetClientIDsForUserEmailIncorrect(t *testing.T) {
	prepareTestDatabase()
	ids := GetClientIDsForUserEmail("rubbish@testing.com")
	fmt.Printf("%+v\n", ids)
	if len(ids) > 0 {
		t.Error("Failed to find user clients")
	}
}

func TestClientStats1(t *testing.T) {
	prepareTestDatabase()

	var statsRequest StatsRequest
	b := time.Now().Add(-time.Hour * 24 * 50)
	b = time.Date(b.Year(), b.Month(), b.Day(), 0, 0, 0, 0, b.Location())
	statsRequest.StartDay = b.String()
	e := time.Now().Add(time.Hour * 24)
	e = time.Date(e.Year(), e.Month(), e.Day(), 0, 0, 0, 0, e.Location())
	statsRequest.EndDay = e.String()
	s, err := ClientStats(statsRequest, 1)
	if err != nil {
		t.Fatal("error getting client stats, err: ", err)
	}
	fmt.Printf("stats: %+v\n", s)
}

func TestClientStats2(t *testing.T) {
	prepareTestDatabase()

	var statsRequest StatsRequest
	b := time.Now().Add(-time.Hour * 24 * 50)
	b = time.Date(b.Year(), b.Month(), b.Day(), 0, 0, 0, 0, b.Location())
	statsRequest.StartDay = b.String()
	e := time.Now().Add(time.Hour * 24)
	e = time.Date(e.Year(), e.Month(), e.Day(), 0, 0, 0, 0, e.Location())
	statsRequest.EndDay = e.String()
	statsRequest.TimeGrouping = "Week"
	s, err := ClientStats(statsRequest, 1)
	if err != nil {
		t.Fatal("error getting stats, err: ", err)
	}
	fmt.Printf("stats: %+v\n", s)
}

func TestClientStats3(t *testing.T) {
	prepareTestDatabase()

	var statsRequest StatsRequest
	b := time.Now().Add(-time.Hour * 24 * 50)
	b = time.Date(b.Year(), b.Month(), b.Day(), 0, 0, 0, 0, b.Location())
	statsRequest.StartDay = b.String()
	e := time.Now().Add(time.Hour * 24)
	e = time.Date(e.Year(), e.Month(), e.Day(), 0, 0, 0, 0, e.Location())
	statsRequest.EndDay = e.String()
	statsRequest.TimeGrouping = "Month"
	s, err := ClientStats(statsRequest, 1)
	if err != nil {
		t.Fatal("error getting stats, err: ", err)
	}
	fmt.Printf("stats: %+v\n", s)
}

func TestClientStats4(t *testing.T) {
	prepareTestDatabase()

	var statsRequest StatsRequest
	b := time.Now().Add(-time.Hour * 24 * 50)
	b = time.Date(b.Year(), b.Month(), b.Day(), 0, 0, 0, 0, b.Location())
	statsRequest.StartDay = b.String()
	e := time.Now().Add(time.Hour * 24)
	e = time.Date(e.Year(), e.Month(), e.Day(), 0, 0, 0, 0, e.Location())
	statsRequest.EndDay = e.String()
	statsRequest.TimeGrouping = "Year"
	s, err := ClientStats(statsRequest, 1)
	if err != nil {
		t.Fatal("error getting stats, err: ", err)
	}
	fmt.Printf("stats: %+v\n", s)
}

func TestConfigFromGoogleMyBusinessLocationNameAndPostalCode(t *testing.T) {
	prepareTestDatabase()
	googleMyBusinessLocationName := "City Taxis"
	googleMyBusinessPostalCode := "S9 2LR"
	grcagmbln := ConfigFromGoogleMyBusinessLocationNameAndPostalCode(googleMyBusinessLocationName, googleMyBusinessPostalCode)
	if grcagmbln.ClientID == 0 {
		t.Fatal("google my business location name", googleMyBusinessLocationName, "not found")
	}
	fmt.Printf("google reviews config fields from google my business location name: %+v\n", grcagmbln)
}

// TestValidateLocationResults tests the NULL and malformed JSON handling
func TestValidateLocationResults(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		expectError bool
		expectEmpty bool
	}{
		{
			name:        "NULL input",
			input:       nil,
			expectError: false,
			expectEmpty: true,
		},
		{
			name:        "Empty byte slice",
			input:       []byte{},
			expectError: false,
			expectEmpty: true,
		},
		{
			name:        "Valid empty JSON array",
			input:       []byte("[]"),
			expectError: false,
			expectEmpty: true,
		},
		{
			name:        "Malformed JSON",
			input:       []byte("{invalid json"),
			expectError: false,
			expectEmpty: true,
		},
		{
			name:        "Valid JSON with analysis result",
			input:       []byte(`[{"analysis":{"overall_summary":{"summary_text":"Test"},"sentiment_analysis":{"total_reviews":1},"key_takeaways":{"strengths":[],"areas_for_improvement":[]},"negative_review_breakdown":{"categories":[],"improvement_recommendations":[]},"training_recommendations":{"for_operators":[],"for_drivers":[]}},"metadata":{"generated_at":"2024-01-01T00:00:00Z","review_count":1,"location_id":"test","location_name":"Test Location","business_name":"Test Business","report_period":{"start_date":"2024-01-01","end_date":"2024-01-31"},"client_id":1,"analyzer_id":"test","analyzer_name":"Test Analyzer","analyzer_model":"gpt-3.5-turbo"}}]`),
			expectError: false,
			expectEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := validateLocationResults(tt.input)
			
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if tt.expectEmpty && len(results) != 0 {
				t.Errorf("Expected empty results but got %d items", len(results))
			}
			if !tt.expectEmpty && len(results) == 0 {
				t.Errorf("Expected non-empty results but got empty slice")
			}
		})
	}
}
