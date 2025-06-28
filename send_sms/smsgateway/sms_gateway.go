package smsgateway

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"send_sms/shared"
)

var smsGatewayMessageTerminationChars = "\r\n"

// sms gateway - default message relative validity time (this is 10 minutes
// == (1 + 1) * 5 minutes - the + 1 is because it starts at 0)
var smsGatewayDefaultMessageRelativeValidityTime = "1"

// LoginResponse - SMS gateway login response
type LoginResponse struct {
	ServerPassword string `json:"server_password"`
	Reply          string `json:"reply"`
	ClientID       string `json:"client_id"`
	MethodReply    string `json:"method_reply"`
	ErrorCode      string `json:"error_code"`
}

// SendMsgResponse - SMS gateway send message response
type SendMsgResponse struct {
	Reply     string `json:"reply"`
	Msg       string `json:"msg"`
	Number    string `json:"number"`
	QueueType string `json:"queue_type"`
	Unicode   string `json:"unicode"`
	Validity  string `json:"validity"`
	ErrorCode string `json:"error_code"`
}

// Send - send SMS through gateway
// returns byte array for use as HTTP response
// and a boolean to indicate whether an email alert should be sent because gateway has errors
func Send(gatewayAddress string, gatewayPort string, gatewayPassword string, gatewaySocketTimeout string, tel string, msg string) ([]byte, bool) {
	// debugging
	// log.Printf("sms_gateway.Send tel: %s, msg: %s\n", tel, msg)
	socketTimeout, err := strconv.Atoi(gatewaySocketTimeout)
	if err != nil {
		socketTimeout = 5000
	}
	timeOutDuration := time.Millisecond * time.Duration(socketTimeout)
	conn, err := net.DialTimeout("tcp", gatewayAddress+":"+gatewayPort, timeOutDuration)
	if err != nil {
		// handle error
		log.Printf("error connecting to gateway (%s:%s): %s\n", gatewayAddress, gatewayPort, err)
		return shared.FailedResponse, true
	}
	// log.Printf("connection to gateway (%s:%s): %+v\n", gatewayAddress, gatewayPort, conn)
	defer conn.Close()
	// set timeout on read / write operations
	conn.SetDeadline(time.Now().Add(timeOutDuration))

	// login to gateway
	// login example:
	// {"method":"authentication", "server_password":"admin","client_id":"id1"}
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	loginStr := "{\"method\":\"authentication\", \"server_password\":\"" + gatewayPassword + "\",\"client_id\":\"id1\"}"
	_, err = rw.WriteString(loginStr + smsGatewayMessageTerminationChars)
	if err != nil {
		log.Printf("error could not send login request string to gateway (%s:%s): %s\n", gatewayAddress, gatewayPort, err)
		return shared.FailedResponse, false
	}
	err = rw.Flush()
	if err != nil {
		log.Printf("error flushing failed: %s\n", err)
		return shared.FailedResponse, false
	}
	message, err := rw.ReadString('\n')
	if err != nil {
		log.Printf("error reading login message from gateway (%s:%s): %s\n", gatewayAddress, gatewayPort, err)
		return shared.FailedResponse, false
	}
	// log.Print("message: " + message)
	// parse reply
	if !LoginSuccessful(message, gatewayAddress, gatewayPort) {
		log.Printf("error gateway (%s:%s) login response failed\n", gatewayAddress, gatewayPort)
		return shared.FailedResponse, false
	}

	// send message
	// message format example:
	// { "number":"6453298", "msg":"6B656C6C6F", "unicode":"0", "queue_type":"master", "validity":"1" }
	gsm := fmt.Sprintf("%x", UTF8ToGsm0338(msg))
	msgStr := "{\"msg\":\"" + gsm + "\",\"number\":\"" + tel + "\",\"queue_type\":\"master\",\"unicode\":\"0\",\"validity\":\"" + smsGatewayDefaultMessageRelativeValidityTime + "\"}"
	// log.Println(msgStr)
	_, err = rw.WriteString(msgStr + smsGatewayMessageTerminationChars)
	if err != nil {
		log.Printf("error could not send message request string to gateway (%s:%s): %s\n", gatewayAddress, gatewayPort, err)
		return shared.FailedResponse, false
	}
	err = rw.Flush()
	if err != nil {
		log.Printf("error flushing failed: %s\n", err)
		return shared.FailedResponse, false
	}
	message, err = rw.ReadString('\n')
	if err != nil {
		log.Printf("error reading send message from gateway (%s:%s): %s\n", gatewayAddress, gatewayPort, err)
		return shared.FailedResponse, false
	}
	// log.Print("message: " + message)
	if !SendMsgSuccessful(message, gatewayAddress, gatewayPort) {
		log.Printf("error gateway (%s:%s) sending message response failed\n", gatewayAddress, gatewayPort)
		return shared.FailedResponse, false
	}

	// debugging
	// log.Printf("sms_gateway.Send tel: %s, response: %s\n", tel, shared.SuccessResponse)

	return shared.SuccessResponse, false
}

