package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

type ticketUseCase struct {
	ticketRepo  domain.TicketRepository
	subjectRepo domain.TicketSubjectRepository
	userRepo    domain.UserRepository
	orderRepo   domain.OrderRepository
	paymentRepo domain.PaymentRepository
	email       domain.EmailProvider
	storage     domain.StorageProvider
}

func NewTicketUseCase(ticketRepo domain.TicketRepository, subjectRepo domain.TicketSubjectRepository, userRepo domain.UserRepository, orderRepo domain.OrderRepository, paymentRepo domain.PaymentRepository, email domain.EmailProvider, storage domain.StorageProvider) domain.TicketUseCase {
	return &ticketUseCase{ticketRepo: ticketRepo, subjectRepo: subjectRepo, userRepo: userRepo, orderRepo: orderRepo, paymentRepo: paymentRepo, email: email, storage: storage}
}

// generateTicketNumber creates a ticket number in format TKT-YYYYMMDD-NNNN.
func (uc *ticketUseCase) generateTicketNumber(ctx context.Context) (string, error) {
	today := time.Now().Format("20060102")
	count, err := uc.ticketRepo.CountOnDate(ctx, time.Now().Format("2006-01-02"))
	if err != nil {
		return "", fmt.Errorf("failed to count tickets: %w", err)
	}
	return fmt.Sprintf("TKT-%s-%04d", today, count+1), nil
}

func (uc *ticketUseCase) PresignTicketAttachment(ctx context.Context, req domain.PresignTicketAttachmentRequest) (*domain.PresignTicketAttachmentResponse, error) {
	ct := normalizeContentType(req.ContentType)
	if !isAllowedTicketAttachmentContentType(ct) {
		return nil, ErrInvalidAttachmentContentType
	}

	objectKey := fmt.Sprintf("ticket-attachments/%s/%s", time.Now().Format("2006/01/02"), uuid.New().String())
	uploadURL, err := uc.storage.GeneratePresignedUploadURL(ctx, objectKey, ct, PresignedUploadExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned upload URL: %w", err)
	}

	return &domain.PresignTicketAttachmentResponse{
		UploadURL:   uploadURL,
		ObjectKey:   objectKey,
		ContentType: ct,
		ExpiresIn:   int(PresignedUploadExpiry.Seconds()),
	}, nil
}

func (uc *ticketUseCase) CreateTicket(ctx context.Context, customerID uuid.UUID, req domain.CreateTicketRequest) (*domain.TicketResponse, error) {
	// Resolve subject from subject_id
	subj, err := uc.subjectRepo.FindByID(ctx, req.SubjectID)
	if err != nil {
		return nil, ErrTicketSubjectNotFound
	}

	// Verify attachment was actually uploaded to S3
	info, err := uc.storage.HeadObject(ctx, req.AttachmentKey)
	if err != nil || info == nil {
		return nil, ErrAttachmentNotUploaded
	}

	// Validate size (<= 5 MB)
	if info.ContentLength > TicketAttachmentMaxBytes {
		return nil, ErrAttachmentTooLarge
	}

	// Validate content type
	ct := normalizeContentType(info.ContentType)
	if !isAllowedTicketAttachmentContentType(ct) {
		return nil, ErrInvalidAttachmentContentType
	}

	attachmentURL := uc.storage.GetURL(req.AttachmentKey)

	ticketNumber, err := uc.generateTicketNumber(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	t := &domain.Ticket{
		ID:                    uuid.New(),
		TicketNumber:          ticketNumber,
		CustomerID:            customerID,
		OrderNumber:           req.OrderNumber,
		Phone:                 req.Phone,
		ReporterName:          req.ReporterName,
		SubjectID:             subj.ID,
		Subject:               subj.Label,
		Detail:                req.Detail,
		Status:                domain.TicketStatusOpen,
		Source:                req.Source,
		AttachmentURL:         &attachmentURL,
		AttachmentContentType: &ct,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	if err := uc.ticketRepo.Create(ctx, t); err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	return toTicketResponse(*t), nil
}

func (uc *ticketUseCase) ListTickets(ctx context.Context, params domain.TicketListParams) ([]domain.TicketResponse, *domain.PaginationMeta, error) {
	tickets, meta, err := uc.ticketRepo.List(ctx, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list tickets: %w", err)
	}
	resp := make([]domain.TicketResponse, len(tickets))
	for i, t := range tickets {
		uc.enrichTicket(ctx, &t)
		resp[i] = *toTicketResponse(t)
	}
	return resp, meta, nil
}

func (uc *ticketUseCase) GetTicket(ctx context.Context, id uuid.UUID) (*domain.TicketResponse, error) {
	t, err := uc.ticketRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTicketNotFound
		}
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}
	uc.enrichTicket(ctx, t)
	return toTicketResponse(*t), nil
}

