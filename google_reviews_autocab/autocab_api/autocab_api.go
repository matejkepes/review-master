package autocab_api

import (
	"encoding/json"
	"log"
	"net/url"
	"strings"
	"time"

	"google_reviews_autocab/client"
)

type AuthorisationUser struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Authorisation struct {
	User   AuthorisationUser
	Secret string `json:"secret"`
}

type Company struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ArchiveBooking struct {
	TelephoneNumber string `json:"telephoneNumber"`
	ArchiveReason   string `json:"archiveReason"`
	BookedAtTime    string `json:"bookedAtTime"`
	PickupDueTime   string `json:"pickupDueTime"`
	// VehicleArrivedAtTime string `json:"vehicleArrivedAtTime"`
	PickedUpAtTime string `json:"pickedUpAtTime"`
	Company        Company
	BookingSource  string `json:"bookingSource"`
}

const BookingSourceMobileApp = "MobileApp"

// GetTokenFromAuthorisationResponse - get the token from authorisation JSON response
func GetTokenFromAuthorisationResponse(jsonResponse string) string {
	var authorisation Authorisation
	err := json.Unmarshal([]byte(jsonResponse), &authorisation)
	if err != nil {
		log.Printf("Error decoding authorisation JSON: %s, error: %+v\n", jsonResponse, err)
		return ""
	}
	return authorisation.Secret
}

// GetArchiveBookingFromResponse - get archive booking from JSON response
func GetArchiveBookingFromResponse(jsonResponse string) []ArchiveBooking {
	archiveBookings := make([]ArchiveBooking, 0)
	err := json.Unmarshal([]byte(jsonResponse), &archiveBookings)
	if err != nil {
		log.Printf("Error decoding archive bookings JSON: %s, error: %+v\n", jsonResponse, err)
	}
	return archiveBookings
}

// GetAuthorisationTokenFromServer - get the autorisation token from the Autocab server
func GetAuthorisationTokenFromServer(serverURL, username, password string) string {
	params := url.Values{}
	params.Add("username", username)
	params.Add("password", password)

	apiURL := serverURL
	if !strings.HasSuffix(apiURL, "/") {
		apiURL += "/"
	}
	apiURL += "api/thirdparty/v1/authenticate"

	resp := client.Send(apiURL, "POST", nil, params, nil)
	log.Println("resp: ", resp)
	token := GetTokenFromAuthorisationResponse(resp)
	log.Println("token: ", token)
	return token
}

// GetArchiveBookingsFromServer - get the archive bookings from the Autocab server
func GetArchiveBookingsFromServer(serverURL, token string, from, to time.Time) []ArchiveBooking {
	headers := map[string]string{"Authentication-Token": token}

	params := url.Values{}
	params.Add("from", from.Format("2006/01/02 15:04"))
	params.Add("to", to.Format("2006/01/02 15:04"))
	params.Add("ArchiveReasons", "Completed")

	apiURL := serverURL
	if !strings.HasSuffix(apiURL, "/") {
		apiURL += "/"
	}
	apiURL += "api/thirdparty/v1/archivedbookings"

	log.Printf("request URL: %s, parameters: %+v\n", apiURL, params)
	resp := client.Send(apiURL, "POST", headers, params, nil)
	log.Println("resp: ", resp)
	archiveBookings := GetArchiveBookingFromResponse(resp)
	log.Printf("archiveBookings: %+v\n", archiveBookings)
	return archiveBookings
}
