package database

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"log"
	"net/url"
	"time"

	"google_reviews/utils"
	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

var (
	Db                         *sql.DB
	ReviewMasterSMSMasterQueue = 0
)

// GoogleReviewsConfig - represents a google reviews config from the token with some checks
type GoogleReviewsConfigFromTokenWithChecks struct {
	MinSendFrequency                     uint
	MaxSendCount                         uint
	MaxDailySendCount                    uint
	TelephoneParameter                   string
	SendFromIcabbiApp                    bool
	AppKey                               string
	SecretKey                            string
	SendURL                              string
	HttpGet                              bool
	SendSuccessResponse                  string
	Start                                string
	End                                  string
	Sunday                               bool
	Monday                               bool
	Tuesday                              bool
	Wednesday                            bool
	Thursday                             bool
	Friday                               bool
	Saturday                             bool
	TimeZone                             string
	ClientID                             uint64
	Country                              string
	MultiMessageEnabled                  uint
	MessageParameter                     string
	MultiMessageSeparator                string
	UseDatabaseMessage                   uint
	Message                              string
	SendDelayEnabled                     bool
	SendDelay                            uint
	DispatcherChecksEnabled              bool
	DispatcherURL                        string
	DispatcherType                       string
	BookingIdParameter                   string
	IsBookingForNowDiffMinutes           uint
	BookingNowPickupToContactMinutes     uint
	PreBookingPickupToContactMinutes     uint
	ReplaceTelephoneCountryCode          bool
	ReplaceTelephoneCountryCodeWith      string
	ReviewMasterSMSGatewayEnabled        bool
	ReviewMasterSMSGatewayUseMasterQueue bool
	ReviewMasterSMSGatewayPairCode       string
	AlternateMessageServiceEnabled       bool
	AlternateMessageService              string
	AlternateMessageServiceSecret1       string
	Companies                            string
	BookingSourceMobileAppState          int
}

// OpenDB - open database connection
func OpenDB(database string, host string, port string, username string, password string) {
	// NOTE: ?parseTime=true which allows DATE and DATETIME database types to be parsed into golang time.Time
	db, err := sql.Open("mysql", username+":"+password+"@tcp("+host+":"+port+")/"+database+"?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	// database connection is pooled and used by many connections, even using defer db.Close() causes issues with these
	// where it indicates the connection is closed.
	// defer db.Close()

	// FIX: issues with this program using all the MySQL connections (max_connections - default 151)
	// this is because there are so many concurrent calls to this program and by default there is
	// no limit on the number of open connections (in-use + idle) at the same time
	//
	// Set the maximum number of concurrently open connections (in-use + idle)
	// to 40. Setting this to less than or equal to 0 will mean there is no
	// maximum limit (which is also the default setting).
	db.SetMaxOpenConns(40)

	// Test connection to the database (this will result in the program ending if unable to establish a connection)
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}
	Db = db
}

