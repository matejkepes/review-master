package process

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	testfixtures "gopkg.in/testfixtures.v2"

	"google_reviews_autocab/autocab_api"
	"google_reviews_autocab/config"
	"google_reviews_autocab/database"
)

var fixtures *testfixtures.Context

func TestMain(m *testing.M) {
	var err error

	// read config
	config.ReadProperties()

	// Open connection with the test database.
	// Do NOT import fixtures in a production database!
	// Existing data would be deleted
	database.OpenDB(TestDbName, TestDbAddress, TestDbPort, TestDbUsername, TestDbPassword)

	// creating the context that hold the fixtures
	// see about all compatible databases in this page below
	fixtures, err = testfixtures.NewFolder(database.Db, &testfixtures.MySQL{}, "../database/fixtures")
	if err != nil {
		log.Fatal(err)
	}

	// set review master SMS gateway master queue ID
	database.SetReviewMasterSMSGatewayMasterQueueID()
	if database.ReviewMasterSMSMasterQueue == 0 {
		log.Fatal("Review Master SMS Gateway Master Queue ID not found")
	}

	os.Exit(m.Run())
}

func prepareTestDatabase() {
	// prevent the error message:
	// Loading aborted because the database name does not contains "test"
	testfixtures.SkipDatabaseNameCheck(true)
	if err := fixtures.Load(); err != nil {
		log.Fatal(err)
	}
}

func TestCheckBooking1(t *testing.T) {
	company := autocab_api.Company{ID: 1, Name: "Driverspay Demo"}
	archiveBooking := autocab_api.ArchiveBooking{TelephoneNumber: "07715527297", ArchiveReason: "Completed", BookedAtTime: "2020-06-29T10:54:43.9659947+01:00", PickupDueTime: "2020-06-29T10:54:43.9359892+01:00", PickedUpAtTime: "2020-06-29T10:55:33.7987942+01:00", Company: company, BookingSource: "Operator"}
	// fmt.Printf("archive booking: %+v\n", archiveBooking)
	grcftwc := database.GoogleReviewsConfigFromTokenWithChecks{MinSendFrequency: 21, MaxSendCount: 10, MaxDailySendCount: 20, TelephoneParameter: "t", SendFromIcabbiApp: false, AppKey: "Digital", SecretKey: "Digicomms1!", SendURL: "", HttpGet: false, SendSuccessResponse: "returnSendSms=success", Start: "00:00", End: "23:59", Sunday: true, Monday: true, Tuesday: true, Wednesday: true, Thursday: true, Friday: true, Saturday: true, TimeZone: "Europe/London", ClientID: 11, Country: "GB", MultiMessageEnabled: 0, MessageParameter: "m", MultiMessageSeparator: "SSSSS", UseDatabaseMessage: 1, Message: "Hope you enjoyed your journey", SendDelayEnabled: false, SendDelay: 0, DispatcherChecksEnabled: false, DispatcherURL: "https://ghost-main-static-b36cb86a19e14a2386de12935fac6526.ghostapi.app:29003/", BookingIdParameter: "b", IsBookingForNowDiffMinutes: 10, BookingNowPickupToContactMinutes: 10, PreBookingPickupToContactMinutes: 3, ReplaceTelephoneCountryCode: false, ReplaceTelephoneCountryCodeWith: "0", ReviewMasterSMSGatewayEnabled: false, ReviewMasterSMSGatewayUseMasterQueue: false, ReviewMasterSMSGatewayPairCode: "1234", Companies: "", BookingSourceMobileAppState: -1}
	// fmt.Printf("grcftwc: %+v\n", grcftwc)
	// fmt.Println(CheckBooking(archiveBooking, grcftwc))
	ok, tel, ts, m, c := CheckBooking(archiveBooking, grcftwc)
	fmt.Printf("%t, %s, %s, %s, %d\n", ok, tel, ts, m, c)
	if !ok {
		t.Fatalf("Error cheking booking should have been ok\n")
	}
}

// replace telephone country code telephone send sms should be different
func TestCheckBooking2(t *testing.T) {
	company := autocab_api.Company{ID: 1, Name: "Driverspay Demo"}
	archiveBooking := autocab_api.ArchiveBooking{TelephoneNumber: "07715527297", ArchiveReason: "Completed", BookedAtTime: "2020-06-29T10:54:43.9659947+01:00", PickupDueTime: "2020-06-29T10:54:43.9359892+01:00", PickedUpAtTime: "2020-06-29T10:55:33.7987942+01:00", Company: company, BookingSource: "Operator"}
	// fmt.Printf("archive booking: %+v\n", archiveBooking)
	grcftwc := database.GoogleReviewsConfigFromTokenWithChecks{MinSendFrequency: 21, MaxSendCount: 10, MaxDailySendCount: 20, TelephoneParameter: "t", SendFromIcabbiApp: false, AppKey: "Digital", SecretKey: "Digicomms1!", SendURL: "", HttpGet: false, SendSuccessResponse: "returnSendSms=success", Start: "00:00", End: "23:59", Sunday: true, Monday: true, Tuesday: true, Wednesday: true, Thursday: true, Friday: true, Saturday: true, TimeZone: "Europe/London", ClientID: 11, Country: "GB", MultiMessageEnabled: 0, MessageParameter: "m", MultiMessageSeparator: "SSSSS", UseDatabaseMessage: 1, Message: "Hope you enjoyed your journey", SendDelayEnabled: false, SendDelay: 0, DispatcherChecksEnabled: false, DispatcherURL: "https://ghost-main-static-b36cb86a19e14a2386de12935fac6526.ghostapi.app:29003/", BookingIdParameter: "b", IsBookingForNowDiffMinutes: 10, BookingNowPickupToContactMinutes: 10, PreBookingPickupToContactMinutes: 3, ReplaceTelephoneCountryCode: true, ReplaceTelephoneCountryCodeWith: "0", ReviewMasterSMSGatewayEnabled: false, ReviewMasterSMSGatewayUseMasterQueue: false, ReviewMasterSMSGatewayPairCode: "1234", Companies: "", BookingSourceMobileAppState: -1}
	// fmt.Printf("grcftwc: %+v\n", grcftwc)
	// fmt.Println(CheckBooking(archiveBooking, grcftwc))
	ok, tel, ts, m, c := CheckBooking(archiveBooking, grcftwc)
	fmt.Printf("%t, %s, %s, %s, %d\n", ok, tel, ts, m, c)
	if !ok {
		t.Fatalf("Error cheking booking should have been ok\n")
	}
}

