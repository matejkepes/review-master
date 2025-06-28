package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"google_reviews/database"
)

func TestReviewMasterSMSGatewayPairingHandler1(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	req, err := http.NewRequest("GET", "/rmsgpair?pairing_code=tpyh17azv43y", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("api-token", testReviewMasterSMSGatewayPairingToken)

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	// using same log token as sending so will match
	handler := http.Handler(ReviewMasterSMSGatewayPairingHandler(testReviewMasterSMSGatewayPairingToken))

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	if rr.Body.String() == `` {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
	if rr.Body.String() != `{"queue_id":"12"}` {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
	fmt.Println(rr.Body.String())
}

func TestReviewMasterSMSGatewayPairingHandler2WrongToken(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	req, err := http.NewRequest("GET", "/rmsgpair?pairing_code=tpyh17azv43y", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("api-token", "rubbishToken")

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	// using same log token as sending so will match
	handler := http.Handler(ReviewMasterSMSGatewayPairingHandler(testReviewMasterSMSGatewayPairingToken))

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

	// Check the response body is what we expect.
	if rr.Body.String() != `` {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
	fmt.Println(rr.Body.String())
}

func TestReviewMasterSMSGatewayPairingHandler3PairingCodeNotFound(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	req, err := http.NewRequest("GET", "/rmsgpair?pairing_code=1234hjdccdgjasjdgskdx8ysbGF", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("api-token", testReviewMasterSMSGatewayPairingToken)

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	// using same log token as sending so will match
	handler := http.Handler(ReviewMasterSMSGatewayPairingHandler(testReviewMasterSMSGatewayPairingToken))

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	if rr.Body.String() == `` {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
	if rr.Body.String() != `{"queue_id":"0"}` {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
	fmt.Println(rr.Body.String())
}

func TestReviewMasterSMSGatewayPairingHandler4WrongParameter(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	req, err := http.NewRequest("GET", "/rmsgpair?pairing=tpyh17azv43y", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("api-token", testReviewMasterSMSGatewayPairingToken)

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	// using same log token as sending so will match
	handler := http.Handler(ReviewMasterSMSGatewayPairingHandler(testReviewMasterSMSGatewayPairingToken))

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	// Check the response body is what we expect.
	if rr.Body.String() != `` {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
	fmt.Println(rr.Body.String())
}

func TestReviewMasterSMSGatewayPairingHandler5MasterQueue(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	req, err := http.NewRequest("GET", "/rmsgpair?pairing_code=tpyh17azv8uf", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("api-token", testReviewMasterSMSGatewayPairingToken)

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	// using same log token as sending so will match
	handler := http.Handler(ReviewMasterSMSGatewayPairingHandler(testReviewMasterSMSGatewayPairingToken))

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	if rr.Body.String() == `` {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
	checkQueue := fmt.Sprintf("{\"queue_id\":\"%d\"}", database.ReviewMasterSMSMasterQueue)
	if rr.Body.String() != checkQueue {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
	fmt.Println(rr.Body.String())
}
