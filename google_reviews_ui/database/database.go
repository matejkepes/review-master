package database

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"strings"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

// Client - represents a client.
type Client struct {
	ID        uint64 `json:"id"`         // id
	Enabled   bool   `json:"enabled"`    // enabled
	Name      string `json:"name"`       // name
	Note      string `json:"note"`       // note
	Country   string `json:"country"`    // country
	PartnerID uint64 `json:"partner_id"` // partner id
}

// GoogleReviewsConfig - represents a google reviews config
type GoogleReviewsConfig struct {
	ID                                            uint64 `json:"id"`                                                     // id
	Enabled                                       bool   `json:"enabled"`                                                // enabled
	MinSendFrequency                              uint   `json:"min_send_frequency,string,omitempty"`                    // min send frequency
	MaxSendCount                                  uint   `json:"max_send_count,string,omitempty"`                        // max send count
	MaxDailySendCount                             uint   `json:"max_daily_send_count,string,omitempty"`                  // max daily send count
	Token                                         string `json:"token"`                                                  // token
	TelephoneParameter                            string `json:"telephone_parameter"`                                    // telephone parameter
	SendFromIcabbiApp                             bool   `json:"send_from_icabbi_app"`                                   // send_from_icabbi_app
	AppKey                                        string `json:"app_key"`                                                // app_key url
	SecretKey                                     string `json:"secret_key"`                                             // secret_key url
	SendURL                                       string `json:"send_url"`                                               // send url
	HttpGet                                       bool   `json:"http_get"`                                               // http get
	SendSuccessResponse                           string `json:"send_success_response"`                                  // send success response
	TimeZone                                      string `json:"time_zone"`                                              // time_zone
	MultiMessageEnabled                           bool   `json:"multi_message_enabled"`                                  // multi message enabled
	MessageParameter                              string `json:"message_parameter"`                                      // message parameter
	MultiMessageSeparator                         string `json:"multi_message_separator"`                                // multi message separator
	UseDatabaseMessage                            bool   `json:"use_database_message"`                                   // use database message
	Message                                       string `json:"message"`                                                // message
	SendDelayEnabled                              bool   `json:"send_delay_enabled"`                                     // send delay enabled
	SendDelay                                     uint   `json:"send_delay,string,omitempty"`                            // send delay
	DispatcherChecksEnabled                       bool   `json:"dispatcher_checks_enabled"`                              // dispatcher checks enabled
	DispatcherType                                string `json:"dispatcher_type"`                                        // dispatcher type
	DispatcherURL                                 string `json:"dispatcher_url"`                                         // dispatcher url
	BookingIdParameter                            string `json:"booking_id_parameter"`                                   // booking id parameter
	IsBookingForNowDiffMinutes                    uint   `json:"is_booking_for_now_diff_minutes,string,omitempty"`       // is booking for now diff minutes
	BookingNowPickupToContactMinutes              uint   `json:"booking_now_pickup_to_contact_minutes,string,omitempty"` // booking now pickup to contact minutes
	PreBookingPickupToContactMinutes              uint   `json:"pre_booking_pickup_to_contact_minutes,string,omitempty"` // pre booking pickup to contact minutes
	ReplaceTelephoneCountryCode                   bool   `json:"replace_telephone_country_code"`                         // replace telephone country code
	ReplaceTelephoneCountryCodeWith               string `json:"replace_telephone_country_code_with"`                    // replace telephone country code with
	ReviewMasterSMSGatewayEnabled                 bool   `json:"review_master_sms_gateway_enabled"`                      // review master sms gateway enabled
	ReviewMasterSMSGatewayUseMasterQueue          bool   `json:"review_master_sms_gateway_use_master_queue"`             // review master sms gateway use master queue
	ReviewMasterSMSGatewayPairCode                string `json:"review_master_sms_gateway_pair_code"`                    // review master sms gateway pair code
	AlternateMessageServiceEnabled                bool   `json:"alternate_message_service_enabled"`                      // alternate message service enabled
	AlternateMessageService                       string `json:"alternate_message_service"`                              // alternate message service
	AlternateMessageServiceSecret1                string `json:"alternate_message_service_secret1"`                      // alternate message service
	Companies                                     string `json:"companies"`                                              // companies
	BookingSourceMobileAppState                   int    `json:"booking_source_mobile_app_state"`                        // booking source mobile app state
	AIResponsesEnabled                            bool   `json:"ai_responses_enabled"`                                   // AI responses enabled
	MonthlyReviewAnalysisEnabled                  bool   `json:"monthly_review_analysis_enabled"`                        // Monthly review analysis enabled
	ContactMethod                                 string `json:"contact_method"`                                         // Contact method
	GoogleMyBusinessReviewReplyEnabled            bool   `json:"google_my_business_review_reply_enabled"`                // Google My Business google my business review reply enabled
	GoogleMyBusinessLocationName                  string `json:"google_my_business_location_name"`                       // Google My Business google my business location name
	GoogleMyBusinessPostalCode                    string `json:"google_my_business_postal_code"`                         // Google My Business google my business postal code address
	GoogleMyBusinessReplyToUnspecfifiedStarRating bool   `json:"google_my_business_reply_to_unspecfified_star_rating"`   // Google My Business reply to unspecified star rating enabled
	GoogleMyBusinessUnspecfifiedStarRatingReply   string `json:"google_my_business_unspecfified_star_rating_reply"`      // Google My Business unspecified star rating reply
	GoogleMyBusinessReplyToOneStarRating          bool   `json:"google_my_business_reply_to_one_star_rating"`            // Google My Business reply to one star rating enabled
	GoogleMyBusinessOneStarRatingReply            string `json:"google_my_business_one_star_rating_reply"`               // Google My Business one star rating reply
	GoogleMyBusinessReplyToTwoStarRating          bool   `json:"google_my_business_reply_to_two_star_rating"`            // Google My Business reply to two star rating enabled
	GoogleMyBusinessTwoStarRatingReply            string `json:"google_my_business_two_star_rating_reply"`               // Google My Business two star rating reply
	GoogleMyBusinessReplyToThreeStarRating        bool   `json:"google_my_business_reply_to_three_star_rating"`          // Google My Business reply to three star rating enabled
	GoogleMyBusinessThreeStarRatingReply          string `json:"google_my_business_three_star_rating_reply"`             // Google My Business three star rating reply
	GoogleMyBusinessReplyToFourStarRating         bool   `json:"google_my_business_reply_to_four_star_rating"`           // Google My Business reply to four star rating enabled
	GoogleMyBusinessFourStarRatingReply           string `json:"google_my_business_four_star_rating_reply"`              // Google My Business four star rating reply
	GoogleMyBusinessReplyToFiveStarRating         bool   `json:"google_my_business_reply_to_five_star_rating"`           // Google My Business reply to five star rating enabled
	GoogleMyBusinessFiveStarRatingReply           string `json:"google_my_business_five_star_rating_reply"`              // Google My Business five star rating reply
	GoogleMyBusinessReportEnabled                 bool   `json:"google_my_business_report_enabled"`                      // Google My Business report enabled
	EmailAddress                                  string `json:"email_address"`                                          // Email address used for reporting
	ClientID                                      uint64 `json:"client_id"`                                              // client id
}

// GoogleReviewsConfigTime - represents a google reviews config time
type GoogleReviewsConfigTime struct {
	ID        uint64 `json:"id"`        // id
	Enabled   bool   `json:"enabled"`   // enabled
	Start     string `json:"start"`     // start
	End       string `json:"end"`       // end
	Sunday    bool   `json:"sunday"`    // sunday
	Monday    bool   `json:"monday"`    // monday
	Tuesday   bool   `json:"tuesday"`   // tuesday
	Wednesday bool   `json:"wednesday"` // wednesday
	Thursday  bool   `json:"thursday"`  // thursday
	Friday    bool   `json:"friday"`    // friday
	Saturday  bool   `json:"saturday"`  // saturday
	// GoogleReviewsConfigID uint64 `json:"google_reviews_config_id"` // google reviews config id
	GoogleReviewsConfigID uint64 `json:"google_reviews_config_id,string,omitempty"` // google reviews config id
}

