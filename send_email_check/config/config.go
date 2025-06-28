package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	SMTPServer     string
	SMTPServerPort string
	EmailPassword  string
	EmailFrom      string
	EmailTo        string
	EmailSubject   string
	EmailMsg       string
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

	config.SMTPServer = viper.Get("smtpServer").(string)
	config.SMTPServerPort = viper.Get("smtpServerPort").(string)
	config.EmailPassword = viper.Get("emailPassword").(string)
	config.EmailFrom = viper.Get("emailFrom").(string)
	config.EmailTo = viper.Get("emailTo").(string)
	config.EmailSubject = viper.Get("emailSubject").(string)
	config.EmailMsg = viper.Get("emailMsg").(string)

	return config
}
