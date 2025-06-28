package google_my_business_api

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"rm_client_portal/database"
	"testing"

	"golang.org/x/oauth2/google"
	"gopkg.in/testfixtures.v2"
)

var fixtures *testfixtures.Context

func TestMain(m *testing.M) {
	var err error

	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/business.manage")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	SetClient(config)
	// fmt.Println(Client)

	// Open connection with the test database.
	// Do NOT import fixtures in a production database!
	// Existing data would be deleted
	database.OpenDB(TestDbName, TestDbAddress, TestDbPort, TestDbUsername, TestDbPassword)

	// creating the context that hold the fixtures
	// see about all compatible databases in this page below
	fixtures, err = testfixtures.NewFolder(database.Db, &testfixtures.MySQL{}, "../database/fixtures")
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

func TestGetAccounts(t *testing.T) {
	accounts := GetAccounts()
	fmt.Println(accounts)
}

func TestGetLocations(t *testing.T) {
	prepareTestDatabase()
	accounts := GetAccounts()
	clientIDs := []uint64{4}
	for _, a := range accounts {
		l := GetLocationsCheckClientID(a, clientIDs)
		fmt.Printf("%+v\n", l)
		// for _, m := range l {
		// 	fmt.Println(m.GoogleMyBusinessPostalCode)
		// }
	}
}

func TestGetLocationsJson(t *testing.T) {
	prepareTestDatabase()
	accounts := GetAccounts()
	clientIDs := []uint64{4}
	var locs []database.GoogleReviewsConfigAndGoogleMyBusinessLocation
	for _, a := range accounts {
		l := GetLocationsCheckClientID(a, clientIDs)
		locs = append(locs, l...)
		fmt.Printf("%+v\n", l)
		// for _, m := range l {
		// 	fmt.Println(m.GoogleMyBusinessPostalCode)
		// }
	}
	j, err := json.Marshal(locs)
	if err != nil {
		t.Fatalf("Error marshalling locations: %+v, to JSON, err: %s", locs, err)
	}
	// fmt.Println(j)
	fmt.Printf("%s\n", j)
}

func TestReportOnReviews(t *testing.T) {
	prepareTestDatabase()
	accounts := GetAccounts()
	clientIDs := []uint64{4}
	for _, a := range accounts {
		l := GetLocationsCheckClientID(a, clientIDs)
		// fmt.Println(l)
		for _, g := range l {
			reviewRatings := ReportOnReviewsWeb(g, "2024-09-01T00:00:00Z", "2024-09-30T00:00:00Z")
			fmt.Printf("%d,%s:%+v\n", g.ClientID, g.GoogleMyBusinessLocationName, reviewRatings)
		}
	}
}

func TestInsights(t *testing.T) {
	prepareTestDatabase()
	accounts := GetAccounts()
	clientIDs := []uint64{4}
	for _, a := range accounts {
		l := GetLocationsCheckClientID(a, clientIDs)
		// fmt.Println(l)
		for _, g := range l {
			// insights
			insights := ReportOnInsightsWeb(g, "2024-09-01T00:00:00Z", "2024-09-30T00:00:00Z")
			fmt.Printf("%d,%s:%+v\n", g.ClientID, g.GoogleMyBusinessLocationName, insights)
		}
	}
}

func TestReportOnReviewsAndInsights(t *testing.T) {
	prepareTestDatabase()
	accounts := GetAccounts()
	clientIDs := []uint64{4}
	var locs []database.GoogleReviewsConfigAndGoogleMyBusinessLocation
	for _, a := range accounts {
		l := GetLocationsCheckClientID(a, clientIDs)
		// fmt.Println(l)
		startTime := "2024-09-01T00:00:00Z"
		endTime := "2024-09-30T00:00:00Z"
		for _, g := range l {
			g.GoogleReviewRatings = ReportOnReviewsWeb(g, startTime, endTime)
			// insights
			g.GoogleInsights = ReportOnInsightsWeb(g, startTime, endTime)
			locs = append(locs, g)
		}
	}
	j, err := json.Marshal(locs)
	if err != nil {
		log.Printf("Error marshalling locations: %+v, to JSON, err: %s", locs, err)
	}
	// fmt.Println(j)
	fmt.Printf("%s\n", j)
}

func TestReportTime(t *testing.T) {
	tm := ReportTime("2024-11-12T08:42:00Z", "Europe/London")
	fmt.Println(tm)
	if tm.Hour() != 8 {
		t.Fatal("Error setting report time")
	}
}
