package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestReadProperties(t *testing.T) {
	config := ReadProperties()
	if config.DbName != "google_reviews" {
		t.Fatal("Error reading DbName")
	}

	fmt.Printf("%+v\n", config)
}

func TestReadPropertiesWithOpenAIKey(t *testing.T) {
	// Save current env var
	oldAPIKey := os.Getenv("OPENAI_API_KEY")
	defer os.Setenv("OPENAI_API_KEY", oldAPIKey)

	// Create a temporary viper instance for testing
	v := viper.New()
	v.Set("db_name", "google_reviews")
	v.Set("db_address", "localhost")
	v.Set("db_port", "3306")
	v.Set("db_username", "test_user")
	v.Set("db_password", "test_password")
	v.Set("reviews_not_before_days", 7)
	v.Set("smtp_server", "smtp.example.com")
	v.Set("smtp_server_port", "587")
	v.Set("email_password", "test_password")
	v.Set("email_from", "test@example.com")
	v.Set("email_subject", "Test Subject")
	v.Set("email_report_to", "report@example.com")
	v.Set("email_csv_report_subject", "CSV Report")
	v.Set("email_name_or_postal_code_not_found_report_subject", "Not Found Report")

	// Test with key in config file
	v.Set("openai_api_key", "test-config-key")
	config := readPropertiesFromViper(v)
	if config.OpenAIAPIKey != "test-config-key" {
		t.Fatalf("Expected OpenAIAPIKey to be 'test-config-key', got '%s'", config.OpenAIAPIKey)
	}

	// Test with key in environment variable but not in config
	v.Set("openai_api_key", "")
	os.Setenv("OPENAI_API_KEY", "test-env-key")
	config = readPropertiesFromViper(v)
	if config.OpenAIAPIKey != "test-env-key" {
		t.Fatalf("Expected OpenAIAPIKey to be 'test-env-key', got '%s'", config.OpenAIAPIKey)
	}
}

func TestReadPropertiesWithSendGridConfig(t *testing.T) {
	// Save current env vars
	oldSendGridKey := os.Getenv("SENDGRID_API_KEY")
	oldSendGridTemplate := os.Getenv("SENDGRID_TEMPLATE_ID")
	defer func() {
		os.Setenv("SENDGRID_API_KEY", oldSendGridKey)
		os.Setenv("SENDGRID_TEMPLATE_ID", oldSendGridTemplate)
	}()

	// Create a temporary viper instance for testing
	v := viper.New()
	v.Set("db_name", "google_reviews")
	v.Set("db_address", "localhost")
	v.Set("db_port", "3306")
	v.Set("db_username", "test_user")
	v.Set("db_password", "test_password")
	v.Set("reviews_not_before_days", 7)
	v.Set("smtp_server", "smtp.example.com")
	v.Set("smtp_server_port", "587")
	v.Set("email_password", "test_password")
	v.Set("email_from", "test@example.com")
	v.Set("email_subject", "Test Subject")
	v.Set("email_report_to", "report@example.com")
	v.Set("email_csv_report_subject", "CSV Report")
	v.Set("email_name_or_postal_code_not_found_report_subject", "Not Found Report")

	// Test with keys in config file
	v.Set("sendgrid_api_key", "test-sendgrid-key")
	v.Set("sendgrid_template_id", "test-template-id")
	config := readPropertiesFromViper(v)
	if config.SendGridAPIKey != "test-sendgrid-key" {
		t.Fatalf("Expected SendGridAPIKey to be 'test-sendgrid-key', got '%s'", config.SendGridAPIKey)
	}
	if config.SendGridTemplateID != "test-template-id" {
		t.Fatalf("Expected SendGridTemplateID to be 'test-template-id', got '%s'", config.SendGridTemplateID)
	}

	// Test with keys in environment variables but not in config
	v.Set("sendgrid_api_key", "")
	v.Set("sendgrid_template_id", "")
	os.Setenv("SENDGRID_API_KEY", "env-sendgrid-key")
	os.Setenv("SENDGRID_TEMPLATE_ID", "env-template-id")
	config = readPropertiesFromViper(v)
	if config.SendGridAPIKey != "env-sendgrid-key" {
		t.Fatalf("Expected SendGridAPIKey to be 'env-sendgrid-key', got '%s'", config.SendGridAPIKey)
	}
	if config.SendGridTemplateID != "env-template-id" {
		t.Fatalf("Expected SendGridTemplateID to be 'env-template-id', got '%s'", config.SendGridTemplateID)
	}
}

