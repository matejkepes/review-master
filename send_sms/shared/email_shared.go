package shared

import (
	"container/ring"
	"time"
)

// // EmailLastSent - used to check when email was last sent so do not send too frquently.
// // Initialise email last sent to a long time ago.
// var EmailLastSent = time.Now().AddDate(-10, 0, 0)

// // FailoverEmailLastSent - used to check when email was last sent so do not send too frquently for failover.
// // Initialise email last sent to a long time ago.
// var FailoverEmailLastSent = time.Now().AddDate(-10, 0, 0)

// EmailSendFreq - frequency with which email should be sent
const EmailSendFreq = time.Minute * 60

// sendSmsLastErrorsLen - length of list of times of last send SMS errors which should initiate an email
const sendSmsLastErrorsLen = 20

// SendSmsLastErrorsCheckPeriod - check period for list of times of last send SMS errors which should initiate an email
const SendSmsLastErrorsCheckPeriod = time.Minute * 60

// // SendSmsLastErrors - list of times of last send SMS errors which should initiate an email
// var SendSmsLastErrors = ring.New(sendSmsLastErrorsLen)

// // FailoverSendSmsLastErrors - list of times of last send SMS errors which should initiate an email for failover
// var FailoverSendSmsLastErrors = ring.New(sendSmsLastErrorsLen)

// func init() {
// 	// initialise last send SMS errors list to a time before check period
// 	for i := 0; i < sendSmsLastErrorsLen; i++ {
// 		SendSmsLastErrors.Value = time.Now().AddDate(-10, 0, 0)
// 		SendSmsLastErrors = SendSmsLastErrors.Next()
// 		FailoverSendSmsLastErrors.Value = time.Now().AddDate(-10, 0, 0)
// 		FailoverSendSmsLastErrors = FailoverSendSmsLastErrors.Next()
// 	}
// }

// GatewaysErrorEmailLastSent - used to check when email was last sent for each gateway so do not send too frquently.
// Should be initialised to a long time ago.
var GatewaysErrorEmailLastSent []time.Time

// GatewaysSendSmsLastErrors - list of times of last send SMS errors for each gateway which should initiate an email
var GatewaysSendSmsLastErrors []*ring.Ring

// Initialise - initialise the gateways values
func Initialise(numberOfGateways int) {
	for j := 0; j < numberOfGateways; j++ {
		// initialise last send SMS errors list for each gateway to a time before check period
		gatewaysSendSmsLastErrors := ring.New(sendSmsLastErrorsLen)
		for i := 0; i < sendSmsLastErrorsLen; i++ {
			gatewaysSendSmsLastErrors.Value = time.Now().AddDate(-10, 0, 0)
			gatewaysSendSmsLastErrors = gatewaysSendSmsLastErrors.Next()
		}
		GatewaysSendSmsLastErrors = append(GatewaysSendSmsLastErrors, gatewaysSendSmsLastErrors)
		// Initialise email last sent to a long time ago.
		GatewaysErrorEmailLastSent = append(GatewaysErrorEmailLastSent, time.Now().AddDate(-10, 0, 0))
	}
}
