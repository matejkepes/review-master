package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestCab9Handler(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"gr_token=" + url.QueryEscape("JE0v6xqCdvIDKXOX_JmUM09p-vYjzkLOUc5iysc9q1vqttMN9s0DfVFD3yv8yLz4") +
			"&telephone=" + testTelephone +
			"&booking_creation_time=" + url.QueryEscape(time.Now().Add(time.Minute*-30).Format(time.RFC3339)) +
			"&booked_for_time=" + url.QueryEscape(time.Now().Add(time.Minute*-25).Format(time.RFC3339)) +
			"&picked_up_time=" + url.QueryEscape(time.Now().Add(time.Minute*-20).Format(time.RFC3339)))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/cab9", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(Cab9Handler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := string(cab9SuccessResponse)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCab9HandlerRepeatRequest(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"gr_token=" + url.QueryEscape("JE0v6xqCdvIDKXOX_JmUM09p-vYjzkLOUc5iysc9q1vqttMN9s0DfVFD3yv8yLz4") +
			"&telephone=" + testTelephone +
			"&booking_creation_time=" + url.QueryEscape(time.Now().Add(time.Minute*-30).Format(time.RFC3339)) +
			"&booked_for_time=" + url.QueryEscape(time.Now().Add(time.Minute*-25).Format(time.RFC3339)) +
			"&picked_up_time=" + url.QueryEscape(time.Now().Add(time.Minute*-20).Format(time.RFC3339)))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/cab9", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(Cab9Handler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := string(cab9SuccessResponse)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// second request
	// create a request to pass to our handler.
	reader = strings.NewReader(
		"gr_token=" + url.QueryEscape("JE0v6xqCdvIDKXOX_JmUM09p-vYjzkLOUc5iysc9q1vqttMN9s0DfVFD3yv8yLz4") +
			"&telephone=" + testTelephone +
			"&booking_creation_time=" + url.QueryEscape(time.Now().Add(time.Minute*-30).Format(time.RFC3339)) +
			"&booked_for_time=" + url.QueryEscape(time.Now().Add(time.Minute*-25).Format(time.RFC3339)) +
			"&picked_up_time=" + url.QueryEscape(time.Now().Add(time.Minute*-20).Format(time.RFC3339)))
	fmt.Println("reader: ", reader)
	req, err = http.NewRequest("POST", "/cab9", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr = httptest.NewRecorder()
	handler = http.Handler(Cab9Handler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	// This is not sent though
	expected = string(cab9SuccessResponse)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCab9HandlerRepeatRequestIgnoreTelephoneChecks(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"gr_token=" + url.QueryEscape("JE0v6xqCdvIDKXOX_JmUM09p-vYjzkLOUc5iysc9q1vqttMN9s0DfVFD3yv8yLz4") +
			"&telephone=" + testTelephone +
			"&booking_creation_time=" + url.QueryEscape(time.Now().Add(time.Minute*-30).Format(time.RFC3339)) +
			"&booked_for_time=" + url.QueryEscape(time.Now().Add(time.Minute*-25).Format(time.RFC3339)) +
			"&picked_up_time=" + url.QueryEscape(time.Now().Add(time.Minute*-20).Format(time.RFC3339)))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/cab9", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(Cab9Handler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := string(cab9SuccessResponse)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// second request
	// create a request to pass to our handler.
	reader = strings.NewReader(
		"gr_token=" + url.QueryEscape("JE0v6xqCdvIDKXOX_JmUM09p-vYjzkLOUc5iysc9q1vqttMN9s0DfVFD3yv8yLz4") +
			"&telephone=" + testTelephone +
			"&booking_creation_time=" + url.QueryEscape(time.Now().Add(time.Minute*-30).Format(time.RFC3339)) +
			"&booked_for_time=" + url.QueryEscape(time.Now().Add(time.Minute*-25).Format(time.RFC3339)) +
			"&picked_up_time=" + url.QueryEscape(time.Now().Add(time.Minute*-20).Format(time.RFC3339)) +
			"&ignore_telephone_checks=1")
	fmt.Println("reader: ", reader)
	req, err = http.NewRequest("POST", "/cab9", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr = httptest.NewRecorder()
	handler = http.Handler(Cab9Handler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected = string(cab9SuccessResponse)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCab9HandlerIncorrectToken(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"gr_token=" + url.QueryEscape("qey-FMF9Wun-dAJQ6Ri1wBv1hh7DsjcH7QRM7WMDv9MUeEdgNF") +
			"&telephone=" + testTelephone +
			"&booking_creation_time=" + url.QueryEscape(time.Now().Add(time.Minute*-30).Format(time.RFC3339)) +
			"&booked_for_time=" + url.QueryEscape(time.Now().Add(time.Minute*-25).Format(time.RFC3339)) +
			"&picked_up_time=" + url.QueryEscape(time.Now().Add(time.Minute*-20).Format(time.RFC3339)))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/cab9", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(Cab9Handler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := string(cab9FailedResponse)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCab9HandlerNotCab9Dispatcher(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"gr_token=" + url.QueryEscape("OYBpBsZ9OhR-nbsupMQU_hCTJuaabtZ1gsfsfp123") +
			"&telephone=" + testTelephone +
			"&booking_creation_time=" + url.QueryEscape(time.Now().Add(time.Minute*-30).Format(time.RFC3339)) +
			"&booked_for_time=" + url.QueryEscape(time.Now().Add(time.Minute*-25).Format(time.RFC3339)) +
			"&picked_up_time=" + url.QueryEscape(time.Now().Add(time.Minute*-20).Format(time.RFC3339)))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/cab9", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(Cab9Handler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := string(cab9FailedResponse)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCab9HandlerTelephoneIncorrect(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"gr_token=" + url.QueryEscape("JE0v6xqCdvIDKXOX_JmUM09p-vYjzkLOUc5iysc9q1vqttMN9s0DfVFD3yv8yLz4") +
			"&telephone=0712345678" +
			"&booking_creation_time=" + url.QueryEscape(time.Now().Add(time.Minute*-30).Format(time.RFC3339)) +
			"&booked_for_time=" + url.QueryEscape(time.Now().Add(time.Minute*-25).Format(time.RFC3339)) +
			"&picked_up_time=" + url.QueryEscape(time.Now().Add(time.Minute*-20).Format(time.RFC3339)))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/cab9", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(Cab9Handler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := string(cab9FailedResponse)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCab9HandlerNoBookingCreationTimeSent(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"gr_token=" + url.QueryEscape("JE0v6xqCdvIDKXOX_JmUM09p-vYjzkLOUc5iysc9q1vqttMN9s0DfVFD3yv8yLz4") +
			"&telephone=" + testTelephone +
			"&booked_for_time=" + url.QueryEscape(time.Now().Add(time.Minute*-25).Format(time.RFC3339)) +
			"&picked_up_time=" + url.QueryEscape(time.Now().Add(time.Minute*-20).Format(time.RFC3339)))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/cab9", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(Cab9Handler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := string(cab9FailedResponse)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCab9HandlerNoBookedForTimeSent(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"gr_token=" + url.QueryEscape("JE0v6xqCdvIDKXOX_JmUM09p-vYjzkLOUc5iysc9q1vqttMN9s0DfVFD3yv8yLz4") +
			"&telephone=" + testTelephone +
			"&booking_creation_time=" + url.QueryEscape(time.Now().Add(time.Minute*-30).Format(time.RFC3339)) +
			"&picked_up_time=" + url.QueryEscape(time.Now().Add(time.Minute*-20).Format(time.RFC3339)))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/cab9", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(Cab9Handler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := string(cab9FailedResponse)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCab9HandlerNoPickuUpTimeSent(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"gr_token=" + url.QueryEscape("JE0v6xqCdvIDKXOX_JmUM09p-vYjzkLOUc5iysc9q1vqttMN9s0DfVFD3yv8yLz4") +
			"&telephone=" + testTelephone +
			"&booking_creation_time=" + url.QueryEscape(time.Now().Add(time.Minute*-30).Format(time.RFC3339)) +
			"&booked_for_time=" + url.QueryEscape(time.Now().Add(time.Minute*-25).Format(time.RFC3339)))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/cab9", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(Cab9Handler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := string(cab9FailedResponse)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCab9HandlerPreBooking(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"gr_token=" + url.QueryEscape("JE0v6xqCdvIDKXOX_JmUM09p-vYjzkLOUc5iysc9q1vqttMN9s0DfVFD3yv8yLz4") +
			"&telephone=" + testTelephone +
			"&booking_creation_time=" + url.QueryEscape(time.Now().Add(time.Hour*-30).Format(time.RFC3339)) +
			"&booked_for_time=" + url.QueryEscape(time.Now().Add(time.Minute*-25).Format(time.RFC3339)) +
			"&picked_up_time=" + url.QueryEscape(time.Now().Add(time.Minute*-24).Format(time.RFC3339)))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/cab9", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(Cab9Handler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := string(cab9SuccessResponse)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCab9HandlerPickupEarlyOk(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"gr_token=" + url.QueryEscape("JE0v6xqCdvIDKXOX_JmUM09p-vYjzkLOUc5iysc9q1vqttMN9s0DfVFD3yv8yLz4") +
			"&telephone=" + testTelephone +
			"&booking_creation_time=" + url.QueryEscape(time.Now().Add(time.Hour*-30).Format(time.RFC3339)) +
			"&booked_for_time=" + url.QueryEscape(time.Now().Add(time.Minute*-25).Format(time.RFC3339)) +
			"&picked_up_time=" + url.QueryEscape(time.Now().Add(time.Minute*-30).Format(time.RFC3339)))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/cab9", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(Cab9Handler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := string(cab9SuccessResponse)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCab9HandlerFailsPickupTime(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"gr_token=" + url.QueryEscape("JE0v6xqCdvIDKXOX_JmUM09p-vYjzkLOUc5iysc9q1vqttMN9s0DfVFD3yv8yLz4") +
			"&telephone=" + testTelephone +
			"&booking_creation_time=" + url.QueryEscape(time.Now().Add(time.Hour*-30).Format(time.RFC3339)) +
			"&booked_for_time=" + url.QueryEscape(time.Now().Add(time.Minute*-25).Format(time.RFC3339)) +
			"&picked_up_time=" + url.QueryEscape(time.Now().Add(time.Minute*-10).Format(time.RFC3339)))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/cab9", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(Cab9Handler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := string(cab9FailedResponse)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
