package server

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"google_reviews/config"
	"google_reviews/database"

	testfixtures "gopkg.in/testfixtures.v2"
)

var fixtures *testfixtures.Context

func TestMain(m *testing.M) {
	var err error

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

	// read config file
	config.ReadProperties()

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
	database.SetReviewMasterSMSGatewayMasterQueueID()
	if database.ReviewMasterSMSMasterQueue == 0 {
		log.Fatal("Review Master SMS Gateway Master Queue ID not found")
	}
}

// NOTE: sends a message to the server set up in config could be the live server
func TestGoogleReviewsHandler(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"token=" + url.QueryEscape(testToken) +
			"&u=amber&p=amber&t=" + testTelephone +
			"&m=" + url.QueryEscape("Your AMBER Car is on its way! Tap here to TRACK YOUR CAR in the APP or CALL YOUR DRIVER enroute.") +
			"&h=0&v=456749&j=1234567&d=123654" +
			"&r=" + url.QueryEscape("AB12 CDE") +
			"&gr_token=" + url.QueryEscape("QxrH0iJc3wv/lj/YKVppNYRad7tN0Z3x"))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/googlereviews", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(GoogleReviewsHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"success":"1"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGoogleReviewsHandlerTelephoneIncorrect(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"token=" + url.QueryEscape(testToken) +
			"&u=amber&p=amber&t=0712345678" +
			"&m=" + url.QueryEscape("Testing with incorrect telephone number.") +
			"&gr_token=" + url.QueryEscape("QxrH0iJc3wv/lj/YKVppNYRad7tN0Z3x"))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/googlereviews", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(GoogleReviewsHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"success":"0"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// NOTE: sends a message to the server set up in config could be the live server
func TestGoogleReviewsHandlerEmptyReturnStrSuccess(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"number=" + testTelephone +
			"&msg=" + url.QueryEscape("Your ALPHA Car is on its way!") +
			"&gr_token=" + url.QueryEscape("z4JXJfxtNmPH575qPAnhlV_FBgzShqPR"))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/googlereviews", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(GoogleReviewsHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := ``
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// NOTE: sends a message to the server set up in config could be the live server
func TestGoogleReviewsHandlerMultiMessage(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"token=" + url.QueryEscape(testToken) +
			"&u=amber&p=amber&t=" + testTelephone +
			"&m=" + url.QueryEscape("Message 1SSSSSMessage2SSSSSMessage3SSSSSMessage4") +
			"&h=0&v=456749&j=1234567&d=123654" +
			"&r=" + url.QueryEscape("AB12 CDE") +
			"&gr_token=" + url.QueryEscape("OYBpBsZ9OhR-nbsupMQU_hCTJuaoqtZw"))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/googlereviews", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(GoogleReviewsHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"success":"1"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// NOTE: sends a message to the server set up in config could be the live server
func TestGoogleReviewsHandlerMultiMessageSingleMessage(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"token=" + url.QueryEscape(testToken) +
			"&u=amber&p=amber&t=" + testTelephone +
			"&m=" + url.QueryEscape("Message 1") +
			"&h=0&v=456749&j=1234567&d=123654" +
			"&r=" + url.QueryEscape("AB12 CDE") +
			"&gr_token=" + url.QueryEscape("OYBpBsZ9OhR-nbsupMQU_hCTJuaoqtZw"))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/googlereviews", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(GoogleReviewsHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"success":"1"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// NOTE: sends a message to the server set up in config could be the live server
