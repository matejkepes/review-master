package server

import (
	"bytes"
	"container/ring"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"send_sms/barred"
	"send_sms/config"
	"send_sms/email"
	"send_sms/shared"
	"send_sms/smsgateway"

	"github.com/EagleChen/restrictor"
	"github.com/dongri/phonenumber"
)

// SendSmsHandler - Send SMS Handler
func SendSmsHandler(config config.Config, bars []string, rateLimiterEnabled bool, rateLimiterRestrictor restrictor.Restrictor) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			// fmt.Printf("ParseForm() err: %+v\n", err)
			log.Printf("ParseForm() err: %+v\n", err)
			w.Write(shared.FailedResponse)
			return
		}
		// debugging
		// log.Printf("SendSmsHandler request parameters: %+v\n", req.PostForm)

		// check server token
		token := req.FormValue(config.TokenParameter)
		// log.Printf("token: %s\n", token)
		// log.Printf("config.Tokens[token]: %d\n", config.Tokens[token])
		// if token != config.Token {
		if config.Tokens[token] == 0 {
			// log.Printf("token %s not found", token)
			w.Write(shared.FailedResponse)
			return
		}

		// get telephone parameter
		tel := req.FormValue(config.TelephoneParameter)
		// log.Printf("tel param: %s\n", tel)
		telephone := phonenumber.Parse(tel, config.Country)
		// log.Printf("telephone: %s\n", telephone)
		if telephone == "" {
			// log.Printf("no telephone found\n")
			w.Write(shared.FailedResponse)
			return
		}
		// check barred telephone prefixes
		if barred.CheckBarred(telephone, bars) {
			// log.Printf("telephone number is barred\n")
			w.Write(shared.FailedResponse)
			return
		}
		// rate limit telephone check
		// if rateLimiterEnabled && !stringInSlice(telephone, config.RateLimiterIgnore) && rateLimiterRestrictor.LimitReached(telephone) {
		reached, _ := rateLimiterRestrictor.LimitReached(telephone)
		if rateLimiterEnabled && !stringInSlice(telephone, config.RateLimiterIgnore) && reached {
			log.Printf("telephone number %s has been rate limited.\nrequest: %+v\n", telephone, req)
			w.Write(shared.FailedResponse)
			return
		}

		// get message parameter
		msg := req.FormValue(config.MessageParameter)
		// log.Printf("msg param: %s\n", msg)
		if msg == "" {
			// log.Printf("no message found\n")
			w.Write(shared.FailedResponse)
			return
		}
		msg = strings.TrimSpace(msg)

		// // send message
		// resp := send(config.GatewayAddress, config.GatewayPort, config.GatewayPassword, config.GatewaySocketTimeout, tel, msg,
		// 	shared.SendSmsLastErrors, &shared.EmailLastSent,
		// 	config.SMTPServer, config.SMTPServerPort, config.EmailPassword, config.EmailFrom, config.EmailTo, config.EmailGatewaySubject, config.EmailGatewayMsg)
		// // log.Printf("resp: %+v\n", resp)
		// if !bytes.Equal(resp, shared.SuccessResponse) {
		// 	resp = send(config.FailoverGatewayAddress, config.FailoverGatewayPort, config.FailoverGatewayPassword, config.FailoverGatewaySocketTimeout, tel, msg,
		// 		shared.FailoverSendSmsLastErrors, &shared.FailoverEmailLastSent,
		// 		config.SMTPServer, config.SMTPServerPort, config.EmailPassword, config.EmailFrom, config.EmailTo, config.EmailFailoverGatewaySubject, config.EmailFailoverGatewayMsg)
		// }

		// determine the gateway to send to (set ignore to -1 so do not ignore any)
		g := SendTo(config, -1)
		// send message
		resp := send(config.Gateways[g].GatewayAddress, config.Gateways[g].GatewayPort, config.Gateways[g].GatewayPassword, config.Gateways[g].GatewaySocketTimeout, tel, msg,
			shared.GatewaysSendSmsLastErrors[g], &shared.GatewaysErrorEmailLastSent[g],
			config.SMTPServer, config.SMTPServerPort, config.EmailPassword, config.EmailFrom, config.EmailTo, config.Gateways[g].EmailSubject, config.Gateways[g].EmailMsg)
		// log.Printf("resp: %+v\n", resp)
		if !bytes.Equal(resp, shared.SuccessResponse) {
			h := SendTo(config, g)
			resp = send(config.Gateways[h].GatewayAddress, config.Gateways[h].GatewayPort, config.Gateways[h].GatewayPassword, config.Gateways[h].GatewaySocketTimeout, tel, msg,
				shared.GatewaysSendSmsLastErrors[h], &shared.GatewaysErrorEmailLastSent[h],
				config.SMTPServer, config.SMTPServerPort, config.EmailPassword, config.EmailFrom, config.EmailTo, config.Gateways[h].EmailSubject, config.Gateways[h].EmailMsg)
		}

		// debugging
		// log.Printf("SendSmsHandler send for telephone %s resp: %+s\n", telephone, resp)

		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		// w.Write(shared.SuccessResponse)
		w.Write([]byte(resp))
	}

	return http.HandlerFunc(fn)
}

