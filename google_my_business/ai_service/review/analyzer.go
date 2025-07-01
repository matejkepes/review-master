package review

import (
	"encoding/json"
	"errors"
	"fmt"
	"google_my_business/ai_service"
	"google_my_business/shared"
	"strings"
	"time"
)

// Review represents a single review for analysis
type Review struct {
	ID     string    // Unique identifier for the review (optional)
	Text   string    // The review text content
	Rating int       // Star rating (1-5)
	Date   time.Time // When the review was posted
}

// ReviewBatch represents a collection of reviews for analysis
type ReviewBatch struct {
	Reviews      []Review // Collection of reviews to analyze
	LocationName string   // Name of the location being reviewed
	LocationID   string   // Unique identifier for the location
	BusinessName string   // Name of the business
	ClientID     int      // Client identifier
	PostalCode   string   // Postal code of the location
	ReportPeriod struct { // Time period for the analysis
		StartDate string `json:"start_date"` // ISO format: "2023-05-01"
		EndDate   string `json:"end_date"`   // ISO format: "2023-05-31"
	}
}

// AnalyzerConfig defines configuration options for the ReviewAnalyzer
type AnalyzerConfig struct {
	// Prompt template settings
	SystemPrompt string // System prompt for the LLM

	// Model settings
	MaxTokens int    // Maximum tokens to generate
	ModelName string // Name of the model being used

	// Analysis options
	IncludeThemes          bool // Include theme extraction in analysis
	IncludeRecommendations bool // Include recommendations in analysis

	// Analyzer identification
	AnalyzerID   string // Unique identifier for this analyzer
	AnalyzerName string // Display name for this analyzer
}

// ReviewAnalyzer analyzes batches of reviews using LLM
type ReviewAnalyzer struct {
	provider ai_service.LLMProvider // Provider for LLM functionality
	config   AnalyzerConfig         // Configuration options
}

// Define the complete analysis schema as a template constant
const analysisSchemaTemplate = `{
  "analysis": {
    "overall_summary": {
      // Provide a general overview of customer feedback based on all reviews
      "summary_text": "string",
      // List 3-5 major positive themes found across multiple reviews
      "positive_themes": ["string"],
      // List 3-5 major negative themes found across multiple reviews
      "negative_themes": ["string"],
      // Describe how passengers perceive the business (punctuality, reliability, friendliness)
      "overall_perception": "string",
      // Average star rating across all reviews
      "average_rating": 4.5
    },
    
    "sentiment_analysis": {
      // Count of positive reviews (4-5 stars)
      "positive_count": 45,
      // Percentage of positive reviews
      "positive_percentage": 75.0,
      // Count of neutral reviews (3 stars)
      "neutral_count": 10,
      // Percentage of neutral reviews
      "neutral_percentage": 16.7,
      // Count of negative reviews (1-2 stars)
      "negative_count": 5,
      // Percentage of negative reviews
      "negative_percentage": 8.3,
      // Total number of reviews analyzed
      "total_reviews": 60,
      // Brief description of the overall sentiment trend (improving, declining, stable)
      "sentiment_trend": "string"
    },
    
    "key_takeaways": {
      // List 3-5 major positive aspects based on reviews
      "strengths": [
        {
          // The category of strength (e.g., punctuality, friendly drivers, affordability)
          "category": "string",
          // Detailed description of this strength and why it matters
          "description": "string",
          // A representative example from an actual review
          "example": "string"
        }
      ],
      // List 3-5 common complaints from negative reviews
      "areas_for_improvement": [
        {
          // The category of issue (e.g., delays, poor communication, driver behavior)
          "category": "string",
          // Detailed description of the issue and its impact
          "description": "string",
          // A representative example from an actual review
          "example": "string"
        }
      ]
    },
    
    "negative_review_breakdown": {
      // Categorize negative reviews (1-2 stars) into these categories
      "categories": [
        {
          // Pre-defined category
          "name": "Missed or Delayed Pre-Bookings",
          // Number of reviews in this category
          "count": 3,
          // Percentage of negative reviews in this category
          "percentage": 60.0
        },
        {
          // Pre-defined category
          "name": "Poor Communication",
          // Number of reviews in this category
          "count": 1,
          // Percentage of negative reviews in this category
          "percentage": 20.0
        },
        {
          // Pre-defined category
          "name": "Driver Behavior",
          // Number of reviews in this category
          "count": 1,
          // Percentage of negative reviews in this category
          "percentage": 20.0
        },
        {
          // Pre-defined category
          "name": "Pricing Concerns",
          // Number of reviews in this category
          "count": 0,
          // Percentage of negative reviews in this category
          "percentage": 0.0
        }
      ],
      // List 3-5 specific, actionable steps to address the negative issues
      "improvement_recommendations": [
        "string", "string", "string"
      ]
    },
    
    "training_recommendations": {
      // List 3-4 strategies for operators (dispatch & customer service)
      // to improve communication and dispatch efficiency
      "for_operators": [
        "string", "string", "string"
      ],
      // List 3-4 best practices for taxi drivers regarding customer
      // service, timely arrivals, and enhancing passenger experience
      "for_drivers": [
        "string", "string", "string"
      ]
    }
  }
}`