// ConfigFromTokenWithChecks - get the config from the token with some checks
// ignoreTimeAndSentCountCheck - ignores the time and daily sent count checks (used for testing on front end)
func ConfigFromTokenWithChecks(token string, ignoreTimeAndSentCountCheck bool) GoogleReviewsConfigFromTokenWithChecks {
	qry := "SELECT config.min_send_frequency, config.max_send_count, config.max_daily_send_count, config.telephone_parameter," +
		" config.send_from_icabbi_app, config.app_key, config.secret_key," +
		" config.send_url, config.http_get, config.send_success_response, times.start, times.end," +
		" times.sunday, times.monday, times.tuesday, times.wednesday, times.thursday, times.friday, times.saturday," +
		" config.time_zone, client.id, client.country," +
		" config.multi_message_enabled, config.message_parameter, config.multi_message_separator," +
		" config.use_database_message, config.message," +
		" config.send_delay_enabled, config.send_delay," +
		" config.dispatcher_checks_enabled, config.dispatcher_url, config.dispatcher_type, config.booking_id_parameter, config.is_booking_for_now_diff_minutes," +
		" config.booking_now_pickup_to_contact_minutes, config.pre_booking_pickup_to_contact_minutes," +
		" config.replace_telephone_country_code, config.replace_telephone_country_code_with," +
		" config.review_master_sms_gateway_enabled, config.review_master_sms_gateway_use_master_queue, config.review_master_sms_gateway_pair_code," +
		" config.alternate_message_service_enabled, config.alternate_message_service, config.alternate_message_service_secret1," +
		" config.companies, config.booking_source_mobile_app_state" +
		" FROM google_reviews_config_times AS times" +
		" JOIN google_reviews_configs AS config ON config.id = times.google_reviews_config_id" +
		" JOIN clients AS client ON client.id = config.client_id" +
		" WHERE config.token = ?" +
		" AND times.enabled = 1" +
		" AND config.enabled = 1" +
		" AND client.enabled = 1"
	rows, err := Db.Query(qry, token)
	var grcftwc GoogleReviewsConfigFromTokenWithChecks
	if err != nil {
		log.Println("Error retrieving token", token, "from database. Error: ", err)
		return grcftwc
	}
	defer rows.Close()
	for rows.Next() {
		if err1 := rows.Scan(&grcftwc.MinSendFrequency, &grcftwc.MaxSendCount, &grcftwc.MaxDailySendCount,
			&grcftwc.TelephoneParameter, &grcftwc.SendFromIcabbiApp, &grcftwc.AppKey, &grcftwc.SecretKey,
			&grcftwc.SendURL, &grcftwc.HttpGet, &grcftwc.SendSuccessResponse, &grcftwc.Start, &grcftwc.End,
			&grcftwc.Sunday, &grcftwc.Monday, &grcftwc.Tuesday, &grcftwc.Wednesday, &grcftwc.Thursday, &grcftwc.Friday,
			&grcftwc.Saturday, &grcftwc.TimeZone, &grcftwc.ClientID, &grcftwc.Country,
			&grcftwc.MultiMessageEnabled, &grcftwc.MessageParameter, &grcftwc.MultiMessageSeparator,
			&grcftwc.UseDatabaseMessage, &grcftwc.Message,
			&grcftwc.SendDelayEnabled, &grcftwc.SendDelay,
			&grcftwc.DispatcherChecksEnabled,
			&grcftwc.DispatcherURL, &grcftwc.DispatcherType, &grcftwc.BookingIdParameter, &grcftwc.IsBookingForNowDiffMinutes,
			&grcftwc.BookingNowPickupToContactMinutes, &grcftwc.PreBookingPickupToContactMinutes,
			&grcftwc.ReplaceTelephoneCountryCode, &grcftwc.ReplaceTelephoneCountryCodeWith,
			&grcftwc.ReviewMasterSMSGatewayEnabled, &grcftwc.ReviewMasterSMSGatewayUseMasterQueue, &grcftwc.ReviewMasterSMSGatewayPairCode,
			&grcftwc.AlternateMessageServiceEnabled, &grcftwc.AlternateMessageService, &grcftwc.AlternateMessageServiceSecret1,
			&grcftwc.Companies, &grcftwc.BookingSourceMobileAppState); err1 != nil {
			log.Println("Error retrieving token", token, "from database whilst reading returned results. Error: ", err1)
			return grcftwc
		}
		if grcftwc.ClientID == 0 {
			log.Printf("token %s not found", token)
			continue
		}
		if !ignoreTimeAndSentCountCheck {
			// check within start and end time and weekday
			if !(utils.CheckTime(grcftwc.Start, grcftwc.End, grcftwc.TimeZone) &&
				utils.CheckWeekday(grcftwc.Sunday, grcftwc.Monday, grcftwc.Tuesday, grcftwc.Wednesday, grcftwc.Thursday, grcftwc.Friday, grcftwc.Saturday, grcftwc.TimeZone)) {
				// log.Printf("token %s not found between times %s and %s for time zone %s for clientID %d", token, start, end, timeZone, clientID)
				grcftwc.MinSendFrequency = 0
				grcftwc.MaxSendCount = 0
				grcftwc.TelephoneParameter = ""
				grcftwc.SendFromIcabbiApp = false
				grcftwc.AppKey = ""
				grcftwc.SecretKey = ""
				grcftwc.SendURL = ""
				grcftwc.HttpGet = false
				grcftwc.SendSuccessResponse = ""
				grcftwc.ClientID = 0
				grcftwc.Country = ""
				grcftwc.MultiMessageEnabled = 0
				grcftwc.MessageParameter = ""
				grcftwc.MultiMessageSeparator = ""
				grcftwc.UseDatabaseMessage = 0
				grcftwc.Message = ""
				grcftwc.SendDelayEnabled = false
				grcftwc.SendDelay = 0
				grcftwc.DispatcherChecksEnabled = false
				grcftwc.DispatcherURL = ""
				grcftwc.DispatcherType = ""
				grcftwc.BookingIdParameter = ""
				grcftwc.IsBookingForNowDiffMinutes = 0
				grcftwc.BookingNowPickupToContactMinutes = 0
				grcftwc.PreBookingPickupToContactMinutes = 0
				grcftwc.ReplaceTelephoneCountryCode = false
				grcftwc.ReplaceTelephoneCountryCodeWith = ""
				grcftwc.ReviewMasterSMSGatewayEnabled = false
				grcftwc.ReviewMasterSMSGatewayUseMasterQueue = false
				grcftwc.ReviewMasterSMSGatewayPairCode = ""
				grcftwc.AlternateMessageServiceEnabled = false
				grcftwc.AlternateMessageService = ""
				grcftwc.AlternateMessageServiceSecret1 = ""
				grcftwc.Companies = ""
				grcftwc.BookingSourceMobileAppState = -1
				continue
			}
			// check whether sent daily allowance
			sentCount := DailySentCount(grcftwc.ClientID)
			if sentCount+1 > grcftwc.MaxDailySendCount {
				// log.Printf("token %s has reached maximum daily send count of %d for clientID %d", token, maxDailySendCount, clientID)
				grcftwc.MinSendFrequency = 0
				grcftwc.MaxSendCount = 0
				grcftwc.TelephoneParameter = ""
				grcftwc.SendFromIcabbiApp = false
				grcftwc.AppKey = ""
				grcftwc.SecretKey = ""
				grcftwc.SendURL = ""
				grcftwc.HttpGet = false
				grcftwc.SendSuccessResponse = ""
				grcftwc.ClientID = 0
				grcftwc.Country = ""
				grcftwc.MultiMessageEnabled = 0
				grcftwc.MessageParameter = ""
				grcftwc.MultiMessageSeparator = ""
				grcftwc.UseDatabaseMessage = 0
				grcftwc.Message = ""
				grcftwc.SendDelayEnabled = false
				grcftwc.SendDelay = 0
				grcftwc.DispatcherChecksEnabled = false
				grcftwc.DispatcherURL = ""
				grcftwc.DispatcherType = ""
				grcftwc.BookingIdParameter = ""
				grcftwc.IsBookingForNowDiffMinutes = 0
				grcftwc.BookingNowPickupToContactMinutes = 0
				grcftwc.PreBookingPickupToContactMinutes = 0
				grcftwc.ReplaceTelephoneCountryCode = false
				grcftwc.ReplaceTelephoneCountryCodeWith = ""
				grcftwc.ReviewMasterSMSGatewayEnabled = false
				grcftwc.ReviewMasterSMSGatewayUseMasterQueue = false
				grcftwc.ReviewMasterSMSGatewayPairCode = ""
				grcftwc.AlternateMessageServiceEnabled = false
				grcftwc.AlternateMessageService = ""
				grcftwc.AlternateMessageServiceSecret1 = ""
				grcftwc.Companies = ""
				grcftwc.BookingSourceMobileAppState = -1
				break
			}
		}
		break
	}
	return grcftwc
}

