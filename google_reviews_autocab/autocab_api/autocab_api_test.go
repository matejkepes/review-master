package autocab_api

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestAuthorisationSuccess(t *testing.T) {
	var jsonResponse = []byte(`{"user":{"id":6,"name":"Digital"},"secret":"M5kuMJd+fzwT3XBS1H1biFDwVwCHNSH1gQqz3FuKDcuudjxelZsCRqtTGpuYI57/YtvqqsVZOjVi5AEVOmESTA=="}`)
	var authorisation Authorisation
	err := json.Unmarshal(jsonResponse, &authorisation)
	if err != nil {
		t.Fatalf("Error decoding authorisation, error: %+v\n", err)
	}
	fmt.Printf("authorisation: %+v\n", authorisation)
}

func TestAuthorisationFailed(t *testing.T) {
	var jsonResponse = []byte(``)
	var authorisation Authorisation
	err := json.Unmarshal(jsonResponse, &authorisation)
	if err == nil {
		t.Fatalf("Error decoding authorisation, error: %+v\n", err)
	}
}

func TestGetTokenFromAuthorisationResponseSuccess(t *testing.T) {
	var jsonResponse = string(`{"user":{"id":6,"name":"Digital"},"secret":"M5kuMJd+fzwT3XBS1H1biFDwVwCHNSH1gQqz3FuKDcuudjxelZsCRqtTGpuYI57/YtvqqsVZOjVi5AEVOmESTA=="}`)
	authorisationToken := GetTokenFromAuthorisationResponse(jsonResponse)
	if authorisationToken == "" {
		t.Fatal("Error getting authorisation token")
	}
	fmt.Printf("authorisationToken: %s\n", authorisationToken)
}

func TestGetTokenFromAuthorisationResponseFailed(t *testing.T) {
	var jsonResponse = string(``)
	authorisationToken := GetTokenFromAuthorisationResponse(jsonResponse)
	if authorisationToken != "" {
		t.Fatal("Getting authorisation token should have failed")
	}
}

// to live test server
func TestGetAuthorisationTokenFromServerLive(t *testing.T) {
	authorisationToken := GetAuthorisationTokenFromServer(testAutocabServerURL, testAutocabUsername, testAutocabPassword)
	if authorisationToken == "" {
		t.Fatal("Error getting authorisation token")
	}
	fmt.Printf("authorisationToken: %s\n", authorisationToken)
}

