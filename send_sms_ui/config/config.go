package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config - config
type Config struct {
	ServerPort string

	Username     string
	UserPassword string

	ServerCert string
	ServerKey  string

	Country string

	MobilePrefix string

	MaxTelephones string

	Token string

	TelephoneParameter string
	MessageParameter   string
	TokenParameter     string

	SendSmsURL string

	SendSmsSuccessResponse string
}

// ReadProperties - read the properties file
func ReadProperties() Config {
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
	var config Config
	config.ServerPort = viper.Get("serverPort").(string)
	config.Username = viper.Get("username").(string)
	config.UserPassword = viper.Get("userPassword").(string)
	config.ServerCert = viper.Get("serverCert").(string)
	config.ServerKey = viper.Get("serverKey").(string)
	config.Country = viper.Get("country").(string)
	config.MobilePrefix = viper.Get("mobilePrefix").(string)
	config.MaxTelephones = viper.Get("maxTelephones").(string)
	config.Token = viper.Get("token").(string)
	config.TelephoneParameter = viper.Get("telephoneParameter").(string)
	config.MessageParameter = viper.Get("messageParameter").(string)
	config.TokenParameter = viper.Get("tokenParameter").(string)
	config.SendSmsURL = viper.Get("sendSmsURL").(string)
	config.SendSmsSuccessResponse = viper.Get("SendSmsSuccessResponse").(string)

	return config
}
