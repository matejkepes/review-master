package database

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	testfixtures "gopkg.in/testfixtures.v2"
)

var fixtures *testfixtures.Context

func TestMain(m *testing.M) {
	var err error

	// Open connection with the test database.
	// Do NOT import fixtures in a production database!
	// Existing data would be deleted
	OpenDB(TestDbName, TestDbAddress, TestDbPort, TestDbUsername, TestDbPassword)

	// creating the context that hold the fixtures
	// see about all compatible databases in this page below
	fixtures, err = testfixtures.NewFolder(Db, &testfixtures.MySQL{}, "fixtures")
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func prepareTestDatabase() {
	// prevent the error message:
	// Loading aborted because the database name does not contain "test"
	testfixtures.SkipDatabaseNameCheck(true)
	if err := fixtures.Load(); err != nil {
		log.Fatal(err)
	}
}

func TestListAllClients(t *testing.T) {
	prepareTestDatabase()
	clients, err := ListAllClients(1)
	if err != nil {
		t.Fatal("error fetching clients, err: ", err)
	}
	fmt.Printf("clients: %+v\n", clients)
	// for _, client := range clients {
	// 	fmt.Printf("%d\n", client.ID)
	// }
	if len(clients) != 5 {
		t.Fatal("error fetching clients there should be 4 but got ", len(clients))
	}
}

func TestGetSimpleClient(t *testing.T) {
	prepareTestDatabase()
	simpleConfig, err := GetSimpleClient(1, 1)
	if err != nil {
		t.Fatal("error fetching simple config for client, err: ", err)
	}
	fmt.Printf("simpleConfig: %+v\n", simpleConfig)
	// for _, client := range clients {
	// 	fmt.Printf("%d\n", client.ID)
	// }
}

func TestUpdateSimpleClient(t *testing.T) {
	prepareTestDatabase()

	var simpleConfig SimpleConfig
	simpleConfig.ClientID = 1
	simpleConfig.ClientEnabled = true
	simpleConfig.ClientName = "Taxi Company 1"
	simpleConfig.ClientNote = "Note 1"
	simpleConfig.ClientCountry = "GB"
	simpleConfig.GoogleReviewsConfigID = 1
	simpleConfig.GoogleReviewsConfigEnabled = true
	simpleConfig.GoogleReviewsConfigMinSendFrequency = 21
	simpleConfig.GoogleReviewsConfigMaxSendCount = 10
	simpleConfig.GoogleReviewsConfigMaxDailySendCount = 2
	simpleConfig.GoogleReviewsConfigToken = "QxrH0iJc3wv/lj/YKVppNYRad7tN0Z3x"
	simpleConfig.GoogleReviewsConfigTelephoneParameter = "t"
	simpleConfig.GoogleReviewsConfigSendFromIcabbiApp = true
	simpleConfig.GoogleReviewsConfigAppKey = TestAppKey
	simpleConfig.GoogleReviewsConfigSecretKey = TestSecretKey
	simpleConfig.GoogleReviewsConfigSendURL = "https://test1.veezu.co.uk:8893/api/v1/sendmessage"
	simpleConfig.GoogleReviewsConfigHttpGet = false
	simpleConfig.GoogleReviewsConfigSendSuccessResponse = "returnSendSms=success"
	simpleConfig.GoogleReviewsConfigTimeZone = "Europe/London"
	simpleConfig.GoogleReviewsConfigMultiMessageEnabled = true
	simpleConfig.GoogleReviewsConfigMessageParameter = "m"
	simpleConfig.GoogleReviewsConfigMultiMessageSeparator = "SSSSS"
	simpleConfig.GoogleReviewsConfigUseDatabaseMessage = false
	simpleConfig.GoogleReviewsConfigMessage = "change me"
	simpleConfig.GoogleReviewsConfigSendDelayEnabled = false
	simpleConfig.GoogleReviewsConfigSendDelay = 0
	simpleConfig.GoogleReviewsConfigDispatcherChecksEnabled = false
	simpleConfig.GoogleReviewsConfigDispatcherType = "ICABBI"
	simpleConfig.GoogleReviewsConfigDispatcherURL = ""
	simpleConfig.GoogleReviewsConfigBookingIdParameter = "b"
	simpleConfig.GoogleReviewsConfigIsBookingForNowDiffMinutes = 10
	simpleConfig.GoogleReviewsConfigBookingNowPickupToContactMinutes = 10
	simpleConfig.GoogleReviewsConfigPreBookingPickupToContactMinutes = 3
	simpleConfig.GoogleReviewsConfigReplaceTelephoneCountryCode = false
	simpleConfig.GoogleReviewsConfigReplaceTelephoneCountryCodeWith = "0"
	simpleConfig.GoogleReviewsConfigReviewMasterSMSGatewayEnabled = false
	simpleConfig.GoogleReviewsConfigReviewMasterSMSGatewayUseMasterQueue = false
	simpleConfig.GoogleReviewsConfigReviewMasterSMSGatewayPairCode = "0"
	simpleConfig.GoogleReviewsConfigAlternateMessageServiceEnabled = false
	simpleConfig.GoogleReviewsConfigAlternateMessageService = ""
	simpleConfig.GoogleReviewsConfigAlternateMessageServiceSecret1 = ""
	simpleConfig.GoogleReviewsConfigCompanies = ""
	simpleConfig.GoogleReviewsConfigBookingSourceMobileAppState = -1
	simpleConfig.GoogleReviewsConfigAIResponsesEnabled = false
	simpleConfig.GoogleReviewsConfigContactMethod = ""
	simpleConfig.GoogleReviewsConfigMonthlyReviewAnalysisEnabled = true
	simpleConfig.GoogleMyBusinessReviewReplyEnabled = false
	simpleConfig.GoogleMyBusinessLocationName = "Taxi Company 1"
	simpleConfig.GoogleMyBusinessPostalCode = "AB1 2CD"
	simpleConfig.GoogleMyBusinessReplyToUnspecfifiedStarRating = false
	simpleConfig.GoogleMyBusinessUnspecfifiedStarRatingReply = "Reply for star rating unspecified"
	simpleConfig.GoogleMyBusinessReplyToOneStarRating = false
	simpleConfig.GoogleMyBusinessOneStarRatingReply = "Reply for star rating 1"
	simpleConfig.GoogleMyBusinessReplyToTwoStarRating = false
	simpleConfig.GoogleMyBusinessTwoStarRatingReply = "Reply for star rating 2"
	simpleConfig.GoogleMyBusinessReplyToThreeStarRating = false
	simpleConfig.GoogleMyBusinessThreeStarRatingReply = "Reply for star rating 3"
	simpleConfig.GoogleMyBusinessReplyToFourStarRating = false
	simpleConfig.GoogleMyBusinessFourStarRatingReply = "Reply for star rating 4"
	simpleConfig.GoogleMyBusinessReplyToFiveStarRating = false
	simpleConfig.GoogleMyBusinessFiveStarRatingReply = "Reply for star rating 5"
	simpleConfig.GoogleMyBusinessReportEnabled = false
	simpleConfig.EmailAddress = "test1@test.com"
	simpleConfig.GoogleReviewsConfigTimeID = 1
	simpleConfig.GoogleReviewsConfigTimeEnabled = true
	simpleConfig.GoogleReviewsConfigTimeStart = "09:00"
	simpleConfig.GoogleReviewsConfigTimeEnd = "12:00"
	simpleConfig.GoogleReviewsConfigTimeSunday = true
	simpleConfig.GoogleReviewsConfigTimeMonday = true
	simpleConfig.GoogleReviewsConfigTimeTuesday = true
	simpleConfig.GoogleReviewsConfigTimeWednesday = true
	simpleConfig.GoogleReviewsConfigTimeThursday = true
	simpleConfig.GoogleReviewsConfigTimeFriday = true
	simpleConfig.GoogleReviewsConfigTimeSaturday = true

	err := UpdateSimpleClient(simpleConfig)
	if err != nil {
		t.Fatal("error updating simple config for client, err: ", err)
	}

	s, err := GetSimpleClient(1, 1)
	if err != nil {
		t.Fatal("error fetching simple config for client, err: ", err)
	}
	if s.GoogleReviewsConfigTimeEnd != "12:00" {
		t.Fatal("error updating simple config for client, has not changed")
	}
	if s.EmailAddress != "test1@test.com" {
		t.Fatal("error updating simple config for client, has not changed")
	}
}

func TestCreateSimpleClient(t *testing.T) {
	prepareTestDatabase()

	var simpleConfig SimpleConfig
	simpleConfig.ClientID = 1
	simpleConfig.ClientEnabled = true
	simpleConfig.ClientName = "Taxi Company 3"
	simpleConfig.ClientNote = "Note 3"
	simpleConfig.ClientCountry = "GB"
	simpleConfig.GoogleReviewsConfigEnabled = true
	simpleConfig.GoogleReviewsConfigMinSendFrequency = 21
	simpleConfig.GoogleReviewsConfigMaxSendCount = 10
	simpleConfig.GoogleReviewsConfigMaxDailySendCount = 2
	simpleConfig.GoogleReviewsConfigToken = "GysCRHdmlJphYKU8XGFOwNLy81gPIQkZ"
	simpleConfig.GoogleReviewsConfigTelephoneParameter = "t"
	simpleConfig.GoogleReviewsConfigSendFromIcabbiApp = true
	simpleConfig.GoogleReviewsConfigAppKey = TestAppKey
	simpleConfig.GoogleReviewsConfigSecretKey = TestSecretKey
	simpleConfig.GoogleReviewsConfigSendURL = "https://test1.veezu.co.uk:8893/api/v1/sendmessage"
	simpleConfig.GoogleReviewsConfigHttpGet = false
	simpleConfig.GoogleReviewsConfigSendSuccessResponse = "returnSendSms=success"
	simpleConfig.GoogleReviewsConfigTimeZone = "Europe/London"
	simpleConfig.GoogleReviewsConfigMultiMessageEnabled = true
	simpleConfig.GoogleReviewsConfigMessageParameter = "m"
	simpleConfig.GoogleReviewsConfigMultiMessageSeparator = "SSSSS"
	simpleConfig.GoogleReviewsConfigUseDatabaseMessage = false
	simpleConfig.GoogleReviewsConfigMessage = "change me"
	simpleConfig.GoogleReviewsConfigSendDelayEnabled = false
	simpleConfig.GoogleReviewsConfigSendDelay = 0
	simpleConfig.GoogleReviewsConfigDispatcherChecksEnabled = false
	simpleConfig.GoogleReviewsConfigDispatcherType = "ICABBI"
	simpleConfig.GoogleReviewsConfigDispatcherURL = ""
	simpleConfig.GoogleReviewsConfigBookingIdParameter = "b"
	simpleConfig.GoogleReviewsConfigIsBookingForNowDiffMinutes = 10
	simpleConfig.GoogleReviewsConfigBookingNowPickupToContactMinutes = 10
	simpleConfig.GoogleReviewsConfigPreBookingPickupToContactMinutes = 3
	simpleConfig.GoogleReviewsConfigReplaceTelephoneCountryCode = false
	simpleConfig.GoogleReviewsConfigReplaceTelephoneCountryCodeWith = "0"
	simpleConfig.GoogleReviewsConfigReviewMasterSMSGatewayEnabled = false
	simpleConfig.GoogleReviewsConfigReviewMasterSMSGatewayUseMasterQueue = false
	simpleConfig.GoogleReviewsConfigReviewMasterSMSGatewayPairCode = "12"
	simpleConfig.GoogleReviewsConfigAlternateMessageServiceEnabled = false
	simpleConfig.GoogleReviewsConfigAlternateMessageService = ""
	simpleConfig.GoogleReviewsConfigAlternateMessageServiceSecret1 = ""
	simpleConfig.GoogleReviewsConfigCompanies = ""
	simpleConfig.GoogleReviewsConfigBookingSourceMobileAppState = -1
	simpleConfig.GoogleReviewsConfigAIResponsesEnabled = false
	simpleConfig.GoogleReviewsConfigContactMethod = ""
	simpleConfig.GoogleReviewsConfigMonthlyReviewAnalysisEnabled = true
	simpleConfig.GoogleMyBusinessReviewReplyEnabled = false
	simpleConfig.GoogleMyBusinessLocationName = "Taxi Company 1"
	simpleConfig.GoogleMyBusinessPostalCode = "AB1 2CD"
	simpleConfig.GoogleMyBusinessReplyToUnspecfifiedStarRating = false
	simpleConfig.GoogleMyBusinessUnspecfifiedStarRatingReply = "Reply for star rating unspecified"
	simpleConfig.GoogleMyBusinessReplyToOneStarRating = false
	simpleConfig.GoogleMyBusinessOneStarRatingReply = "Reply for star rating 1"
	simpleConfig.GoogleMyBusinessReplyToTwoStarRating = false
	simpleConfig.GoogleMyBusinessTwoStarRatingReply = "Reply for star rating 2"
	simpleConfig.GoogleMyBusinessReplyToThreeStarRating = false
	simpleConfig.GoogleMyBusinessThreeStarRatingReply = "Reply for star rating 3"
	simpleConfig.GoogleMyBusinessReplyToFourStarRating = false
	simpleConfig.GoogleMyBusinessFourStarRatingReply = "Reply for star rating 4"
	simpleConfig.GoogleMyBusinessReplyToFiveStarRating = false
	simpleConfig.GoogleMyBusinessFiveStarRatingReply = "Reply for star rating 5"
	simpleConfig.GoogleMyBusinessReportEnabled = false
	simpleConfig.EmailAddress = "test2@test.com"
	simpleConfig.GoogleReviewsConfigTimeEnabled = true
	simpleConfig.GoogleReviewsConfigTimeStart = "09:00"
	simpleConfig.GoogleReviewsConfigTimeEnd = "14:00"
	simpleConfig.GoogleReviewsConfigTimeSunday = true
	simpleConfig.GoogleReviewsConfigTimeMonday = true
	simpleConfig.GoogleReviewsConfigTimeTuesday = true
	simpleConfig.GoogleReviewsConfigTimeWednesday = true
	simpleConfig.GoogleReviewsConfigTimeThursday = true
	simpleConfig.GoogleReviewsConfigTimeFriday = true
	simpleConfig.GoogleReviewsConfigTimeSaturday = true

	err := CreateSimpleClient(simpleConfig, 1)
	if err != nil {
		t.Fatal("error creating simple config for client, err: ", err)
	}

	s, err := GetSimpleClient(1, 1)
	if err != nil {
		t.Fatal("error fetching simple config for client, err: ", err)
	}
	fmt.Printf("Config: %+v\n", s)
	// if s.GoogleReviewsConfigTimeEnd != "14:00" {
	// 	t.Fatal("error creating simple config for client, has not changed, it is: ", s.GoogleReviewsConfigTimeEnd)
	// }
	// if s.EmailAddress != "test2@test.com" {
	// 	t.Fatal("error creating simple config for client, has not changed, it is: ", s.EmailAddress)
	// }
}

func TestGetClient(t *testing.T) {
	prepareTestDatabase()
	config, err := GetClient(1, 1)
	if err != nil {
		t.Fatal("error fetching config for client, err: ", err)
	}
	fmt.Printf("Config: %+v\n", config)
}

func TestUpdateClient(t *testing.T) {
	prepareTestDatabase()

	var client Client
	client.ID = 1
	client.Enabled = true
	client.Name = "Taxi Company 1"
	client.Note = "Note 1"
	client.Country = "GB"
	client.PartnerID = 1

	var grc GoogleReviewsConfig
	grc.ID = 1
	grc.Enabled = true
	grc.MinSendFrequency = 21
	grc.MaxSendCount = 10
	grc.MaxDailySendCount = 2
	grc.Token = "QxrH0iJc3wv/lj/YKVppNYRad7tN0Z3x"
	grc.TelephoneParameter = "t"
	grc.SendFromIcabbiApp = false
	grc.AppKey = ""
	grc.SecretKey = ""
	grc.SendURL = "https://test1.veezu.co.uk:8893/api/v1/sendmessage"
	grc.HttpGet = true
	grc.SendSuccessResponse = "returnSendSms=success"
	grc.TimeZone = "Europe/London"
	grc.MultiMessageEnabled = true
	grc.MessageParameter = "m"
	grc.MultiMessageSeparator = "SSSSS"
	grc.UseDatabaseMessage = false
	grc.Message = "change me"
	grc.SendDelayEnabled = false
	grc.SendDelay = 0
	grc.DispatcherChecksEnabled = false
	grc.DispatcherType = "ICABBI"
	grc.DispatcherURL = ""
	grc.BookingIdParameter = "b"
	grc.IsBookingForNowDiffMinutes = 10
	grc.BookingNowPickupToContactMinutes = 10
	grc.PreBookingPickupToContactMinutes = 3
	grc.ReplaceTelephoneCountryCode = false
	grc.ReplaceTelephoneCountryCodeWith = "0"
	grc.ReviewMasterSMSGatewayEnabled = true
	grc.ReviewMasterSMSGatewayUseMasterQueue = false
	grc.ReviewMasterSMSGatewayPairCode = "123"
	grc.AlternateMessageServiceEnabled = false
	grc.AlternateMessageService = ""
	grc.AlternateMessageServiceSecret1 = ""
	grc.Companies = ""
	grc.BookingSourceMobileAppState = -1
	grc.AIResponsesEnabled = false
	grc.ContactMethod = ""
	grc.MonthlyReviewAnalysisEnabled = true
	grc.GoogleMyBusinessReviewReplyEnabled = false
	grc.GoogleMyBusinessLocationName = "Taxi Company 1"
	grc.GoogleMyBusinessPostalCode = "AB1 2CD"
	grc.GoogleMyBusinessReplyToUnspecfifiedStarRating = false
	grc.GoogleMyBusinessUnspecfifiedStarRatingReply = "Reply for star rating unspecified"
	grc.GoogleMyBusinessReplyToOneStarRating = false
	grc.GoogleMyBusinessOneStarRatingReply = "Reply for star rating 1"
	grc.GoogleMyBusinessReplyToTwoStarRating = false
	grc.GoogleMyBusinessTwoStarRatingReply = "Reply for star rating 2"
	grc.GoogleMyBusinessReplyToThreeStarRating = false
	grc.GoogleMyBusinessThreeStarRatingReply = "Reply for star rating 3"
	grc.GoogleMyBusinessReplyToFourStarRating = false
	grc.GoogleMyBusinessFourStarRatingReply = "Reply for star rating 4"
	grc.GoogleMyBusinessReplyToFiveStarRating = false
	grc.GoogleMyBusinessFiveStarRatingReply = "Reply for star rating 5"
	grc.GoogleMyBusinessReportEnabled = false
	grc.EmailAddress = "test@test.com"

	var grct GoogleReviewsConfigTime
	grct.ID = 1
	grct.Enabled = true
	grct.Start = "09:00"
	grct.End = "12:00"
	grct.Sunday = true
	grct.Monday = true
	grct.Tuesday = true
	grct.Wednesday = true
	grct.Thursday = true
	grct.Friday = true
	grct.Saturday = false

	grcts := []GoogleReviewsConfigTime{}
	grcts = append(grcts, grct)

	var config Config
	config.GoogleReviewsConfig = grc
	config.GoogleReviewsConfigTimes = grcts

	configs := []Config{}
	configs = append(configs, config)

	var clientConfig ClientConfig
	clientConfig.Client = client
	clientConfig.Configs = configs

	err := UpdateClient(clientConfig, 1)
	if err != nil {
		t.Fatal("error updating config for client, err: ", err)
	}

	c, err := GetClient(1, 1)
	if err != nil {
		t.Fatal("error fetching config for client, err: ", err)
	}
	fmt.Printf("clientConfig: %+v\n", c)
	if c.Configs[0].GoogleReviewsConfigTimes[0].End != "12:00" {
		t.Fatal("error updating config for client, has not changed")
	}
	if c.Configs[0].GoogleReviewsConfigTimes[0].Saturday != false {
		t.Fatal("error updating config for client, has not changed")
	}
	if c.Configs[0].GoogleReviewsConfig.ReviewMasterSMSGatewayEnabled != true {
		t.Fatal("error updating config for client, has not changed")
	}
	if c.Configs[0].GoogleReviewsConfig.ReviewMasterSMSGatewayUseMasterQueue != false {
		t.Fatal("error updating config for client, has not changed")
	}
	if c.Configs[0].GoogleReviewsConfig.ReviewMasterSMSGatewayPairCode != "123" {
		t.Fatal("error updating config for client, has not changed")
	}
}

func TestCreateClient(t *testing.T) {
	prepareTestDatabase()

	var client Client
	client.ID = 1
	client.Enabled = true
	client.Name = "Taxi Company 3"
	client.Note = "Note 3"
	client.Country = "GB"
	client.PartnerID = 1

	var grc GoogleReviewsConfig
	grc.ID = 1
	grc.Enabled = true
	grc.MinSendFrequency = 21
	grc.MaxSendCount = 10
	grc.MaxDailySendCount = 2
	grc.Token = "QxrH0iJc3wv/lj/GysCRHdmlJphYKU8XGFOwNLy81gPIQkZ"
	grc.TelephoneParameter = "t"
	grc.SendFromIcabbiApp = false
	grc.AppKey = ""
	grc.SecretKey = ""
	grc.SendURL = "https://test1.veezu.co.uk:8893/api/v1/sendmessage"
	grc.HttpGet = false
	grc.SendSuccessResponse = "returnSendSms=success"
	grc.TimeZone = "Europe/London"
	grc.MultiMessageEnabled = true
	grc.MessageParameter = "m"
	grc.MultiMessageSeparator = "SSSSS"
	grc.UseDatabaseMessage = false
	grc.Message = "change me"
	grc.SendDelayEnabled = false
	grc.SendDelay = 0
	grc.DispatcherChecksEnabled = false
	grc.DispatcherType = "ICABBI"
	grc.DispatcherURL = ""
	grc.BookingIdParameter = "b"
	grc.IsBookingForNowDiffMinutes = 10
	grc.BookingNowPickupToContactMinutes = 10
	grc.PreBookingPickupToContactMinutes = 3
	grc.ReplaceTelephoneCountryCode = false
	grc.ReplaceTelephoneCountryCodeWith = "0"
	grc.ReviewMasterSMSGatewayEnabled = false
	grc.ReviewMasterSMSGatewayUseMasterQueue = false
	grc.ReviewMasterSMSGatewayPairCode = "12345"
	grc.AlternateMessageServiceEnabled = false
	grc.AlternateMessageService = ""
	grc.AlternateMessageServiceSecret1 = ""
	grc.Companies = ""
	grc.BookingSourceMobileAppState = -1
	grc.AIResponsesEnabled = false
	grc.ContactMethod = ""
	grc.MonthlyReviewAnalysisEnabled = true
	grc.GoogleMyBusinessReviewReplyEnabled = false
	grc.GoogleMyBusinessLocationName = "Taxi Company 1"
	grc.GoogleMyBusinessPostalCode = "AB1 2CD"
	grc.GoogleMyBusinessReplyToUnspecfifiedStarRating = false
	grc.GoogleMyBusinessUnspecfifiedStarRatingReply = "Reply for star rating unspecified"
	grc.GoogleMyBusinessReplyToOneStarRating = false
	grc.GoogleMyBusinessOneStarRatingReply = "Reply for star rating 1"
	grc.GoogleMyBusinessReplyToTwoStarRating = false
	grc.GoogleMyBusinessTwoStarRatingReply = "Reply for star rating 2"
	grc.GoogleMyBusinessReplyToThreeStarRating = false
	grc.GoogleMyBusinessThreeStarRatingReply = "Reply for star rating 3"
	grc.GoogleMyBusinessReplyToFourStarRating = false
	grc.GoogleMyBusinessFourStarRatingReply = "Reply for star rating 4"
	grc.GoogleMyBusinessReplyToFiveStarRating = false
	grc.GoogleMyBusinessFiveStarRatingReply = "Reply for star rating 5"
	grc.GoogleMyBusinessReportEnabled = false
	grc.EmailAddress = "test@test.com"

	var grct GoogleReviewsConfigTime
	grct.ID = 1
	grct.Enabled = true
	grct.Start = "09:00"
	grct.End = "14:00"
	grct.Sunday = false
	grct.Monday = true
	grct.Tuesday = true
	grct.Wednesday = true
	grct.Thursday = true
	grct.Friday = true
	grct.Saturday = false

	grcts := []GoogleReviewsConfigTime{}
	grcts = append(grcts, grct)

	var config Config
	config.GoogleReviewsConfig = grc
	config.GoogleReviewsConfigTimes = grcts

	configs := []Config{}
	configs = append(configs, config)

	var clientConfig ClientConfig
	clientConfig.Client = client
	clientConfig.Configs = configs

	err := CreateClient(clientConfig, 1)
	if err != nil {
		t.Fatal("error creating config for client, err: ", err)
	}

	c, err := GetClient(1, 1)
	if err != nil {
		t.Fatal("error fetching config for client, err: ", err)
	}
	fmt.Printf("clientConfig: %+v\n", c)
	if c.Configs[0].GoogleReviewsConfigTimes[0].End != "14:00" {
		t.Fatal("error creating config for client")
	}
}

func TestCreateGRConfig(t *testing.T) {
	prepareTestDatabase()

	var grc GoogleReviewsConfig
	grc.ID = 1
	grc.Enabled = true
	grc.MinSendFrequency = 21
	grc.MaxSendCount = 10
	grc.MaxDailySendCount = 2
	grc.Token = "QxrH0iJc3wv/lj/GysCRHdmlJphYKU8XGFOwNLy81gPIQkZ098765"
	grc.TelephoneParameter = "t"
	grc.SendFromIcabbiApp = false
	grc.AppKey = ""
	grc.SecretKey = ""
	grc.SendURL = "https://test1.veezu.co.uk:8893/api/v1/sendmessage"
	grc.HttpGet = false
	grc.SendSuccessResponse = "returnSendSms=success"
	grc.TimeZone = "Europe/London"
	grc.MultiMessageEnabled = true
	grc.MessageParameter = "m"
	grc.MultiMessageSeparator = "SSSSS"
	grc.UseDatabaseMessage = false
	grc.Message = "change me"
	grc.SendDelayEnabled = false
	grc.SendDelay = 0
	grc.DispatcherChecksEnabled = false
	grc.DispatcherType = "ICABBI"
	grc.DispatcherURL = ""
	grc.BookingIdParameter = "b"
	grc.IsBookingForNowDiffMinutes = 10
	grc.BookingNowPickupToContactMinutes = 10
	grc.PreBookingPickupToContactMinutes = 3
	grc.ReplaceTelephoneCountryCode = false
	grc.ReplaceTelephoneCountryCodeWith = "0"
	grc.ReviewMasterSMSGatewayEnabled = false
	grc.ReviewMasterSMSGatewayUseMasterQueue = false
	grc.ReviewMasterSMSGatewayPairCode = "12345dhdhgs"
	grc.AlternateMessageServiceEnabled = false
	grc.AlternateMessageService = ""
	grc.AlternateMessageServiceSecret1 = ""
	grc.Companies = ""
	grc.BookingSourceMobileAppState = -1
	grc.AIResponsesEnabled = false
	grc.ContactMethod = ""
	grc.MonthlyReviewAnalysisEnabled = true
	grc.GoogleMyBusinessReviewReplyEnabled = false
	grc.GoogleMyBusinessLocationName = "Taxi Company 1"
	grc.GoogleMyBusinessPostalCode = "AB1 2CD"
	grc.GoogleMyBusinessReplyToUnspecfifiedStarRating = false
	grc.GoogleMyBusinessUnspecfifiedStarRatingReply = "Reply for star rating unspecified"
	grc.GoogleMyBusinessReplyToOneStarRating = false
	grc.GoogleMyBusinessOneStarRatingReply = "Reply for star rating 1"
	grc.GoogleMyBusinessReplyToTwoStarRating = false
	grc.GoogleMyBusinessTwoStarRatingReply = "Reply for star rating 2"
	grc.GoogleMyBusinessReplyToThreeStarRating = false
	grc.GoogleMyBusinessThreeStarRatingReply = "Reply for star rating 3"
	grc.GoogleMyBusinessReplyToFourStarRating = false
	grc.GoogleMyBusinessFourStarRatingReply = "Reply for star rating 4"
	grc.GoogleMyBusinessReplyToFiveStarRating = false
	grc.GoogleMyBusinessFiveStarRatingReply = "Reply for star rating 5"
	grc.GoogleMyBusinessReportEnabled = false
	grc.EmailAddress = "test@test.com"
	grc.ClientID = 1

	err := CreateGRConfig(grc)
	if err != nil {
		t.Fatal("error creating config, err: ", err)
	}

	c, err := GetClient(1, 1)
	if err != nil {
		t.Fatal("error fetching config for client, err: ", err)
	}
	fmt.Printf("clientConfig: %+v\n", c)
	if len(c.Configs) < 2 {
		t.Fatal("error creating config, there should be at least 2 configs")
	}
}

func TestStats1(t *testing.T) {
	prepareTestDatabase()

	var statsRequest StatsRequest
	b := time.Now().Add(-time.Hour * 24 * 50)
	b = time.Date(b.Year(), b.Month(), b.Day(), 0, 0, 0, 0, b.Location())
	statsRequest.StartDay = b.String()
	e := time.Now().Add(time.Hour * 24)
	e = time.Date(e.Year(), e.Month(), e.Day(), 0, 0, 0, 0, e.Location())
	statsRequest.EndDay = e.String()
	s, err := Stats(statsRequest, 1)
	if err != nil {
		t.Fatal("error getting stats, err: ", err)
	}
	fmt.Printf("stats: %+v\n", s)
}

func TestStats2(t *testing.T) {
	prepareTestDatabase()

	var statsRequest StatsRequest
	b := time.Now().Add(-time.Hour * 24 * 50)
	b = time.Date(b.Year(), b.Month(), b.Day(), 0, 0, 0, 0, b.Location())
	statsRequest.StartDay = b.String()
	e := time.Now().Add(time.Hour * 24)
	e = time.Date(e.Year(), e.Month(), e.Day(), 0, 0, 0, 0, e.Location())
	statsRequest.EndDay = e.String()
	statsRequest.TimeGrouping = "Week"
	s, err := Stats(statsRequest, 1)
	if err != nil {
		t.Fatal("error getting stats, err: ", err)
	}
	fmt.Printf("stats: %+v\n", s)
}

func TestStats3(t *testing.T) {
	prepareTestDatabase()

	var statsRequest StatsRequest
	b := time.Now().Add(-time.Hour * 24 * 50)
	b = time.Date(b.Year(), b.Month(), b.Day(), 0, 0, 0, 0, b.Location())
	statsRequest.StartDay = b.String()
	e := time.Now().Add(time.Hour * 24)
	e = time.Date(e.Year(), e.Month(), e.Day(), 0, 0, 0, 0, e.Location())
	statsRequest.EndDay = e.String()
	statsRequest.TimeGrouping = "Month"
	s, err := Stats(statsRequest, 1)
	if err != nil {
		t.Fatal("error getting stats, err: ", err)
	}
	fmt.Printf("stats: %+v\n", s)
}

func TestStats4(t *testing.T) {
	prepareTestDatabase()

	var statsRequest StatsRequest
	b := time.Now().Add(-time.Hour * 24 * 50)
	b = time.Date(b.Year(), b.Month(), b.Day(), 0, 0, 0, 0, b.Location())
	statsRequest.StartDay = b.String()
	e := time.Now().Add(time.Hour * 24)
	e = time.Date(e.Year(), e.Month(), e.Day(), 0, 0, 0, 0, e.Location())
	statsRequest.EndDay = e.String()
	statsRequest.TimeGrouping = "Year"
	s, err := Stats(statsRequest, 1)
	if err != nil {
		t.Fatal("error getting stats, err: ", err)
	}
	fmt.Printf("stats: %+v\n", s)
}

func TestStatsNew1(t *testing.T) {
	prepareTestDatabase()

	var statsRequest StatsRequest
	b := time.Now().Add(-time.Hour * 24 * 50)
	b = time.Date(b.Year(), b.Month(), b.Day(), 0, 0, 0, 0, b.Location())
	statsRequest.StartDay = b.String()
	e := time.Now().Add(time.Hour * 24)
	e = time.Date(e.Year(), e.Month(), e.Day(), 0, 0, 0, 0, e.Location())
	statsRequest.EndDay = e.String()
	s, err := StatsNew(statsRequest, 1)
	if err != nil {
		t.Fatal("error getting stats, err: ", err)
	}
	fmt.Printf("stats: %+v\n", s)
}

func TestStatsNew2(t *testing.T) {
	prepareTestDatabase()

	var statsRequest StatsRequest
	b := time.Now().Add(-time.Hour * 24 * 50)
	b = time.Date(b.Year(), b.Month(), b.Day(), 0, 0, 0, 0, b.Location())
	statsRequest.StartDay = b.String()
	e := time.Now().Add(time.Hour * 24)
	e = time.Date(e.Year(), e.Month(), e.Day(), 0, 0, 0, 0, e.Location())
	statsRequest.EndDay = e.String()
	statsRequest.TimeGrouping = "Week"
	s, err := StatsNew(statsRequest, 1)
	if err != nil {
		t.Fatal("error getting stats, err: ", err)
	}
	fmt.Printf("stats: %+v\n", s)
}

func TestStatsNew3(t *testing.T) {
	prepareTestDatabase()

	var statsRequest StatsRequest
	b := time.Now().Add(-time.Hour * 24 * 50)
	b = time.Date(b.Year(), b.Month(), b.Day(), 0, 0, 0, 0, b.Location())
	statsRequest.StartDay = b.String()
	e := time.Now().Add(time.Hour * 24)
	e = time.Date(e.Year(), e.Month(), e.Day(), 0, 0, 0, 0, e.Location())
	statsRequest.EndDay = e.String()
	statsRequest.TimeGrouping = "Month"
	s, err := StatsNew(statsRequest, 1)
	if err != nil {
		t.Fatal("error getting stats, err: ", err)
	}
	fmt.Printf("stats: %+v\n", s)
}

func TestStatsNew4(t *testing.T) {
	prepareTestDatabase()

	var statsRequest StatsRequest
	b := time.Now().Add(-time.Hour * 24 * 50)
	b = time.Date(b.Year(), b.Month(), b.Day(), 0, 0, 0, 0, b.Location())
	statsRequest.StartDay = b.String()
	e := time.Now().Add(time.Hour * 24)
	e = time.Date(e.Year(), e.Month(), e.Day(), 0, 0, 0, 0, e.Location())
	statsRequest.EndDay = e.String()
	statsRequest.TimeGrouping = "Year"
	s, err := StatsNew(statsRequest, 1)
	if err != nil {
		t.Fatal("error getting stats, err: ", err)
	}
	fmt.Printf("stats: %+v\n", s)
}

func TestCheckNothingSent1(t *testing.T) {
	prepareTestDatabase()

	// s, err := CheckNothingSent(24, 1)
	s, err := CheckNothingSent(0, 1)
	if err != nil {
		t.Fatal("error getting stats, err: ", err)
	}
	fmt.Printf("stats: %+v\n", s)
}

func TestListAllUsers(t *testing.T) {
	prepareTestDatabase()
	userClients, err := ListAllUsers(1)
	if err != nil {
		t.Fatal("error fetching user clients, err: ", err)
	}
	fmt.Printf("user clients: %+v\n", userClients)
	if len(userClients) != 2 {
		t.Fatal("error fetching user clients there should be 2 but got", len(userClients))
	}
}

func TestGetUser(t *testing.T) {
	prepareTestDatabase()
	user, err := GetUser(1, 1)
	if err != nil {
		t.Fatal("error fetching user, err: ", err)
	}
	fmt.Printf("user: %+v\n", user)
	if len(user.Clients) != 2 {
		t.Fatal("error fetching user clients there should be 2 but got", len(user.Clients))
	}
}

func TestGetUserNotExist(t *testing.T) {
	prepareTestDatabase()
	user, err := GetUser(100, 1)
	if err != nil {
		t.Fatal("error fetching user, err: ", err)
	}
	fmt.Printf("user: %+v\n", user)
	if user.User.ID != 0 {
		t.Fatal("error fetching user clients it should not have been found user ID", user.User.ID)
	}
}

func TestGetUserWrongPartner(t *testing.T) {
	prepareTestDatabase()
	user, err := GetUser(2, 2)
	if err != nil {
		t.Fatal("error fetching user, err: ", err)
	}
	fmt.Printf("user: %+v\n", user)
	if user.User.ID != 0 {
		t.Fatal("error fetching user clients it should not have been found user ID", user.User.ID)
	}
}

func TestUpdateUser(t *testing.T) {
	prepareTestDatabase()
	var ucs UserClients
	var user User
	user.ID = 1
	user.Email = "test@testing.com"
	user.Password = "somerubbish"
	var clients []Client
	var client1 Client
	client1.ID = 1
	client1.Enabled = true
	client1.Name = "Taxi Company 1"
	client1.Note = ""
	client1.Country = "GB"
	clients = append(clients, client1)
	var client2 Client
	client2.ID = 2
	client2.Enabled = true
	client2.Name = "Taxi Company 2 DiSC"
	client2.Note = ""
	client2.Country = "GB"
	clients = append(clients, client2)
	ucs.User = user
	ucs.Clients = clients
	err := UpdateUser(ucs, 1)
	if err != nil {
		t.Fatalf("error: %v\n", err)
	}
}

func TestCreateUser(t *testing.T) {
	prepareTestDatabase()
	var ucs UserClients
	var user User
	user.Email = "test2@testing.com"
	user.Password = "nothing"
	var clients []Client
	var client1 Client
	client1.ID = 1
	clients = append(clients, client1)
	var client2 Client
	client2.ID = 6
	clients = append(clients, client2)
	ucs.User = user
	ucs.Clients = clients
	err := CreateUser(ucs, 1)
	if err != nil {
		t.Fatalf("error: %v\n", err)
	}
}

func TestCreateUserClientNotFound(t *testing.T) {
	prepareTestDatabase()
	var ucs UserClients
	var user User
	user.Email = "test2@testing.com"
	user.Password = "nothing"
	var clients []Client
	var client1 Client
	client1.ID = 1
	clients = append(clients, client1)
	var client2 Client
	client2.ID = 5
	clients = append(clients, client2)
	ucs.User = user
	ucs.Clients = clients
	err := CreateUser(ucs, 1)
	if err == nil {
		t.Fatalf("error should have been incorrect client\n")
	}
}

func TestDeleteUser(t *testing.T) {
	prepareTestDatabase()
	err := DeleteUser(1, 1)
	if err != nil {
		t.Fatalf("error: %v\n", err)
	}
}

func TestDeleteUserWrongPartner(t *testing.T) {
	prepareTestDatabase()
	err := DeleteUser(2, 2)
	if err == nil {
		t.Fatalf("error: %v\n", err)
	}
}
