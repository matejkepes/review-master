package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func testServer(success bool) *httptest.Server {
	// start a local HTTP server
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("req: ", req)
		// send response to be tested
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		// data := []byte(`[{"success":"1"}]`)
		data := string(`{"success":"1"}`)
		if !success {
			data = string(`{"success":"0"}`)
		}
		json.NewEncoder(w).Encode(data)
	}))

	return testServer
}

func TestSend(t *testing.T) {
	// start a local HTTP server
	testServer := testServer(true)
	// close the server when test finishes
	defer testServer.Close()

	params := url.Values{}
	params.Add("token", testToken)
	params.Add("t", "447123456789")
	params.Add("m", "Your Car is on its way! Tap here to TRACK YOUR CAR in the APP or CALL YOUR DRIVER enroute.")

	// use Client & URL from local test server
	testURL := testServer.URL
	resp := Send(testURL, params)
	fmt.Println("resp: ", resp)
}

// Not running a test server, send direct to server specified (live)
func TestSendLive(t *testing.T) {

	params := url.Values{}
	params.Add("token", testToken)
	params.Add("t", testTelephone)
	params.Add("m", "Your Car is on its way! Tap here to TRACK YOUR CAR in the APP or CALL YOUR DRIVER enroute.")

	// use Client & URL from local test server
	resp := Send(testURL, params)
	fmt.Println("resp: ", resp)
}

// Not running a test server, send direct to server specified (live - our own server)
func TestSendLive2(t *testing.T) {

	params := url.Values{}
	params.Add("token", testToken)
	params.Add("t", testTelephone)
	params.Add("m", "Your TEST Car is on its way!")

	resp := Send(testURL, params)
	fmt.Println("resp: ", resp)
}
