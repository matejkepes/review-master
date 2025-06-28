package database

import (
	"encoding/json"
	"fmt"
	"google_reviews_autocab/autocab_api_v1"
	"log"
	"net/url"
	"os"
	"testing"
	"time"

	testfixtures "gopkg.in/testfixtures.v2"
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
	// set review master SMS gateway master queue ID
	SetReviewMasterSMSGatewayMasterQueueID()
	if ReviewMasterSMSMasterQueue == 0 {
		log.Fatal("Review Master SMS Gateway Master Queue ID not found")
	}
}

func TestConfigFromTokenWithChecks(t *testing.T) {
	prepareTestDatabase()
	token := "QxrH0iJc3wv/lj/YKVppNYRad7tN0Z3x"
	grcftwc := ConfigFromTokenWithChecks(token, false)
	if grcftwc.ClientID == 0 {
		t.Fatal("token ", token, " not found")
	}
	fmt.Printf("google reviews config with checks: %+v\n", grcftwc)
}

func TestConfigFromTokenWithChecksWrongToken(t *testing.T) {
	prepareTestDatabase()
	token := "QxrH0iJc3wv/lj/YKVppNYRad7tN0123"
	grcftwc := ConfigFromTokenWithChecks(token, false)
	if grcftwc.Country != "" {
		t.Fatal("token ", token, " was found and it should not have been")
	}
}

func TestConfigFromTokenWithChecksMultiMessage(t *testing.T) {
	prepareTestDatabase()
	token := "OYBpBsZ9OhR-nbsupMQU_hCTJuaoqtZw"
	grcftwc := ConfigFromTokenWithChecks(token, false)
	if grcftwc.ClientID == 0 {
		t.Fatal("token ", token, " not found")
	}
	fmt.Printf("google reviews config with checks: %+v\n", grcftwc)
}

func TestConfigFromTokenWithChecksSendFromIcabbiApp(t *testing.T) {
	prepareTestDatabase()
	token := "OYBpBsZ9OhR-nbsupMQU_hCTJuaabtZ1"
	grcftwc := ConfigFromTokenWithChecks(token, false)
	if grcftwc.ClientID == 0 {
		t.Fatal("token ", token, " not found")
	}
	fmt.Printf("google reviews config with checks: %+v\n", grcftwc)
}

func TestConfigFromTokenWithChecksReplaceTelephoneCountryCode(t *testing.T) {
	prepareTestDatabase()
	token := "OYBpBsZ9OhR-nbsupMQU_hCTJuaabtZ4"
	grcftwc := ConfigFromTokenWithChecks(token, false)
	if grcftwc.ClientID == 0 {
		t.Fatal("token ", token, " not found")
	}
	fmt.Printf("google reviews config with checks: %+v\n", grcftwc)
	if !grcftwc.ReplaceTelephoneCountryCode {
		t.Fatal("replace telephone country code should be true")
	}
}

func TestConfigFromTokenWithChecksWeekdaysAllFalse(t *testing.T) {
	prepareTestDatabase()
	token := "OYBpBsZ9OhR-nbsupMQU_hCTJuaabtZ5"
	grcftwc := ConfigFromTokenWithChecks(token, false)
	if grcftwc.ClientID != 0 {
		t.Fatal("clientID should be set to 0 and not returned")
	}
}

func TestConfigFromTokenWithChecksDrive(t *testing.T) {
	prepareTestDatabase()
	token := "gDVTnJ7peYBX8WvIsomu9dkBPULR1rcwZrURQNIqdce3kRIXPQLiyj8IILuXZmhN"
	grcftwc := ConfigFromTokenWithChecks(token, false)
	if grcftwc.ClientID == 0 {
		t.Fatal("token ", token, " not found")
	}
	fmt.Printf("google reviews config with checks: %+v\n", grcftwc)
}