// LoginSuccessful - check to see if successfully logged into SMS gateway
func LoginSuccessful(message string, gatewayAddress string, gatewayPort string) bool {
	var loginResponse LoginResponse
	err := json.Unmarshal(cleanResponse(message), &loginResponse)
	if err != nil {
		log.Printf("error unmarshalling login response from gateway (%s:%s): %s, err: %s\nresponse: %s\n", gatewayAddress, gatewayPort, message, err, message)
		return false
	}
	// log.Printf("login json response: %+v\n", loginResponse)
	// Normal response reply will contain ok
	// but there is another successful response an exampler of which is:
	// {"notification": "replacing old connection IP - 35.176.91.183"}
	// Therefore check for error.
	if loginResponse.Reply == "error" {
		log.Printf("error logging into gateway (%s:%s) with error_code: %s - %s\nresponse: %s\n", gatewayAddress, gatewayPort, loginResponse.ErrorCode, errorCodeMeaning(loginResponse.ErrorCode), message)
		return false
	}
	return true
}

// SendMsgSuccessful - check to see if successfully sent message to SMS gateway
func SendMsgSuccessful(message string, gatewayAddress string, gatewayPort string) bool {
	var sendMsgResponse SendMsgResponse
	err := json.Unmarshal(cleanResponse(message), &sendMsgResponse)
	if err != nil {
		log.Printf("error unmarshalling send message response from gateway (%s:%s): %s, err: %s\nresponse: %s\n", gatewayAddress, gatewayPort, message, err, message)
		return false
	}
	// log.Printf("send message json response: %+v\n", sendMsgResponse)
	// The reply can be proceeding, ok or confirmation for a successfully sent message
	// and error if there was an error sending the SMS with another field error_code
	// which will contain the error code received from the network.
	// if sendMsgResponse.Reply != "proceeding" {
	if sendMsgResponse.Reply == "error" {
		log.Printf("error sending message to gateway (%s:%s) with error_code: %s - %s\nresponse: %s\n", gatewayAddress, gatewayPort, sendMsgResponse.ErrorCode, errorCodeMeaning(sendMsgResponse.ErrorCode), message)
		return false
	}
	return true
}

// clean the SMS gateway response for further processing
func cleanResponse(message string) []byte {
	message = strings.Replace(message, "\r", "", 1)
	message = strings.Replace(message, "\n", "", 1)
	// error_code key is not quoted, check it does not include quotes and add quotes
	// this also facilitates for the case where it is already quted just incase
	if !strings.Contains(message, "\"error_code\":") {
		message = strings.Replace(message, "error_code:", "\"error_code\":", 1)
	}
	return []byte(message)
}

