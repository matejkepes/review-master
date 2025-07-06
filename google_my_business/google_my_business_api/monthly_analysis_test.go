package google_my_business_api

import (
	"database/sql"
	"errors"
	"google_my_business/ai_service/review"
	"google_my_business/database"
	"shared-templates"
	"net/http"
	"testing"
	"time"
)

// Mock variables
var (
	originalGetAccounts          = GetAccountsFunc
	originalGetLocations         = GetLocationsFunc
	originalFetchReviewsForMonth = FetchReviewsForMonthFunc
	originalConvertHTMLToPDF     = ConvertHTMLToPDFFunc
)

// Mock GetAccounts for testing
func mockGetAccounts(client *http.Client) []string {
	return []string{"account1", "account2"}
}

// Mock GetLocations for testing
// This function is patched by the test to return locations that already have client IDs set
// This eliminates the need for the database lookup inside AnalyzeClientReviews
func mockGetLocations(client *http.Client, account string, db *sql.DB, lookupMode int) []database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode {
	// Return mock locations
	return []database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode{
		{
			GoogleMyBusinessLocationName:                  "Test Location 1",
			GoogleMyBusinessLocationPath:                  "accounts/123/locations/456",
			GoogleMyBusinessPostalCode:                    "12345",
			GoogleMyBusinessReplyToUnspecfifiedStarRating: true,
			GoogleMyBusinessUnspecfifiedStarRatingReply:   "Thank you for your feedback",
			GoogleMyBusinessReplyToOneStarRating:          true,
			GoogleMyBusinessOneStarRatingReply:            "I'm sorry to hear that",
			GoogleMyBusinessReplyToTwoStarRating:          true,
			GoogleMyBusinessTwoStarRatingReply:            "We apologize for the inconvenience",
			GoogleMyBusinessReplyToThreeStarRating:        true,
			GoogleMyBusinessThreeStarRatingReply:          "Thank you for your feedback",
			GoogleMyBusinessReplyToFourStarRating:         true,
			GoogleMyBusinessFourStarRatingReply:           "We appreciate your positive feedback",
			GoogleMyBusinessReplyToFiveStarRating:         true,
			GoogleMyBusinessFiveStarRatingReply:           "Thank you for your wonderful review",
			GoogleMyBusinessReportEnabled:                 true,
			EmailAddress:                                  "test@example.com",
			TimeZone:                                      "UTC",
			ClientID:                                      1,
			MultiMessageSeparator:                         "SSSSS",
		},
		{
			GoogleMyBusinessLocationName:                  "Test Location 2",
			GoogleMyBusinessLocationPath:                  "accounts/123/locations/789",
			GoogleMyBusinessPostalCode:                    "67890",
			GoogleMyBusinessReplyToUnspecfifiedStarRating: true,
			GoogleMyBusinessUnspecfifiedStarRatingReply:   "Thank you for your feedback",
			GoogleMyBusinessReplyToOneStarRating:          true,
			GoogleMyBusinessOneStarRatingReply:            "I'm sorry to hear that",
			GoogleMyBusinessReplyToTwoStarRating:          true,
			GoogleMyBusinessTwoStarRatingReply:            "We apologize for the inconvenience",
			GoogleMyBusinessReplyToThreeStarRating:        true,
			GoogleMyBusinessThreeStarRatingReply:          "Thank you for your feedback",
			GoogleMyBusinessReplyToFourStarRating:         true,
			GoogleMyBusinessFourStarRatingReply:           "We appreciate your positive feedback",
			GoogleMyBusinessReplyToFiveStarRating:         true,
			GoogleMyBusinessFiveStarRatingReply:           "Thank you for your wonderful review",
			GoogleMyBusinessReportEnabled:                 true,
			EmailAddress:                                  "test@example.com",
			TimeZone:                                      "UTC",
			ClientID:                                      1,
			MultiMessageSeparator:                         "SSSSS",
		},
	}
}