// landline, should fail check
func TestCheckBooking3(t *testing.T) {
	company := autocab_api.Company{ID: 1, Name: "Driverspay Demo"}
	archiveBooking := autocab_api.ArchiveBooking{TelephoneNumber: "01132345678", ArchiveReason: "Completed", BookedAtTime: "2020-06-29T10:54:43.9659947+01:00", PickupDueTime: "2020-06-29T10:54:43.9359892+01:00", PickedUpAtTime: "2020-06-29T10:55:33.7987942+01:00", Company: company, BookingSource: "Operator"}
	// fmt.Printf("archive booking: %+v\n", archiveBooking)
	grcftwc := database.GoogleReviewsConfigFromTokenWithChecks{MinSendFrequency: 21, MaxSendCount: 10, MaxDailySendCount: 20, TelephoneParameter: "t", SendFromIcabbiApp: false, AppKey: "Digital", SecretKey: "Digicomms1!", SendURL: "", HttpGet: false, SendSuccessResponse: "returnSendSms=success", Start: "00:00", End: "23:59", Sunday: true, Monday: true, Tuesday: true, Wednesday: true, Thursday: true, Friday: true, Saturday: true, TimeZone: "Europe/London", ClientID: 11, Country: "GB", MultiMessageEnabled: 0, MessageParameter: "m", MultiMessageSeparator: "SSSSS", UseDatabaseMessage: 1, Message: "Hope you enjoyed your journey", SendDelayEnabled: false, SendDelay: 0, DispatcherChecksEnabled: false, DispatcherURL: "https://ghost-main-static-b36cb86a19e14a2386de12935fac6526.ghostapi.app:29003/", BookingIdParameter: "b", IsBookingForNowDiffMinutes: 10, BookingNowPickupToContactMinutes: 10, PreBookingPickupToContactMinutes: 3, ReplaceTelephoneCountryCode: false, ReplaceTelephoneCountryCodeWith: "0", ReviewMasterSMSGatewayEnabled: false, ReviewMasterSMSGatewayUseMasterQueue: false, ReviewMasterSMSGatewayPairCode: "1234", Companies: "", BookingSourceMobileAppState: -1}
	// fmt.Printf("grcftwc: %+v\n", grcftwc)
	// fmt.Println(CheckBooking(archiveBooking, grcftwc))
	ok, tel, ts, m, c := CheckBooking(archiveBooking, grcftwc)
	fmt.Printf("%t, %s, %s, %s, %d\n", ok, tel, ts, m, c)
	if ok {
		t.Fatalf("Error cheking booking should have NOT been ok\n")
	}
}

// booking source mobile app state set to 1 - must be mobile app
func TestCheckBooking4(t *testing.T) {
	company := autocab_api.Company{ID: 1, Name: "Driverspay Demo"}
	archiveBooking := autocab_api.ArchiveBooking{TelephoneNumber: "07715527297", ArchiveReason: "Completed", BookedAtTime: "2020-06-29T10:54:43.9659947+01:00", PickupDueTime: "2020-06-29T10:54:43.9359892+01:00", PickedUpAtTime: "2020-06-29T10:55:33.7987942+01:00", Company: company, BookingSource: "MobileApp"}
	// fmt.Printf("archive booking: %+v\n", archiveBooking)
	grcftwc := database.GoogleReviewsConfigFromTokenWithChecks{MinSendFrequency: 21, MaxSendCount: 10, MaxDailySendCount: 20, TelephoneParameter: "t", SendFromIcabbiApp: false, AppKey: "Digital", SecretKey: "Digicomms1!", SendURL: "", HttpGet: false, SendSuccessResponse: "returnSendSms=success", Start: "00:00", End: "23:59", Sunday: true, Monday: true, Tuesday: true, Wednesday: true, Thursday: true, Friday: true, Saturday: true, TimeZone: "Europe/London", ClientID: 11, Country: "GB", MultiMessageEnabled: 0, MessageParameter: "m", MultiMessageSeparator: "SSSSS", UseDatabaseMessage: 1, Message: "Hope you enjoyed your journey", SendDelayEnabled: false, SendDelay: 0, DispatcherChecksEnabled: false, DispatcherURL: "https://ghost-main-static-b36cb86a19e14a2386de12935fac6526.ghostapi.app:29003/", BookingIdParameter: "b", IsBookingForNowDiffMinutes: 10, BookingNowPickupToContactMinutes: 10, PreBookingPickupToContactMinutes: 3, ReplaceTelephoneCountryCode: false, ReplaceTelephoneCountryCodeWith: "0", ReviewMasterSMSGatewayEnabled: false, ReviewMasterSMSGatewayUseMasterQueue: false, ReviewMasterSMSGatewayPairCode: "1234", Companies: "", BookingSourceMobileAppState: 1}
	// fmt.Printf("grcftwc: %+v\n", grcftwc)
	// fmt.Println(CheckBooking(archiveBooking, grcftwc))
	ok, tel, ts, m, c := CheckBooking(archiveBooking, grcftwc)
	fmt.Printf("%t, %s, %s, %s, %d\n", ok, tel, ts, m, c)
	if !ok {
		t.Fatalf("Error cheking booking should have been ok\n")
	}
}