// TakeTicket lets a CS self-assign an open/unassigned ticket (first-come-first-served).
func (uc *ticketUseCase) TakeTicket(ctx context.Context, csID, ticketID uuid.UUID) (*domain.TicketResponse, error) {
	t, err := uc.ticketRepo.FindByID(ctx, ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTicketNotFound
		}
		return nil, fmt.Errorf("failed to find ticket: %w", err)
	}

	// Already closed → read-only
	if t.Status == domain.TicketStatusClosed {
		return nil, ErrTicketClosed
	}

	// Already assigned to someone else → reject
	if t.AssignedCSID != nil && *t.AssignedCSID != csID {
		return nil, ErrTicketAlreadyTaken
	}

	now := time.Now()
	t.AssignedCSID = &csID
	t.Status = domain.TicketStatusOnProgress
	t.UpdatedAt = now

	if err := uc.ticketRepo.Update(ctx, t); err != nil {
		return nil, fmt.Errorf("failed to take ticket: %w", err)
	}

	// Log status change
	oldStatus := domain.TicketStatusOpen
	statusLog := &domain.TicketStatusLog{
		ID:        uuid.New(),
		TicketID:  ticketID,
		ChangedBy: csID,
		OldStatus: &oldStatus,
		NewStatus: domain.TicketStatusOnProgress,
		CreatedAt: now,
	}
	_ = uc.ticketRepo.CreateStatusLog(ctx, statusLog)

	uc.enrichTicket(ctx, t)
	return toTicketResponse(*t), nil
}

func (uc *ticketUseCase) UpdateTicketStatus(ctx context.Context, csID, ticketID uuid.UUID, req domain.UpdateTicketStatusRequest) (*domain.TicketResponse, error) {
	t, err := uc.ticketRepo.FindByID(ctx, ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTicketNotFound
		}
		return nil, fmt.Errorf("failed to find ticket: %w", err)
	}

	// Only the assigned CS can change status
	if t.AssignedCSID == nil || *t.AssignedCSID != csID {
		return nil, ErrTicketNotAssignedToYou
	}

	// Already closed → no further status changes
	if t.Status == domain.TicketStatusClosed {
		return nil, ErrTicketClosed
	}

	oldStatus := t.Status
	now := time.Now()
	t.Status = req.Status
	t.UpdatedAt = now

	if req.Status == domain.TicketStatusClosed {
		t.ClosedAt = &now
	}

	if err := uc.ticketRepo.Update(ctx, t); err != nil {
		return nil, fmt.Errorf("failed to update ticket: %w", err)
	}

	statusLog := &domain.TicketStatusLog{
		ID:        uuid.New(),
		TicketID:  ticketID,
		ChangedBy: csID,
		OldStatus: &oldStatus,
		NewStatus: req.Status,
		Notes:     req.Notes,
		CreatedAt: now,
	}
	_ = uc.ticketRepo.CreateStatusLog(ctx, statusLog)

	uc.enrichTicket(ctx, t)
	return toTicketResponse(*t), nil
}

func (uc *ticketUseCase) AddTicketMessage(ctx context.Context, senderID uuid.UUID, isCS bool, ticketID uuid.UUID, req domain.AddTicketMessageRequest) (*domain.TicketMessageResponse, error) {
	t, err := uc.ticketRepo.FindByID(ctx, ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTicketNotFound
		}
		return nil, fmt.Errorf("failed to find ticket: %w", err)
	}

	// Only the assigned CS can send messages on this ticket
	if isCS {
		if t.AssignedCSID == nil || *t.AssignedCSID != senderID {
			return nil, ErrTicketNotAssignedToYou
		}
	}

	// Closed tickets are read-only
	if t.Status == domain.TicketStatusClosed {
		return nil, ErrTicketClosed
	}

	msg := &domain.TicketMessage{
		ID:             uuid.New(),
		TicketID:       ticketID,
		SenderID:       senderID,
		Message:        req.Message,
		IsFromCS:       isCS,
		IsInternalNote: req.IsInternalNote,
		CreatedAt:      time.Now(),
	}

	if err := uc.ticketRepo.CreateMessage(ctx, msg); err != nil {
		return nil, fmt.Errorf("failed to add message: %w", err)
	}

	// Send email to customer (non-blocking, best-effort)
	if isCS && !req.IsInternalNote {
		go uc.sendTicketReplyEmail(t, req.Message)
	}

	return toTicketMessageResponse(*msg), nil
}

