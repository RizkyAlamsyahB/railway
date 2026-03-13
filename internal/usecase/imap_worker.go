package usecase

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/infrastructure/email"
	"gorm.io/gorm"
)

// IMAPWorker periodically checks an IMAP inbox for customer replies to
// ticket emails and creates corresponding ticket messages automatically.
type IMAPWorker struct {
	reader     *email.IMAPReader
	ticketRepo domain.TicketRepository
	userRepo   domain.UserRepository
	interval   time.Duration
	stopCh     chan struct{}
}

// NewIMAPWorker creates a new background IMAP inbox poller.
func NewIMAPWorker(reader *email.IMAPReader, ticketRepo domain.TicketRepository, userRepo domain.UserRepository, pollInterval time.Duration) *IMAPWorker {
	return &IMAPWorker{
		reader:     reader,
		ticketRepo: ticketRepo,
		userRepo:   userRepo,
		interval:   pollInterval,
		stopCh:     make(chan struct{}),
	}
}

// Start begins polling the IMAP inbox in a background goroutine.
// It blocks until Stop() is called or ctx is cancelled.
func (w *IMAPWorker) Start(ctx context.Context) {
	log.Printf("[imap-worker] starting inbox poller (interval=%s)", w.interval)

	// Initial poll on startup
	w.poll(ctx)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.poll(ctx)
		case <-w.stopCh:
			log.Println("[imap-worker] stopped")
			return
		case <-ctx.Done():
			log.Println("[imap-worker] context cancelled, stopping")
			return
		}
	}
}

// Stop signals the worker to stop polling.
func (w *IMAPWorker) Stop() {
	close(w.stopCh)
}

// poll fetches unread emails and processes them.
func (w *IMAPWorker) poll(ctx context.Context) {
	emails, err := w.reader.FetchUnread(ctx)
	if err != nil {
		log.Printf("[imap-worker] fetch error: %v", err)
		return
	}

	if len(emails) == 0 {
		return
	}

	log.Printf("[imap-worker] processing %d ticket reply emails", len(emails))

	for _, incoming := range emails {
		if err := w.processEmail(ctx, incoming); err != nil {
			log.Printf("[imap-worker] failed to process email (ticket=%s, from=%s): %v",
				incoming.TicketNumber, incoming.From, err)
		}
	}
}

// processEmail matches an incoming email to a ticket and creates a message.
func (w *IMAPWorker) processEmail(ctx context.Context, incoming email.IncomingEmail) error {
	// Find the ticket by number
	ticket, err := w.ticketRepo.FindByTicketNumber(ctx, incoming.TicketNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[imap-worker] no ticket found for %s, skipping", incoming.TicketNumber)
			return nil
		}
		return err
	}

	// Don't add messages to closed tickets
	if ticket.Status == domain.TicketStatusClosed {
		log.Printf("[imap-worker] ticket %s is closed, skipping email from %s", incoming.TicketNumber, incoming.From)
		return nil
	}

	// Find the sender (customer) by email
	user, err := w.userRepo.FindByEmail(ctx, incoming.From)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[imap-worker] unknown sender %s for ticket %s, skipping", incoming.From, incoming.TicketNumber)
			return nil
		}
		return err
	}

	// Only accept replies from the ticket's customer
	if user.ID != ticket.CustomerID {
		log.Printf("[imap-worker] sender %s is not the ticket owner for %s, skipping", incoming.From, incoming.TicketNumber)
		return nil
	}

	body := incoming.Body
	if body == "" {
		body = "(email kosong)"
	}

	msg := &domain.TicketMessage{
		ID:             uuid.New(),
		TicketID:       ticket.ID,
		SenderID:       user.ID,
		Message:        body,
		IsFromCS:       false,
		IsInternalNote: false,
		CreatedAt:      time.Now(),
	}

	if err := w.ticketRepo.CreateMessage(ctx, msg); err != nil {
		return err
	}

	log.Printf("[imap-worker] created message for ticket %s from %s (email reply)",
		incoming.TicketNumber, incoming.From)
	return nil
}