// Mock FetchReviewsForMonth for testing
func mockFetchReviewsForMonth(client *http.Client, location database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode, targetMonth time.Time) ([]map[string]interface{}, time.Time, time.Time, error) {
	startTime := time.Date(targetMonth.Year(), targetMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(targetMonth.Year(), targetMonth.Month()+1, 1, 0, 0, 0, 0, time.UTC)

	// Create a few mock reviews
	reviews := []map[string]interface{}{
		{
			"name": "accounts/123/locations/456/reviews/review1",
			"comment": map[string]interface{}{
				"text": "Great service, very happy with the taxi ride!",
			},
			"starRating": map[string]interface{}{
				"count": float64(5),
			},
			"createTime": startTime.Add(24 * time.Hour).Format(time.RFC3339),
		},
		{
			"name": "accounts/123/locations/456/reviews/review2",
			"comment": map[string]interface{}{
				"text": "Driver was very polite but arrived a bit late.",
			},
			"starRating": map[string]interface{}{
				"count": float64(4),
			},
			"createTime": startTime.Add(48 * time.Hour).Format(time.RFC3339),
		},
	}

	// Return an error for a specific location to test error handling
	if location.GoogleMyBusinessLocationPath == "accounts/123/locations/error" {
		return nil, startTime, endTime, errors.New("error fetching reviews")
	}

	// Return empty reviews for a specific location to test empty reviews handling
	if location.GoogleMyBusinessLocationPath == "accounts/123/locations/empty" {
		return []map[string]interface{}{}, startTime, endTime, nil
	}

	return reviews, startTime, endTime, nil
}

// Mock ConvertHTMLToPDF for testing
func mockConvertHTMLToPDF(htmlContent string) ([]byte, error) {
	// Check for error case
	if htmlContent == "error_case" {
		return nil, errors.New("pdf generation error")
	}

	// Return mock PDF data
	return []byte("mock PDF content"), nil
}

// Reset mocks after tests
func resetMocks() {
	GetAccountsFunc = originalGetAccounts
	GetLocationsFunc = originalGetLocations
	FetchReviewsForMonthFunc = originalFetchReviewsForMonth
	ConvertHTMLToPDFFunc = originalConvertHTMLToPDF
}

// MockReviewAnalyzer represents a mock review analyzer
type MockReviewAnalyzer struct {
	analyzeWasCalled bool
	returnError      bool
}

// Analyze is a mock implementation
func (m *MockReviewAnalyzer) Analyze(batch review.ReviewBatch) (*shared_templates.AnalysisResult, error) {
	m.analyzeWasCalled = true

	if m.returnError {
		return nil, errors.New("analyzer error")
	}

	// Create a simple analysis result
	return &shared_templates.AnalysisResult{
		Analysis: shared_templates.Analysis{
			OverallSummary: shared_templates.OverallSummary{
				SummaryText:       "This is a test summary",
				PositiveThemes:    []string{"Service", "Punctuality"},
				NegativeThemes:    []string{"Communication"},
				OverallPerception: "Positive",
				AverageRating:     4.5,
			},
			SentimentAnalysis: shared_templates.SentimentAnalysis{
				PositiveCount:      1,
				PositivePercentage: 50.0,
				NeutralCount:       1,
				NeutralPercentage:  50.0,
				NegativeCount:      0,
				NegativePercentage: 0.0,
				TotalReviews:       2,
				SentimentTrend:     "Stable",
			},
		},
		Metadata: shared_templates.AnalysisMetadata{
			GeneratedAt:   time.Now(),
			ReviewCount:   2,
			LocationID:    batch.LocationID,
			LocationName:  batch.LocationName,
			BusinessName:  batch.BusinessName,
			ClientID:      batch.ClientID,
			AnalyzerID:    "mock-analyzer",
			AnalyzerName:  "Mock Analyzer",
			AnalyzerModel: "mock-model",
		},
	}, nil
}

// MockDB implements a minimal DB interface for testing
type MockDB struct {
	clients        []database.ClientWithMonthlyReviewAnalysis
	existingReport *shared_templates.ClientReportData
	deleteError    error
	saveError      error
	savedReportID  int64

	// Track if delete and save were called
	deleteWasCalled bool
	saveWasCalled   bool
	savedData       []byte
}

// Mock DB.Client() method
func (m *MockDB) Client() *sql.DB {
	// Create a mock DB connection that won't be used but will pass nil checks
	db, _ := sql.Open("mysql", "mock:mock@/mock")
	return db
}

func (m *MockDB) GetClientsWithMonthlyReviewAnalysisEnabled() ([]database.ClientWithMonthlyReviewAnalysis, error) {
	return m.clients, nil
}

func (m *MockDB) GetClientReportByClientAndPeriod(clientID int, periodStart, periodEnd time.Time) (*shared_templates.ClientReportData, error) {
	return m.existingReport, nil
}

func (m *MockDB) DeleteClientReport(reportID int64) error {
	m.deleteWasCalled = true
	return m.deleteError
}

func (m *MockDB) SaveClientReport(clientID int, periodStart, periodEnd time.Time, locationResults []byte) (int64, error) {
	m.saveWasCalled = true
	m.savedData = locationResults
	return m.savedReportID, m.saveError
}

// MockEmailService implements email_service.EmailService
type MockEmailService struct {
	returnError     bool
	sendWasCalled   bool
	sentEmails      int
	lastPDFContent  []byte
	lastClientName  string
	lastEmail       string
	lastMonth       string
	lastSubject     string
	lastTextContent string
}

// SendMonthlyReport is a mock implementation
func (m *MockEmailService) SendMonthlyReport(clientName string, emailAddress string, month string, pdfReport []byte) error {
	m.sendWasCalled = true
	m.lastClientName = clientName
	m.lastEmail = emailAddress
	m.lastMonth = month
	m.lastPDFContent = pdfReport

	if m.returnError {
		return errors.New("mock email sending error")
	}

	m.sentEmails++
	return nil
}

// SendPlainTextEmail is a mock implementation for the new method
func (m *MockEmailService) SendPlainTextEmail(subject string, recipient string, recipientName string, textContent string) error {
	m.sendWasCalled = true
	m.lastSubject = subject
	m.lastEmail = recipient
	m.lastClientName = recipientName
	m.lastTextContent = textContent

	if m.returnError {
		return errors.New("mock email sending error")
	}

	m.sentEmails++
	return nil
}

func TestAnalyzeClientReviews(t *testing.T) {
	tests := []struct {
		name                string
		db                  *MockDB
		forceReprocess      bool
		analyzerReturnError bool
		pdfGenerationError  bool
		emailServiceError   bool
		wantError           bool
		expectDeleteCalled  bool
		expectSaveCalled    bool
		expectSucceeded     int
		expectFailed        int
		expectSkipped       int
		expectReviewCount   int
		expectPDFCount      int
		expectEmailCount    int
		skipRun             bool
	}{
		{
			name: "successful processing with PDF",
			db: &MockDB{
				clients: []database.ClientWithMonthlyReviewAnalysis{
					{
						ClientID:     1,
						ClientName:   "Test Client",
						EmailAddress: "test@example.com",
						ReportEmailAddress: sql.NullString{
							String: "reports1@example.com",
							Valid:  true,
						},
					},
				},
				savedReportID: 123,
			},
			forceReprocess:      false,
			analyzerReturnError: false,
			pdfGenerationError:  false,
			emailServiceError:   false,
			wantError:           false,
			expectDeleteCalled:  false,
			expectSaveCalled:    true,
			expectSucceeded:     1,
			expectFailed:        0,
			expectSkipped:       0,
			expectReviewCount:   8, // 2 reviews * 2 locations * 2 accounts
			expectPDFCount:      1,
			expectEmailCount:    1, // Expect one email to be sent
		},
		{
			name: "skip existing report",
			db: &MockDB{
				clients: []database.ClientWithMonthlyReviewAnalysis{
					{
						ClientID:     1,
						ClientName:   "Test Client",
						EmailAddress: "test@example.com",
						ReportEmailAddress: sql.NullString{
							String: "reports1@example.com",
							Valid:  true,
						},
					},
				},
				existingReport: &shared_templates.ClientReportData{
					ReportID: 456,
				},
			},
			forceReprocess:      false,
			analyzerReturnError: false,
			wantError:           false,
			expectDeleteCalled:  false,
			expectSaveCalled:    false,
			expectSucceeded:     0,
			expectFailed:        0,
			expectSkipped:       1,
			expectReviewCount:   0,
		},
		{
			name: "force_reprocess",
			db: &MockDB{
				clients: []database.ClientWithMonthlyReviewAnalysis{
					{
						ClientID:     1,
						ClientName:   "Test Client",
						EmailAddress: "test@example.com",
						ReportEmailAddress: sql.NullString{
							String: "reports1@example.com",
							Valid:  true,
						},
					},
				},
				existingReport: &shared_templates.ClientReportData{
					ReportID:    123,
					ClientID:    1,
					GeneratedAt: time.Now(),
				},
				savedReportID: 456,
			},
			forceReprocess:      true,
			analyzerReturnError: false,
			pdfGenerationError:  false,
			emailServiceError:   false,
			wantError:           false,
			expectDeleteCalled:  true,
			expectSaveCalled:    true,
			expectSucceeded:     1,
			expectFailed:        0,
			expectSkipped:       0,
			expectReviewCount:   8,
			expectPDFCount:      1, // Expect PDF to be generated
			expectEmailCount:    1, // Expect email to be sent
		},
		{
			name: "no clients",
			db: &MockDB{
				clients: []database.ClientWithMonthlyReviewAnalysis{},
			},
			forceReprocess:      false,
			analyzerReturnError: false,
			wantError:           false,
			expectDeleteCalled:  false,
			expectSaveCalled:    false,
			expectSucceeded:     0,
			expectFailed:        0,
			expectSkipped:       0,
			expectReviewCount:   0,
			skipRun:             true,
		},
		{
			name: "error saving report",
			db: &MockDB{
				clients: []database.ClientWithMonthlyReviewAnalysis{
					{
						ClientID:     1,
						ClientName:   "Test Client",
						EmailAddress: "test@example.com",
						ReportEmailAddress: sql.NullString{
							String: "reports1@example.com",
							Valid:  true,
						},
					},
				},
				saveError: errors.New("failed to save report"),
			},
			forceReprocess:      false,
			analyzerReturnError: false,
			wantError:           false, // The error is captured in the summary, not returned
			expectDeleteCalled:  false,
			expectSaveCalled:    true,
			expectSucceeded:     0, // The function now properly marks as failed if save errors
			expectFailed:        1,
			expectSkipped:       0,
			expectReviewCount:   8, // 2 reviews * 2 locations * 2 accounts
		},
		{
			name: "error deleting report",
			db: &MockDB{
				clients: []database.ClientWithMonthlyReviewAnalysis{
					{
						ClientID:     1,
						ClientName:   "Test Client",
						EmailAddress: "test@example.com",
						ReportEmailAddress: sql.NullString{
							String: "reports1@example.com",
							Valid:  true,
						},
					},
				},
				existingReport: &shared_templates.ClientReportData{
					ReportID: 456,
				},
				deleteError: errors.New("failed to delete report"),
			},
			forceReprocess:      true,
			analyzerReturnError: false,
			wantError:           false, // The error is captured in the summary, not returned
			expectDeleteCalled:  true,
			expectSaveCalled:    false,
			expectSucceeded:     0,
			expectFailed:        1,
			expectSkipped:       0,
			expectReviewCount:   0,
		},
		{
			name: "analyzer error",
			db: &MockDB{
				clients: []database.ClientWithMonthlyReviewAnalysis{
					{
						ClientID:     1,
						ClientName:   "Test Client",
						EmailAddress: "test@example.com",
						ReportEmailAddress: sql.NullString{
							String: "reports1@example.com",
							Valid:  true,
						},
					},
				},
			},
			forceReprocess:      false,
			analyzerReturnError: true,
			wantError:           false, // The error is captured in the summary, not returned
			expectDeleteCalled:  false,
			expectSaveCalled:    false, // No save called when analyzer fails
			expectSucceeded:     0,
			expectFailed:        1,
			expectSkipped:       0,
			expectReviewCount:   0,
		},
		{
			name: "successful processing with PDF error",
			db: &MockDB{
				clients: []database.ClientWithMonthlyReviewAnalysis{
					{
						ClientID:     1,
						ClientName:   "Test Client",
						EmailAddress: "test@example.com",
						ReportEmailAddress: sql.NullString{
							String: "reports1@example.com",
							Valid:  true,
						},
					},
				},
				savedReportID: 123,
			},
			forceReprocess:      false,
			analyzerReturnError: false,
			pdfGenerationError:  true,
			wantError:           false,
			expectDeleteCalled:  false,
			expectSaveCalled:    true,
			expectSucceeded:     1, // Still succeeds because PDF failure is not critical
			expectFailed:        0,
			expectSkipped:       0,
			expectReviewCount:   8,
			expectPDFCount:      0, // No PDFs generated due to error
		},
		{
			name: "successful processing with email error",
			db: &MockDB{
				clients: []database.ClientWithMonthlyReviewAnalysis{
					{
						ClientID:     1,
						ClientName:   "Test Client",
						EmailAddress: "test@example.com",
						ReportEmailAddress: sql.NullString{
							String: "reports1@example.com",
							Valid:  true,
						},
					},
				},
				savedReportID: 123,
			},
			forceReprocess:      false,
			analyzerReturnError: false,
			pdfGenerationError:  false,
			emailServiceError:   true,
			wantError:           false,
			expectDeleteCalled:  false,
			expectSaveCalled:    true,
			expectSucceeded:     1, // Still succeeds because email failure is not critical
			expectFailed:        0,
			expectSkipped:       0,
			expectReviewCount:   8,
			expectPDFCount:      1,
			expectEmailCount:    0, // No emails sent due to error
		},
		{
			name: "missing report email address",
			db: &MockDB{
				clients: []database.ClientWithMonthlyReviewAnalysis{
					{
						ClientID:     1,
						ClientName:   "Test Client",
						EmailAddress: "test@example.com", // This will be ignored
						ReportEmailAddress: sql.NullString{
							Valid: false, // Null in database
						},
					},
				},
				savedReportID: 123,
			},
			forceReprocess:      false,
			analyzerReturnError: false,
			pdfGenerationError:  false,
			emailServiceError:   false,
			wantError:           false,
			expectDeleteCalled:  false,
			expectSaveCalled:    true,
			expectSucceeded:     1,
			expectFailed:        0,
			expectSkipped:       0,
			expectReviewCount:   8,
			expectPDFCount:      1,
			expectEmailCount:    0, // No email should be sent if the report_email_address is missing
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mocks
			GetAccountsFunc = mockGetAccounts
			GetLocationsFunc = mockGetLocations
			FetchReviewsForMonthFunc = mockFetchReviewsForMonth

			// Set up PDF mock
			if tt.pdfGenerationError {
				ConvertHTMLToPDFFunc = func(htmlContent string) ([]byte, error) {
					return nil, errors.New("mock PDF generation error")
				}
			} else {
				ConvertHTMLToPDFFunc = mockConvertHTMLToPDF
			}

			// Create a special version of GetLocations that directly returns locations
			// with the correct client ID, bypassing the database lookup
			GetLocationsFunc = func(client *http.Client, account string, db *sql.DB, lookupMode int) []database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode {
				// For tests, ensure locations have the correct client ID already set
				clientID := uint64(1) // This matches our test client ID
				return []database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode{
					{
						GoogleMyBusinessLocationName:  "Location 1",
						GoogleMyBusinessLocationPath:  "accounts/123/locations/456",
						GoogleMyBusinessPostalCode:    "12345",
						ClientID:                      clientID,
						GoogleMyBusinessReportEnabled: true,
						TimeZone:                      "UTC",
					},
					{
						GoogleMyBusinessLocationName:  "Location 2",
						GoogleMyBusinessLocationPath:  "accounts/123/locations/789",
						GoogleMyBusinessPostalCode:    "67890",
						ClientID:                      clientID,
						GoogleMyBusinessReportEnabled: true,
						TimeZone:                      "UTC",
					},
				}
			}

			defer resetMocks()

			// Set up test parameters
			httpClient := &http.Client{}
			targetMonth := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)
			targetMonthPtr := &targetMonth
			mockAnalyzer := &MockReviewAnalyzer{
				returnError: tt.analyzerReturnError,
			}

			// Create mock email service
			mockEmailSvc := &MockEmailService{
				returnError: tt.emailServiceError,
			}

			// Run test
			summary, err := AnalyzeClientReviews(tt.db, httpClient, mockAnalyzer, targetMonthPtr, tt.forceReprocess, mockEmailSvc)

			// Skip run checks
			if tt.skipRun {
				return
			}

			// Check error
			if (err != nil) != tt.wantError {
				t.Errorf("AnalyzeClientReviews() error = %v, wantError %v", err, tt.wantError)
			}

			// Check if methods were called
			if tt.db.deleteWasCalled != tt.expectDeleteCalled {
				t.Errorf("DeleteClientReport called = %v, expectDeleteCalled %v", tt.db.deleteWasCalled, tt.expectDeleteCalled)
			}

			if tt.db.saveWasCalled != tt.expectSaveCalled {
				t.Errorf("SaveClientReport called = %v, expectSaveCalled %v", tt.db.saveWasCalled, tt.expectSaveCalled)
			}

			// Check summary
			if summary.ClientsSucceeded != tt.expectSucceeded {
				t.Errorf("ClientsSucceeded = %v, expectSucceeded %v", summary.ClientsSucceeded, tt.expectSucceeded)
			}

			if summary.ClientsFailed != tt.expectFailed {
				t.Errorf("ClientsFailed = %v, expectFailed %v", summary.ClientsFailed, tt.expectFailed)
			}

			if summary.ClientsSkipped != tt.expectSkipped {
				t.Errorf("ClientsSkipped = %v, expectSkipped %v", summary.ClientsSkipped, tt.expectSkipped)
			}

			if summary.TotalReviews != tt.expectReviewCount && !tt.analyzerReturnError {
				t.Errorf("TotalReviews = %v, expectReviewCount %v", summary.TotalReviews, tt.expectReviewCount)
			}

			// Check PDF generation count
			if summary.PDFsGenerated != tt.expectPDFCount {
				t.Errorf("PDFsGenerated = %v, expectPDFCount %v", summary.PDFsGenerated, tt.expectPDFCount)
			}

			// Check email sent count
			if summary.EmailsSent != tt.expectEmailCount {
				t.Errorf("EmailsSent = %v, expectEmailCount %v", summary.EmailsSent, tt.expectEmailCount)
			}

			// For successful cases, verify the saved data
			if tt.expectSaveCalled && tt.expectSucceeded > 0 {
				if len(tt.db.savedData) == 0 {
					t.Errorf("Expected saved data but none was found")
				}

				// We could add more detailed validation of the saved JSON data here if needed
			}

			// For successful email cases, check if email service was called properly
			if tt.expectEmailCount > 0 {
				if !mockEmailSvc.sendWasCalled {
					t.Error("Expected email service to be called, but it wasn't")
				}

				if mockEmailSvc.sentEmails != tt.expectEmailCount {
					t.Errorf("Email service sent %d emails, expected %d", mockEmailSvc.sentEmails, tt.expectEmailCount)
				}
			}
		})
	}
}

