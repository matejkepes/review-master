package server

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"send_sms/config"

	"github.com/EagleChen/restrictor"
)

// Server - Google reviews server
// The server configuration should return a perfect SSL Labs score when using correct certificates for site
func Server(config config.Config, bars []string, rateLimiterEnabled bool, rateLimiterRestrictor restrictor.Restrictor) {
	mux := http.NewServeMux()
	mux.Handle("/sendsms", SendSmsHandler(config, bars, rateLimiterEnabled, rateLimiterRestrictor))

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
		Addr:         ":" + config.ServerPort,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	log.Fatal(srv.ListenAndServeTLS(config.ServerCert, config.ServerKey))
}