// GetAutocabConfigsWithChecks - get list of Autocab configs with some checks, these are then used to make request to dispatchers (polling)
// ignoreTimeAndSentCountCheck - ignores the time and daily sent count checks (used for testing on front end)
func GetAutocabConfigsWithChecks(ignoreTimeAndSentCountCheck bool) []GoogleReviewsConfigFromTokenWithChecks {
	qry := "SELECT config.min_send_frequency, config.max_send_count, config.max_daily_send_count, config.telephone_parameter," +
		" config.send_from_icabbi_app, config.app_key, config.secret_key," +
		" config.send_url, config.http_get, config.send_success_response, times.start, times.end," +
		" times.sunday, times.monday, times.tuesday, times.wednesday, times.thursday, times.friday, times.saturday," +
		" config.time_zone, client.id, client.country," +
		" config.multi_message_enabled, config.message_parameter, config.multi_message_separator," +
		" config.use_database_message, config.message," +
		" config.send_delay_enabled, config.send_delay," +
		" config.dispatcher_checks_enabled, config.dispatcher_url, config.dispatcher_type, config.booking_id_parameter, config.is_booking_for_now_diff_minutes," +
		" config.booking_now_pickup_to_contact_minutes, config.pre_booking_pickup_to_contact_minutes," +
		" config.replace_telephone_country_code, config.replace_telephone_country_code_with," +
		" config.review_master_sms_gateway_enabled, config.review_master_sms_gateway_use_master_queue, config.review_master_sms_gateway_pair_code," +
		" config.alternate_message_service_enabled, config.alternate_message_service, config.alternate_message_service_secret1," +
		" config.companies, config.booking_source_mobile_app_state" +
		" FROM google_reviews_config_times AS times" +
		" JOIN google_reviews_configs AS config ON config.id = times.google_reviews_config_id" +
		" JOIN clients AS client ON client.id = config.client_id" +
		" WHERE config.dispatcher_type = 'AUTOCAB'" +
		" AND times.enabled = 1" +
		" AND config.enabled = 1" +
		" AND client.enabled = 1"
	rows, err := Db.Query(qry)
	var grcftwcs []GoogleReviewsConfigFromTokenWithChecks
	if err != nil {
		log.Println("Error retrieving configs for Autocab from database. Error: ", err)
		return grcftwcs
	}
	defer rows.Close()
	for rows.Next() {
		var grcftwc GoogleReviewsConfigFromTokenWithChecks
		if err1 := rows.Scan(&grcftwc.MinSendFrequency, &grcftwc.MaxSendCount, &grcftwc.MaxDailySendCount,
			&grcftwc.TelephoneParameter, &grcftwc.SendFromIcabbiApp, &grcftwc.AppKey, &grcftwc.SecretKey,
			&grcftwc.SendURL, &grcftwc.HttpGet, &grcftwc.SendSuccessResponse, &grcftwc.Start, &grcftwc.End,
			&grcftwc.Sunday, &grcftwc.Monday, &grcftwc.Tuesday, &grcftwc.Wednesday, &grcftwc.Thursday, &grcftwc.Friday,
			&grcftwc.Saturday, &grcftwc.TimeZone, &grcftwc.ClientID, &grcftwc.Country,
			&grcftwc.MultiMessageEnabled, &grcftwc.MessageParameter, &grcftwc.MultiMessageSeparator,
			&grcftwc.UseDatabaseMessage, &grcftwc.Message,
			&grcftwc.SendDelayEnabled, &grcftwc.SendDelay,
			&grcftwc.DispatcherChecksEnabled,
			&grcftwc.DispatcherURL, &grcftwc.DispatcherType, &grcftwc.BookingIdParameter, &grcftwc.IsBookingForNowDiffMinutes,
			&grcftwc.BookingNowPickupToContactMinutes, &grcftwc.PreBookingPickupToContactMinutes,
			&grcftwc.ReplaceTelephoneCountryCode, &grcftwc.ReplaceTelephoneCountryCodeWith,
			&grcftwc.ReviewMasterSMSGatewayEnabled, &grcftwc.ReviewMasterSMSGatewayUseMasterQueue, &grcftwc.ReviewMasterSMSGatewayPairCode,
			&grcftwc.AlternateMessageServiceEnabled, &grcftwc.AlternateMessageService, &grcftwc.AlternateMessageServiceSecret1,
			&grcftwc.Companies, &grcftwc.BookingSourceMobileAppState); err1 != nil {
			log.Println("Error retrieving configs for Autocab from database whilst reading returned results. Error: ", err1)
			return grcftwcs
		}
		if grcftwc.ClientID == 0 {
			continue
		}
		if !ignoreTimeAndSentCountCheck {
			// check within start and end time and weekday
			if !(utils.CheckTime(grcftwc.Start, grcftwc.End, grcftwc.TimeZone) &&
				utils.CheckWeekday(grcftwc.Sunday, grcftwc.Monday, grcftwc.Tuesday, grcftwc.Wednesday, grcftwc.Thursday, grcftwc.Friday, grcftwc.Saturday, grcftwc.TimeZone)) {
				// log.Printf("Autocab config not found between times %s and %s for time zone %s for clientID %d", start, end, timeZone, clientID)
				continue
			}
			// check whether sent daily allowance
			sentCount := DailySentCount(grcftwc.ClientID)
			if sentCount+1 > grcftwc.MaxDailySendCount {
				// log.Printf("Autocab config has reached maximum daily send count of %d for clientID %d", maxDailySendCount, clientID)
				continue
			}
		}
		grcftwcs = append(grcftwcs, grcftwc)
	}
	return grcftwcs
}

