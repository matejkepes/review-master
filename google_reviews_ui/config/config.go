package config

import (
	"encoding/json"
	"log"

	"github.com/spf13/viper"
)

var Conf Config

// Config - config
type Config struct {
	ServerPort string

	AuthRealm     string
	AuthSecretKey string

	// AdminUsername string
	// AdminUserPassword  string
	// AdminUserFirstName string
	// AdminUserLastName  string
	// TestUsername       string
	// TestUserPassword   string
	// TestUserFirstName  string
	// TestUserLastName   string

	AdminRole string
	UserRole  string

	Users []User

	ServerCert string
	ServerKey  string

	DbName     string
	DbAddress  string
	DbPort     string
	DbUsername string
	DbPassword string

	LogServers []LogServer

	GoogleMyBusinessDirectory string
}

// User - user
type User struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	PartnerID string `json:"partner_id"`
}

// LogServer - log server to interogate for logs
type LogServer struct {
	URL      string `json:"url"`
	LogToken string `json:"log_token"`
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
	Conf.ServerPort = viper.Get("server_port").(string)
	Conf.AuthRealm = viper.Get("auth_realm").(string)
	Conf.AuthSecretKey = viper.Get("auth_secret_key").(string)
	// Conf.AdminUsername = viper.Get("admin_username").(string)
	// Conf.AdminUserPassword = viper.Get("admin_user_password").(string)
	// Conf.AdminUserFirstName = viper.Get("admin_user_first_name").(string)
	// Conf.AdminUserLastName = viper.Get("admin_user_last_name").(string)
	// Conf.TestUsername = viper.Get("test_username").(string)
	// Conf.TestUserPassword = viper.Get("test_user_password").(string)
	// Conf.TestUserFirstName = viper.Get("test_user_first_name").(string)
	// Conf.TestUserLastName = viper.Get("test_user_last_name").(string)
	Conf.AdminRole = viper.Get("admin_role").(string)
	Conf.UserRole = viper.Get("user_role").(string)
	Conf.ServerCert = viper.Get("server_cert").(string)
	Conf.ServerKey = viper.Get("server_key").(string)
	Conf.DbName = viper.Get("db_name").(string)
	Conf.DbAddress = viper.Get("db_address").(string)
	Conf.DbPort = viper.Get("db_port").(string)
	Conf.DbUsername = viper.Get("db_username").(string)
	Conf.DbPassword = viper.Get("db_password").(string)

	// var users []User
	// us := viper.Get("users").(string)
	// a := regexp.MustCompile("\\}\\s*\\,\\s*\\{").Split(us, -1)
	// for i := 0; i < len(a); i++ {
	// 	b := a[i]
	// 	if !strings.HasPrefix(b, "{") {
	// 		b = "{" + b
	// 	}
	// 	if !strings.HasSuffix(b, "}") {
	// 		b += "}"
	// 	}
	// 	// fmt.Printf("%s\n", b)
	// 	u := User{}
	// 	json.Unmarshal([]byte(b), &u)
	// 	// fmt.Printf("%v\n", u)
	// 	users = append(users, u)
	// }
	// // fmt.Printf("users: %v\n", users)
	// Conf.Users = users
	var users []User
	us := viper.Get("users").(string)
	json.Unmarshal([]byte(us), &users)
	Conf.Users = users

	logServers := make([]LogServer, 0)
	lss := viper.Get("logservers").(string)
	json.Unmarshal([]byte(lss), &logServers)
	Conf.LogServers = logServers

	Conf.GoogleMyBusinessDirectory = viper.GetString("google_my_business_directory")
}
