package google_my_business_api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ConvertHTMLToPDF converts HTML content to PDF using an external API
// It takes HTML content as a string and returns the PDF as a byte slice
func ConvertHTMLToPDF(htmlContent string) ([]byte, error) {
	// API configuration
	apiURL := "https://icixfqne85.execute-api.eu-west-2.amazonaws.com/prod/generate"
	apiKey := "vTdgRaAWIlaIDXxleBhoxaRU2DZjKJKm1P9Rkg96"

	// Create request to PDF generation API
	req, err := http.NewRequest("POST", apiURL, bytes.NewReader([]byte(htmlContent)))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("Content-Type", "text/html")
	req.Header.Set("Accept", "application/pdf")

	// Send request
	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request to PDF API: %v", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("PDF API returned non-OK status: %d", resp.StatusCode)
	}

	// Check content length
	contentLength := resp.ContentLength
	if contentLength < 100 {
		return nil, fmt.Errorf("PDF content length too small: %d bytes", contentLength)
	}

	// Read response body
	pdfData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading PDF response: %v", err)
	}

	return pdfData, nil
}
