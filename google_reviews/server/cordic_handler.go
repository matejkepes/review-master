package server

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"google_reviews/database"
	"google_reviews/utils"
)

// cordic request parameters
const cordicGoogleReviewTokenParameter = "gr_token"
const cordicPassengerIDParameter = "passenger_id"
const cordicBookingCreationTimeParameter = "booking_creation_time"
const cordicBookedForTimeParameter = "booked_for_time"
const cordicPickedUpTimeParameter = "picked_up_time"
const cordicIgnorePassengerIDChecksParameter = "ignore_passenger_id_checks"

// cordic failed response
var cordicFailedResponse = []byte(`{"message":""}`)

// CordicHandler - Cordic Google Reviews Handler
func CordicHandler() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

		if err := req.ParseForm(); err != nil {
			// fmt.Printf("ParseForm() err: %v\n", err)
			log.Printf("ParseForm() err: %v\n", err)
			w.Write(cordicFailedResponse)
			return
		}
		// log.Printf("r.PostForm = %v\n", req.PostForm)

		// check google reviews token
		grToken := strings.TrimSpace(req.FormValue(cordicGoogleReviewTokenParameter))
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
			w.Write(cordicFailedResponse)
			return
		}
		// check database dispatcher_type is set to cordic
		if grcftwc.DispatcherType != "CORDIC" {
			log.Printf("Dispatcher type set to: %s should be CORDIC for clientID: %d", grcftwc.DispatcherType, grcftwc.ClientID)
			w.Write(cordicFailedResponse)
			return
		}

		// The cordic passenger identifier is unique and is treated like a telephone number
		passengerID := strings.TrimSpace(req.FormValue(cordicPassengerIDParameter))
		if passengerID == "" {
			log.Printf("no passenger ID parameter sent in request for clientID: %d\n", grcftwc.ClientID)
			// update stats
			database.UpdateStatsCanUseToken(grcftwc.ClientID, grToken, false)
			w.Write(cordicFailedResponse)
			return
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
			w.Write(cordicFailedResponse)
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
			bookingCreationTime := strings.TrimSpace(req.FormValue(cordicBookingCreationTimeParameter))
			if bookingCreationTime == "" {
				log.Printf("no booking creation time parameter sent in request for clientID: %d\n", grcftwc.ClientID)
				// update stats
				database.UpdateStatsCanUseToken(grcftwc.ClientID, grToken, false)
				w.Write(cordicFailedResponse)
				return
			}
			bookedForTime := strings.TrimSpace(req.FormValue(cordicBookedForTimeParameter))
			if bookedForTime == "" {
				log.Printf("no booked for time parameter sent in request for clientID: %d\n", grcftwc.ClientID)
				// update stats
				database.UpdateStatsCanUseToken(grcftwc.ClientID, grToken, false)
				w.Write(cordicFailedResponse)
				return
			}
			pickedUpTime := strings.TrimSpace(req.FormValue(cordicPickedUpTimeParameter))
			if pickedUpTime == "" {
				log.Printf("no picked up time parameter sent in request for clientID: %d\n", grcftwc.ClientID)
				// update stats
				database.UpdateStatsCanUseToken(grcftwc.ClientID, grToken, false)
				w.Write(cordicFailedResponse)
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
				// log.Printf("failed dispatcher test for clientID: %d, tripID: %s\n", clientID, tripID)
				// update stats
				database.UpdateStatsCanUseToken(grcftwc.ClientID, grToken, false)
				w.Write(cordicFailedResponse)
				return
			}
		}

		lastSent, sentCount, stop, found := database.LastSentFromTelephoneAndClient(passengerID, grcftwc.ClientID)
		// check whether should ignore passenger ID checks (used for testing on front end)
		ignorePassengerIDChecks := strings.TrimSpace(req.FormValue(cordicIgnorePassengerIDChecksParameter))
		if ignorePassengerIDChecks != "1" {
			// check if stop set (do not send)
			if stop {
				// log.Printf("stop on passenger ID: %s for clientID: %d\n", passengerID, clientID)
				w.Write(cordicFailedResponse)
				return
			}
			// check found record
			if found {
				// check last sent greater than min send frequency
				if lastSent.After(time.Now().AddDate(0, 0, int(-grcftwc.MinSendFrequency))) {
					// log.Printf("Last sent too recent for passenger ID: %s for clientID: %d\n", passengerID, clientID)
					// update stats
					database.UpdateStatsCanUseToken(grcftwc.ClientID, grToken, false)
					w.Write(cordicFailedResponse)
					return
				}
				// check sent count
				if int(sentCount) > int(grcftwc.MaxSendCount) {
					// log.Printf("Reached maximum number of sends for passenger ID: %s for clientID: %d\n", passengerID, clientID)
					// update stats
					database.UpdateStatsCanUseToken(grcftwc.ClientID, grToken, false)
					w.Write(cordicFailedResponse)
					return
				}
			}
		}

		// success
		// update last sent in database
		// log.Printf("updating last sent using passenger id for telephone: %s\n", passengerID)
		database.UpdateLastSent(passengerID, grcftwc.ClientID, sentCount+1)
		// update stats
		database.UpdateStatsCanUseToken(grcftwc.ClientID, grToken, true)
		// return the message to send
		w.Write([]byte(`{"message":"` + message + `"}`))
	}

	return http.HandlerFunc(fn)
}
