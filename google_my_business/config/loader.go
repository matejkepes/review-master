package config

import (
	"time"
)

// LoadSystemConfig loads system configuration from environment and defaults
func LoadSystemConfig() (*SystemConfig, error) {
	// Get configuration from properties file
	props := ReadProperties()

	return LoadSystemConfigWithProps(props)
}

// LoadSystemConfigWithProps loads system configuration using provided properties
func LoadSystemConfigWithProps(props Config) (*SystemConfig, error) {
	config := &SystemConfig{
		AI: AIConfig{
			OpenAI: OpenAIConfig{
				APIKey:      props.OpenAIAPIKey,
				Model:       "gpt-4",
				MaxTokens:   150,
				Temperature: 0.7,
				MaxRetries:  3,
				RetryDelay:  time.Second,
			},
			Generator: GeneratorConfig{
				BasePrompt: `You are a professional customer service representative.
Your task is to generate thoughtful responses to customer reviews.

Guidelines:
1. Be polite and professional
2. Address specific points mentioned in the review
3. For negative reviews:
   - Show empathy and understanding
   - Only offer solutions that are explicitly provided
   - Never make promises not listed in business policies
4. For positive reviews:
   - Express genuine gratitude
   - Encourage future visits
5. Keep responses concise and genuine`,
				MaxLength:      300,
				DefaultTimeout: 30 * time.Second,
			},
			Validator: ValidatorConfig{
				BasePrompt: `You are a review response validator.
Evaluate if the response is appropriate and professional.
Consider:
1. Professionalism and tone
2. Relevance to the original review
3. Completeness of response
4. Safety and appropriateness
5. Only approved solutions are offered
6. No unauthorized promises made`,
				DefaultMinScore: 0.7,
				Timeout:         15 * time.Second,
			},
		},
		Email: EmailConfig{
			SendGrid: SendGridConfig{
				APIKey:     props.SendGridAPIKey,
				TemplateID: props.SendGridTemplateID,
			},
		},
	}

	return config, config.Validate()
}

// GetDefaultClientConfig returns default client configuration
func GetDefaultClientConfig() ClientConfig {
	return ClientConfig{
		Solutions: BusinessSolutions{
			AllowedActions: []string{
				"Invite customer to contact customer service",
				"Direct to help center",
				"Speak with a manager during next visit",
			},
			ProhibitedActions: []string{
				"Offer refunds",
				"Promise compensation",
				"Make price adjustments",
			},
			ContactChannels: []ContactChannel{
				{Type: "email", Details: "support@business.com"},
				{Type: "phone", Details: "(555) 123-4567"},
			},
		},
		Style: ResponseStyle{
			Tone:      "professional",
			MaxLength: 250,
			MinLength: 50,
		},
		Validation: ValidationRules{
			MinScore: 0.7,
			ProhibitedTerms: []string{
				"refund",
				"discount",
				"compensation",
				"guarantee",
			},
		},
	}
}

// MergeConfigs combines system and client configurations
func MergeConfigs(system *SystemConfig, client *ClientConfig) *CombinedConfig {
	return &CombinedConfig{
		System: *system,
		Client: *client,
	}
}