// LastSentFromTelephoneAndClient - get the last sent from telephone and client
func LastSentFromTelephoneAndClient(telephone string, clientID uint64) (time.Time, uint, bool, bool) {
	qry := "SELECT last_sent, sent_count, stop FROM google_reviews_last_sents WHERE telephone = ? AND client_id = ?"
	rows, err := Db.Query(qry, telephone, clientID)
	if err != nil {
		log.Println(err)
		return time.Now().AddDate(-3, 0, 0), 0, false, false
	}
	defer rows.Close()
	rows.Next()
	var (
		lastSent  time.Time
		sentCount uint
		stop      bool
	)
	if err1 := rows.Scan(&lastSent, &sentCount, &stop); err1 != nil {
		// fmt.Println("err1: ", err1)
		// time 3 years earlier which will be before the minimum send frequency
		return time.Now().AddDate(-3, 0, 0), 0, false, false
	}
	return lastSent, sentCount, stop, true
}

// UpdateLastSent - update the last sent
// If this has been requested then the stop (sending) will always be false
func UpdateLastSent(telephone string, clientID uint64, sentCount uint) {
	qry := "INSERT INTO google_reviews_last_sents" +
		" (telephone, client_id, last_sent, last_sent_date, sent_count, stop)" +
		" VALUES (?, ?, NOW(), CURDATE(), ?, FALSE)" +
		" ON DUPLICATE KEY UPDATE" +
		" last_sent = NOW()," +
		" last_sent_date = CURDATE()," +
		" sent_count = ?"
	_, err := Db.Exec(qry, telephone, clientID, sentCount, sentCount)
	if err != nil {
		log.Println(err)
	}
}