// SimpleConfig - 1 client, 1 config, 1 time
type SimpleConfig struct {
	ClientID                                                uint64 `json:"client_id"`                                                                    // client id
	ClientEnabled                                           bool   `json:"client_enabled"`                                                               // client enabled
	ClientName                                              string `json:"client_name"`                                                                  // client name
	ClientNote                                              string `json:"client_note"`                                                                  // client note
	ClientCountry                                           string `json:"client_country"`                                                               // client country
	GoogleReviewsConfigID                                   uint64 `json:"google_reviews_config_id"`                                                     // google reviews config id
	GoogleReviewsConfigEnabled                              bool   `json:"google_reviews_config_enabled"`                                                // google reviews config enabled
	GoogleReviewsConfigMinSendFrequency                     uint   `json:"google_reviews_config_min_send_frequency,string,omitempty"`                    // google reviews config min send frequency
	GoogleReviewsConfigMaxSendCount                         uint   `json:"google_reviews_config_max_send_count,string,omitempty"`                        // google reviews config max send count
	GoogleReviewsConfigMaxDailySendCount                    uint   `json:"google_reviews_config_max_daily_send_count,string,omitempty"`                  // google reviews config max daily send count
	GoogleReviewsConfigToken                                string `json:"google_reviews_config_token"`                                                  // google reviews config token
	GoogleReviewsConfigTelephoneParameter                   string `json:"google_reviews_config_telephone_parameter"`                                    // google reviews config telephone parameter
	GoogleReviewsConfigSendFromIcabbiApp                    bool   `json:"google_reviews_config_send_from_icabbi_app"`                                   // google_reviews_config_send_from_icabbi_app
	GoogleReviewsConfigAppKey                               string `json:"google_reviews_config_app_key"`                                                // google_reviews_config_app_key url
	GoogleReviewsConfigSecretKey                            string `json:"google_reviews_config_secret_key"`                                             // google_reviews_config_secret_key url
	GoogleReviewsConfigSendURL                              string `json:"google_reviews_config_send_url"`                                               // google reviews config send url
	GoogleReviewsConfigHttpGet                              bool   `json:"google_reviews_config_http_get"`                                               // google reviews config http get
	GoogleReviewsConfigSendSuccessResponse                  string `json:"google_reviews_config_send_success_response"`                                  // google reviews config send success response
	GoogleReviewsConfigTimeZone                             string `json:"google_reviews_config_time_zone"`                                              // google reviews config time zone
	GoogleReviewsConfigMultiMessageEnabled                  bool   `json:"google_reviews_config_multi_message_enabled"`                                  // google reviews config multi message enabled
	GoogleReviewsConfigMessageParameter                     string `json:"google_reviews_config_message_parameter"`                                      // google reviews config message parameter
	GoogleReviewsConfigMultiMessageSeparator                string `json:"google_reviews_config_multi_message_separator"`                                // google reviews config multi message separator
	GoogleReviewsConfigUseDatabaseMessage                   bool   `json:"google_reviews_config_use_database_message"`                                   // google reviews config use database message
	GoogleReviewsConfigMessage                              string `json:"google_reviews_config_message"`                                                // google reviews config message
	GoogleReviewsConfigSendDelayEnabled                     bool   `json:"google_reviews_config_send_delay_enabled"`                                     // google reviews config send delay enabled
	GoogleReviewsConfigSendDelay                            uint   `json:"google_reviews_config_send_delay,string,omitempty"`                            // google reviews config send delay
	GoogleReviewsConfigDispatcherChecksEnabled              bool   `json:"google_reviews_config_dispatcher_checks_enabled"`                              // google reviews config dispatcher checks enabled
	GoogleReviewsConfigDispatcherType                       string `json:"google_reviews_config_dispatcher_type"`                                        // google reviews config dispatcher type
	GoogleReviewsConfigDispatcherURL                        string `json:"google_reviews_config_dispatcher_url"`                                         // google reviews config dispatcher url
	GoogleReviewsConfigBookingIdParameter                   string `json:"google_reviews_config_booking_id_parameter"`                                   // google reviews config booking id parameter
	GoogleReviewsConfigIsBookingForNowDiffMinutes           uint   `json:"google_reviews_config_is_booking_for_now_diff_minutes,string,omitempty"`       // google reviews config is booking for now diff minutes
	GoogleReviewsConfigBookingNowPickupToContactMinutes     uint   `json:"google_reviews_config_booking_now_pickup_to_contact_minutes,string,omitempty"` // google reviews config booking now pickup to contact minutes
	GoogleReviewsConfigPreBookingPickupToContactMinutes     uint   `json:"google_reviews_config_pre_booking_pickup_to_contact_minutes,string,omitempty"` // google reviews config pre booking pickup to contact minutes
	GoogleReviewsConfigReplaceTelephoneCountryCode          bool   `json:"google_reviews_config_replace_telephone_country_code"`                         // google reviews config replace telephone country code
	GoogleReviewsConfigReplaceTelephoneCountryCodeWith      string `json:"google_reviews_config_replace_telephone_country_code_with"`                    // google reviews config replace telephone country code with
	GoogleReviewsConfigReviewMasterSMSGatewayEnabled        bool   `json:"google_reviews_config_review_master_sms_gateway_enabled"`                      // google reviews config review master sms gateway enabled
	GoogleReviewsConfigReviewMasterSMSGatewayUseMasterQueue bool   `json:"google_reviews_config_review_master_sms_gateway_use_master_queue"`             // oogle reviews config review master sms gateway use master queue
	GoogleReviewsConfigReviewMasterSMSGatewayPairCode       string `json:"google_reviews_config_review_master_sms_gateway_pair_code"`                    // google reviews config review master sms gateway pair code
	GoogleReviewsConfigAlternateMessageServiceEnabled       bool   `json:"google_reviews_config_alternate_message_service_enabled"`                      // google reviews config alternate message service enabled
	GoogleReviewsConfigAlternateMessageService              string `json:"google_reviews_config_alternate_message_service"`                              // google reviews config alternate message service
	GoogleReviewsConfigAlternateMessageServiceSecret1       string `json:"google_reviews_config_alternate_message_service_secret1"`                      // google reviews config alternate message service secret1
	GoogleReviewsConfigCompanies                            string `json:"google_reviews_config_companies"`                                              // google reviews config review companies
	GoogleReviewsConfigBookingSourceMobileAppState          int    `json:"google_reviews_config_booking_source_mobile_app_state"`                        // google reviews config review booking source mobile app state
	GoogleReviewsConfigAIResponsesEnabled                   bool   `json:"google_reviews_config_ai_responses_enabled"`                                   // google reviews config AI responses enabled
	GoogleReviewsConfigMonthlyReviewAnalysisEnabled         bool   `json:"google_reviews_config_monthly_review_analysis_enabled"`                        // google reviews config monthly review analysis enabled
	GoogleReviewsConfigContactMethod                        string `json:"google_reviews_config_contact_method"`                                         // google reviews config contact method
	GoogleMyBusinessReviewReplyEnabled                      bool   `json:"google_my_business_review_reply_enabled"`                                      // Google My Business google my business review reply enabled
	GoogleMyBusinessLocationName                            string `json:"google_my_business_location_name"`                                             // Google My Business google my business location name
	GoogleMyBusinessPostalCode                              string `json:"google_my_business_postal_code"`                                               // Google My Business google my business postal code address
	GoogleMyBusinessReplyToUnspecfifiedStarRating           bool   `json:"google_my_business_reply_to_unspecfified_star_rating"`                         // Google My Business reply to unspecified star rating enabled
	GoogleMyBusinessUnspecfifiedStarRatingReply             string `json:"google_my_business_unspecfified_star_rating_reply"`                            // Google My Business unspecified star rating reply
	GoogleMyBusinessReplyToOneStarRating                    bool   `json:"google_my_business_reply_to_one_star_rating"`                                  // Google My Business reply to one star rating enabled
	GoogleMyBusinessOneStarRatingReply                      string `json:"google_my_business_one_star_rating_reply"`                                     // Google My Business one star rating reply
	GoogleMyBusinessReplyToTwoStarRating                    bool   `json:"google_my_business_reply_to_two_star_rating"`                                  // Google My Business reply to two star rating enabled
	GoogleMyBusinessTwoStarRatingReply                      string `json:"google_my_business_two_star_rating_reply"`                                     // Google My Business two star rating reply
	GoogleMyBusinessReplyToThreeStarRating                  bool   `json:"google_my_business_reply_to_three_star_rating"`                                // Google My Business reply to three star rating enabled
	GoogleMyBusinessThreeStarRatingReply                    string `json:"google_my_business_three_star_rating_reply"`                                   // Google My Business three star rating reply
	GoogleMyBusinessReplyToFourStarRating                   bool   `json:"google_my_business_reply_to_four_star_rating"`                                 // Google My Business reply to four star rating enabled
	GoogleMyBusinessFourStarRatingReply                     string `json:"google_my_business_four_star_rating_reply"`                                    // Google My Business four star rating reply
	GoogleMyBusinessReplyToFiveStarRating                   bool   `json:"google_my_business_reply_to_five_star_rating"`                                 // Google My Business reply to five star rating enabled
	GoogleMyBusinessFiveStarRatingReply                     string `json:"google_my_business_five_star_rating_reply"`                                    // Google My Business five star rating reply
	GoogleMyBusinessReportEnabled                           bool   `json:"google_my_business_report_enabled"`                                            // Google My Business report enabled
	EmailAddress                                            string `json:"email_address"`                                                                // Email address used for reporting
	GoogleReviewsConfigTimeID                               uint64 `json:"google_reviews_config_time_id"`                                                // google reviews config time id
	GoogleReviewsConfigTimeEnabled                          bool   `json:"google_reviews_config_time_enabled"`                                           // google reviews config time enabled
	GoogleReviewsConfigTimeStart                            string `json:"google_reviews_config_time_start"`                                             // google reviews config time start
	GoogleReviewsConfigTimeEnd                              string `json:"google_reviews_config_time_end"`                                               // google reviews config time end
	GoogleReviewsConfigTimeSunday                           bool   `json:"google_reviews_config_time_sunday"`
	GoogleReviewsConfigTimeMonday                           bool   `json:"google_reviews_config_time_monday"`
	GoogleReviewsConfigTimeTuesday                          bool   `json:"google_reviews_config_time_tuesday"`
	GoogleReviewsConfigTimeWednesday                        bool   `json:"google_reviews_config_time_wednesday"`
	GoogleReviewsConfigTimeThursday                         bool   `json:"google_reviews_config_time_thursday"`
	GoogleReviewsConfigTimeFriday                           bool   `json:"google_reviews_config_time_friday"`
	GoogleReviewsConfigTimeSaturday                         bool   `json:"google_reviews_config_time_saturday"`
}

// Config - 1 config, many times
type Config struct {
	GoogleReviewsConfig      GoogleReviewsConfig       `json:"google_reviews_config"`
	GoogleReviewsConfigTimes []GoogleReviewsConfigTime `json:"google_reviews_config_times"`
}

// ClientConfig - 1 client, many configs with many times
type ClientConfig struct {
	Client  Client   `json:"client"`
	Configs []Config `json:"configs"`
}

// StatsRequest - represents a statistic request.
type StatsRequest struct {
	StartDay     string `json:"start_day"`     // start day
	EndDay       string `json:"end_day"`       // end day
	TimeGrouping string `json:"time_grouping"` // time grouping
}

// StatsResult - represents a statistic result.
type StatsResult struct {
	ClientID    uint64 `json:"client_id"`    // client id
	ClientName  string `json:"client_name"`  // client name
	Sent        uint64 `json:"sent"`         // sent
	GroupPeriod string `json:"group_period"` // group period
}

// StatsNewResult - represents a statistic (new) result.
type StatsNewResult struct {
	ClientID    uint64 `json:"client_id"`    // client id
	ClientName  string `json:"client_name"`  // client name
	Sent        uint64 `json:"sent"`         // sent
	Requested   uint64 `json:"requested"`    // requested
	GroupPeriod string `json:"group_period"` // group period
}

// NothingSentResult - represents a check for no messages sent result.
type NothingSentResult struct {
	ClientID   uint64 `json:"client_id"`   // client id
	ClientName string `json:"client_name"` // client name
}

// User - represents a user of the frontend
type User struct {
	ID       uint64 `json:"id"`       // id
	Email    string `json:"email"`    // email
	Password string `json:"password"` // password
}

// UserClients - represents a user and list of clients for display for the frontend
type UserClients struct {
	User    User     `json:"user"`    // user
	Clients []Client `json:"clients"` // clients
}

// UserClientList - represents a user and clients list for display for the frontend
type UserClientList struct {
	UserID     uint64 `json:"user_id"` // user id
	Email      string `json:"email"`   // email
	ClientList string `json:"clients"` // clients
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
	// NOTE: The above about concurrent calls is not correct for this program but restrict anyway.
	//
	// Set the maximum number of concurrently open connections (in-use + idle)
	// to 10. Setting this to less than or equal to 0 will mean there is no
	// maximum limit (which is also the default setting).
	db.SetMaxOpenConns(10)

	// Test connection to the database (this will result in the program ending if unable to establish a connection)
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}
	Db = db
}