func TestConfigFromTokenWithChecksWeekdaysAllFalseAndIgnoreTimeAndSentCountChecks(t *testing.T) {
	prepareTestDatabase()
	token := "OYBpBsZ9OhR-nbsupMQU_hCTJuaabtZ5"
	grcftwc := ConfigFromTokenWithChecks(token, true)
	if grcftwc.ClientID == 0 {
		t.Fatal("token ", token, " not found")
	}
	fmt.Printf("google reviews config with checks: %+v\n", grcftwc)
}

func TestGetAutocabConfigsWithChecks(t *testing.T) {
	prepareTestDatabase()
	grcftwcs := GetAutocabConfigsWithChecks(false)
	if len(grcftwcs) == 0 {
		t.Fatal("There should be at least one Autocab config found")
	}
	fmt.Printf("Autocab google reviews configs with checks: %+v\n", grcftwcs)
}

func TestConfigFromTokenWithChecksReviewsMasterQueue(t *testing.T) {
	prepareTestDatabase()
	token := "OYBpBsZ9OhR-nbsupMQU_hCTJuaabtZ1gsfsfp123"
	grcftwc := ConfigFromTokenWithChecks(token, false)
	if grcftwc.ClientID == 0 {
		t.Fatal("token ", token, " not found")
	}
	fmt.Printf("google reviews config with checks: %+v\n", grcftwc)
	if !grcftwc.ReviewMasterSMSGatewayUseMasterQueue {
		t.Fatal("Review Master SMS Gateway Use Master Queue should be true")
	}
}

func TestLastSentFromTelephoneAndClient(t *testing.T) {
	prepareTestDatabase()
	var clientID uint64 = 1
	telephone := "447123456789"
	lastSent, sentCount, stop, found := LastSentFromTelephoneAndClient(telephone, clientID)
	fmt.Println("lastSent: ", lastSent, " sentCount: ", sentCount, " stop: ", stop, " found: ", found)
	// check less than 3 years ago
	if !lastSent.After(time.Now().AddDate(-3, 0, 0)) {
		t.Fatal("last sent not found")
	}
	// check found record
	if !found {
		t.Fatal("record not found")
	}
}

func TestLastSentFromTelephoneAndClientFailsToFind(t *testing.T) {
	prepareTestDatabase()
	var clientID uint64 = 1
	telephone := "447123456780"
	lastSent, sentCount, stop, found := LastSentFromTelephoneAndClient(telephone, clientID)
	fmt.Println("lastSent: ", lastSent, " sentCount: ", sentCount, " stop: ", stop, " found: ", found)
	// check less than 3 years ago
	if lastSent.After(time.Now().AddDate(-3, 0, 0)) {
		t.Fatal("last sent not found")
	}
	// check found record
	if found {
		t.Fatal("record found")
	}
}

func TestLastSentFromTelephoneAndClientStop(t *testing.T) {
	prepareTestDatabase()
	var clientID uint64 = 1
	telephone := "447123456788"
	lastSent, sentCount, stop, found := LastSentFromTelephoneAndClient(telephone, clientID)
	fmt.Println("lastSent: ", lastSent, " sentCount: ", sentCount, " stop: ", stop, " found: ", found)
	// check less than 3 years ago
	if !lastSent.After(time.Now().AddDate(-3, 0, 0)) {
		t.Fatal("last sent not found")
	}
	// check stop is true
	if !stop {
		t.Fatal("stop should be set to true")
	}
	// check found record
	if !found {
		t.Fatal("record not found")
	}
}

