package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config - config
type Config struct {
	ServerPort string
	ServerCert string
	ServerKey  string

	LogFileName string

	WebsiteUser     string
	WebsitePassword string
}

// ReadProperties - read the properties file
func ReadProperties() Config {
	// fileName := "./config.properties"
	// viper.SetConfigFile(fileName)
	viper.SetConfigName("config")   // name of config file (without extension)
	viper.AddConfigPath("./config") // look for config in the working directory config path
	viper.AddConfigPath(".")        // look for config in the working directory
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	config.ServerPort = viper.Get("serverport").(string)
	config.ServerCert = viper.Get("servercert").(string)
	config.ServerKey = viper.Get("serverkey").(string)

	config.LogFileName = viper.Get("logfilename").(string)

	config.WebsiteUser = viper.Get("websiteuser").(string)
	config.WebsitePassword = viper.Get("websitepassword").(string)

	return config
}