// NewAnalyzer creates a new ReviewAnalyzer with default configuration
func NewAnalyzer(provider ai_service.LLMProvider, config AnalyzerConfig) *ReviewAnalyzer {
	// Apply default values for empty config fields
	if config.SystemPrompt == "" {
		config.SystemPrompt = "You are a business analytics expert analyzing customer reviews. Provide insights in JSON format."
	}

	if config.MaxTokens <= 0 {
		config.MaxTokens = 2000 // Default to a reasonably large context
	}

	return &ReviewAnalyzer{
		provider: provider,
		config:   config,
	}
}

// preparePrompt creates a prompt for the LLM based on a batch of reviews
func (a *ReviewAnalyzer) preparePrompt(batch ReviewBatch) string {
	// Format reviews as a string
	var reviewsText strings.Builder
	for i, review := range batch.Reviews {
		reviewsText.WriteString(fmt.Sprintf("### Review %d\n", i+1))
		reviewsText.WriteString(fmt.Sprintf("Rating: %d/5 stars\n", review.Rating))
		reviewsText.WriteString(fmt.Sprintf("Date: %s\n", review.Date.Format("2006-01-02")))
		reviewsText.WriteString(fmt.Sprintf("Text: %s\n\n", review.Text))
	}

	// Create the prompt using a backtick string literal with the complete JSON schema
	return fmt.Sprintf(`## Business Information
Business Name: %s
Location: %s
Location ID: %s
Postal Code: %s
Report Period: %s to %s

## Reviews
%s

## Analysis Instructions
Please analyze these reviews and provide a comprehensive analysis including:
1. Overall summary and themes
2. Sentiment breakdown (positive, neutral, negative)
3. Key strengths and areas for improvement
4. Recommendations for addressing issues
5. Training suggestions for staff

Provide your analysis in the exact JSON format specified below:

%s`,
		batch.BusinessName,
		batch.LocationName,
		batch.LocationID,
		batch.PostalCode,
		batch.ReportPeriod.StartDate, batch.ReportPeriod.EndDate,
		reviewsText.String(),
		analysisSchemaTemplate)
}

