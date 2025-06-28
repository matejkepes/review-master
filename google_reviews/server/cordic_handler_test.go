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

func TestCordicHandler(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"gr_token=" + url.QueryEscape("qey-FMF9Wun-dAJQ6Ri1wBv1hh7DsjcH7QRM7WMDv9MUeEdgNFwgW4pYxedTjfV8") +
			"&passenger_id=abcdefghijklmno" +
			"&booking_creation_time=" + url.QueryEscape(time.Now().Add(time.Minute*-30).Format(time.RFC3339)) +
			"&booked_for_time=" + url.QueryEscape(time.Now().Add(time.Minute*-25).Format(time.RFC3339)) +
			"&picked_up_time=" + url.QueryEscape(time.Now().Add(time.Minute*-20).Format(time.RFC3339)))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/cordic", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(CordicHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is not the empty message (failed)
	emptyMessage := `{"message":""}`
	if rr.Body.String() == emptyMessage {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	} else {
		fmt.Println("response: ", rr.Body.String())
	}
}

func TestCordicHandlerRepeatRequest(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"gr_token=" + url.QueryEscape("qey-FMF9Wun-dAJQ6Ri1wBv1hh7DsjcH7QRM7WMDv9MUeEdgNFwgW4pYxedTjfV8") +
			"&passenger_id=abcdefghijklmno" +
			"&booking_creation_time=" + url.QueryEscape(time.Now().Add(time.Minute*-30).Format(time.RFC3339)) +
			"&booked_for_time=" + url.QueryEscape(time.Now().Add(time.Minute*-25).Format(time.RFC3339)) +
			"&picked_up_time=" + url.QueryEscape(time.Now().Add(time.Minute*-20).Format(time.RFC3339)))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/cordic", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(CordicHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is not the empty message (failed)
	emptyMessage := `{"message":""}`
	if rr.Body.String() == emptyMessage {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	} else {
		fmt.Println("response: ", rr.Body.String())
	}

	// second request
	// create a request to pass to our handler.
	reader = strings.NewReader(
		"gr_token=" + url.QueryEscape("qey-FMF9Wun-dAJQ6Ri1wBv1hh7DsjcH7QRM7WMDv9MUeEdgNFwgW4pYxedTjfV8") +
			"&passenger_id=abcdefghijklmno" +
			"&booking_creation_time=" + url.QueryEscape(time.Now().Add(time.Minute*-30).Format(time.RFC3339)) +
			"&booked_for_time=" + url.QueryEscape(time.Now().Add(time.Minute*-25).Format(time.RFC3339)) +
			"&picked_up_time=" + url.QueryEscape(time.Now().Add(time.Minute*-20).Format(time.RFC3339)))
	fmt.Println("reader: ", reader)
	req, err = http.NewRequest("POST", "/cordic", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr = httptest.NewRecorder()
	handler = http.Handler(CordicHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is not the empty message (failed)
	emptyMessage = `{"message":""}`
	if rr.Body.String() == emptyMessage {
		fmt.Println("response: ", rr.Body.String())
	} else {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
}

func TestCordicHandlerRepeatRequestIgnorePassengerIDCheck(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"gr_token=" + url.QueryEscape("qey-FMF9Wun-dAJQ6Ri1wBv1hh7DsjcH7QRM7WMDv9MUeEdgNFwgW4pYxedTjfV8") +
			"&passenger_id=abcdefghijklmno" +
			"&booking_creation_time=" + url.QueryEscape(time.Now().Add(time.Minute*-30).Format(time.RFC3339)) +
			"&booked_for_time=" + url.QueryEscape(time.Now().Add(time.Minute*-25).Format(time.RFC3339)) +
			"&picked_up_time=" + url.QueryEscape(time.Now().Add(time.Minute*-20).Format(time.RFC3339)))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/cordic", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(CordicHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is not the empty message (failed)
	emptyMessage := `{"message":""}`
	if rr.Body.String() == emptyMessage {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	} else {
		fmt.Println("response: ", rr.Body.String())
	}

	// second request ignore passenger request checks
	// create a request to pass to our handler.
	reader = strings.NewReader(
		"gr_token=" + url.QueryEscape("qey-FMF9Wun-dAJQ6Ri1wBv1hh7DsjcH7QRM7WMDv9MUeEdgNFwgW4pYxedTjfV8") +
			"&passenger_id=abcdefghijklmno" +
			"&booking_creation_time=" + url.QueryEscape(time.Now().Add(time.Minute*-30).Format(time.RFC3339)) +
			"&booked_for_time=" + url.QueryEscape(time.Now().Add(time.Minute*-25).Format(time.RFC3339)) +
			"&picked_up_time=" + url.QueryEscape(time.Now().Add(time.Minute*-20).Format(time.RFC3339)) +
			"&ignore_passenger_id_checks=1")
	fmt.Println("reader: ", reader)
	req, err = http.NewRequest("POST", "/cordic", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr = httptest.NewRecorder()
	handler = http.Handler(CordicHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is not the empty message (failed)
	emptyMessage = `{"message":""}`
	if rr.Body.String() == emptyMessage {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	} else {
		fmt.Println("response: ", rr.Body.String())
	}
}

func TestCordicHandlerIncorrectToken(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"gr_token=" + url.QueryEscape("qey-FMF9Wun-dAJQ6Ri1wBv1hh7DsjcH7QRM7WMDv9MUeEdgNF") +
			"&passenger_id=abcdefghijklmno" +
			"&booking_creation_time=" + url.QueryEscape(time.Now().Add(time.Minute*-30).Format(time.RFC3339)) +
			"&booked_for_time=" + url.QueryEscape(time.Now().Add(time.Minute*-25).Format(time.RFC3339)) +
			"&picked_up_time=" + url.QueryEscape(time.Now().Add(time.Minute*-20).Format(time.RFC3339)))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/cordic", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(CordicHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is not the empty message (failed)
	emptyMessage := `{"message":""}`
	if rr.Body.String() == emptyMessage {
		fmt.Println("response: ", rr.Body.String())
	} else {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
}

func TestCordicHandlerNotCordicDispatcher(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"gr_token=" + url.QueryEscape("OYBpBsZ9OhR-nbsupMQU_hCTJuaabtZ1gsfsfp123") +
			"&passenger_id=abcdefghijklmno" +
			"&booking_creation_time=" + url.QueryEscape(time.Now().Add(time.Minute*-30).Format(time.RFC3339)) +
			"&booked_for_time=" + url.QueryEscape(time.Now().Add(time.Minute*-25).Format(time.RFC3339)) +
			"&picked_up_time=" + url.QueryEscape(time.Now().Add(time.Minute*-20).Format(time.RFC3339)))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/cordic", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(CordicHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is not the empty message (failed)
	emptyMessage := `{"message":""}`
	if rr.Body.String() == emptyMessage {
		fmt.Println("response: ", rr.Body.String())
	} else {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
}

func TestCordicHandlerNoPassengerID(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"gr_token=" + url.QueryEscape("qey-FMF9Wun-dAJQ6Ri1wBv1hh7DsjcH7QRM7WMDv9MUeEdgNFwgW4pYxedTjfV8") +
			"&booking_creation_time=" + url.QueryEscape(time.Now().Add(time.Minute*-30).Format(time.RFC3339)) +
			"&booked_for_time=" + url.QueryEscape(time.Now().Add(time.Minute*-25).Format(time.RFC3339)) +
			"&picked_up_time=" + url.QueryEscape(time.Now().Add(time.Minute*-20).Format(time.RFC3339)))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/cordic", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(CordicHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is not the empty message (failed)
	emptyMessage := `{"message":""}`
	if rr.Body.String() == emptyMessage {
		fmt.Println("response: ", rr.Body.String())
	} else {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
}

func TestCordicHandlerNoBookingCreationTimeSent(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"gr_token=" + url.QueryEscape("qey-FMF9Wun-dAJQ6Ri1wBv1hh7DsjcH7QRM7WMDv9MUeEdgNFwgW4pYxedTjfV8") +
			"&passenger_id=abcdefghijklmno" +
			"&booked_for_time=" + url.QueryEscape(time.Now().Add(time.Minute*-25).Format(time.RFC3339)) +
			"&picked_up_time=" + url.QueryEscape(time.Now().Add(time.Minute*-20).Format(time.RFC3339)))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/cordic", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(CordicHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is not the empty message (failed)
	emptyMessage := `{"message":""}`
	if rr.Body.String() == emptyMessage {
		fmt.Println("response: ", rr.Body.String())
	} else {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
}

func TestCordicHandlerNoBookedForTimeSent(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"gr_token=" + url.QueryEscape("qey-FMF9Wun-dAJQ6Ri1wBv1hh7DsjcH7QRM7WMDv9MUeEdgNFwgW4pYxedTjfV8") +
			"&passenger_id=abcdefghijklmno" +
			"&booking_creation_time=" + url.QueryEscape(time.Now().Add(time.Minute*-30).Format(time.RFC3339)) +
			"&picked_up_time=" + url.QueryEscape(time.Now().Add(time.Minute*-20).Format(time.RFC3339)))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/cordic", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(CordicHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is not the empty message (failed)
	emptyMessage := `{"message":""}`
	if rr.Body.String() == emptyMessage {
		fmt.Println("response: ", rr.Body.String())
	} else {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
}

func TestCordicHandlerNoPickedUpTimeSent(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"gr_token=" + url.QueryEscape("qey-FMF9Wun-dAJQ6Ri1wBv1hh7DsjcH7QRM7WMDv9MUeEdgNFwgW4pYxedTjfV8") +
			"&passenger_id=abcdefghijklmno" +
			"&booking_creation_time=" + url.QueryEscape(time.Now().Add(time.Minute*-30).Format(time.RFC3339)) +
			"&booked_for_time=" + url.QueryEscape(time.Now().Add(time.Minute*-25).Format(time.RFC3339)))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/cordic", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(CordicHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is not the empty message (failed)
	emptyMessage := `{"message":""}`
	if rr.Body.String() == emptyMessage {
		fmt.Println("response: ", rr.Body.String())
	} else {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
}

func TestCordicHandlerPreBooking(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"gr_token=" + url.QueryEscape("qey-FMF9Wun-dAJQ6Ri1wBv1hh7DsjcH7QRM7WMDv9MUeEdgNFwgW4pYxedTjfV8") +
			"&passenger_id=abcdefghijklmno" +
			"&booking_creation_time=" + url.QueryEscape(time.Now().Add(time.Hour*-30).Format(time.RFC3339)) +
			"&booked_for_time=" + url.QueryEscape(time.Now().Add(time.Minute*-25).Format(time.RFC3339)) +
			"&picked_up_time=" + url.QueryEscape(time.Now().Add(time.Minute*-24).Format(time.RFC3339)))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/cordic", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(CordicHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is not the empty message (failed)
	emptyMessage := `{"message":""}`
	if rr.Body.String() == emptyMessage {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	} else {
		fmt.Println("response: ", rr.Body.String())
	}
}

func TestCordicHandlerPickupEarlyOk(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"gr_token=" + url.QueryEscape("qey-FMF9Wun-dAJQ6Ri1wBv1hh7DsjcH7QRM7WMDv9MUeEdgNFwgW4pYxedTjfV8") +
			"&passenger_id=abcdefghijklmno" +
			"&booking_creation_time=" + url.QueryEscape(time.Now().Add(time.Hour*-30).Format(time.RFC3339)) +
			"&booked_for_time=" + url.QueryEscape(time.Now().Add(time.Minute*-25).Format(time.RFC3339)) +
			"&picked_up_time=" + url.QueryEscape(time.Now().Add(time.Minute*-30).Format(time.RFC3339)))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/cordic", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(CordicHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is not the empty message (failed)
	emptyMessage := `{"message":""}`
	if rr.Body.String() == emptyMessage {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	} else {
		fmt.Println("response: ", rr.Body.String())
	}
}

func TestCordicHandlerFailsPickupTime(t *testing.T) {
	prepareTestDatabase()
	// create a request to pass to our handler.
	reader := strings.NewReader(
		"gr_token=" + url.QueryEscape("qey-FMF9Wun-dAJQ6Ri1wBv1hh7DsjcH7QRM7WMDv9MUeEdgNFwgW4pYxedTjfV8") +
			"&passenger_id=abcdefghijklmno" +
			"&booking_creation_time=" + url.QueryEscape(time.Now().Add(time.Hour*-30).Format(time.RFC3339)) +
			"&booked_for_time=" + url.QueryEscape(time.Now().Add(time.Minute*-25).Format(time.RFC3339)) +
			"&picked_up_time=" + url.QueryEscape(time.Now().Add(time.Minute*-10).Format(time.RFC3339)))
	fmt.Println("reader: ", reader)
	req, err := http.NewRequest("POST", "/cordic", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(CordicHandler())

	// handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is not the empty message (failed)
	emptyMessage := `{"message":""}`
	if rr.Body.String() == emptyMessage {
		fmt.Println("response: ", rr.Body.String())
	} else {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
}
