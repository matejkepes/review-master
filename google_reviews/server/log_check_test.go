package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestCheckLogHandler1(t *testing.T) {
	// create a request to pass to our handler.
	req, err := http.NewRequest("GET", "/checklogs?log_token="+url.QueryEscape(testLogToken)+"&hours_back=2400", nil)
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	// using same log token as sending so will match
	handler := http.Handler(CheckLogHandler("../google_reviews.log", testLogToken))

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
	fmt.Println(rr.Body.String())
}

func TestCheckLogHandler2WrongToken(t *testing.T) {
	// create a request to pass to our handler.
	req, err := http.NewRequest("GET", "/checklogs?log_token="+url.QueryEscape("rubbishToken")+"&hours_back=2400", nil)
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	// using same log token as sending so will match
	handler := http.Handler(CheckLogHandler("../google_reviews.log", testLogToken))

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

	// Check the response body is what we expect.
	if rr.Body == nil {
		t.Errorf("handler returned unexpected body: got %+v", rr.Body)
	}
	fmt.Println(rr.Body.String())
}

func TestCheckLog1(t *testing.T) {
	logs := CheckLog("../google_reviews.log", 2400)
	fmt.Println(logs)
	// server.CheckLog("../google_reviews.log", 240)
}