// booking source mobile app state set to 0 - must be mobile NOT app
func TestCheckBooking5(t *testing.T) {
	company := autocab_api.Company{ID: 1, Name: "Driverspay Demo"}
	archiveBooking := autocab_api.ArchiveBooking{TelephoneNumber: "07715527297", ArchiveReason: "Completed", BookedAtTime: "2020-06-29T10:54:43.9659947+01:00", PickupDueTime: "2020-06-29T10:54:43.9359892+01:00", PickedUpAtTime: "2020-06-29T10:55:33.7987942+01:00", Company: company, BookingSource: "Operator"}
	// fmt.Printf("archive booking: %+v\n", archiveBooking)
	grcftwc := database.GoogleReviewsConfigFromTokenWithChecks{MinSendFrequency: 21, MaxSendCount: 10, MaxDailySendCount: 20, TelephoneParameter: "t", SendFromIcabbiApp: false, AppKey: "Digital", SecretKey: "Digicomms1!", SendURL: "", HttpGet: false, SendSuccessResponse: "returnSendSms=success", Start: "00:00", End: "23:59", Sunday: true, Monday: true, Tuesday: true, Wednesday: true, Thursday: true, Friday: true, Saturday: true, TimeZone: "Europe/London", ClientID: 11, Country: "GB", MultiMessageEnabled: 0, MessageParameter: "m", MultiMessageSeparator: "SSSSS", UseDatabaseMessage: 1, Message: "Hope you enjoyed your journey", SendDelayEnabled: false, SendDelay: 0, DispatcherChecksEnabled: false, DispatcherURL: "https://ghost-main-static-b36cb86a19e14a2386de12935fac6526.ghostapi.app:29003/", BookingIdParameter: "b", IsBookingForNowDiffMinutes: 10, BookingNowPickupToContactMinutes: 10, PreBookingPickupToContactMinutes: 3, ReplaceTelephoneCountryCode: false, ReplaceTelephoneCountryCodeWith: "0", ReviewMasterSMSGatewayEnabled: false, ReviewMasterSMSGatewayUseMasterQueue: false, ReviewMasterSMSGatewayPairCode: "1234", Companies: "", BookingSourceMobileAppState: 0}
	// fmt.Printf("grcftwc: %+v\n", grcftwc)
	// fmt.Println(CheckBooking(archiveBooking, grcftwc))
	ok, tel, ts, m, c := CheckBooking(archiveBooking, grcftwc)
	fmt.Printf("%t, %s, %s, %s, %d\n", ok, tel, ts, m, c)
	if !ok {
		t.Fatalf("Error cheking booking should have been ok\n")
	}
}

// booking source mobile app state set to 1 - must be mobile app but the booking source is operator
func TestCheckBooking6(t *testing.T) {
	company := autocab_api.Company{ID: 1, Name: "Driverspay Demo"}
	archiveBooking := autocab_api.ArchiveBooking{TelephoneNumber: "07715527297", ArchiveReason: "Completed", BookedAtTime: "2020-06-29T10:54:43.9659947+01:00", PickupDueTime: "2020-06-29T10:54:43.9359892+01:00", PickedUpAtTime: "2020-06-29T10:55:33.7987942+01:00", Company: company, BookingSource: "Operator"}
	// fmt.Printf("archive booking: %+v\n", archiveBooking)
	grcftwc := database.GoogleReviewsConfigFromTokenWithChecks{MinSendFrequency: 21, MaxSendCount: 10, MaxDailySendCount: 20, TelephoneParameter: "t", SendFromIcabbiApp: false, AppKey: "Digital", SecretKey: "Digicomms1!", SendURL: "", HttpGet: false, SendSuccessResponse: "returnSendSms=success", Start: "00:00", End: "23:59", Sunday: true, Monday: true, Tuesday: true, Wednesday: true, Thursday: true, Friday: true, Saturday: true, TimeZone: "Europe/London", ClientID: 11, Country: "GB", MultiMessageEnabled: 0, MessageParameter: "m", MultiMessageSeparator: "SSSSS", UseDatabaseMessage: 1, Message: "Hope you enjoyed your journey", SendDelayEnabled: false, SendDelay: 0, DispatcherChecksEnabled: false, DispatcherURL: "https://ghost-main-static-b36cb86a19e14a2386de12935fac6526.ghostapi.app:29003/", BookingIdParameter: "b", IsBookingForNowDiffMinutes: 10, BookingNowPickupToContactMinutes: 10, PreBookingPickupToContactMinutes: 3, ReplaceTelephoneCountryCode: false, ReplaceTelephoneCountryCodeWith: "0", ReviewMasterSMSGatewayEnabled: false, ReviewMasterSMSGatewayUseMasterQueue: false, ReviewMasterSMSGatewayPairCode: "1234", Companies: "", BookingSourceMobileAppState: 1}
	// fmt.Printf("grcftwc: %+v\n", grcftwc)
	// fmt.Println(CheckBooking(archiveBooking, grcftwc))
	ok, tel, ts, m, c := CheckBooking(archiveBooking, grcftwc)
	fmt.Printf("%t, %s, %s, %s, %d\n", ok, tel, ts, m, c)
	if ok {
		t.Fatalf("Error cheking booking should have NOT been ok\n")
	}
}

