package config

import (
	"fmt"
	"time"
)

// SystemConfig holds all system-level configurations
type SystemConfig struct {
	AI    AIConfig
	Email EmailConfig
	// Future system configs (e.g., database, server) would go here
}

// EmailConfig holds all email-related configuration
type EmailConfig struct {
	SendGrid SendGridConfig
}

// SendGridConfig holds SendGrid-specific configuration
type SendGridConfig struct {
	APIKey     string
	TemplateID string
}

// AIConfig holds all AI-related system configuration
type AIConfig struct {
	OpenAI    OpenAIConfig
	Generator GeneratorConfig
	Validator ValidatorConfig
}

// OpenAIConfig holds OpenAI-specific settings
type OpenAIConfig struct {
	APIKey      string
	Model       string
	MaxTokens   int
	Temperature float32
	MaxRetries  int
	RetryDelay  time.Duration
}

// GeneratorConfig holds base generator settings
type GeneratorConfig struct {
	BasePrompt     string
	MaxLength      int
	DefaultTimeout time.Duration
}

// ValidatorConfig holds base validator settings
type ValidatorConfig struct {
	BasePrompt      string
	DefaultMinScore float32
	Timeout         time.Duration
}

// ClientConfig holds all client-specific configurations
type ClientConfig struct {
	ID                 int
	Name               string
	Solutions          BusinessSolutions
	Style              ResponseStyle
	Validation         ValidationRules
	CreatedAt          time.Time
	UpdatedAt          time.Time
	CompanyName        string
	ContactMethod      string // e.g., "Call us at (555) 123-4567" or "Book through our app at example.com"
	AIResponsesEnabled bool   // Whether to use AI-generated responses or fall back to templates
}

// BusinessSolutions defines available solutions for customer issues
type BusinessSolutions struct {
	AllowedActions    []string // e.g., "contact_support", "speak_to_manager"
	ProhibitedActions []string // e.g., "refund", "compensation"
	ContactChannels   []ContactChannel
	CustomPolicies    []string
}

// ContactChannel represents a way to contact the business
type ContactChannel struct {
	Type    string // "email", "phone", "in_person"
	Details string // actual contact information
}

// ResponseStyle defines how responses should be formatted
type ResponseStyle struct {
	Tone        string // "professional", "friendly", "formal"
	MaxLength   int
	MinLength   int
	CustomIntro string
	CustomClose string
}

// ValidationRules defines client-specific validation requirements
type ValidationRules struct {
	MinScore        float32
	RequiredPoints  []string
	ProhibitedTerms []string
}

// CombinedConfig represents the merged system and client configurations
type CombinedConfig struct {
	System SystemConfig
	Client ClientConfig
}

// Validate checks if the system configuration is valid
func (c *SystemConfig) Validate() error {
	if err := c.AI.Validate(); err != nil {
		return fmt.Errorf("AI config validation failed: %w", err)
	}

	if err := c.Email.Validate(); err != nil {
		return fmt.Errorf("email config validation failed: %w", err)
	}

	return nil
}

// Validate checks if the email configuration is valid
func (c *EmailConfig) Validate() error {
	if c.SendGrid.APIKey == "" {
		return fmt.Errorf("SendGrid API key is required")
	}

	if c.SendGrid.TemplateID == "" {
		return fmt.Errorf("SendGrid template ID is required")
	}

	return nil
}

// Validate checks if the AI configuration is valid
func (c *AIConfig) Validate() error {
	if c.OpenAI.APIKey == "" {
		return fmt.Errorf("OpenAI API key is required")
	}
	if c.OpenAI.MaxTokens <= 0 {
		return fmt.Errorf("max tokens must be positive")
	}
	if c.OpenAI.Temperature < 0 || c.OpenAI.Temperature > 1 {
		return fmt.Errorf("temperature must be between 0 and 1")
	}
	return nil
}

// Validate checks if the client configuration is valid
func (c *ClientConfig) Validate() error {
	if c.Style.MaxLength <= c.Style.MinLength {
		return fmt.Errorf("max length must be greater than min length")
	}
	if c.Style.MaxLength <= 0 {
		return fmt.Errorf("max length must be positive")
	}
	if len(c.Solutions.AllowedActions) == 0 {
		return fmt.Errorf("at least one allowed action is required")
	}
	return nil
}
