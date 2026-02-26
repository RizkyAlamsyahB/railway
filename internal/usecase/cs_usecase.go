package usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

// ============================================================
// Reply Template UseCase
// ============================================================

var reSubKeyNormalizer = regexp.MustCompile(`[^a-z0-9]+`)

func normalizeSubKey(raw string) string {
	s := strings.ToLower(raw)
	s = reSubKeyNormalizer.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

func buildShortcut(category, subKey string) string {
	return "/" + category + "/" + normalizeSubKey(subKey)
}

type replyTemplateUseCase struct {
	templateRepo domain.ReplyTemplateRepository
}

func NewReplyTemplateUseCase(templateRepo domain.ReplyTemplateRepository) domain.ReplyTemplateUseCase {
	return &replyTemplateUseCase{templateRepo: templateRepo}
}

func (uc *replyTemplateUseCase) Create(ctx context.Context, csID uuid.UUID, req domain.CreateReplyTemplateRequest) (*domain.ReplyTemplateResponse, error) {
	if _, valid := domain.ValidReplyTemplateCategories[req.Category]; !valid {
		return nil, ErrInvalidReplyTemplateCategory
	}

	shortcut := buildShortcut(req.Category, req.SubKey)

	exists, err := uc.templateRepo.ExistsShortcut(ctx, shortcut, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check shortcut: %w", err)
	}
	if exists {
		return nil, ErrShortcutAlreadyExists
	}

	now := time.Now()
	tpl := &domain.ReplyTemplate{
		ID:        uuid.New(),
		Title:     req.Title,
		Category:  req.Category,
		Shortcut:  shortcut,
		Content:   req.Content,
		IsActive:  true,
		CreatedBy: csID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := uc.templateRepo.Create(ctx, tpl); err != nil {
		return nil, fmt.Errorf("failed to create template: %w", err)
	}
	return toReplyTemplateResponse(*tpl), nil
}

func (uc *replyTemplateUseCase) List(ctx context.Context, params domain.ReplyTemplateListParams) ([]domain.ReplyTemplateResponse, *domain.PaginationMeta, error) {
	tpls, meta, err := uc.templateRepo.List(ctx, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list templates: %w", err)
	}
	resp := make([]domain.ReplyTemplateResponse, len(tpls))
	for i, t := range tpls {
		resp[i] = *toReplyTemplateResponse(t)
	}
	return resp, meta, nil
}

func (uc *replyTemplateUseCase) GetByID(ctx context.Context, id uuid.UUID) (*domain.ReplyTemplateResponse, error) {
	tpl, err := uc.templateRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReplyTemplateNotFound
		}
		return nil, fmt.Errorf("failed to get template: %w", err)
	}
	return toReplyTemplateResponse(*tpl), nil
}

func (uc *replyTemplateUseCase) Update(ctx context.Context, csID, id uuid.UUID, req domain.UpdateReplyTemplateRequest) (*domain.ReplyTemplateResponse, error) {
	tpl, err := uc.templateRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReplyTemplateNotFound
		}
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	if req.Title != nil {
		tpl.Title = *req.Title
	}
	if req.Content != nil {
		tpl.Content = *req.Content
	}
	if req.IsActive != nil {
		tpl.IsActive = *req.IsActive
	}

	// Recalculate shortcut if SubKey changed (category is immutable)
	if req.SubKey != nil {
		newShortcut := buildShortcut(tpl.Category, *req.SubKey)
		if newShortcut != tpl.Shortcut {
			excludeID := tpl.ID
			exists, err := uc.templateRepo.ExistsShortcut(ctx, newShortcut, &excludeID)
			if err != nil {
				return nil, fmt.Errorf("failed to check shortcut: %w", err)
			}
			if exists {
				return nil, ErrShortcutAlreadyExists
			}
			tpl.Shortcut = newShortcut
		}
	}

	tpl.UpdatedBy = &csID
	tpl.UpdatedAt = time.Now()

	if err := uc.templateRepo.Update(ctx, tpl); err != nil {
		return nil, fmt.Errorf("failed to update template: %w", err)
	}
	return toReplyTemplateResponse(*tpl), nil
}

func (uc *replyTemplateUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := uc.templateRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrReplyTemplateNotFound
		}
		return fmt.Errorf("failed to get template: %w", err)
	}
	return uc.templateRepo.Delete(ctx, id)
}

