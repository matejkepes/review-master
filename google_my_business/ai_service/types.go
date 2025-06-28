package ai_service

import (
    "context"
    "github.com/sashabaranov/go-openai"
    "time"
)

// Review represents the input review data
type Review struct {
    Text     string
    Rating   int
    Location string
    Author   string
}

// Response represents the generated response
type Response struct {
    Text       string
    Validated  bool
    RetryCount int
    Error      error
}

// OpenAIConfig holds the configuration for OpenAI
type OpenAIConfig struct {
    APIKey      string
    Model       string
    MaxTokens   int
    Temperature float32
    MaxRetries  int
    RetryDelay  time.Duration
}

// OpenAIClientInterface defines the methods we need from the OpenAI client
type OpenAIClientInterface interface {
    CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}

// LLMProvider interface defines methods that any LLM service must implement
type LLMProvider interface {
    Generate(text string) (string, error)
    GenerateWithPrompt(systemPrompt string, userPrompt string) (string, error)
}

// OpenAIProvider implements LLMProvider for OpenAI
type OpenAIProvider struct {
    client OpenAIClientInterface
    config OpenAIConfig
} 