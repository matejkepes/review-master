package server

import (
	"crypto/tls"
	"log"
	"net/http"
	"rm_client_portal/config"
	"rm_client_portal/database"
	"shared_templates"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

// jwt login
type login struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// jwt identity key
var identityKey = "id"

// User -jwt user
type User struct {
	Email string
	Role  string
}

// Server - server
func Server(logFilename string, local bool) {
	r := SetupRouter(logFilename)

	srv := &http.Server{
		Addr:        ":" + config.Conf.ServerPort,
		Handler:     r,
		ReadTimeout: 21 * time.Second,
		// The write timeout has been extended to facilitate for the Google my business reports, if too short the request is repeated 2 more times
		// WriteTimeout: 21 * time.Second,
		WriteTimeout: 8 * time.Minute,
	}

	if local {
		// Run HTTP server for local development
		log.Printf("Starting HTTP server on port %s for local development", config.Conf.ServerPort)
		log.Fatal(srv.ListenAndServe())
	} else {
		// Run HTTPS server for production
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
		srv.TLSConfig = cfg
		srv.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler))
		log.Printf("Starting HTTPS server on port %s for production", config.Conf.ServerPort)
		log.Fatal(srv.ListenAndServeTLS(config.Conf.ServerCert, config.Conf.ServerKey))
	}
}

// SetupRouter - setup the router
func SetupRouter(logFilename string) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS for https://localhost:8080 origins, allowing:
	r.Use(CORSMiddleware())

	// API middleware for passing log filename to handlers
	r.Use(APIMiddleware(logFilename))

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

		// fetch user stats from stats table
		auth.GET("/userstats", StatsUserHandler)

		// fetch reviews and insights from Google
		auth.GET("/reviews", ReportOnReviewsAndInsights)

		// reports endpoints
		auth.GET("/reports", ReportsListHandler)
		auth.GET("/reports/:id/html", ReportHTMLHandler)
	}

	return r
}

// CORSMiddleware - CORS issue with OPTIONS sent rather than HTTP method
// NOTE: This is required since the frontend resides on a separate server.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Allow both HTTP and HTTPS localhost origins for development
		allowedOrigins := []string{
			config.Conf.RemoteFrontendURL,
			"http://localhost:9000",
			"https://localhost:9000",
			"http://localhost:8080",
			"https://localhost:8080",
		}

		originAllowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				originAllowed = true
				break
			}
		}

		if originAllowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

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
					identityKey: v.Email,
					"role":      config.Conf.UserRole,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				Email: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Email
			password := loginVals.Password
			// fmt.Printf("loginVals: %v\n", loginVals)

			e := database.GetUser(userID, password)
			if e != "" {
				return &User{
					Email: userID,
					Role:  config.Conf.UserRole,
				}, nil
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
				if claims["role"] == config.Conf.UserRole {
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

// APIMiddleware will add the log filename to the context
func APIMiddleware(logFilename string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(shared_templates.LoggerFilename, logFilename)
		c.Next()
	}
}
