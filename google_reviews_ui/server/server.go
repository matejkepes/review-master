package server

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"

	"google_reviews_ui/client"
	"google_reviews_ui/config"
	"google_reviews_ui/shared"
)

// jwt login
type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// jwt identity key
var identityKey = "id"

// User - jwt user
type User struct {
	UserName  string
	FirstName string
	LastName  string
	Role      string
	PartnerID string
}

// google reviews frontend using vue.js
func feHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "fe/index.tmpl", gin.H{})
}

// send test parameters
type sendTestParameters struct {
	SendURL    string `form:"sendURL" json:"send_url" binding:"required"`
	HttpGet    bool   `form:"httpGet" json:"http_get"`
	Parameters string `form:"parameters" json:"parameters" binding:"required"`
}

// ErrorResult - represents an error result.
type ErrorResult struct {
	ClientID  int       `json:"client_id"`  // client id
	Frequency int       `json:"frequency"`  // frequency
	LastError time.Time `json:"last_error"` // last error
}

// send test to get response from remote send SMS service
func sendTestHandler(c *gin.Context) {
	success := true
	var errStr string
	var resp string
	var sendTestParams sendTestParameters
	if err := c.ShouldBind(&sendTestParams); err != nil {
		log.Printf("err: %v\n", err)
		errStr = fmt.Sprintf("Problem with sent parameters: %v", err)
		success = false
	} else {
		// parameters are separated by a newline
		p := regexp.MustCompile("\r?\n").Split(sendTestParams.Parameters, -1)
		params := url.Values{}
		paramSep := regexp.MustCompile("=")
		for i := 0; i < len(p); i++ {
			a := paramSep.Split(p[i], -1)
			if len(a) == 2 {
				params.Add(a[0], a[1])
			}
		}
		httpMethod := "POST"
		if sendTestParams.HttpGet {
			httpMethod = "GET"
		}
		log.Printf("sendURL: %s, httpMethod: %s, params: %v\n", sendTestParams.SendURL, httpMethod, params)
		resp = client.Send(sendTestParams.SendURL, httpMethod, params)
	}
	c.JSON(200, gin.H{
		"success":  success,
		"err":      errStr,
		"response": resp,
	})
}

// get check log info from google review server(s)
func checkLogHandler(c *gin.Context) {
	success := true
	var errStr string
	var resp string
	errorResults := make([]ErrorResult, 0)
	// get the log (review) servers are setup in config from context otherwise will panic (can't do this anyway without)
	if logServers, ok := c.MustGet("logServers").([]config.LogServer); ok {
		hoursBack := c.DefaultQuery("hours_back", "24")
		params := url.Values{}
		params.Add("hours_back", hoursBack)
		// multiple servers need to combine results
		for _, logServer := range logServers {
			params.Set("log_token", logServer.LogToken)
			resp = client.Send(logServer.URL+"/checklogs", "GET", params)
			errResults := make([]ErrorResult, 0)
			json.Unmarshal([]byte(resp), &errResults)
			eResults := errorResults
			for _, er := range errResults {
				found := false
				for i, e := range eResults {
					if e.ClientID == er.ClientID {
						found = true
						errorResults[i].Frequency = e.Frequency + er.Frequency
						if er.LastError.After(e.LastError) {
							errorResults[i].LastError = er.LastError
						}
						break
					}
				}
				if !found {
					errorResults = append(errorResults, er)
				}
			}
		}
	} else {
		log.Printf("err: getting log (review) servers from config (and or context)\n")
		errStr = "unable to communicate with servers"
		success = false
	}
	c.JSON(200, gin.H{
		"success": success,
		"err":     errStr,
		"errors":  errorResults,
	})
}

// Server - server
// The server configuration should return a perfect SSL Labs score when using correct certificates for site
func Server(logFilename string) {
	r := SetupRouter(logFilename, "templates/*")

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
		Addr:        ":" + config.Conf.ServerPort,
		Handler:     r,
		ReadTimeout: 21 * time.Second,
		// The write timeout has been extended to facilitate for the Google my business reports, if too short the request is repeated 2 more times
		// WriteTimeout: 21 * time.Second,
		WriteTimeout: 8 * time.Minute,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	log.Fatal(srv.ListenAndServeTLS(config.Conf.ServerCert, config.Conf.ServerKey))
}

