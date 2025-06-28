package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

type AuthorisationUser struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Authorisation struct {
	User   AuthorisationUser
	Secret string `json:"secret"`
}

func testServer(success bool, requestType string) *httptest.Server {
	// start a local HTTP server
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("req: ", req)
		// send response to be tested
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		var data string
		switch requestType {
		case "authorise":
			data = string(`{"user":{"id":6,"name":"Digital"},"secret":"M5kuMJd+fzwT3XBS1H1biFDwVwCHNSH1gQqz3FuKDcuudjxelZsCRqtTGpuYI57/YtvqqsVZOjVi5AEVOmESTA=="}`)
			if !success {
				data = ""
			}
		case "archives":
			// NOTE: This is a cancelled booking. Only interested in Completed bookings, just easier to create for a sample
			// data = string(`[{"id":609,"bookingId":583,"docketNumber":0,"archiveTime":"2020-06-25T23:12:25+01:00","archiveReason":"Cancelled","pickupDueTime":"2020-06-25T23:10:39+01:00","dropOffDueTime":null,"pickup":{"address":"Royal Crescent, Cheadle, SK8 3BF","zone":{"id":120,"descriptor":"118","name":"Heald Green"},"coordinates":{"longitude":-2.219810277777778,"latitude":53.378326111111114}},"vias":null,"destination":{"address":"Cheadle Heath, Stockport, UK, Stockport","zone":{"id":119,"descriptor":"117","name":"Cheadle"},"coordinates":{"longitude":-2.1846388888888888,"latitude":53.39294}},"company":{"id":1,"name":"Driverspay Demo"},"account":null,"loyaltyCardID":null,"loyaltyCardCostValue":0.0,"accountType":null,"paymentType":"Cash","name":"","telephoneNumber":"","customerEmail":"","capabilities":null,"priority":0,"passengers":0,"luggage":0,"driverNote":"","officeNote":"","yourReferences":{"yourReference1":"","yourReference2":"","yourReference3":"","yourReference4":"","yourReference5":"","yourReference6":"","yourReference7":"","yourReference8":""},"ourReference":"","flightDetails":"","driver":{"id":2108,"callsign":"101","forename":"DG","surname":"DG","badgeNumber":""},"vehicle":{"id":1080,"callsign":"101","registration":"","plateNumber":""},"pricing":{"fare":0.0,"cost":0.0,"price":0.0,"waitingTime":0,"waitingTimeChargeable":0,"gratuityAmount":0.0,"costSource":"","pricingTariff":""},"distance":3.4,"systemDistance":3.4,"meterDistance":0.0,"typeOfBooking":"ASAP","bookingSource":"OperatorWeb","bookedAtTime":"2020-06-25T23:10:38.8809104+01:00","bookedBy":"admin","estimatedPickupTime":null,"dispatchedAtTime":"2020-06-25T23:11:01.2210988+01:00","pickedUpAtTime":null,"vehicleArrivedAtTime":null,"completedAtTime":null,"cabExchangeAgentBookingRef":""},{"id":610,"bookingId":586,"docketNumber":247,"archiveTime":"2020-06-25T23:29:48+01:00","archiveReason":"Cancelled","pickupDueTime":"2020-06-25T23:14:20+01:00","dropOffDueTime":null,"pickup":{"address":"Cheadle Royal Business Park, Cheadle, UK, Cheadle, SK8 3GP","zone":{"id":120,"descriptor":"118","name":"Heald Green"},"coordinates":{"longitude":-2.2241458333333335,"latitude":53.378203333333332}},"vias":null,"destination":{"address":"Manchester Airport (MAN), Manchester, M90 1QX","zone":{"id":127,"descriptor":"125","name":"M.I.A"},"coordinates":{"longitude":-2.2727302777777778,"latitude":53.3588025}},"company":{"id":1,"name":"Driverspay Demo"},"account":null,"loyaltyCardID":null,"loyaltyCardCostValue":0.0,"accountType":"Unknown","paymentType":"Card","name":"test","telephoneNumber":"447685932724","customerEmail":"","capabilities":null,"priority":0,"passengers":0,"luggage":0,"driverNote":"","officeNote":"","yourReferences":{"yourReference1":"","yourReference2":"","yourReference3":"","yourReference4":"","yourReference5":"","yourReference6":"","yourReference7":"","yourReference8":""},"ourReference":"","flightDetails":"","driver":{"id":2108,"callsign":"101","forename":"DG","surname":"DG","badgeNumber":""},"vehicle":{"id":1080,"callsign":"101","registration":"","plateNumber":""},"pricing":{"fare":0.0,"cost":0.0,"price":0.0,"waitingTime":0,"waitingTimeChargeable":0,"gratuityAmount":0.0,"costSource":"","pricingTariff":""},"distance":5.4,"systemDistance":5.4,"meterDistance":0.0,"typeOfBooking":"ASAP","bookingSource":"OperatorWeb","bookedAtTime":"2020-06-25T23:14:20.0373504+01:00","bookedBy":"admin","estimatedPickupTime":null,"dispatchedAtTime":"2020-06-25T23:14:52.3448569+01:00","pickedUpAtTime":null,"vehicleArrivedAtTime":null,"completedAtTime":null,"cabExchangeAgentBookingRef":""},{"id":612,"bookingId":588,"docketNumber":0,"archiveTime":"2020-06-26T13:56:54+01:00","archiveReason":"Cancelled","pickupDueTime":"2020-06-26T13:54:43+01:00","dropOffDueTime":null,"pickup":{"address":"1 Somerton Road, Bolton, BL2 6QA","zone":{"id":128,"descriptor":"126","name":"O O N"},"coordinates":{"longitude":-2.3753191666666669,"latitude":53.57869}},"vias":null,"destination":{"address":"2 Elsdon Drive, Manchester, M18 8WG","zone":{"id":69,"descriptor":"067","name":"Openshaw/Gorton"},"coordinates":{"longitude":-2.1697383333333335,"latitude":53.466091666666664}},"company":{"id":1,"name":"Driverspay Demo"},"account":null,"loyaltyCardID":null,"loyaltyCardCostValue":0.0,"accountType":null,"paymentType":"Cash","name":"Test","telephoneNumber":"","customerEmail":"","capabilities":null,"priority":0,"passengers":0,"luggage":0,"driverNote":"","officeNote":"","yourReferences":{"yourReference1":"","yourReference2":"","yourReference3":"","yourReference4":"","yourReference5":"","yourReference6":"","yourReference7":"","yourReference8":""},"ourReference":"","flightDetails":"","driver":null,"vehicle":null,"pricing":{"fare":0.0,"cost":0.0,"price":0.0,"waitingTime":0,"waitingTimeChargeable":0,"gratuityAmount":0.0,"costSource":"","pricingTariff":""},"distance":13.2,"systemDistance":13.2,"meterDistance":0.0,"typeOfBooking":"ASAP","bookingSource":"Operator","bookedAtTime":"2020-06-26T13:54:44.5176693+01:00","bookedBy":"Digital","estimatedPickupTime":null,"dispatchedAtTime":null,"pickedUpAtTime":null,"vehicleArrivedAtTime":null,"completedAtTime":null,"cabExchangeAgentBookingRef":""}]`)
			// NOTE: This is made up with just the fields that am interested in
			// data = string(`[{"archiveReason":"Completed","pickupDueTime":"2020-06-25T23:10:39+01:00","telephoneNumber":"447123456789","bookedAtTime":"2020-06-25T23:10:38.8809104+01:00","pickedUpAtTime":"2020-06-25T23:10:48.8809104+01:00"},{"archiveReason":"Completed","pickupDueTime":"2020-06-25T23:14:20+01:00","telephoneNumber":"447685932724","bookedAtTime":"2020-06-25T23:14:20.0373504+01:00","pickedUpAtTime":"2020-06-25T23:18:20.0373504+01:00"},{"pickupDueTime":"2020-06-26T13:54:43+01:00","telephoneNumber":"447123456788","bookedAtTime":"2020-06-26T13:54:44.5176693+01:00","pickedUpAtTime":"2020-06-26T13:58:44.5176693+01:00"}]`)
			data = string(`[{"id":626,"bookingId":599,"docketNumber":259,"archiveTime":"2020-06-28T13:29:21.0673972+01:00","archiveReason":"Completed","pickupDueTime":"2020-06-28T13:27:21.3696155+01:00","dropOffDueTime":null,"pickup":{"address":"Ossy Tyres, Accrington, BB5 0EP","zone":{"id":132,"descriptor":"130","name":"National"},"coordinates":{"longitude":-2.389951388888889,"latitude":53.747869444444447}},"vias":null,"destination":{"address":"Accrington, UK, Lancashire","zone":{"id":132,"descriptor":"130","name":"National"},"coordinates":{"longitude":-2.37218,"latitude":53.753608888888891}},"company":{"id":1,"name":"Driverspay Demo"},"account":null,"loyaltyCardID":null,"loyaltyCardCostValue":0.0,"accountType":"Unknown","paymentType":"Card","name":"Andy","telephoneNumber":"07715527297","customerEmail":"","capabilities":null,"priority":0,"passengers":0,"luggage":0,"driverNote":"","officeNote":"","yourReferences":{"yourReference1":"","yourReference2":"","yourReference3":"","yourReference4":"","yourReference5":"","yourReference6":"","yourReference7":"","yourReference8":""},"ourReference":"","flightDetails":"","driver":{"id":2105,"callsign":"001","forename":"Andrew","surname":"Baxter","badgeNumber":""},"vehicle":{"id":1027,"callsign":"001","registration":"AD13 NOD","plateNumber":"NOD574R"},"pricing":{"fare":5.0,"cost":5.0,"price":5.0,"waitingTime":0,"waitingTimeChargeable":0,"gratuityAmount":0.0,"costSource":"Shout","pricingTariff":""},"distance":1.0,"systemDistance":1.0,"meterDistance":0.0,"typeOfBooking":"ASAP","bookingSource":"OperatorWeb","bookedAtTime":"2020-06-28T13:27:21.4046179+01:00","bookedBy":"Andy","estimatedPickupTime":"2020-06-28T13:32:19.1302782+01:00","dispatchedAtTime":"2020-06-28T13:27:23.5301663+01:00","pickedUpAtTime":"2020-06-28T13:28:37.0691283+01:00","vehicleArrivedAtTime":"2020-06-28T13:28:32.5975099+01:00","completedAtTime":"2020-06-28T13:29:21.0673972+01:00","cabExchangeAgentBookingRef":""},` +
				`{"id":627,"bookingId":600,"docketNumber":260,"archiveTime":"2020-06-28T13:32:31.750631+01:00","archiveReason":"Completed","pickupDueTime":"2020-06-28T13:30:20.6830325+01:00","dropOffDueTime":null,"pickup":{"address":"27 Wittlewood Drive, Accrington, BB5 5DJ","zone":{"id":132,"descriptor":"130","name":"National"},"coordinates":{"longitude":-2.372962777777778,"latitude":53.767095}},"vias":null,"destination":{"address":"Accrington, UK, Lancashire","zone":{"id":132,"descriptor":"130","name":"National"},"coordinates":{"longitude":-2.37218,"latitude":53.753608888888891}},"company":{"id":1,"name":"Driverspay Demo"},"account":null,"loyaltyCardID":null,"loyaltyCardCostValue":0.0,"accountType":"Unknown","paymentType":"Card","name":"c","telephoneNumber":"07715527297","customerEmail":"","capabilities":null,"priority":0,"passengers":0,"luggage":0,"driverNote":"","officeNote":"","yourReferences":{"yourReference1":"","yourReference2":"","yourReference3":"","yourReference4":"","yourReference5":"","yourReference6":"","yourReference7":"","yourReference8":""},"ourReference":"","flightDetails":"","driver":{"id":2105,"callsign":"001","forename":"Andrew","surname":"Baxter","badgeNumber":""},"vehicle":{"id":1027,"callsign":"001","registration":"AD13 NOD","plateNumber":"NOD574R"},"pricing":{"fare":10.0,"cost":10.0,"price":10.0,"waitingTime":0,"waitingTimeChargeable":0,"gratuityAmount":0.0,"costSource":"Shout","pricingTariff":""},"distance":1.4,"systemDistance":1.4,"meterDistance":0.0,"typeOfBooking":"ASAP","bookingSource":"OperatorWeb","bookedAtTime":"2020-06-28T13:30:20.7130882+01:00","bookedBy":"Andy","estimatedPickupTime":"2020-06-28T13:31:30.8485978+01:00","dispatchedAtTime":"2020-06-28T13:30:21.9986559+01:00","pickedUpAtTime":"2020-06-28T13:32:00.8404127+01:00","vehicleArrivedAtTime":"2020-06-28T13:31:50.4020964+01:00","completedAtTime":"2020-06-28T13:32:31.750631+01:00","cabExchangeAgentBookingRef":""},` +
				`{"id":628,"bookingId":601,"docketNumber":261,"archiveTime":"2020-06-28T13:34:43.6172213+01:00","archiveReason":"Completed","pickupDueTime":"2020-06-28T13:33:40.6168373+01:00","dropOffDueTime":null,"pickup":{"address":"Ossy Tyres, Accrington, BB5 0EP","zone":{"id":132,"descriptor":"130","name":"National"},"coordinates":{"longitude":-2.3899508333333332,"latitude":53.747868888888888}},"vias":null,"destination":{"address":"Accrington, UK, Lancashire","zone":{"id":132,"descriptor":"130","name":"National"},"coordinates":{"longitude":-2.37218,"latitude":53.753608888888891}},"company":{"id":1,"name":"Driverspay Demo"},"account":null,"loyaltyCardID":null,"loyaltyCardCostValue":0.0,"accountType":"Unknown","paymentType":"Card","name":"Andy","telephoneNumber":"07715527297","customerEmail":"","capabilities":null,"priority":0,"passengers":0,"luggage":0,"driverNote":"","officeNote":"","yourReferences":{"yourReference1":"","yourReference2":"","yourReference3":"","yourReference4":"","yourReference5":"","yourReference6":"","yourReference7":"","yourReference8":""},"ourReference":"","flightDetails":"","driver":{"id":2105,"callsign":"001","forename":"Andrew","surname":"Baxter","badgeNumber":""},"vehicle":{"id":1027,"callsign":"001","registration":"AD13 NOD","plateNumber":"NOD574R"},"pricing":{"fare":0.05,"cost":0.05,"price":0.05,"waitingTime":0,"waitingTimeChargeable":0,"gratuityAmount":0.0,"costSource":"Shout","pricingTariff":""},"distance":1.0,"systemDistance":1.0,"meterDistance":0.0,"typeOfBooking":"ASAP","bookingSource":"OperatorWeb","bookedAtTime":"2020-06-28T13:33:40.6418335+01:00","bookedBy":"Andy","estimatedPickupTime":"2020-06-28T13:38:46.5447248+01:00","dispatchedAtTime":"2020-06-28T13:33:41.5974097+01:00","pickedUpAtTime":"2020-06-28T13:33:53.8684533+01:00","vehicleArrivedAtTime":"2020-06-28T13:33:50.5498015+01:00","completedAtTime":"2020-06-28T13:34:43.6172213+01:00","cabExchangeAgentBookingRef":""},` +
				`{"id":629,"bookingId":602,"docketNumber":262,"archiveTime":"2020-06-28T13:50:29.2423057+01:00","archiveReason":"Completed","pickupDueTime":"2020-06-28T13:43:31.0119572+01:00","dropOffDueTime":null,"pickup":{"address":"27 Wittlewood Drive, Accrington, BB5 5DJ","zone":{"id":132,"descriptor":"130","name":"National"},"coordinates":{"longitude":-2.372962777777778,"latitude":53.767095}},"vias":null,"destination":{"address":"Accrington, UK, Lancashire","zone":{"id":132,"descriptor":"130","name":"National"},"coordinates":{"longitude":-2.37218,"latitude":53.753608888888891}},"company":{"id":1,"name":"Driverspay Demo"},"account":null,"loyaltyCardID":null,"loyaltyCardCostValue":0.0,"accountType":"Unknown","paymentType":"Card","name":"c","telephoneNumber":"07715527297","customerEmail":"","capabilities":null,"priority":0,"passengers":0,"luggage":0,"driverNote":"","officeNote":"","yourReferences":{"yourReference1":"","yourReference2":"","yourReference3":"","yourReference4":"","yourReference5":"","yourReference6":"","yourReference7":"","yourReference8":""},"ourReference":"","flightDetails":"","driver":{"id":2105,"callsign":"001","forename":"Andrew","surname":"Baxter","badgeNumber":""},"vehicle":{"id":1027,"callsign":"001","registration":"AD13 NOD","plateNumber":"NOD574R"},"pricing":{"fare":0.05,"cost":1.05,"price":1.05,"waitingTime":10,"waitingTimeChargeable":0,"gratuityAmount":0.0,"costSource":"Shout","pricingTariff":""},"distance":2.5,"systemDistance":1.4,"meterDistance":2.5,"typeOfBooking":"ASAP","bookingSource":"OperatorWeb","bookedAtTime":"2020-06-28T13:43:31.1969936+01:00","bookedBy":"Andy","estimatedPickupTime":"2020-06-28T13:43:33.9101872+01:00","dispatchedAtTime":"2020-06-28T13:43:32.9288928+01:00","pickedUpAtTime":"2020-06-28T13:49:32.9111614+01:00","vehicleArrivedAtTime":"2020-06-28T13:49:23.4956889+01:00","completedAtTime":"2020-06-28T13:50:29.2423057+01:00","cabExchangeAgentBookingRef":""},` +
				`{"id":630,"bookingId":603,"docketNumber":0,"archiveTime":"2020-06-28T14:08:51.167328+01:00","archiveReason":"Completed","pickupDueTime":"2020-06-28T14:05:11.7527201+01:00","dropOffDueTime":null,"pickup":{"address":"27 Wittlewood Drive, Accrington, BB5 5DJ","zone":{"id":132,"descriptor":"130","name":"National"},"coordinates":{"longitude":-2.372962777777778,"latitude":53.767095}},"vias":null,"destination":{"address":"Accrington, UK, Lancashire","zone":{"id":132,"descriptor":"130","name":"National"},"coordinates":{"longitude":-2.37218,"latitude":53.753608888888891}},"company":{"id":1,"name":"Driverspay Demo"},"account":null,"loyaltyCardID":null,"loyaltyCardCostValue":0.0,"accountType":null,"paymentType":"Cash","name":"c","telephoneNumber":"07715527297","customerEmail":"","capabilities":null,"priority":0,"passengers":0,"luggage":0,"driverNote":"","officeNote":"","yourReferences":{"yourReference1":"","yourReference2":"","yourReference3":"","yourReference4":"","yourReference5":"","yourReference6":"","yourReference7":"","yourReference8":""},"ourReference":"","flightDetails":"","driver":{"id":2105,"callsign":"001","forename":"Andrew","surname":"Baxter","badgeNumber":""},"vehicle":{"id":1027,"callsign":"001","registration":"AD13 NOD","plateNumber":"NOD574R"},"pricing":{"fare":0.0,"cost":0.0,"price":0.0,"waitingTime":0,"waitingTimeChargeable":0,"gratuityAmount":0.0,"costSource":"Shout","pricingTariff":""},"distance":1.4,"systemDistance":1.4,"meterDistance":0.0,"typeOfBooking":"ASAP","bookingSource":"OperatorWeb","bookedAtTime":"2020-06-28T14:05:11.8413076+01:00","bookedBy":"Andy","estimatedPickupTime":"2020-06-28T14:05:14.1594982+01:00","dispatchedAtTime":"2020-06-28T14:05:12.9744395+01:00","pickedUpAtTime":"2020-06-28T14:08:45.7087884+01:00","vehicleArrivedAtTime":"2020-06-28T14:08:42.3445656+01:00","completedAtTime":"2020-06-28T14:08:51.167328+01:00","cabExchangeAgentBookingRef":""},` +
				`{"id":631,"bookingId":604,"docketNumber":264,"archiveTime":"2020-06-28T14:11:23.9332623+01:00","archiveReason":"Completed","pickupDueTime":"2020-06-28T14:09:33.329504+01:00","dropOffDueTime":null,"pickup":{"address":"27 Wittlewood Drive, Accrington, BB5 5DJ","zone":{"id":132,"descriptor":"130","name":"National"},"coordinates":{"longitude":-2.372962777777778,"latitude":53.767095}},"vias":null,"destination":{"address":"Accrington, UK, Lancashire","zone":{"id":132,"descriptor":"130","name":"National"},"coordinates":{"longitude":-2.37218,"latitude":53.753608888888891}},"company":{"id":1,"name":"Driverspay Demo"},"account":null,"loyaltyCardID":null,"loyaltyCardCostValue":0.0,"accountType":"Unknown","paymentType":"Card","name":"c","telephoneNumber":"07715527297","customerEmail":"","capabilities":null,"priority":0,"passengers":0,"luggage":0,"driverNote":"","officeNote":"","yourReferences":{"yourReference1":"","yourReference2":"","yourReference3":"","yourReference4":"","yourReference5":"","yourReference6":"","yourReference7":"","yourReference8":""},"ourReference":"","flightDetails":"","driver":{"id":2105,"callsign":"001","forename":"Andrew","surname":"Baxter","badgeNumber":""},"vehicle":{"id":1027,"callsign":"001","registration":"AD13 NOD","plateNumber":"NOD574R"},"pricing":{"fare":0.05,"cost":0.05,"price":0.05,"waitingTime":0,"waitingTimeChargeable":0,"gratuityAmount":0.0,"costSource":"Shout","pricingTariff":""},"distance":1.4,"systemDistance":1.4,"meterDistance":0.0,"typeOfBooking":"ASAP","bookingSource":"OperatorWeb","bookedAtTime":"2020-06-28T14:09:33.4562645+01:00","bookedBy":"Andy","estimatedPickupTime":"2020-06-28T14:09:35.1824348+01:00","dispatchedAtTime":"2020-06-28T14:09:34.477472+01:00","pickedUpAtTime":"2020-06-28T14:10:27.8827447+01:00","vehicleArrivedAtTime":"2020-06-28T14:10:24.0051819+01:00","completedAtTime":"2020-06-28T14:11:23.9283284+01:00","cabExchangeAgentBookingRef":""},` +
				`{"id":632,"bookingId":605,"docketNumber":265,"archiveTime":"2020-06-28T14:13:50.9951973+01:00","archiveReason":"Completed","pickupDueTime":"2020-06-28T14:12:16.1651572+01:00","dropOffDueTime":null,"pickup":{"address":"27 Wittlewood Drive, Accrington, BB5 5DJ","zone":{"id":132,"descriptor":"130","name":"National"},"coordinates":{"longitude":-2.372962777777778,"latitude":53.767095}},"vias":null,"destination":{"address":"Accrington, UK, Lancashire","zone":{"id":132,"descriptor":"130","name":"National"},"coordinates":{"longitude":-2.37218,"latitude":53.753608888888891}},"company":{"id":1,"name":"Driverspay Demo"},"account":null,"loyaltyCardID":null,"loyaltyCardCostValue":0.0,"accountType":"Unknown","paymentType":"Card","name":"c","telephoneNumber":"07715527297","customerEmail":"","capabilities":null,"priority":0,"passengers":0,"luggage":0,"driverNote":"","officeNote":"","yourReferences":{"yourReference1":"","yourReference2":"","yourReference3":"","yourReference4":"","yourReference5":"","yourReference6":"","yourReference7":"","yourReference8":""},"ourReference":"","flightDetails":"","driver":{"id":2105,"callsign":"001","forename":"Andrew","surname":"Baxter","badgeNumber":""},"vehicle":{"id":1027,"callsign":"001","registration":"AD13 NOD","plateNumber":"NOD574R"},"pricing":{"fare":0.5,"cost":0.5,"price":0.5,"waitingTime":0,"waitingTimeChargeable":0,"gratuityAmount":0.0,"costSource":"Shout","pricingTariff":""},"distance":1.4,"systemDistance":1.4,"meterDistance":0.0,"typeOfBooking":"ASAP","bookingSource":"OperatorWeb","bookedAtTime":"2020-06-28T14:12:16.3252341+01:00","bookedBy":"Andy","estimatedPickupTime":"2020-06-28T14:12:18.3216565+01:00","dispatchedAtTime":"2020-06-28T14:12:17.5252892+01:00","pickedUpAtTime":"2020-06-28T14:13:04.3432805+01:00","vehicleArrivedAtTime":"2020-06-28T14:13:00.8408371+01:00","completedAtTime":"2020-06-28T14:13:50.9951973+01:00","cabExchangeAgentBookingRef":""},` +
				`{"id":633,"bookingId":606,"docketNumber":0,"archiveTime":"2020-06-28T15:26:56.94929+01:00","archiveReason":"Completed","pickupDueTime":"2020-06-28T15:25:41.7442367+01:00","dropOffDueTime":null,"pickup":{"address":"27 Wittlewood Drive, Accrington, BB5 5DJ","zone":{"id":132,"descriptor":"130","name":"National"},"coordinates":{"longitude":-2.372962777777778,"latitude":53.767095}},"vias":null,"destination":{"address":"Accrington, UK, Lancashire","zone":{"id":132,"descriptor":"130","name":"National"},"coordinates":{"longitude":-2.37218,"latitude":53.753608888888891}},"company":{"id":1,"name":"Driverspay Demo"},"account":null,"loyaltyCardID":null,"loyaltyCardCostValue":0.0,"accountType":null,"paymentType":"Cash","name":"c","telephoneNumber":"07715527297","customerEmail":"","capabilities":null,"priority":0,"passengers":0,"luggage":0,"driverNote":"","officeNote":"","yourReferences":{"yourReference1":"","yourReference2":"","yourReference3":"","yourReference4":"","yourReference5":"","yourReference6":"","yourReference7":"","yourReference8":""},"ourReference":"","flightDetails":"","driver":{"id":2105,"callsign":"001","forename":"Andrew","surname":"Baxter","badgeNumber":""},"vehicle":{"id":1027,"callsign":"001","registration":"AD13 NOD","plateNumber":"NOD574R"},"pricing":{"fare":0.0,"cost":0.0,"price":0.0,"waitingTime":0,"waitingTimeChargeable":0,"gratuityAmount":0.0,"costSource":"Shout","pricingTariff":""},"distance":1.4,"systemDistance":1.4,"meterDistance":0.0,"typeOfBooking":"ASAP","bookingSource":"OperatorWeb","bookedAtTime":"2020-06-28T15:25:41.8092453+01:00","bookedBy":"Andy","estimatedPickupTime":"2020-06-28T15:26:52.8427846+01:00","dispatchedAtTime":"2020-06-28T15:25:43.7628424+01:00","pickedUpAtTime":"2020-06-28T15:26:54.2936365+01:00","vehicleArrivedAtTime":"2020-06-28T15:26:51.5878772+01:00","completedAtTime":"2020-06-28T15:26:56.94929+01:00","cabExchangeAgentBookingRef":""},` +
				`{"id":634,"bookingId":607,"docketNumber":267,"archiveTime":"2020-06-28T15:28:52.0526348+01:00","archiveReason":"Completed","pickupDueTime":"2020-06-28T15:27:23.8454349+01:00","dropOffDueTime":null,"pickup":{"address":"27 Wittlewood Drive, Accrington, BB5 5DJ","zone":{"id":132,"descriptor":"130","name":"National"},"coordinates":{"longitude":-2.372962777777778,"latitude":53.767095}},"vias":null,"destination":{"address":"Accrington, UK, Lancashire","zone":{"id":132,"descriptor":"130","name":"National"},"coordinates":{"longitude":-2.37218,"latitude":53.753608888888891}},"company":{"id":1,"name":"Driverspay Demo"},"account":null,"loyaltyCardID":null,"loyaltyCardCostValue":0.0,"accountType":"Unknown","paymentType":"Card","name":"c","telephoneNumber":"07715527297","customerEmail":"","capabilities":null,"priority":0,"passengers":0,"luggage":0,"driverNote":"","officeNote":"","yourReferences":{"yourReference1":"","yourReference2":"","yourReference3":"","yourReference4":"","yourReference5":"","yourReference6":"","yourReference7":"","yourReference8":""},"ourReference":"","flightDetails":"","driver":{"id":2105,"callsign":"001","forename":"Andrew","surname":"Baxter","badgeNumber":""},"vehicle":{"id":1027,"callsign":"001","registration":"AD13 NOD","plateNumber":"NOD574R"},"pricing":{"fare":0.55,"cost":0.55,"price":0.55,"waitingTime":0,"waitingTimeChargeable":0,"gratuityAmount":0.0,"costSource":"Shout","pricingTariff":""},"distance":1.4,"systemDistance":1.4,"meterDistance":0.0,"typeOfBooking":"ASAP","bookingSource":"OperatorWeb","bookedAtTime":"2020-06-28T15:27:23.8704364+01:00","bookedBy":"Andy","estimatedPickupTime":"2020-06-28T15:28:33.5040768+01:00","dispatchedAtTime":"2020-06-28T15:27:24.4291364+01:00","pickedUpAtTime":"2020-06-28T15:28:22.825393+01:00","vehicleArrivedAtTime":"2020-06-28T15:28:18.2061259+01:00","completedAtTime":"2020-06-28T15:28:52.0526348+01:00","cabExchangeAgentBookingRef":""},` +
				`{"id":636,"bookingId":609,"docketNumber":0,"archiveTime":"2020-06-28T17:13:01.5056108+01:00","archiveReason":"Completed","pickupDueTime":"2020-06-28T17:12:13.7222659+01:00","dropOffDueTime":null,"pickup":{"address":"27 Wittlewood Drive, Accrington, BB5 5DJ","zone":{"id":132,"descriptor":"130","name":"National"},"coordinates":{"longitude":-2.372962777777778,"latitude":53.767095}},"vias":null,"destination":{"address":"Accrington, UK, Lancashire","zone":{"id":132,"descriptor":"130","name":"National"},"coordinates":{"longitude":-2.37218,"latitude":53.753608888888891}},"company":{"id":1,"name":"Driverspay Demo"},"account":null,"loyaltyCardID":null,"loyaltyCardCostValue":0.0,"accountType":null,"paymentType":"Cash","name":"c","telephoneNumber":"07715527297","customerEmail":"","capabilities":null,"priority":0,"passengers":0,"luggage":0,"driverNote":"","officeNote":"","yourReferences":{"yourReference1":"","yourReference2":"","yourReference3":"","yourReference4":"","yourReference5":"","yourReference6":"","yourReference7":"","yourReference8":""},"ourReference":"","flightDetails":"","driver":{"id":2105,"callsign":"001","forename":"Andrew","surname":"Baxter","badgeNumber":""},"vehicle":{"id":1027,"callsign":"001","registration":"AD13 NOD","plateNumber":"NOD574R"},"pricing":{"fare":0.0,"cost":0.0,"price":0.0,"waitingTime":0,"waitingTimeChargeable":0,"gratuityAmount":0.0,"costSource":"Shout","pricingTariff":""},"distance":1.4,"systemDistance":1.4,"meterDistance":0.0,"typeOfBooking":"ASAP","bookingSource":"OperatorWeb","bookedAtTime":"2020-06-28T17:12:13.8677714+01:00","bookedBy":"Andy","estimatedPickupTime":"2020-06-28T17:13:24.5792601+01:00","dispatchedAtTime":"2020-06-28T17:12:15.1970364+01:00","pickedUpAtTime":"2020-06-28T17:12:59.1176905+01:00","vehicleArrivedAtTime":"2020-06-28T17:12:55.5399221+01:00","completedAtTime":"2020-06-28T17:13:01.5056108+01:00","cabExchangeAgentBookingRef":""},` +
				`{"id":637,"bookingId":610,"docketNumber":0,"archiveTime":"2020-06-28T17:14:16.6100955+01:00","archiveReason":"Completed","pickupDueTime":"2020-06-28T17:13:15.7422387+01:00","dropOffDueTime":null,"pickup":{"address":"27 Wittlewood Drive, Accrington, BB5 5DJ","zone":{"id":132,"descriptor":"130","name":"National"},"coordinates":{"longitude":-2.372962777777778,"latitude":53.767095}},"vias":null,"destination":{"address":"Accrington, UK, Lancashire","zone":{"id":132,"descriptor":"130","name":"National"},"coordinates":{"longitude":-2.37218,"latitude":53.753608888888891}},"company":{"id":1,"name":"Driverspay Demo"},"account":null,"loyaltyCardID":null,"loyaltyCardCostValue":0.0,"accountType":null,"paymentType":"Cash","name":"c","telephoneNumber":"07715527297","customerEmail":"","capabilities":null,"priority":0,"passengers":0,"luggage":0,"driverNote":"","officeNote":"","yourReferences":{"yourReference1":"","yourReference2":"","yourReference3":"","yourReference4":"","yourReference5":"","yourReference6":"","yourReference7":"","yourReference8":""},"ourReference":"","flightDetails":"","driver":{"id":2105,"callsign":"001","forename":"Andrew","surname":"Baxter","badgeNumber":""},"vehicle":{"id":1027,"callsign":"001","registration":"AD13 NOD","plateNumber":"NOD574R"},"pricing":{"fare":30.0,"cost":30.0,"price":30.0,"waitingTime":0,"waitingTimeChargeable":0,"gratuityAmount":0.0,"costSource":"Shout","pricingTariff":""},"distance":1.4,"systemDistance":1.4,"meterDistance":0.0,"typeOfBooking":"ASAP","bookingSource":"OperatorWeb","bookedAtTime":"2020-06-28T17:13:15.9042087+01:00","bookedBy":"Andy","estimatedPickupTime":"2020-06-28T17:14:26.571544+01:00","dispatchedAtTime":"2020-06-28T17:13:16.9715773+01:00","pickedUpAtTime":"2020-06-28T17:14:02.8619933+01:00","vehicleArrivedAtTime":"2020-06-28T17:13:59.6393845+01:00","completedAtTime":"2020-06-28T17:14:16.6100955+01:00","cabExchangeAgentBookingRef":""},` +
				`{"id":638,"bookingId":611,"docketNumber":0,"archiveTime":"2020-06-28T17:16:41.0923121+01:00","archiveReason":"Completed","pickupDueTime":"2020-06-28T17:15:09.028227+01:00","dropOffDueTime":null,"pickup":{"address":"27 Wittlewood Drive, Accrington, BB5 5DJ","zone":{"id":132,"descriptor":"130","name":"National"},"coordinates":{"longitude":-2.372962777777778,"latitude":53.767095}},"vias":null,"destination":{"address":"Accrington, UK, Lancashire","zone":{"id":132,"descriptor":"130","name":"National"},"coordinates":{"longitude":-2.37218,"latitude":53.753608888888891}},"company":{"id":1,"name":"Driverspay Demo"},"account":null,"loyaltyCardID":null,"loyaltyCardCostValue":0.0,"accountType":null,"paymentType":"Cash","name":"c","telephoneNumber":"07715527297","customerEmail":"","capabilities":null,"priority":0,"passengers":0,"luggage":0,"driverNote":"","officeNote":"","yourReferences":{"yourReference1":"","yourReference2":"","yourReference3":"","yourReference4":"","yourReference5":"","yourReference6":"","yourReference7":"","yourReference8":""},"ourReference":"","flightDetails":"","driver":{"id":2105,"callsign":"001","forename":"Andrew","surname":"Baxter","badgeNumber":""},"vehicle":{"id":1027,"callsign":"001","registration":"AD13 NOD","plateNumber":"NOD574R"},"pricing":{"fare":0.0,"cost":0.0,"price":0.0,"waitingTime":0,"waitingTimeChargeable":0,"gratuityAmount":0.0,"costSource":"Shout","pricingTariff":""},"distance":1.4,"systemDistance":1.4,"meterDistance":0.0,"typeOfBooking":"ASAP","bookingSource":"OperatorWeb","bookedAtTime":"2020-06-28T17:15:09.0533047+01:00","bookedBy":"Andy","estimatedPickupTime":"2020-06-28T17:16:19.5736905+01:00","dispatchedAtTime":"2020-06-28T17:15:10.1576804+01:00","pickedUpAtTime":"2020-06-28T17:16:39.2289048+01:00","vehicleArrivedAtTime":"2020-06-28T17:16:10.3954418+01:00","completedAtTime":"2020-06-28T17:16:41.0923121+01:00","cabExchangeAgentBookingRef":""},` +
				`{"id":639,"bookingId":612,"docketNumber":0,"archiveTime":"2020-06-28T17:55:33.1802975+01:00","archiveReason":"Completed","pickupDueTime":"2020-06-28T17:17:01.3268246+01:00","dropOffDueTime":null,"pickup":{"address":"27 Wittlewood Drive, Accrington, BB5 5DJ","zone":{"id":132,"descriptor":"130","name":"National"},"coordinates":{"longitude":-2.372962777777778,"latitude":53.767095}},"vias":null,"destination":{"address":"Accrington, UK, Lancashire","zone":{"id":132,"descriptor":"130","name":"National"},"coordinates":{"longitude":-2.37218,"latitude":53.753608888888891}},"company":{"id":1,"name":"Driverspay Demo"},"account":null,"loyaltyCardID":null,"loyaltyCardCostValue":0.0,"accountType":null,"paymentType":"Cash","name":"c","telephoneNumber":"07715527297","customerEmail":"","capabilities":null,"priority":0,"passengers":0,"luggage":0,"driverNote":"","officeNote":"","yourReferences":{"yourReference1":"","yourReference2":"","yourReference3":"","yourReference4":"","yourReference5":"","yourReference6":"","yourReference7":"","yourReference8":""},"ourReference":"","flightDetails":"","driver":{"id":2105,"callsign":"001","forename":"Andrew","surname":"Baxter","badgeNumber":""},"vehicle":{"id":1027,"callsign":"001","registration":"AD13 NOD","plateNumber":"NOD574R"},"pricing":{"fare":0.0,"cost":0.0,"price":0.0,"waitingTime":0,"waitingTimeChargeable":0,"gratuityAmount":0.0,"costSource":"Shout","pricingTariff":""},"distance":1.4,"systemDistance":1.4,"meterDistance":0.0,"typeOfBooking":"ASAP","bookingSource":"OperatorWeb","bookedAtTime":"2020-06-28T17:17:01.5638856+01:00","bookedBy":"Andy","estimatedPickupTime":"2020-06-28T17:18:13.1884543+01:00","dispatchedAtTime":"2020-06-28T17:17:03.4654428+01:00","pickedUpAtTime":"2020-06-28T17:55:31.2941184+01:00","vehicleArrivedAtTime":"2020-06-28T17:55:26.5561323+01:00","completedAtTime":"2020-06-28T17:55:33.1802975+01:00","cabExchangeAgentBookingRef":""},` +
				`{"id":641,"bookingId":614,"docketNumber":273,"archiveTime":"2020-06-29T10:53:22.8478573+01:00","archiveReason":"Completed","pickupDueTime":"2020-06-29T10:50:21.5729978+01:00","dropOffDueTime":null,"pickup":{"address":"27 Wittlewood Drive, Accrington, BB5 5DJ","zone":{"id":132,"descriptor":"130","name":"National"},"coordinates":{"longitude":-2.372962777777778,"latitude":53.767095}},"vias":null,"destination":{"address":"Accrington, UK, Lancashire","zone":{"id":132,"descriptor":"130","name":"National"},"coordinates":{"longitude":-2.37218,"latitude":53.753608888888891}},"company":{"id":1,"name":"Driverspay Demo"},"account":null,"loyaltyCardID":null,"loyaltyCardCostValue":0.0,"accountType":"Unknown","paymentType":"Card","name":"c","telephoneNumber":"07715527297","customerEmail":"","capabilities":null,"priority":0,"passengers":0,"luggage":0,"driverNote":"","officeNote":"","yourReferences":{"yourReference1":"","yourReference2":"","yourReference3":"","yourReference4":"","yourReference5":"","yourReference6":"","yourReference7":"","yourReference8":""},"ourReference":"","flightDetails":"","driver":{"id":2105,"callsign":"001","forename":"Andrew","surname":"Baxter","badgeNumber":""},"vehicle":{"id":1027,"callsign":"001","registration":"AD13 NOD","plateNumber":"NOD574R"},"pricing":{"fare":5.08,"cost":5.08,"price":5.08,"waitingTime":0,"waitingTimeChargeable":0,"gratuityAmount":0.0,"costSource":"Shout","pricingTariff":""},"distance":1.4,"systemDistance":1.4,"meterDistance":0.0,"typeOfBooking":"ASAP","bookingSource":"OperatorWeb","bookedAtTime":"2020-06-29T10:50:21.7181196+01:00","bookedBy":"Andy","estimatedPickupTime":"2020-06-29T10:53:32.8053516+01:00","dispatchedAtTime":"2020-06-29T10:52:22.8685851+01:00","pickedUpAtTime":"2020-06-29T10:52:55.4303739+01:00","vehicleArrivedAtTime":"2020-06-29T10:52:49.8717353+01:00","completedAtTime":"2020-06-29T10:53:22.8478573+01:00","cabExchangeAgentBookingRef":""},` +
				`{"id":642,"bookingId":615,"docketNumber":274,"archiveTime":"2020-06-29T10:55:53.8327078+01:00","archiveReason":"Completed","pickupDueTime":"2020-06-29T10:54:43.9359892+01:00","dropOffDueTime":null,"pickup":{"address":"27 Wittlewood Drive, Accrington, BB5 5DJ","zone":{"id":132,"descriptor":"130","name":"National"},"coordinates":{"longitude":-2.372962777777778,"latitude":53.767095}},"vias":null,"destination":{"address":"Accrington, UK, Lancashire","zone":{"id":132,"descriptor":"130","name":"National"},"coordinates":{"longitude":-2.37218,"latitude":53.753608888888891}},"company":{"id":1,"name":"Driverspay Demo"},"account":null,"loyaltyCardID":null,"loyaltyCardCostValue":0.0,"accountType":"Unknown","paymentType":"Card","name":"c","telephoneNumber":"07715527297","customerEmail":"","capabilities":null,"priority":0,"passengers":0,"luggage":0,"driverNote":"","officeNote":"","yourReferences":{"yourReference1":"","yourReference2":"","yourReference3":"","yourReference4":"","yourReference5":"","yourReference6":"","yourReference7":"","yourReference8":""},"ourReference":"","flightDetails":"","driver":{"id":2105,"callsign":"001","forename":"Andrew","surname":"Baxter","badgeNumber":""},"vehicle":{"id":1027,"callsign":"001","registration":"AD13 NOD","plateNumber":"NOD574R"},"pricing":{"fare":5.0,"cost":5.0,"price":5.0,"waitingTime":0,"waitingTimeChargeable":0,"gratuityAmount":0.0,"costSource":"Shout","pricingTariff":""},"distance":1.4,"systemDistance":1.4,"meterDistance":0.0,"typeOfBooking":"ASAP","bookingSource":"OperatorWeb","bookedAtTime":"2020-06-29T10:54:43.9659947+01:00","bookedBy":"Andy","estimatedPickupTime":"2020-06-29T10:55:54.2565874+01:00","dispatchedAtTime":"2020-06-29T10:54:44.7016442+01:00","pickedUpAtTime":"2020-06-29T10:55:33.7987942+01:00","vehicleArrivedAtTime":"2020-06-29T10:54:59.9333557+01:00","completedAtTime":"2020-06-29T10:55:53.8327078+01:00","cabExchangeAgentBookingRef":""}]`)
			if !success {
				data = "[]"
			}
		}
		io.WriteString(w, data)
	}))

	return testServer
}