func TestUpdateLastSentUpdate(t *testing.T) {
	prepareTestDatabase()
	var clientID uint64 = 1
	telephone := "447123456789"
	lastSent, sentCount, stop, found := LastSentFromTelephoneAndClient(telephone, clientID)
	fmt.Println("lastSent: ", lastSent, " sentCount: ", sentCount, " stop: ", stop, " found: ", found)
	// check less than 3 years ago
	if !lastSent.After(time.Now().AddDate(-3, 0, 0)) {
		t.Fatal("last sent not found")
	}
	// check found record
	if !found {
		t.Fatal("record not found")
	}
	UpdateLastSent(telephone, clientID, sentCount+1)
	lastSent1, sentCount1, stop1, found1 := LastSentFromTelephoneAndClient(telephone, clientID)
	if !found1 {
		t.Fatal("record not found after update")
	}
	if stop != stop1 {
		t.Fatal("stop should not have changed after update")
	}
	if sentCount+1 != sentCount1 {
		t.Fatal("sent count should be one more after update")
	}
	if !lastSent1.After(time.Now().Add(-time.Minute * 1)) {
		t.Fatal("last sent has not updated after update")
	}
}

func TestUpdateLastSentNew(t *testing.T) {
	prepareTestDatabase()
	var clientID uint64 = 1
	telephone := "447123456787"
	lastSent, sentCount, stop, found := LastSentFromTelephoneAndClient(telephone, clientID)
	fmt.Println("lastSent: ", lastSent, " sentCount: ", sentCount, " stop: ", stop, " found: ", found)
	// check less than 3 years ago
	if lastSent.After(time.Now().AddDate(-3, 0, 0)) {
		t.Fatal("last sent not found")
	}
	// check found record
	if found {
		t.Fatal("record found")
	}
	UpdateLastSent(telephone, clientID, sentCount+1)
	lastSent1, sentCount1, stop1, found1 := LastSentFromTelephoneAndClient(telephone, clientID)
	if !found1 {
		t.Fatal("record not found after update")
	}
	if stop != stop1 {
		t.Fatal("stop should not have changed after update")
	}
	if sentCount+1 != sentCount1 {
		t.Fatal("sent count should be one more after update")
	}
	if !lastSent1.After(time.Now().Add(-time.Minute * 1)) {
		t.Fatal("last sent has not updated after update")
	}
}

func TestStopSending(t *testing.T) {
	prepareTestDatabase()
	var clientID uint64 = 1
	telephone := "447123456789"
	StopSending(telephone, clientID)
	lastSent, sentCount, stop, found := LastSentFromTelephoneAndClient(telephone, clientID)
	fmt.Println("lastSent: ", lastSent, " sentCount: ", sentCount, " stop: ", stop, " found: ", found)
	// check stop is true
	if !stop {
		t.Fatal("stop should be set to true")
	}
}

func TestStopSendingNoRecord(t *testing.T) {
	prepareTestDatabase()
	var clientID uint64 = 1
	telephone := "447123456786"
	StopSending(telephone, clientID)
	lastSent, sentCount, stop, found := LastSentFromTelephoneAndClient(telephone, clientID)
	fmt.Println("lastSent: ", lastSent, " sentCount: ", sentCount, " stop: ", stop, " found: ", found)
	// check found is false
	if found {
		t.Fatal("stop should be set to true")
	}
}

func TestDailySentCount(t *testing.T) {
	prepareTestDatabase()
	count := DailySentCount(1)
	fmt.Println("count: ", count)
	if count != 1 {
		t.Fatal("daily sent count incorrect should be 1")
	}
}

func TestAddSendLater1(t *testing.T) {
	prepareTestDatabase()
	var clientID uint64 = 1
	telephone := "447123456787"

	body, _ := json.Marshal(map[string]string{
		"queue_id":  "81",
		"telephone": "+" + telephone,
		"message":   "testing",
	})
	// fmt.Println("body:", string(body))

	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
		"Api-Token":    "GGK8dkags0EYe0r0UuPowQBV79DeUJE/Lu6190iyyG5PZ+3v8c8Bs8g",
	}

	params1 := url.Values{}
	params1.Add("t", telephone)
	params1.Add("m", "some message")

	AddSendLater(telephone, clientID, 5,
		"https://api.messagemedia.com/v1/messages", "POST", "12348GYxCGv5abcES7GF", "2HGFUd9i987Gv6018D1234cNhHDRONH",
		headers, params1, body, true, false, false, "", false, "success=1", 20)
}

