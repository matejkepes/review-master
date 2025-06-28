//go:build test_ai_responses

package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"google_my_business/ai_service"
	"google_my_business/ai_service/review"
	"google_my_business/config"
)

// This is the entry point for the command-line tool
func main() {
	testAIResponses()
}

func testAIResponses() {
	// Get CSV file path from command line or environment
	csvPath := os.Getenv("CSV_PATH")
	if len(os.Args) > 1 {
		csvPath = os.Args[1]
	}
	if csvPath == "" {
		csvPath = "Drive Doncaster_2091990859.csv"
	}

	// Get output file path
	outputPath := os.Getenv("OUTPUT_CSV_PATH")
	if outputPath == "" {
		outputPath = strings.TrimSuffix(csvPath, ".csv") + "_with_responses.csv"
	}

	// Check if OpenAI API key is set
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: OPENAI_API_KEY environment variable is not set")
		os.Exit(1)
	}

	// Process the CSV file
	err := processCSV(csvPath, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully processed reviews and saved to %s\n", outputPath)
}

func processCSV(csvPath, outputPath string) error {
	// Read CSV file
	file, err := os.Open(csvPath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Configure the reader to handle double quotes
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV file: %w", err)
	}

	if len(records) < 2 {
		return fmt.Errorf("CSV file has insufficient data")
	}

	// Find column indexes
	headers := records[0]
	reviewTextIdx := findColumnIndex(headers, "Review")
	ratingIdx := findColumnIndex(headers, "Stars")
	reviewerNameIdx := findColumnIndex(headers, "User_Name")

	if reviewTextIdx == -1 || ratingIdx == -1 {
		return fmt.Errorf("required columns not found in CSV. Looking for 'Review' and 'Stars' columns")
	}

	// Add AI Response column if it doesn't exist
	aiResponseIdx := findColumnIndex(headers, "AI Response")
	if aiResponseIdx == -1 {
		headers = append(headers, "AI Response")
		aiResponseIdx = len(headers) - 1
		records[0] = headers
	}

	// Get business info from environment or use defaults
	businessName := getEnvOrDefault("BUSINESS_NAME", "Drive Doncaster")
	location := getEnvOrDefault("LOCATION", "Doncaster")
	contactMethod := getEnvOrDefault("CONTACT_METHOD", "email us at info@drive.co.uk")

	// Get configuration from properties file
	props := config.ReadProperties()

	// Create OpenAI provider using the centralized configuration
	provider, err := ai_service.NewOpenAIProvider(ai_service.GetConfigForUseCase(ai_service.ReviewResponse, props.OpenAIAPIKey))
	if err != nil {
		return fmt.Errorf("failed to create OpenAI provider: %w", err)
	}

	// Create generator using existing code
	generator := review.NewGenerator(provider, review.GeneratorConfig{})

	// Process each review
	for i := 1; i < len(records); i++ {
		row := records[i]

		// Ensure row has enough columns
		for len(row) < len(headers) {
			row = append(row, "")
		}

		// Get review text and rating
		reviewText := ""
		if reviewTextIdx < len(row) {
			reviewText = row[reviewTextIdx]
		}

		// Get rating
		rating := 0
		if ratingIdx < len(row) {
			ratingStr := row[ratingIdx]
			ratingVal, err := strconv.Atoi(ratingStr)
			if err == nil {
				rating = ratingVal
			}
		}

		// Get reviewer name
		reviewerName := ""
		if reviewerNameIdx != -1 && reviewerNameIdx < len(row) {
			reviewerName = row[reviewerNameIdx]
		}

		// Create review context
		reviewCtx := review.ReviewContext{
			Text:          reviewText,
			Rating:        rating,
			Author:        reviewerName,
			BusinessName:  businessName,
			Location:      location,
			ContactMethod: contactMethod,
		}

		// Generate response using the existing generator
		response, err := generator.Generate(reviewCtx)
		if err != nil {
			fmt.Printf("Error generating response for row %d: %v\n", i, err)
			response = fmt.Sprintf("Error: %v", err)
		}

		fmt.Printf("Response: %s\n", response)

		// Add response to the row
		row[aiResponseIdx] = response
		records[i] = row
	}

	// Write updated CSV
	return writeCSV(outputPath, records)
}

func writeCSV(filePath string, records [][]string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return writer.WriteAll(records)
}

func findColumnIndex(headers []string, name string) int {
	for i, header := range headers {
		if strings.EqualFold(header, name) {
			return i
		}
	}
	return -1
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "..."
}
