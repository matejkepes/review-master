package ai_service

import "time"

// UseCase represents different AI use cases in the application
type UseCase string

const (
	// ReviewResponse is for generating responses to customer reviews
	ReviewResponse UseCase = "review_response"
	// MonthlyAnalysis is for analyzing monthly review data
	MonthlyAnalysis UseCase = "monthly_analysis"
	// TestResponse is for testing AI responses
	TestResponse UseCase = "test_response"
)

// GetConfigForUseCase returns the OpenAIConfig for a specific use case
func GetConfigForUseCase(useCase UseCase, apiKey string) OpenAIConfig {
	switch useCase {
	case ReviewResponse:
		return OpenAIConfig{
			APIKey:      apiKey,
			Model:       "gpt-4o-mini",
			MaxTokens:   500,
			Temperature: 0.7,
			MaxRetries:  3,
			RetryDelay:  time.Second * 2,
		}
	case MonthlyAnalysis:
		return OpenAIConfig{
			APIKey:      apiKey,
			Model:       "gpt-4.1",
			MaxTokens:   2000,
			Temperature: 0.7,
			MaxRetries:  3,
			RetryDelay:  time.Second * 2,
		}
	case TestResponse:
		return OpenAIConfig{
			APIKey:      apiKey,
			Model:       "gpt-4o-mini",
			MaxTokens:   500,
			Temperature: 0.7,
			MaxRetries:  3,
			RetryDelay:  time.Second * 2,
		}
	default:
		// Default to review response config
		return OpenAIConfig{
			APIKey:      apiKey,
			Model:       "gpt-4o-mini",
			MaxTokens:   500,
			Temperature: 0.7,
			MaxRetries:  3,
			RetryDelay:  time.Second * 2,
		}
	}
}
