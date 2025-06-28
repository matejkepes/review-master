package autocab_api_v2

import (
	"encoding/json"
	"google_reviews_autocab/autocab_api"
	"google_reviews_autocab/client"
	"log"
	"strings"
	"time"
)

type Company struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ArchivedBooking struct {
	PickedUpAtTime string `json:"pickedUpAtTime"`
	Reason         string `json:"reason"`
	// Updated to use booked by company rather than completed by company because some multi
	// companies don't have any cars on some companies that the job is booked by but the job
	// is transferred and completed by another one of their companies.
	// Also a customer knows the company they book with not necessarily the company that did the job.
	// Company        Company `json:"completedByCompanyID"`
	Company Company `json:"bookedByCompanyID"`
}

type Booking struct {
	TelephoneNumber string `json:"telephoneNumber"`
	BookedAtTime    string `json:"bookedAtTime"`
	PickupDueTime   string `json:"pickupDueTime"`
	BookingSource   string `json:"bookingSource"`
	ArchivedBooking ArchivedBooking
}

type BookingResponse struct {
	ContinuationToken string    `json:"continuationToken"`
	Bookings          []Booking `json:"bookings"`
}

type SearchBookingRequest struct {
	From              string   `json:"from"`
	To                string   `json:"to"`
	Types             []string `json:"types"`
	ContinuationToken string   `json:"continuationToken"`
}

// GetBookingFromResponse - get archive booking from JSON response
func GetBookingFromResponse(jsonResponse string) ([]Booking, string) {
	var bookingResponse BookingResponse
	err := json.Unmarshal([]byte(jsonResponse), &bookingResponse)
	if err != nil {
		log.Printf("Error decoding bookings JSON: %s, error: %+v\n", jsonResponse, err)
	}
	return bookingResponse.Bookings, bookingResponse.ContinuationToken
}

// GetBookingsFromServer - get the bookings from the Autocab server
func GetBookingsFromServer(serverURL, key string, from, to time.Time) []Booking {
	headers := map[string]string{
		"Content-Type":              "application/json",
		"Accept":                    "application/json",
		"Ocp-Apim-Subscription-Key": key,
	}

	body, _ := json.Marshal(SearchBookingRequest{
		From:              from.Format("2006/01/02 15:04"),
		To:                to.Format("2006/01/02 15:04"),
		Types:             []string{"Completed"},
		ContinuationToken: "",
	})

	apiURL := serverURL
	if !strings.HasSuffix(apiURL, "/") {
		apiURL += "/"
	}
	apiURL += "booking/v1/1.2/search"

	getMoreBookings := true
	bookings := make([]Booking, 0)
	for ok := true; ok; ok = getMoreBookings {
		log.Printf("request URL: %s, json body: %s\n", apiURL, body)
		resp := client.Send(apiURL, "POST", headers, nil, body)
		log.Println("resp: ", resp)
		bkings, continuationToken := GetBookingFromResponse(resp)
		bookings = append(bookings, bkings...)
		if continuationToken != "" {
			body, _ = json.Marshal(SearchBookingRequest{
				From:              from.Format("2006/01/02 15:04"),
				To:                to.Format("2006/01/02 15:04"),
				Types:             []string{"Completed"},
				ContinuationToken: continuationToken,
			})
		} else {
			getMoreBookings = false
		}
	}
	log.Printf("bookings: %+v\n", bookings)
	return bookings
}

// TranslateBookingsToArchiveBookings - translate bookings to archive bookings
// This makes it simpler so can use same code for both
func TranslateBookingsToArchiveBookings(bookings []Booking) []autocab_api.ArchiveBooking {
	var archiveBookings []autocab_api.ArchiveBooking
	for _, booking := range bookings {
		// log.Printf("booking: %+v\n", booking)
		var archiveBooking autocab_api.ArchiveBooking
		archiveBooking.TelephoneNumber = booking.TelephoneNumber
		archiveBooking.ArchiveReason = booking.ArchivedBooking.Reason
		archiveBooking.BookedAtTime = booking.BookedAtTime
		archiveBooking.PickupDueTime = booking.PickupDueTime
		archiveBooking.PickedUpAtTime = booking.ArchivedBooking.PickedUpAtTime
		archiveBooking.Company.ID = booking.ArchivedBooking.Company.ID
		archiveBooking.Company.Name = booking.ArchivedBooking.Company.Name
		archiveBooking.BookingSource = booking.BookingSource
		archiveBookings = append(archiveBookings, archiveBooking)
	}
	return archiveBookings
}
