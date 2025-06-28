package autocab_api_v1

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

type SearchBookingRequest struct {
	From  string   `json:"from"`
	To    string   `json:"to"`
	Types []string `json:"types"`
}

type SendSMSRequest struct {
	SenderName      string   `json:"senderName"`
	RecipientTelnos []string `json:"recipientTelnos"`
	Message         string   `json:"message"`
}

type SendSMSResponse struct {
	SentRecipients   []string `json:"sentRecipients"`
	UnsentRecipients []string `json:"unsentRecipients"`
}

// GetBookingFromResponse - get archive booking from JSON response
func GetBookingFromResponse(jsonResponse string) []Booking {
	bookings := make([]Booking, 0)
	err := json.Unmarshal([]byte(jsonResponse), &bookings)
	if err != nil {
		log.Printf("Error decoding bookings JSON: %s, error: %+v\n", jsonResponse, err)
	}
	return bookings
}

// GetBookingsFromServer - get the bookings from the Autocab server
func GetBookingsFromServer(serverURL, key string, from, to time.Time) []Booking {
	headers := map[string]string{
		"Content-Type":              "application/json",
		"Accept":                    "application/json",
		"Ocp-Apim-Subscription-Key": key,
	}

	body, _ := json.Marshal(SearchBookingRequest{
		From:  from.Format("2006/01/02 15:04"),
		To:    to.Format("2006/01/02 15:04"),
		Types: []string{"Completed"},
	})

	apiURL := serverURL
	if !strings.HasSuffix(apiURL, "/") {
		apiURL += "/"
	}
	apiURL += "booking/v1/search"

	log.Printf("request URL: %s, json body: %s\n", apiURL, body)
	resp := client.Send(apiURL, "POST", headers, nil, body)
	log.Println("resp: ", resp)
	bookings := GetBookingFromResponse(resp)
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

// GetSendSMSResponse - get sent SMS from JSON response
func GetSendSMSResponse(jsonResponse string) SendSMSResponse {
	var sendSMSResponse SendSMSResponse
	err := json.Unmarshal([]byte(jsonResponse), &sendSMSResponse)
	if err != nil {
		log.Printf("Error decoding send SMS JSON: %s, error: %+v\n", jsonResponse, err)
	}
	return sendSMSResponse
}

// SendSMS - send SMS
// responds with success or failure
func SendSMS(serverURL, key, telephone, message, senderName string) bool {
	headers := map[string]string{
		"Content-Type":              "application/json",
		"Accept":                    "application/json",
		"Ocp-Apim-Subscription-Key": key,
	}

	body, _ := json.Marshal(SendSMSRequest{
		SenderName:      senderName,
		RecipientTelnos: []string{telephone},
		Message:         message,
	})

	apiURL := serverURL
	if !strings.HasSuffix(apiURL, "/") {
		apiURL += "/"
	}
	apiURL += "sms/v1/send"

	// For a successfully sent SMS (probably via SIMs since received SMS shows a telephone number) response appears to be empty i.e.:
	// response:
	// In other words it looks like the documentation may be incorrect and will not be able to decode SendSMSResponse json
	// For failed may get reason why e.g.:
	// response:  { "statusCode": 401, "message": "Access denied due to invalid subscription key. Make sure to provide a valid key for an active subscription." }

	// // serialize headers just to get raw hex for testing send later
	// h := new(bytes.Buffer)
	// he := gob.NewEncoder(h)
	// err := he.Encode(headers)
	// if err != nil {
	// 	h = nil
	// }
	// httpHeaders := h.Bytes()
	// log.Printf("request URL: %s, headers: %X, json body: %X\n", apiURL, httpHeaders, body)

	// log.Printf("request URL: %s, headers: %s, json body: %s\n", apiURL, headers, body)
	resp := client.Send(apiURL, "POST", headers, nil, body)
	// log.Println("resp: ", resp)

	telephoneSent := false
	if len(resp) == 0 {
		telephoneSent = true
	} else {
		// putting this here just incase there are any that comply with the documentation
		sentSMS := GetSendSMSResponse(resp)
		// log.Printf("sentSMS: %+v\n", sentSMS)
		for _, t := range sentSMS.SentRecipients {
			if t == telephone {
				telephoneSent = true
				break
			}
		}
	}
	if !telephoneSent {
		log.Printf("Error sending SMS via Autocab Send SMS to telephone: %s, message: %s, response: %s", telephone, message, resp)
	} else {
		log.Printf("Sent SMS via Autocab Send SMS to telephone: %s, message: %s, response: %s", telephone, message, resp)
	}

	return telephoneSent
}
