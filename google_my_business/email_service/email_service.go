package email_service

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// EmailService defines the interface for sending emails
type EmailService interface {
	SendMonthlyReport(clientName string, emailAddress string, month string, pdfReport []byte) error
	SendPlainTextEmail(subject string, recipient string, recipientName string, textContent string) error
}

// SendGridClient defines the interface for the SendGrid client
type SendGridClient interface {
	Send(email *mail.SGMailV3) (*Response, error)
}

// Response represents a SendGrid API response
type Response struct {
	StatusCode int
	Body       string
}

// SendGridEmailService implements EmailService using SendGrid
type SendGridEmailService struct {
	client     SendGridClient
	fromEmail  string
	fromName   string
	templateID string
}

// NewSendGridEmailService creates a new SendGrid email service
func NewSendGridEmailService(apiKey string, fromEmail string, fromName string, templateID string) *SendGridEmailService {
	return &SendGridEmailService{
		client:     &RealSendGridClient{apiKey: apiKey},
		fromEmail:  fromEmail,
		fromName:   fromName,
		templateID: templateID,
	}
}

// RealSendGridClient implements SendGridClient using the real SendGrid client
type RealSendGridClient struct {
	apiKey string
}

func (c *RealSendGridClient) Send(email *mail.SGMailV3) (*Response, error) {
	client := sendgrid.NewSendClient(c.apiKey)
	resp, err := client.Send(email)
	if err != nil {
		return nil, err
	}
	return &Response{StatusCode: resp.StatusCode, Body: resp.Body}, nil
}

// SendMonthlyReport sends a monthly review analysis report to a client
func (s *SendGridEmailService) SendMonthlyReport(clientName string, emailAddress string, month string, pdfReport []byte) error {
	if emailAddress == "" {
		return errors.New("email address cannot be empty")
	}

	if len(pdfReport) == 0 {
		return errors.New("PDF report cannot be empty")
	}

	// Create email
	from := mail.NewEmail(s.fromName, s.fromEmail)
	to := mail.NewEmail(clientName, emailAddress)

	// Create email message
	message := mail.NewV3Mail()
	message.SetFrom(from)
	message.SetTemplateID(s.templateID)

	// Add personalization
	p := mail.NewPersonalization()
	p.AddTos(to)
	p.SetDynamicTemplateData("client_name", clientName)
	p.SetDynamicTemplateData("month", month)
	message.AddPersonalizations(p)

	// Add PDF attachment
	attachment := mail.NewAttachment()
	// Properly base64 encode the PDF content
	encodedContent := base64.StdEncoding.EncodeToString(pdfReport)
	attachment.SetContent(encodedContent)
	attachment.SetType("application/pdf")
	attachment.SetFilename(fmt.Sprintf("review-analysis-%s.pdf", month))
	attachment.SetDisposition("attachment")
	message.AddAttachment(attachment)

	// Send email
	response, err := s.client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("sendgrid API returned non-2xx status code: %d, body: %s", response.StatusCode, response.Body)
	}

	return nil
}

// SendPlainTextEmail sends a plain text email without using a template
func (s *SendGridEmailService) SendPlainTextEmail(subject string, recipient string, recipientName string, textContent string) error {
	if recipient == "" {
		return errors.New("recipient email address cannot be empty")
	}

	if textContent == "" {
		return errors.New("email content cannot be empty")
	}

	// Create email
	from := mail.NewEmail(s.fromName, s.fromEmail)
	to := mail.NewEmail(recipientName, recipient)

	// Create a simple V3 mail with plain text content
	message := mail.NewSingleEmail(from, subject, to, textContent, textContent)

	// Send email
	response, err := s.client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("sendgrid API returned non-2xx status code: %d, body: %s", response.StatusCode, response.Body)
	}

	return nil
}
