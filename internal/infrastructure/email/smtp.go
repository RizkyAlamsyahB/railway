package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"mime"
	"net"
	"net/smtp"
	"strings"

	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/config"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type smtpSender struct {
	host      string
	port      int
	username  string
	password  string
	fromEmail string
	fromName  string
}

// NewSMTPSender creates a new EmailProvider backed by SMTP with STARTTLS.
func NewSMTPSender(cfg config.SMTPConfig) (domain.EmailProvider, error) {
	if cfg.Host == "" {
		return nil, fmt.Errorf("SMTP_HOST is required")
	}
	if cfg.Port == 0 {
		return nil, fmt.Errorf("SMTP_PORT is required")
	}
	if cfg.Username == "" {
		return nil, fmt.Errorf("SMTP_USERNAME is required")
	}
	if cfg.Password == "" {
		return nil, fmt.Errorf("SMTP_PASSWORD is required")
	}
	if cfg.FromEmail == "" {
		return nil, fmt.Errorf("SMTP_FROM_EMAIL is required")
	}

	return &smtpSender{
		host:      cfg.Host,
		port:      cfg.Port,
		username:  cfg.Username,
		password:  cfg.Password,
		fromEmail: cfg.FromEmail,
		fromName:  cfg.FromName,
	}, nil
}

// Send delivers an email message via SMTP with STARTTLS authentication.
func (s *smtpSender) Send(ctx context.Context, msg domain.EmailMessage) error {
	if len(msg.To) == 0 {
		return fmt.Errorf("at least one recipient (To) is required")
	}

	addr := net.JoinHostPort(s.host, fmt.Sprintf("%d", s.port))

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server %s: %w", addr, err)
	}

	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	tlsConfig := &tls.Config{ServerName: s.host}
	if err := client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("STARTTLS failed: %w", err)
	}

	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP authentication failed: %w", err)
	}

	if err := client.Mail(s.fromEmail); err != nil {
		return fmt.Errorf("MAIL FROM failed: %w", err)
	}

	allRecipients := s.collectRecipients(msg)
	for _, rcpt := range allRecipients {
		if err := client.Rcpt(rcpt); err != nil {
			return fmt.Errorf("RCPT TO <%s> failed: %w", rcpt, err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("DATA command failed: %w", err)
	}

	body := s.buildMessage(msg)
	if _, err := w.Write([]byte(body)); err != nil {
		return fmt.Errorf("failed to write email body: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	return client.Quit()
}

// collectRecipients merges To, CC, and BCC into a single deduplicated slice.
func (s *smtpSender) collectRecipients(msg domain.EmailMessage) []string {
	seen := make(map[string]struct{})
	var result []string

	for _, list := range [][]string{msg.To, msg.CC, msg.BCC} {
		for _, addr := range list {
			lower := strings.ToLower(strings.TrimSpace(addr))
			if _, exists := seen[lower]; !exists {
				seen[lower] = struct{}{}
				result = append(result, addr)
			}
		}
	}

	return result
}

// buildMessage constructs the full RFC 2822 email message with MIME headers.
func (s *smtpSender) buildMessage(msg domain.EmailMessage) string {
	var b strings.Builder

	// From header
	if s.fromName != "" {
		b.WriteString(fmt.Sprintf("From: %s <%s>\r\n",
			mime.QEncoding.Encode("utf-8", s.fromName), s.fromEmail))
	} else {
		b.WriteString(fmt.Sprintf("From: %s\r\n", s.fromEmail))
	}

	// To header
	b.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(msg.To, ", ")))

	// CC header (BCC is intentionally omitted from headers)
	if len(msg.CC) > 0 {
		b.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(msg.CC, ", ")))
	}

	// Subject with UTF-8 encoding
	b.WriteString(fmt.Sprintf("Subject: %s\r\n",
		mime.QEncoding.Encode("utf-8", msg.Subject)))

	// MIME version
	b.WriteString("MIME-Version: 1.0\r\n")

	// Content-Type
	if msg.IsHTML {
		b.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	} else {
		b.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	}

	// Blank line separating headers from body
	b.WriteString("\r\n")
	b.WriteString(msg.Body)

	return b.String()
}
