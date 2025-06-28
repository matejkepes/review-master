package review

import (
    "encoding/json"
    "fmt"
    "google_my_business/ai_service"
)

// ValidationResult holds the structured validation output
type ValidationResult struct {
    IsValid     bool     `json:"isValid"`
    Score       float32  `json:"score"`
    Reasons     []string `json:"reasons"`
    Suggestions []string `json:"suggestions,omitempty"`
}

// ValidatorConfig holds configuration for the response validator
type ValidatorConfig struct {
    SystemPrompt string
    MinScore     float32
}

// ResponseValidator validates generated review responses
type ResponseValidator struct {
    provider ai_service.LLMProvider
    config   ValidatorConfig
}

// NewValidator creates a new response validator
func NewValidator(provider ai_service.LLMProvider, config ValidatorConfig) *ResponseValidator {
    // Set default system prompt if not provided
    if config.SystemPrompt == "" {
        config.SystemPrompt = `You are a review response validator.
Evaluate if the response is appropriate and professional.
Consider:
1. Professionalism and tone
2. Relevance to the original review
3. Completeness of response
4. Safety and appropriateness
5. Only approved solutions are offered
6. No unauthorized promises made

Return a JSON object with the following structure:
{
  "isValid": true/false,
  "score": 0.0-1.0,
  "reasons": ["reason1", "reason2"],
  "suggestions": ["suggestion1", "suggestion2"]
}`
    }

    // Set default minimum score if not provided
    if config.MinScore == 0 {
        config.MinScore = 0.7
    }

    return &ResponseValidator{
        provider: provider,
        config:   config,
    }
}

// Validate checks if a generated response is appropriate
func (v *ResponseValidator) Validate(ctx ReviewContext, response string) (ValidationResult, error) {
    prompt := v.formatValidationPrompt(ctx, response)
    
    // Get validation response from LLM
    result, err := v.provider.GenerateWithPrompt(v.config.SystemPrompt, prompt)
    if err != nil {
        return ValidationResult{}, fmt.Errorf("validation generation failed: %w", err)
    }

    // Parse JSON response
    var validationResult ValidationResult
    if err := json.Unmarshal([]byte(result), &validationResult); err != nil {
        return ValidationResult{}, fmt.Errorf("failed to parse validation result: %w", err)
    }

    // Override IsValid based on minimum score if necessary
    if validationResult.Score < v.config.MinScore {
        validationResult.IsValid = false
        validationResult.Reasons = append(validationResult.Reasons, 
            fmt.Sprintf("Score %.2f below minimum threshold %.2f", validationResult.Score, v.config.MinScore))
    }

    return validationResult, nil
}

// formatValidationPrompt creates a structured prompt for the validator
func (v *ResponseValidator) formatValidationPrompt(ctx ReviewContext, response string) string {
    return fmt.Sprintf(`Please evaluate this review response:

Original Review:
"%s"
Rating: %d stars
Author: %s
Location: %s
Business: %s
Contact Method: %s

Generated Response:
"%s"

Evaluate the response and return a JSON object with your assessment.`, 
        ctx.Text, ctx.Rating, ctx.Author, ctx.Location, ctx.BusinessName, ctx.ContactMethod, response)
} 