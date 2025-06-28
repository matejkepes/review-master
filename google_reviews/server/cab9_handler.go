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

// cab 9 request parameters
const cab9GoogleReviewTokenParameter = "gr_token"
const cab9BookingCreationTimeParameter = "booking_creation_time"
const cab9BookedForTimeParameter = "booked_for_time"
const cab9PickedUpTimeParameter = "picked_up_time"

// cab 9 success response
var cab9SuccessResponse = []byte(`{"success":"1"}`)

// cab 9 failed response - only send if request parameters are wrong, otherwise send success response
var cab9FailedResponse = []byte(`{"success":"0"}`)

// Cab9Handler - cab 9 Google Reviews Handler
func Cab9Handler() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

		if err := req.ParseForm(); err != nil {
			// fmt.Printf("ParseForm() err: %v\n", err)
			log.Printf("ParseForm() err: %v\n", err)
			w.Write(cab9FailedResponse)
			return
		}
		// log.Printf("r.PostForm = %v\n", req.PostForm)

		// check google reviews token
		grToken := strings.TrimSpace(req.FormValue(cab9GoogleReviewTokenParameter))
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
			w.Write(cab9FailedResponse)
			return
		}
		// check database dispatcher_type is set to cab 9
		if grcftwc.DispatcherType != "CAB 9" {
			log.Printf("Dispatcher type set to: %s should be CAB 9 for clientID: %d", grcftwc.DispatcherType, grcftwc.ClientID)
			w.Write(cab9FailedResponse)
			return
		}

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
			log.Printf("telephone number %s is barred\n", telephone)
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

		// get initial message
		message := ""
		if grcftwc.UseDatabaseMessage == 1 {
			message = grcftwc.Message
		} else {
			message = strings.TrimSpace(req.FormValue(grcftwc.MessageParameter))
		}

		// check message is not empty
		if message == "" {
			log.Printf("no message sent in request or found in database for clientID: %d\n", grcftwc.ClientID)
			// update stats
			database.UpdateStatsCanUseToken(grcftwc.ClientID, grToken, false)
			w.Write(cab9SuccessResponse)
			return
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

		// see whether should do dispatcher checks
		// log.Println(dispatcherChecksEnabled, dispatcherURL, bookingIdParameter, isBookingForNowDiffMinutes, bookingNowPickupToContactMinutes, preBookingPickupToContactMinutes)
		if ignoreDispatcherChecks != "1" {
			bookingCreationTime := strings.TrimSpace(req.FormValue(cab9BookingCreationTimeParameter))
			if bookingCreationTime == "" {
				log.Printf("no booking creation time parameter sent in request for clientID: %d\n", grcftwc.ClientID)
				// update stats
				database.UpdateStatsCanUseToken(grcftwc.ClientID, grToken, false)
				w.Write(cab9FailedResponse)
				return
			}
			bookedForTime := strings.TrimSpace(req.FormValue(cab9BookedForTimeParameter))
			if bookedForTime == "" {
				log.Printf("no booked for time parameter sent in request for clientID: %d\n", grcftwc.ClientID)
				// update stats
				database.UpdateStatsCanUseToken(grcftwc.ClientID, grToken, false)
				w.Write(cab9FailedResponse)
				return
			}
			pickedUpTime := strings.TrimSpace(req.FormValue(cab9PickedUpTimeParameter))
			if pickedUpTime == "" {
				log.Printf("no picked up time parameter sent in request for clientID: %d\n", grcftwc.ClientID)
				// update stats
				database.UpdateStatsCanUseToken(grcftwc.ClientID, grToken, false)
				w.Write(cab9FailedResponse)
				return
			}

			bookingForNow := utils.CheckDiffTimeRFC3339(bookingCreationTime, bookedForTime, strconv.Itoa(int(grcftwc.IsBookingForNowDiffMinutes)))
			dispatcherCheckPassed := false
			if bookingForNow {
				dispatcherCheckPassed = utils.CheckDiffTimeRFC3339(bookedForTime, pickedUpTime, strconv.Itoa(int(grcftwc.BookingNowPickupToContactMinutes)))
			} else {
				dispatcherCheckPassed = utils.CheckDiffTimeRFC3339(bookedForTime, pickedUpTime, strconv.Itoa(int(grcftwc.PreBookingPickupToContactMinutes)))
			}

			if !dispatcherCheckPassed {
				// log.Printf("failed dispatcher test for clientID: %d\n", grcftwc.ClientID)
				// update stats
				database.UpdateStatsCanUseToken(grcftwc.ClientID, grToken, false)
				w.Write(cab9FailedResponse)
				return
			}
		}

		lastSent, sentCount, stop, found := database.LastSentFromTelephoneAndClient(telephone, grcftwc.ClientID)
		// check whether should ignore telephone checks (used for testing on front end)
		ignoreTelephoneChecks := strings.TrimSpace(req.FormValue("ignore_telephone_checks"))
		if ignoreTelephoneChecks != "1" {
			// check if stop set (do not send)
			if stop {
				log.Printf("stop on telephone: %s\n", telephone)
				// update stats
				database.UpdateStatsCanUseToken(grcftwc.ClientID, grToken, false)
				w.Write(cab9SuccessResponse)
				return
			}
			// check found record
			if found {
				// check last sent greater than min send frequency
				if lastSent.After(time.Now().AddDate(0, 0, int(-grcftwc.MinSendFrequency))) {
					log.Printf("Last sent too recent for telephone: %s\n", telephone)
					// update stats
					database.UpdateStatsCanUseToken(grcftwc.ClientID, grToken, false)
					w.Write(cab9SuccessResponse)
					return
				}
				// check sent count
				if int(sentCount) > int(grcftwc.MaxSendCount) {
					log.Printf("Reached maximum number of sends for telephone: %s\n", telephone)
					// update stats
					database.UpdateStatsCanUseToken(grcftwc.ClientID, grToken, false)
					w.Write(cab9SuccessResponse)
					return
				}
			}
		}

		// send SMS via Review Master SMS Gateway (currently the only option)
		//
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
		sendMessageURL := ""
		params := url.Values{}
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
			params = nil
		} else {
			log.Printf("Review Master SMS Gateway not enabled for clientID: %d\n", grcftwc.ClientID)
			// update stats
			database.UpdateStatsCanUseToken(grcftwc.ClientID, grToken, false)
			w.Write(cab9SuccessResponse)
			return
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
			// currently only option is to send SMS via Review Master SMS Gateway
			if grcftwc.ReviewMasterSMSGatewayEnabled {
				headers = map[string]string{
					"Content-Type": "application/json",
					"Accept":       "application/json",
					"Api-Token":    config.Conf.ReviewMasterSMSGatewayApiToken,
				}
			}
			// store request in database
			// database.AddSendLater(telephone, grcftwc.ClientID, int(grcftwc.SendDelay),
			// 	sendMessageURL, httpMethod, "", "",
			// 	headers, params, body, grcftwc.SendFromIcabbiApp,
			// 	grcftwc.ReviewMasterSMSGatewayEnabled, grcftwc.AlternateMessageServiceEnabled,
			// 	grcftwc.AlternateMessageService, false, string(cab9SuccessResponse), grcftwc.MaxDailySendCount)
			database.AddSendLater(telephone, grcftwc.ClientID, int(grcftwc.SendDelay),
				sendMessageURL, "POST", "", "",
				headers, params, body, false,
				grcftwc.ReviewMasterSMSGatewayEnabled, false,
				"", false, string(cab9SuccessResponse), grcftwc.MaxDailySendCount)

			// send success response
			resp = string(cab9SuccessResponse)
		} else {
			// send now
			// resp := client.Send(sendMessageURL, httpMethod, "", "", params)
			// currently only option is to send SMS via Review Master SMS Gateway
			if grcftwc.ReviewMasterSMSGatewayEnabled {
				resp = client.Send(sendMessageURL, httpMethod, "", "", config.Conf.ReviewMasterSMSGatewayApiToken, nil, body, false, "", "")
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
					resp = string(cab9SuccessResponse)
				}
			}

			// update last sent in database
			if resp == string(cab9SuccessResponse) {
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