func TestSendAuthorisationTestServer(t *testing.T) {
	// start a local HTTP server
	testServer := testServer(true, "authorise")
	// close the server when test finishes
	defer testServer.Close()

	params := url.Values{}
	params.Add("username", testAutocabUsername)
	params.Add("password", testAutocabPassword)

	testURL := testServer.URL
	resp := Send(testURL, "POST", nil, params, nil)
	fmt.Println("resp: ", resp)
	if resp == "" {
		t.Fatal("Error authorising user")
	}
}

func TestSendAuthorisationWrongCredentialsTestServer(t *testing.T) {
	// start a local HTTP server
	testServer := testServer(false, "authorise")
	// close the server when test finishes
	defer testServer.Close()

	params := url.Values{}
	params.Add("username", "rubbish")
	params.Add("password", testAutocabPassword)

	testURL := testServer.URL
	resp := Send(testURL, "POST", nil, params, nil)
	fmt.Println("resp: ", resp)
	if resp != "" {
		t.Fatal("Error authorising user should have failed")
	}
}

func TestSendArchiveBookingsTestServer(t *testing.T) {
	// start a local HTTP server
	testServer := testServer(true, "archives")
	// close the server when test finishes
	defer testServer.Close()

	headers := map[string]string{"Authentication-Token": "does_not_matter_for_test_server"}

	// request for archive bookings (do not matter for test server)
	params1 := url.Values{}
	params1.Add("from", "2020/06/25 23:12")
	params1.Add("to", "2020/06/25 23:13")
	params1.Add("ArchiveReasons", "Completed")

	testURL := testServer.URL
	resp := Send(testURL, "POST", headers, params1, nil)
	fmt.Println("resp: ", resp)
	if resp == "" {
		t.Fatal("Error gettting archive bookings")
	}
}