// SMS gateway error code meaning
func errorCodeMeaning(errorCode string) string {
	// error code is like err-n where n is a number
	switch errorCode {
	case "err-8":
		return "Operator determined barring. This cause indicates that the MS has tried to send a mobile originating short message when the MS's network operator or service provider has forbidden such transactions."
	case "err-10":
		return "Call barred. This cause indicates that the outgoing call barred service applies to the short message service for the called destination."
	case "err-21":
		return "Short message transfer rejected. This cause indicates that the equipment sending this cause does not wish to accept this short message, although it could have accepted the short message since the equipment sending this cause is neither busy nor incompatible."
	case "err-27":
		return "Destination out of service. This cause indicates that the destination indicated by the Mobile Station cannot be reached because the interface to the destination is not functioning correctly. The term \"not functioning correctly\" indicates that a signaling message was unable to be delivered to the remote user e.g., a physical layer or data link layer failure at the remote user, user equipment off-line, etc."
	case "err-28":
		return "Unidentified subscriber. This cause indicates that the subscriber is not registered in the PLMN (i.e. IMSI not known)."
	case "err-29":
		return "Facility rejected. This cause indicates that the facility requested by the Mobile Station is not supported by the PLMN."
	case "err-30":
		return "Unknown subscriber. This cause indicates that the subscriber is not registered in the HLR (i.e. IMSI or directory number is not allocated to a subscriber)."
	case "err-31":
		return "Normal unspecified. The GSM engine refused to send the message but no reason was stated. Note that this can also be the result of a message that was recently sent to the card, before a reply was received for the previous message."
	case "err-34":
		return "Module Error. Module either has no SIM, no reception, is faulty or is still handling the sending of a previous message."
	case "err-38":
		return "Network out of order. This cause indicates that the network is not functioning correctly and that the condition is likely to last a relatively long period of time e.g., immediately reattempting the short message transfer is not likely to be successful."
	case "err-41":
		return "Temporary failure. This cause indicates that the network is not functioning correctly and that the condition is not likely to last a long period of time e.g., the Mobile Station may wish to try another short message transfer attempt almost immediately."
	case "err-42":
		return "Congestion. This cause indicates that the short message service cannot be serviced because of high traffic."
	case "err-47":
		return "Resources unavailable, unspecified. This cause is used to report a resource unavailable event only when no other cause applies."
	case "err-50":
		return "Requested facility not subscribed. This cause indicates that the requested short message service could not be provided by the network because the user has not completed the necessary administrative arrangements with its supporting networks."
	case "err-69":
		return "Requested facility not implemented. This cause indicates that the network is unable to provide the requested short message service."
	case "err-81":
		return "Invalid short message transfer reference value. This cause indicates that the equipment sending this cause has received a message with a short message reference which is not currently in use on the MS-network interface."
	case "err-95":
		return "Invalid message, unspecified. This cause is used to report an invalid message event only when no other cause in the invalid message class applies."
	case "err-96":
		return "Invalid mandatory information. This cause indicates that the equipment sending this cause has received a message where a mandatory information element is missing and/or has a content error (the two cases are indistinguishable)."
	case "err-97":
		return "Message type non-existent or not implemented. This cause indicates that the equipment sending this cause has received a message with a message type it does not recognize either because this is a message not defined or defined but not implemented by the equipment sending this cause."
	case "err-98":
		return "Message not compatible with short message protocol state. This cause indicates that the equipment sending this cause has received a message such that the procedures do not indicate that this is a permissible message to receive while in the short message transfer state."
	case "err-99":
		return "Information element non-existent or not implemented. This cause indicates that the equipment sending this cause has received a message which includes information elements not recognized because the information element identifier is not defined or it is defined but not implemented by the equipment sending the cause. However, the information element is not required to be present in the message in order for the equipment sending the cause to process the message."
	case "err-102":
		return "Timer expiry. Sending failed due to a timeout, and during this time the GSM network didn't return any specific error code."
	case "err-111":
		return "Protocol error, unspecified. This cause is used to report a protocol error event only when no other cause applies."
	case "err-127":
		return "Interworking, unspecified. This cause indicates that there has been interworking with a network which does not provide causes for actions it takes thus, the precise cause for a message which is being sent cannot be ascertained."
	case "err-128":
		return "Telematic interworking not supported"
	case "err-129":
		return "Short message Type 0 not supported"
	case "err-130":
		return "Cannot replace short message"
	case "err-143":
		return "Unspecified TP-PID error"
	case "err-144":
		return "Data coding scheme (alphabet) not supported"
	case "err-145":
		return "Message class not supported"
	case "err-159":
		return "Unspecified TP-DCS error"
	case "err-160":
		return "Command cannot be auctioned"
	case "err-161":
		return "Command unsupported"
	case "err-175":
		return "Unspecified TP-Command error"
	case "err-176":
		return "TPDU not supported"
	case "err-192":
		return "SC busy"
	case "err-193":
		return "No SC subscription"
	case "err-194":
		return "SC system failure"
	case "err-195":
		return "Invalid SME address"
	case "err-196":
		return "Destination SME barred"
	case "err-197":
		return "SM Rejected-Duplicate SM"
	case "err-198":
		return "TP-VPF not supported"
	case "err-199":
		return "TP-VP not supported"
	case "err-208":
		return "D0 SIM SMS storage full"
	case "err-209":
		return "No SMS storage capability in SIM"
	case "err-210":
		return "Error in MS"
	case "err-211":
		return "Memory Capacity Exceeded"
	case "err-212":
		return "SIM Application Toolkit Busy"
	case "err-213":
		return "SIM data download error"
	case "err-224":
		return "Card reply timeout error"
	case "err-225":
		return "SIM reply timeout error"
	case "err-226":
		return "Missing SIM"
	case "err-255":
		return "Unspecified error cause"
	case "err-300":
		return "ME failure"
	case "err-301":
		return "SMS service of ME reserved"
	case "err-302":
		return "Operation not allowed"
	case "err-303":
		return "Operation not supported"
	case "err-304":
		return "Invalid PDU mode parameter"
	case "err-305":
		return "Invalid text mode parameter"
	case "err-310":
		return "SIM not inserted"
	case "err-311":
		return "SIM PIN required"
	case "err-312":
		return "PH-SIM PIN required"
	case "err-313":
		return "SIM failure"
	case "err-314":
		return "SIM busy"
	case "err-315":
		return "SIM wrong"
	case "err-316":
		return "SIM PUK required"
	case "err-317":
		return "SIM PIN2 required"
	case "err-318":
		return "SIM PUK2 required"
	case "err-320":
		return "Memory failure"
	case "err-321":
		return "Invalid memory index"
	case "err-322":
		return "Memory full"
	case "err-330":
		return "SMSC address unknown"
	case "err-331":
		return "No network service"
	case "err-332":
		return "Network timeout"
	case "err-340":
		return "NO +CNMA ACK EXPECTED"
	case "err-500":
		return "Unknown error"
	case "err-512":
		return "MM establishment failure"
	case "err-513":
		return "Lower layer failure"
	case "err-514":
		return "CP error"
	default:
		return "Default unknown error"
	}
}
