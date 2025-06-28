package client

import (
	"fmt"
	"net/url"
	"testing"
)

// Test constants - update these with actual test values
const (
	testTelephone            = "07700900000" // Example test phone number
	testURL                  = "http://example.com/sms"
	testWheelsServerPassword = "testpassword"
	testWheelsURL            = "http://example.com/wheels/sms"
)

// Not running a test server, send direct to server specified (live)
func TestSendLive(t *testing.T) {

	params := url.Values{}
	params.Add("u", "ambercars")
	params.Add("p", "chrislovesamber")
	params.Add("t", testTelephone)
	params.Add("m", "Your Amber Car is on its way it is a  Red  Skoda Octavia reg  AB12 CDE")

	resp := Send(testURL, "POST", params)
	fmt.Println("resp: ", resp)
}

// Not running a test server, send direct to server specified (live)
func TestSendLiveWheels(t *testing.T) {

	params := url.Values{}
	params.Add("server_password", testWheelsServerPassword)
	params.Add("number", testTelephone)
	params.Add("msg", "Testing")

	resp := Send(testWheelsURL, "POST", params)
	fmt.Println("resp: ", resp)
}
