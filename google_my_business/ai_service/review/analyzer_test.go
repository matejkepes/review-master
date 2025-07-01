package review

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

// MockLLMProvider implements the LLMProvider interface for testing
type MockLLMProvider struct {
	responses map[string]mockResponse
}

type mockResponse struct {
	result string
	err    error
}

// Generate implements the LLMProvider interface
func (m *MockLLMProvider) Generate(text string) (string, error) {
	if response, ok := m.responses[text]; ok {
		return response.result, response.err
	}
	return "", errors.New("no mock response for given text")
}

// GenerateWithPrompt implements the LLMProvider interface
func (m *MockLLMProvider) GenerateWithPrompt(systemPrompt, userPrompt string) (string, error) {
	// Check for exact match first
	key := systemPrompt + ":::" + userPrompt
	if response, ok := m.responses[key]; ok {
		return response.result, response.err
	}

	// Check for wildcard matches where 'any' is used
	wildcardSystemKey := systemPrompt + ":::" + "any"
	if response, ok := m.responses[wildcardSystemKey]; ok {
		return response.result, response.err
	}

	wildcardUserKey := "any" + ":::" + userPrompt
	if response, ok := m.responses[wildcardUserKey]; ok {
		return response.result, response.err
	}

	// Check for complete wildcard
	completeWildcard := "any:::any"
	if response, ok := m.responses[completeWildcard]; ok {
		return response.result, response.err
	}

	return "", errors.New("no mock response for given prompts")
}

// NewMockLLMProvider creates a new mock provider with predefined responses
func NewMockLLMProvider() *MockLLMProvider {
	return &MockLLMProvider{
		responses: make(map[string]mockResponse),
	}
}

// AddMockResponse adds a mock response for a given prompt combination
func (m *MockLLMProvider) AddMockResponse(systemPrompt, userPrompt, result string, err error) {
	key := systemPrompt + ":::" + userPrompt
	m.responses[key] = mockResponse{result: result, err: err}
}

// TestStructure tests the basic structure and initialization of the analyzer
func TestStructure(t *testing.T) {
	// Create a mock provider
	mockProvider := NewMockLLMProvider()

	// Test with default config
	config := AnalyzerConfig{}
	analyzer := NewAnalyzer(mockProvider, config)

	// Check that the analyzer was created properly
	if analyzer == nil {
		t.Fatal("NewAnalyzer() returned nil")
	}

	// Check that the provider was properly assigned
	if analyzer.provider == nil {
		t.Error("Provider not assigned")
	}

	// Check that default values were applied
	if analyzer.config.SystemPrompt == "" {
		t.Error("Default SystemPrompt not applied")
	}

	if analyzer.config.MaxTokens <= 0 {
		t.Error("Default MaxTokens not applied")
	}

	// Test with custom config values
	customConfig := AnalyzerConfig{
		SystemPrompt: "Custom prompt",
		MaxTokens:    1000,
	}

	customAnalyzer := NewAnalyzer(mockProvider, customConfig)

	// Check that custom values were retained
	if customAnalyzer.config.SystemPrompt != "Custom prompt" {
		t.Errorf("SystemPrompt = %s, want %s", customAnalyzer.config.SystemPrompt, "Custom prompt")
	}

	if customAnalyzer.config.MaxTokens != 1000 {
		t.Errorf("MaxTokens = %d, want %d", customAnalyzer.config.MaxTokens, 1000)
	}
}