// send message with parameters to SMS gateway
// and send email alert if necessary
func send(gatewayAddress string, gatewayPort string, gatewayPassword string, gatewaySocketTimeout string, tel string, msg string,
	sendSmsLastErrors *ring.Ring, emailLastSent *time.Time,
	smtpServer string, smtpServerPort string, emailPassword string, emailFrom string, emailTo string, emailGatewaySubject string, emailGatewayMsg string) []byte {
	resp, sendEmail := smsgateway.Send(gatewayAddress, gatewayPort, gatewayPassword, gatewaySocketTimeout, tel, msg)
	if sendEmail {
		if email.CheckSend(sendSmsLastErrors, emailLastSent) {
			if email.Send(smtpServer, smtpServerPort, emailPassword, emailFrom, emailTo, emailGatewaySubject, emailGatewayMsg) {
				*emailLastSent = time.Now()
			}
		}
	}
	return resp
}

// total number of sims in all gateways
var totalNumberOfSims = 0

// gateway to send to position by value
var sendToGateway []int

// SendTo - determine which gateway to send to. Returns the config gateway array index
// weighted by number of sims in each gateway
// ignore will ignore the gateway
func SendTo(config config.Config, ignore int) int {
	if totalNumberOfSims == 0 {
		initialise(config)
	}
	if totalNumberOfSims == 0 {
		log.Printf("Unable to determine total number of sims using first gateway\n")
		// return the first gateway
		return 0
	}
	r := rand.Intn(totalNumberOfSims)
	// log.Printf("totalNumberOfSims = %d, sendToGateway = %+v, r = %d, ignore = %d\n", totalNumberOfSims, sendToGateway, r, ignore)
	for i := 0; i < len(sendToGateway); i++ {
		if r < sendToGateway[i] {
			if ignore == i {
				// select next gateway
				if i < len(sendToGateway)-1 {
					return i + 1
				} else if i == 0 {
					// have no option but to return the ignored gateway
					return 0
				} else {
					return i - 1
				}
			}
			return i
		}
	}
	return 0
}

// initialise total number of sims in all gateways and send to position
func initialise(config config.Config) {
	// log.Printf("initialise\n")
	for i := 0; i < len(config.Gateways); i++ {
		v, err := strconv.Atoi(config.Gateways[i].NumberOfSims)
		if err == nil {
			totalNumberOfSims += v
			sendToGateway = append(sendToGateway, totalNumberOfSims)
		}
	}
}

// check whether a string is in a slice of strings (used to check if a telephone is in the rate limiter ignore list)
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if strings.TrimSpace(b) == a {
			return true
		}
	}
	return false
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