// to live test server
func TestSendAuthorisationLive(t *testing.T) {

	params := url.Values{}
	params.Add("username", testAutocabUsername)
	params.Add("password", testAutocabPassword)

	resp := Send(testAutocabServerURL+"api/thirdparty/v1/authenticate", "POST", nil, params, nil)
	fmt.Println("resp: ", resp)
	if resp == "" {
		t.Fatal("Error authorising user")
	}
}

// to live test server
func TestSendAuthorisationWrongCredentialsLive(t *testing.T) {

	params := url.Values{}
	params.Add("username", "rubbish")
	params.Add("password", testAutocabPassword)

	resp := Send(testAutocabServerURL+"api/thirdparty/v1/authenticate", "POST", nil, params, nil)
	fmt.Println("resp: ", resp)
	if resp != "" {
		t.Fatal("Error authorising user should have failed")
	}
}

// to live test server
func TestSendArchiveBookingsLive(t *testing.T) {

	// Need authrise user for token first
	params := url.Values{}
	params.Add("username", testAutocabUsername)
	params.Add("password", testAutocabPassword)

	resp := Send(testAutocabServerURL+"api/thirdparty/v1/authenticate", "POST", nil, params, nil)
	fmt.Println("resp: ", resp)
	if resp == "" {
		t.Fatal("Error authorising user")
	}
	var jsonResp = []byte(resp)
	var authorisation Authorisation
	err := json.Unmarshal(jsonResp, &authorisation)
	if err != nil {
		t.Fatalf("Error decoding authorisation, error: %+v\n", err)
	}
	fmt.Printf("%+v\n", authorisation)
	headers := map[string]string{"Authentication-Token": authorisation.Secret}

	// current time
	loc, err := time.LoadLocation("Europe/London")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().In(loc)
	to := now.Format("2006/01/02 15:04")
	dayEarlier := now.Add(-time.Hour * 24)
	from := dayEarlier.Format("2006/01/02 15:04")

	// Testing specific times to determine that it is the archive time that to and from works on.
	// Can therefore use the previous to time as the from time for next call.
	// If overlap the conditions on a telephone for send frequency will prevent an SMS being sent.
	// to := "2020/06/25 23:13"
	// from := "2020/06/25 23:12"

	// request for archive bookings
	params1 := url.Values{}
	params1.Add("from", from)
	params1.Add("to", to)
	// params1.Add("ArchiveReasons", "Completed")

	resp1 := Send(testAutocabServerURL+"api/thirdparty/v1/archivedbookings", "POST", headers, params1, nil)
	fmt.Println("resp: ", resp1)
}

// Review Master SMS Gateway
func TestSendReviewMasterSMSGateway(t *testing.T) {
	body, _ := json.Marshal(map[string]string{
		"queue_id":  "81",
		"telephone": "+" + testTelephone,
		"message":   "testing",
	})
	fmt.Println("body:", string(body))

	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
		"Api-Token":    testReviewMasterSMSGatewayApiToken,
	}

	resp := Send(testReviewMasterSMSGatewayURL, "POST", headers, nil, body)
	fmt.Println("resp: ", resp)
}
