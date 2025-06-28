package server

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"send_sms_ui/config"
)

// Server - server
// The server configuration should return a perfect SSL Labs score when using correct certificates for site
func Server(config config.Config) {
	r := SetupRouter(config, "templates/**/*")

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
		Handler:      r,
		ReadTimeout:  21 * time.Second,
		WriteTimeout: 21 * time.Second,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	log.Fatal(srv.ListenAndServeTLS(config.ServerCert, config.ServerKey))
}

// SetupRouter - setup the router
func SetupRouter(config config.Config, templatesPattern string) *gin.Engine {
	// GIN release mode (remove if want debug mode) - can also be set in systemd script with
	// Environment=GIN_MODE=release
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.LoadHTMLGlob(templatesPattern)

	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		config.Username: config.UserPassword,
	}))

	authorized.GET("/sendsms", ShowSendSmsHandler(r, config.MaxTelephones))
	authorized.POST("/sendsms", SendSmsHandler(r, config.Country, config.MobilePrefix,
		config.MaxTelephones, config.SendSmsURL, config.SendSmsSuccessResponse,
		config.TelephoneParameter, config.MessageParameter, config.TokenParameter, config.Token))

	return r
}
