package server

import (
	"bufio"
	"crypto/subtle"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"google_reviews/config"
)

// log viewer variables
var (
	logFileName     = ""
	websiteUser     = "djhSLgukORf88dIfkdsfk"
	websitePassword = "hbhdf77ms7gcGDJJff"
)

// Display log file
func logviewer(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintln(w, "Logs")
	// fmt.Fprintf(w, "logFileName: %s\n", logFileName)
	f, err := os.Open(logFileName) // For read access.
	if err != nil {
		log.Printf("Error, opening log file: %s to view logs with error: %v\n", logFileName, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fmt.Fprintln(w, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		log.Printf("Error, reading log file: %s to view logs with error: %v\n", logFileName, err)
	}
}

// basicAuth wraps a handler requiring HTTP basic auth for it using the given
// username and password and the specified realm, which shouldn't contain quotes.
//
// Most web browser display a dialog with something like:
//
//	The website says: "<realm>"
//
// Which is really stupid so you may want to set the realm to a message rather than
// an actual realm.
func basicAuth(handler http.HandlerFunc, realm string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || subtle.ConstantTimeCompare([]byte(user),
			[]byte(websiteUser)) != 1 || subtle.ConstantTimeCompare([]byte(pass),
			[]byte(websitePassword)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)
			w.Write([]byte("You are Unauthorized to access the application.\n"))
			return
		}
		handler(w, r)
	}
}

// Server - Google reviews server
// The server configuration should return a perfect SSL Labs score when using correct certificates for site
func Server(logFilename string) {

	// define the log file name to display the logs
	logFileName = logFilename
	// username and password to access website / view logs
	websiteUser = config.Conf.WebsiteUser
	websitePassword = config.Conf.WebsitePassword

	mux := http.NewServeMux()
	// mux.Handle("/googlereviews", googleReviewsHandler{})
	// mux.Handle("/googlereviews", http.HandlerFunc(GoogleReviewsHandler))
	// mux.HandleFunc("/googlereviews", GoogleReviewsHandler)
	mux.Handle("/googlereviews", GoogleReviewsHandler())
	// mux.HandleFunc("/googlereviews", func(w http.ResponseWriter, req *http.Request) {
	// 	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	// 	w.Write([]byte("This is the response from server.\n"))
	// })
	mux.HandleFunc("/logs", basicAuth(logviewer, "Please enter your username and password"))
	mux.Handle("/checklogs", CheckLogHandler(logFileName, config.Conf.LogToken))
	mux.Handle("/rmsgpair", ReviewMasterSMSGatewayPairingHandler(config.Conf.ReviewMasterSMSGatewayPairingToken))
	mux.Handle("/cordic", CordicHandler())
	mux.Handle("/cab9", Cab9Handler())

	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
	srv := &http.Server{
		Addr:         ":443",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	log.Fatal(srv.ListenAndServeTLS(config.Conf.ServerCert, config.Conf.ServerKey))
}
