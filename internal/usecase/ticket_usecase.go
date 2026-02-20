package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

type ticketUseCase struct {
	ticketRepo domain.TicketRepository
}

func NewTicketUseCase(ticketRepo domain.TicketRepository) domain.TicketUseCase {
	return &ticketUseCase{ticketRepo: ticketRepo}
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

func (uc *ticketUseCase) CreateTicket(ctx context.Context, customerID uuid.UUID, req domain.CreateTicketRequest) (*domain.TicketResponse, error) {
	ticketNumber, err := uc.generateTicketNumber(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	t := &domain.Ticket{
		ID:           uuid.New(),
		TicketNumber: ticketNumber,
		CustomerID:   customerID,
		OrderNumber:  req.OrderNumber,
		Phone:        req.Phone,
		ReporterName: req.ReporterName,
		Subject:      req.Subject,
		Detail:       req.Detail,
		Status:       domain.TicketStatusOpen,
		Source:       req.Source,
		CreatedAt:    now,
		UpdatedAt:    now,
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

	oldStatus := t.Status
	now := time.Now()
	t.Status = req.Status
	t.UpdatedAt = now

	switch req.Status {
	case domain.TicketStatusResolved:
		t.ResolvedAt = &now
	case domain.TicketStatusClosed:
		t.ClosedAt = &now
	}

	if err := uc.ticketRepo.Update(ctx, t); err != nil {
		return nil, fmt.Errorf("failed to update ticket: %w", err)
	}

	log := &domain.TicketStatusLog{
		ID:        uuid.New(),
		TicketID:  ticketID,
		ChangedBy: csID,
		OldStatus: &oldStatus,
		NewStatus: req.Status,
		Notes:     req.Notes,
		CreatedAt: now,
	}
	_ = uc.ticketRepo.CreateStatusLog(ctx, log)

	return toTicketResponse(*t), nil
}

func (uc *ticketUseCase) AssignTicket(ctx context.Context, csID, ticketID uuid.UUID, req domain.AssignTicketRequest) (*domain.TicketResponse, error) {
	t, err := uc.ticketRepo.FindByID(ctx, ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTicketNotFound
		}
		return nil, fmt.Errorf("failed to find ticket: %w", err)
	}

	t.AssignedCSID = &req.AssignedCSID
	t.UpdatedAt = time.Now()

	if err := uc.ticketRepo.Update(ctx, t); err != nil {
		return nil, fmt.Errorf("failed to assign ticket: %w", err)
	}

	return toTicketResponse(*t), nil
}

func (uc *ticketUseCase) AddTicketMessage(ctx context.Context, senderID uuid.UUID, isCS bool, ticketID uuid.UUID, req domain.AddTicketMessageRequest) (*domain.TicketMessageResponse, error) {
	_, err := uc.ticketRepo.FindByID(ctx, ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTicketNotFound
		}
		return nil, fmt.Errorf("failed to find ticket: %w", err)
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
