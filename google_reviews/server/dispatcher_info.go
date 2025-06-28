package server

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"google_reviews/client"
	"google_reviews/utils"
)

// BookingResponse
type BookingResponse struct {
	Code string      `json:"code"`
	Body BookingBody `json:"body"`
}

type BookingBody struct {
	Booking BookingInfo `json:"booking"`
}

type BookingInfo struct {
	Status      string `json:"status"`
	CreatedDate string `json:"created_date"`
	PickupDate  string `json:"pickup_date"`
	// ArriveDate  string `json:"arrive_date"`
	ContactDate string `json:"contact_date"`
}

// RetrieveBooking - get booking details from the dispatcher
func RetrieveBooking(dispatcherURL string, dispatcherAppKey string, dispatcherSecretKey string, tripID string) string {
	// add correct API call to dispatcherURL
	dURL := dispatcherURL
	if !strings.HasSuffix(dURL, "/") {
		dURL += "/"
	}
	dURL += "bookings/index/" + tripID

	resp := client.Send(dURL, "GET", dispatcherAppKey, dispatcherSecretKey, "", nil, nil, false, "", "")
	// fmt.Println("resp: ", resp)

	return resp
}

// CheckBooking - get booking details from the dispatcher and check against passed criteria
// return true passed criteria else false
func CheckBooking(dispatcherResp string, tripID string, isBookingForNowDiffMinutes int,
	bookingNowPickupToContactMinutes int, preBookingPickupToContactMinutes int,
	clientID uint64) bool {
	rJson := BookingResponse{}
	if errJson := json.Unmarshal([]byte(dispatcherResp), &rJson); errJson != nil {
		log.Printf("Error unmarshalling JSON for clientID: %d, booking id: %s from iCabbi API, response from server: %v, error: %v", clientID, tripID, dispatcherResp, errJson)
		return false
	}
	if rJson.Code != "0" {
		log.Printf("Error retrieving for clientID: %d, booking id: %s from iCabbi API, code: %s", clientID, tripID, rJson.Code)
		return false
	}
	// check status (COMPLETED)
	status := rJson.Body.Booking.Status
	if status != "COMPLETED" {
		return false
	}
	// determine whether booking is prebooked using a time difference between created and pickup time
	createdDate := rJson.Body.Booking.CreatedDate
	pickupDate := rJson.Body.Booking.PickupDate
	// arriveDate := rJson.Body.Booking.ArriveDate
	contactDate := rJson.Body.Booking.ContactDate
	bookingForNow := utils.CheckDiffTimeRFC3339(createdDate, pickupDate, strconv.Itoa(isBookingForNowDiffMinutes))
	if bookingForNow {
		return utils.CheckDiffTimeRFC3339(pickupDate, contactDate, strconv.Itoa(bookingNowPickupToContactMinutes))
	} else {
		return utils.CheckDiffTimeRFC3339(pickupDate, contactDate, strconv.Itoa(preBookingPickupToContactMinutes))
	}
}

// BookingOk - get booking details from the dispatcher and check against passed criteria
// return true passed criteria else false
func BookingOk(dispatcherURL string, dispatcherAppKey string, dispatcherSecretKey string, tripID string,
	isBookingForNowDiffMinutes int, bookingNowPickupToContactMinutes int, preBookingPickupToContactMinutes int,
	clientID uint64) bool {
	// add correct API call to dispatcherURL
	dURL := dispatcherURL
	if !strings.HasSuffix(dURL, "/") {
		dURL += "/"
	}
	dURL += "bookings/index/" + tripID

	resp := RetrieveBooking(dispatcherURL, dispatcherAppKey, dispatcherSecretKey, tripID)
	// fmt.Println("resp: ", resp)

	return CheckBooking(resp, tripID, isBookingForNowDiffMinutes, bookingNowPickupToContactMinutes, preBookingPickupToContactMinutes, clientID)
}
