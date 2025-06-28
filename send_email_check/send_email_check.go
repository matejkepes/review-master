// Send email check
//
// This is used to check that the email app is still able to send emails.
// This is necessary because Google mail will automatically turn the setting
// which is turned on for less secure app off if ot's not used.
//
// This will be run as a cron job in crontab
// It will just send an email to stop Google turning the allow sending emails
// from less secure app off
//
// To allow sending emails from apps it is necessary to turn on less secure app
// access on Google see: https://support.google.com/accounts/answer/6010255
// NOTE: Google will automatically turn this setting off if it’s not being used.
// This is quite useful also: https://gist.github.com/jpillora/cb46d183eca0710d909a
//
// compile for 64bit Linux production using:
// $ env GOOS=linux GOARCH=amd64 go build send_email_check.go
//
// Need to add:
// 		config/config.properties file
//

package main

import (
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/natefinch/lumberjack.v2"

	"send_email_check/config"
)

// send - send email
// returns true if successfully sent
func send(smtpServer string, smtpServerPort string, password string, from string, to string, subject string, msg string) bool {
	// Set up authentication information.
	auth := smtp.PlainAuth("", from, password, smtpServer)

	// Connect to the server, authenticate, set the sender and recipient
	// and send the email all in one step.

	// check recipient
	// sendTo := strings.Split(to, ",")
	toSplit := regexp.MustCompile(` *, *`)
	sendTo := toSplit.Split(to, -1)
	// check that sendTo is not empty (could be array of empty strings)
	sendToEmpty := true
	for i := 0; i < len(sendTo); i++ {
		if len(sendTo[i]) > 0 {
			sendToEmpty = false
			break
		}
	}
	if sendToEmpty {
		return true
	}

	message := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		msg + "\r\n")

	err := smtp.SendMail(smtpServer+":"+smtpServerPort, auth, from, sendTo, message)
	if err != nil {
		log.Printf("error sending email: %s\n", err)
		return false
	}
	return true
}

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
		logFilename := filepath.Dir(ex) + "/send_email_check.log"
		// f, err := os.OpenFile(filepath.Dir(ex)+"/send_email_check.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		f, err := os.OpenFile(logFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		// log.SetOutput(f)
		log.SetOutput(&lumberjack.Logger{
			Filename:   logFilename,
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28, //days
		})
	}

	// read config file
	config := config.ReadProperties()

	// send email
	send(config.SMTPServer, config.SMTPServerPort, config.EmailPassword, config.EmailFrom, config.EmailTo, config.EmailSubject, config.EmailMsg)
}
