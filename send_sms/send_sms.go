// Send SMS
//
// compile for 64bit Linux production using:
// $ env GOOS=linux GOARCH=amd64 go build send_sms.go
//
// Need to add:
// 		config/config.properties file
// 		config/barred_telephone_prefixes.txt file (name configured in config.properties)
// 		certs/server.rsa.crt - or whatever is configured in the config file
// 		certs/server.rsa.key - or whatever is configured in the config file
//
// The barred telephone prefixes file is read on program startup so any changes to this file
// will require a restart.
//
// Useful for token generation, use Elixir iex:
// iex(1)> length = 64
// 64
// iex(2)> :crypto.strong_rand_bytes(length) |> Base.url_encode64 |> binary_part(0, length)
// "lWCNWOzaoIPGhbxvPx7BUKJioG63JIh2W-Hc8YTyvgBBNg_vkRSQUXv--fBzKyqb"
//
// Example CURL statement for testing
// curl -k -X POST -d 'token=DOkTxeI8SkxO-KRaX2YsHkZ6XJ81ln7_InNTv4p-kjXgMri_KJ1W-wmurgSMBf_s&t=07123456789&m=testing' 'https://localhost/sendsms'
//

package main

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"

	"send_sms/barred"
	"send_sms/config"
	"send_sms/rate_limiter"
	"send_sms/server"
	"send_sms/shared"
)

func main() {
	// log.Printf("len(os.Args) = %d\n", len(os.Args))
	// test first argument is used to change certain behaviour e.g. turn off logging to file so can see output in terminal
	test := false
	if len(os.Args) > 1 && os.Args[1] == "test" {
		test = true
	}

	if !test {
		// create log file in same directory as executable
		ex, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(filepath.Dir(ex))
		logFilename := filepath.Dir(ex) + "/send_sms.log"
		// f, err := os.OpenFile(filepath.Dir(ex)+"/send_sms.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		f, err := os.OpenFile(logFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		// log.SetOutput(f)
		log.SetOutput(&lumberjack.Logger{
			Filename:   logFilename,
			MaxSize:    50, // megabytes
			MaxBackups: 10,
			MaxAge:     28, //days
		})
	}

	// read config file
	config := config.ReadProperties()

	// initialise email for each gateway
	shared.Initialise(len(config.Gateways))

	// read barred telephone numbers file
	bars, _ := barred.ReadBarredFile(config.BarredTelephonePrefixFile)

	// rate limiter
	rateLimiterEnabled, rateLimiterRestrictor := rate_limiter.RateLimiter(config)

	// run http server
	server.Server(config, bars, rateLimiterEnabled, rateLimiterRestrictor)
}