// Additional tests for more specific scenarios would be added here

// TestSendReportEmail specifically tests the sendReportEmail function
func TestSendReportEmail(t *testing.T) {
	tests := []struct {
		name         string
		clientInfo   database.ClientWithMonthlyReviewAnalysis
		pdfContent   []byte
		mockEmail    *MockEmailService
		expectError  bool
		expectCalled bool
	}{
		{
			name: "successful_email",
			clientInfo: database.ClientWithMonthlyReviewAnalysis{
				ClientID:     1,
				ClientName:   "Test Client",
				EmailAddress: "test@example.com", // This will be ignored
				ReportEmailAddress: sql.NullString{
					Valid:  true,
					String: "reports1@example.com", // This should be used
				},
			},
			pdfContent:   []byte("mock PDF content"),
			mockEmail:    &MockEmailService{},
			expectError:  false,
			expectCalled: true,
		},
		{
			name: "empty_email_address",
			clientInfo: database.ClientWithMonthlyReviewAnalysis{
				ClientID:     1,
				ClientName:   "Test Client",
				EmailAddress: "", // This will be ignored
				// No ReportEmailAddress set - should result in an error
				ReportEmailAddress: sql.NullString{
					Valid: false,
				},
			},
			pdfContent:   []byte("mock PDF content"),
			mockEmail:    &MockEmailService{},
			expectError:  true,
			expectCalled: false,
		},
		{
			name: "empty_PDF_content",
			clientInfo: database.ClientWithMonthlyReviewAnalysis{
				ClientID:     1,
				ClientName:   "Test Client",
				EmailAddress: "test@example.com", // This will be ignored
				ReportEmailAddress: sql.NullString{
					Valid:  true,
					String: "reports1@example.com", // Not used due to empty PDF
				},
			},
			pdfContent:   []byte{},
			mockEmail:    &MockEmailService{},
			expectError:  true,
			expectCalled: false,
		},
		{
			name: "email_service_error",
			clientInfo: database.ClientWithMonthlyReviewAnalysis{
				ClientID:     1,
				ClientName:   "Test Client",
				EmailAddress: "test@example.com", // This will be ignored
				ReportEmailAddress: sql.NullString{
					Valid:  true,
					String: "reports1@example.com", // This will be used
				},
			},
			pdfContent:   []byte("mock PDF content"),
			mockEmail:    &MockEmailService{returnError: true},
			expectError:  true,
			expectCalled: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Reset mock for each test
			tc.mockEmail.sendWasCalled = false
			tc.mockEmail.sentEmails = 0

			// Call the function
			err := sendReportEmail(tc.mockEmail, tc.clientInfo, time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC), tc.pdfContent)

			// Check error expectation
			if (err != nil) != tc.expectError {
				t.Errorf("sendReportEmail() error = %v, expectError %v", err, tc.expectError)
			}

			// Check if email service was called as expected
			if tc.mockEmail.sendWasCalled != tc.expectCalled {
				t.Errorf("Email service called = %v, expectCalled %v", tc.mockEmail.sendWasCalled, tc.expectCalled)
			}

			// Additional checks if we expect the call to succeed
			if tc.expectCalled && !tc.expectError {
				// Check if the correct email address was used
				expectedEmail := tc.clientInfo.ReportEmailAddress.String
				if tc.mockEmail.lastEmail != expectedEmail {
					t.Errorf("Wrong email address: got %s, want %s", tc.mockEmail.lastEmail, expectedEmail)
				}

				// Check if the correct client name was used
				if tc.mockEmail.lastClientName != tc.clientInfo.ClientName {
					t.Errorf("Wrong client name: got %s, want %s", tc.mockEmail.lastClientName, tc.clientInfo.ClientName)
				}

				// Check if the month is correct
				expectedMonth := "March 2025"
				if tc.mockEmail.lastMonth != expectedMonth {
					t.Errorf("Wrong month: got %s, want %s", tc.mockEmail.lastMonth, expectedMonth)
				}
			}
		})
	}
}