// SetupRouter - setup the router
func SetupRouter(logFilename string, templatesPattern string) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.LoadHTMLGlob(templatesPattern)

	// CORS for https://localhost:8080 origins, allowing:
	// TODO: Tighten up
	// - GET, POST, PUT, DELETE, PATCH methods
	// - Origin header
	// - Credentials share
	// - Preflight requests cached for 12 hours
	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:8080"},
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
	// 	AllowHeaders:     []string{"Origin"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	AllowOriginFunc: func(origin string) bool {
	// 		return origin == "https://localhost:8080"
	// 	},
	// 	MaxAge: 12 * time.Hour,
	// }))
	r.Use(CORSMiddleware())

	// API middleware for passing database and log filename to handlers
	r.Use(APIMiddleware(logFilename))

	// reviewServerMiddleware will add the review servers to the context
	r.Use(reviewServerMiddleware())

	// the jwt middleware
	authMiddleware, err := authJWTMiddleware()

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	r.POST("/login", authMiddleware.LoginHandler)

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	auth := r.Group("/auth")
	// Refresh time can be longer than token timeout
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		// fetch clients
		auth.GET("/clients", ClientsHandler)

		// fetch simple config for a client
		auth.GET("/configsimple", SimpleGetClientHandler)
		// update simple config for a client
		auth.PUT("/configsimple", SimpleUpdateClientHandler)
		// create simple config for a client
		auth.POST("/configsimple", SimpleCreateClientHandler)

		// fetch config for a client
		auth.GET("/config", GetClientHandler)
		// update config for a client
		auth.PUT("/config", UpdateClientHandler)
		// create config for a client
		auth.POST("/config", CreateClientHandler)
		// add config for a client
		auth.POST("/configadd", CreateClientConfigHandler)
		// create config time for a client
		auth.POST("/configtime", CreateClientConfigTimeHandler)

		// fetch stats
		auth.GET("/stats", StatsHandler)
		// fetch stats from the log file
		auth.GET("/statsfromlog", StatsFromLogHandler)
		// fetch stats from the stats file
		auth.GET("/statsfromfile", StatsFromFileHandler)

		// fetch stats from stats table
		auth.GET("/statsnew", StatsNewHandler)

		// send test
		auth.POST("/sendtest", sendTestHandler)

		// check logs
		auth.GET("/checklog", checkLogHandler)

		// check no message sent for compaies for specific period
		auth.GET("/checknothingsent", CheckNothingSentHandler)

		// run Google My Business report (adhoc)
		auth.GET("/gmybusinessreport", GoogleMyBusinessReportHandler)

		// fetch user
		auth.GET("/users", UsersHandler)
		// fetch user and associated clients
		auth.GET("/user", GetUserHandler)
		// update user and associated client
		auth.PUT("/user", UpdateUserHandler)
		// create user and associated client
		auth.POST("/user", CreateUserHandler)
		// delete user and associated client
		auth.DELETE("/user", DeleteUserHandler)
	}

	// google reviews frontend using vue.js
	r.GET("/fe", feHandler)

	// serve up static content
	// r.Static("/static", "./static")
	r.Static("/js", "./static/dist/spa/js")
	r.Static("/css", "./static/dist/spa/css")
	r.StaticFile("/favicon.ico", "./static/dist/spa/icons/favicon.ico")
	r.Static("/statics", "./static/dist/spa/statics")
	r.Static("/fonts", "./static/dist/spa/fonts")
	r.Static("/icons", "./static/dist/spa/icons")
	r.Static("/favicon.ico", "./static/favicon.ico")

	return r
}

// CORSMiddleware - CORS issue with OPTIONS sent rather than HTTP method
// NOTE: This is probably not required for deployment with the dist directory added to this server
// setup outlined in the main file and using templates and static content outlined here.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		// c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// This facilitates for OPTIONS being sent before HTTP method
		if c.Request.Method == "OPTIONS" {
			// fmt.Println("OPTIONS")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

// the authorization jwt middleware
func authJWTMiddleware() (*jwt.GinJWTMiddleware, error) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       config.Conf.AuthRealm,
		Key:         []byte(config.Conf.AuthSecretKey),
		Timeout:     time.Hour * 6,
		MaxRefresh:  time.Hour * 6,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey:   v.UserName,
					"role":        v.Role,
					"displayName": v.FirstName + " " + v.LastName,
					"partnerID":   v.PartnerID,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				UserName: claims["id"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password
			// fmt.Printf("loginVals: %v\n", loginVals)

			// if userID == config.Conf.AdminUsername && password == config.Conf.AdminUserPassword {
			// 	return &User{
			// 		UserName:  userID,
			// 		LastName:  config.Conf.AdminUserLastName,
			// 		FirstName: config.Conf.AdminUserFirstName,
			// 		Role:      config.Conf.AdminRole,
			// 	}, nil
			// } else if userID == config.Conf.TestUsername && password == config.Conf.TestUserPassword {
			// 	return &User{
			// 		UserName:  userID,
			// 		LastName:  config.Conf.TestUserLastName,
			// 		FirstName: config.Conf.TestUserFirstName,
			// 		Role:      config.Conf.UserRole,
			// 	}, nil
			// }
			// log.Printf("config.Conf.Users: %v\n", config.Conf.Users)
			for i := 0; i < len(config.Conf.Users); i++ {
				if userID == config.Conf.Users[i].Username && password == config.Conf.Users[i].Password {
					return &User{
						UserName:  userID,
						LastName:  config.Conf.Users[i].LastName,
						FirstName: config.Conf.Users[i].FirstName,
						Role:      config.Conf.Users[i].Role,
						PartnerID: config.Conf.Users[i].PartnerID,
					}, nil
				}
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			// if v, ok := data.(*User); ok && v.UserName == config.Conf.AdminUsername {
			// 	return true
			// }
			if _, ok := data.(*User); ok {
				claims := jwt.ExtractClaims(c)
				// log.Printf("claims: %v\n", claims)
				if claims["role"] == config.Conf.AdminRole {
					return true
				}
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
	return authMiddleware, err
}

// APIMiddleware will add the db connection and log filename to the context
func APIMiddleware(logFilename string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(shared.LoggerFilename, logFilename)
		c.Next()
	}
}

// reviewServerMiddleware will add the review servers to the context
func reviewServerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("logServers", config.Conf.LogServers)
		c.Next()
	}
}