// Analyze processes a batch of reviews and returns a structured analysis
func (a *ReviewAnalyzer) Analyze(batch ReviewBatch) (*shared.AnalysisResult, error) {
	// Filter out reviews with empty text, but keep track of all ratings for average calculation
	var filteredReviews []Review
	var allRatings []int

	for _, review := range batch.Reviews {
		// Always track the rating for average calculation
		allRatings = append(allRatings, review.Rating)

		// Only include reviews with actual text for analysis
		if strings.TrimSpace(review.Text) != "" {
			filteredReviews = append(filteredReviews, review)
		}
	}

	// Create a filtered batch for analysis
	filteredBatch := batch
	filteredBatch.Reviews = filteredReviews

	// Check if we have any reviews with text
	if len(filteredBatch.Reviews) == 0 {
		// If no reviews have text but we have ratings, create a minimal analysis result
		if len(allRatings) > 0 {
			// Calculate average rating
			var total float64
			for _, rating := range allRatings {
				total += float64(rating)
			}
			avgRating := total / float64(len(allRatings))

			// Create a minimal analysis result with metadata
			result := &shared.AnalysisResult{
				Analysis: shared.Analysis{
					OverallSummary: shared.OverallSummary{
						SummaryText:       "No text content in reviews to analyze. Analysis based on ratings only.",
						AverageRating:     avgRating,
						PositiveThemes:    []string{"Not available - no review text"},
						NegativeThemes:    []string{"Not available - no review text"},
						OverallPerception: "Analysis limited to numeric ratings only.",
					},
					SentimentAnalysis: shared.SentimentAnalysis{
						TotalReviews:   len(allRatings),
						SentimentTrend: "Unable to determine from ratings only",
					},
					KeyTakeaways: shared.KeyTakeaways{
						Strengths: []shared.Insight{
							{Category: "Limited Data", Description: "Unable to identify strengths without review text", Example: ""},
						},
						AreasForImprovement: []shared.Insight{
							{Category: "Limited Data", Description: "Unable to identify areas for improvement without review text", Example: ""},
						},
					},
				},
				Metadata: shared.AnalysisMetadata{
					GeneratedAt:   time.Now(),
					ReviewCount:   len(allRatings),
					LocationID:    batch.LocationID,
					LocationName:  batch.LocationName,
					BusinessName:  batch.BusinessName,
					ClientID:      batch.ClientID,
					AnalyzerID:    a.config.AnalyzerID,
					AnalyzerName:  a.config.AnalyzerName,
					AnalyzerModel: a.config.ModelName,
				},
			}

			// Set report period separately due to the json tags
			result.Metadata.ReportPeriod.StartDate = batch.ReportPeriod.StartDate
			result.Metadata.ReportPeriod.EndDate = batch.ReportPeriod.EndDate

			return result, nil
		}

		return nil, errors.New("cannot analyze empty review batch")
	}

	// Prepare the prompt for the LLM
	prompt := a.preparePrompt(filteredBatch)

	// Generate response using the LLM provider
	response, err := a.provider.GenerateWithPrompt(a.config.SystemPrompt, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate analysis: %w", err)
	}

	// Extract JSON from the response (in case LLM returns additional text)
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}")

	if jsonStart == -1 || jsonEnd == -1 || jsonEnd <= jsonStart {
		return nil, fmt.Errorf("failed to extract valid JSON from response: %s", response)
	}

	jsonData := response[jsonStart : jsonEnd+1]

	// Parse the LLM response into a structured analysis result
	result := &shared.AnalysisResult{}
	if err := json.Unmarshal([]byte(jsonData), result); err != nil {
		// print the response
		fmt.Println(response)
		return nil, fmt.Errorf("failed to parse analysis response: %w", err)
	}

	// Validate required fields
	if result.Analysis.OverallSummary.SummaryText == "" {
		return nil, errors.New("analysis result missing summary")
	}

	if result.Analysis.SentimentAnalysis.SentimentTrend == "" {
		return nil, errors.New("analysis result missing sentiment trend")
	}

	if len(result.Analysis.KeyTakeaways.Strengths) == 0 {
		return nil, errors.New("analysis result missing strengths")
	}

	if len(result.Analysis.KeyTakeaways.AreasForImprovement) == 0 {
		return nil, errors.New("analysis result missing areas for improvement")
	}

	// Always calculate average rating using all reviews (including those without text)
	// This ensures consistency regardless of what the AI provides
	var total float64
	for _, rating := range allRatings {
		total += float64(rating)
	}
	result.Analysis.OverallSummary.AverageRating = total / float64(len(allRatings))

	// Add metadata to the result
	result.Metadata = shared.AnalysisMetadata{
		GeneratedAt:   time.Now(),
		ReviewCount:   len(allRatings), // Use the total count including non-text reviews
		LocationID:    batch.LocationID,
		LocationName:  batch.LocationName,
		BusinessName:  batch.BusinessName,
		ClientID:      batch.ClientID,
		AnalyzerID:    a.config.AnalyzerID,
		AnalyzerName:  a.config.AnalyzerName,
		AnalyzerModel: a.config.ModelName,
	}

	// Set the report period separately to handle the struct with json tags
	result.Metadata.ReportPeriod.StartDate = batch.ReportPeriod.StartDate
	result.Metadata.ReportPeriod.EndDate = batch.ReportPeriod.EndDate

	// Ensure TotalReviews reflects all reviews, not just those with text content
	// This must be done after all validation to avoid breaking the analysis logic
	result.Analysis.SentimentAnalysis.TotalReviews = len(allRatings)

	// Recalculate percentages based on the corrected total to ensure mathematical consistency
	// The AI calculated percentages based only on reviews with text, but we need percentages
	// based on all reviews (including those without text content)
	if result.Analysis.SentimentAnalysis.TotalReviews > 0 {
		totalReviews := float64(result.Analysis.SentimentAnalysis.TotalReviews)
		result.Analysis.SentimentAnalysis.PositivePercentage = (float64(result.Analysis.SentimentAnalysis.PositiveCount) / totalReviews) * 100.0
		result.Analysis.SentimentAnalysis.NeutralPercentage = (float64(result.Analysis.SentimentAnalysis.NeutralCount) / totalReviews) * 100.0
		result.Analysis.SentimentAnalysis.NegativePercentage = (float64(result.Analysis.SentimentAnalysis.NegativeCount) / totalReviews) * 100.0
	}

	return result, nil
}