// StopSending - stop sending messages
func StopSending(telephone string, clientID uint64) {
	qry := "UPDATE google_reviews_last_sents" +
		" SET stop = TRUE" +
		" WHERE telephone = ? and client_id = ?"
	_, err := Db.Exec(qry, telephone, clientID)
	if err != nil {
		log.Println(err)
	}
}

// DailySentCount - daily sent count - used for throttling, prevent too many SMS being sent
func DailySentCount(clientID uint64) uint {
	qry := "SELECT COUNT(id) FROM google_reviews_last_sents WHERE last_sent_date = CURDATE() AND client_id = ?"
	// row := db.QueryRow(qry, clientID)
	var count uint
	// err := row.Scan(&count)
	err := Db.QueryRow(qry, clientID).Scan(&count)
	switch {
	case err == sql.ErrNoRows:
		log.Println(err)
		return 0
	case err != nil:
		log.Println(err)
		return 0
	default:
		return count
	}
}

// QueueIDFromReviewMasterPairCode - get the queue ID (which normally is the client ID) from the
// Review Master SMS Pairing Code with some checks.
// When use master queue is enabled then the queue ID is set to the master queue.
func QueueIDFromReviewMasterPairCode(reviewMasterSmsGatewayPairCode string) int {
	qry := "SELECT client.id, config.review_master_sms_gateway_use_master_queue" +
		" FROM google_reviews_configs AS config" +
		" JOIN clients AS client ON client.id = config.client_id" +
		" WHERE config.review_master_sms_gateway_pair_code = ?" +
		" AND config.review_master_sms_gateway_enabled = 1" +
		" AND config.enabled = 1" +
		" AND client.enabled = 1"
	rows, err := Db.Query(qry, reviewMasterSmsGatewayPairCode)
	var id int
	var useMasterQueue bool
	if err != nil {
		log.Println("Error retrieving Review Master SMS Gateway Pair Code", reviewMasterSmsGatewayPairCode, "from database. Error: ", err)
		return 0
	}
	defer rows.Close()
	for rows.Next() {
		if err1 := rows.Scan(&id, &useMasterQueue); err1 != nil {
			log.Println("Error retrieving Review Master SMS Gateway Pair Code", reviewMasterSmsGatewayPairCode, "from database whilst reading returned results. Error: ", err1)
			return 0
		}
		if id == 0 {
			log.Printf("Review Master SMS Gateway Pair Code %s not found", reviewMasterSmsGatewayPairCode)
			continue
		}
	}
	if useMasterQueue {
		return ReviewMasterSMSMasterQueue
	}
	return id
}