// ListAllClients - get all the clients for a specific partner
func ListAllClients(partnerID int) ([]Client, error) {
	var err error

	const qry = "SELECT client.id, client.enabled, client.name, client.country" +
		" FROM clients AS client" +
		" WHERE partner_id = ?"
	rows, err := Db.Query(qry, partnerID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	var clients []Client
	for rows.Next() {
		var c Client
		if err := rows.Scan(&c.ID, &c.Enabled, &c.Name, &c.Country); err != nil {
			log.Printf("Error getting all clients: %v\n", err)
		}
		clients = append(clients, c)
	}
	return clients, nil
}

// GetSimpleClient - get a client, simple assume only 1 config and time
func GetSimpleClient(clientID int, partnerID int) (SimpleConfig, error) {
	// var err error
	var s SimpleConfig

	const qry = "SELECT client.id, client.enabled, client.name, client.note, client.country," +
		" config.id, config.enabled, config.min_send_frequency, config.max_send_count," +
		" config.max_daily_send_count, config.token, config.telephone_parameter," +
		" config.send_from_icabbi_app, config.app_key, config.secret_key," +
		" config.send_url, config.http_get, config.send_success_response, config.time_zone," +
		" config.multi_message_enabled, config.message_parameter, config.multi_message_separator," +
		" config.use_database_message, config.message," +
		" config.send_delay_enabled, config.send_delay," +
		" config.dispatcher_checks_enabled, config.dispatcher_type, config.dispatcher_url," +
		" config.booking_id_parameter, config.is_booking_for_now_diff_minutes," +
		" config.booking_now_pickup_to_contact_minutes, config.pre_booking_pickup_to_contact_minutes," +
		" config.replace_telephone_country_code, config.replace_telephone_country_code_with," +
		" config.review_master_sms_gateway_enabled," +
		" config.review_master_sms_gateway_use_master_queue," +
		" config.review_master_sms_gateway_pair_code," +
		" config.alternate_message_service_enabled, config.alternate_message_service, config.alternate_message_service_secret1," +
		" config.companies, config.booking_source_mobile_app_state," +
		" IFNULL(config.ai_responses_enabled, 0), IFNULL(config.contact_method, '')," +
		" IFNULL(config.monthly_review_analysis_enabled, 0)," +
		" config.google_my_business_review_reply_enabled," +
		" config.google_my_business_location_name," +
		" config.google_my_business_postal_code," +
		" config.google_my_business_reply_to_unspecfified_star_rating," +
		" config.google_my_business_unspecfified_star_rating_reply," +
		" config.google_my_business_reply_to_one_star_rating," +
		" config.google_my_business_one_star_rating_reply," +
		" config.google_my_business_reply_to_two_star_rating," +
		" config.google_my_business_two_star_rating_reply," +
		" config.google_my_business_reply_to_three_star_rating," +
		" config.google_my_business_three_star_rating_reply," +
		" config.google_my_business_reply_to_four_star_rating," +
		" config.google_my_business_four_star_rating_reply," +
		" config.google_my_business_reply_to_five_star_rating," +
		" config.google_my_business_five_star_rating_reply," +
		" config.google_my_business_report_enabled," +
		" config.email_address," +
		" times.id, times.enabled, times.start, times.end," +
		" times.sunday, times.monday, times.tuesday, times.wednesday, times.thursday, times.friday, times.saturday" +
		" FROM google_reviews_config_times AS times" +
		" JOIN google_reviews_configs AS config ON config.id = times.google_reviews_config_id" +
		" JOIN clients AS client ON client.id = config.client_id" +
		" WHERE client.id = ? AND partner_id = ?"
	rows, err := Db.Query(qry, clientID, partnerID)
	if err != nil {
		log.Println(err)
		return s, err
	}
	defer rows.Close()
	// assume only one (therefore retrieving the first one)
	if rows.Next() {
		if err := rows.Scan(&s.ClientID, &s.ClientEnabled, &s.ClientName, &s.ClientNote, &s.ClientCountry,
			&s.GoogleReviewsConfigID, &s.GoogleReviewsConfigEnabled, &s.GoogleReviewsConfigMinSendFrequency, &s.GoogleReviewsConfigMaxSendCount,
			&s.GoogleReviewsConfigMaxDailySendCount, &s.GoogleReviewsConfigToken, &s.GoogleReviewsConfigTelephoneParameter,
			&s.GoogleReviewsConfigSendFromIcabbiApp, &s.GoogleReviewsConfigAppKey, &s.GoogleReviewsConfigSecretKey,
			&s.GoogleReviewsConfigSendURL, &s.GoogleReviewsConfigHttpGet, &s.GoogleReviewsConfigSendSuccessResponse, &s.GoogleReviewsConfigTimeZone,
			&s.GoogleReviewsConfigMultiMessageEnabled, &s.GoogleReviewsConfigMessageParameter, &s.GoogleReviewsConfigMultiMessageSeparator,
			&s.GoogleReviewsConfigUseDatabaseMessage, &s.GoogleReviewsConfigMessage,
			&s.GoogleReviewsConfigSendDelayEnabled, &s.GoogleReviewsConfigSendDelay,
			&s.GoogleReviewsConfigDispatcherChecksEnabled, &s.GoogleReviewsConfigDispatcherType, &s.GoogleReviewsConfigDispatcherURL,
			&s.GoogleReviewsConfigBookingIdParameter, &s.GoogleReviewsConfigIsBookingForNowDiffMinutes,
			&s.GoogleReviewsConfigBookingNowPickupToContactMinutes, &s.GoogleReviewsConfigPreBookingPickupToContactMinutes,
			&s.GoogleReviewsConfigReplaceTelephoneCountryCode, &s.GoogleReviewsConfigReplaceTelephoneCountryCodeWith,
			&s.GoogleReviewsConfigReviewMasterSMSGatewayEnabled,
			&s.GoogleReviewsConfigReviewMasterSMSGatewayUseMasterQueue,
			&s.GoogleReviewsConfigReviewMasterSMSGatewayPairCode,
			&s.GoogleReviewsConfigAlternateMessageServiceEnabled, &s.GoogleReviewsConfigAlternateMessageService, &s.GoogleReviewsConfigAlternateMessageServiceSecret1,
			&s.GoogleReviewsConfigCompanies, &s.GoogleReviewsConfigBookingSourceMobileAppState,
			&s.GoogleReviewsConfigAIResponsesEnabled, &s.GoogleReviewsConfigContactMethod,
			&s.GoogleReviewsConfigMonthlyReviewAnalysisEnabled,
			&s.GoogleMyBusinessReviewReplyEnabled,
			&s.GoogleMyBusinessLocationName,
			&s.GoogleMyBusinessPostalCode,
			&s.GoogleMyBusinessReplyToUnspecfifiedStarRating,
			&s.GoogleMyBusinessUnspecfifiedStarRatingReply,
			&s.GoogleMyBusinessReplyToOneStarRating,
			&s.GoogleMyBusinessOneStarRatingReply,
			&s.GoogleMyBusinessReplyToTwoStarRating,
			&s.GoogleMyBusinessTwoStarRatingReply,
			&s.GoogleMyBusinessReplyToThreeStarRating,
			&s.GoogleMyBusinessThreeStarRatingReply,
			&s.GoogleMyBusinessReplyToFourStarRating,
			&s.GoogleMyBusinessFourStarRatingReply,
			&s.GoogleMyBusinessReplyToFiveStarRating,
			&s.GoogleMyBusinessFiveStarRatingReply,
			&s.GoogleMyBusinessReportEnabled,
			&s.EmailAddress,
			&s.GoogleReviewsConfigTimeID, &s.GoogleReviewsConfigTimeEnabled, &s.GoogleReviewsConfigTimeStart, &s.GoogleReviewsConfigTimeEnd,
			&s.GoogleReviewsConfigTimeSunday, &s.GoogleReviewsConfigTimeMonday, &s.GoogleReviewsConfigTimeTuesday, &s.GoogleReviewsConfigTimeWednesday, &s.GoogleReviewsConfigTimeThursday, &s.GoogleReviewsConfigTimeFriday, &s.GoogleReviewsConfigTimeSaturday); err != nil {
			log.Printf("Error getting configs for client: %v\n", err)
		}
	}
	return s, nil
}

// UpdateSimpleClient - update a client and config and time, simple assume only 1 config and time
func UpdateSimpleClient(simpleConfig SimpleConfig) error {
	const clientQry = "UPDATE clients SET enabled = ?, name = ?, note = ?, country = ? WHERE id = ?"
	const googleReviewsConfigQry = "UPDATE google_reviews_configs SET enabled = ?," +
		" min_send_frequency = ?, max_send_count = ?," +
		" max_daily_send_count = ?, token = ?, telephone_parameter = ?," +
		" send_from_icabbi_app = ?, app_key = ?, secret_key = ?," +
		" send_url = ?, http_get = ?, send_success_response = ?, time_zone = ?," +
		" multi_message_enabled = ?, message_parameter = ?, multi_message_separator = ?," +
		" use_database_message = ?, message = ?," +
		" send_delay_enabled = ?, send_delay = ?," +
		" dispatcher_checks_enabled = ?, dispatcher_type = ?, dispatcher_url = ?," +
		" booking_id_parameter = ?, is_booking_for_now_diff_minutes = ?," +
		" booking_now_pickup_to_contact_minutes = ?, pre_booking_pickup_to_contact_minutes = ?," +
		" replace_telephone_country_code = ?, replace_telephone_country_code_with = ?," +
		" review_master_sms_gateway_enabled = ?," +
		" review_master_sms_gateway_use_master_queue = ?," +
		" review_master_sms_gateway_pair_code = ?," +
		" alternate_message_service_enabled = ?, alternate_message_service = ?, alternate_message_service_secret1 = ?," +
		" companies = ?, booking_source_mobile_app_state = ?," +
		" ai_responses_enabled = ?," +
		" contact_method = NULLIF(?, '')," + // Use NULLIF to convert empty string to NULL
		" monthly_review_analysis_enabled = ?," + // Use NULLIF to convert empty string to NULL
		" google_my_business_review_reply_enabled = ?," +
		" google_my_business_location_name = ?," +
		" google_my_business_postal_code = ?," +
		" google_my_business_reply_to_unspecfified_star_rating = ?," +
		" google_my_business_unspecfified_star_rating_reply = ?," +
		" google_my_business_reply_to_one_star_rating = ?," +
		" google_my_business_one_star_rating_reply = ?," +
		" google_my_business_reply_to_two_star_rating = ?," +
		" google_my_business_two_star_rating_reply = ?," +
		" google_my_business_reply_to_three_star_rating = ?," +
		" google_my_business_three_star_rating_reply = ?," +
		" google_my_business_reply_to_four_star_rating = ?," +
		" google_my_business_four_star_rating_reply = ?," +
		" google_my_business_reply_to_five_star_rating = ?," +
		" google_my_business_five_star_rating_reply = ?," +
		" google_my_business_report_enabled = ?," +
		" email_address = ?" +
		" WHERE id = ?"
	const googleReviewsConfigTimeQry = "UPDATE google_reviews_config_times SET enabled = ?," +
		" start = ?, end = ?," +
		" sunday = ?, monday = ?, tuesday = ?, wednesday = ?, thursday = ?, friday = ?, saturday = ?" +
		" WHERE id = ?"

	tx, err := Db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}

	_, execErr := tx.Exec(clientQry, simpleConfig.ClientEnabled, strings.TrimSpace(simpleConfig.ClientName),
		strings.TrimSpace(simpleConfig.ClientNote), strings.TrimSpace(simpleConfig.ClientCountry), simpleConfig.ClientID)
	if execErr != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("update failed: %v, unable to rollback: %v\n", execErr, rollbackErr)
			return execErr
		}
		log.Printf("update failed: %v", execErr)
		return execErr
	}

	_, execErr = tx.Exec(googleReviewsConfigQry, simpleConfig.GoogleReviewsConfigEnabled,
		simpleConfig.GoogleReviewsConfigMinSendFrequency, simpleConfig.GoogleReviewsConfigMaxSendCount,
		simpleConfig.GoogleReviewsConfigMaxDailySendCount, strings.TrimSpace(simpleConfig.GoogleReviewsConfigToken),
		strings.TrimSpace(simpleConfig.GoogleReviewsConfigTelephoneParameter), simpleConfig.GoogleReviewsConfigSendFromIcabbiApp,
		strings.TrimSpace(simpleConfig.GoogleReviewsConfigAppKey), strings.TrimSpace(simpleConfig.GoogleReviewsConfigSecretKey),
		strings.TrimSpace(simpleConfig.GoogleReviewsConfigSendURL), simpleConfig.GoogleReviewsConfigHttpGet,
		strings.TrimSpace(simpleConfig.GoogleReviewsConfigSendSuccessResponse), strings.TrimSpace(simpleConfig.GoogleReviewsConfigTimeZone),
		simpleConfig.GoogleReviewsConfigMultiMessageEnabled, strings.TrimSpace(simpleConfig.GoogleReviewsConfigMessageParameter),
		strings.TrimSpace(simpleConfig.GoogleReviewsConfigMultiMessageSeparator),
		simpleConfig.GoogleReviewsConfigUseDatabaseMessage, strings.TrimSpace(simpleConfig.GoogleReviewsConfigMessage),
		simpleConfig.GoogleReviewsConfigSendDelayEnabled, simpleConfig.GoogleReviewsConfigSendDelay,
		simpleConfig.GoogleReviewsConfigDispatcherChecksEnabled, strings.TrimSpace(simpleConfig.GoogleReviewsConfigDispatcherType),
		strings.TrimSpace(simpleConfig.GoogleReviewsConfigDispatcherURL),
		strings.TrimSpace(simpleConfig.GoogleReviewsConfigBookingIdParameter),
		simpleConfig.GoogleReviewsConfigIsBookingForNowDiffMinutes,
		simpleConfig.GoogleReviewsConfigBookingNowPickupToContactMinutes, simpleConfig.GoogleReviewsConfigPreBookingPickupToContactMinutes,
		simpleConfig.GoogleReviewsConfigReplaceTelephoneCountryCode, strings.TrimSpace(simpleConfig.GoogleReviewsConfigReplaceTelephoneCountryCodeWith),
		simpleConfig.GoogleReviewsConfigReviewMasterSMSGatewayEnabled,
		simpleConfig.GoogleReviewsConfigReviewMasterSMSGatewayUseMasterQueue,
		strings.TrimSpace(simpleConfig.GoogleReviewsConfigReviewMasterSMSGatewayPairCode),
		simpleConfig.GoogleReviewsConfigAlternateMessageServiceEnabled, simpleConfig.GoogleReviewsConfigAlternateMessageService,
		strings.TrimSpace(simpleConfig.GoogleReviewsConfigAlternateMessageServiceSecret1),
		strings.TrimSpace(simpleConfig.GoogleReviewsConfigCompanies), simpleConfig.GoogleReviewsConfigBookingSourceMobileAppState,
		simpleConfig.GoogleReviewsConfigAIResponsesEnabled,
		strings.TrimSpace(simpleConfig.GoogleReviewsConfigContactMethod),
		simpleConfig.GoogleReviewsConfigMonthlyReviewAnalysisEnabled,
		simpleConfig.GoogleMyBusinessReviewReplyEnabled,
		strings.TrimSpace(simpleConfig.GoogleMyBusinessLocationName),
		strings.TrimSpace(simpleConfig.GoogleMyBusinessPostalCode),
		simpleConfig.GoogleMyBusinessReplyToUnspecfifiedStarRating,
		strings.TrimSpace(simpleConfig.GoogleMyBusinessUnspecfifiedStarRatingReply),
		simpleConfig.GoogleMyBusinessReplyToOneStarRating,
		strings.TrimSpace(simpleConfig.GoogleMyBusinessOneStarRatingReply),
		simpleConfig.GoogleMyBusinessReplyToTwoStarRating,
		strings.TrimSpace(simpleConfig.GoogleMyBusinessTwoStarRatingReply),
		simpleConfig.GoogleMyBusinessReplyToThreeStarRating,
		strings.TrimSpace(simpleConfig.GoogleMyBusinessThreeStarRatingReply),
		simpleConfig.GoogleMyBusinessReplyToFourStarRating,
		strings.TrimSpace(simpleConfig.GoogleMyBusinessFourStarRatingReply),
		simpleConfig.GoogleMyBusinessReplyToFiveStarRating,
		strings.TrimSpace(simpleConfig.GoogleMyBusinessFiveStarRatingReply),
		simpleConfig.GoogleMyBusinessReportEnabled,
		strings.TrimSpace(simpleConfig.EmailAddress),
		simpleConfig.GoogleReviewsConfigID)
	if execErr != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("update failed: %v, unable to rollback: %v\n", execErr, rollbackErr)
			return execErr
		}
		log.Printf("update failed: %v", execErr)
		return execErr
	}

	_, execErr = tx.Exec(googleReviewsConfigTimeQry, simpleConfig.GoogleReviewsConfigTimeEnabled,
		strings.TrimSpace(simpleConfig.GoogleReviewsConfigTimeStart), strings.TrimSpace(simpleConfig.GoogleReviewsConfigTimeEnd),
		simpleConfig.GoogleReviewsConfigTimeSunday, simpleConfig.GoogleReviewsConfigTimeMonday,
		simpleConfig.GoogleReviewsConfigTimeTuesday, simpleConfig.GoogleReviewsConfigTimeWednesday,
		simpleConfig.GoogleReviewsConfigTimeThursday, simpleConfig.GoogleReviewsConfigTimeFriday,
		simpleConfig.GoogleReviewsConfigTimeSaturday,
		simpleConfig.GoogleReviewsConfigTimeID)
	if execErr != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("update failed: %v, unable to rollback: %v\n", execErr, rollbackErr)
			return execErr
		}
		log.Printf("update failed: %v", execErr)
		return execErr
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// CreateSimpleClient - create a client and config and time, simple assume only 1 config and time
func CreateSimpleClient(simpleConfig SimpleConfig, partnerID int) error {
	// var err error

	const clientIDQry = "SELECT MAX(ID) + 1 FROM clients"
	const clientQry = "INSERT INTO clients (id, enabled, name, note, country, partner_id) VALUES (?, ?, ?, ?, ?, ?)"
	const googleReviewsConfigQry = "INSERT INTO google_reviews_configs (enabled," +
		" min_send_frequency, max_send_count," +
		" max_daily_send_count, token, telephone_parameter," +
		" send_from_icabbi_app, app_key, secret_key," +
		" send_url, http_get, send_success_response, time_zone," +
		" multi_message_enabled, message_parameter, multi_message_separator," +
		" use_database_message, message," +
		" send_delay_enabled, send_delay," +
		" dispatcher_checks_enabled, dispatcher_type, dispatcher_url," +
		" booking_id_parameter, is_booking_for_now_diff_minutes," +
		" booking_now_pickup_to_contact_minutes, pre_booking_pickup_to_contact_minutes," +
		" replace_telephone_country_code, replace_telephone_country_code_with," +
		" review_master_sms_gateway_enabled," +
		" review_master_sms_gateway_use_master_queue," +
		" review_master_sms_gateway_pair_code," +
		" alternate_message_service_enabled, alternate_message_service, alternate_message_service_secret1," +
		" companies, booking_source_mobile_app_state," +
		" ai_responses_enabled, contact_method, monthly_review_analysis_enabled," +
		" google_my_business_review_reply_enabled," +
		" google_my_business_location_name," +
		" google_my_business_postal_code," +
		" google_my_business_reply_to_unspecfified_star_rating," +
		" google_my_business_unspecfified_star_rating_reply," +
		" google_my_business_reply_to_one_star_rating," +
		" google_my_business_one_star_rating_reply," +
		" google_my_business_reply_to_two_star_rating," +
		" google_my_business_two_star_rating_reply," +
		" google_my_business_reply_to_three_star_rating," +
		" google_my_business_three_star_rating_reply," +
		" google_my_business_reply_to_four_star_rating," +
		" google_my_business_four_star_rating_reply," +
		" google_my_business_reply_to_five_star_rating," +
		" google_my_business_five_star_rating_reply," +
		" google_my_business_report_enabled," +
		" email_address," +
		" client_id)" +
		" VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	const googleReviewsConfigTimeQry = "INSERT INTO google_reviews_config_times (enabled," +
		" start, end, sunday, monday, tuesday, wednesday, thursday, friday, saturday, google_reviews_config_id)" +
		" VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	tx, err := Db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}

	// First, check if the client exists
	var clientExists bool

	// Generate a new client ID if simpleConfig.ClientID is 0
	if simpleConfig.ClientID == 0 {
		var clientNullID sql.NullInt64
		err = tx.QueryRow(clientIDQry).Scan(&clientNullID)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("query error: %v\n", err)
			return err
		}

		if clientNullID.Valid {
			simpleConfig.ClientID = uint64(clientNullID.Int64)
		} else {
			simpleConfig.ClientID = 1
		}

		clientExists = false
	} else {
		// Check if the specified client exists
		err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM clients WHERE id = ?)", simpleConfig.ClientID).Scan(&clientExists)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("query failed: %v, unable to rollback: %v\n", err, rollbackErr)
				return err
			}
			log.Printf("query failed: %v", err)
			return err
		}
	}

	// If client doesn't exist, create it
	if !clientExists {
		_, execErr := tx.Exec(clientQry, simpleConfig.ClientID, simpleConfig.ClientEnabled,
			strings.TrimSpace(simpleConfig.ClientName), strings.TrimSpace(simpleConfig.ClientNote),
			strings.TrimSpace(simpleConfig.ClientCountry), partnerID)
		if execErr != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("create failed: %v, unable to rollback: %v\n", execErr, rollbackErr)
				return execErr
			}
			log.Printf("create failed: %v", execErr)
			return execErr
		}
	}

	// Now create the Google Reviews Config
	res, execErr := tx.Exec(googleReviewsConfigQry, simpleConfig.GoogleReviewsConfigEnabled,
		simpleConfig.GoogleReviewsConfigMinSendFrequency, simpleConfig.GoogleReviewsConfigMaxSendCount,
		simpleConfig.GoogleReviewsConfigMaxDailySendCount, strings.TrimSpace(simpleConfig.GoogleReviewsConfigToken),
		strings.TrimSpace(simpleConfig.GoogleReviewsConfigTelephoneParameter), simpleConfig.GoogleReviewsConfigSendFromIcabbiApp,
		strings.TrimSpace(simpleConfig.GoogleReviewsConfigAppKey), strings.TrimSpace(simpleConfig.GoogleReviewsConfigSecretKey),
		strings.TrimSpace(simpleConfig.GoogleReviewsConfigSendURL), simpleConfig.GoogleReviewsConfigHttpGet,
		strings.TrimSpace(simpleConfig.GoogleReviewsConfigSendSuccessResponse), strings.TrimSpace(simpleConfig.GoogleReviewsConfigTimeZone),
		simpleConfig.GoogleReviewsConfigMultiMessageEnabled, strings.TrimSpace(simpleConfig.GoogleReviewsConfigMessageParameter),
		strings.TrimSpace(simpleConfig.GoogleReviewsConfigMultiMessageSeparator),
		simpleConfig.GoogleReviewsConfigUseDatabaseMessage, strings.TrimSpace(simpleConfig.GoogleReviewsConfigMessage),
		simpleConfig.GoogleReviewsConfigSendDelayEnabled, simpleConfig.GoogleReviewsConfigSendDelay,
		simpleConfig.GoogleReviewsConfigDispatcherChecksEnabled, strings.TrimSpace(simpleConfig.GoogleReviewsConfigDispatcherType),
		strings.TrimSpace(simpleConfig.GoogleReviewsConfigDispatcherURL),
		strings.TrimSpace(simpleConfig.GoogleReviewsConfigBookingIdParameter),
		simpleConfig.GoogleReviewsConfigIsBookingForNowDiffMinutes,
		simpleConfig.GoogleReviewsConfigBookingNowPickupToContactMinutes, simpleConfig.GoogleReviewsConfigPreBookingPickupToContactMinutes,
		simpleConfig.GoogleReviewsConfigReplaceTelephoneCountryCode, strings.TrimSpace(simpleConfig.GoogleReviewsConfigReplaceTelephoneCountryCodeWith),
		simpleConfig.GoogleReviewsConfigReviewMasterSMSGatewayEnabled,
		simpleConfig.GoogleReviewsConfigReviewMasterSMSGatewayUseMasterQueue,
		simpleConfig.GoogleReviewsConfigReviewMasterSMSGatewayPairCode,
		simpleConfig.GoogleReviewsConfigAlternateMessageServiceEnabled, simpleConfig.GoogleReviewsConfigAlternateMessageService,
		strings.TrimSpace(simpleConfig.GoogleReviewsConfigAlternateMessageServiceSecret1),
		strings.TrimSpace(simpleConfig.GoogleReviewsConfigCompanies), simpleConfig.GoogleReviewsConfigBookingSourceMobileAppState,
		simpleConfig.GoogleReviewsConfigAIResponsesEnabled,
		strings.TrimSpace(simpleConfig.GoogleReviewsConfigContactMethod), simpleConfig.GoogleReviewsConfigMonthlyReviewAnalysisEnabled,
		simpleConfig.GoogleMyBusinessReviewReplyEnabled,
		strings.TrimSpace(simpleConfig.GoogleMyBusinessLocationName),
		strings.TrimSpace(simpleConfig.GoogleMyBusinessPostalCode),
		simpleConfig.GoogleMyBusinessReplyToUnspecfifiedStarRating, strings.TrimSpace(simpleConfig.GoogleMyBusinessUnspecfifiedStarRatingReply),
		simpleConfig.GoogleMyBusinessReplyToOneStarRating, strings.TrimSpace(simpleConfig.GoogleMyBusinessOneStarRatingReply),
		simpleConfig.GoogleMyBusinessReplyToTwoStarRating, strings.TrimSpace(simpleConfig.GoogleMyBusinessTwoStarRatingReply),
		simpleConfig.GoogleMyBusinessReplyToThreeStarRating, strings.TrimSpace(simpleConfig.GoogleMyBusinessThreeStarRatingReply),
		simpleConfig.GoogleMyBusinessReplyToFourStarRating, strings.TrimSpace(simpleConfig.GoogleMyBusinessFourStarRatingReply),
		simpleConfig.GoogleMyBusinessReplyToFiveStarRating, strings.TrimSpace(simpleConfig.GoogleMyBusinessFiveStarRatingReply),
		simpleConfig.GoogleMyBusinessReportEnabled, strings.TrimSpace(simpleConfig.EmailAddress),
		simpleConfig.ClientID)
	if execErr != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("create failed: %v, unable to rollback: %v\n", execErr, rollbackErr)
			return execErr
		}
		log.Printf("create failed: %v", execErr)
		return execErr
	}

	// Get the ID of the newly created Google Reviews Config
	googleReviewsConfigID, err := res.LastInsertId()
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("create failed: %v, unable to rollback: %v\n", err, rollbackErr)
			return err
		}
		log.Printf("create failed: %v", err)
		return err
	}

	// Create the Google Reviews Config Time
	_, execErr = tx.Exec(googleReviewsConfigTimeQry, simpleConfig.GoogleReviewsConfigTimeEnabled,
		strings.TrimSpace(simpleConfig.GoogleReviewsConfigTimeStart), strings.TrimSpace(simpleConfig.GoogleReviewsConfigTimeEnd),
		simpleConfig.GoogleReviewsConfigTimeSunday, simpleConfig.GoogleReviewsConfigTimeMonday,
		simpleConfig.GoogleReviewsConfigTimeTuesday, simpleConfig.GoogleReviewsConfigTimeWednesday,
		simpleConfig.GoogleReviewsConfigTimeThursday, simpleConfig.GoogleReviewsConfigTimeFriday,
		simpleConfig.GoogleReviewsConfigTimeSaturday,
		googleReviewsConfigID)
	if execErr != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("create failed: %v, unable to rollback: %v\n", execErr, rollbackErr)
			return execErr
		}
		log.Printf("create failed: %v", execErr)
		return execErr
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// GetClient - get a client, with configs and times
func GetClient(clientID int, partnerID int) (ClientConfig, error) {
	var c ClientConfig

	const clientQry = "SELECT id, enabled, name, note, country, partner_id" +
		" FROM clients" +
		" WHERE id = ? AND partner_id = ?"

	const configQry = "SELECT" +
		" id, enabled, min_send_frequency, max_send_count," +
		" max_daily_send_count, token, telephone_parameter," +
		" send_from_icabbi_app, app_key, secret_key," +
		" send_url, http_get, send_success_response, time_zone," +
		" multi_message_enabled, message_parameter, multi_message_separator," +
		" use_database_message, message," +
		" send_delay_enabled, send_delay," +
		" dispatcher_checks_enabled, dispatcher_type, dispatcher_url," +
		" booking_id_parameter, is_booking_for_now_diff_minutes," +
		" booking_now_pickup_to_contact_minutes, pre_booking_pickup_to_contact_minutes," +
		" replace_telephone_country_code, replace_telephone_country_code_with," +
		" review_master_sms_gateway_enabled," +
		" review_master_sms_gateway_use_master_queue," +
		" review_master_sms_gateway_pair_code," +
		" alternate_message_service_enabled, alternate_message_service, alternate_message_service_secret1," +
		" companies, booking_source_mobile_app_state," +
		" IFNULL(ai_responses_enabled, 0) as ai_responses_enabled, IFNULL(contact_method, '') as contact_method," +
		" IFNULL(monthly_review_analysis_enabled, 0) as monthly_review_analysis_enabled," +
		" google_my_business_review_reply_enabled," +
		" google_my_business_location_name," +
		" google_my_business_postal_code," +
		" google_my_business_reply_to_unspecfified_star_rating," +
		" google_my_business_unspecfified_star_rating_reply," +
		" google_my_business_reply_to_one_star_rating," +
		" google_my_business_one_star_rating_reply," +
		" google_my_business_reply_to_two_star_rating," +
		" google_my_business_two_star_rating_reply," +
		" google_my_business_reply_to_three_star_rating," +
		" google_my_business_three_star_rating_reply," +
		" google_my_business_reply_to_four_star_rating," +
		" google_my_business_four_star_rating_reply," +
		" google_my_business_reply_to_five_star_rating," +
		" google_my_business_five_star_rating_reply," +
		" google_my_business_report_enabled," +
		" email_address," +
		" client_id" +
		" FROM  google_reviews_configs" +
		" WHERE client_id = ?"

	const timeQry = "SELECT" +
		" id, enabled, start, end, sunday, monday, tuesday, wednesday, thursday, friday, saturday, google_reviews_config_id" +
		" FROM google_reviews_config_times" +
		" WHERE google_reviews_config_id = ?"

	// client - There will only be one result
	client := Client{}
	err := Db.QueryRow(clientQry, clientID, partnerID).Scan(
		&client.ID,
		&client.Enabled,
		&client.Name,
		&client.Note,
		&client.Country,
		&client.PartnerID)
	if err != nil {
		log.Println(err)
		return c, err
	}

	// configs
	configRows, err := Db.Query(configQry, clientID)
	if err != nil {
		log.Println(err)
		return c, err
	}
	defer configRows.Close()
	configRes := []Config{}
	for configRows.Next() {
		grc := GoogleReviewsConfig{}
		conf := Config{}
		err = configRows.Scan(&grc.ID, &grc.Enabled, &grc.MinSendFrequency, &grc.MaxSendCount, &grc.MaxDailySendCount,
			&grc.Token, &grc.TelephoneParameter, &grc.SendFromIcabbiApp, &grc.AppKey, &grc.SecretKey, &grc.SendURL,
			&grc.HttpGet, &grc.SendSuccessResponse, &grc.TimeZone, &grc.MultiMessageEnabled, &grc.MessageParameter,
			&grc.MultiMessageSeparator, &grc.UseDatabaseMessage, &grc.Message,
			&grc.SendDelayEnabled, &grc.SendDelay, &grc.DispatcherChecksEnabled,
			&grc.DispatcherType, &grc.DispatcherURL, &grc.BookingIdParameter, &grc.IsBookingForNowDiffMinutes,
			&grc.BookingNowPickupToContactMinutes, &grc.PreBookingPickupToContactMinutes,
			&grc.ReplaceTelephoneCountryCode, &grc.ReplaceTelephoneCountryCodeWith,
			&grc.ReviewMasterSMSGatewayEnabled, &grc.ReviewMasterSMSGatewayUseMasterQueue,
			&grc.ReviewMasterSMSGatewayPairCode,
			&grc.AlternateMessageServiceEnabled, &grc.AlternateMessageService, &grc.AlternateMessageServiceSecret1,
			&grc.Companies, &grc.BookingSourceMobileAppState,
			&grc.AIResponsesEnabled,
			&grc.ContactMethod,
			&grc.MonthlyReviewAnalysisEnabled,
			&grc.GoogleMyBusinessReviewReplyEnabled,
			&grc.GoogleMyBusinessLocationName,
			&grc.GoogleMyBusinessPostalCode,
			&grc.GoogleMyBusinessReplyToUnspecfifiedStarRating,
			&grc.GoogleMyBusinessUnspecfifiedStarRatingReply,
			&grc.GoogleMyBusinessReplyToOneStarRating,
			&grc.GoogleMyBusinessOneStarRatingReply,
			&grc.GoogleMyBusinessReplyToTwoStarRating,
			&grc.GoogleMyBusinessTwoStarRatingReply,
			&grc.GoogleMyBusinessReplyToThreeStarRating,
			&grc.GoogleMyBusinessThreeStarRatingReply,
			&grc.GoogleMyBusinessReplyToFourStarRating,
			&grc.GoogleMyBusinessFourStarRatingReply,
			&grc.GoogleMyBusinessReplyToFiveStarRating,
			&grc.GoogleMyBusinessFiveStarRatingReply,
			&grc.GoogleMyBusinessReportEnabled,
			&grc.EmailAddress, &grc.ClientID)
		if err != nil {
			return c, err
		}

		// times
		timeRows, err := Db.Query(timeQry, grc.ID)
		if err != nil {
			return c, err
		}
		defer timeRows.Close()
		timeRes := []GoogleReviewsConfigTime{}
		for timeRows.Next() {
			grct := GoogleReviewsConfigTime{}
			err = timeRows.Scan(&grct.ID, &grct.Enabled, &grct.Start, &grct.End, &grct.Sunday, &grct.Monday, &grct.Tuesday, &grct.Wednesday, &grct.Thursday, &grct.Friday, &grct.Saturday, &grct.GoogleReviewsConfigID)
			if err != nil {
				return c, err
			}
			timeRes = append(timeRes, grct)
		}
		conf.GoogleReviewsConfigTimes = timeRes
		conf.GoogleReviewsConfig = grc

		configRes = append(configRes, conf)
	}

	c.Client = client
	c.Configs = configRes
	return c, nil
}

