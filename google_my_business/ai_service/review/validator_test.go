package review

import (
    "encoding/json"
    "testing"
    "google_my_business/ai_service"
)

// mockValidatorProvider simulates LLM responses for validation
type mockValidatorProvider struct {
    response string
    err      error
}

var _ ai_service.LLMProvider = (*mockValidatorProvider)(nil)

func (m *mockValidatorProvider) Generate(text string) (string, error) {
    return m.response, m.err
}

func (m *mockValidatorProvider) GenerateWithPrompt(systemPrompt, userPrompt string) (string, error) {
    return m.response, m.err
}

func TestNewValidator(t *testing.T) {
    provider := &mockValidatorProvider{}
    
    tests := []struct {
        name   string
        config ValidatorConfig
    }{
        {
            name:   "with default config",
            config: ValidatorConfig{},
        },
        {
            name: "with custom config",
            config: ValidatorConfig{
                SystemPrompt: "Custom prompt",
                MinScore:    0.8,
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            validator := NewValidator(provider, tt.config)
            if validator == nil {
                t.Error("NewValidator() returned nil")
            }
            if tt.config.SystemPrompt != "" && validator.config.SystemPrompt != tt.config.SystemPrompt {
                t.Errorf("NewValidator() system prompt = %v, want %v", 
                    validator.config.SystemPrompt, tt.config.SystemPrompt)
            }
        })
    }
}

func TestResponseValidator_Validate(t *testing.T) {
    tests := []struct {
        name        string
        ctx         ReviewContext
        response    string
        mockResult  string
        wantErr     bool
    }{
        {
            name: "valid response",
            ctx: ReviewContext{
                Text:          "Great service!",
                Rating:        5,
                Author:        "John",
                BusinessName:  "ACME Corp",
                ContactMethod: "Call us at (555) 123-4567",
            },
            response: "Thank you for your kind feedback, John!",
            mockResult: `{"isValid":true,"score":0.9,"reasons":["Professional","Personalized","Appropriate tone"]}`,
            wantErr: false,
        },
        {
            name: "invalid response",
            ctx: ReviewContext{
                Text:          "Terrible service!",
                Rating:        1,
                Author:        "Jane",
                BusinessName:  "ACME Corp",
                ContactMethod: "Book through our app at example.com",
            },
            response: "You're wrong about our service.",
            mockResult: `{"isValid":false,"score":0.3,"reasons":["Defensive tone","Unprofessional"],"suggestions":["Show empathy","Offer solution"]}`,
            wantErr: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            provider := &mockValidatorProvider{
                response: tt.mockResult,
            }
            
            validator := NewValidator(provider, ValidatorConfig{MinScore: 0.7})
            
            result, err := validator.Validate(tt.ctx, tt.response)
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            
            if !tt.wantErr {
                var expectedResult ValidationResult
                if err := json.Unmarshal([]byte(tt.mockResult), &expectedResult); err != nil {
                    t.Errorf("Failed to parse expected result: %v", err)
                    return
                }
                
                if result.IsValid != expectedResult.IsValid {
                    t.Errorf("Validate() isValid = %v, want %v", result.IsValid, expectedResult.IsValid)
                }
                if result.Score != expectedResult.Score {
                    t.Errorf("Validate() score = %v, want %v", result.Score, expectedResult.Score)
                }
            }
        })
    }
} 