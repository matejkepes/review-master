package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

var Conf Config

// Config - config
type Config struct {
	DbName     string
	DbAddress  string
	DbPort     string
	DbUsername string
	DbPassword string

	PollPeriod   int
	LastPollTime time.Time

	SendSmsURL                string
	SendSmsToken              string
	SendSmsTokenParameter     string
	SendSmsTelephoneParameter string
	SendSmsMessageParameter   string
	SendSmsSuccessResponse    string
	SendSmsFailureResponse    string

	ReviewMasterSMSGatewayURL      string
	ReviewMasterSMSGatewayApiToken string

	AutocabSendSMSSenderName string

	BarredTelephonePrefixFile string
}

// ReadProperties - read the properties file
func ReadProperties() {
	// fileName := "./config.properties"
	// viper.SetConfigFile(fileName)
	viper.SetConfigName("config")    // name of config file (without extension)
	viper.AddConfigPath("./config")  // look for config in the working directory config path
	viper.AddConfigPath(".")         // look for config in the working directory
	viper.AddConfigPath("../config") // Added for testing
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	Conf.DbName = viper.Get("dbname").(string)
	Conf.DbAddress = viper.Get("dbaddress").(string)
	Conf.DbPort = viper.Get("dbport").(string)
	Conf.DbUsername = viper.Get("dbUsername").(string)
	Conf.DbPassword = viper.Get("dbpassword").(string)

	Conf.PollPeriod = viper.GetInt("pollperiod")
	Conf.LastPollTime = viper.GetTime("lastpolltime")

	Conf.SendSmsURL = viper.GetString("sendsmsurl")
	Conf.SendSmsToken = viper.GetString("sendsmstoken")
	Conf.SendSmsTokenParameter = viper.GetString("sendsmstokenparameter")
	Conf.SendSmsTelephoneParameter = viper.GetString("sendsmstelephoneparameter")
	Conf.SendSmsMessageParameter = viper.GetString("sendsmsmessageparameter")
	Conf.SendSmsSuccessResponse = viper.GetString("sendsmssuccessresponse")
	Conf.SendSmsFailureResponse = viper.GetString("sendsmsfailureresponse")

	Conf.ReviewMasterSMSGatewayURL = viper.GetString("review_master_sms_gateway_url")
	Conf.ReviewMasterSMSGatewayApiToken = viper.GetString("review_master_sms_gateway_api_token")

	Conf.AutocabSendSMSSenderName = viper.GetString("autocab_send_sms_sender_name")

	Conf.BarredTelephonePrefixFile = viper.Get("barred_telephone_prefix_file").(string)
}

// UpdateProperties - update properties file
func UpdateProperties(lastPollTime time.Time) {
	viper.Set("LastPollTime", lastPollTime.Format("2006-01-02T15:04:05Z"))
	viper.WriteConfig()
}