func TestGoogleReviewsHandlerMultiMessageRealFailure(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"gr_token=" + url.QueryEscape("OYBpBsZ9OhR-nbsupMQU_hCTJuaoqtZw") +
			"&token=" + url.QueryEscape(testToken) +
			// "&token=" + url.QueryEscape("DOkTxeI8SkxO-KRaX2YsHkZ6XJ81ln7_InNTv4p-kjXgMri_KJ1W-wmurgSMBf_s") +
			"&t=" + testTelephone +
			"&m=" + url.QueryEscape("We hope you enjoyed your Drive, we would love it if you left us a review, thank you https://goo.gl/ZEw96RSSSSSWe hope you enjoyed your Drive, we would love it if you left us a review, thank you https://goo.gl/ZEw96RSSSSSWe hope you enjoyed your Drive, we would love it if you left us a review, thank you https://goo.gl/ZEw96RSSSSSWe hope you enjoyed your Drive, we would love it if you left us a review, thank you https://cutt.ly/XGpqLpSSSSSWe hope you enjoyed your Drive, we would love it if you left us a review, thank you https://cutt.ly/4Gpq9QSSSSSWe hope you enjoyed your Drive, we would love it if you left us a review, thank you https://cutt.ly/yGpwQCSSSSSWe hope you enjoyed your Drive, we would love it if you left us a review, thank you https://cutt.ly/eGpwSeSSSSSWe hope you enjoyed your Drive, we would love it if you left us a review, thank you https://cutt.ly/qGpwLYSSSSSWe hope you enjoyed your Drive, we would love it if you left us a review, thank you https://cutt.ly/IGpwBTSSSSSWe hope you enjoyed your Drive, we would love it if you left us a review, thank you https://cutt.ly/UGpw3e"))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/googlereviews", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(GoogleReviewsHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"success":"1"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// NOTE: sends a message to the server set up in config could be the live server
func TestGoogleReviewsHandlerSendFromIcabbiApp(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"t=" + testTelephone +
			"&message=" + url.QueryEscape("Message 1") +
			"&gr_token=" + url.QueryEscape("OYBpBsZ9OhR-nbsupMQU_hCTJuaabtZ1"))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/googlereviews", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(GoogleReviewsHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	// TODO: Get correct response when testing with a live iCabbi app
	expected := "ok"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// NOTE: sends a message to the server set up in config could be the live server
func TestGoogleReviewsHandlerSendFromIcabbiAppMessageFromDB(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"t=" + testTelephone +
			"&message=" + url.QueryEscape("Message 1") +
			"&gr_token=" + url.QueryEscape("OYBpBsZ9OhR-nbsupMQU_hCTJuaabtZ2"))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/googlereviews", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(GoogleReviewsHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	// TODO: Get correct response when testing with a live iCabbi app
	expected := "ok"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// NOTE: sends a message to the server set up in config could be the live server
func TestGoogleReviewsHandlerSendFromIcabbiAppMessageFromDBDoDispatcherCheck(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"t=" + testTelephone +
			"&message=" + url.QueryEscape("Message 1") +
			"&gr_token=" + url.QueryEscape("OYBpBsZ9OhR-nbsupMQU_hCTJuaabtZ3") +
			"&b=" + url.QueryEscape(testIcabbiTripID))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/googlereviews", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(GoogleReviewsHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	// TODO: Get correct response when testing with a live iCabbi app
	expected := "ok"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// NOTE: sends a message to the server set up in config could be the live server
func TestGoogleReviewsHandlerSendFromIcabbiAppMessageFromDBDoDispatcherCheckFails(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	// tripID does not exist so will fail dispatcher test
	reader := strings.NewReader(
		"t=" + testTelephone +
			"&message=" + url.QueryEscape("Message 1") +
			"&gr_token=" + url.QueryEscape("OYBpBsZ9OhR-nbsupMQU_hCTJuaabtZ3") +
			"&b=" + url.QueryEscape("12374975395936756"))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/googlereviews", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(GoogleReviewsHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect i.e. NOT ok.
	expected := "ok"
	if rr.Body.String() == expected {
		t.Errorf("handler returned unexpected body: got %v it should not have been %v",
			rr.Body.String(), expected)
	}
}