// SetReviewMasterSMSGatewayMasterQueueID - set the Review Master SMS Gateway Master Queue ID
func SetReviewMasterSMSGatewayMasterQueueID() {
	qry := "SELECT MAX(id) FROM review_master_sms_gateway_master_queues"
	var id int
	row := Db.QueryRow(qry)
	if err := row.Scan(&id); err != nil {
		log.Println("Error retrieving Review Master SMS Gateway Master Queue ID from database, error: ", err)
	}
	ReviewMasterSMSMasterQueue = id
}

// AddSendLater - add send later for messages that are delayed (send later)
func AddSendLater(telephone string, clientID uint64, sendAfterMinutes int,
	sendURL string, method string, appKey string, secretKey string, headers map[string]string,
	params url.Values, body []byte, sendFromIcabbiApp bool, reviewMasterSMSGatewayEnabled bool,
	alternateMessageServiceEnabled bool, alternateMessageService string,
	sendFromOwnSMSGatewayEnabled bool, sendSuccessResponse string, maxDailySendCount uint) {

	// serialize headers
	h := new(bytes.Buffer)
	he := gob.NewEncoder(h)
	err := he.Encode(headers)
	if err != nil {
		h = nil
	}
	httpHeaders := h.Bytes()

	// serialize params
	httpParams := params.Encode()

	qry := "INSERT INTO google_reviews_send_laters" +
		" (telephone, send_after, send_url, http_method, app_key, secret_key," +
		" http_headers, http_params, http_body, send_from_icabbi_app," +
		" review_master_sms_gateway_enabled, alternate_message_service_enabled," +
		" alternate_message_service, send_from_own_sms_gateway_enabled," +
		" send_success_response, max_daily_send_count, client_id)" +
		" VALUES (?, DATE_ADD(NOW(), INTERVAL ? MINUTE), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)" +
		" ON DUPLICATE KEY UPDATE" +
		" send_after = DATE_ADD(NOW(), INTERVAL ? MINUTE)," +
		" send_url = ?," +
		" http_method = ?," +
		" app_key = ?," +
		" secret_key = ?," +
		" http_headers = ?," +
		" http_params = ?," +
		" http_body = ?," +
		" send_from_icabbi_app = ?," +
		" review_master_sms_gateway_enabled = ?," +
		" alternate_message_service_enabled = ?," +
		" alternate_message_service = ?," +
		" send_from_own_sms_gateway_enabled = ?," +
		" send_success_response	= ?," +
		" max_daily_send_count = ?"
	_, err = Db.Exec(qry, telephone, sendAfterMinutes, sendURL, method, appKey, secretKey,
		httpHeaders, httpParams, body, sendFromIcabbiApp, reviewMasterSMSGatewayEnabled,
		alternateMessageServiceEnabled, alternateMessageService,
		sendFromOwnSMSGatewayEnabled, sendSuccessResponse, maxDailySendCount, clientID,
		sendAfterMinutes, sendURL, method, appKey, secretKey, httpHeaders, httpParams,
		body, sendFromIcabbiApp, reviewMasterSMSGatewayEnabled,
		alternateMessageServiceEnabled, alternateMessageService,
		sendFromOwnSMSGatewayEnabled, sendSuccessResponse, maxDailySendCount)
	if err != nil {
		log.Println(err)
	}
}

