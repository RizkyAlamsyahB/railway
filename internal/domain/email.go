package domain

import "context"

// EmailMessage represents a single email to be sent.
type EmailMessage struct {
	// To is the list of primary recipient email addresses.
	To []string

	// CC is the list of carbon-copy recipient email addresses.
	CC []string

	// BCC is the list of blind-carbon-copy recipient email addresses.
	BCC []string

	// Subject is the email subject line.
	Subject string

	// Body is the email content (plain text or HTML depending on IsHTML).
	Body string

	// IsHTML indicates whether Body contains HTML content.
	// When true, the Content-Type header is set to "text/html; charset=UTF-8",
	// otherwise "text/plain; charset=UTF-8".
	IsHTML bool
}

// EmailProvider defines the interface for sending emails.
// Implementations may target SMTP, SendGrid, SES, or any other email delivery service.
type EmailProvider interface {
	// Send delivers an email message. Returns an error if delivery fails.
	Send(ctx context.Context, msg EmailMessage) error
}
