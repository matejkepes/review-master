package ai_service

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"time"
)

// NewOpenAIProvider creates a new OpenAI provider instance
func NewOpenAIProvider(config OpenAIConfig) (*OpenAIProvider, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("API key is required")
	}
	
	client := openai.NewClient(config.APIKey)
	
	// Set default values if not provided
	if config.Model == "" {
		config.Model = openai.GPT4
	}
	if config.MaxTokens == 0 {
		config.MaxTokens = 150
	}
	if config.Temperature == 0 {
		config.Temperature = 0.7
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.RetryDelay == 0 {
		config.RetryDelay = time.Second
	}

	return &OpenAIProvider{
		client: client,
		config: config,
	}, nil
}

// GenerateWithPrompt implements LLMProvider interface with custom system prompt
func (p *OpenAIProvider) GenerateWithPrompt(systemPrompt string, userPrompt string) (string, error) {
	var lastErr error
	for attempt := 0; attempt <= p.config.MaxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(p.config.RetryDelay)
		}

		resp, err := p.client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: p.config.Model,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleSystem,
						Content: systemPrompt,
					},
					{
						Role:    openai.ChatMessageRoleUser,
						Content: userPrompt,
					},
				},
				MaxTokens:   p.config.MaxTokens,
				Temperature: p.config.Temperature,
			},
		)

		// If we get an error, store it and continue to next retry
		if err != nil {
			lastErr = err
			continue
		}

		// Check for empty response
		if len(resp.Choices) == 0 {
			lastErr = fmt.Errorf("no response generated")
			continue
		}

		// Success case
		return resp.Choices[0].Message.Content, nil
	}

	// If we get here, we've exhausted all retries
	return "", fmt.Errorf("failed after %d retries. Last error: %v", p.config.MaxRetries, lastErr)
}

// Generate is a convenience method that uses a default system prompt
func (p *OpenAIProvider) Generate(userPrompt string) (string, error) {
	return p.GenerateWithPrompt(
		"You are a professional customer service representative. Generate a polite and constructive response to the customer review.",
		userPrompt,
	)
} 