// Google Reviews Autocab - Autocab do not have iCabbi hooks equivalent so polls each server.
// It determines whether to send a message to the telephone based on certain criteria.
// The responsibility of sending the message is passed to the configured service for this task.
//
// compile for 64bit Linux production using:
// $ env GOOS=linux GOARCH=amd64 go build google_reviews_autocab.go
//
// Need to add:
// 		config/config.properties file
// 		config/barred_telephone_prefixes.txt file (name configured in config.properties)
//
// The barred telephone prefixes file is read on program startup so any changes to this file
// will require a restart.
// The barred telephone prefixes must include the country prefix.
//

package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"

	"google_reviews_autocab/barred"
	"google_reviews_autocab/config"
	"google_reviews_autocab/database"
	"google_reviews_autocab/process"
	"google_reviews_autocab/utils"
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
		logFilename = filepath.Dir(ex) + "/google_reviews_autocab.log"
		// f, err := os.OpenFile(filepath.Dir(ex)+"/google_reviews_autocab.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		f, err := os.OpenFile(logFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		// log.SetOutput(f)
		log.SetOutput(&lumberjack.Logger{
			Filename:   logFilename,
			MaxSize:    20, // megabytes
			MaxBackups: 10,
			MaxAge:     2, //days
		})
	}

	// read config file
	config.ReadProperties()
	// log.Printf("poll period: %d\n", conf.PollPeriod)
	// log.Printf("last poll time: %v+\n", conf.LastPollTime)

	// read barred telephone numbers file
	bars, err := barred.ReadBarredFile(config.Conf.BarredTelephonePrefixFile)
	if err != nil {
		log.Printf("Error: %+v, reading barred telephone prefix file: %s\n", err, config.Conf.BarredTelephonePrefixFile)
	} else {
		process.Bars = bars
	}

	// database
	database.OpenDB(config.Conf.DbName, config.Conf.DbAddress, config.Conf.DbPort, config.Conf.DbUsername, config.Conf.DbPassword)

	// set the Review Master SMS Gateway master queue ID
	database.SetReviewMasterSMSGatewayMasterQueueID()

	// initialise last poll time, read from the config file
	// NOTE: If this is more than 24 hours ago then Autocab will error so will need to make sure it is
	// less than 24 hours at start up, although the following poll it will catch up
	lastPollTime := config.Conf.LastPollTime

	// infinite loop
	for {
		// use this time in the to time for request for the archive bookings
		// use last poll time from config properties for the from time for request for the archive bookings
		// store this time as last poll time in config properties when complete
		// It does not actually matter if there is an overlap with a previous request as the checks will prevent sending an SMS to the same number
		// startPollTime := time.Now()
		// time in UTC
		startPollTime := utils.ConvertToTimeZone(time.Now(), "UTC")
		log.Printf("startPollTime: %+v\n", startPollTime)

		log.Printf("lastPollTime: %+v\n", lastPollTime)

		// get and process archive bookings
		process.PollAutocab(lastPollTime, startPollTime)

		// store last poll time to config properties
		config.UpdateProperties(startPollTime)

		// reset last poll time to start time
		lastPollTime = startPollTime

		// wait for: poll period
		time.Sleep(time.Duration(config.Conf.PollPeriod) * time.Second)
	}
}
