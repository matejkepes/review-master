package server

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"google_reviews/barred"
	"google_reviews/client"
	"google_reviews/config"
	"google_reviews/database"
	"google_reviews/utils"

	"github.com/dongri/phonenumber"
)

// var successResponse = []byte(`{"success":"1"}`)
var failedResponse = []byte(`{"success":"0"}`)
var Bars []string

// type IcabbiAppResponse struct {
// 	Code int `form:"code" json:"code" binding:"required"`
// }

// type googleReviewsHandler struct{}

// func (googleReviewsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
// 	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
// 	w.Write([]byte("This is the response from server.\n"))
// }

// // GoogleReviewsHandler - Google Reviews Handler
// func GoogleReviewsHandler(w http.ResponseWriter, req *http.Request) {
// 	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
// 	w.Write([]byte("This is the response from server.\n"))
// }

// // GoogleReviewsHandler - Google Reviews Handler
// func GoogleReviewsHandler() http.Handler {
// 	fn := func(w http.ResponseWriter, req *http.Request) {
// 		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
// 		w.Write([]byte("This is the response from server.\n"))
// 	}
// 	return http.HandlerFunc(fn)
// }

// // GoogleReviewsHandler - Google Reviews Handler
// func GoogleReviewsHandler() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
// 		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
// 		w.Write([]byte("This is the response from server.\n"))
// 	})
// }