// booking source mobile app state set to 0 - must be mobile NOT app but the booking source is mobile app
func TestCheckBooking7(t *testing.T) {
	company := autocab_api.Company{ID: 1, Name: "Driverspay Demo"}
	archiveBooking := autocab_api.ArchiveBooking{TelephoneNumber: "07715527297", ArchiveReason: "Completed", BookedAtTime: "2020-06-29T10:54:43.9659947+01:00", PickupDueTime: "2020-06-29T10:54:43.9359892+01:00", PickedUpAtTime: "2020-06-29T10:55:33.7987942+01:00", Company: company, BookingSource: "MobileApp"}
	// fmt.Printf("archive booking: %+v\n", archiveBooking)
	grcftwc := database.GoogleReviewsConfigFromTokenWithChecks{MinSendFrequency: 21, MaxSendCount: 10, MaxDailySendCount: 20, TelephoneParameter: "t", SendFromIcabbiApp: false, AppKey: "Digital", SecretKey: "Digicomms1!", SendURL: "", HttpGet: false, SendSuccessResponse: "returnSendSms=success", Start: "00:00", End: "23:59", Sunday: true, Monday: true, Tuesday: true, Wednesday: true, Thursday: true, Friday: true, Saturday: true, TimeZone: "Europe/London", ClientID: 11, Country: "GB", MultiMessageEnabled: 0, MessageParameter: "m", MultiMessageSeparator: "SSSSS", UseDatabaseMessage: 1, Message: "Hope you enjoyed your journey", SendDelayEnabled: false, SendDelay: 0, DispatcherChecksEnabled: false, DispatcherURL: "https://ghost-main-static-b36cb86a19e14a2386de12935fac6526.ghostapi.app:29003/", BookingIdParameter: "b", IsBookingForNowDiffMinutes: 10, BookingNowPickupToContactMinutes: 10, PreBookingPickupToContactMinutes: 3, ReplaceTelephoneCountryCode: false, ReplaceTelephoneCountryCodeWith: "0", ReviewMasterSMSGatewayEnabled: false, ReviewMasterSMSGatewayUseMasterQueue: false, ReviewMasterSMSGatewayPairCode: "1234", Companies: "", BookingSourceMobileAppState: 0}
	// fmt.Printf("grcftwc: %+v\n", grcftwc)
	// fmt.Println(CheckBooking(archiveBooking, grcftwc))
	ok, tel, ts, m, c := CheckBooking(archiveBooking, grcftwc)
	fmt.Printf("%t, %s, %s, %s, %d\n", ok, tel, ts, m, c)
	if ok {
		t.Fatalf("Error cheking booking should have NOT been ok\n")
	}
}

// companies config is not an empty string and is set to a list of acceptable company id's
func TestCheckBooking8(t *testing.T) {
	company := autocab_api.Company{ID: 1, Name: "Driverspay Demo"}
	archiveBooking := autocab_api.ArchiveBooking{TelephoneNumber: "07715527297", ArchiveReason: "Completed", BookedAtTime: "2020-06-29T10:54:43.9659947+01:00", PickupDueTime: "2020-06-29T10:54:43.9359892+01:00", PickedUpAtTime: "2020-06-29T10:55:33.7987942+01:00", Company: company, BookingSource: "Operator"}
	// fmt.Printf("archive booking: %+v\n", archiveBooking)
	grcftwc := database.GoogleReviewsConfigFromTokenWithChecks{MinSendFrequency: 21, MaxSendCount: 10, MaxDailySendCount: 20, TelephoneParameter: "t", SendFromIcabbiApp: false, AppKey: "Digital", SecretKey: "Digicomms1!", SendURL: "", HttpGet: false, SendSuccessResponse: "returnSendSms=success", Start: "00:00", End: "23:59", Sunday: true, Monday: true, Tuesday: true, Wednesday: true, Thursday: true, Friday: true, Saturday: true, TimeZone: "Europe/London", ClientID: 11, Country: "GB", MultiMessageEnabled: 0, MessageParameter: "m", MultiMessageSeparator: "SSSSS", UseDatabaseMessage: 1, Message: "Hope you enjoyed your journey", SendDelayEnabled: false, SendDelay: 0, DispatcherChecksEnabled: false, DispatcherURL: "https://ghost-main-static-b36cb86a19e14a2386de12935fac6526.ghostapi.app:29003/", BookingIdParameter: "b", IsBookingForNowDiffMinutes: 10, BookingNowPickupToContactMinutes: 10, PreBookingPickupToContactMinutes: 3, ReplaceTelephoneCountryCode: false, ReplaceTelephoneCountryCodeWith: "0", ReviewMasterSMSGatewayEnabled: false, ReviewMasterSMSGatewayUseMasterQueue: false, ReviewMasterSMSGatewayPairCode: "1234", Companies: "1", BookingSourceMobileAppState: -1}
	// fmt.Printf("grcftwc: %+v\n", grcftwc)
	// fmt.Println(CheckBooking(archiveBooking, grcftwc))
	ok, tel, ts, m, c := CheckBooking(archiveBooking, grcftwc)
	fmt.Printf("%t, %s, %s, %s, %d\n", ok, tel, ts, m, c)
	if !ok {
		t.Fatalf("Error cheking booking should have been ok\n")
	}
}