func TestAddSendLater2(t *testing.T) {
	prepareTestDatabase()
	var clientID uint64 = 4
	telephone := testTelephone
	message := "testing"

	headers := map[string]string{
		"Content-Type":              "application/json",
		"Accept":                    "application/json",
		"Ocp-Apim-Subscription-Key": testAutocabKeySendSMS,
	}

	body, _ := json.Marshal(autocab_api_v1.SendSMSRequest{
		SenderName:      testAutocabSendSMSSenderName,
		RecipientTelnos: []string{telephone},
		Message:         message,
	})

	AddSendLater(telephone, clientID, 5,
		"https://autocab-api.azure-api.net/sms/v1/send", "POST", "", "",
		headers, nil, body, false, false, true, "AUTOCAB_V1", false, "ok", 20)
}

func TestUpdateStatsWithCountsSent(t *testing.T) {
	prepareTestDatabase()
	var clientID uint64 = 1

	qry := "SELECT sent_count, requested_count FROM google_reviews_stats WHERE client_id = ? AND stats_date = CURDATE()"
	var sentCount1 int
	var requestedCount1 int
	row := Db.QueryRow(qry, clientID)
	if err := row.Scan(&sentCount1, &requestedCount1); err != nil {
		fmt.Println("Error retrieving stats for clientID", clientID, "from database. Error: ", err)
	}

	// sent
	UpdateStatsWithCounts(clientID, 1, 1)

	qry = "SELECT sent_count, requested_count FROM google_reviews_stats WHERE client_id = ? AND stats_date = CURDATE()"
	var sentCount2 int
	var requestedCount2 int
	row = Db.QueryRow(qry, clientID)
	if err := row.Scan(&sentCount2, &requestedCount2); err != nil {
		t.Fatal("Error retrieving stats for clientID", clientID, "from database. Error: ", err)
	}

	if sentCount1+1 != sentCount2 {
		t.Fatal("sent count should have been incremented for stats")
	}
	if requestedCount1+1 != requestedCount2 {
		t.Fatal("requested count should have been incremented for stats")
	}
}

func TestUpdateStatsWithCountsNotSent(t *testing.T) {
	prepareTestDatabase()
	var clientID uint64 = 1

	qry := "SELECT sent_count, requested_count FROM google_reviews_stats WHERE client_id = ? AND stats_date = CURDATE()"
	var sentCount1 int
	var requestedCount1 int
	row := Db.QueryRow(qry, clientID)
	if err := row.Scan(&sentCount1, &requestedCount1); err != nil {
		fmt.Println("Error retrieving stats for clientID", clientID, "from database. Error: ", err)
	}

	// not sent
	UpdateStatsWithCounts(clientID, 0, 1)

	qry = "SELECT sent_count, requested_count FROM google_reviews_stats WHERE client_id = ? AND stats_date = CURDATE()"
	var sentCount2 int
	var requestedCount2 int
	row = Db.QueryRow(qry, clientID)
	if err := row.Scan(&sentCount2, &requestedCount2); err != nil {
		t.Fatal("Error retrieving stats for clientID", clientID, "from database. Error: ", err)
	}

	if sentCount1 != sentCount2 {
		t.Fatal("sent count should not have been incremented for stats")
	}
	if requestedCount1+1 != requestedCount2 {
		t.Fatal("requested count should have been incremented for stats")
	}
}

func TestUpdateStatsWithCountsNotSentClientID0(t *testing.T) {
	prepareTestDatabase()
	var clientID uint64 = 0

	// not sent
	UpdateStatsWithCounts(clientID, 0, 1)

	// should not be an entry in the stats for client_id 0
	qry := "SELECT id FROM google_reviews_stats WHERE client_id = ? AND stats_date = CURDATE()"
	cnt := 0
	row := Db.QueryRow(qry, clientID)
	if err := row.Scan(&cnt); err == nil {
		t.Fatal("Error there should be no results for stats for clientID", clientID, "from database. Error: ", err)
	}
}