// TestPreparePrompt tests the prompt generation for review analysis
func TestPreparePrompt(t *testing.T) {
	// Create analyzer with mock provider
	mockProvider := NewMockLLMProvider()

	// Create a config that includes themes and recommendations
	config := AnalyzerConfig{
		IncludeThemes:          true,
		IncludeRecommendations: true,
	}
	analyzer := NewAnalyzer(mockProvider, config)

	// Create a test review batch
	batch := ReviewBatch{
		Reviews:      createTestReviews(),
		LocationName: "Downtown Taxi Office",
		LocationID:   "loc_123",
		BusinessName: "City Taxi",
		ClientID:     42,
		PostalCode:   "10001",
	}

	// Set time period
	batch.ReportPeriod.StartDate = "2023-05-01"
	batch.ReportPeriod.EndDate = "2023-05-31"

	// Generate the prompt
	prompt := analyzer.preparePrompt(batch)

	// Verify prompt contains essential business information
	if !containsString(prompt, "City Taxi") {
		t.Error("Prompt missing business name")
	}

	if !containsString(prompt, "Downtown Taxi Office") {
		t.Error("Prompt missing location name")
	}

	if !containsString(prompt, "2023-05-01") || !containsString(prompt, "2023-05-31") {
		t.Error("Prompt missing report period dates")
	}

	// Verify prompt contains review information
	for _, review := range batch.Reviews {
		if !containsString(prompt, review.Text) {
			t.Errorf("Prompt missing review text: %s", review.Text)
		}

		// Check that rating is included
		ratingStr := fmt.Sprintf("%d", review.Rating)
		if !containsString(prompt, ratingStr+" star") && !containsString(prompt, ratingStr+"/5") {
			t.Errorf("Prompt missing review rating: %d", review.Rating)
		}
	}

	// Verify prompt contains JSON schema instructions
	if !containsString(prompt, "JSON") {
		t.Error("Prompt missing JSON format instructions")
	}

	// Verify key analysis sections are requested
	analysisElements := []string{
		"summary",
		"sentiment",
		"themes",
		"recommendations",
	}

	for _, element := range analysisElements {
		if !containsString(prompt, element) {
			t.Errorf("Prompt missing request for '%s' analysis", element)
		}
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// createTestReviews creates a set of test reviews for use in tests
func createTestReviews() []Review {
	return []Review{
		{
			ID:     "rev_1",
			Text:   "Great service, the driver was very friendly",
			Rating: 5,
			Date:   time.Date(2023, 5, 10, 14, 30, 0, 0, time.UTC),
		},
		{
			ID:     "rev_2",
			Text:   "Car was clean but the driver was late",
			Rating: 3,
			Date:   time.Date(2023, 5, 15, 9, 45, 0, 0, time.UTC),
		},
		{
			ID:     "rev_3",
			Text:   "Terrible experience. Driver was rude and car smelled bad.",
			Rating: 1,
			Date:   time.Date(2023, 5, 20, 18, 15, 0, 0, time.UTC),
		},
	}
}

// TestAnalyze tests the full analysis flow
func TestAnalyze(t *testing.T) {
	// Create mock provider with pre-defined response
	mockProvider := NewMockLLMProvider()

	// Define the system prompt that will be used
	systemPrompt := "You are a business analytics expert analyzing customer reviews."

	mockResponse := `{
		"analysis": {
			"overall_summary": {
				"summary_text": "Mixed reviews with average satisfaction. Some customers appreciate friendly service while others complain about timeliness and cleanliness.",
				"positive_themes": ["Friendly drivers", "Comfortable rides", "Good customer service"],
				"negative_themes": ["Timeliness issues", "Vehicle cleanliness", "Communication problems"],
				"overall_perception": "Passengers generally see the service as reliable but inconsistent",
				"average_rating": 3.0
			},
			"sentiment_analysis": {
				"positive_count": 1,
				"positive_percentage": 33.3,
				"neutral_count": 1,
				"neutral_percentage": 33.3,
				"negative_count": 1,
				"negative_percentage": 33.3,
				"total_reviews": 3,
				"sentiment_trend": "stable"
			},
			"key_takeaways": {
				"strengths": [
					{
						"category": "Driver Friendliness",
						"description": "Customers frequently mention friendly and professional drivers",
						"example": "Great service, the driver was very friendly"
					}
				],
				"areas_for_improvement": [
					{
						"category": "Timeliness",
						"description": "Some customers experienced delays with their taxi service",
						"example": "Car was clean but the driver was late"
					},
					{
						"category": "Vehicle Cleanliness",
						"description": "A few customers complained about the condition of vehicles",
						"example": "Car smelled bad"
					}
				]
			},
			"negative_review_breakdown": {
				"categories": [
					{
						"name": "Driver Behavior",
						"count": 1,
						"percentage": 100.0
					}
				],
				"improvement_recommendations": [
					"Implement a better system for tracking and improving on-time performance",
					"Establish vehicle cleanliness standards and regular checks",
					"Provide additional customer service training for drivers"
				]
			},
			"training_recommendations": {
				"for_operators": [
					"Improve dispatch efficiency to reduce wait times",
					"Implement better communication systems to alert customers of delays",
					"Create a system to verify vehicle cleanliness before shifts"
				],
				"for_drivers": [
					"Maintain professional and friendly demeanor with all customers",
					"Communicate proactively about any potential delays",
					"Ensure vehicle cleanliness before each shift"
				]
			}
		}
	}`

	// Configure the mock response with the specific system prompt
	// The user prompt will be matched regardless of content
	mockProvider.AddMockResponse(systemPrompt, "any", mockResponse, nil)

	// Create analyzer with fully specified config
	config := AnalyzerConfig{
		SystemPrompt:           "You are a business analytics expert analyzing customer reviews.",
		MaxTokens:              2000,
		ModelName:              "gpt-4",
		IncludeThemes:          true,
		IncludeRecommendations: true,
		AnalyzerID:             "taxi-review-analyzer-v1",
		AnalyzerName:           "Taxi Service Review Analyzer",
	}
	analyzer := NewAnalyzer(mockProvider, config)

	// Create a test review batch
	batch := ReviewBatch{
		Reviews:      createTestReviews(),
		LocationName: "Downtown Taxi Office",
		LocationID:   "loc_123",
		BusinessName: "City Taxi",
		ClientID:     42,
		PostalCode:   "10001",
	}
	batch.ReportPeriod.StartDate = "2023-05-01"
	batch.ReportPeriod.EndDate = "2023-05-31"

	// Execute analysis
	result, err := analyzer.Analyze(batch)

	// Verify no errors occurred
	if err != nil {
		t.Fatalf("Analyze failed with error: %v", err)
	}

	// Check that result is not nil
	if result == nil {
		t.Fatal("Analyze returned nil result")
	}

	// Verify analysis content matches expected response
	if result.Analysis.OverallSummary.SummaryText == "" {
		t.Error("Analysis missing overall summary")
	}
	if result.Analysis.SentimentAnalysis.SentimentTrend != "stable" {
		t.Errorf("Unexpected sentiment trend: got %s, want %s", result.Analysis.SentimentAnalysis.SentimentTrend, "stable")
	}
	if len(result.Analysis.KeyTakeaways.Strengths) == 0 {
		t.Error("Analysis missing strengths")
	}
	if len(result.Analysis.KeyTakeaways.AreasForImprovement) == 0 {
		t.Error("Analysis missing areas for improvement")
	}
	if len(result.Analysis.NegativeReviewBreakdown.Categories) == 0 {
		t.Error("Analysis missing negative review breakdown")
	}
	if len(result.Analysis.TrainingRecommendations.ForOperators) == 0 {
		t.Error("Analysis missing training recommendations for operators")
	}
	if len(result.Analysis.TrainingRecommendations.ForDrivers) == 0 {
		t.Error("Analysis missing training recommendations for drivers")
	}
	if result.Analysis.OverallSummary.AverageRating != 3.0 {
		t.Errorf("Incorrect average rating: got %.1f, want %.1f", result.Analysis.OverallSummary.AverageRating, 3.0)
	}

	// Verify metadata
	if result.Metadata.ReviewCount != len(batch.Reviews) {
		t.Errorf("Incorrect review count in metadata: got %d, want %d", result.Metadata.ReviewCount, len(batch.Reviews))
	}
	if result.Metadata.LocationID != batch.LocationID {
		t.Errorf("Incorrect location ID in metadata: got %s, want %s", result.Metadata.LocationID, batch.LocationID)
	}
	if result.Metadata.BusinessName != batch.BusinessName {
		t.Errorf("Incorrect business name in metadata: got %s, want %s", result.Metadata.BusinessName, batch.BusinessName)
	}
	if result.Metadata.AnalyzerID != config.AnalyzerID {
		t.Errorf("Incorrect analyzer ID in metadata: got %s, want %s", result.Metadata.AnalyzerID, config.AnalyzerID)
	}
	if result.Metadata.AnalyzerName != config.AnalyzerName {
		t.Errorf("Incorrect analyzer name in metadata: got %s, want %s", result.Metadata.AnalyzerName, config.AnalyzerName)
	}
	if result.Metadata.AnalyzerModel != config.ModelName {
		t.Errorf("Incorrect model name in metadata: got %s, want %s", result.Metadata.AnalyzerModel, config.ModelName)
	}

	// Test error handling for empty review batch
	emptyBatch := ReviewBatch{
		Reviews:      []Review{},
		LocationName: "Empty Location",
		LocationID:   "loc_empty",
	}

	_, emptyErr := analyzer.Analyze(emptyBatch)
	if emptyErr == nil {
		t.Error("Expected error for empty review batch, got nil")
	}

	// Test error handling for LLM provider error
	errorProvider := NewMockLLMProvider()
	errorProvider.AddMockResponse("any", "any", "", errors.New("LLM service unavailable"))

	errorAnalyzer := NewAnalyzer(errorProvider, config)
	_, providerErr := errorAnalyzer.Analyze(batch)

	if providerErr == nil {
		t.Error("Expected error when LLM provider fails, got nil")
	}
}

// TestAnalyzeWithMixedReviews tests analysis with both text and non-text reviews
func TestAnalyzeWithMixedReviews(t *testing.T) {
	// Create a mock provider that returns a structured response
	mockProvider := NewMockLLMProvider()
	mockProvider.AddMockResponse("any", "any", `{
		"analysis": {
			"overall_summary": {
				"summary_text": "Mixed review analysis based on available text",
				"average_rating": 3.5,
				"positive_themes": ["Good service"],
				"negative_themes": ["Some issues"],
				"overall_perception": "Generally positive with areas for improvement"
			},
			"sentiment_analysis": {
				"total_reviews": 2,
				"positive_count": 1,
				"neutral_count": 0,
				"negative_count": 1,
				"positive_percentage": 50.0,
				"neutral_percentage": 0.0,
				"negative_percentage": 50.0,
				"sentiment_trend": "Mixed sentiment"
			},
			"key_takeaways": {
				"strengths": [
					{
						"category": "Service",
						"description": "Good customer service",
						"example": "Staff was helpful"
					}
				],
				"areas_for_improvement": [
					{
						"category": "Timeliness",
						"description": "Need to improve punctuality",
						"example": "Driver was late"
					}
				]
			},
			"negative_review_breakdown": {
				"categories": [
					{
						"name": "Poor Communication",
						"count": 1,
						"percentage": 100.0
					}
				],
				"improvement_recommendations": [
					"Improve communication with customers"
				]
			},
			"training_recommendations": {
				"for_operators": [
					"Enhance customer service training"
				],
				"for_drivers": [
					"Focus on punctuality training"
				]
			}
		}
	}`, nil)

	config := AnalyzerConfig{
		SystemPrompt:           "Test system prompt",
		MaxTokens:              1000,
		ModelName:              "test-model",
		IncludeThemes:          true,
		IncludeRecommendations: true,
		AnalyzerID:             "test-analyzer",
		AnalyzerName:           "Test Analyzer",
	}

	analyzer := NewAnalyzer(mockProvider, config)

	// Create batch with 4 reviews: 2 with text, 2 without text
	batch := ReviewBatch{
		Reviews: []Review{
			{
				ID:     "rev_1",
				Text:   "Great service, very friendly staff",
				Rating: 5,
				Date:   time.Date(2023, 5, 10, 14, 30, 0, 0, time.UTC),
			},
			{
				ID:     "rev_2",
				Text:   "Driver was late, not happy",
				Rating: 2,
				Date:   time.Date(2023, 5, 15, 9, 45, 0, 0, time.UTC),
			},
			{
				ID:     "rev_3",
				Text:   "", // No text content
				Rating: 4,
				Date:   time.Date(2023, 5, 20, 18, 15, 0, 0, time.UTC),
			},
			{
				ID:     "rev_4",
				Text:   "", // No text content
				Rating: 3,
				Date:   time.Date(2023, 5, 25, 12, 0, 0, 0, time.UTC),
			},
		},
		LocationName: "Test Location",
		LocationID:   "test-location-123",
		BusinessName: "Test Business",
		ClientID:     1,
		PostalCode:   "12345",
	}
	batch.ReportPeriod.StartDate = "2023-05-01"
	batch.ReportPeriod.EndDate = "2023-05-31"

	// Run analysis
	result, err := analyzer.Analyze(batch)
	if err != nil {
		t.Fatalf("Analysis failed: %v", err)
	}

	// Verify that TotalReviews includes ALL reviews (both with and without text)
	expectedTotalReviews := 4 // All 4 reviews should be counted
	if result.Analysis.SentimentAnalysis.TotalReviews != expectedTotalReviews {
		t.Errorf("TotalReviews should include all reviews (with and without text): got %d, want %d",
			result.Analysis.SentimentAnalysis.TotalReviews, expectedTotalReviews)
	}

	// Verify that metadata ReviewCount also includes all reviews
	if result.Metadata.ReviewCount != expectedTotalReviews {
		t.Errorf("Metadata ReviewCount should include all reviews: got %d, want %d",
			result.Metadata.ReviewCount, expectedTotalReviews)
	}

	// Verify that average rating is calculated from all reviews
	expectedAverage := (5 + 2 + 4 + 3) / 4.0 // 3.5
	if result.Analysis.OverallSummary.AverageRating != expectedAverage {
		t.Errorf("Average rating should include all reviews: got %f, want %f",
			result.Analysis.OverallSummary.AverageRating, expectedAverage)
	}

	expectedPositiveCount := 2 // 4 and 5 star reviews
	if result.Analysis.SentimentAnalysis.PositiveCount != expectedPositiveCount {
		t.Errorf("Expected PositiveCount %d, got %d", expectedPositiveCount, result.Analysis.SentimentAnalysis.PositiveCount)
	}

	expectedNeutralCount := 1 // 3 star reviews
	if result.Analysis.SentimentAnalysis.NeutralCount != expectedNeutralCount {
		t.Errorf("Expected NeutralCount %d, got %d", expectedNeutralCount, result.Analysis.SentimentAnalysis.NeutralCount)
	}

	expectedNegativeCount := 1 // 1 and 2 star reviews
	if result.Analysis.SentimentAnalysis.NegativeCount != expectedNegativeCount {
		t.Errorf("Expected NegativeCount %d, got %d", expectedNegativeCount, result.Analysis.SentimentAnalysis.NegativeCount)
	}

	// Verify that percentages are calculated correctly from counts
	expectedPositivePercentage := float64(expectedPositiveCount) / float64(expectedTotalReviews) * 100 // 2/4 * 100 = 50%
	if result.Analysis.SentimentAnalysis.PositivePercentage != expectedPositivePercentage {
		t.Errorf("Expected PositivePercentage %.1f%%, got %.1f%%", expectedPositivePercentage, result.Analysis.SentimentAnalysis.PositivePercentage)
	}

	expectedNeutralPercentage := float64(expectedNeutralCount) / float64(expectedTotalReviews) * 100 // 1/4 * 100 = 25%
	if result.Analysis.SentimentAnalysis.NeutralPercentage != expectedNeutralPercentage {
		t.Errorf("Expected NeutralPercentage %.1f%%, got %.1f%%", expectedNeutralPercentage, result.Analysis.SentimentAnalysis.NeutralPercentage)
	}

	expectedNegativePercentage := float64(expectedNegativeCount) / float64(expectedTotalReviews) * 100 // 1/4 * 100 = 25%
	if result.Analysis.SentimentAnalysis.NegativePercentage != expectedNegativePercentage {
		t.Errorf("Expected NegativePercentage %.1f%%, got %.1f%%", expectedNegativePercentage, result.Analysis.SentimentAnalysis.NegativePercentage)
	}

	// Verify percentages add up to 100%
	totalPercentage := result.Analysis.SentimentAnalysis.PositivePercentage +
		result.Analysis.SentimentAnalysis.NeutralPercentage +
		result.Analysis.SentimentAnalysis.NegativePercentage
	if totalPercentage != 100.0 {
		t.Errorf("Percentages should add up to 100%%, got %.1f%%", totalPercentage)
	}

	// Log for debugging
	t.Logf("Analysis completed - TotalReviews: %d, Metadata ReviewCount: %d, Average Rating: %f",
		result.Analysis.SentimentAnalysis.TotalReviews,
		result.Metadata.ReviewCount,
		result.Analysis.OverallSummary.AverageRating)
	t.Logf("Sentiment percentages - Positive: %.1f%%, Neutral: %.1f%%, Negative: %.1f%%, Total: %.1f%%",
		result.Analysis.SentimentAnalysis.PositivePercentage,
		result.Analysis.SentimentAnalysis.NeutralPercentage,
		result.Analysis.SentimentAnalysis.NegativePercentage,
		totalPercentage)
}
