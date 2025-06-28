package config

import (
	"log"

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

	ServerCert string
	ServerKey  string

	WebsiteUser     string
	WebsitePassword string

	LogToken string

	ReviewMasterSMSGatewayPairingToken string
	ReviewMasterSMSGatewayURL          string
	ReviewMasterSMSGatewayApiToken     string

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
	Conf.DbName = viper.Get("db_name").(string)
	Conf.DbAddress = viper.Get("db_address").(string)
	Conf.DbPort = viper.Get("db_port").(string)
	Conf.DbUsername = viper.Get("db_username").(string)
	Conf.DbPassword = viper.Get("db_password").(string)

	Conf.ServerCert = viper.Get("server_cert").(string)
	Conf.ServerKey = viper.Get("server_key").(string)

	Conf.WebsiteUser = viper.Get("website_user").(string)
	Conf.WebsitePassword = viper.Get("website_password").(string)

	Conf.LogToken = viper.Get("log_token").(string)

	Conf.ReviewMasterSMSGatewayPairingToken = viper.GetString("review_master_sms_gateway_pairing_token")
	Conf.ReviewMasterSMSGatewayURL = viper.GetString("review_master_sms_gateway_url")
	Conf.ReviewMasterSMSGatewayApiToken = viper.GetString("review_master_sms_gateway_api_token")

	Conf.BarredTelephonePrefixFile = viper.Get("barred_telephone_prefix_file").(string)
}