// NOTE: sends a message to the server set up in config could be the live server
func TestGoogleReviewsHandlerReplaceTelephoneCountryCode(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"token=" + url.QueryEscape(testToken) +
			"&u=amber&p=amber&t=" + testTelephone +
			"&m=" + url.QueryEscape("Your AMBER Car is on its way! Tap here to TRACK YOUR CAR in the APP or CALL YOUR DRIVER enroute.") +
			"&h=0&v=456749&j=1234567&d=123654" +
			"&r=" + url.QueryEscape("AB12 CDE") +
			"&gr_token=" + url.QueryEscape("OYBpBsZ9OhR-nbsupMQU_hCTJuaabtZ4"))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/googlereviews", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(GoogleReviewsHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"success":"1"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// NOTE: sends a message to the server set up in config could be the live server
func TestGoogleReviewsHandlerSendFromIcabbiAppMultiMessageFromDBDoDispatcherCheck(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"t=" + testTelephone +
			"&message=" + url.QueryEscape("Message 1") +
			"&gr_token=" + url.QueryEscape("OYBpBsZ9OhR-nbsupMQU_hCTJuaabtZ6") +
			"&b=" + url.QueryEscape(testIcabbiTripID))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/googlereviews", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(GoogleReviewsHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	// TODO: Get correct response when testing with a live iCabbi app
	expected := "ok"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// NOTE: sends a message to the server set up in config could be the live server
func TestGoogleReviewsHandlerDriveSetup(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"t=" + testTelephone +
			"&gr_token=" + url.QueryEscape("gDVTnJ7peYBX8WvIsomu9dkBPULR1rcwZrURQNIqdce3kRIXPQLiyj8IILuXZmhN") +
			"&token=" + url.QueryEscape("sKFvGpy_izCTVmP-azFNESUaB-dUw2dycrZh9RB4GczUHCZoB5j9xcgK1tYSJu4f") +
			"&b=" + url.QueryEscape(testIcabbiTripID) +
			"&m=" + url.QueryEscape("We hope you enjoyed your Drive, we would love it if you left us a review, thank you https://goo.gl/ZEw96RSSSSSWe hope you enjoyed your Drive, we would love it if you left us a review, thank you https://goo.gl/ZEw96RSSSSSWe hope you enjoyed your Drive, we would love it if you left us a review, thank you https://goo.gl/ZEw96RSSSSSWe hope you enjoyed your Drive, we would love it if you left us a review, thank you https://cutt.ly/XGpqLpSSSSSWe hope you enjoyed your Drive, we would love it if you left us a review, thank you https://cutt.ly/4Gpq9QSSSSSWe hope you enjoyed your Drive, we would love it if you left us a review, thank you https://cutt.ly/yGpwQCSSSSSWe hope you enjoyed your Drive, we would love it if you left us a review, thank you https://cutt.ly/eGpwSeSSSSSWe hope you enjoyed your Drive, we would love it if you left us a review, thank you https://cutt.ly/qGpwLYSSSSSWe hope you enjoyed your Drive, we would love it if you left us a review, thank you https://cutt.ly/IGpwBTSSSSSWe hope you enjoyed your Drive, we would love it if you left us a review, thank you https://cutt.ly/UGpw3e"))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/googlereviews", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(GoogleReviewsHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	// TODO: Get correct response when testing with a live iCabbi app
	expected := "{\"success\":\"1\"}"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// Update last sent so will not send a message but will get the success response from the database