// companies config is not an empty string and is set to a list of acceptable company id's fails because company NOT in list
func TestCheckBooking9(t *testing.T) {
	company := autocab_api.Company{ID: 2, Name: "Driverspay Demo"}
	archiveBooking := autocab_api.ArchiveBooking{TelephoneNumber: "07715527297", ArchiveReason: "Completed", BookedAtTime: "2020-06-29T10:54:43.9659947+01:00", PickupDueTime: "2020-06-29T10:54:43.9359892+01:00", PickedUpAtTime: "2020-06-29T10:55:33.7987942+01:00", Company: company, BookingSource: "Operator"}
	// fmt.Printf("archive booking: %+v\n", archiveBooking)
	grcftwc := database.GoogleReviewsConfigFromTokenWithChecks{MinSendFrequency: 21, MaxSendCount: 10, MaxDailySendCount: 20, TelephoneParameter: "t", SendFromIcabbiApp: false, AppKey: "Digital", SecretKey: "Digicomms1!", SendURL: "", HttpGet: false, SendSuccessResponse: "returnSendSms=success", Start: "00:00", End: "23:59", Sunday: true, Monday: true, Tuesday: true, Wednesday: true, Thursday: true, Friday: true, Saturday: true, TimeZone: "Europe/London", ClientID: 11, Country: "GB", MultiMessageEnabled: 0, MessageParameter: "m", MultiMessageSeparator: "SSSSS", UseDatabaseMessage: 1, Message: "Hope you enjoyed your journey", SendDelayEnabled: false, SendDelay: 0, DispatcherChecksEnabled: false, DispatcherURL: "https://ghost-main-static-b36cb86a19e14a2386de12935fac6526.ghostapi.app:29003/", BookingIdParameter: "b", IsBookingForNowDiffMinutes: 10, BookingNowPickupToContactMinutes: 10, PreBookingPickupToContactMinutes: 3, ReplaceTelephoneCountryCode: false, ReplaceTelephoneCountryCodeWith: "0", ReviewMasterSMSGatewayEnabled: false, ReviewMasterSMSGatewayUseMasterQueue: false, ReviewMasterSMSGatewayPairCode: "1234", Companies: "1", BookingSourceMobileAppState: -1}
	// fmt.Printf("grcftwc: %+v\n", grcftwc)
	// fmt.Println(CheckBooking(archiveBooking, grcftwc))
	ok, tel, ts, m, c := CheckBooking(archiveBooking, grcftwc)
	fmt.Printf("%t, %s, %s, %s, %d\n", ok, tel, ts, m, c)
	if ok {
		t.Fatalf("Error cheking booking should have NOT been ok\n")
	}
}

// companies config is not an empty string and is set to a list of acceptable company id's more than one in list
func TestCheckBooking10(t *testing.T) {
	company := autocab_api.Company{ID: 2, Name: "Driverspay Demo"}
	archiveBooking := autocab_api.ArchiveBooking{TelephoneNumber: "07715527297", ArchiveReason: "Completed", BookedAtTime: "2020-06-29T10:54:43.9659947+01:00", PickupDueTime: "2020-06-29T10:54:43.9359892+01:00", PickedUpAtTime: "2020-06-29T10:55:33.7987942+01:00", Company: company, BookingSource: "Operator"}
	// fmt.Printf("archive booking: %+v\n", archiveBooking)
	grcftwc := database.GoogleReviewsConfigFromTokenWithChecks{MinSendFrequency: 21, MaxSendCount: 10, MaxDailySendCount: 20, TelephoneParameter: "t", SendFromIcabbiApp: false, AppKey: "Digital", SecretKey: "Digicomms1!", SendURL: "", HttpGet: false, SendSuccessResponse: "returnSendSms=success", Start: "00:00", End: "23:59", Sunday: true, Monday: true, Tuesday: true, Wednesday: true, Thursday: true, Friday: true, Saturday: true, TimeZone: "Europe/London", ClientID: 11, Country: "GB", MultiMessageEnabled: 0, MessageParameter: "m", MultiMessageSeparator: "SSSSS", UseDatabaseMessage: 1, Message: "Hope you enjoyed your journey", SendDelayEnabled: false, SendDelay: 0, DispatcherChecksEnabled: false, DispatcherURL: "https://ghost-main-static-b36cb86a19e14a2386de12935fac6526.ghostapi.app:29003/", BookingIdParameter: "b", IsBookingForNowDiffMinutes: 10, BookingNowPickupToContactMinutes: 10, PreBookingPickupToContactMinutes: 3, ReplaceTelephoneCountryCode: false, ReplaceTelephoneCountryCodeWith: "0", ReviewMasterSMSGatewayEnabled: false, ReviewMasterSMSGatewayUseMasterQueue: false, ReviewMasterSMSGatewayPairCode: "1234", Companies: "1,2,3,4", BookingSourceMobileAppState: -1}
	// fmt.Printf("grcftwc: %+v\n", grcftwc)
	// fmt.Println(CheckBooking(archiveBooking, grcftwc))
	ok, tel, ts, m, c := CheckBooking(archiveBooking, grcftwc)
	fmt.Printf("%t, %s, %s, %s, %d\n", ok, tel, ts, m, c)
	if !ok {
		t.Fatalf("Error cheking booking should have been ok\n")
	}
}

