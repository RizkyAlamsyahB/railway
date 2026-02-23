package repository

import (
	"context"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

// ============================================================
// GORM model structs
// ============================================================

type ticketModel struct {
	ID                    string     `gorm:"column:id;primaryKey"`
	TicketNumber          string     `gorm:"column:ticket_number"`
	CustomerID            string     `gorm:"column:customer_id"`
	AssignedCSID          *string    `gorm:"column:assigned_cs_id"`
	OrderNumber           string     `gorm:"column:order_number"`
	Phone                 string     `gorm:"column:phone"`
	ReporterName          string     `gorm:"column:reporter_name"`
	Subject               string     `gorm:"column:subject"`
	Detail                string     `gorm:"column:detail"`
	Status                string     `gorm:"column:status"`
	Source                string     `gorm:"column:source"`
	AttachmentURL         *string    `gorm:"column:attachment_url"`
	AttachmentContentType *string    `gorm:"column:attachment_content_type"`
	ResolvedAt            *time.Time `gorm:"column:resolved_at"`
	ClosedAt              *time.Time `gorm:"column:closed_at"`
	CreatedAt             time.Time  `gorm:"column:created_at"`
	UpdatedAt             time.Time  `gorm:"column:updated_at"`
}

func (ticketModel) TableName() string { return "tickets" }

type ticketMessageModel struct {
	ID             string    `gorm:"column:id;primaryKey"`
	TicketID       string    `gorm:"column:ticket_id"`
	SenderID       string    `gorm:"column:sender_id"`
	Message        string    `gorm:"column:message"`
	IsFromCS       bool      `gorm:"column:is_from_cs"`
	IsInternalNote bool      `gorm:"column:is_internal_note"`
	CreatedAt      time.Time `gorm:"column:created_at"`
}

func (ticketMessageModel) TableName() string { return "ticket_messages" }

type ticketStatusLogModel struct {
	ID        string    `gorm:"column:id;primaryKey"`
	TicketID  string    `gorm:"column:ticket_id"`
	ChangedBy string    `gorm:"column:changed_by"`
	OldStatus *string   `gorm:"column:old_status"`
	NewStatus string    `gorm:"column:new_status"`
	Notes     *string   `gorm:"column:notes"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (ticketStatusLogModel) TableName() string { return "ticket_status_logs" }

type chatConversationModel struct {
	ID            string     `gorm:"column:id;primaryKey"`
	InitiatorID   string     `gorm:"column:initiator_id"`
	ParticipantID string     `gorm:"column:participant_id"`
	Status        string     `gorm:"column:status"`
	LastMessageAt *time.Time `gorm:"column:last_message_at"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at"`
}

func (chatConversationModel) TableName() string { return "chat_conversations" }

type chatMessageModel struct {
	ID             string     `gorm:"column:id;primaryKey"`
	ConversationID string     `gorm:"column:conversation_id"`
	SenderID       string     `gorm:"column:sender_id"`
	Message        string     `gorm:"column:message"`
	AttachmentURL  *string    `gorm:"column:attachment_url"`
	IsRead         bool       `gorm:"column:is_read"`
	ReadAt         *time.Time `gorm:"column:read_at"`
	CreatedAt      time.Time  `gorm:"column:created_at"`
}

func (chatMessageModel) TableName() string { return "chat_messages" }

type replyTemplateModel struct {
	ID        string    `gorm:"column:id;primaryKey"`
	Title     string    `gorm:"column:title"`
	Shortcut  string    `gorm:"column:shortcut"`
	Category  string    `gorm:"column:category"`
	Content   string    `gorm:"column:content"`
	IsActive  bool      `gorm:"column:is_active"`
	CreatedBy string    `gorm:"column:created_by"`
	UpdatedBy *string   `gorm:"column:updated_by"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (replyTemplateModel) TableName() string { return "reply_templates" }

// ============================================================
// Ticket Repository
// ============================================================

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) domain.TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) Create(ctx context.Context, t *domain.Ticket) error {
	m := ticketModel{
		ID:                    t.ID.String(),
		TicketNumber:          t.TicketNumber,
		CustomerID:            t.CustomerID.String(),
		OrderNumber:           t.OrderNumber,
		Phone:                 t.Phone,
		ReporterName:          t.ReporterName,
		Subject:               t.Subject,
		Detail:                t.Detail,
		Status:                t.Status,
		Source:                t.Source,
		AttachmentURL:         t.AttachmentURL,
		AttachmentContentType: t.AttachmentContentType,
		CreatedAt:             t.CreatedAt,
		UpdatedAt:             t.UpdatedAt,
	}
	if t.AssignedCSID != nil {
		s := t.AssignedCSID.String()
		m.AssignedCSID = &s
	}
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *ticketRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Ticket, error) {
	var m ticketModel
	if err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&m).Error; err != nil {
		return nil, err
	}
	return toTicketDomain(m), nil
}

func (r *ticketRepository) FindByTicketNumber(ctx context.Context, number string) (*domain.Ticket, error) {
	var m ticketModel
	if err := r.db.WithContext(ctx).Where("ticket_number = ?", number).First(&m).Error; err != nil {
		return nil, err
	}
	return toTicketDomain(m), nil
}

func (r *ticketRepository) List(ctx context.Context, p domain.TicketListParams) ([]domain.Ticket, *domain.PaginationMeta, error) {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 || p.Limit > 100 {
		p.Limit = 10
	}
	offset := (p.Page - 1) * p.Limit

	q := r.db.WithContext(ctx).Model(&ticketModel{})
	if p.Status != "" {
		q = q.Where("status = ?", p.Status)
	}
	if p.AssignedCSID != nil {
		q = q.Where("assigned_cs_id = ?", p.AssignedCSID.String())
	}
	if p.CustomerID != nil {
		q = q.Where("customer_id = ?", p.CustomerID.String())
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	var rows []ticketModel
	if err := q.Order("created_at DESC").Offset(offset).Limit(p.Limit).Find(&rows).Error; err != nil {
		return nil, nil, err
	}

	items := make([]domain.Ticket, len(rows))
	for i, row := range rows {
		items[i] = *toTicketDomain(row)
	}
	return items, &domain.PaginationMeta{
		Page:       p.Page,
		Limit:      p.Limit,
		TotalItems: total,
		TotalPages: int(math.Ceil(float64(total) / float64(p.Limit))),
	}, nil
}

func (r *ticketRepository) Update(ctx context.Context, t *domain.Ticket) error {
	m := ticketModel{
		ID:                    t.ID.String(),
		TicketNumber:          t.TicketNumber,
		CustomerID:            t.CustomerID.String(),
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
	if t.AssignedCSID != nil {
		s := t.AssignedCSID.String()
		m.AssignedCSID = &s
	}
	return r.db.WithContext(ctx).Save(&m).Error
}

func (r *ticketRepository) CountByStatus(ctx context.Context) (map[string]int64, error) {
	type row struct {
		Status string
		Count  int64
	}
	var rows []row
	if err := r.db.WithContext(ctx).
		Model(&ticketModel{}).
		Select("status, count(*) as count").
		Group("status").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	result := make(map[string]int64)
	for _, row := range rows {
		result[row.Status] = row.Count
	}
	return result, nil
}

func (r *ticketRepository) CreateMessage(ctx context.Context, msg *domain.TicketMessage) error {
	m := ticketMessageModel{
		ID:             msg.ID.String(),
		TicketID:       msg.TicketID.String(),
		SenderID:       msg.SenderID.String(),
		Message:        msg.Message,
		IsFromCS:       msg.IsFromCS,
		IsInternalNote: msg.IsInternalNote,
		CreatedAt:      msg.CreatedAt,
	}
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *ticketRepository) ListMessages(ctx context.Context, ticketID uuid.UUID) ([]domain.TicketMessage, error) {
	var rows []ticketMessageModel
	if err := r.db.WithContext(ctx).
		Where("ticket_id = ?", ticketID.String()).
		Order("created_at ASC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	msgs := make([]domain.TicketMessage, len(rows))
	for i, row := range rows {
		msgs[i] = toTicketMessageDomain(row)
	}
	return msgs, nil
}

func (r *ticketRepository) CreateStatusLog(ctx context.Context, log *domain.TicketStatusLog) error {
	m := ticketStatusLogModel{
		ID:        log.ID.String(),
		TicketID:  log.TicketID.String(),
		ChangedBy: log.ChangedBy.String(),
		OldStatus: log.OldStatus,
		NewStatus: log.NewStatus,
		Notes:     log.Notes,
		CreatedAt: log.CreatedAt,
	}
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *ticketRepository) CountOnDate(ctx context.Context, date string) (int64, error) {
	// date format: "2006-01-02"
	var count int64
	err := r.db.WithContext(ctx).
		Model(&ticketModel{}).
		Where("DATE(created_at) = ?", date).
		Count(&count).Error
	return count, err
}

// ============================================================
// Chat Repository
// ============================================================

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) domain.ChatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) FindConversation(ctx context.Context, userA, userB uuid.UUID) (*domain.ChatConversation, error) {
	var m chatConversationModel
	err := r.db.WithContext(ctx).
		Where(
			"(initiator_id = ? AND participant_id = ?) OR (initiator_id = ? AND participant_id = ?)",
			userA.String(), userB.String(), userB.String(), userA.String(),
		).
		First(&m).Error
	if err != nil {
		return nil, err
	}
	return toChatConvDomain(m), nil
}

func (r *chatRepository) CreateConversation(ctx context.Context, conv *domain.ChatConversation) error {
	m := chatConversationModel{
		ID:            conv.ID.String(),
		InitiatorID:   conv.InitiatorID.String(),
		ParticipantID: conv.ParticipantID.String(),
		Status:        conv.Status,
		LastMessageAt: conv.LastMessageAt,
		CreatedAt:     conv.CreatedAt,
		UpdatedAt:     conv.UpdatedAt,
	}
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *chatRepository) ListConversations(ctx context.Context, p domain.ConversationListParams) ([]domain.ChatConversation, *domain.PaginationMeta, error) {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 || p.Limit > 100 {
		p.Limit = 10
	}
	offset := (p.Page - 1) * p.Limit

	q := r.db.WithContext(ctx).Model(&chatConversationModel{}).
		Where("initiator_id = ? OR participant_id = ?", p.UserID.String(), p.UserID.String())
	if p.Status != "" {
		q = q.Where("status = ?", p.Status)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	var rows []chatConversationModel
	if err := q.Order("last_message_at DESC NULLS LAST").Offset(offset).Limit(p.Limit).Find(&rows).Error; err != nil {
		return nil, nil, err
	}

	items := make([]domain.ChatConversation, len(rows))
	for i, row := range rows {
		items[i] = *toChatConvDomain(row)
	}
	return items, &domain.PaginationMeta{
		Page:       p.Page,
		Limit:      p.Limit,
		TotalItems: total,
		TotalPages: int(math.Ceil(float64(total) / float64(p.Limit))),
	}, nil
}

func (r *chatRepository) FindConversationByID(ctx context.Context, id uuid.UUID) (*domain.ChatConversation, error) {
	var m chatConversationModel
	if err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&m).Error; err != nil {
		return nil, err
	}
	return toChatConvDomain(m), nil
}

func (r *chatRepository) UpdateConversation(ctx context.Context, conv *domain.ChatConversation) error {
	m := chatConversationModel{
		ID:            conv.ID.String(),
		InitiatorID:   conv.InitiatorID.String(),
		ParticipantID: conv.ParticipantID.String(),
		Status:        conv.Status,
		LastMessageAt: conv.LastMessageAt,
		CreatedAt:     conv.CreatedAt,
		UpdatedAt:     conv.UpdatedAt,
	}
	return r.db.WithContext(ctx).Save(&m).Error
}

func (r *chatRepository) CountActiveConversations(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&chatConversationModel{}).
		Where("status = ?", domain.ChatConvStatusOpen).
		Count(&count).Error
	return count, err
}

func (r *chatRepository) CreateMessage(ctx context.Context, msg *domain.ChatMessage) error {
	m := chatMessageModel{
		ID:             msg.ID.String(),
		ConversationID: msg.ConversationID.String(),
		SenderID:       msg.SenderID.String(),
		Message:        msg.Message,
		AttachmentURL:  msg.AttachmentURL,
		IsRead:         msg.IsRead,
		ReadAt:         msg.ReadAt,
		CreatedAt:      msg.CreatedAt,
	}
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *chatRepository) ListMessages(ctx context.Context, p domain.ChatMessageListParams) ([]domain.ChatMessage, *domain.PaginationMeta, error) {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 || p.Limit > 100 {
		p.Limit = 10
	}
	offset := (p.Page - 1) * p.Limit

	q := r.db.WithContext(ctx).Model(&chatMessageModel{}).
		Where("conversation_id = ?", p.ConversationID.String())

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	var rows []chatMessageModel
	if err := q.Order("created_at ASC").Offset(offset).Limit(p.Limit).Find(&rows).Error; err != nil {
		return nil, nil, err
	}

	items := make([]domain.ChatMessage, len(rows))
	for i, row := range rows {
		items[i] = toChatMessageDomain(row)
	}
	return items, &domain.PaginationMeta{
		Page:       p.Page,
		Limit:      p.Limit,
		TotalItems: total,
		TotalPages: int(math.Ceil(float64(total) / float64(p.Limit))),
	}, nil
}

func (r *chatRepository) MarkMessagesRead(ctx context.Context, conversationID, readerID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&chatMessageModel{}).
		Where("conversation_id = ? AND sender_id <> ? AND is_read = FALSE", conversationID.String(), readerID.String()).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

func (r *chatRepository) CountUnreadMessages(ctx context.Context, userID uuid.UUID) (int64, error) {
	// count unread messages in all conversations where userID is a participant but NOT the sender
	var count int64
	err := r.db.WithContext(ctx).
		Model(&chatMessageModel{}).
		Joins("JOIN chat_conversations cc ON cc.id = chat_messages.conversation_id").
		Where(
			"(cc.initiator_id = ? OR cc.participant_id = ?) AND chat_messages.sender_id <> ? AND chat_messages.is_read = FALSE",
			userID.String(), userID.String(), userID.String(),
		).
		Count(&count).Error
	return count, err
}

// ============================================================
// Reply Template Repository
// ============================================================

type replyTemplateRepository struct {
	db *gorm.DB
}

func NewReplyTemplateRepository(db *gorm.DB) domain.ReplyTemplateRepository {
	return &replyTemplateRepository{db: db}
}

func (r *replyTemplateRepository) Create(ctx context.Context, tpl *domain.ReplyTemplate) error {
	m := replyTemplateModel{
		ID:        tpl.ID.String(),
		Title:     tpl.Title,
		Shortcut:  tpl.Shortcut,
		Category:  tpl.Category,
		Content:   tpl.Content,
		IsActive:  tpl.IsActive,
		CreatedBy: tpl.CreatedBy.String(),
		CreatedAt: tpl.CreatedAt,
		UpdatedAt: tpl.UpdatedAt,
	}
	if tpl.UpdatedBy != nil {
		s := tpl.UpdatedBy.String()
		m.UpdatedBy = &s
	}
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *replyTemplateRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.ReplyTemplate, error) {
	var m replyTemplateModel
	if err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&m).Error; err != nil {
		return nil, err
	}
	return toReplyTemplateDomain(m), nil
}

func (r *replyTemplateRepository) ExistsShortcut(ctx context.Context, shortcut string, excludeID *uuid.UUID) (bool, error) {
	q := r.db.WithContext(ctx).Model(&replyTemplateModel{}).Where("shortcut = ?", shortcut)
	if excludeID != nil {
		q = q.Where("id != ?", excludeID.String())
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *replyTemplateRepository) List(ctx context.Context, p domain.ReplyTemplateListParams) ([]domain.ReplyTemplate, *domain.PaginationMeta, error) {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 || p.Limit > 100 {
		p.Limit = 10
	}
	offset := (p.Page - 1) * p.Limit

	q := r.db.WithContext(ctx).Model(&replyTemplateModel{})
	if p.IsActive != nil {
		q = q.Where("is_active = ?", *p.IsActive)
	}
	if p.Category != "" {
		q = q.Where("category = ?", p.Category)
	}
	if p.Search != "" {
		like := "%" + p.Search + "%"
		q = q.Where("title ILIKE ? OR shortcut ILIKE ?", like, like)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	var rows []replyTemplateModel
	if err := q.Order("created_at DESC").Offset(offset).Limit(p.Limit).Find(&rows).Error; err != nil {
		return nil, nil, err
	}

	items := make([]domain.ReplyTemplate, len(rows))
	for i, row := range rows {
		items[i] = *toReplyTemplateDomain(row)
	}
	return items, &domain.PaginationMeta{
		Page:       p.Page,
		Limit:      p.Limit,
		TotalItems: total,
		TotalPages: int(math.Ceil(float64(total) / float64(p.Limit))),
	}, nil
}

func (r *replyTemplateRepository) Update(ctx context.Context, tpl *domain.ReplyTemplate) error {
	m := replyTemplateModel{
		ID:        tpl.ID.String(),
		Title:     tpl.Title,
		Shortcut:  tpl.Shortcut,
		Category:  tpl.Category,
		Content:   tpl.Content,
		IsActive:  tpl.IsActive,
		CreatedBy: tpl.CreatedBy.String(),
		CreatedAt: tpl.CreatedAt,
		UpdatedAt: tpl.UpdatedAt,
	}
	if tpl.UpdatedBy != nil {
		s := tpl.UpdatedBy.String()
		m.UpdatedBy = &s
	}
	return r.db.WithContext(ctx).Save(&m).Error
}

func (r *replyTemplateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id.String()).
		Delete(&replyTemplateModel{}).Error
}

// ============================================================
// Domain mappers
// ============================================================

func toTicketDomain(m ticketModel) *domain.Ticket {
	t := &domain.Ticket{
		ID:                    mustParseUUID(m.ID),
		TicketNumber:          m.TicketNumber,
		CustomerID:            mustParseUUID(m.CustomerID),
		OrderNumber:           m.OrderNumber,
		Phone:                 m.Phone,
		ReporterName:          m.ReporterName,
		Subject:               m.Subject,
		Detail:                m.Detail,
		Status:                m.Status,
		Source:                m.Source,
		AttachmentURL:         m.AttachmentURL,
		AttachmentContentType: m.AttachmentContentType,
		ResolvedAt:            m.ResolvedAt,
		ClosedAt:              m.ClosedAt,
		CreatedAt:             m.CreatedAt,
		UpdatedAt:             m.UpdatedAt,
	}
	if m.AssignedCSID != nil {
		id := mustParseUUID(*m.AssignedCSID)
		t.AssignedCSID = &id
	}
	return t
}

func toTicketMessageDomain(m ticketMessageModel) domain.TicketMessage {
	return domain.TicketMessage{
		ID:             mustParseUUID(m.ID),
		TicketID:       mustParseUUID(m.TicketID),
		SenderID:       mustParseUUID(m.SenderID),
		Message:        m.Message,
		IsFromCS:       m.IsFromCS,
		IsInternalNote: m.IsInternalNote,
		CreatedAt:      m.CreatedAt,
	}
}

func toChatConvDomain(m chatConversationModel) *domain.ChatConversation {
	return &domain.ChatConversation{
		ID:            mustParseUUID(m.ID),
		InitiatorID:   mustParseUUID(m.InitiatorID),
		ParticipantID: mustParseUUID(m.ParticipantID),
		Status:        m.Status,
		LastMessageAt: m.LastMessageAt,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func toChatMessageDomain(m chatMessageModel) domain.ChatMessage {
	return domain.ChatMessage{
		ID:             mustParseUUID(m.ID),
		ConversationID: mustParseUUID(m.ConversationID),
		SenderID:       mustParseUUID(m.SenderID),
		Message:        m.Message,
		AttachmentURL:  m.AttachmentURL,
		IsRead:         m.IsRead,
		ReadAt:         m.ReadAt,
		CreatedAt:      m.CreatedAt,
	}
}

func toReplyTemplateDomain(m replyTemplateModel) *domain.ReplyTemplate {
	tpl := &domain.ReplyTemplate{
		ID:        mustParseUUID(m.ID),
		Title:     m.Title,
		Shortcut:  m.Shortcut,
		Category:  m.Category,
		Content:   m.Content,
		IsActive:  m.IsActive,
		CreatedBy: mustParseUUID(m.CreatedBy),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
	if m.UpdatedBy != nil {
		id := mustParseUUID(*m.UpdatedBy)
		tpl.UpdatedBy = &id
	}
	return tpl
}

// ============================================================
// Helpers
// ============================================================

func mustParseUUID(s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return id
}
