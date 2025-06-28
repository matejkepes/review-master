// Google Reviews - Called from dispatcher hooks.
// It determines whether to send a message to the telephone based on certain criteria.
// The responsibility of sending the message is passed to the configured service for this task.
//
// compile for 64bit Linux production using:
// $ env GOOS=linux GOARCH=amd64 go build google_reviews.go
//
// Need to add:
// 		config/config.properties file
// 		config/barred_telephone_prefixes.txt file (name configured in config.properties)
// 		certs/server.rsa.crt - or whatever is configured in the config file
// 		certs/server.rsa.key - or whatever is configured in the config file
//
// The barred telephone prefixes file is read on program startup so any changes to this file
// will require a restart.
// The barred telephone prefixes must include the country prefix.
//
// Useful for database token generation, use Elixir iex:
// iex(1)> length = 64
// 64
// iex(2)> :crypto.strong_rand_bytes(length) |> Base.url_encode64 |> binary_part(0, length)
// "lWCNWOzaoIPGhbxvPx7BUKJioG63JIh2W-Hc8YTyvgBBNg_vkRSQUXv--fBzKyqb"
//
// Google Reviews Link Generator (short link) use:
// 		https://whitespark.ca/google-review-link-generator/
//
// Example CURL statements for testing:
//
// pairing code:
// curl -k -X GET -H "api-token: FjyFuBCyM11199VvPqjArxYYv1fCPG8XX2rzO4ycMkcHh6dl5oe-Ea7c11sIHfkr" 'https://localhost/rmsgpair?pairing_code=1234'
//

package main

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"

	"google_reviews/barred"
	"google_reviews/config"
	"google_reviews/database"
	"google_reviews/server"
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
		logFilename = filepath.Dir(ex) + "/google_reviews.log"
		// f, err := os.OpenFile(filepath.Dir(ex)+"/google_reviews.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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
			MaxAge:     7, //days
		})
	}

	// read config file
	config.ReadProperties()

	// read barred telephone numbers file
	bars, err := barred.ReadBarredFile(config.Conf.BarredTelephonePrefixFile)
	if err != nil {
		log.Printf("Error: %+v, reading barred telephone prefix file: %s\n", err, config.Conf.BarredTelephonePrefixFile)
	} else {
		server.Bars = bars
	}

	// database
	database.OpenDB(config.Conf.DbName, config.Conf.DbAddress, config.Conf.DbPort, config.Conf.DbUsername, config.Conf.DbPassword)

	// set the Review Master SMS Gateway master queue ID
	database.SetReviewMasterSMSGatewayMasterQueueID()

	// run http server
	server.Server(logFilename)
}