func TestGoogleReviewsHandlerShouldNotSendDueToLastSent(t *testing.T) {
	prepareTestDatabase()

	// add last sent for test telephone
	grcftwc := database.ConfigFromTokenWithChecks("QxrH0iJc3wv/lj/YKVppNYRad7tN0Z3x", false)
	database.UpdateLastSent(testTelephone, grcftwc.ClientID, 1)

	// create a request to pass to our handler.
	reader := strings.NewReader(
		"token=" + url.QueryEscape(testToken) +
			"&u=amber&p=amber&t=" + testTelephone +
			"&m=" + url.QueryEscape("Your AMBER Car is on its way! Tap here to TRACK YOUR CAR in the APP or CALL YOUR DRIVER enroute.") +
			"&h=0&v=456749&j=1234567&d=123654" +
			"&r=" + url.QueryEscape("AB12 CDE") +
			"&gr_token=" + url.QueryEscape("QxrH0iJc3wv/lj/YKVppNYRad7tN0Z3x"))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/googlereviews", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(GoogleReviewsHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	// expected := `{"success":"1"}`
	expected := string(grcftwc.SendSuccessResponse)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// NOTE: sends a message to the server set up in config could be the live server
// Ignore telephone checks
func TestGoogleReviewsHandlerIgnoreTelephoneChecks(t *testing.T) {
	prepareTestDatabase()

	// add last sent for test telephone
	grcftwc := database.ConfigFromTokenWithChecks("QxrH0iJc3wv/lj/YKVppNYRad7tN0Z3x", false)
	database.UpdateLastSent(testTelephone, grcftwc.ClientID, 1)

	// create a request to pass to our handler.
	reader := strings.NewReader(
		"token=" + url.QueryEscape(testToken) +
			"&u=amber&p=amber&t=" + testTelephone +
			"&m=" + url.QueryEscape("Your AMBER Car is on its way! Tap here to TRACK YOUR CAR in the APP or CALL YOUR DRIVER enroute.") +
			"&h=0&v=456749&j=1234567&d=123654" +
			"&r=" + url.QueryEscape("AB12 CDE") +
			"&gr_token=" + url.QueryEscape("QxrH0iJc3wv/lj/YKVppNYRad7tN0Z3x") +
			"&ignore_telephone_checks=1")
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/googlereviews", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(GoogleReviewsHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	// expected := `{"success":"1"}`
	expected := string(grcftwc.SendSuccessResponse)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// NOTE: sends a message to the server set up in config could be the live server
// Ignore telephone checks and ignore dispatcher checks
func TestGoogleReviewsHandlerIgnoreTelephoneChecksAndIgnoreDispatcherChecks(t *testing.T) {
	prepareTestDatabase()

	// add last sent for test telephone
	grcftwc := database.ConfigFromTokenWithChecks("QxrH0iJc3wv/lj/YKVppNYRad7tN0Z3x", false)
	database.UpdateLastSent(testTelephone, grcftwc.ClientID, 1)

	// create a request to pass to our handler.
	reader := strings.NewReader(
		"token=" + url.QueryEscape(testToken) +
			"&u=amber&p=amber&t=" + testTelephone +
			"&m=" + url.QueryEscape("Your AMBER Car is on its way! Tap here to TRACK YOUR CAR in the APP or CALL YOUR DRIVER enroute.") +
			"&h=0&v=456749&j=1234567&d=123654" +
			"&r=" + url.QueryEscape("AB12 CDE") +
			"&gr_token=" + url.QueryEscape("QxrH0iJc3wv/lj/YKVppNYRad7tN0Z3x") +
			"&ignore_telephone_checks=1" +
			"&ignore_dispatcher_checks=1")
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/googlereviews", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(GoogleReviewsHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	// expected := `{"success":"1"}`
	expected := string(grcftwc.SendSuccessResponse)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// NOTE: sends a message to the server set up in config could be the live server
// Ignore telephone checks and ignore dispatcher checks and time and sent cont checks and ignore time and sent count checks
func TestGoogleReviewsHandlerIgnoreTelephoneChecksAndIgnoreDispatcherChecksAndIgnoreTimeAndSentCountChecks(t *testing.T) {
	prepareTestDatabase()

	// add last sent for test telephone
	grcftwc := database.ConfigFromTokenWithChecks("QxrH0iJc3wv/lj/YKVppNYRad7tN0Z3x", false)
	database.UpdateLastSent(testTelephone, grcftwc.ClientID, 1)

	// create a request to pass to our handler.
	reader := strings.NewReader(
		"token=" + url.QueryEscape(testToken) +
			"&u=amber&p=amber&t=" + testTelephone +
			"&m=" + url.QueryEscape("Your AMBER Car is on its way! Tap here to TRACK YOUR CAR in the APP or CALL YOUR DRIVER enroute.") +
			"&h=0&v=456749&j=1234567&d=123654" +
			"&r=" + url.QueryEscape("AB12 CDE") +
			"&gr_token=" + url.QueryEscape("QxrH0iJc3wv/lj/YKVppNYRad7tN0Z3x") +
			"&ignore_telephone_checks=1" +
			"&ignore_dispatcher_checks=1" +
			"&ignore_time_and_sent_count_checks=1")
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/googlereviews", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(GoogleReviewsHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	// expected := `{"success":"1"}`
	expected := string(grcftwc.SendSuccessResponse)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// NOTE: sends a message to the server set up in config could be the live server
// For Review Master SMS Gateway
func TestGoogleReviewsHandlerReviewMasterSMSGateway(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"t=" + testTelephone +
			"&message=" + url.QueryEscape("Test Message 1 from Review Master SMS Gateway") +
			"&gr_token=" + url.QueryEscape("OYBpBsZ9OhR-nbsupMQU_hCTJuaabtZ1gsfsfp"))

	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/googlereviews", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(GoogleReviewsHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "ok"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// NOTE: sends a message to the server set up in config could be the live server
// For Review Master SMS Gateway on the master queue
func TestGoogleReviewsHandlerReviewMasterSMSGatewayMasterQueue(t *testing.T) {
	prepareTestDatabase()

	// create a request to pass to our handler.
	reader := strings.NewReader(
		"t=" + testTelephone +
			"&message=" + url.QueryEscape("Test Message 1 from Review Master SMS Gateway") +
			"&gr_token=" + url.QueryEscape("OYBpBsZ9OhR-nbsupMQU_hCTJuaabtZ1gsfsfp123"))

	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/googlereviews", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(GoogleReviewsHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "ok"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// NOTE: sends a message to the server set up in config could be the live server
// For alternate message service Message Media
func TestGoogleReviewsHandlerReviewMessageMedia(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"t=" + testTelephone +
			"&m=" + url.QueryEscape("Test Message 1 from Message Media") +
			"&gr_token=" + url.QueryEscape("OYBpBsZ9OhR-nbsupMQU_hCTJuaabtf5esfsfp") +
			"&api_key=" + testMessageMediaApiKey +
			"&api_secret=" + testMessageMediaApiSecret)

	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/googlereviews", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(GoogleReviewsHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "ok"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// NOTE: sends a message to the server set up in config could be the live server
func TestGoogleReviewsHandlerDutchCompanyTelephoneNotFound(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"u=amber&p=amber&t=" + "0031615363999" +
			"&m=" + url.QueryEscape("Testing") +
			"&h=0&v=456749&j=1234567&d=123654" +
			"&gr_token=" + url.QueryEscape("9niammrACc18lklup-rQolaZmqtwnhptB8PNYqaMXPx-TKuCJ49dTb_8uHaeEPge"))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/googlereviews", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(GoogleReviewsHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"success":"1"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGoogleReviewsHandlerDisabledClientAndConfig(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"token=" + url.QueryEscape(testToken) +
			"&u=amber&p=amber&t=" + testTelephone +
			"&m=" + url.QueryEscape("Test") +
			"&h=0&v=456749&j=1234567&d=123654" +
			"&r=" + url.QueryEscape("AB12 CDE") +
			"&gr_token=" + url.QueryEscape("r6KAQpqdLhHnxZtUmFupDLA6zkL0LjdpkJCn-rBQ7og35i1Sxg-SQ0HxUrERDE3_"))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/googlereviews", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(GoogleReviewsHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"success":"0"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
