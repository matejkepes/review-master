package review

// GeneratorConfig holds configuration for the review response generator
type GeneratorConfig struct {
	SystemPrompt string
	MaxLength    int
	Style        string // e.g., "professional", "friendly", "formal"
}

// ReviewContext contains all relevant information about a review
type ReviewContext struct {
	Text          string
	Rating        int    // 1-5 stars
	Location      string // business location
	Author        string
	BusinessName  string
	ContactMethod string // Client's preferred contact method
}

// Generator interface defines methods for generating review responses
type Generator interface {
	Generate(ctx ReviewContext) (string, error)
}
