package email_service

import (
	"errors"
	"testing"

	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// MockSendGridClient implements a mock SendGrid client for testing
type MockSendGridClient struct {
	sendFunc func(*mail.SGMailV3) (*Response, error)
}

func (m *MockSendGridClient) Send(email *mail.SGMailV3) (*Response, error) {
	if m.sendFunc == nil {
		return nil, errors.New("sendFunc not implemented")
	}
	return m.sendFunc(email)
}

func TestSendMonthlyReport(t *testing.T) {
	tests := []struct {
		name          string
		clientName    string
		emailAddress  string
		month         string
		pdfReport     []byte
		sendFunc      func(*mail.SGMailV3) (*Response, error)
		expectedError string
	}{
		{
			name:         "successful send",
			clientName:   "Test Client",
			emailAddress: "test@example.com",
			month:        "March 2024",
			pdfReport:    []byte("test pdf content"),
			sendFunc: func(email *mail.SGMailV3) (*Response, error) {
				return &Response{StatusCode: 202}, nil
			},
			expectedError: "",
		},
		{
			name:          "empty email address",
			clientName:    "Test Client",
			emailAddress:  "",
			month:         "March 2024",
			pdfReport:     []byte("test pdf content"),
			sendFunc:      nil,
			expectedError: "email address cannot be empty",
		},
		{
			name:          "empty PDF report",
			clientName:    "Test Client",
			emailAddress:  "test@example.com",
			month:         "March 2024",
			pdfReport:     []byte{},
			sendFunc:      nil,
			expectedError: "PDF report cannot be empty",
		},
		{
			name:         "sendgrid API error",
			clientName:   "Test Client",
			emailAddress: "test@example.com",
			month:        "March 2024",
			pdfReport:    []byte("test pdf content"),
			sendFunc: func(email *mail.SGMailV3) (*Response, error) {
				return nil, errors.New("API error")
			},
			expectedError: "failed to send email: API error",
		},
		{
			name:         "non-2xx status code",
			clientName:   "Test Client",
			emailAddress: "test@example.com",
			month:        "March 2024",
			pdfReport:    []byte("test pdf content"),
			sendFunc: func(email *mail.SGMailV3) (*Response, error) {
				return &Response{StatusCode: 400, Body: ""}, nil
			},
			expectedError: "sendgrid API returned non-2xx status code: 400, body: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock client
			mockClient := &MockSendGridClient{
				sendFunc: tt.sendFunc,
			}

			// Create service with mock client
			service := &SendGridEmailService{
				client:     mockClient,
				fromEmail:  "noreply@example.com",
				fromName:   "Review Master",
				templateID: "d-123456789",
			}

			// Send report
			err := service.SendMonthlyReport(tt.clientName, tt.emailAddress, tt.month, tt.pdfReport)

			// Check error
			if tt.expectedError == "" {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected error %q, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError {
					t.Errorf("expected error %q, got %q", tt.expectedError, err.Error())
				}
			}
		})
	}
}
