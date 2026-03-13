package email

import (
	"context"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/emersion/go-message/mail"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/config"
)

// IncomingEmail represents a parsed email fetched from the IMAP inbox.
type IncomingEmail struct {
	UID        imap.UID
	From       string
	Subject    string
	Body       string
	ReceivedAt time.Time
	// TicketNumber extracted from the subject line (e.g. "TKT-20260310-0025").
	TicketNumber string
}

// IMAPReader connects to an IMAP mailbox and fetches unread emails.
type IMAPReader struct {
	host     string
	port     int
	username string
	password string
}

// NewIMAPReader creates a new IMAP inbox reader from config.
func NewIMAPReader(cfg config.IMAPConfig) (*IMAPReader, error) {
	if cfg.Host == "" {
		return nil, fmt.Errorf("IMAP_HOST is required")
	}
	if cfg.Port == 0 {
		return nil, fmt.Errorf("IMAP_PORT is required")
	}
	if cfg.Username == "" {
		return nil, fmt.Errorf("IMAP_USERNAME is required")
	}
	if cfg.Password == "" {
		return nil, fmt.Errorf("IMAP_PASSWORD is required")
	}

	return &IMAPReader{
		host:     cfg.Host,
		port:     cfg.Port,
		username: cfg.Username,
		password: cfg.Password,
	}, nil
}

// ticketNumberRe matches ticket numbers like TKT-20260310-0025 or TKT-RWY-20260310-025.
var ticketNumberRe = regexp.MustCompile(`TKT-[A-Z0-9-]+-\d{3,4}`)

// FetchUnread connects to the IMAP server, fetches all UNSEEN messages from
// INBOX whose subject contains a ticket number, marks them as \Seen, and
// returns the parsed emails.
func (r *IMAPReader) FetchUnread(ctx context.Context) ([]IncomingEmail, error) {
	addr := fmt.Sprintf("%s:%d", r.host, r.port)

	c, err := imapclient.DialTLS(addr, nil)
	if err != nil {
		return nil, fmt.Errorf("IMAP dial %s failed: %w", addr, err)
	}
	defer c.Close()

	loginCmd := c.Login(r.username, r.password)
	if err := loginCmd.Wait(); err != nil {
		return nil, fmt.Errorf("IMAP login failed: %w", err)
	}
	log.Println("[imap-reader] login successful")

	selectCmd := c.Select("INBOX", nil)
	selectedMbox, err := selectCmd.Wait()
	if err != nil {
		return nil, fmt.Errorf("IMAP select INBOX failed: %w", err)
	}
	log.Printf("[imap-reader] INBOX selected (%d messages)", selectedMbox.NumMessages)

	// Search for unseen messages
	searchData, err := c.UIDSearch(&imap.SearchCriteria{
		NotFlag: []imap.Flag{imap.FlagSeen},
	}, nil).Wait()
	if err != nil {
		return nil, fmt.Errorf("IMAP UID SEARCH failed: %w", err)
	}

	uids := searchData.AllUIDs()
	if len(uids) == 0 {
		return nil, nil
	}

	log.Printf("[imap-reader] found %d unseen messages", len(uids))

	// Fetch envelope + full body for the unseen messages
	uidSet := imap.UIDSetNum(uids...)
	fetchOptions := &imap.FetchOptions{
		Envelope: true,
		UID:      true,
		BodySection: []*imap.FetchItemBodySection{
			{}, // empty = full RFC 822 message (headers + body)
		},
	}

	fetchCmd := c.Fetch(uidSet, fetchOptions)

	var emails []IncomingEmail
	var seenUIDs []imap.UID

	for {
		msg := fetchCmd.Next()
		if msg == nil {
			break
		}

		envelope, bodyBytes, uid := parseMessage(msg)
		if envelope == nil {
			continue
		}

		// Extract ticket number from subject
		ticketNum := ticketNumberRe.FindString(envelope.Subject)
		if ticketNum == "" {
			// Not a ticket reply — skip but still mark as seen to avoid re-processing
			log.Printf("[imap-reader] skipping non-ticket email: %q", envelope.Subject)
			seenUIDs = append(seenUIDs, uid)
			continue
		}

		// Parse sender
		from := ""
		if len(envelope.From) > 0 {
			from = envelope.From[0].Addr()
		}

		// Read body text — parse MIME to extract text/plain part
		body := ""
		if len(bodyBytes) > 0 {
			parsed := parseMIMEBody(bodyBytes)
			if parsed != "" {
				body = extractPlainText(parsed)
			}
			log.Printf("[imap-reader] UID %d ticket=%s bodyLen=%d body=%q", uid, ticketNum, len(bodyBytes), truncate(body, 200))
		} else {
			log.Printf("[imap-reader] UID %d ticket=%s WARNING: empty body bytes", uid, ticketNum)
		}

		emails = append(emails, IncomingEmail{
			UID:          uid,
			From:         from,
			Subject:      envelope.Subject,
			Body:         body,
			ReceivedAt:   envelope.Date,
			TicketNumber: ticketNum,
		})
		seenUIDs = append(seenUIDs, uid)
	}

	if err := fetchCmd.Close(); err != nil {
		return nil, fmt.Errorf("IMAP FETCH failed: %w", err)
	}

	// Mark processed messages as \Seen
	if len(seenUIDs) > 0 {
		markSet := imap.UIDSetNum(seenUIDs...)
		storeCmd := c.Store(markSet, &imap.StoreFlags{
			Op:    imap.StoreFlagsAdd,
			Flags: []imap.Flag{imap.FlagSeen},
		}, nil)
		// Drain the store command results
		for {
			if storeCmd.Next() == nil {
				break
			}
		}
		if err := storeCmd.Close(); err != nil {
			log.Printf("[imap-reader] failed to mark messages as seen: %v", err)
		}
	}

	// Logout gracefully
	if err := c.Logout().Wait(); err != nil {
		log.Printf("[imap-reader] logout error (non-fatal): %v", err)
	}

	return emails, nil
}

