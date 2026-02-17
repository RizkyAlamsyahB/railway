package email

import (
	"context"
	"strings"
	"testing"

	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/config"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

func validConfig() config.SMTPConfig {
	return config.SMTPConfig{
		Host:      "smtp.gmail.com",
		Port:      587,
		Username:  "test@gmail.com",
		Password:  "test-app-password",
		FromEmail: "noreply@example.com",
		FromName:  "Test Sender",
	}
}

func newTestSender(t *testing.T) *smtpSender {
	t.Helper()
	provider, err := NewSMTPSender(validConfig())
	if err != nil {
		t.Fatalf("NewSMTPSender() unexpected error: %v", err)
	}
	return provider.(*smtpSender)
}

func TestNewSMTPSender_Success(t *testing.T) {
	provider, err := NewSMTPSender(validConfig())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if provider == nil {
		t.Fatal("expected non-nil provider")
	}
}

func TestNewSMTPSender_Validation(t *testing.T) {
	tests := []struct {
		name    string
		modify  func(*config.SMTPConfig)
		wantErr string
	}{
		{
			name:    "missing host",
			modify:  func(c *config.SMTPConfig) { c.Host = "" },
			wantErr: "SMTP_HOST is required",
		},
		{
			name:    "missing port",
			modify:  func(c *config.SMTPConfig) { c.Port = 0 },
			wantErr: "SMTP_PORT is required",
		},
		{
			name:    "missing username",
			modify:  func(c *config.SMTPConfig) { c.Username = "" },
			wantErr: "SMTP_USERNAME is required",
		},
		{
			name:    "missing password",
			modify:  func(c *config.SMTPConfig) { c.Password = "" },
			wantErr: "SMTP_PASSWORD is required",
		},
		{
			name:    "missing from email",
			modify:  func(c *config.SMTPConfig) { c.FromEmail = "" },
			wantErr: "SMTP_FROM_EMAIL is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := validConfig()
			tt.modify(&cfg)

			_, err := NewSMTPSender(cfg)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("expected error containing %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

func TestBuildMessage_PlainText(t *testing.T) {
	s := newTestSender(t)
	msg := domain.EmailMessage{
		To:      []string{"recipient@example.com"},
		Subject: "Test Subject",
		Body:    "Hello, this is a test.",
		IsHTML:  false,
	}

	result := s.buildMessage(msg)

	checks := []struct {
		name     string
		contains string
	}{
		{"from with name", "From: Test Sender <noreply@example.com>"},
		{"to", "To: recipient@example.com"},
		{"subject", "Subject:"},
		{"mime version", "MIME-Version: 1.0"},
		{"content type", "Content-Type: text/plain; charset=UTF-8"},
		{"body", "Hello, this is a test."},
	}

	for _, c := range checks {
		if !strings.Contains(result, c.contains) {
			t.Errorf("%s: expected message to contain %q, got:\n%s", c.name, c.contains, result)
		}
	}

	// Verify body separator
	if !strings.Contains(result, "\r\n\r\n") {
		t.Error("expected CRLF CRLF separator between headers and body")
	}
}

func TestBuildMessage_HTML(t *testing.T) {
	s := newTestSender(t)
	msg := domain.EmailMessage{
		To:      []string{"recipient@example.com"},
		Subject: "HTML Test",
		Body:    "<h1>Hello</h1>",
		IsHTML:  true,
	}

	result := s.buildMessage(msg)

	if !strings.Contains(result, "Content-Type: text/html; charset=UTF-8") {
		t.Errorf("expected HTML content type, got:\n%s", result)
	}
}

func TestBuildMessage_WithCC(t *testing.T) {
	s := newTestSender(t)
	msg := domain.EmailMessage{
		To:      []string{"to@example.com"},
		CC:      []string{"cc1@example.com", "cc2@example.com"},
		Subject: "CC Test",
		Body:    "test",
	}

	result := s.buildMessage(msg)

	if !strings.Contains(result, "Cc: cc1@example.com, cc2@example.com") {
		t.Errorf("expected Cc header, got:\n%s", result)
	}
}

func TestBuildMessage_BCC_NotInHeaders(t *testing.T) {
	s := newTestSender(t)
	msg := domain.EmailMessage{
		To:      []string{"to@example.com"},
		BCC:     []string{"secret@example.com"},
		Subject: "BCC Test",
		Body:    "test",
	}

	result := s.buildMessage(msg)

	if strings.Contains(result, "secret@example.com") {
		t.Error("BCC address should not appear in message headers")
	}
	if strings.Contains(result, "Bcc") {
		t.Error("Bcc header should not be present")
	}
}

func TestBuildMessage_NoFromName(t *testing.T) {
	cfg := validConfig()
	cfg.FromName = ""
	provider, err := NewSMTPSender(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := provider.(*smtpSender)

	msg := domain.EmailMessage{
		To:      []string{"to@example.com"},
		Subject: "Test",
		Body:    "test",
	}

	result := s.buildMessage(msg)

	if !strings.Contains(result, "From: noreply@example.com\r\n") {
		t.Errorf("expected From without display name, got:\n%s", result)
	}
}

func TestCollectRecipients_AllFields(t *testing.T) {
	s := newTestSender(t)
	msg := domain.EmailMessage{
		To:  []string{"to@example.com"},
		CC:  []string{"cc@example.com"},
		BCC: []string{"bcc@example.com"},
	}

	result := s.collectRecipients(msg)

	if len(result) != 3 {
		t.Fatalf("expected 3 recipients, got %d: %v", len(result), result)
	}

	expected := map[string]bool{
		"to@example.com":  false,
		"cc@example.com":  false,
		"bcc@example.com": false,
	}
	for _, r := range result {
		expected[r] = true
	}
	for addr, found := range expected {
		if !found {
			t.Errorf("expected recipient %q not found", addr)
		}
	}
}

func TestCollectRecipients_Deduplication(t *testing.T) {
	s := newTestSender(t)
	msg := domain.EmailMessage{
		To:  []string{"same@example.com"},
		CC:  []string{"Same@Example.com"},
		BCC: []string{"SAME@EXAMPLE.COM"},
	}

	result := s.collectRecipients(msg)

	if len(result) != 1 {
		t.Errorf("expected 1 deduplicated recipient, got %d: %v", len(result), result)
	}
}

func TestSend_EmptyTo(t *testing.T) {
	s := newTestSender(t)
	msg := domain.EmailMessage{
		Subject: "No Recipients",
		Body:    "test",
	}

	err := s.Send(context.Background(), msg)
	if err == nil {
		t.Fatal("expected error for empty To, got nil")
	}
	if !strings.Contains(err.Error(), "at least one recipient") {
		t.Errorf("expected recipient error, got: %v", err)
	}
}