func TestArchiveBookingsSuccessAllCancelled(t *testing.T) {
	var jsonResponse = []byte(`[{"id":609,"bookingId":583,"docketNumber":0,"archiveTime":"2020-06-25T23:12:25+01:00","archiveReason":"Cancelled","pickupDueTime":"2020-06-25T23:10:39+01:00","dropOffDueTime":null,"pickup":{"address":"Royal Crescent, Cheadle, SK8 3BF","zone":{"id":120,"descriptor":"118","name":"Heald Green"},"coordinates":{"longitude":-2.219810277777778,"latitude":53.378326111111114}},"vias":null,"destination":{"address":"Cheadle Heath, Stockport, UK, Stockport","zone":{"id":119,"descriptor":"117","name":"Cheadle"},"coordinates":{"longitude":-2.1846388888888888,"latitude":53.39294}},"company":{"id":1,"name":"Driverspay Demo"},"account":null,"loyaltyCardID":null,"loyaltyCardCostValue":0.0,"accountType":null,"paymentType":"Cash","name":"","telephoneNumber":"","customerEmail":"","capabilities":null,"priority":0,"passengers":0,"luggage":0,"driverNote":"","officeNote":"","yourReferences":{"yourReference1":"","yourReference2":"","yourReference3":"","yourReference4":"","yourReference5":"","yourReference6":"","yourReference7":"","yourReference8":""},"ourReference":"","flightDetails":"","driver":{"id":2108,"callsign":"101","forename":"DG","surname":"DG","badgeNumber":""},"vehicle":{"id":1080,"callsign":"101","registration":"","plateNumber":""},"pricing":{"fare":0.0,"cost":0.0,"price":0.0,"waitingTime":0,"waitingTimeChargeable":0,"gratuityAmount":0.0,"costSource":"","pricingTariff":""},"distance":3.4,"systemDistance":3.4,"meterDistance":0.0,"typeOfBooking":"ASAP","bookingSource":"OperatorWeb","bookedAtTime":"2020-06-25T23:10:38.8809104+01:00","bookedBy":"admin","estimatedPickupTime":null,"dispatchedAtTime":"2020-06-25T23:11:01.2210988+01:00","pickedUpAtTime":null,"vehicleArrivedAtTime":null,"completedAtTime":null,"cabExchangeAgentBookingRef":""},{"id":610,"bookingId":586,"docketNumber":247,"archiveTime":"2020-06-25T23:29:48+01:00","archiveReason":"Cancelled","pickupDueTime":"2020-06-25T23:14:20+01:00","dropOffDueTime":null,"pickup":{"address":"Cheadle Royal Business Park, Cheadle, UK, Cheadle, SK8 3GP","zone":{"id":120,"descriptor":"118","name":"Heald Green"},"coordinates":{"longitude":-2.2241458333333335,"latitude":53.378203333333332}},"vias":null,"destination":{"address":"Manchester Airport (MAN), Manchester, M90 1QX","zone":{"id":127,"descriptor":"125","name":"M.I.A"},"coordinates":{"longitude":-2.2727302777777778,"latitude":53.3588025}},"company":{"id":1,"name":"Driverspay Demo"},"account":null,"loyaltyCardID":null,"loyaltyCardCostValue":0.0,"accountType":"Unknown","paymentType":"Card","name":"test","telephoneNumber":"447685932724","customerEmail":"","capabilities":null,"priority":0,"passengers":0,"luggage":0,"driverNote":"","officeNote":"","yourReferences":{"yourReference1":"","yourReference2":"","yourReference3":"","yourReference4":"","yourReference5":"","yourReference6":"","yourReference7":"","yourReference8":""},"ourReference":"","flightDetails":"","driver":{"id":2108,"callsign":"101","forename":"DG","surname":"DG","badgeNumber":""},"vehicle":{"id":1080,"callsign":"101","registration":"","plateNumber":""},"pricing":{"fare":0.0,"cost":0.0,"price":0.0,"waitingTime":0,"waitingTimeChargeable":0,"gratuityAmount":0.0,"costSource":"","pricingTariff":""},"distance":5.4,"systemDistance":5.4,"meterDistance":0.0,"typeOfBooking":"ASAP","bookingSource":"OperatorWeb","bookedAtTime":"2020-06-25T23:14:20.0373504+01:00","bookedBy":"admin","estimatedPickupTime":null,"dispatchedAtTime":"2020-06-25T23:14:52.3448569+01:00","pickedUpAtTime":null,"vehicleArrivedAtTime":null,"completedAtTime":null,"cabExchangeAgentBookingRef":""},{"id":612,"bookingId":588,"docketNumber":0,"archiveTime":"2020-06-26T13:56:54+01:00","archiveReason":"Cancelled","pickupDueTime":"2020-06-26T13:54:43+01:00","dropOffDueTime":null,"pickup":{"address":"1 Somerton Road, Bolton, BL2 6QA","zone":{"id":128,"descriptor":"126","name":"O O N"},"coordinates":{"longitude":-2.3753191666666669,"latitude":53.57869}},"vias":null,"destination":{"address":"2 Elsdon Drive, Manchester, M18 8WG","zone":{"id":69,"descriptor":"067","name":"Openshaw/Gorton"},"coordinates":{"longitude":-2.1697383333333335,"latitude":53.466091666666664}},"company":{"id":1,"name":"Driverspay Demo"},"account":null,"loyaltyCardID":null,"loyaltyCardCostValue":0.0,"accountType":null,"paymentType":"Cash","name":"Test","telephoneNumber":"","customerEmail":"","capabilities":null,"priority":0,"passengers":0,"luggage":0,"driverNote":"","officeNote":"","yourReferences":{"yourReference1":"","yourReference2":"","yourReference3":"","yourReference4":"","yourReference5":"","yourReference6":"","yourReference7":"","yourReference8":""},"ourReference":"","flightDetails":"","driver":null,"vehicle":null,"pricing":{"fare":0.0,"cost":0.0,"price":0.0,"waitingTime":0,"waitingTimeChargeable":0,"gratuityAmount":0.0,"costSource":"","pricingTariff":""},"distance":13.2,"systemDistance":13.2,"meterDistance":0.0,"typeOfBooking":"ASAP","bookingSource":"Operator","bookedAtTime":"2020-06-26T13:54:44.5176693+01:00","bookedBy":"Digital","estimatedPickupTime":null,"dispatchedAtTime":null,"pickedUpAtTime":null,"vehicleArrivedAtTime":null,"completedAtTime":null,"cabExchangeAgentBookingRef":""}]`)
	archiveBookings := make([]ArchiveBooking, 0)
	err := json.Unmarshal(jsonResponse, &archiveBookings)
	if err != nil {
		t.Fatalf("Error decoding archive bookings, error: %+v\n", err)
	}
	fmt.Printf("archive bookings: %+v\n", archiveBookings)
}

