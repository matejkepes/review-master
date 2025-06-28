package config

import (
	"log"

	"github.com/spf13/viper"
)

var Conf Config

// Config - config
type Config struct {
	ServerPort string

	AuthRealm     string
	AuthSecretKey string

	ServerCert string
	ServerKey  string

	RemoteFrontendURL string

	AdminRole string
	UserRole  string

	DbName     string
	DbAddress  string
	DbPort     string
	DbUsername string
	DbPassword string
}

// ReadProperties - read the properties file
func ReadProperties() {
	viper.SetConfigName("config")    // name of config file (without extension)
	viper.AddConfigPath("./config")  // look for config in the working directory config path
	viper.AddConfigPath(".")         // look for config in the working directory
	viper.AddConfigPath("../config") // Added for testing
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	Conf.ServerPort = viper.Get("server_port").(string)
	Conf.AuthRealm = viper.Get("auth_realm").(string)
	Conf.AuthSecretKey = viper.Get("auth_secret_key").(string)
	Conf.ServerCert = viper.Get("server_cert").(string)
	Conf.ServerKey = viper.Get("server_key").(string)

	Conf.RemoteFrontendURL = viper.Get("remote_frontend_url").(string)

	Conf.AdminRole = viper.Get("admin_role").(string)
	Conf.UserRole = viper.Get("user_role").(string)

	Conf.DbName = viper.Get("db_name").(string)
	Conf.DbAddress = viper.Get("db_address").(string)
	Conf.DbPort = viper.Get("db_port").(string)
	Conf.DbUsername = viper.Get("db_username").(string)
	Conf.DbPassword = viper.Get("db_password").(string)
}
