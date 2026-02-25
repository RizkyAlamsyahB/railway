package usecase

import "github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"

// toUserResponse converts a domain User entity (with preloaded Role) into
// a UserResponse DTO suitable for API output.
func toUserResponse(u *domain.User) *domain.UserResponse {
	role := ""
	if u.Role != nil {
		role = u.Role.Code
	}
	return &domain.UserResponse{
		ID:              u.ID,
		Email:           u.Email,
		FullName:        u.FullName,
		BirthDate:       u.BirthDate,
		Phone:           u.Phone,
		Status:          u.Status,
		EmailVerifiedAt: u.EmailVerifiedAt,
		Role:            role,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}

// ============================================================
// CS mappers
// ============================================================

func toTicketResponse(t domain.Ticket) *domain.TicketResponse {
	return &domain.TicketResponse{
		ID:                    t.ID,
		TicketNumber:          t.TicketNumber,
		CustomerID:            t.CustomerID,
		CustomerEmail:         t.CustomerEmail,
		AssignedCSID:          t.AssignedCSID,
		AssignedCSName:        t.AssignedCSName,
		OrderNumber:           t.OrderNumber,
		Phone:                 t.Phone,
		ReporterName:          t.ReporterName,
		Subject:               t.Subject,
		Detail:                t.Detail,
		Status:                t.Status,
		Source:                t.Source,
		AttachmentURL:         t.AttachmentURL,
		AttachmentContentType: t.AttachmentContentType,
		ResolvedAt:            t.ResolvedAt,
		ClosedAt:              t.ClosedAt,
		CreatedAt:             t.CreatedAt,
		UpdatedAt:             t.UpdatedAt,
	}
}

func toTicketMessageResponse(m domain.TicketMessage) *domain.TicketMessageResponse {
	return &domain.TicketMessageResponse{
		ID:             m.ID,
		TicketID:       m.TicketID,
		SenderID:       m.SenderID,
		Message:        m.Message,
		IsFromCS:       m.IsFromCS,
		IsInternalNote: m.IsInternalNote,
		CreatedAt:      m.CreatedAt,
	}
}

func toConversationResponse(c domain.ChatConversation) *domain.ConversationResponse {
	return &domain.ConversationResponse{
		ID:            c.ID,
		InitiatorID:   c.InitiatorID,
		ParticipantID: c.ParticipantID,
		Status:        c.Status,
		LastMessageAt: c.LastMessageAt,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
	}
}

func toChatMessageResponse(m domain.ChatMessage) *domain.ChatMessageResponse {
	return &domain.ChatMessageResponse{
		ID:             m.ID,
		ConversationID: m.ConversationID,
		SenderID:       m.SenderID,
		Message:        m.Message,
		AttachmentURL:  m.AttachmentURL,
		IsRead:         m.IsRead,
		ReadAt:         m.ReadAt,
		CreatedAt:      m.CreatedAt,
	}
}

func toReplyTemplateResponse(t domain.ReplyTemplate) *domain.ReplyTemplateResponse {
	return &domain.ReplyTemplateResponse{
		ID:        t.ID,
		Title:     t.Title,
		Shortcut:  t.Shortcut,
		Category:  t.Category,
		Content:   t.Content,
		IsActive:  t.IsActive,
		CreatedBy: t.CreatedBy,
		UpdatedBy: t.UpdatedBy,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
