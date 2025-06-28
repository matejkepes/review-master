package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

// Config - config
type Config struct {
	DbName     string
	DbAddress  string
	DbPort     string
	DbUsername string
	DbPassword string

	ReviewsNotBeforeDays int

	SMTPServer     string
	SMTPServerPort string
	EmailPassword  string
	EmailFrom      string
	// EmailTo        string
	EmailSubject                               string
	EmailReportTo                              string
	EmailCsvReportSubject                      string
	EmailNameOrPostalCodeNotFoundReportSubject string

	// AI Service Configuration
	OpenAIAPIKey string

	// SendGrid Configuration
	SendGridAPIKey     string
	SendGridTemplateID string
}

// ReadProperties - read the properties file
func ReadProperties() Config {
	// fileName := "./config.properties"
	// viper.SetConfigFile(fileName)
	viper.SetConfigName("config")   // name of config file (without extension)
	viper.AddConfigPath("./config") // look for config in the working directory config path
	viper.AddConfigPath(".")        // look for config in the working directory
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	config.DbName = viper.Get("db_name").(string)
	config.DbAddress = viper.Get("db_address").(string)
	config.DbPort = viper.Get("db_port").(string)
	config.DbUsername = viper.Get("db_username").(string)
	config.DbPassword = viper.Get("db_password").(string)

	config.ReviewsNotBeforeDays = viper.GetInt("reviews_not_before_days")

	config.SMTPServer = viper.Get("smtp_server").(string)
	config.SMTPServerPort = viper.Get("smtp_server_port").(string)
	config.EmailPassword = viper.Get("email_password").(string)
	config.EmailFrom = viper.Get("email_from").(string)
	// config.EmailTo = viper.Get("email_to").(string)
	config.EmailSubject = viper.Get("email_subject").(string)
	config.EmailReportTo = viper.Get("email_report_to").(string)
	config.EmailCsvReportSubject = viper.Get("email_csv_report_subject").(string)
	config.EmailNameOrPostalCodeNotFoundReportSubject = viper.Get("email_name_or_postal_code_not_found_report_subject").(string)

	// Try to get OpenAI API key from config file, fallback to environment variable
	if viper.IsSet("openai_api_key") {
		config.OpenAIAPIKey = viper.Get("openai_api_key").(string)
	} else {
		// Fallback to environment variable
		config.OpenAIAPIKey = os.Getenv("OPENAI_API_KEY")
	}

	// Try to get SendGrid API key from config file, fallback to environment variable
	if viper.IsSet("sendgrid_api_key") {
		config.SendGridAPIKey = viper.Get("sendgrid_api_key").(string)
	} else {
		// Fallback to environment variable
		config.SendGridAPIKey = os.Getenv("SENDGRID_API_KEY")
	}

	// Try to get SendGrid template ID from config file, fallback to environment variable
	if viper.IsSet("sendgrid_template_id") {
		config.SendGridTemplateID = viper.Get("sendgrid_template_id").(string)
	} else {
		// Fallback to environment variable
		config.SendGridTemplateID = os.Getenv("SENDGRID_TEMPLATE_ID")
	}

	return config
}