// companies config is not an empty string and is set to a list of acceptable company id's more than one in list with spaces
func TestCheckBooking11(t *testing.T) {
	company := autocab_api.Company{ID: 2, Name: "Driverspay Demo"}
	archiveBooking := autocab_api.ArchiveBooking{TelephoneNumber: "07715527297", ArchiveReason: "Completed", BookedAtTime: "2020-06-29T10:54:43.9659947+01:00", PickupDueTime: "2020-06-29T10:54:43.9359892+01:00", PickedUpAtTime: "2020-06-29T10:55:33.7987942+01:00", Company: company, BookingSource: "Operator"}
	// fmt.Printf("archive booking: %+v\n", archiveBooking)
	grcftwc := database.GoogleReviewsConfigFromTokenWithChecks{MinSendFrequency: 21, MaxSendCount: 10, MaxDailySendCount: 20, TelephoneParameter: "t", SendFromIcabbiApp: false, AppKey: "Digital", SecretKey: "Digicomms1!", SendURL: "", HttpGet: false, SendSuccessResponse: "returnSendSms=success", Start: "00:00", End: "23:59", Sunday: true, Monday: true, Tuesday: true, Wednesday: true, Thursday: true, Friday: true, Saturday: true, TimeZone: "Europe/London", ClientID: 11, Country: "GB", MultiMessageEnabled: 0, MessageParameter: "m", MultiMessageSeparator: "SSSSS", UseDatabaseMessage: 1, Message: "Hope you enjoyed your journey", SendDelayEnabled: false, SendDelay: 0, DispatcherChecksEnabled: false, DispatcherURL: "https://ghost-main-static-b36cb86a19e14a2386de12935fac6526.ghostapi.app:29003/", BookingIdParameter: "b", IsBookingForNowDiffMinutes: 10, BookingNowPickupToContactMinutes: 10, PreBookingPickupToContactMinutes: 3, ReplaceTelephoneCountryCode: false, ReplaceTelephoneCountryCodeWith: "0", ReviewMasterSMSGatewayEnabled: false, ReviewMasterSMSGatewayUseMasterQueue: false, ReviewMasterSMSGatewayPairCode: "1234", Companies: "1, 2 , 3 , 4 ", BookingSourceMobileAppState: -1}
	// fmt.Printf("grcftwc: %+v\n", grcftwc)
	// fmt.Println(CheckBooking(archiveBooking, grcftwc))
	ok, tel, ts, m, c := CheckBooking(archiveBooking, grcftwc)
	fmt.Printf("%t, %s, %s, %s, %d\n", ok, tel, ts, m, c)
	if !ok {
		t.Fatalf("Error cheking booking should have been ok\n")
	}
}

// companies config is not an empty string but spaces
func TestCheckBooking12(t *testing.T) {
	company := autocab_api.Company{ID: 1, Name: "Driverspay Demo"}
	archiveBooking := autocab_api.ArchiveBooking{TelephoneNumber: "07715527297", ArchiveReason: "Completed", BookedAtTime: "2020-06-29T10:54:43.9659947+01:00", PickupDueTime: "2020-06-29T10:54:43.9359892+01:00", PickedUpAtTime: "2020-06-29T10:55:33.7987942+01:00", Company: company, BookingSource: "Operator"}
	// fmt.Printf("archive booking: %+v\n", archiveBooking)
	grcftwc := database.GoogleReviewsConfigFromTokenWithChecks{MinSendFrequency: 21, MaxSendCount: 10, MaxDailySendCount: 20, TelephoneParameter: "t", SendFromIcabbiApp: false, AppKey: "Digital", SecretKey: "Digicomms1!", SendURL: "", HttpGet: false, SendSuccessResponse: "returnSendSms=success", Start: "00:00", End: "23:59", Sunday: true, Monday: true, Tuesday: true, Wednesday: true, Thursday: true, Friday: true, Saturday: true, TimeZone: "Europe/London", ClientID: 11, Country: "GB", MultiMessageEnabled: 0, MessageParameter: "m", MultiMessageSeparator: "SSSSS", UseDatabaseMessage: 1, Message: "Hope you enjoyed your journey", SendDelayEnabled: false, SendDelay: 0, DispatcherChecksEnabled: false, DispatcherURL: "https://ghost-main-static-b36cb86a19e14a2386de12935fac6526.ghostapi.app:29003/", BookingIdParameter: "b", IsBookingForNowDiffMinutes: 10, BookingNowPickupToContactMinutes: 10, PreBookingPickupToContactMinutes: 3, ReplaceTelephoneCountryCode: false, ReplaceTelephoneCountryCodeWith: "0", ReviewMasterSMSGatewayEnabled: false, ReviewMasterSMSGatewayUseMasterQueue: false, ReviewMasterSMSGatewayPairCode: "1234", Companies: " ", BookingSourceMobileAppState: -1}
	// fmt.Printf("grcftwc: %+v\n", grcftwc)
	// fmt.Println(CheckBooking(archiveBooking, grcftwc))
	ok, tel, ts, m, c := CheckBooking(archiveBooking, grcftwc)
	fmt.Printf("%t, %s, %s, %s, %d\n", ok, tel, ts, m, c)
	if !ok {
		t.Fatalf("Error cheking booking should have been ok\n")
	}
}

// companies config is not an empty string and is set to a list of acceptable company id's but has commas with spaces and some cannot be converted to an integer
func TestCheckBooking13(t *testing.T) {
	company := autocab_api.Company{ID: 1, Name: "Driverspay Demo"}
	archiveBooking := autocab_api.ArchiveBooking{TelephoneNumber: "07715527297", ArchiveReason: "Completed", BookedAtTime: "2020-06-29T10:54:43.9659947+01:00", PickupDueTime: "2020-06-29T10:54:43.9359892+01:00", PickedUpAtTime: "2020-06-29T10:55:33.7987942+01:00", Company: company, BookingSource: "Operator"}
	// fmt.Printf("archive booking: %+v\n", archiveBooking)
	grcftwc := database.GoogleReviewsConfigFromTokenWithChecks{MinSendFrequency: 21, MaxSendCount: 10, MaxDailySendCount: 20, TelephoneParameter: "t", SendFromIcabbiApp: false, AppKey: "Digital", SecretKey: "Digicomms1!", SendURL: "", HttpGet: false, SendSuccessResponse: "returnSendSms=success", Start: "00:00", End: "23:59", Sunday: true, Monday: true, Tuesday: true, Wednesday: true, Thursday: true, Friday: true, Saturday: true, TimeZone: "Europe/London", ClientID: 11, Country: "GB", MultiMessageEnabled: 0, MessageParameter: "m", MultiMessageSeparator: "SSSSS", UseDatabaseMessage: 1, Message: "Hope you enjoyed your journey", SendDelayEnabled: false, SendDelay: 0, DispatcherChecksEnabled: false, DispatcherURL: "https://ghost-main-static-b36cb86a19e14a2386de12935fac6526.ghostapi.app:29003/", BookingIdParameter: "b", IsBookingForNowDiffMinutes: 10, BookingNowPickupToContactMinutes: 10, PreBookingPickupToContactMinutes: 3, ReplaceTelephoneCountryCode: false, ReplaceTelephoneCountryCodeWith: "0", ReviewMasterSMSGatewayEnabled: false, ReviewMasterSMSGatewayUseMasterQueue: false, ReviewMasterSMSGatewayPairCode: "1234", Companies: " , 1, ", BookingSourceMobileAppState: -1}
	// fmt.Printf("grcftwc: %+v\n", grcftwc)
	// fmt.Println(CheckBooking(archiveBooking, grcftwc))
	ok, tel, ts, m, c := CheckBooking(archiveBooking, grcftwc)
	fmt.Printf("%t, %s, %s, %s, %d\n", ok, tel, ts, m, c)
	if !ok {
		t.Fatalf("Error cheking booking should have been ok\n")
	}
}

