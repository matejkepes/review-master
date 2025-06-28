// Review Master Client Portal back end
//
// compile for 64bit Linux production using:
// $ env GOOS=linux GOARCH=amd64 go build rm_client_portal.go
//
// By default, runs in HTTPS mode (production)
// For local development with HTTP, use:
// $ go run rm_client_portal.go local
// or
// $ go run rm_client_portal.go dev
//
// For testing with CURL (need to login first and use the returned token):
//
// HTTPS (production - default):
// curl -k -X POST -d email=<email> -d password=<password> https://localhost:8443/login
// or json:
// curl -k -X POST -H "Content-Type: application/json" -d '{"email":"<email>","password":"<password>"}' https://localhost:8443/login
//
// HTTP (local development):
// curl -X POST -d email=<email> -d password=<password> http://localhost:8443/login
// or json:
// curl -X POST -H "Content-Type: application/json" -d '{"email":"<email>","password":"<password>"}' http://localhost:8443/login
//
// replace <email> and <password>
//
// this will return a token use this token to replace the <token> below for any other requests e.g.:
//
// HTTPS: curl -k -H 'Accept: application/json' -H "Authorization: Bearer <token>" 'https://localhost:8443/auth/userstats?start_day=2024-11-01&end_day=2024-11-12&time_grouping=day'
// HTTP:  curl -H 'Accept: application/json' -H "Authorization: Bearer <token>" 'http://localhost:8443/auth/userstats?start_day=2024-11-01&end_day=2024-11-12&time_grouping=day'
//
// NOTE: decoding the jwt token can be done using jq (install on MacOS using: brew install jq)
// Decode the token with jq using:
// echo <token> | jq -R 'split(".") | .[0],.[1] | @base64d | fromjson'
//
// -------------------------
//
// This project calls Google My Business which requires credentials.
// These are acquired from logging into the Google API Console: https://console.developers.google.com/apis/credentials
// First select the correct project near the top of the page (the project_id can be got from here).
// The create a OAuth 2.0 Client ID using the CREATE CREDENTIATIALS option, selecting OAauth client ID and
// selecting Desktop app (IMPORTANT) as the Application type.
//
// NOTE: The credentials.json file can be downloaded from the Google API Console under credentials menu.
//
// The credentials.json file is placed in the root directory.
// To generate a new token.json file use the function in the rm_client_portal_test.go file
//
// When running the tests in google_my_business_api_test.go it is necessary to put the credentials.json file in the same directory.
//
// -------------------------
//

package main

import (
	"log"
	"os"
	"path/filepath"
	"rm_client_portal/config"
	"rm_client_portal/database"
	"rm_client_portal/google_my_business_api"
	"rm_client_portal/server"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gopkg.in/natefinch/lumberjack.v2"
)

const googleBusinessManageAuthAPIURL = "https://www.googleapis.com/auth/business.manage"

// GetGoogleCredentials - get the Google API credentials / config
func GetGoogleCredentials() *oauth2.Config {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, googleBusinessManageAuthAPIURL)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	return config
}

// Init - The functionality of main() has been moved here to make it possible to run the this in debug mode
func Init() {
	// test first argument is used to change certain behaviour e.g. turn off logging to file so can see output in terminal
	test := false
	// Default to HTTPS (production mode), only use HTTP for local development when explicitly requested
	local := false
	if len(os.Args) > 1 {
		if os.Args[1] == "test" {
			test = true
		} else if os.Args[1] == "local" || os.Args[1] == "dev" {
			local = true
			log.Println("Running in local development mode (HTTP)")
		}
	}

	if !local {
		log.Println("Running in production mode (HTTPS)")
	}

	logFilename := ""
	if !test {
		// create log file in same directory as executable
		ex, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(filepath.Dir(ex))
		logFilename = filepath.Dir(ex) + "/rm_client_portal.log"
		// f, err := os.OpenFile(filepath.Dir(ex)+"/rm_client_portal.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		f, err := os.OpenFile(logFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		// log.SetOutput(f)
		log.SetOutput(&lumberjack.Logger{
			Filename:   logFilename,
			MaxSize:    10, // megabytes
			MaxBackups: 3,
			MaxAge:     28, //days
		})
	}

	// read config file
	config.ReadProperties()

	// database
	database.OpenDB(config.Conf.DbName, config.Conf.DbAddress, config.Conf.DbPort, config.Conf.DbUsername, config.Conf.DbPassword)

	// google my business
	config := GetGoogleCredentials()
	// get the client to call the Google APIs
	google_my_business_api.SetClient(config)

	// run http server
	server.Server(logFilename, local)
}

func main() {
	Init()
}
