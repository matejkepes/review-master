package client

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

// Clients and Transports are safe for concurrent use by multiple goroutines and for efficiency
// should only be created once and re-used.
var httpClient *http.Client

const (
	maxIdleConnections    int           = 20
	idleConnTimeout       time.Duration = time.Duration(10) * time.Second
	expectContinueTimeout time.Duration = time.Duration(10) * time.Second
	timeout               time.Duration = time.Duration(10) * time.Second
)

// Reuse the connection
func createHTTPClient() *http.Client {
	tr := &http.Transport{
		MaxIdleConnsPerHost:   maxIdleConnections,
		IdleConnTimeout:       idleConnTimeout,
		ExpectContinueTimeout: expectContinueTimeout,
		// Added for sending to servers with invalid certificate
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: timeout}

	return client
}

// Send - send HTTP POST request
func Send(sendURL string, params url.Values) string {
	baseURL, err := url.Parse(sendURL)
	if err != nil {
		log.Println(err)
		return ""
	}

	baseURL.RawQuery = params.Encode()
	// log.Println(baseURL)

	req, _ := http.NewRequest("POST", baseURL.String(), nil)
	req.Header.Set("Content-Type", "text/plain")

	if httpClient == nil {
		httpClient = createHTTPClient()
	}

	resp, err := httpClient.Do(req)
	if resp != nil {
		// close the connection to reuse it
		defer resp.Body.Close()
	}
	if err != nil {
		// get here if HTTP response code is not 2xx
		// fmt.Println(resp.StatusCode)
		log.Println(err)
		return ""
	}

	txt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	return string(txt)
}
