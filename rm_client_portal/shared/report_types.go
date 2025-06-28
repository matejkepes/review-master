package shared

import (
	"time"
)

// ClientReportData represents the data structure for a client report
type ClientReportData struct {
	ReportID        int64
	ClientID        int
	ClientName      string
	PeriodStart     time.Time
	PeriodEnd       time.Time
	GeneratedAt     time.Time
	LocationResults []AnalysisResult
}

// AnalysisResult represents the structured output of review analysis
type AnalysisResult struct {
	Analysis Analysis         `json:"analysis"`
	Metadata AnalysisMetadata `json:"metadata"`
}

// Analysis contains the main components of the review analysis
type Analysis struct {
	OverallSummary          OverallSummary          `json:"overall_summary"`
	SentimentAnalysis       SentimentAnalysis       `json:"sentiment_analysis"`
	KeyTakeaways            KeyTakeaways            `json:"key_takeaways"`
	NegativeReviewBreakdown NegativeReviewBreakdown `json:"negative_review_breakdown"`
	TrainingRecommendations TrainingRecommendations `json:"training_recommendations"`
}

// OverallSummary provides a high-level overview of the reviews
type OverallSummary struct {
	SummaryText       string   `json:"summary_text"`
	PositiveThemes    []string `json:"positive_themes"`
	NegativeThemes    []string `json:"negative_themes"`
	OverallPerception string   `json:"overall_perception"`
	AverageRating     float64  `json:"average_rating"`
}

// SentimentAnalysis provides statistics about review sentiment
type SentimentAnalysis struct {
	PositiveCount      int     `json:"positive_count"`
	PositivePercentage float64 `json:"positive_percentage"`
	NeutralCount       int     `json:"neutral_count"`
	NeutralPercentage  float64 `json:"neutral_percentage"`
	NegativeCount      int     `json:"negative_count"`
	NegativePercentage float64 `json:"negative_percentage"`
	TotalReviews       int     `json:"total_reviews"`
	SentimentTrend     string  `json:"sentiment_trend"`
}

// KeyTakeaways highlights the main strengths and areas for improvement
type KeyTakeaways struct {
	Strengths           []Insight `json:"strengths"`
	AreasForImprovement []Insight `json:"areas_for_improvement"`
}

// Insight represents a specific finding with supporting details
type Insight struct {
	Category    string `json:"category"`
	Description string `json:"description"`
	Example     string `json:"example"`
}

// NegativeReviewBreakdown categorizes negative reviews
type NegativeReviewBreakdown struct {
	Categories                 []ReviewCategory `json:"categories"`
	ImprovementRecommendations []string         `json:"improvement_recommendations"`
}

// ReviewCategory represents a group of similar negative reviews
type ReviewCategory struct {
	Name       string  `json:"name"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

// TrainingRecommendations provides suggestions for staff development
type TrainingRecommendations struct {
	ForOperators []string `json:"for_operators"`
	ForDrivers   []string `json:"for_drivers"`
}

// AnalysisMetadata contains information about the analysis process
type AnalysisMetadata struct {
	GeneratedAt  time.Time `json:"generated_at"`
	ReviewCount  int       `json:"review_count"`
	LocationID   string    `json:"location_id"`
	LocationName string    `json:"location_name"`
	BusinessName string    `json:"business_name"`
	ReportPeriod struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	} `json:"report_period"`
	ClientID      int    `json:"client_id"`
	AnalyzerID    string `json:"analyzer_id"`
	AnalyzerName  string `json:"analyzer_name"`
	AnalyzerModel string `json:"analyzer_model"`
}