func TestSendReportEmailWithMissingReportEmailAddress(t *testing.T) {
	mockEmail := &MockEmailService{}

	// Test case with null report_email_address
	err := sendReportEmail(
		mockEmail,
		database.ClientWithMonthlyReviewAnalysis{
			ClientID:     1,
			ClientName:   "Test Client",
			EmailAddress: "test@example.com", // This is set but should be ignored
			ReportEmailAddress: sql.NullString{
				Valid: false, // Null in database
			},
		},
		time.Now(),
		[]byte("mock PDF content"),
	)

	// Should have error and no email sent
	if err == nil {
		t.Error("Expected error when report_email_address is null")
	}
	if mockEmail.sendWasCalled {
		t.Error("Email service should not be called when report_email_address is null")
	}

	// Reset mock
	mockEmail = &MockEmailService{}

	// Test case with empty report_email_address
	err = sendReportEmail(
		mockEmail,
		database.ClientWithMonthlyReviewAnalysis{
			ClientID:     1,
			ClientName:   "Test Client",
			EmailAddress: "test@example.com", // This is set but should be ignored
			ReportEmailAddress: sql.NullString{
				Valid:  true,
				String: "", // Empty string
			},
		},
		time.Now(),
		[]byte("mock PDF content"),
	)

	// Should have error and no email sent
	if err == nil {
		t.Error("Expected error when report_email_address is empty")
	}
	if mockEmail.sendWasCalled {
		t.Error("Email service should not be called when report_email_address is empty")
	}
}