// companies config is not an empty string but spaces with commas
func TestCheckBooking14(t *testing.T) {
	company := autocab_api.Company{ID: 1, Name: "Driverspay Demo"}
	archiveBooking := autocab_api.ArchiveBooking{TelephoneNumber: "07715527297", ArchiveReason: "Completed", BookedAtTime: "2020-06-29T10:54:43.9659947+01:00", PickupDueTime: "2020-06-29T10:54:43.9359892+01:00", PickedUpAtTime: "2020-06-29T10:55:33.7987942+01:00", Company: company, BookingSource: "Operator"}
	// fmt.Printf("archive booking: %+v\n", archiveBooking)
	grcftwc := database.GoogleReviewsConfigFromTokenWithChecks{MinSendFrequency: 21, MaxSendCount: 10, MaxDailySendCount: 20, TelephoneParameter: "t", SendFromIcabbiApp: false, AppKey: "Digital", SecretKey: "Digicomms1!", SendURL: "", HttpGet: false, SendSuccessResponse: "returnSendSms=success", Start: "00:00", End: "23:59", Sunday: true, Monday: true, Tuesday: true, Wednesday: true, Thursday: true, Friday: true, Saturday: true, TimeZone: "Europe/London", ClientID: 11, Country: "GB", MultiMessageEnabled: 0, MessageParameter: "m", MultiMessageSeparator: "SSSSS", UseDatabaseMessage: 1, Message: "Hope you enjoyed your journey", SendDelayEnabled: false, SendDelay: 0, DispatcherChecksEnabled: false, DispatcherURL: "https://ghost-main-static-b36cb86a19e14a2386de12935fac6526.ghostapi.app:29003/", BookingIdParameter: "b", IsBookingForNowDiffMinutes: 10, BookingNowPickupToContactMinutes: 10, PreBookingPickupToContactMinutes: 3, ReplaceTelephoneCountryCode: false, ReplaceTelephoneCountryCodeWith: "0", ReviewMasterSMSGatewayEnabled: false, ReviewMasterSMSGatewayUseMasterQueue: false, ReviewMasterSMSGatewayPairCode: "1234", Companies: " , , ", BookingSourceMobileAppState: -1}
	// fmt.Printf("grcftwc: %+v\n", grcftwc)
	// fmt.Println(CheckBooking(archiveBooking, grcftwc))
	ok, tel, ts, m, c := CheckBooking(archiveBooking, grcftwc)
	fmt.Printf("%t, %s, %s, %s, %d\n", ok, tel, ts, m, c)
	if !ok {
		t.Fatalf("Error cheking booking should have been ok\n")
	}
}

// companies config is not an empty string but spaces with commas and rubbish
func TestCheckBooking15(t *testing.T) {
	company := autocab_api.Company{ID: 1, Name: "Driverspay Demo"}
	archiveBooking := autocab_api.ArchiveBooking{TelephoneNumber: "07715527297", ArchiveReason: "Completed", BookedAtTime: "2020-06-29T10:54:43.9659947+01:00", PickupDueTime: "2020-06-29T10:54:43.9359892+01:00", PickedUpAtTime: "2020-06-29T10:55:33.7987942+01:00", Company: company, BookingSource: "Operator"}
	// fmt.Printf("archive booking: %+v\n", archiveBooking)
	grcftwc := database.GoogleReviewsConfigFromTokenWithChecks{MinSendFrequency: 21, MaxSendCount: 10, MaxDailySendCount: 20, TelephoneParameter: "t", SendFromIcabbiApp: false, AppKey: "Digital", SecretKey: "Digicomms1!", SendURL: "", HttpGet: false, SendSuccessResponse: "returnSendSms=success", Start: "00:00", End: "23:59", Sunday: true, Monday: true, Tuesday: true, Wednesday: true, Thursday: true, Friday: true, Saturday: true, TimeZone: "Europe/London", ClientID: 11, Country: "GB", MultiMessageEnabled: 0, MessageParameter: "m", MultiMessageSeparator: "SSSSS", UseDatabaseMessage: 1, Message: "Hope you enjoyed your journey", SendDelayEnabled: false, SendDelay: 0, DispatcherChecksEnabled: false, DispatcherURL: "https://ghost-main-static-b36cb86a19e14a2386de12935fac6526.ghostapi.app:29003/", BookingIdParameter: "b", IsBookingForNowDiffMinutes: 10, BookingNowPickupToContactMinutes: 10, PreBookingPickupToContactMinutes: 3, ReplaceTelephoneCountryCode: false, ReplaceTelephoneCountryCodeWith: "0", ReviewMasterSMSGatewayEnabled: false, ReviewMasterSMSGatewayUseMasterQueue: false, ReviewMasterSMSGatewayPairCode: "1234", Companies: " , A, B, @,,, ", BookingSourceMobileAppState: -1}
	// fmt.Printf("grcftwc: %+v\n", grcftwc)
	// fmt.Println(CheckBooking(archiveBooking, grcftwc))
	ok, tel, ts, m, c := CheckBooking(archiveBooking, grcftwc)
	fmt.Printf("%t, %s, %s, %s, %d\n", ok, tel, ts, m, c)
	if !ok {
		t.Fatalf("Error cheking booking should have been ok\n")
	}
}

