package usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

// allowedChat defines which roles may initiate chat with which other roles.
// Keys and values are role codes as stored in the JWT / users table.
var allowedChat = map[string][]string{
	domain.RoleCustomer: {domain.RoleUMKM, domain.RoleCS},
	domain.RoleUMKM:     {domain.RoleCustomer, domain.RoleAdmin, domain.RoleCS, domain.RoleFinance},
	domain.RoleAdmin:    {domain.RoleUMKM, domain.RoleAdmin, domain.RoleCS, domain.RoleFinance},
	domain.RoleCS:       {domain.RoleCustomer, domain.RoleUMKM, domain.RoleAdmin, domain.RoleCS, domain.RoleFinance},
	domain.RoleFinance:  {domain.RoleUMKM, domain.RoleAdmin, domain.RoleCS, domain.RoleFinance},
}

func canChat(initiatorRole, participantRole string) bool {
	allowed, ok := allowedChat[initiatorRole]
	if !ok {
		return false
	}
	for _, r := range allowed {
		if r == participantRole {
			return true
		}
	}
	return false
}

type chatUseCase struct {
	chatRepo domain.ChatRepository
	userRepo domain.UserRepository
}

func NewChatUseCase(chatRepo domain.ChatRepository, userRepo domain.UserRepository) domain.ChatUseCase {
	return &chatUseCase{chatRepo: chatRepo, userRepo: userRepo}
}

func (uc *chatUseCase) StartOrGetConversation(ctx context.Context, initiatorID uuid.UUID, initiatorRole string, req domain.StartConversationRequest) (*domain.ConversationResponse, error) {
	if initiatorID == req.ParticipantID {
		return nil, ErrConversationNotAllowed
	}

	// Look up participant role
	participant, err := uc.userRepo.FindByID(ctx, req.ParticipantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCSUserNotFound
		}
		return nil, fmt.Errorf("failed to find participant: %w", err)
	}

	participantRole := ""
	if participant.Role != nil {
		participantRole = participant.Role.Code
	}
	if !canChat(initiatorRole, participantRole) {
		return nil, ErrConversationNotAllowed
	}

	// Try to find existing conversation (bi-directional)
	conv, err := uc.chatRepo.FindConversation(ctx, initiatorID, req.ParticipantID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check conversation: %w", err)
	}
	if conv != nil {
		return toConversationResponse(*conv), nil
	}

	// Create new conversation
	now := time.Now()
	conv = &domain.ChatConversation{
		ID:            uuid.New(),
		InitiatorID:   initiatorID,
		ParticipantID: req.ParticipantID,
		Status:        domain.ChatConvStatusOpen,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := uc.chatRepo.CreateConversation(ctx, conv); err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}
	return toConversationResponse(*conv), nil
}

func (uc *chatUseCase) ListConversations(ctx context.Context, params domain.ConversationListParams) ([]domain.ConversationResponse, *domain.PaginationMeta, error) {
	convs, meta, err := uc.chatRepo.ListConversations(ctx, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list conversations: %w", err)
	}
	resp := make([]domain.ConversationResponse, len(convs))
	for i, c := range convs {
		resp[i] = *toConversationResponse(c)
	}
	return resp, meta, nil
}

func (uc *chatUseCase) GetConversation(ctx context.Context, conversationID, userID uuid.UUID) (*domain.ConversationResponse, error) {
	conv, err := uc.chatRepo.FindConversationByID(ctx, conversationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrConversationNotFound
		}
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	// Verify the caller is a participant
	if conv.InitiatorID != userID && conv.ParticipantID != userID {
		return nil, ErrConversationUnauthorized
	}

	return toConversationResponse(*conv), nil
}