// parseMessage extracts the envelope, raw body bytes, and UID from a fetch message.
// IMPORTANT: The body section literal MUST be fully read before calling msg.Next()
// again, because it is a streamed reader from the IMAP connection.
func parseMessage(msg *imapclient.FetchMessageData) (*imap.Envelope, []byte, imap.UID) {
	var envelope *imap.Envelope
	var bodyBytes []byte
	var uid imap.UID

	for {
		item := msg.Next()
		if item == nil {
			break
		}
		switch data := item.(type) {
		case imapclient.FetchItemDataEnvelope:
			envelope = data.Envelope
		case imapclient.FetchItemDataBodySection:
			// Must read the literal immediately — it streams from the connection
			raw, err := io.ReadAll(data.Literal)
			if err != nil {
				log.Printf("[imap-reader] error reading body literal: %v", err)
			} else {
				bodyBytes = raw
			}
		case imapclient.FetchItemDataUID:
			uid = data.UID
		}
	}

	return envelope, bodyBytes, uid
}

// extractPlainText strips common email reply prefixes/quotes and HTML tags,
// returning only the new content written by the customer.
func extractPlainText(raw string) string {
	// Try to find the cut-off point for quoted text
	// Common patterns: "On ... wrote:", "> quoted line", "-----Original Message-----"
	lines := strings.Split(raw, "\n")
	var result []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Stop at quoted reply markers
		if strings.HasPrefix(trimmed, ">") {
			break
		}
		if strings.Contains(trimmed, "wrote:") && strings.Contains(trimmed, "On ") {
			break
		}
		// Gmail Indonesian: "Pada Sel, 10 Mar 2026 pukul 13.35 UmrahStore menulis:"
		if strings.Contains(trimmed, "menulis:") && strings.Contains(trimmed, "Pada ") {
			break
		}
		if strings.HasPrefix(trimmed, "-----") {
			break
		}
		if strings.Contains(trimmed, "Pesan ini dikirim otomatis dari sistem") {
			break
		}

		result = append(result, line)
	}

	text := strings.Join(result, "\n")

	// Strip HTML tags if present
	text = stripHTMLTags(text)

	return strings.TrimSpace(text)
}

// stripHTMLTags removes HTML tags from a string (simple regex-based).
var htmlTagRe = regexp.MustCompile(`<[^>]*>`)

func stripHTMLTags(s string) string {
	return htmlTagRe.ReplaceAllString(s, "")
}

// parseMIMEBody reads a full RFC 822 message (as raw bytes) and extracts the
// text/plain part. Falls back to text/html (stripped) if no text/plain is found.
// If MIME parsing fails entirely, returns the raw body as-is.
func parseMIMEBody(raw []byte) string {
	mr, err := mail.CreateReader(strings.NewReader(string(raw)))
	if err != nil {
		log.Printf("[imap-reader] MIME parse error: %v (raw len=%d)", err, len(raw))
		// Not a proper MIME message — return raw as-is
		return strings.TrimSpace(string(raw))
	}
	defer mr.Close()

	var plainText, htmlText string

	for {
		p, err := mr.NextPart()
		if err != nil {
			break
		}

		switch p.Header.(type) {
		case *mail.InlineHeader:
			ct, _, _ := p.Header.(*mail.InlineHeader).ContentType()
			bodyBytes, readErr := io.ReadAll(p.Body)
			if readErr != nil {
				continue
			}
			content := string(bodyBytes)

			if ct == "text/plain" && plainText == "" {
				plainText = content
			} else if ct == "text/html" && htmlText == "" {
				htmlText = content
			}
		}
	}

	// Prefer text/plain
	if plainText != "" {
		return strings.TrimSpace(plainText)
	}
	// Fall back to HTML with tags stripped
	if htmlText != "" {
		return strings.TrimSpace(stripHTMLTags(htmlText))
	}

	return ""
}

// truncate shortens a string for logging.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
