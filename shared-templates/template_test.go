package shared_templates

import (
	"bytes"
	"html/template"
	"testing"
	"time"
)

func TestMonthlyReportTemplateParses(t *testing.T) {
	// Test that the template parses without errors
	tmpl, err := template.New("monthly_report").Parse(MonthlyReportTemplate)
	if err != nil {
		t.Fatalf("Template parsing failed: %v", err)
	}

	if tmpl == nil {
		t.Fatal("Template is nil after parsing")
	}
}

func TestMonthlyReportTemplateExecutes(t *testing.T) {
	// Create minimal test data
	testData := ClientReportData{
		ReportID:    1,
		ClientID:    1,
		ClientName:  "Test Client",
		PeriodStart: time.Now().AddDate(0, -1, 0),
		PeriodEnd:   time.Now(),
		GeneratedAt: time.Now(),
		LocationResults: []AnalysisResult{
			{
				Metadata: AnalysisMetadata{
					LocationName: "Test Location",
					BusinessName: "Test Business",
				},
				Analysis: Analysis{
					OverallSummary: OverallSummary{
						AverageRating: 4.5,
					},
					SentimentAnalysis: SentimentAnalysis{
						TotalReviews: 10,
					},
				},
			},
		},
	}

	// Test that the template executes without errors
	tmpl, err := template.New("monthly_report").Parse(MonthlyReportTemplate)
	if err != nil {
		t.Fatalf("Template parsing failed: %v", err)
	}

	// Execute template to a buffer
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, testData)
	if err != nil {
		t.Fatalf("Template execution failed: %v", err)
	}
}