// UpdateClient - update a client and configs and times
func UpdateClient(clientConfig ClientConfig, partnerID int) error {

	const clientQry = "UPDATE clients SET enabled = ?, name = ?, note = ?, country = ?" +
		" WHERE id = ? AND partner_id = ?"
	const googleReviewsConfigQry = "UPDATE google_reviews_configs SET enabled = ?," +
		" min_send_frequency = ?, max_send_count = ?," +
		" max_daily_send_count = ?, token = ?, telephone_parameter = ?," +
		" send_from_icabbi_app = ?, app_key = ?, secret_key = ?," +
		" send_url = ?, http_get = ?, send_success_response = ?, time_zone = ?," +
		" multi_message_enabled = ?, message_parameter = ?, multi_message_separator = ?," +
		" use_database_message = ?, message = ?," +
		" send_delay_enabled = ?, send_delay = ?," +
		" dispatcher_checks_enabled = ?, dispatcher_type = ?, dispatcher_url = ?," +
		" booking_id_parameter = ?, is_booking_for_now_diff_minutes = ?," +
		" booking_now_pickup_to_contact_minutes = ?, pre_booking_pickup_to_contact_minutes = ?," +
		" replace_telephone_country_code = ?, replace_telephone_country_code_with = ?," +
		" review_master_sms_gateway_enabled = ?," +
		" review_master_sms_gateway_use_master_queue = ?," +
		" review_master_sms_gateway_pair_code = ?," +
		" alternate_message_service_enabled = ?, alternate_message_service = ?, alternate_message_service_secret1 = ?," +
		" companies = ?, booking_source_mobile_app_state = ?," +
		" ai_responses_enabled = ?," +
		" contact_method = NULLIF(?, '')," + // Use NULLIF to convert empty string to NULL
		" monthly_review_analysis_enabled = ?," +
		" google_my_business_review_reply_enabled = ?," +
		" google_my_business_location_name = ?," +
		" google_my_business_postal_code = ?," +
		" google_my_business_reply_to_unspecfified_star_rating = ?," +
		" google_my_business_unspecfified_star_rating_reply = ?," +
		" google_my_business_reply_to_one_star_rating = ?," +
		" google_my_business_one_star_rating_reply = ?," +
		" google_my_business_reply_to_two_star_rating = ?," +
		" google_my_business_two_star_rating_reply = ?," +
		" google_my_business_reply_to_three_star_rating = ?," +
		" google_my_business_three_star_rating_reply = ?," +
		" google_my_business_reply_to_four_star_rating = ?," +
		" google_my_business_four_star_rating_reply = ?," +
		" google_my_business_reply_to_five_star_rating = ?," +
		" google_my_business_five_star_rating_reply = ?," +
		" google_my_business_report_enabled = ?," +
		" email_address = ?" +
		" WHERE id = ?"
	const googleReviewsConfigTimeQry = "UPDATE google_reviews_config_times SET enabled = ?," +
		" start = ?, end = ?," +
		" sunday = ?, monday = ?, tuesday = ?, wednesday = ?, thursday = ?, friday = ?, saturday = ?" +
		" WHERE id = ?"

	tx, err := Db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}

	_, execErr := tx.Exec(clientQry, clientConfig.Client.Enabled, clientConfig.Client.Name,
		clientConfig.Client.Note, clientConfig.Client.Country, clientConfig.Client.ID, partnerID)
	if execErr != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("update failed: %v, unable to rollback: %v\n", execErr, rollbackErr)
			return execErr
		}
		log.Printf("update failed: %v", execErr)
		return execErr
	}

	for _, config := range clientConfig.Configs {
		// Add debug logging
		log.Printf("Updating config ID: %d with ai_responses_enabled: %v, contact_method: %q\n",
			config.GoogleReviewsConfig.ID,
			config.GoogleReviewsConfig.AIResponsesEnabled,
			config.GoogleReviewsConfig.ContactMethod)

		_, execErr = tx.Exec(googleReviewsConfigQry, config.GoogleReviewsConfig.Enabled,
			config.GoogleReviewsConfig.MinSendFrequency, config.GoogleReviewsConfig.MaxSendCount,
			config.GoogleReviewsConfig.MaxDailySendCount, strings.TrimSpace(config.GoogleReviewsConfig.Token),
			strings.TrimSpace(config.GoogleReviewsConfig.TelephoneParameter), config.GoogleReviewsConfig.SendFromIcabbiApp,
			strings.TrimSpace(config.GoogleReviewsConfig.AppKey), strings.TrimSpace(config.GoogleReviewsConfig.SecretKey),
			strings.TrimSpace(config.GoogleReviewsConfig.SendURL), config.GoogleReviewsConfig.HttpGet,
			strings.TrimSpace(config.GoogleReviewsConfig.SendSuccessResponse), strings.TrimSpace(config.GoogleReviewsConfig.TimeZone),
			config.GoogleReviewsConfig.MultiMessageEnabled, strings.TrimSpace(config.GoogleReviewsConfig.MessageParameter),
			strings.TrimSpace(config.GoogleReviewsConfig.MultiMessageSeparator),
			config.GoogleReviewsConfig.UseDatabaseMessage, strings.TrimSpace(config.GoogleReviewsConfig.Message),
			config.GoogleReviewsConfig.SendDelayEnabled, config.GoogleReviewsConfig.SendDelay,
			config.GoogleReviewsConfig.DispatcherChecksEnabled, strings.TrimSpace(config.GoogleReviewsConfig.DispatcherType),
			strings.TrimSpace(config.GoogleReviewsConfig.DispatcherURL),
			strings.TrimSpace(config.GoogleReviewsConfig.BookingIdParameter),
			config.GoogleReviewsConfig.IsBookingForNowDiffMinutes,
			config.GoogleReviewsConfig.BookingNowPickupToContactMinutes, config.GoogleReviewsConfig.PreBookingPickupToContactMinutes,
			config.GoogleReviewsConfig.ReplaceTelephoneCountryCode, strings.TrimSpace(config.GoogleReviewsConfig.ReplaceTelephoneCountryCodeWith),
			config.GoogleReviewsConfig.ReviewMasterSMSGatewayEnabled,
			config.GoogleReviewsConfig.ReviewMasterSMSGatewayUseMasterQueue,
			strings.TrimSpace(config.GoogleReviewsConfig.ReviewMasterSMSGatewayPairCode),
			config.GoogleReviewsConfig.AlternateMessageServiceEnabled, strings.TrimSpace(config.GoogleReviewsConfig.AlternateMessageService),
			strings.TrimSpace(config.GoogleReviewsConfig.AlternateMessageServiceSecret1),
			strings.TrimSpace(config.GoogleReviewsConfig.Companies), config.GoogleReviewsConfig.BookingSourceMobileAppState,
			config.GoogleReviewsConfig.AIResponsesEnabled,
			strings.TrimSpace(config.GoogleReviewsConfig.ContactMethod),
			config.GoogleReviewsConfig.MonthlyReviewAnalysisEnabled,
			config.GoogleReviewsConfig.GoogleMyBusinessReviewReplyEnabled,
			strings.TrimSpace(config.GoogleReviewsConfig.GoogleMyBusinessLocationName),
			strings.TrimSpace(config.GoogleReviewsConfig.GoogleMyBusinessPostalCode),
			config.GoogleReviewsConfig.GoogleMyBusinessReplyToUnspecfifiedStarRating, strings.TrimSpace(config.GoogleReviewsConfig.GoogleMyBusinessUnspecfifiedStarRatingReply),
			config.GoogleReviewsConfig.GoogleMyBusinessReplyToOneStarRating, strings.TrimSpace(config.GoogleReviewsConfig.GoogleMyBusinessOneStarRatingReply),
			config.GoogleReviewsConfig.GoogleMyBusinessReplyToTwoStarRating, strings.TrimSpace(config.GoogleReviewsConfig.GoogleMyBusinessTwoStarRatingReply),
			config.GoogleReviewsConfig.GoogleMyBusinessReplyToThreeStarRating, strings.TrimSpace(config.GoogleReviewsConfig.GoogleMyBusinessThreeStarRatingReply),
			config.GoogleReviewsConfig.GoogleMyBusinessReplyToFourStarRating, strings.TrimSpace(config.GoogleReviewsConfig.GoogleMyBusinessFourStarRatingReply),
			config.GoogleReviewsConfig.GoogleMyBusinessReplyToFiveStarRating, strings.TrimSpace(config.GoogleReviewsConfig.GoogleMyBusinessFiveStarRatingReply),
			config.GoogleReviewsConfig.GoogleMyBusinessReportEnabled, strings.TrimSpace(config.GoogleReviewsConfig.EmailAddress),
			config.GoogleReviewsConfig.ID)
		if execErr != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("update failed: %v, unable to rollback: %v\n", execErr, rollbackErr)
				return execErr
			}
			log.Printf("update failed: %v", execErr)
			return execErr
		}

		for _, configTime := range config.GoogleReviewsConfigTimes {
			_, execErr = tx.Exec(googleReviewsConfigTimeQry, configTime.Enabled,
				strings.TrimSpace(configTime.Start), strings.TrimSpace(configTime.End),
				configTime.Sunday, configTime.Monday, configTime.Tuesday, configTime.Wednesday,
				configTime.Thursday, configTime.Friday, configTime.Saturday,
				configTime.ID)
			if execErr != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					log.Printf("update failed: %v, unable to rollback: %v\n", execErr, rollbackErr)
					return execErr
				}
				log.Printf("update failed: %v", execErr)
				return execErr
			}
		}
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// CreateClient - create a client and configs and times
func CreateClient(clientConfig ClientConfig, partnerID int) error {
	// var err error

	const clientIDQry = "SELECT MAX(ID) + 1 FROM clients"
	const clientQry = "INSERT INTO clients (id, enabled, name, note, country, partner_id) VALUES (?, ?, ?, ?, ?, ?)"
	const googleReviewsConfigQry = "INSERT INTO google_reviews_configs (enabled," +
		" min_send_frequency, max_send_count," +
		" max_daily_send_count, token, telephone_parameter," +
		" send_from_icabbi_app, app_key, secret_key," +
		" send_url, http_get, send_success_response, time_zone," +
		" multi_message_enabled, message_parameter, multi_message_separator," +
		" use_database_message, message," +
		" send_delay_enabled, send_delay," +
		" dispatcher_checks_enabled, dispatcher_type, dispatcher_url, booking_id_parameter," +
		" is_booking_for_now_diff_minutes," +
		" booking_now_pickup_to_contact_minutes, pre_booking_pickup_to_contact_minutes," +
		" replace_telephone_country_code, replace_telephone_country_code_with, " +
		" review_master_sms_gateway_enabled," +
		" review_master_sms_gateway_use_master_queue," +
		" review_master_sms_gateway_pair_code, " +
		" alternate_message_service_enabled, alternate_message_service, alternate_message_service_secret1, " +
		" companies, booking_source_mobile_app_state, " +
		" ai_responses_enabled, contact_method, monthly_review_analysis_enabled, " +
		" google_my_business_review_reply_enabled," +
		" google_my_business_location_name," +
		" google_my_business_postal_code," +
		" google_my_business_reply_to_unspecfified_star_rating," +
		" google_my_business_unspecfified_star_rating_reply," +
		" google_my_business_reply_to_one_star_rating," +
		" google_my_business_one_star_rating_reply," +
		" google_my_business_reply_to_two_star_rating," +
		" google_my_business_two_star_rating_reply," +
		" google_my_business_reply_to_three_star_rating," +
		" google_my_business_three_star_rating_reply," +
		" google_my_business_reply_to_four_star_rating," +
		" google_my_business_four_star_rating_reply," +
		" google_my_business_reply_to_five_star_rating," +
		" google_my_business_five_star_rating_reply," +
		" google_my_business_report_enabled," +
		" email_address," +
		" client_id)" +
		" VALUES(?,?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	const googleReviewsConfigTimeQry = "INSERT INTO google_reviews_config_times (enabled," +
		" start, end, sunday, monday, tuesday, wednesday, thursday, friday, saturday, google_reviews_config_id)" +
		" VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	tx, err := Db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}

	var clientID uint64
	var clientNullID sql.NullInt64
	err = tx.QueryRow(clientIDQry).Scan(&clientNullID)
	switch {
	case err == sql.ErrNoRows:
		clientID = 1
	case err != nil:
		log.Printf("query error: %v\n", err)
		return err
	default:
		// clientID = clientID
	}
	if clientNullID.Valid {
		clientID = uint64(clientNullID.Int64)
	} else {
		clientID = 1
	}

	_, execErr := tx.Exec(clientQry, clientID, clientConfig.Client.Enabled, strings.TrimSpace(clientConfig.Client.Name),
		strings.TrimSpace(clientConfig.Client.Note), strings.TrimSpace(clientConfig.Client.Country), partnerID)
	if execErr != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("create failed: %v, unable to rollback: %v\n", execErr, rollbackErr)
			return execErr
		}
		log.Printf("create failed: %v", execErr)
		return execErr
	}

	for _, config := range clientConfig.Configs {
		res, execErr := tx.Exec(googleReviewsConfigQry, config.GoogleReviewsConfig.Enabled,
			config.GoogleReviewsConfig.MinSendFrequency, config.GoogleReviewsConfig.MaxSendCount,
			config.GoogleReviewsConfig.MaxDailySendCount, strings.TrimSpace(config.GoogleReviewsConfig.Token),
			strings.TrimSpace(config.GoogleReviewsConfig.TelephoneParameter), config.GoogleReviewsConfig.SendFromIcabbiApp,
			strings.TrimSpace(config.GoogleReviewsConfig.AppKey), strings.TrimSpace(config.GoogleReviewsConfig.SecretKey),
			strings.TrimSpace(config.GoogleReviewsConfig.SendURL), config.GoogleReviewsConfig.HttpGet,
			strings.TrimSpace(config.GoogleReviewsConfig.SendSuccessResponse), strings.TrimSpace(config.GoogleReviewsConfig.TimeZone),
			config.GoogleReviewsConfig.MultiMessageEnabled, strings.TrimSpace(config.GoogleReviewsConfig.MessageParameter),
			strings.TrimSpace(config.GoogleReviewsConfig.MultiMessageSeparator),
			config.GoogleReviewsConfig.UseDatabaseMessage, strings.TrimSpace(config.GoogleReviewsConfig.Message),
			config.GoogleReviewsConfig.SendDelayEnabled, config.GoogleReviewsConfig.SendDelay,
			config.GoogleReviewsConfig.DispatcherChecksEnabled, config.GoogleReviewsConfig.DispatcherType,
			config.GoogleReviewsConfig.DispatcherURL,
			config.GoogleReviewsConfig.BookingIdParameter,
			config.GoogleReviewsConfig.IsBookingForNowDiffMinutes,
			config.GoogleReviewsConfig.BookingNowPickupToContactMinutes, config.GoogleReviewsConfig.PreBookingPickupToContactMinutes,
			config.GoogleReviewsConfig.ReplaceTelephoneCountryCode, strings.TrimSpace(config.GoogleReviewsConfig.ReplaceTelephoneCountryCodeWith),
			config.GoogleReviewsConfig.ReviewMasterSMSGatewayEnabled,
			config.GoogleReviewsConfig.ReviewMasterSMSGatewayUseMasterQueue,
			strings.TrimSpace(config.GoogleReviewsConfig.ReviewMasterSMSGatewayPairCode),
			config.GoogleReviewsConfig.AlternateMessageServiceEnabled, strings.TrimSpace(config.GoogleReviewsConfig.AlternateMessageService),
			strings.TrimSpace(config.GoogleReviewsConfig.AlternateMessageServiceSecret1),
			strings.TrimSpace(config.GoogleReviewsConfig.Companies), config.GoogleReviewsConfig.BookingSourceMobileAppState,
			config.GoogleReviewsConfig.AIResponsesEnabled,
			strings.TrimSpace(config.GoogleReviewsConfig.ContactMethod),
			config.GoogleReviewsConfig.MonthlyReviewAnalysisEnabled,
			config.GoogleReviewsConfig.GoogleMyBusinessReviewReplyEnabled,
			strings.TrimSpace(config.GoogleReviewsConfig.GoogleMyBusinessLocationName),
			strings.TrimSpace(config.GoogleReviewsConfig.GoogleMyBusinessPostalCode),
			config.GoogleReviewsConfig.GoogleMyBusinessReplyToUnspecfifiedStarRating, strings.TrimSpace(config.GoogleReviewsConfig.GoogleMyBusinessUnspecfifiedStarRatingReply),
			config.GoogleReviewsConfig.GoogleMyBusinessReplyToOneStarRating, strings.TrimSpace(config.GoogleReviewsConfig.GoogleMyBusinessOneStarRatingReply),
			config.GoogleReviewsConfig.GoogleMyBusinessReplyToTwoStarRating, strings.TrimSpace(config.GoogleReviewsConfig.GoogleMyBusinessTwoStarRatingReply),
			config.GoogleReviewsConfig.GoogleMyBusinessReplyToThreeStarRating, strings.TrimSpace(config.GoogleReviewsConfig.GoogleMyBusinessThreeStarRatingReply),
			config.GoogleReviewsConfig.GoogleMyBusinessReplyToFourStarRating, strings.TrimSpace(config.GoogleReviewsConfig.GoogleMyBusinessFourStarRatingReply),
			config.GoogleReviewsConfig.GoogleMyBusinessReplyToFiveStarRating, strings.TrimSpace(config.GoogleReviewsConfig.GoogleMyBusinessFiveStarRatingReply),
			config.GoogleReviewsConfig.GoogleMyBusinessReportEnabled, strings.TrimSpace(config.GoogleReviewsConfig.EmailAddress),
			clientID)
		if execErr != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("create failed: %v, unable to rollback: %v\n", execErr, rollbackErr)
				return execErr
			}
			log.Printf("create failed: %v", execErr)
			return execErr
		}
		googleReviewsConfigID, err := res.LastInsertId()
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("create failed: %v, unable to rollback: %v\n", err, rollbackErr)
				return err
			}
			log.Printf("create failed: %v", err)
			return err
		}

		for _, configTime := range config.GoogleReviewsConfigTimes {
			_, execErr = tx.Exec(googleReviewsConfigTimeQry, configTime.Enabled,
				strings.TrimSpace(configTime.Start), strings.TrimSpace(configTime.End),
				configTime.Sunday, configTime.Monday, configTime.Tuesday, configTime.Wednesday,
				configTime.Thursday, configTime.Friday, configTime.Saturday,
				googleReviewsConfigID)
			if execErr != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					log.Printf("create failed: %v, unable to rollback: %v\n", execErr, rollbackErr)
					return execErr
				}
				log.Printf("create failed: %v", execErr)
				return execErr
			}
		}
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// CreateGRConfig - create a google reviews config
func CreateGRConfig(googleReviewsConfig GoogleReviewsConfig) error {
	// var err error

	const googleReviewsConfigQry = "INSERT INTO google_reviews_configs (enabled," +
		" min_send_frequency, max_send_count," +
		" max_daily_send_count, token, telephone_parameter," +
		" send_from_icabbi_app, app_key, secret_key," +
		" send_url, http_get, send_success_response, time_zone," +
		" multi_message_enabled, message_parameter, multi_message_separator," +
		" use_database_message, message," +
		" send_delay_enabled, send_delay," +
		" dispatcher_checks_enabled, dispatcher_type, dispatcher_url, booking_id_parameter," +
		" is_booking_for_now_diff_minutes," +
		" booking_now_pickup_to_contact_minutes, pre_booking_pickup_to_contact_minutes," +
		" replace_telephone_country_code, replace_telephone_country_code_with, " +
		" review_master_sms_gateway_enabled," +
		" review_master_sms_gateway_use_master_queue," +
		" review_master_sms_gateway_pair_code, " +
		" alternate_message_service_enabled, alternate_message_service, alternate_message_service_secret1, " +
		" companies, booking_source_mobile_app_state, " +
		" ai_responses_enabled, contact_method, monthly_review_analysis_enabled, " +
		" google_my_business_review_reply_enabled," +
		" google_my_business_location_name," +
		" google_my_business_postal_code," +
		" google_my_business_reply_to_unspecfified_star_rating," +
		" google_my_business_unspecfified_star_rating_reply," +
		" google_my_business_reply_to_one_star_rating," +
		" google_my_business_one_star_rating_reply," +
		" google_my_business_reply_to_two_star_rating," +
		" google_my_business_two_star_rating_reply," +
		" google_my_business_reply_to_three_star_rating," +
		" google_my_business_three_star_rating_reply," +
		" google_my_business_reply_to_four_star_rating," +
		" google_my_business_four_star_rating_reply," +
		" google_my_business_reply_to_five_star_rating," +
		" google_my_business_five_star_rating_reply," +
		" google_my_business_report_enabled," +
		" email_address," +
		" client_id)" +
		" VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	tx, err := Db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}

	_, execErr := tx.Exec(googleReviewsConfigQry, googleReviewsConfig.Enabled,
		googleReviewsConfig.MinSendFrequency, googleReviewsConfig.MaxSendCount,
		googleReviewsConfig.MaxDailySendCount, strings.TrimSpace(googleReviewsConfig.Token),
		strings.TrimSpace(googleReviewsConfig.TelephoneParameter), googleReviewsConfig.SendFromIcabbiApp,
		strings.TrimSpace(googleReviewsConfig.AppKey), strings.TrimSpace(googleReviewsConfig.SecretKey),
		strings.TrimSpace(googleReviewsConfig.SendURL), googleReviewsConfig.HttpGet,
		strings.TrimSpace(googleReviewsConfig.SendSuccessResponse), strings.TrimSpace(googleReviewsConfig.TimeZone),
		googleReviewsConfig.MultiMessageEnabled, strings.TrimSpace(googleReviewsConfig.MessageParameter),
		strings.TrimSpace(googleReviewsConfig.MultiMessageSeparator),
		googleReviewsConfig.UseDatabaseMessage, strings.TrimSpace(googleReviewsConfig.Message),
		googleReviewsConfig.SendDelayEnabled, googleReviewsConfig.SendDelay,
		googleReviewsConfig.DispatcherChecksEnabled, googleReviewsConfig.DispatcherType,
		googleReviewsConfig.DispatcherURL,
		googleReviewsConfig.BookingIdParameter,
		googleReviewsConfig.IsBookingForNowDiffMinutes,
		googleReviewsConfig.BookingNowPickupToContactMinutes, googleReviewsConfig.PreBookingPickupToContactMinutes,
		googleReviewsConfig.ReplaceTelephoneCountryCode, strings.TrimSpace(googleReviewsConfig.ReplaceTelephoneCountryCodeWith),
		googleReviewsConfig.ReviewMasterSMSGatewayEnabled,
		googleReviewsConfig.ReviewMasterSMSGatewayUseMasterQueue,
		strings.TrimSpace(googleReviewsConfig.ReviewMasterSMSGatewayPairCode),
		googleReviewsConfig.AlternateMessageServiceEnabled, strings.TrimSpace(googleReviewsConfig.AlternateMessageService),
		strings.TrimSpace(googleReviewsConfig.AlternateMessageServiceSecret1),
		strings.TrimSpace(googleReviewsConfig.Companies), googleReviewsConfig.BookingSourceMobileAppState,
		googleReviewsConfig.AIResponsesEnabled,
		strings.TrimSpace(googleReviewsConfig.ContactMethod),
		googleReviewsConfig.MonthlyReviewAnalysisEnabled,
		googleReviewsConfig.GoogleMyBusinessReviewReplyEnabled,
		strings.TrimSpace(googleReviewsConfig.GoogleMyBusinessLocationName),
		strings.TrimSpace(googleReviewsConfig.GoogleMyBusinessPostalCode),
		googleReviewsConfig.GoogleMyBusinessReplyToUnspecfifiedStarRating, strings.TrimSpace(googleReviewsConfig.GoogleMyBusinessUnspecfifiedStarRatingReply),
		googleReviewsConfig.GoogleMyBusinessReplyToOneStarRating, strings.TrimSpace(googleReviewsConfig.GoogleMyBusinessOneStarRatingReply),
		googleReviewsConfig.GoogleMyBusinessReplyToTwoStarRating, strings.TrimSpace(googleReviewsConfig.GoogleMyBusinessTwoStarRatingReply),
		googleReviewsConfig.GoogleMyBusinessReplyToThreeStarRating, strings.TrimSpace(googleReviewsConfig.GoogleMyBusinessThreeStarRatingReply),
		googleReviewsConfig.GoogleMyBusinessReplyToFourStarRating, strings.TrimSpace(googleReviewsConfig.GoogleMyBusinessFourStarRatingReply),
		googleReviewsConfig.GoogleMyBusinessReplyToFiveStarRating, strings.TrimSpace(googleReviewsConfig.GoogleMyBusinessFiveStarRatingReply),
		googleReviewsConfig.GoogleMyBusinessReportEnabled, strings.TrimSpace(googleReviewsConfig.EmailAddress),
		googleReviewsConfig.ClientID)
	if execErr != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("create failed: %v, unable to rollback: %v\n", execErr, rollbackErr)
			return execErr
		}
		log.Printf("create failed: %v", execErr)
		return execErr
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// CreateGRCTime - create a google review time
func CreateGRCTime(googleReviewsConfigTime GoogleReviewsConfigTime) error {
	// var err error

	const googleReviewsConfigTimeQry = "INSERT INTO google_reviews_config_times (enabled," +
		" start, end, sunday, monday, tuesday, wednesday, thursday, friday, saturday, google_reviews_config_id)" +
		" VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	tx, err := Db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}

	_, execErr := tx.Exec(googleReviewsConfigTimeQry, googleReviewsConfigTime.Enabled,
		strings.TrimSpace(googleReviewsConfigTime.Start), strings.TrimSpace(googleReviewsConfigTime.End),
		googleReviewsConfigTime.Sunday, googleReviewsConfigTime.Monday, googleReviewsConfigTime.Tuesday,
		googleReviewsConfigTime.Wednesday, googleReviewsConfigTime.Thursday, googleReviewsConfigTime.Friday,
		googleReviewsConfigTime.Saturday,
		googleReviewsConfigTime.GoogleReviewsConfigID)
	if execErr != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("create failed: %v, unable to rollback: %v\n", execErr, rollbackErr)
			return execErr
		}
		log.Printf("create failed: %v", execErr)
		return execErr
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// // Stats - get some stats
// func Stats(db *sql.DB, statsRequest StatsRequest, partnerID int) ([]StatsResult, error) {
// 	var qry = "SELECT g.client_id AS client_id, c.name AS client_name, " +
// 		" count(g.id) AS sent, "
// 	const qryEnd = " AS group_period" +
// 		" FROM clients AS c" +
// 		" LEFT JOIN google_reviews_last_sents AS g ON c.id = g.client_id" +
// 		" WHERE c.partner_id = ? AND g.last_sent > ? AND g.last_sent < ?" +
// 		" GROUP BY client_id, group_period"
// 	timeGrouping := "DATE(g.last_sent)"
// 	switch statsRequest.TimeGrouping {
// 	case "Week":
// 		timeGrouping = "WEEK(g.last_sent)"
// 	case "Month":
// 		timeGrouping = "MONTH(g.last_sent)"
// 	case "Year":
// 		timeGrouping = "YEAR(g.last_sent)"
// 	}
// 	rows, err := db.Query(qry+timeGrouping+qryEnd, partnerID, statsRequest.StartDay, statsRequest.EndDay)
// 	if err != nil {
// 		log.Println(err)
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	var statsResult []StatsResult
// 	for rows.Next() {
// 		var s StatsResult
// 		if err := rows.Scan(&s.ClientID, &s.ClientName, &s.Sent, &s.GroupPeriod); err != nil {
// 			log.Printf("Error getting stats results: %v\n", err)
// 		}
// 		statsResult = append(statsResult, s)
// 	}
// 	return statsResult, nil
// }

