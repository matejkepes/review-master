package email

import (
	"container/ring"
	"log"
	"net/smtp"
	"regexp"
	"time"

	"send_sms/shared"
)

// Send - send email
// returns true if successfully sent
func Send(smtpServer string, smtpServerPort string, password string, from string, to string, subject string, msg string) bool {
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

// CheckSend - check whether to send email.
// This is called when an error occurs during sending SMS that should initiate an email.
func CheckSend(sendSmsLastErrors *ring.Ring, emailLastSent *time.Time) bool {
	sendSms := true

	// only send if all list error times are after the check time period
	oldestValue := time.Now()
	oldestIndex := 0
	checkPeriod := time.Now().Add(-shared.SendSmsLastErrorsCheckPeriod)
	// log.Printf("checkPeriod = %+v\n", checkPeriod)
	for i := 0; i < sendSmsLastErrors.Len(); i++ {
		// log.Printf("oldestValue = %+v, sendSmsLastErrors.Value = %+v\n", oldestValue, sendSmsLastErrors.Value)
		// get oldest
		if oldestValue.After(sendSmsLastErrors.Value.(time.Time)) {
			oldestValue = sendSmsLastErrors.Value.(time.Time)
			oldestIndex = i
		}
		// all must be after check period
		if checkPeriod.After(sendSmsLastErrors.Value.(time.Time)) {
			sendSms = false
		}
		sendSmsLastErrors = sendSmsLastErrors.Next()
	}
	// log.Printf("sendSms = %t\n", sendSms)
	// log.Printf("oldestIndex = %+v, oldestValue = %+v\n", oldestIndex, oldestValue)

	// check last sent
	checkTime := emailLastSent.Add(shared.EmailSendFreq)
	if checkTime.Before(time.Now()) && sendSms {
		sendSms = true
	} else {
		sendSms = false
	}

	// Add error to list making sure it replaces the oldest
	for i := 0; i < sendSmsLastErrors.Len(); i++ {
		if i == oldestIndex {
			sendSmsLastErrors.Value = time.Now()
		}
		sendSmsLastErrors = sendSmsLastErrors.Next()
	}

	return sendSms
}
