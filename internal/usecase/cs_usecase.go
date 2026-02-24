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

func (uc *csDashboardUseCase) GetDashboard(ctx context.Context) (*domain.CSDashboardResponse, error) {
	statusCounts, err := uc.ticketRepo.CountByStatus(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count tickets: %w", err)
	}

	// We don't have a specific "CS user" here for unread count, so we use 0 as placeholder.
	// The handler should inject the caller's userID via context if needed.
	activeConvs, err := uc.chatRepo.CountActiveConversations(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count conversations: %w", err)
	}

	total := statusCounts[domain.TicketStatusOpen] +
		statusCounts[domain.TicketStatusOnProgress] +
		statusCounts[domain.TicketStatusResolved] +
		statusCounts[domain.TicketStatusClosed]

	return &domain.CSDashboardResponse{
		TotalTickets:        total,
		OpenTickets:         statusCounts[domain.TicketStatusOpen],
		OnProgressTickets:   statusCounts[domain.TicketStatusOnProgress],
		ResolvedTickets:     statusCounts[domain.TicketStatusResolved],
		ClosedTickets:       statusCounts[domain.TicketStatusClosed],
		ActiveConversations: activeConvs,
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
		Role:   params.Role,
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
		resp[i] = *toTicketResponse(t)
	}
	return resp, meta, nil
}