func TestPollAutocab(t *testing.T) {
	prepareTestDatabase()
	PollAutocab(config.Conf.LastPollTime, time.Now())
}

func TestSendReviewMasterSMSGateway1(t *testing.T) {
	grcftwc := database.GoogleReviewsConfigFromTokenWithChecks{MinSendFrequency: 21, MaxSendCount: 10, MaxDailySendCount: 20, TelephoneParameter: "t", SendFromIcabbiApp: false, AppKey: "Digital", SecretKey: "Digicomms1!", SendURL: "", HttpGet: false, SendSuccessResponse: "returnSendSms=success", Start: "00:00", End: "23:59", Sunday: true, Monday: true, Tuesday: true, Wednesday: true, Thursday: true, Friday: true, Saturday: true, TimeZone: "Europe/London", ClientID: 11, Country: "GB", MultiMessageEnabled: 0, MessageParameter: "m", MultiMessageSeparator: "SSSSS", UseDatabaseMessage: 1, Message: "Hope you enjoyed your journey", SendDelayEnabled: false, SendDelay: 0, DispatcherChecksEnabled: false, DispatcherURL: "https://ghost-main-static-b36cb86a19e14a2386de12935fac6526.ghostapi.app:29003/", BookingIdParameter: "b", IsBookingForNowDiffMinutes: 10, BookingNowPickupToContactMinutes: 10, PreBookingPickupToContactMinutes: 3, ReplaceTelephoneCountryCode: false, ReplaceTelephoneCountryCodeWith: "0", ReviewMasterSMSGatewayEnabled: true, ReviewMasterSMSGatewayUseMasterQueue: false, ReviewMasterSMSGatewayPairCode: "1234", Companies: "", BookingSourceMobileAppState: -1}
	// fmt.Printf("grcftwc: %+v\n", grcftwc)
	fmt.Println(SendReviewMasterSMSGateway("+4471234567890", "testing", grcftwc))
}

func TestSendReviewMasterSMSGateway2(t *testing.T) {
	// Send to Master Queue
	grcftwc := database.GoogleReviewsConfigFromTokenWithChecks{MinSendFrequency: 21, MaxSendCount: 10, MaxDailySendCount: 20, TelephoneParameter: "t", SendFromIcabbiApp: false, AppKey: "Digital", SecretKey: "Digicomms1!", SendURL: "", HttpGet: false, SendSuccessResponse: "returnSendSms=success", Start: "00:00", End: "23:59", Sunday: true, Monday: true, Tuesday: true, Wednesday: true, Thursday: true, Friday: true, Saturday: true, TimeZone: "Europe/London", ClientID: 11, Country: "GB", MultiMessageEnabled: 0, MessageParameter: "m", MultiMessageSeparator: "SSSSS", UseDatabaseMessage: 1, Message: "Hope you enjoyed your journey", SendDelayEnabled: false, SendDelay: 0, DispatcherChecksEnabled: false, DispatcherURL: "https://ghost-main-static-b36cb86a19e14a2386de12935fac6526.ghostapi.app:29003/", BookingIdParameter: "b", IsBookingForNowDiffMinutes: 10, BookingNowPickupToContactMinutes: 10, PreBookingPickupToContactMinutes: 3, ReplaceTelephoneCountryCode: false, ReplaceTelephoneCountryCodeWith: "0", ReviewMasterSMSGatewayEnabled: true, ReviewMasterSMSGatewayUseMasterQueue: true, ReviewMasterSMSGatewayPairCode: "1234", Companies: "", BookingSourceMobileAppState: -1}
	// fmt.Printf("grcftwc: %+v\n", grcftwc)
	fmt.Println(SendReviewMasterSMSGateway("+"+testTelephone, "testing", grcftwc))
}

func TestSendSMSServer1(t *testing.T) {
	grcftwc := database.GoogleReviewsConfigFromTokenWithChecks{MinSendFrequency: 21, MaxSendCount: 10, MaxDailySendCount: 20, TelephoneParameter: "t", SendFromIcabbiApp: false, AppKey: "Digital", SecretKey: "Digicomms1!", SendURL: "", HttpGet: false, SendSuccessResponse: "returnSendSms=success", Start: "00:00", End: "23:59", Sunday: true, Monday: true, Tuesday: true, Wednesday: true, Thursday: true, Friday: true, Saturday: true, TimeZone: "Europe/London", ClientID: 11, Country: "GB", MultiMessageEnabled: 0, MessageParameter: "m", MultiMessageSeparator: "SSSSS", UseDatabaseMessage: 1, Message: "Hope you enjoyed your journey", SendDelayEnabled: false, SendDelay: 0, DispatcherChecksEnabled: false, DispatcherURL: "https://ghost-main-static-b36cb86a19e14a2386de12935fac6526.ghostapi.app:29003/", BookingIdParameter: "b", IsBookingForNowDiffMinutes: 10, BookingNowPickupToContactMinutes: 10, PreBookingPickupToContactMinutes: 3, ReplaceTelephoneCountryCode: false, ReplaceTelephoneCountryCodeWith: "0", ReviewMasterSMSGatewayEnabled: false, ReviewMasterSMSGatewayUseMasterQueue: false, ReviewMasterSMSGatewayPairCode: "1234", Companies: "", BookingSourceMobileAppState: -1}
	// fmt.Printf("grcftwc: %+v\n", grcftwc)
	fmt.Println(SendSMSServer("+4471234567890", "testing", grcftwc))
}
