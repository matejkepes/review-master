package config

import (
	// "fmt"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Token - token
type Token struct {
	Value   string `json:"value"`
	Comment string `json:"comment"`
}

// Gateway - gateway
type Gateway struct {
	GatewayAddress       string `json:"gateway_address"`
	GatewayPort          string `json:"gateway_port"`
	GatewayPassword      string `json:"gateway_password"`
	GatewaySocketTimeout string `json:"gateway_socket_timeout"`
	NumberOfSims         string `json:"number_of_sims"`
	EmailSubject         string `json:"email_subject"`
	EmailMsg             string `json:"email_msg"`
}

// Config - config
type Config struct {
	// Token string
	// only require a map with key as the token value for a quick look up to check that the token is valid
	Tokens map[string]int

	// GatewayAddress       string
	// GatewayPort          string
	// GatewayPassword      string
	// GatewaySocketTimeout string

	// FailoverGatewayAddress       string
	// FailoverGatewayPort          string
	// FailoverGatewayPassword      string
	// FailoverGatewaySocketTimeout string

	TelephoneParameter string
	MessageParameter   string
	TokenParameter     string

	Country string

	BarredTelephonePrefixFile string

	ServerPort string
	ServerCert string
	ServerKey  string

	SMTPServer     string
	SMTPServerPort string
	EmailPassword  string
	EmailFrom      string
	EmailTo        string
	// EmailGatewaySubject         string
	// EmailGatewayMsg             string
	// EmailFailoverGatewaySubject string
	// EmailFailoverGatewayMsg     string

	Gateways []Gateway

	// Rate Limiter
	RateLimiterEnabled       bool
	RateLimiterWindowMinutes time.Duration
	RateLimiterUpperLimit    int
	RateLimiterBucketSpan    int
	RateLimiterIgnore        []string
}

// ReadProperties - read the properties file
func ReadProperties() Config {
	// fileName := "./config.properties"
	// viper.SetConfigFile(fileName)
	viper.SetConfigName("config")    // name of config file (without extension)
	viper.AddConfigPath("./config")  // look for config in the working directory config path
	viper.AddConfigPath(".")         // look for config in the working directory
	viper.AddConfigPath("../config") // Added for testing
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	var config Config

	// config.Token = viper.Get("token").(string)
	var tokens = make(map[string]int)
	ts := viper.Get("tokens").(string)
	// ta := regexp.MustCompile("\\}\\s*\\,\\s*\\{").Split(ts, -1)
	// for i := 0; i < len(ta); i++ {
	// 	b := ta[i]
	// 	if !strings.HasPrefix(b, "{") {
	// 		b = "{" + b
	// 	}
	// 	if !strings.HasSuffix(b, "}") {
	// 		b += "}"
	// 	}
	// 	// fmt.Printf("%s\n", b)
	// 	t := Token{}
	// 	json.Unmarshal([]byte(b), &t)
	// 	// fmt.Printf("%+v\n", t)
	// 	// set value to 1 because a key which is not found will return the zero value which for int will be 0
	// 	tokens[t.Value] = 1
	// }
	var tks []Token
	json.Unmarshal([]byte(ts), &tks)
	for i := 0; i < len(tks); i++ {
		// set value to 1 because a key which is not found will return the zero value which for int will be 0
		tokens[tks[i].Value] = 1
	}
	// fmt.Printf("tokens: %+v\n", tokens)
	config.Tokens = tokens

	// config.GatewayAddress = viper.Get("gateway_address").(string)
	// config.GatewayPort = viper.Get("gateway_port").(string)
	// config.GatewayPassword = viper.Get("gateway_password").(string)
	// config.GatewaySocketTimeout = viper.Get("gateway_socket_timeout").(string)
	// config.FailoverGatewayAddress = viper.Get("failover_gateway_address").(string)
	// config.FailoverGatewayPort = viper.Get("failover_gateway_port").(string)
	// config.FailoverGatewayPassword = viper.Get("failover_gateway_password").(string)
	// config.FailoverGatewaySocketTimeout = viper.Get("failover_gateway_socket_timeout").(string)
	config.TelephoneParameter = viper.Get("telephone_parameter").(string)
	config.MessageParameter = viper.Get("message_parameter").(string)
	config.TokenParameter = viper.Get("token_parameter").(string)
	config.Country = viper.Get("country").(string)
	config.BarredTelephonePrefixFile = viper.Get("barred_telephone_prefix_file").(string)
	config.ServerPort = viper.Get("server_port").(string)
	config.ServerCert = viper.Get("server_cert").(string)
	config.ServerKey = viper.Get("server_key").(string)
	config.SMTPServer = viper.Get("smtp_server").(string)
	config.SMTPServerPort = viper.Get("smtp_server_port").(string)
	config.EmailPassword = viper.Get("email_password").(string)
	config.EmailFrom = viper.Get("email_from").(string)
	config.EmailTo = viper.Get("email_to").(string)
	// config.EmailGatewaySubject = viper.Get("email_gateway_subject").(string)
	// config.EmailGatewayMsg = viper.Get("email_gateway_msg").(string)
	// config.EmailFailoverGatewaySubject = viper.Get("email_failover_gateway_subject").(string)
	// config.EmailFailoverGatewayMsg = viper.Get("email_failover_gateway_msg").(string)

	var gateways []Gateway
	gs := viper.Get("gateways").(string)
	// a := regexp.MustCompile("\\}\\s*\\,\\s*\\{").Split(gs, -1)
	// for i := 0; i < len(a); i++ {
	// 	b := a[i]
	// 	if !strings.HasPrefix(b, "{") {
	// 		b = "{" + b
	// 	}
	// 	if !strings.HasSuffix(b, "}") {
	// 		b += "}"
	// 	}
	// 	// fmt.Printf("%s\n", b)
	// 	g := Gateway{}
	// 	json.Unmarshal([]byte(b), &g)
	// 	// fmt.Printf("%+v\n", g)
	// 	gateways = append(gateways, g)
	// }
	json.Unmarshal([]byte(gs), &gateways)
	// fmt.Printf("gateways: %+v\n", gateways)
	config.Gateways = gateways

	// Rate Limiter
	config.RateLimiterEnabled = viper.GetBool("rate_limiter_enabled")
	config.RateLimiterWindowMinutes = viper.GetDuration("rate_limiter_window_minutes") * time.Minute
	config.RateLimiterUpperLimit = viper.GetInt("rate_limiter_upper_limit")
	config.RateLimiterBucketSpan = viper.GetInt("rate_limiter_bucket_span")
	// config.RateLimiterIgnore = viper.GetStringSlice("rate_limiter_ignore")
	rateLimiterIgnore := viper.GetString("rate_limiter_ignore")
	config.RateLimiterIgnore = strings.Split(rateLimiterIgnore, ",")

	return config
}