func TestSendReportEmailWithValidReportEmailAddress(t *testing.T) {
	mockEmail := &MockEmailService{}

	// Test case with valid report_email_address
	clientInfo := database.ClientWithMonthlyReviewAnalysis{
		ClientID:     1,
		ClientName:   "Test Client",
		EmailAddress: "test@example.com", // This should be ignored
		ReportEmailAddress: sql.NullString{
			Valid:  true,
			String: "reports@example.com", // This should be used
		},
	}

	periodStart := time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC)
	err := sendReportEmail(mockEmail, clientInfo, periodStart, []byte("mock PDF content"))

	// Should have no error and email sent
	if err != nil {
		t.Errorf("Unexpected error when report_email_address is valid: %v", err)
	}
	if !mockEmail.sendWasCalled {
		t.Error("Email service should be called when report_email_address is valid")
	}
	if mockEmail.lastEmail != "reports@example.com" {
		t.Errorf("Wrong email address used. Expected 'reports@example.com', got '%s'", mockEmail.lastEmail)
	}
	if mockEmail.lastClientName != clientInfo.ClientName {
		t.Errorf("Wrong client name used. Expected '%s', got '%s'", clientInfo.ClientName, mockEmail.lastClientName)
	}
	if mockEmail.lastMonth != "April 2023" {
		t.Errorf("Wrong month format. Expected 'April 2023', got '%s'", mockEmail.lastMonth)
	}
}
