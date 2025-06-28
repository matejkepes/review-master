package server

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/dongri/phonenumber"
	"github.com/gin-gonic/gin"

	"send_sms_ui/client"
)

type sendSmsForm struct {
	Tels string `form:"tels"`
	Msg  string `form:"msg"`
}

// ShowSendSmsHandler - show send SMS form
func ShowSendSmsHandler(r *gin.Engine, maxTelephones string) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.HTML(http.StatusOK, "send_sms/form.tmpl", gin.H{
			"title":   "Send SMS",
			"content": "Maximum telephone numbers that can be sent is: " + maxTelephones,
		})
	}
	return gin.HandlerFunc(fn)
}

// SendSmsHandler - send SMS form
func SendSmsHandler(r *gin.Engine, country string, mobilePrefix string, maxTelephones string,
	sendSmsURL string, sendSmsSuccessResponse string, telephoneParameter string,
	messageParameter string, tokenParameter string, token string) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var f sendSmsForm
		c.ShouldBind(&f)

		// do not send more than maximum configured
		maxTels, err := strconv.Atoi(maxTelephones)
		if err != nil {
			maxTels = 20
		}

		// check telephone numbers (only send mobile telephone numbers)
		tels := strings.Split(f.Tels, "\n")
		space := regexp.MustCompile(`\s+`)
		// use map to unique numbers
		phoneMap := make(map[string]bool)
		for _, tel := range tels {
			tel := space.ReplaceAllString(tel, " ")
			t := phonenumber.Parse(tel, country)
			if t != "" && !phoneMap[t] && strings.HasPrefix(t, mobilePrefix) {
				phoneMap[t] = true
				if len(phoneMap) == maxTels {
					break
				}
			}
		}

		status := ""

		if len(phoneMap) == 0 {
			status = "No correctly configured telephone numbers to send, telephones: " + f.Tels
		}

		// check message
		msg := strings.Trim(f.Msg, " ")
		if len(status) == 0 {
			if len(msg) < 1 {
				status = "No message to send"
			}
		}

		if len(status) == 0 {
			// send
			params := url.Values{}
			params.Set(tokenParameter, token)
			params.Set(messageParameter, msg)
			for phone := range phoneMap {
				params.Set(telephoneParameter, phone)
				resp := client.Send(sendSmsURL, params)
				log.Printf("resp: %v\n", resp)
				// check response
				if strings.HasPrefix(strings.Trim(resp, " "), sendSmsSuccessResponse) {
					status = status + "<br />" + phone + " - sent"
				} else {
					status = status + "<br />" + phone + " - NOT sent"
				}
			}
		}

		c.HTML(http.StatusOK, "send_sms/output.tmpl", gin.H{
			"title":   "Send SMS",
			"content": "Any telephone numbers not listed did not get sent",
			"status":  template.HTML("status: " + status),
			"message": "message: " + msg,
		})
	}
	return gin.HandlerFunc(fn)
}