// Helper function to read properties from a specific viper instance
func readPropertiesFromViper(v *viper.Viper) Config {
	var config Config
	config.DbName = v.GetString("db_name")
	config.DbAddress = v.GetString("db_address")
	config.DbPort = v.GetString("db_port")
	config.DbUsername = v.GetString("db_username")
	config.DbPassword = v.GetString("db_password")

	config.ReviewsNotBeforeDays = v.GetInt("reviews_not_before_days")

	config.SMTPServer = v.GetString("smtp_server")
	config.SMTPServerPort = v.GetString("smtp_server_port")
	config.EmailPassword = v.GetString("email_password")
	config.EmailFrom = v.GetString("email_from")
	config.EmailSubject = v.GetString("email_subject")
	config.EmailReportTo = v.GetString("email_report_to")
	config.EmailCsvReportSubject = v.GetString("email_csv_report_subject")
	config.EmailNameOrPostalCodeNotFoundReportSubject = v.GetString("email_name_or_postal_code_not_found_report_subject")

	// Try to get OpenAI API key from config file, fallback to environment variable
	if v.IsSet("openai_api_key") && v.GetString("openai_api_key") != "" {
		config.OpenAIAPIKey = v.GetString("openai_api_key")
	} else {
		// Fallback to environment variable
		config.OpenAIAPIKey = os.Getenv("OPENAI_API_KEY")
	}

	// Try to get SendGrid API key from config file, fallback to environment variable
	if v.IsSet("sendgrid_api_key") && v.GetString("sendgrid_api_key") != "" {
		config.SendGridAPIKey = v.GetString("sendgrid_api_key")
	} else {
		// Fallback to environment variable
		config.SendGridAPIKey = os.Getenv("SENDGRID_API_KEY")
	}

	// Try to get SendGrid template ID from config file, fallback to environment variable
	if v.IsSet("sendgrid_template_id") && v.GetString("sendgrid_template_id") != "" {
		config.SendGridTemplateID = v.GetString("sendgrid_template_id")
	} else {
		// Fallback to environment variable
		config.SendGridTemplateID = os.Getenv("SENDGRID_TEMPLATE_ID")
	}

	return config
}

func TestLoadSystemConfig(t *testing.T) {
	// Create a test config with a known API key
	testConfig := Config{
		DbName:             "test-db",
		DbAddress:          "test-host",
		DbPort:             "3306",
		DbUsername:         "test-user",
		DbPassword:         "test-pass",
		OpenAIAPIKey:       "test-system-config-key",
		EmailFrom:          "test@example.com",
		SendGridAPIKey:     "test-sendgrid-key",
		SendGridTemplateID: "test-template-id",
	}

	// Load system config with our test properties
	config, err := LoadSystemConfigWithProps(testConfig)
	if err != nil {
		t.Fatalf("Failed to load system config: %v", err)
	}

	// Verify the OpenAI API key was loaded correctly
	if config.AI.OpenAI.APIKey != "test-system-config-key" {
		t.Fatalf("Expected OpenAI API key to be 'test-system-config-key', got '%s'", config.AI.OpenAI.APIKey)
	}

	// Verify the SendGrid configuration was loaded correctly
	if config.Email.SendGrid.APIKey != "test-sendgrid-key" {
		t.Fatalf("Expected SendGrid API key to be 'test-sendgrid-key', got '%s'", config.Email.SendGrid.APIKey)
	}

	if config.Email.SendGrid.TemplateID != "test-template-id" {
		t.Fatalf("Expected SendGrid template ID to be 'test-template-id', got '%s'", config.Email.SendGrid.TemplateID)
	}
}

func TestSystemConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  SystemConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: SystemConfig{
				AI: AIConfig{
					OpenAI: OpenAIConfig{
						APIKey:      "test-key",
						Model:       "gpt-4",
						MaxTokens:   150,
						Temperature: 0.7,
					},
				},
				Email: EmailConfig{
					SendGrid: SendGridConfig{
						APIKey:     "test-sendgrid-key",
						TemplateID: "test-template-id",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid temperature",
			config: SystemConfig{
				AI: AIConfig{
					OpenAI: OpenAIConfig{
						APIKey:      "test-key",
						Temperature: 1.5,
					},
				},
				Email: EmailConfig{
					SendGrid: SendGridConfig{
						APIKey:     "test-sendgrid-key",
						TemplateID: "test-template-id",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid max tokens",
			config: SystemConfig{
				AI: AIConfig{
					OpenAI: OpenAIConfig{
						APIKey:    "test-key",
						MaxTokens: -1,
					},
				},
				Email: EmailConfig{
					SendGrid: SendGridConfig{
						APIKey:     "test-sendgrid-key",
						TemplateID: "test-template-id",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing sendgrid api key",
			config: SystemConfig{
				AI: AIConfig{
					OpenAI: OpenAIConfig{
						APIKey:      "test-key",
						Model:       "gpt-4",
						MaxTokens:   150,
						Temperature: 0.7,
					},
				},
				Email: EmailConfig{
					SendGrid: SendGridConfig{
						APIKey:     "",
						TemplateID: "test-template-id",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing sendgrid template id",
			config: SystemConfig{
				AI: AIConfig{
					OpenAI: OpenAIConfig{
						APIKey:      "test-key",
						Model:       "gpt-4",
						MaxTokens:   150,
						Temperature: 0.7,
					},
				},
				Email: EmailConfig{
					SendGrid: SendGridConfig{
						APIKey:     "test-sendgrid-key",
						TemplateID: "",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.config.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("SystemConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClientConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  ClientConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: ClientConfig{
				Style: ResponseStyle{
					MaxLength: 300,
					MinLength: 50,
				},
				Solutions: BusinessSolutions{
					AllowedActions: []string{"contact_support"},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid length range",
			config: ClientConfig{
				Style: ResponseStyle{
					MaxLength: 50,
					MinLength: 100,
				},
				Solutions: BusinessSolutions{
					AllowedActions: []string{"contact_support"},
				},
			},
			wantErr: true,
		},
		{
			name: "no allowed actions",
			config: ClientConfig{
				Style: ResponseStyle{
					MaxLength: 300,
					MinLength: 50,
				},
				Solutions: BusinessSolutions{
					AllowedActions: []string{},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.config.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("ClientConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetDefaultClientConfig(t *testing.T) {
	config := GetDefaultClientConfig()

	// Test default values
	if len(config.Solutions.AllowedActions) == 0 {
		t.Error("Expected default allowed actions, got empty slice")
	}
	if len(config.Solutions.ProhibitedActions) == 0 {
		t.Error("Expected default prohibited actions, got empty slice")
	}
	if len(config.Solutions.ContactChannels) == 0 {
		t.Error("Expected default contact channels, got empty slice")
	}
	if config.Style.Tone != "professional" {
		t.Errorf("Expected default tone 'professional', got %s", config.Style.Tone)
	}
}

func TestMergeConfigs(t *testing.T) {
	sysConfig := &SystemConfig{
		AI: AIConfig{
			OpenAI: OpenAIConfig{
				APIKey:    "test-key",
				Model:     "gpt-4",
				MaxTokens: 150,
			},
		},
	}

	clientConfig := &ClientConfig{
		ID:   1,
		Name: "Test Client",
		Style: ResponseStyle{
			Tone:      "professional",
			MaxLength: 300,
			MinLength: 50,
		},
	}

	combined := MergeConfigs(sysConfig, clientConfig)

	if combined.System.AI.OpenAI.APIKey != "test-key" {
		t.Error("System config not properly merged")
	}
	if combined.Client.ID != 1 {
		t.Error("Client config not properly merged")
	}
}
