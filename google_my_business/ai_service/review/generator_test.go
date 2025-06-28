package review

import (
	"testing"
	"google_my_business/ai_service"
)

// mockProvider simulates LLM responses for testing
type mockProvider struct {
	response string
	err      error
}

var _ ai_service.LLMProvider = (*mockProvider)(nil)

func (m *mockProvider) Generate(text string) (string, error) {
	return m.response, m.err
}

func (m *mockProvider) GenerateWithPrompt(systemPrompt, userPrompt string) (string, error) {
	return m.response, m.err
}

func TestNewGenerator(t *testing.T) {
	provider := &mockProvider{}
	
	tests := []struct {
		name   string
		config GeneratorConfig
	}{
		{
			name:   "with default config",
			config: GeneratorConfig{},
		},
		{
			name: "with custom config",
			config: GeneratorConfig{
				SystemPrompt: "Custom prompt",
				MaxLength:    200,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generator := NewGenerator(provider, tt.config)
			if generator == nil {
				t.Error("NewGenerator() returned nil")
			}
			if tt.config.SystemPrompt != "" && generator.config.SystemPrompt != tt.config.SystemPrompt {
				t.Errorf("NewGenerator() system prompt = %v, want %v", 
					generator.config.SystemPrompt, tt.config.SystemPrompt)
			}
		})
	}
}

func TestResponseGenerator_Generate(t *testing.T) {
	tests := []struct {
		name     string
		ctx      ReviewContext
		mockResp string
		wantErr  bool
	}{
		{
			name: "positive review",
			ctx: ReviewContext{
				Text:          "Great service! Really enjoyed my visit.",
				Rating:        5,
				Location:      "Downtown Branch",
				Author:        "John",
				BusinessName:  "ACME Corp",
				ContactMethod: "Call us at (555) 123-4567",
			},
			mockResp: "Thank you for your wonderful feedback, John!",
			wantErr:  false,
		},
		{
			name: "negative review",
			ctx: ReviewContext{
				Text:          "Poor service, long wait times.",
				Rating:        2,
				Location:      "Mall Branch",
				Author:        "Jane",
				BusinessName:  "ACME Corp",
				ContactMethod: "Book through our app at example.com",
			},
			mockResp: "We apologize for your experience, Jane.",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &mockProvider{
				response: tt.mockResp,
			}
			
			generator := NewGenerator(provider, GeneratorConfig{})
			
			response, err := generator.Generate(tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if response != tt.mockResp {
				t.Errorf("Generate() response = %v, want %v", response, tt.mockResp)
			}
		})
	}
} 