// UpdateStatsCanUseToken - update the stats
// set the clientID to 0 if not known and send the token to try and retrieve the clientID
func UpdateStatsCanUseToken(clientID uint64, token string, sent bool) {
	// if clientID is not set then try and retrieve it from the token (this will be the case when the token checks fail)
	cID := clientID
	if cID == 0 {
		qry := "SELECT client.id" +
			" FROM google_reviews_configs AS config" +
			" JOIN clients AS client ON client.id = config.client_id" +
			" WHERE config.token = ?" +
			" AND config.enabled = 1" +
			" AND client.enabled = 1"
		rows, err := Db.Query(qry, token)
		if err != nil {
			log.Println("Error retrieving token", token, "from database, for stats. Error: ", err)
			return
		}
		defer rows.Close()
		for rows.Next() {
			if err1 := rows.Scan(&cID); err1 != nil {
				log.Println("Error retrieving token", token, "from database whilst reading returned results, for stats. Error: ", err1)
				return
			}
			if cID == 0 {
				log.Printf("token %s not found, for stats", token)
				continue
			}
			break
		}
	}
	if cID == 0 {
		return
	}

	// sent
	qry := "INSERT INTO google_reviews_stats" +
		" (client_id, stats_date, sent_count, requested_count)" +
		" VALUES (?, CURDATE(), 1, 1)" +
		" ON DUPLICATE KEY UPDATE" +
		" sent_count = sent_count + 1," +
		" requested_count = requested_count + 1"
	if !sent {
		// not sent
		qry = "INSERT INTO google_reviews_stats" +
			" (client_id, stats_date, sent_count, requested_count)" +
			" VALUES (?, CURDATE(), 0, 1)" +
			" ON DUPLICATE KEY UPDATE" +
			" requested_count = requested_count + 1"
	}
	_, err := Db.Exec(qry, cID)
	if err != nil {
		log.Println(err)
	}
}
