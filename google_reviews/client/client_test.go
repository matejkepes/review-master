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
	params.Add("u", "amber")
	params.Add("p", "amber")
	params.Add("t", "447123456789")
	params.Add("m", "Your AMBER Car is on its way! Tap here to TRACK YOUR CAR in the APP or CALL YOUR DRIVER enroute.")
	params.Add("h", "0")
	params.Add("v", "456749")
	params.Add("j", "1234567")
	params.Add("d", "123654")
	params.Add("r", "AB12 CDE")
	params.Add("x", "{\"b\":\"987654\",\"t\":\"2\"}")

	// use Client & URL from local test server
	// https://appmessages.veezu.co.uk/api/v1/sendmessage
	// "https://test1.veezu.co.uk:8893/api/v1/sendmessage"
	testURL := testServer.URL
	resp := Send(testURL, "POST", "", "", "", params, nil, false, "", "")
	fmt.Println("resp: ", resp)
}

// Not running a test server, send direct to server specified (live)
func TestSendLive(t *testing.T) {

	params := url.Values{}
	params.Add("token", testToken)
	params.Add("u", "amber")
	params.Add("p", "amber")
	params.Add("t", testTelephone)
	params.Add("m", "Your AMBER Car is on its way! Tap here to TRACK YOUR CAR in the APP or CALL YOUR DRIVER enroute.")
	params.Add("h", "0")
	params.Add("v", "456749")
	params.Add("j", "1234567")
	params.Add("d", "123654")
	params.Add("r", "AB12 CDE")
	params.Add("x", "{\"b\":\"987654\",\"t\":\"2\"}")

	// use Client & URL from local test server
	// https://appmessages.veezu.co.uk/api/v1/sendmessage
	// "https://test1.veezu.co.uk:8893/api/v1/sendmessage"
	// resp := Send("https://appmessages.veezu.co.uk/api/v1/sendmessage", params)
	resp := Send(testURL, "POST", "", "", "", params, nil, false, "", "")
	fmt.Println("resp: ", resp)
}

// Not running a test server, send direct to server specified (live - our own server)
func TestSendLive2(t *testing.T) {

	params := url.Values{}
	params.Add("token", testToken)
	params.Add("t", testTelephone)
	params.Add("m", "Your TEST Car is on its way!")

	resp := Send(testURL, "POST", "", "", "", params, nil, false, "", "")
	fmt.Println("resp: ", resp)
}

// Not running a test server, send direct to server specified (live - Alpha (uses DiSC) server)
func TestSendLive3(t *testing.T) {

	params := url.Values{}
	params.Add("number", testTelephone)
	params.Add("msg", "Your TEST Car is on its way!")

	resp := Send(testAlphaURL, "GET", "", "", "", params, nil, false, "", "")
	fmt.Println("resp: ", resp)
}

func TestSendRetrieveBooking1(t *testing.T) {
	resp := Send(testIcabbiURL+"/bookings/index/"+testIcabbiTripID, "GET", testIcabbiAppKey, testIcabbiSecretKey, "", nil, nil, false, "", "")
	fmt.Println("resp: ", resp)
}

// Booking does not exist
func TestSendRetrieveBooking2(t *testing.T) {
	resp := Send(testIcabbiURL+"/bookings/index/12374975395936756", "GET", testIcabbiAppKey, testIcabbiSecretKey, "", nil, nil, false, "", "")
	fmt.Println("resp: ", resp)
}

// Has a rating
func TestSendCustomerRating1(t *testing.T) {
	resp := Send(testIcabbiURL+"/customerrating/check?phone_number="+testIcabbiPhone, "GET", testIcabbiAppKey, testIcabbiSecretKey, "", nil, nil, false, "", "")
	fmt.Println("resp: ", resp)
}

// User does not exist
func TestSendCustomerRating2(t *testing.T) {
	resp := Send(testIcabbiURL+"/customerrating/check?phone_number="+testIcabbiPhoneDoesNotExist, "GET", testIcabbiAppKey, testIcabbiSecretKey, "", nil, nil, false, "", "")
	fmt.Println("resp: ", resp)
}

// Review Master SMS Gateway
func TestSendReviewMasterSMSGateway(t *testing.T) {
	body, _ := json.Marshal(map[string]string{
		"queue_id":  "81",
		"telephone": "+" + testTelephone,
		"message":   "testing",
	})
	fmt.Println("body:", string(body))

	resp := Send(testReviewMasterSMSGatewayURL, "POST", "", "", testReviewMasterSMSGatewayApiToken, nil, body, false, "", "")
	fmt.Println("resp: ", resp)
}

// Message Media send SMS Rest API (https://support.messagemedia.com/hc/en-us/articles/4413635760527-Messaging-API)
func TestSendMessageMediaSMSRestAPI(t *testing.T) {
	type message struct {
		Content           string `json:"content"`
		DestinationNumber string `json:"destination_number"`
		Format            string `json:"format"`
		DeliveryReport    string `json:"delivery_report"`
	}
	type messages struct {
		Messages []message `json:"messages"`
	}
	ms := messages{
		Messages: []message{
			{
				Content:           "testing",
				DestinationNumber: "+" + testTelephone,
				Format:            "SMS",
				DeliveryReport:    "true",
			},
		},
	}
	body, _ := json.Marshal(ms)
	fmt.Println("body:", string(body))

	// example successful response:
	// {"messages":[{"callback_url":null,"delivery_report":true,"destination_number":"+447889525579","format":"SMS","message_expiry_timestamp":null,"message_flags":[],"message_id":"13d1a0ba-be11-401d-b4e8-7f27d8ccce99","metadata":null,"scheduled":null,"status":"queued","content":"testing","source_number":null,"rich_link":null,"media":null,"subject":null}]}
	// uses basic authenticaion (therefore set appKey and secretKey) and set json to true
	resp := Send("https://api.messagemedia.com/v1/messages", "POST", testMessageMediaApiKey, testMessageMediaApiSecret, "", nil, body, true, "Message Media", "")
	fmt.Println("resp: ", resp)
}

// Veezu send message API (https://messages.veezu.com/api/messages)
func TestSendVeezuMessageAPI(t *testing.T) {
	type veezuMessage struct {
		Message   string `json:"message"`
		Telephone string `json:"telephone"`
	}
	ms := veezuMessage{
		Message:   "This is a test",
		Telephone: testTelephone,
	}
	body, _ := json.Marshal(ms)
	fmt.Println("body:", string(body))

	// example successful response:
	// {"success":"1"}
	// uses auth_token in header for authenticaion (therefore set alternate_message_service_secret1) and set json to true
	resp := Send("https://messages.veezu.com/api/messages", "POST", "", "", "", nil, body, true, "Veezu", testAuthToken)
	fmt.Println("resp: ", resp)

	// decode the response
	type veezuMessageResponse struct {
		Success string `json:"success"`
	}
	var vmr veezuMessageResponse
	json.Unmarshal([]byte(resp), &vmr)
	if vmr.Success != "1" {
		t.Fatalf("Error sending message for Veezu alternate message service response from send server: %+v", resp)
	}
}