func (uc *chatUseCase) SendMessage(ctx context.Context, conversationID, senderID uuid.UUID, req domain.SendChatMessageRequest) (*domain.ChatMessageResponse, error) {
	conv, err := uc.chatRepo.FindConversationByID(ctx, conversationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrConversationNotFound
		}
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}
	if conv.InitiatorID != senderID && conv.ParticipantID != senderID {
		return nil, ErrConversationUnauthorized
	}

	now := time.Now()
	msg := &domain.ChatMessage{
		ID:             uuid.New(),
		ConversationID: conversationID,
		SenderID:       senderID,
		Message:        req.Message,
		AttachmentURL:  req.AttachmentURL,
		IsRead:         false,
		CreatedAt:      now,
	}
	if err := uc.chatRepo.CreateMessage(ctx, msg); err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	// Update last_message_at on conversation
	conv.LastMessageAt = &now
	conv.UpdatedAt = now
	_ = uc.chatRepo.UpdateConversation(ctx, conv)

	return toChatMessageResponse(*msg), nil
}

func (uc *chatUseCase) ListMessages(ctx context.Context, params domain.ChatMessageListParams, userID uuid.UUID) ([]domain.ChatMessageResponse, *domain.PaginationMeta, error) {
	// Verify caller is participant
	conv, err := uc.chatRepo.FindConversationByID(ctx, params.ConversationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrConversationNotFound
		}
		return nil, nil, fmt.Errorf("failed to get conversation: %w", err)
	}
	if conv.InitiatorID != userID && conv.ParticipantID != userID {
		return nil, nil, ErrConversationUnauthorized
	}

	msgs, meta, err := uc.chatRepo.ListMessages(ctx, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list messages: %w", err)
	}
	resp := make([]domain.ChatMessageResponse, len(msgs))
	for i, m := range msgs {
		resp[i] = *toChatMessageResponse(m)
	}
	return resp, meta, nil
}

func (uc *chatUseCase) MarkRead(ctx context.Context, conversationID, readerID uuid.UUID) error {
	conv, err := uc.chatRepo.FindConversationByID(ctx, conversationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrConversationNotFound
		}
		return fmt.Errorf("failed to get conversation: %w", err)
	}
	if conv.InitiatorID != readerID && conv.ParticipantID != readerID {
		return ErrConversationUnauthorized
	}
	return uc.chatRepo.MarkMessagesRead(ctx, conversationID, readerID)
}
func (uc *chatUseCase) SearchChatableUsers(ctx context.Context, initiatorID uuid.UUID, initiatorRole string, params domain.ChatableUsersParams) ([]domain.ChatableUserResponse, *domain.PaginationMeta, error) {
	// Determine target roles: allowed roles intersected with requested filter
	allowed := allowedChat[initiatorRole]
	targetRoles := allowed
	if len(params.Roles) > 0 {
		var intersection []string
		for _, r := range params.Roles {
			for _, a := range allowed {
				if r == a {
					intersection = append(intersection, r)
					break
				}
			}
		}
		targetRoles = intersection
	}

	page := params.Page
	if page < 1 {
		page = 1
	}
	limit := params.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}

	if len(targetRoles) == 0 {
		return []domain.ChatableUserResponse{}, &domain.PaginationMeta{Page: page, Limit: limit}, nil
	}

	// Fetch up to 200 users matching the query across all target roles, then paginate in-memory.
	var collected []domain.ChatableUserResponse
	seen := map[uuid.UUID]struct{}{}
	for _, role := range targetRoles {
		users, _, err := uc.userRepo.List(ctx, domain.UserListParams{
			Page: 1, Limit: 200, Role: role, Search: params.Q,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to list users: %w", err)
		}
		for _, u := range users {
			if u.ID == initiatorID {
				continue
			}
			if _, dup := seen[u.ID]; dup {
				continue
			}
			seen[u.ID] = struct{}{}
			roleCode := ""
			if u.Role != nil {
				roleCode = u.Role.Code
			}
			collected = append(collected, domain.ChatableUserResponse{
				ID:       u.ID,
				FullName: u.FullName,
				Email:    u.Email,
				Role:     roleCode,
			})
		}
	}

	total := int64(len(collected))
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if totalPages < 1 {
		totalPages = 1
	}
	start := (page - 1) * limit
	end := start + limit
	if start > len(collected) {
		start = len(collected)
	}
	if end > len(collected) {
		end = len(collected)
	}

	meta := &domain.PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: totalPages,
	}
	return collected[start:end], meta, nil
}
