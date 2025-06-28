// Google Reviews UI
//
// compile for 64bit Linux production using:
// $ env GOOS=linux GOARCH=amd64 go build google_reviews_ui.go
//
// Need to add:
// 		config/config.properties file
// 		certs/server.rsa.crt - or whatever is configured in the config file
// 		certs/server.rsa.key - or whatever is configured in the config file
// 		templates/fe.gohtml - the front end template files (for vue.js), if there are more copy them all
// 		static - all this directory, it is the static content that is used for the frontend
//
// 		The frontend is produced separately using quasar which uses vue.js and a distribution is created using:
// 			$ quasar build
// 		This produces a directory dist/spa in the root directory of the quasar application.
// 		Then the whole of the dist directory is copied to this /static directory.
// 		Then /static/dist/spa/index.html content (only one line) is copied to templates/fe.gohtml between:
// 			{{ define "fe/index.tmpl" }}
// 		and
// 			{{ end }}
//
// NOTE: This needs the Google My Business code running on the same machine for the Google My Business Report to work.
//
// During testing amd using a self signed certificate could get message like:
//
// http: TLS handshake error from [::1]:50439: remote error: tls: unknown certificate
//
// In Chrome type following in URL:
//
// chrome://flags/#allow-insecure-localhost
//
// and change from Disabled to Enabled
//
// For testing with CURL (need to login first and use the returned token) use:
//
// curl -k -X POST -d '{"username":"<username>","password":"<password>"}' https://localhost:8443/login
//
// replace <username> and <password>
//
// this will return a token use this token to replace the <token> below for any other requests e.g.:
//
//  curl -k -H 'Accept: application/json' -H "Authorization: Bearer <token>" 'https://localhost:8443/auth/stats?start_day=2021-09-01&end_day=2021-09-02&time_grouping=day'
//

package main

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"

	"google_reviews_ui/config"
	"google_reviews_ui/database"
	"google_reviews_ui/server"
)

func main() {
	// log.Printf("len(os.Args) = %d\n", len(os.Args))
	// test first argument is used to change certain behaviour e.g. turn off logging to file so can see output in terminal
	test := false
	if len(os.Args) > 1 && os.Args[1] == "test" {
		test = true
	}

	logFilename := ""
	if !test {
		// create log file in same directory as executable
		ex, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(filepath.Dir(ex))
		logFilename = filepath.Dir(ex) + "/google_reviews_ui.log"
		// f, err := os.OpenFile(filepath.Dir(ex)+"/google_reviews_ui.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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

	// run http server
	server.Server(logFilename)
}