func TestArchiveBookingsNoEntries(t *testing.T) {
	var jsonResponse = []byte(`[]`)
	archiveBookings := make([]ArchiveBooking, 0)
	err := json.Unmarshal(jsonResponse, &archiveBookings)
	if err != nil {
		t.Fatalf("Error decoding archive bookings, error: %+v\n", err)
	}
	fmt.Printf("archive bookings: %+v\n", archiveBookings)
}

func TestArchiveBookingsSuccessAllCompletedMadeUpJustFieldsInterestedIn(t *testing.T) {
	var jsonResponse = []byte(`[{"archiveReason":"Completed","pickupDueTime":"2020-06-25T23:10:39+01:00","telephoneNumber":"447123456789","bookedAtTime":"2020-06-25T23:10:38.8809104+01:00","pickedUpAtTime":"2020-06-25T23:10:48.8809104+01:00","company":{"id":1,"name":"Driverspay Demo"},"bookingSource":"OperatorWeb"},{"archiveReason":"Completed","pickupDueTime":"2020-06-25T23:14:20+01:00","telephoneNumber":"447685932724","bookedAtTime":"2020-06-25T23:14:20.0373504+01:00","pickedUpAtTime":"2020-06-25T23:18:20.0373504+01:00"},{"pickupDueTime":"2020-06-26T13:54:43+01:00","telephoneNumber":"447123456788","bookedAtTime":"2020-06-26T13:54:44.5176693+01:00","pickedUpAtTime":"2020-06-26T13:58:44.5176693+01:00","company":{"id":1,"name":"Driverspay Demo"},"bookingSource":"OperatorWeb"}]`)
	archiveBookings := make([]ArchiveBooking, 0)
	err := json.Unmarshal(jsonResponse, &archiveBookings)
	if err != nil {
		t.Fatalf("Error decoding archive bookings, error: %+v\n", err)
	}
	fmt.Printf("archive bookings: %+v\n", archiveBookings)
}

func TestGetArchiveBookingFromResponseSuccessAllCompletedMadeUpJustFieldsInterestedIn(t *testing.T) {
	var jsonResponse = string(`[{"archiveReason":"Completed","pickupDueTime":"2020-06-25T23:10:39+01:00","telephoneNumber":"447123456789","bookedAtTime":"2020-06-25T23:10:38.8809104+01:00","pickedUpAtTime":"2020-06-25T23:10:48.8809104+01:00","company":{"id":1,"name":"Driverspay Demo"},"bookingSource":"OperatorWeb"},{"archiveReason":"Completed","pickupDueTime":"2020-06-25T23:14:20+01:00","telephoneNumber":"447685932724","bookedAtTime":"2020-06-25T23:14:20.0373504+01:00","pickedUpAtTime":"2020-06-25T23:18:20.0373504+01:00"},{"pickupDueTime":"2020-06-26T13:54:43+01:00","telephoneNumber":"447123456788","bookedAtTime":"2020-06-26T13:54:44.5176693+01:00","pickedUpAtTime":"2020-06-26T13:58:44.5176693+01:00","company":{"id":1,"name":"Driverspay Demo"},"bookingSource":"OperatorWeb"}]`)
	archiveBookings := GetArchiveBookingFromResponse(jsonResponse)
	if len(archiveBookings) == 0 {
		t.Fatal("Error getting archive bookings")
	}
	fmt.Printf("archive bookings: %+v\n", archiveBookings)
}

func TestGetArchiveBookingFromResponseNoEntries(t *testing.T) {
	var jsonResponse = string(`[]`)
	archiveBookings := GetArchiveBookingFromResponse(jsonResponse)
	if len(archiveBookings) != 0 {
		t.Fatal("Error getting archive bookings should be empty")
	}
	fmt.Printf("archive bookings: %+v\n", archiveBookings)
}

// to live test server
func TestGetArchiveBookingsFromServerLive(t *testing.T) {
	authorisationToken := GetAuthorisationTokenFromServer(testAutocabServerURL, testAutocabUsername, testAutocabPassword)
	if authorisationToken == "" {
		t.Fatal("Error getting authorisation token")
	}
	fmt.Printf("authorisationToken: %s\n", authorisationToken)

	// current time
	loc, err := time.LoadLocation("Europe/London")
	if err != nil {
		t.Fatal(err)
	}
	to := time.Now().In(loc)
	from := to.Add(-time.Hour * 24)
	archiveBookings := GetArchiveBookingsFromServer(testAutocabServerURL, authorisationToken, from, to)
	fmt.Printf("archive bookings: %+v\n", archiveBookings)
}