// ============================================================
// CS Dashboard UseCase
// ============================================================

type csDashboardUseCase struct {
	ticketRepo domain.TicketRepository
	chatRepo   domain.ChatRepository
}

func NewCSDashboardUseCase(ticketRepo domain.TicketRepository, chatRepo domain.ChatRepository) domain.CSDashboardUseCase {
	return &csDashboardUseCase{ticketRepo: ticketRepo, chatRepo: chatRepo}
}

func (uc *csDashboardUseCase) GetDashboard(ctx context.Context, params domain.CSDashboardParams) (*domain.CSDashboardResponse, error) {
	todayTickets, err := uc.ticketRepo.CountTodayTickets(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count today tickets: %w", err)
	}

	waitingTickets, err := uc.ticketRepo.CountWaitingTickets(ctx, params.Year, params.Month)
	if err != nil {
		return nil, fmt.Errorf("failed to count waiting tickets: %w", err)
	}

	doneTickets, err := uc.ticketRepo.CountDoneTickets(ctx, params.Year, params.Month)
	if err != nil {
		return nil, fmt.Errorf("failed to count done tickets: %w", err)
	}

	recentTickets, err := uc.ticketRepo.ListRecentUnassigned(ctx, 3)
	if err != nil {
		return nil, fmt.Errorf("failed to list recent unassigned tickets: %w", err)
	}

	recent := make([]domain.DashboardRecentTicket, len(recentTickets))
	for i, t := range recentTickets {
		recent[i] = domain.DashboardRecentTicket{
			ID:           t.ID,
			TicketNumber: t.TicketNumber,
			Phone:        t.Phone,
			ReporterName: t.ReporterName,
			Subject:      t.Subject,
			Status:       t.Status,
			CreatedAt:    t.CreatedAt,
		}
	}

	return &domain.CSDashboardResponse{
		TodayTickets:   todayTickets,
		WaitingTickets: waitingTickets,
		DoneTickets:    doneTickets,
		RecentTickets:  recent,
	}, nil
}

// ============================================================
// CS User UseCase  (data pengguna)
// ============================================================

type csUserUseCase struct {
	userRepo   domain.UserRepository
	ticketRepo domain.TicketRepository
}

func NewCSUserUseCase(userRepo domain.UserRepository, ticketRepo domain.TicketRepository) domain.CSUserUseCase {
	return &csUserUseCase{userRepo: userRepo, ticketRepo: ticketRepo}
}

func (uc *csUserUseCase) ListUsers(ctx context.Context, params domain.CSUserListParams) ([]domain.UserResponse, *domain.PaginationMeta, error) {
	users, total, err := uc.userRepo.List(ctx, domain.UserListParams{
		Page:   params.Page,
		Limit:  params.Limit,
		Role:   "customer",
		Search: params.Query,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list users: %w", err)
	}

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 || params.Limit > 100 {
		params.Limit = 10
	}

	resp := make([]domain.UserResponse, len(users))
	for i := range users {
		resp[i] = *toUserResponse(&users[i])
	}
	return resp, &domain.PaginationMeta{
		Page:       params.Page,
		Limit:      params.Limit,
		TotalItems: total,
		TotalPages: int(math.Ceil(float64(total) / float64(params.Limit))),
	}, nil
}

func (uc *csUserUseCase) GetUser(ctx context.Context, userID uuid.UUID) (*domain.UserResponse, error) {
	u, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCSUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return toUserResponse(u), nil
}

func (uc *csUserUseCase) GetUserTickets(ctx context.Context, customerID uuid.UUID, params domain.TicketListParams) ([]domain.TicketResponse, *domain.PaginationMeta, error) {
	params.CustomerID = &customerID
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

// enrichTicket populates CustomerEmail and AssignedCSName from the users table.
func (uc *csUserUseCase) enrichTicket(ctx context.Context, t *domain.Ticket) {
	if customer, err := uc.userRepo.FindByID(ctx, t.CustomerID); err == nil && customer != nil {
		t.CustomerEmail = customer.Email
	}
	if t.AssignedCSID != nil {
		if cs, err := uc.userRepo.FindByID(ctx, *t.AssignedCSID); err == nil && cs != nil {
			t.AssignedCSName = &cs.FullName
		}
	}
}