// GoogleReviewsHandler - Google Reviews Handler
func GoogleReviewsHandler() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			// fmt.Printf("ParseForm() err: %v\n", err)
			log.Printf("ParseForm() err: %v\n", err)
			w.Write(failedResponse)
			return
		}
		// log.Printf("r.PostForm = %v\n", req.PostForm)

		// check google reviews token
		grToken := strings.TrimSpace(req.FormValue("gr_token"))
		// log.Printf("grToken: %s\n", grToken)
		// check whether should ignore the time and daily sent count checks (used for testing on front end)
		ignoreTimeAndSentCountCheck := false
		ignoreTimeAndSentCountCheckReq := strings.TrimSpace(req.FormValue("ignore_time_and_sent_count_checks"))
		if ignoreTimeAndSentCountCheckReq == "1" {
			ignoreTimeAndSentCountCheck = true
		}
		grcftwc := database.ConfigFromTokenWithChecks(grToken, ignoreTimeAndSentCountCheck)
		if grcftwc.ClientID == 0 {
			// log.Printf("token %s does not meet criteria", grToken)
			// update stats
			database.UpdateStatsCanUseToken(0, grToken, false)
			w.Write(failedResponse)
			return
		}

		// TODO: maybe send this as one of own parameters e.g. gr_phone
		tel := strings.TrimSpace(req.FormValue(grcftwc.TelephoneParameter))
		// log.Printf("t param: %s\n", tel)
		// telephone := phonenumber.Parse(tel, grcftwc.Country)
		// // log.Printf("telephone: %s\n", telephone)
		// // HACK: The following was necessary for Dutch telephone numbers (may be others also) for phonenumber library to parse correctly.
		// // i.e. replacing the 00 prefix with a +
		// if telephone == "" && strings.HasPrefix(tel, "00") {
		// 	tel1 := strings.Replace(tel, "00", "+", 1)
		// 	telephone = phonenumber.Parse(tel1, grcftwc.Country)
		// }
		telephone := utils.TelephoneParse(tel, grcftwc.Country)
		// log.Printf("telephone: %s\n", telephone)
		if telephone == "" {
			log.Printf("no telephone found (sent telephone parameter: %s) for clientID: %d\n", tel, grcftwc.ClientID)
			// update stats
			database.UpdateStatsCanUseToken(grcftwc.ClientID, grToken, false)
			w.Write(failedResponse)
			return
		}
		// check barred telephone prefixes
		if barred.CheckBarred(telephone, Bars) {
			// log.Printf("telephone number is barred\n")
			// update stats
			database.UpdateStatsCanUseToken(grcftwc.ClientID, grToken, false)
			w.Write(failedResponse)
			return
		}
		// Some SIMs are configured not to send international numbers and when the telephone is
		// configured to E.164 format with the local country code this is determined to be international
		// so the SMS is not sent.
		// Therefore have to replace the country code to make it a national number this is normally with a 0.
		telephoneSendSMS := telephone
		if grcftwc.ReplaceTelephoneCountryCode {
			countryForTelephone := phonenumber.GetISO3166ByNumber(telephone, false)
			// log.Println(countryForTelephone.CountryCode)
			telephoneSendSMS = strings.Replace(telephone, countryForTelephone.CountryCode, grcftwc.ReplaceTelephoneCountryCodeWith, 1)
		}

		lastSent, sentCount, stop, found := database.LastSentFromTelephoneAndClient(telephone, grcftwc.ClientID)
		// check whether should ignore telephone checks (used for testing on front end)
		ignoreTelephoneChecks := strings.TrimSpace(req.FormValue("ignore_telephone_checks"))
		if ignoreTelephoneChecks != "1" {
			// replace default success response with the one from the database
			successResponseReplacement := []byte("")
			if string(grcftwc.SendSuccessResponse) != "EMPTY" {
				successResponseReplacement = []byte(grcftwc.SendSuccessResponse)
			}
			// check if stop set (do not send)
			if stop {
				// log.Printf("stop on telephone: %s\n", telephone)
				// w.Write(successResponse)
				w.Write(successResponseReplacement)
				return
			}
			// check found record
			if found {
				// check last sent greater than min send frequency
				if lastSent.After(time.Now().AddDate(0, 0, int(-grcftwc.MinSendFrequency))) {
					// log.Printf("Last sent too recent for telephone: %s\n", telephone)
					// update stats
					database.UpdateStatsCanUseToken(grcftwc.ClientID, grToken, false)
					// w.Write(successResponse)
					w.Write(successResponseReplacement)
					return
				}
				// check sent count
				if int(sentCount) > int(grcftwc.MaxSendCount) {
					// log.Printf("Reached maximum number of sends for telephone: %s\n", telephone)
					// update stats
					database.UpdateStatsCanUseToken(grcftwc.ClientID, grToken, false)
					// w.Write(successResponse)
					w.Write(successResponseReplacement)
					return
				}
			}
		}

		// get initial message
		message := ""
		if grcftwc.UseDatabaseMessage == 1 {
			message = grcftwc.Message
		} else {
			message = strings.TrimSpace(req.FormValue(grcftwc.MessageParameter))
		}

		// check for multi message
		if grcftwc.MultiMessageEnabled == 1 {
			// split message by the separator
			if message == "" {
				log.Printf("no message found for clientID: %d\n", grcftwc.ClientID)
				message = ""
			}
			sep := ""
			if message != "" {
				sep = strings.Trim(grcftwc.MultiMessageSeparator, " ")
				if sep == "" {
					log.Printf("no separator specified for multi message for clientID: %d\n", grcftwc.ClientID)
					message = ""
					sep = ""
				}
			}
			if message != "" && sep != "" {
				ms := strings.Split(message, sep)
				r := rand.Intn(len(ms))
				message = ms[r]
			}
			if message == "" {
				log.Printf("no message found for multi message after randomising found message array for clientID: %d\n", grcftwc.ClientID)
				message = ""
			}
		}

		// check whether should ignore dispatcher checks (used for testing on front end)
		ignoreDispatcherChecks := strings.TrimSpace(req.FormValue("ignore_dispatcher_checks"))

		// send message with parameters passed in this request removing own parameters
		params := req.PostForm
		params.Del("gr_token")
		params.Del("ignore_telephone_checks")
		params.Del("ignore_dispatcher_checks")
		params.Del("ignore_time_and_sent_count_checks")
		if message != "" {
			// send found message replacing the original message
			// NOTE: In some cases the message may be passed straight through without knowing what it is or what the message parameter is.
			params.Set(grcftwc.MessageParameter, message)
		}

		// set the send params to the following if sendFromIcabbi true:
		// app_key: from iCabbi
		// secret_key: from iCabbi
		// recipient: the telephone
		// body: the message
		sendMessageURL := grcftwc.SendURL
		if grcftwc.SendFromIcabbiApp && message != "" {
			params = url.Values{}
			params.Set("app_key", grcftwc.AppKey)
			params.Set("secret_key", grcftwc.SecretKey)
			params.Set("recipient", telephoneSendSMS)
			params.Set("body", message)

			// add correct API call to sendMessageURL
			if !strings.HasSuffix(sendMessageURL, "/") {
				sendMessageURL += "/"
			}
			sendMessageURL += "sms/add"
		}

		// set the send params to nil and create the following json if Review Master SMS Gateway enabled true:
		//
		// 		{"queue_id": "81", "message": "Hello world", "telephone": "+441234567890"}
		//
		// 		queue_id: this is the same as the client id, unless set to master queue, which is sent in the pairing process
		// 		telephone: the telephone (should be in E.164 format prepended with a + sign)
		// 		message: the message
		//
		// NOTE: header needs (this is done in the http client that sends the request):
		// 		api-token set
		// 		content-type: application/json
		// 		accept: application/json
		// Also need to set the sendMessageURL
		var body []byte
		if grcftwc.ReviewMasterSMSGatewayEnabled {
			reviewMasterSMSGatewayQueueID := grcftwc.ClientID
			if grcftwc.ReviewMasterSMSGatewayUseMasterQueue {
				reviewMasterSMSGatewayQueueID = uint64(database.ReviewMasterSMSMasterQueue)
			}
			body, _ = json.Marshal(map[string]string{
				"queue_id":  strconv.FormatUint(reviewMasterSMSGatewayQueueID, 10),
				"telephone": "+" + telephoneSendSMS,
				"message":   message,
			})
			sendMessageURL = config.Conf.ReviewMasterSMSGatewayURL
		}

		// check if send alternate message service
		sendJson := false
		apiKey := ""
		apiSecret := ""
		alternateSecret1 := ""
		if grcftwc.AlternateMessageServiceEnabled {
			if grcftwc.AlternateMessageService == "Message Media" {
				// use message media to send message (see: https://support.messagemedia.com/hc/en-us/articles/4413635760527-Messaging-API)
				// This uses json which is the main reason why need to create a specific set up, due to json body
				//
				// NOTE: Put the api_key and api_secret as parameters in the request.
				// 		These are passed as the appKey and secretKey to the client that makes the request to the message service
				// 		and is used for basic authentication in the header.
				//
				// NOTE: header needs (this is done in the http client that sends the request):
				// 		api-token set
				// 		content-type: application/json
				// 		accept: application/json
				type messageMediaMessage struct {
					Content           string `json:"content"`
					DestinationNumber string `json:"destination_number"`
					Format            string `json:"format"`
					DeliveryReport    string `json:"delivery_report"`
				}
				type messageMediaMessages struct {
					Messages []messageMediaMessage `json:"messages"`
				}
				ms := messageMediaMessages{
					Messages: []messageMediaMessage{
						{
							Content:           message,
							DestinationNumber: "+" + telephoneSendSMS,
							Format:            "SMS",
							DeliveryReport:    "true",
						},
					},
				}
				body, _ = json.Marshal(ms)
				apiKey = strings.TrimSpace(req.FormValue("api_key"))
				apiSecret = strings.TrimSpace(req.FormValue("api_secret"))
				params = nil
				sendJson = true
			} else if grcftwc.AlternateMessageService == "Veezu" {
				// use Veezu to send message (see: https://messages.veezu.com/api/messages)
				// This uses json which is the main reason why need to create a specific set up, due to json body
				//
				// NOTE: Put the alternate message service secret1 as parameters in the request.
				// 		This is used in the header as the auth_token.
				//
				// NOTE: header needs (this is done in the http client that sends the request):
				// 		auth_token set
				// 		content-type: application/json
				// 		accept: application/json
				type veezuMessage struct {
					Message   string `json:"message"`
					Telephone string `json:"telephone"`
				}
				ms := veezuMessage{
					Message:   message,
					Telephone: telephoneSendSMS,
				}
				body, _ = json.Marshal(ms)
				alternateSecret1 = strings.TrimSpace(grcftwc.AlternateMessageServiceSecret1)
				params = nil
				sendJson = true
			}
		}

		// see whether should do dispatcher checks
		// log.Println(dispatcherChecksEnabled, dispatcherURL, bookingIdParameter, isBookingForNowDiffMinutes, bookingNowPickupToContactMinutes, preBookingPickupToContactMinutes)
		if ignoreDispatcherChecks != "1" {
			if grcftwc.DispatcherChecksEnabled && grcftwc.DispatcherURL != "" && grcftwc.AppKey != "" && grcftwc.SecretKey != "" && grcftwc.BookingIdParameter != "" {
				// get the booking / trip ID
				tripID := strings.TrimSpace(req.FormValue(grcftwc.BookingIdParameter))
				dispatcherCheckPassed := BookingOk(grcftwc.DispatcherURL, grcftwc.AppKey, grcftwc.SecretKey, tripID,
					int(grcftwc.IsBookingForNowDiffMinutes), int(grcftwc.BookingNowPickupToContactMinutes), int(grcftwc.PreBookingPickupToContactMinutes),
					grcftwc.ClientID)
				if !dispatcherCheckPassed {
					// log.Printf("failed dispatcher test for clientID: %d, tripID: %s\n", clientID, tripID)
					// update stats
					database.UpdateStatsCanUseToken(grcftwc.ClientID, grToken, false)
					w.Write(failedResponse)
					return
				}
			}
		}

		// log.Printf("params: %v\n", params)
		httpMethod := "POST"
		if grcftwc.HttpGet {
			httpMethod = "GET"
		}

		var resp string

		// check if send later
		if grcftwc.SendDelayEnabled && grcftwc.SendDelay > 0 {
			// send later
			var headers map[string]string
			if grcftwc.ReviewMasterSMSGatewayEnabled {
				headers = map[string]string{
					"Content-Type": "application/json",
					"Accept":       "application/json",
					"Api-Token":    config.Conf.ReviewMasterSMSGatewayApiToken,
				}
			} else if grcftwc.AlternateMessageService == "Veezu" {
				headers = map[string]string{
					"auth_token":   alternateSecret1,
					"Content-Type": "application/json; charset=utf-8",
					"Accept":       "application/json",
				}
			} else if sendJson {
				headers = map[string]string{
					"Content-Type": "application/json",
					"Accept":       "application/json",
				}
			} else if httpMethod == "POST" {
				headers = map[string]string{
					"Content-Type": "application/x-www-form-urlencoded",
				}
			} else {
				headers = map[string]string{
					"Content-Type": "text/plain",
				}
			}
			// store request in database
			database.AddSendLater(telephone, grcftwc.ClientID, int(grcftwc.SendDelay),
				sendMessageURL, httpMethod, apiKey, apiSecret,
				headers, params, body, grcftwc.SendFromIcabbiApp,
				grcftwc.ReviewMasterSMSGatewayEnabled, grcftwc.AlternateMessageServiceEnabled,
				grcftwc.AlternateMessageService, false,
				grcftwc.SendSuccessResponse, grcftwc.MaxDailySendCount)

			// send success response
			resp = grcftwc.SendSuccessResponse
		} else {
			// send now
			// resp := client.Send(sendMessageURL, httpMethod, "", "", params)
			if grcftwc.ReviewMasterSMSGatewayEnabled {
				resp = client.Send(sendMessageURL, httpMethod, "", "", config.Conf.ReviewMasterSMSGatewayApiToken, nil, body, false, grcftwc.AlternateMessageService, "")
			} else {
				// resp = client.Send(sendMessageURL, httpMethod, "", "", "", params, nil, false)
				// resp = client.Send(sendMessageURL, httpMethod, apiKey, apiSecret, "", params, body, sendJson)
				resp = client.Send(sendMessageURL, httpMethod, apiKey, apiSecret, "", params, body, sendJson, grcftwc.AlternateMessageService, alternateSecret1)
			}
			// log.Printf("resp: %v\n", resp)

			// if sent from iCabbi App check response to make sure it has been sent successfully
			if grcftwc.SendFromIcabbiApp {
				// check sent successfully from iCabbi App
				// rJson := IcabbiAppResponse{}
				var rJson map[string]interface{}
				if errJson := json.Unmarshal([]byte(resp), &rJson); errJson != nil {
					log.Printf("Error unmarshalling JSON from sending message from iCabbi APP for clientID: %d, sending message to %s with params: %+v, response from send server: %+v, error: %+v", grcftwc.ClientID, sendMessageURL, params, resp, errJson)
				}
				// if rJson.Code != 0 {
				// type conversion is a string value "0" when successul else type int HTTP code (e.g. 404)
				code, ok := rJson["code"].(string)
				if !ok {
					log.Printf("Error sending message for clientID: %d to %s with params: %+v, response from send server: %+v", grcftwc.ClientID, sendMessageURL, params, resp)
				} else if code != "0" {
					// if rJson["code"].(string) != "0" {
					log.Printf("Error sending message for clientID: %d to %s with params: %+v, response from send server: %+v", grcftwc.ClientID, sendMessageURL, params, resp)
				} else {
					// set response to expected configured response which can be anything
					resp = grcftwc.SendSuccessResponse
				}
			}

			// if sent from Review Master SMS Gateway check response to make sure it has been sent successfully:
			// example possible responses:
			// 		{"id":76} - success with message id
			// 		{"error":"unauthorized"} - unsuccessful unauthorized
			// 		{"errors":{"message":["can't be blank"]}} - unsuccessful errors
			// 		{"errors":{"telephone":["can't be blank"]}} - unsuccessful errors
			// 		{"errors":{"detail":"Bad Request"}} - unsuccessful errors
			if grcftwc.ReviewMasterSMSGatewayEnabled {
				var rJson map[string]interface{}
				if errJson := json.Unmarshal([]byte(resp), &rJson); errJson != nil {
					log.Printf("Error unmarshalling JSON from sending message from Review Master SMS Gateway APP for clientID: %d, sending message to %s with params: %+v, response from send server: %+v, error: %+v", grcftwc.ClientID, sendMessageURL, params, resp, errJson)
				}
				if _, ok := rJson["id"]; !ok {
					log.Printf("Error sending message for clientID: %d to %s with params: %+v, response from send server: %+v", grcftwc.ClientID, sendMessageURL, params, resp)
				} else {
					// set response to expected configured response which can be anything
					resp = grcftwc.SendSuccessResponse
				}
			}

			// sent via alternate message service
			if grcftwc.AlternateMessageServiceEnabled {
				if grcftwc.AlternateMessageService == "Message Media" {
					// use message media to send message (see: https://support.messagemedia.com/hc/en-us/articles/4413635760527-Messaging-API)
					// example successful response:
					// {"messages":[{"callback_url":null,"delivery_report":true,"destination_number":"+447889525579","format":"SMS","message_expiry_timestamp":null,"message_flags":[],"message_id":"13d1a0ba-be11-401d-b4e8-7f27d8ccce99","metadata":null,"scheduled":null,"status":"queued","content":"testing","source_number":null,"rich_link":null,"media":null,"subject":null}]}
					type messageMediaMessageResponse struct {
						Content           string `json:"content"`
						DestinationNumber string `json:"destination_number"`
						Format            string `json:"format"`
						MessageID         string `json:"message_id"`
						Status            string `json:"status"`
					}
					type messageMediaMessagesResponse struct {
						Messages []messageMediaMessageResponse `json:"messages"`
					}
					var mmmr messageMediaMessagesResponse
					// NOTE: Only one message is sent at a time so there should only be one.
					json.Unmarshal([]byte(resp), &mmmr)
					if len(mmmr.Messages) < 1 {
						log.Printf("Error unmarshalling JSON from sending message from Message Media for clientID: %d, sending message to %s with body: %+v, response from send server: %+v", grcftwc.ClientID, sendMessageURL, body, resp)
					} else if len(mmmr.Messages[0].MessageID) < 1 {
						log.Printf("Error sending message for clientID: %d to %s with body: %+v, response from send server: %+v", grcftwc.ClientID, sendMessageURL, body, resp)
					} else if mmmr.Messages[0].Status != "queued" {
						log.Printf("Error sending message for clientID: %d to %s with body: %+v, response from send server: %+v", grcftwc.ClientID, sendMessageURL, body, resp)
					} else {
						// set response to expected configured response which can be anything
						resp = grcftwc.SendSuccessResponse
					}
				} else if grcftwc.AlternateMessageService == "Veezu" {
					// use Veezu to send message (see: https://messages.veezu.com/api/messages)
					// example successful response:
					// {"success":"1"}
					type veezuMessageResponse struct {
						Success string `json:"success"`
					}
					var vmr veezuMessageResponse
					json.Unmarshal([]byte(resp), &vmr)
					if vmr.Success != "1" {
						log.Printf("Error sending message for clientID: %d to %s with body: %+v, response from send server: %+v", grcftwc.ClientID, sendMessageURL, body, resp)
					} else {
						// set response to expected configured response which can be anything
						resp = grcftwc.SendSuccessResponse
					}
				}
			}

			// update last sent in database
			// NOTE: the word EMPTY is put in the database for the send_success_response field when an empty string is returned
			if (string(grcftwc.SendSuccessResponse) == "EMPTY" && resp == "") ||
				(strings.HasPrefix(strings.Trim(resp, " "), string(grcftwc.SendSuccessResponse))) {
				// log.Printf("updating last sent for telephone: %s\n", telephone)
				database.UpdateLastSent(telephone, grcftwc.ClientID, sentCount+1)
				// update stats
				database.UpdateStatsCanUseToken(grcftwc.ClientID, grToken, true)
			} else {
				log.Printf("Error sending message for clientID: %d to %s with params: %v, response from send server: %v", grcftwc.ClientID, sendMessageURL, params, resp)
				// update stats
				database.UpdateStatsCanUseToken(grcftwc.ClientID, grToken, false)
			}
		}

		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		// w.Write(successResponse)
		w.Write([]byte(resp))
	}

	return http.HandlerFunc(fn)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