func (uc *ticketUseCase) ListTicketMessages(ctx context.Context, ticketID uuid.UUID) ([]domain.TicketMessageResponse, error) {
	_, err := uc.ticketRepo.FindByID(ctx, ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTicketNotFound
		}
		return nil, fmt.Errorf("failed to find ticket: %w", err)
	}

	msgs, err := uc.ticketRepo.ListMessages(ctx, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}

	resp := make([]domain.TicketMessageResponse, len(msgs))
	for i, m := range msgs {
		resp[i] = *toTicketMessageResponse(m)
	}
	return resp, nil
}

// ============================================================
// Private helpers
// ============================================================

// enrichTicket populates joined fields (CustomerEmail, AssignedCSName) by
// looking up the related users. Errors are swallowed — the fields will simply
// remain zero-valued if the lookup fails.
func (uc *ticketUseCase) enrichTicket(ctx context.Context, t *domain.Ticket) {
	if customer, err := uc.userRepo.FindByID(ctx, t.CustomerID); err == nil {
		t.CustomerEmail = customer.Email
	}
	if t.AssignedCSID != nil {
		if cs, err := uc.userRepo.FindByID(ctx, *t.AssignedCSID); err == nil {
			t.AssignedCSName = &cs.FullName
		}
	}
	// Enrich order info (status, payment method, grand total)
	if t.OrderNumber != "" {
		if order, err := uc.orderRepo.FindByOrderNo(ctx, t.OrderNumber); err == nil && order != nil {
			t.OrderStatus = order.OrderStatus
			t.GrandTotal = order.GrandTotal
			if invoice, err := uc.paymentRepo.FindInvoiceByOrderID(ctx, order.ID); err == nil && invoice != nil {
				t.PaymentMethod = invoice.PaymentMethod
			}
		}
	}
}

// sendTicketReplyEmail sends the CS reply to the customer's email via SMTP.
// It runs in a background goroutine and logs errors rather than propagating them.
func (uc *ticketUseCase) sendTicketReplyEmail(t *domain.Ticket, message string) {
	// Look up customer email
	customer, err := uc.userRepo.FindByID(context.Background(), t.CustomerID)
	if err != nil {
		log.Printf("[ticket-email] failed to find customer %s: %v", t.CustomerID, err)
		return
	}

	// Replace template placeholders with actual ticket/customer data
	message = resolveTemplatePlaceholders(message, t, customer)

	subject := fmt.Sprintf("Balasan Tiket %s — %s", t.TicketNumber, t.Subject)
	body := fmt.Sprintf(`<html><body>
<p>Assalamu'alaikum <strong>%s</strong>,</p>
<p>Berikut balasan dari tim Customer Service kami untuk tiket <strong>%s</strong>:</p>
<hr/>
<div style="padding:12px;background:#f9f6f0;border-left:4px solid #d4a853;margin:16px 0;white-space:pre-wrap;">%s</div>
<hr/>
<p><em>Pesan ini dikirim otomatis dari sistem Haji &amp; Umrah Store. Anda dapat membalas langsung ke email ini untuk melanjutkan percakapan.</em></p>
<p>Jazakumullahu khairan,<br/>Tim Customer Service<br/>Haji &amp; Umrah Store</p>
</body></html>`, customer.FullName, t.TicketNumber, message)

	emailMsg := domain.EmailMessage{
		To:      []string{customer.Email},
		Subject: subject,
		Body:    body,
		IsHTML:  true,
	}

	if err := uc.email.Send(context.Background(), emailMsg); err != nil {
		log.Printf("[ticket-email] failed to send email for ticket %s to %s: %v", t.TicketNumber, customer.Email, err)
	} else {
		log.Printf("[ticket-email] email sent for ticket %s to %s", t.TicketNumber, customer.Email)
	}
}

// resolveTemplatePlaceholders replaces {variable} placeholders in a reply
// template message with actual data from the ticket and customer.
// Unknown placeholders (e.g. {tracking_number} when no shipping data exists)
// are left as-is so CS can manually fill them before sending, or they serve
// as a visual cue that the data wasn't available.
func resolveTemplatePlaceholders(message string, t *domain.Ticket, customer *domain.User) string {
	replacements := map[string]string{
		// Ticket data
		"{ticket_number}": t.TicketNumber,
		"{order_number}":  t.OrderNumber,
		"{subject}":       t.Subject,
		"{detail}":        t.Detail,
		"{status}":        t.Status,
		"{source}":        t.Source,
		"{phone}":         t.Phone,

		// Customer data
		"{customer_name}":  customer.FullName,
		"{customer_email}": customer.Email,
		"{reporter_name}":  t.ReporterName,
	}

	// Add customer phone if available
	if customer.Phone != nil {
		replacements["{customer_phone}"] = *customer.Phone
	}

	for placeholder, value := range replacements {
		message = strings.ReplaceAll(message, placeholder, value)
	}
	return message
}