// Stats - get some stats
// The query has been modified to take advantage of indexes
func Stats(statsRequest StatsRequest, partnerID int) ([]StatsResult, error) {
	var outerQryStart = "SELECT a.client_id, a.client_name, a.sent, a.group_period FROM ("
	var outerQryEnd = ") AS a WHERE a.partner_id = ?"
	var qry = "SELECT g.client_id AS client_id, c.name AS client_name, c.partner_id AS partner_id," +
		" count(g.id) AS sent, "
	const qryEnd = " AS group_period" +
		" FROM google_reviews_last_sents AS g" +
		" LEFT JOIN clients AS c ON g.client_id = c.id" +
		" WHERE g.last_sent_date BETWEEN ? AND ?" +
		" GROUP BY client_id, group_period"
	timeGrouping := "DATE(g.last_sent)"
	switch statsRequest.TimeGrouping {
	case "Week":
		timeGrouping = "WEEK(g.last_sent)"
	case "Month":
		timeGrouping = "MONTH(g.last_sent)"
	case "Year":
		timeGrouping = "YEAR(g.last_sent)"
	}
	rows, err := Db.Query(outerQryStart+qry+timeGrouping+qryEnd+outerQryEnd, statsRequest.StartDay, statsRequest.EndDay, partnerID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	var statsResult []StatsResult
	for rows.Next() {
		var s StatsResult
		if err := rows.Scan(&s.ClientID, &s.ClientName, &s.Sent, &s.GroupPeriod); err != nil {
			log.Printf("Error getting stats results: %v\n", err)
		}
		statsResult = append(statsResult, s)
	}
	return statsResult, nil
}

// StatsNew - get some stats using the stats table
func StatsNew(statsRequest StatsRequest, partnerID int) ([]StatsNewResult, error) {
	var outerQryStart = "SELECT a.client_id, a.client_name, a.sent, a.requested, a.group_period FROM ("
	var outerQryEnd = ") AS a WHERE a.partner_id = ?"
	var qry = "SELECT g.client_id AS client_id, c.name AS client_name, c.partner_id AS partner_id," +
		" sum(g.sent_count) AS sent,  sum(g.requested_count) AS requested, "
	const qryEnd = " AS group_period" +
		" FROM google_reviews_stats AS g" +
		" LEFT JOIN clients AS c ON g.client_id = c.id" +
		" WHERE g.stats_date BETWEEN ? AND ?" +
		" GROUP BY client_id, group_period"
	timeGrouping := "DATE(g.stats_date)"
	switch statsRequest.TimeGrouping {
	case "Week":
		timeGrouping = "WEEK(g.stats_date)"
	case "Month":
		timeGrouping = "MONTH(g.stats_date)"
	case "Year":
		timeGrouping = "YEAR(g.stats_date)"
	}
	rows, err := Db.Query(outerQryStart+qry+timeGrouping+qryEnd+outerQryEnd, statsRequest.StartDay, statsRequest.EndDay, partnerID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	var statsResult []StatsNewResult
	for rows.Next() {
		var s StatsNewResult
		if err := rows.Scan(&s.ClientID, &s.ClientName, &s.Sent, &s.Requested, &s.GroupPeriod); err != nil {
			log.Printf("Error getting stats results: %v\n", err)
		}
		statsResult = append(statsResult, s)
	}
	return statsResult, nil
}

// // CheckNothingSent - check to see if no messages have been sent for companies in a set time period
// func CheckNothingSent(db *sql.DB, hoursBack int, partnerID int) ([]NothingSentResult, error) {
// 	var qry = "SELECT c.id AS client_id, c.name AS client_name" +
// 		" FROM clients AS c" +
// 		" WHERE c.partner_id = ?" +
// 		" AND c.enabled" +
// 		" AND c.id NOT IN" +
// 		" (" +
// 		" SELECT g.client_id" +
// 		" FROM google_reviews_last_sents AS g" +
// 		" WHERE g.last_sent > NOW() - INTERVAL ? HOUR" +
// 		" )"
// 	// fmt.Println(qry)
// 	rows, err := db.Query(qry, partnerID, hoursBack)
// 	if err != nil {
// 		log.Println(err)
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	var nothingSentResults []NothingSentResult
// 	for rows.Next() {
// 		var n NothingSentResult
// 		if err := rows.Scan(&n.ClientID, &n.ClientName); err != nil {
// 			log.Printf("Error getting nothing sent check results: %+v\n", err)
// 		}
// 		nothingSentResults = append(nothingSentResults, n)
// 	}
// 	return nothingSentResults, nil
// }

// CheckNothingSent - check to see if no messages have been sent for companies in a set time period
// The query has been modified to take advantage of indexes
func CheckNothingSent(daysBack int, partnerID int) ([]NothingSentResult, error) {
	var qry = "SELECT c.id AS client_id, c.name AS client_name" +
		" FROM clients AS c" +
		" WHERE c.partner_id = ?" +
		" AND c.enabled" +
		" AND c.id NOT IN" +
		" (" +
		" SELECT g.client_id" +
		" FROM google_reviews_last_sents AS g" +
		" WHERE g.last_sent_date > DATE(NOW()) - INTERVAL ? DAY" +
		" )"
	// fmt.Println(qry)
	rows, err := Db.Query(qry, partnerID, daysBack)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	var nothingSentResults []NothingSentResult
	for rows.Next() {
		var n NothingSentResult
		if err := rows.Scan(&n.ClientID, &n.ClientName); err != nil {
			log.Printf("Error getting nothing sent check results: %+v\n", err)
		}
		nothingSentResults = append(nothingSentResults, n)
	}
	return nothingSentResults, nil
}

// ListAllUsers - get all the users for a specific partner
func ListAllUsers(partnerID int) ([]UserClientList, error) {
	var err error
	const qry = "SELECT u.id AS user_id, u.email," +
		" c.id AS client_id, c.name" +
		" FROM users_clients AS u_c" +
		" JOIN users AS u ON u.id = u_c.user_id" +
		" JOIN clients AS c ON c.id = u_c.client_id" +
		" WHERE c.partner_id = ?" +
		" ORDER BY u.id"
	rows, err := Db.Query(qry, partnerID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	type ucs struct {
		UserID     uint64
		Email      string
		ClientID   uint64
		ClientName string
	}
	var userClients []UserClientList
	var lastUserID uint64
	var ucl UserClientList
	for rows.Next() {
		var uc ucs
		if err := rows.Scan(&uc.UserID, &uc.Email, &uc.ClientID, &uc.ClientName); err != nil {
			log.Printf("Error getting all user clients: %v\n", err)
		}
		if lastUserID != uc.UserID {
			if lastUserID != 0 {
				userClients = append(userClients, ucl)
			}
			lastUserID = uc.UserID
			ucl.UserID = uc.UserID
			ucl.Email = uc.Email
			ucl.ClientList = uc.ClientName + " (" + strconv.FormatUint(uc.ClientID, 10) + ")"
		} else {
			ucl.ClientList += ", " + uc.ClientName + " (" + strconv.FormatUint(uc.ClientID, 10) + ")"
		}
	}
	if lastUserID != 0 {
		userClients = append(userClients, ucl)
	}
	return userClients, nil
}

// GetUser - get user of frontend
func GetUser(userID int, partnerID int) (UserClients, error) {
	var ucs UserClients
	const qry = "SELECT u.id AS user_id, u.email, u.password," +
		" c.id AS client_id, c.enabled, c.name, c.note, c.country, c.partner_id" +
		" FROM users_clients AS u_c" +
		" JOIN users AS u ON u.id = u_c.user_id" +
		" JOIN clients AS c ON c.id = u_c.client_id" +
		" WHERE u.id = ? AND c.partner_id = ?"
	rows, err := Db.Query(qry, userID, partnerID)
	if err != nil {
		log.Println(err)
		return ucs, err
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		var client Client
		if err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&client.ID,
			&client.Enabled,
			&client.Name,
			&client.Note,
			&client.Country,
			&client.PartnerID); err != nil {
			log.Printf("Error getting user clients for user id: %d, partner id: %d, err: %v\n", userID, partnerID, err)
		} else {
			ucs.User = user
			ucs.Clients = append(ucs.Clients, client)
		}
	}
	return ucs, nil
}

// UpdateUser - update a user and clients
func UpdateUser(userClients UserClients, partnerID int) error {
	const userQry = "SELECT u.id, u.email, u.password FROM users AS u" +
		" WHERE u.id = ?"
	const userUpdateQry1 = "UPDATE users SET email = ?"
	var userUpdateQry2 = ", password = ?"
	const userUpdateQry3 = " WHERE id = ?"
	// Fix: Remove the alias 'AS u' from the DELETE statement
	const userClientsDeleteQry = "DELETE FROM users_clients" +
		" WHERE user_id = ?"
	const userClientsInsertQry = "INSERT INTO users_clients (user_id, client_id) VALUES (?, ?)"
	const clientQry = "SELECT COUNT(c.id) FROM clients AS c" +
		" WHERE c.id = ? and c.partner_id = ?"

	// check that all the clients are associated with the partner
	clients := userClients.Clients
	if len(clients) == 0 {
		err := errors.New("there must be at least one client selected")
		return err
	}
	for _, client := range clients {
		row := Db.QueryRow(clientQry, client.ID, partnerID)
		var count int
		if err := row.Scan(&count); err != nil {
			log.Printf("Error getting client ID: %d for partner ID: %d for user ID: %d, err: %v\n", client.ID, partnerID, userClients.User.ID, err)
			return err
		}
		if count == 0 {
			err := errors.New("Client cannot be found")
			return err
		}
	}

	// check user exists
	row := Db.QueryRow(userQry, userClients.User.ID)
	var user User
	if err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
	); err != nil {
		log.Printf("Error getting user ID: %d, err: %v\n", userClients.User.ID, err)
		return err
	}
	if user.ID == 0 {
		log.Printf("Error no user found for user ID: %d\n", userClients.User.ID)
		err := errors.New("User cannot be found")
		return err
	}

	tx, err := Db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}

	//  check if password has changed
	pw := user.Password
	if user.Password != userClients.User.Password {
		pw = userClients.User.Password
		userUpdateQry2 = ", password = SHA2(?, 224)"
	}
	// update user
	userUpdateQry := userUpdateQry1 + userUpdateQry2 + userUpdateQry3
	_, execErr := tx.Exec(userUpdateQry, strings.TrimSpace(userClients.User.Email), pw, user.ID)
	if execErr != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("update failed: %v, unable to rollback: %v\n", execErr, rollbackErr)
			return execErr
		}
		log.Printf("update failed: %v", execErr)
		return execErr
	}
	// update user clients: by deleting user clients followed by inserting user clients
	_, execErr = tx.Exec(userClientsDeleteQry, user.ID)
	if execErr != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("update failed: %v, unable to rollback: %v\n", execErr, rollbackErr)
			return execErr
		}
		log.Printf("update failed: %v", execErr)
		return execErr
	}
	for _, client := range clients {
		_, execErr = tx.Exec(userClientsInsertQry, user.ID, client.ID)
		if execErr != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("update failed: %v, unable to rollback: %v\n", execErr, rollbackErr)
				return execErr
			}
			log.Printf("update failed: %v", execErr)
			return execErr
		}
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// CreateUser - Create a user with associated clients
func CreateUser(userClients UserClients, partnerID int) error {
	const userQry = "SELECT COUNT(u.id) FROM users AS u" +
		" WHERE u.email = ?"
	const userCreateQry = "INSERT INTO users (email, password) VALUES (?, SHA2(?, 224))"
	const userClientsInsertQry = "INSERT INTO users_clients (user_id, client_id) VALUES (?, ?)"
	const clientQry = "SELECT COUNT(c.id) FROM clients AS c" +
		" WHERE c.id = ? and c.partner_id = ?"

	// check that all the clients are associated with the partner
	clients := userClients.Clients
	if len(clients) == 0 {
		err := errors.New("there must be at least one client selected")
		return err
	}
	for _, client := range clients {
		row := Db.QueryRow(clientQry, client.ID, partnerID)
		var count int
		if err := row.Scan(&count); err != nil {
			log.Printf("Error getting client ID: %d for partner ID: %d for user ID: %d, err: %v\n", client.ID, partnerID, userClients.User.ID, err)
			return err
		}
		if count == 0 {
			err := errors.New("Client cannot be found")
			return err
		}
	}

	// check user does not exists
	row := Db.QueryRow(userQry, userClients.User.Email)
	var count int
	if err := row.Scan(&count); err != nil {
		log.Printf("Error checking if user exists, err: %v\n", err)
		return err
	}
	if count != 0 {
		log.Printf("Error user already exists user: %s\n", userClients.User.Email)
		err := errors.New("User already exists")
		return err
	}

	tx, err := Db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}

	// create user
	data, execErr := tx.Exec(userCreateQry, strings.TrimSpace(userClients.User.Email), userClients.User.Password)
	if execErr != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("create failed: %v, unable to rollback: %v\n", execErr, rollbackErr)
			return execErr
		}
		log.Printf("create failed: %v", execErr)
		return execErr
	}
	userID, _ := data.LastInsertId()
	// insert user clients
	for _, client := range clients {
		_, execErr = tx.Exec(userClientsInsertQry, userID, client.ID)
		if execErr != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("create failed: %v, unable to rollback: %v\n", execErr, rollbackErr)
				return execErr
			}
			log.Printf("create failed: %v", execErr)
			return execErr
		}
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// DeleteUser - delete a user and clients
func DeleteUser(userID int, partnerID int) error {
	// Fix: Remove the alias 'AS u' from the DELETE statement
	const userQry = "DELETE FROM users" +
		" WHERE id = ?"

	// check if user exists this will also check the partner has access to this user
	userClients, error := GetUser(userID, partnerID)
	if error != nil {
		return error
	}
	if userClients.User.ID == 0 {
		return errors.New("user does not exist")
	}

	tx, err := Db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}

	_, execErr := tx.Exec(userQry, userID)
	if execErr != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("delete failed: %v, unable to rollback: %v\n", execErr, rollbackErr)
			return execErr
		}
		log.Printf("delete failed: %v", execErr)
		return execErr
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
