package ai_service

import (
    "context"
    "errors"
    "testing"
    "time"
    "github.com/sashabaranov/go-openai"
)

// mockOpenAIClient implements necessary OpenAI client methods for testing
type mockOpenAIClient struct {
    shouldFail bool
    response   string
}

func (m *mockOpenAIClient) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
    if m.shouldFail {
        return openai.ChatCompletionResponse{}, errors.New("mock API error")
    }

    return openai.ChatCompletionResponse{
        Choices: []openai.ChatCompletionChoice{
            {
                Message: openai.ChatCompletionMessage{
                    Content: m.response,
                },
            },
        },
    }, nil
}

type mockOpenAIClientWithRetry struct {
    failCount int
    curCount  int
    response  string
}

func (m *mockOpenAIClientWithRetry) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
    m.curCount++
    if m.curCount <= m.failCount {
        return openai.ChatCompletionResponse{}, errors.New("temporary error")
    }

    return openai.ChatCompletionResponse{
        Choices: []openai.ChatCompletionChoice{
            {
                Message: openai.ChatCompletionMessage{
                    Content: m.response,
                },
            },
        },
    }, nil
}

// TestNewOpenAIProvider tests the provider creation
func TestNewOpenAIProvider(t *testing.T) {
    tests := []struct {
        name    string
        config  OpenAIConfig
        wantErr bool
    }{
        {
            name: "valid config",
            config: OpenAIConfig{
                APIKey: "test-key",
            },
            wantErr: false,
        },
        {
            name: "missing API key",
            config: OpenAIConfig{},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            provider, err := NewOpenAIProvider(tt.config)
            if (err != nil) != tt.wantErr {
                t.Errorf("NewOpenAIProvider() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !tt.wantErr && provider == nil {
                t.Error("NewOpenAIProvider() returned nil provider without error")
            }
        })
    }
}

// TestOpenAIProvider_Generate tests the response generation
func TestOpenAIProvider_Generate(t *testing.T) {
    tests := []struct {
        name       string
        mockResp   string
        shouldFail bool
        wantErr    bool
    }{
        {
            name:       "successful generation",
            mockResp:   "Thank you for your feedback. We apologize for your experience.",
            shouldFail: false,
            wantErr:    false,
        },
        {
            name:       "API error",
            mockResp:   "",
            shouldFail: true,
            wantErr:    true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            provider := &OpenAIProvider{
                client: &mockOpenAIClient{
                    shouldFail: tt.shouldFail,
                    response:   tt.mockResp,
                },
                config: OpenAIConfig{
                    Model:       "gpt-4",
                    MaxTokens:   150,
                    Temperature: 0.7,
                },
            }

            response, err := provider.Generate("Test review")
            if (err != nil) != tt.wantErr {
                t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !tt.wantErr && response != tt.mockResp {
                t.Errorf("Generate() response = %v, want %v", response, tt.mockResp)
            }
        })
    }
}

func TestOpenAIProvider_GenerateWithRetry(t *testing.T) {
    tests := []struct {
        name      string
        failCount int
        maxRetry  int
        response  string
        wantErr   bool
    }{
        {
            name:      "success after 2 retries",
            failCount: 2,
            maxRetry:  3,
            response:  "Thank you for your feedback",
            wantErr:   false,
        },
        {
            name:      "failure after max retries",
            failCount: 4,
            maxRetry:  3,
            response:  "Thank you for your feedback",
            wantErr:   true,
        },
        {
            name:      "immediate success",
            failCount: 0,
            maxRetry:  3,
            response:  "Thank you for your feedback",
            wantErr:   false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockClient := &mockOpenAIClientWithRetry{
                failCount: tt.failCount,
                response:  tt.response,
            }

            provider := &OpenAIProvider{
                client: mockClient,
                config: OpenAIConfig{
                    Model:       "gpt-4",
                    MaxTokens:   150,
                    Temperature: 0.7,
                    MaxRetries:  tt.maxRetry,
                    RetryDelay:  time.Millisecond, // Short delay for tests
                },
            }

            response, err := provider.Generate("Test review")
            if (err != nil) != tt.wantErr {
                t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !tt.wantErr {
                if response != tt.response {
                    t.Errorf("Generate() response = %v, want %v", response, tt.response)
                }
                if mockClient.curCount != tt.failCount+1 {
                    t.Errorf("Generate() retry count = %v, want %v", mockClient.curCount, tt.failCount+1)
                }
            }
        })
    }
}

func TestOpenAIProvider_GenerateWithPrompt(t *testing.T) {
    tests := []struct {
        name         string
        systemPrompt string
        userPrompt   string
        mockResp     string
        shouldFail   bool
        wantErr      bool
    }{
        {
            name:         "custom system prompt",
            systemPrompt: "You are a helpful assistant.",
            userPrompt:   "Hello!",
            mockResp:     "Hi! How can I help you today?",
            shouldFail:   false,
            wantErr:      false,
        },
        {
            name:         "empty system prompt",
            systemPrompt: "",
            userPrompt:   "Hello!",
            mockResp:     "Hi!",
            shouldFail:   false,
            wantErr:      false,
        },
        {
            name:         "API error",
            systemPrompt: "You are a helpful assistant.",
            userPrompt:   "Hello!",
            mockResp:     "",
            shouldFail:   true,
            wantErr:      true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            provider := &OpenAIProvider{
                client: &mockOpenAIClient{
                    shouldFail: tt.shouldFail,
                    response:   tt.mockResp,
                },
                config: OpenAIConfig{
                    Model:       "gpt-4",
                    MaxTokens:   150,
                    Temperature: 0.7,
                },
            }

            response, err := provider.GenerateWithPrompt(tt.systemPrompt, tt.userPrompt)
            if (err != nil) != tt.wantErr {
                t.Errorf("GenerateWithPrompt() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !tt.wantErr && response != tt.mockResp {
                t.Errorf("GenerateWithPrompt() response = %v, want %v", response, tt.mockResp)
            }
        })
    }
} 