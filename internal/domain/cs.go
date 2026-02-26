package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ============================================================
// Entities
// ============================================================

type Ticket struct {
	ID                    uuid.UUID  `json:"id"`
	TicketNumber          string     `json:"ticket_number"`
	CustomerID            uuid.UUID  `json:"customer_id"`
	AssignedCSID          *uuid.UUID `json:"assigned_cs_id"`
	OrderNumber           string     `json:"order_number"`
	Phone                 string     `json:"phone"`
	ReporterName          string     `json:"reporter_name"`
	Subject               string     `json:"subject"`
	Detail                string     `json:"detail"`
	Status                string     `json:"status"`
	Source                string     `json:"source"`
	AttachmentURL         *string    `json:"attachment_url"`
	AttachmentContentType *string    `json:"attachment_content_type"`
	ResolvedAt            *time.Time `json:"resolved_at"`
	ClosedAt              *time.Time `json:"closed_at"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`

	// Joined fields (not persisted directly on tickets table)
	CustomerEmail  string  `json:"-" gorm:"-"`
	AssignedCSName *string `json:"-" gorm:"-"`
}

type TicketMessage struct {
	ID             uuid.UUID `json:"id"`
	TicketID       uuid.UUID `json:"ticket_id"`
	SenderID       uuid.UUID `json:"sender_id"`
	Message        string    `json:"message"`
	IsFromCS       bool      `json:"is_from_cs"`
	IsInternalNote bool      `json:"is_internal_note"`
	CreatedAt      time.Time `json:"created_at"`
}

type TicketAttachment struct {
	ID        uuid.UUID  `json:"id"`
	TicketID  uuid.UUID  `json:"ticket_id"`
	MessageID *uuid.UUID `json:"message_id"`
	FileURL   string     `json:"file_url"`
	FileName  string     `json:"file_name"`
	FileType  string     `json:"file_type"`
	FileSize  int        `json:"file_size"`
	CreatedAt time.Time  `json:"created_at"`
}

type TicketStatusLog struct {
	ID        uuid.UUID `json:"id"`
	TicketID  uuid.UUID `json:"ticket_id"`
	ChangedBy uuid.UUID `json:"changed_by"`
	OldStatus *string   `json:"old_status"`
	NewStatus string    `json:"new_status"`
	Notes     *string   `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatConversation struct {
	ID            uuid.UUID  `json:"id"`
	InitiatorID   uuid.UUID  `json:"initiator_id"`
	ParticipantID uuid.UUID  `json:"participant_id"`
	Status        string     `json:"status"`
	LastMessageAt *time.Time `json:"last_message_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type ChatMessage struct {
	ID             uuid.UUID  `json:"id"`
	ConversationID uuid.UUID  `json:"conversation_id"`
	SenderID       uuid.UUID  `json:"sender_id"`
	Message        string     `json:"message"`
	AttachmentURL  *string    `json:"attachment_url"`
	IsRead         bool       `json:"is_read"`
	ReadAt         *time.Time `json:"read_at"`
	CreatedAt      time.Time  `json:"created_at"`
}

type ReplyTemplate struct {
	ID        uuid.UUID  `json:"id"`
	Title     string     `json:"title"`
	Shortcut  string     `json:"shortcut"`
	Category  string     `json:"category"`
	Content   string     `json:"content"`
	IsActive  bool       `json:"is_active"`
	CreatedBy uuid.UUID  `json:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// ============================================================
// Request / Response DTOs
// ============================================================

// --- Ticket ---

type PresignTicketAttachmentRequest struct {
	ContentType string `json:"content_type" binding:"required"`
}

type PresignTicketAttachmentResponse struct {
	UploadURL   string `json:"upload_url"`
	ObjectKey   string `json:"object_key"`
	ContentType string `json:"content_type"`
	ExpiresIn   int    `json:"expires_in"` // seconds
}

type CreateTicketRequest struct {
	OrderNumber   string `json:"order_number"   binding:"required"`
	Phone         string `json:"phone"          binding:"required"`
	ReporterName  string `json:"reporter_name"  binding:"required"`
	Subject       string `json:"subject"        binding:"required"`
	Detail        string `json:"detail"         binding:"required"`
	Source        string `json:"source"         binding:"required,oneof=app web"`
	AttachmentKey string `json:"attachment_key" binding:"required"`
}

type UpdateTicketStatusRequest struct {
	Status string  `json:"status" binding:"required,oneof=open on_progress resolved closed"`
	Notes  *string `json:"notes"`
}

type AddTicketMessageRequest struct {
	Message        string `json:"message"          binding:"required"`
	IsInternalNote bool   `json:"is_internal_note"`
}

type TicketListParams struct {
	Page         int
	Limit        int
	Status       string
	AssignedCSID *uuid.UUID
	CustomerID   *uuid.UUID
}

type TicketResponse struct {
	ID                    uuid.UUID  `json:"id"`
	TicketNumber          string     `json:"ticket_number"`
	CustomerID            uuid.UUID  `json:"customer_id"`
	CustomerEmail         string     `json:"customer_email"`
	AssignedCSID          *uuid.UUID `json:"assigned_cs_id"`
	AssignedCSName        *string    `json:"assigned_cs_name"`
	OrderNumber           string     `json:"order_number"`
	Phone                 string     `json:"phone"`
	ReporterName          string     `json:"reporter_name"`
	Subject               string     `json:"subject"`
	Detail                string     `json:"detail"`
	Status                string     `json:"status"`
	Source                string     `json:"source"`
	AttachmentURL         *string    `json:"attachment_url"`
	AttachmentContentType *string    `json:"attachment_content_type"`
	ResolvedAt            *time.Time `json:"resolved_at"`
	ClosedAt              *time.Time `json:"closed_at"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

type TicketMessageResponse struct {
	ID             uuid.UUID `json:"id"`
	TicketID       uuid.UUID `json:"ticket_id"`
	SenderID       uuid.UUID `json:"sender_id"`
	Message        string    `json:"message"`
	IsFromCS       bool      `json:"is_from_cs"`
	IsInternalNote bool      `json:"is_internal_note"`
	CreatedAt      time.Time `json:"created_at"`
}

// --- Chat ---

type StartConversationRequest struct {
	ParticipantID uuid.UUID `json:"participant_id" binding:"required"`
}

type SendChatMessageRequest struct {
	Message       string  `json:"message"        binding:"required"`
	AttachmentURL *string `json:"attachment_url"`
}

type MarkMessagesReadRequest struct {
	// intentionally empty — marks all unread messages in conversation as read
}

type ChatMessageListParams struct {
	Page           int
	Limit          int
	ConversationID uuid.UUID
}

type ConversationListParams struct {
	Page   int
	Limit  int
	Status string
	UserID uuid.UUID
}

type ChatableUsersParams struct {
	Q     string
	Roles []string // subset of allowed roles; empty = all allowed
	Page  int
	Limit int
}

type ChatableUserResponse struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"full_name"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
}

type ConversationResponse struct {
	ID            uuid.UUID  `json:"id"`
	InitiatorID   uuid.UUID  `json:"initiator_id"`
	ParticipantID uuid.UUID  `json:"participant_id"`
	Status        string     `json:"status"`
	LastMessageAt *time.Time `json:"last_message_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type ChatMessageResponse struct {
	ID             uuid.UUID  `json:"id"`
	ConversationID uuid.UUID  `json:"conversation_id"`
	SenderID       uuid.UUID  `json:"sender_id"`
	Message        string     `json:"message"`
	AttachmentURL  *string    `json:"attachment_url"`
	IsRead         bool       `json:"is_read"`
	ReadAt         *time.Time `json:"read_at"`
	CreatedAt      time.Time  `json:"created_at"`
}

// --- Reply Templates ---

// ValidReplyTemplateCategories is the allowed set for reply template categories.
var ValidReplyTemplateCategories = map[string]struct{}{
	"general": {}, "order": {}, "payment": {}, "refund": {}, "shipping": {},
	"product": {}, "account": {}, "complaint": {}, "return": {}, "promo": {},
	"voucher": {}, "technical": {}, "verification": {}, "vendor": {}, "stock": {},
	"cancellation": {}, "delivery": {}, "pickup": {}, "subscription": {},
	"review": {}, "fraud": {}, "other": {},
}

type CreateReplyTemplateRequest struct {
	Title    string `json:"title"    binding:"required,max=150"`
	Category string `json:"category" binding:"required"`
	SubKey   string `json:"sub_key"  binding:"required"`
	Content  string `json:"content"  binding:"required"`
}

type UpdateReplyTemplateRequest struct {
	Title    *string `json:"title"     binding:"omitempty,max=150"`
	SubKey   *string `json:"sub_key"` // recalculates shortcut; category unchanged
	Content  *string `json:"content"`
	IsActive *bool   `json:"is_active"`
}

type ReplyTemplateListParams struct {
	Page     int
	Limit    int
	Category string
	Search   string
	IsActive *bool
}

type ReplyTemplateResponse struct {
	ID        uuid.UUID  `json:"id"`
	Title     string     `json:"title"`
	Shortcut  string     `json:"shortcut"`
	Category  string     `json:"category"`
	Content   string     `json:"content"`
	IsActive  bool       `json:"is_active"`
	CreatedBy uuid.UUID  `json:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// --- Dashboard ---

type CSDashboardResponse struct {
	TotalTickets        int64 `json:"total_tickets"`
	OpenTickets         int64 `json:"open_tickets"`
	OnProgressTickets   int64 `json:"on_progress_tickets"`
	ResolvedTickets     int64 `json:"resolved_tickets"`
	ClosedTickets       int64 `json:"closed_tickets"`
	UnreadMessages      int64 `json:"unread_messages"`
	ActiveConversations int64 `json:"active_conversations"`
}

// --- Reports ---

type CSReportSummaryResponse struct {
	TotalTickets    int64   `json:"total_tickets"`
	ResolvedTickets int64   `json:"resolved_tickets"`
	ClosedTickets   int64   `json:"closed_tickets"`
	AvgResolutionHr float64 `json:"avg_resolution_hours"`
}

type CSTicketReportRow struct {
	TicketNumber string     `json:"ticket_number"`
	Subject      string     `json:"subject"`
	Status       string     `json:"status"`
	Source       string     `json:"source"`
	ReporterName string     `json:"reporter_name"`
	AssignedCS   *string    `json:"assigned_cs"`
	ResolvedAt   *time.Time `json:"resolved_at"`
	CreatedAt    time.Time  `json:"created_at"`
}

type CSReportTicketParams struct {
	Page      int
	Limit     int
	Status    string
	StartDate *time.Time
	EndDate   *time.Time
}

// --- CS User List (data pengguna) ---

type CSUserListParams struct {
	Page  int
	Limit int
	Query string
}

// ============================================================
// Repository Interfaces
// ============================================================

type TicketRepository interface {
	Create(ctx context.Context, ticket *Ticket) error
	FindByID(ctx context.Context, id uuid.UUID) (*Ticket, error)
	FindByTicketNumber(ctx context.Context, number string) (*Ticket, error)
	List(ctx context.Context, params TicketListParams) ([]Ticket, *PaginationMeta, error)
	Update(ctx context.Context, ticket *Ticket) error
	CountByStatus(ctx context.Context) (map[string]int64, error)

	// Messages
	CreateMessage(ctx context.Context, msg *TicketMessage) error
	ListMessages(ctx context.Context, ticketID uuid.UUID) ([]TicketMessage, error)

	// Status log
	CreateStatusLog(ctx context.Context, log *TicketStatusLog) error

	// Ticket count for number generation
	CountOnDate(ctx context.Context, date string) (int64, error)
}

type ChatRepository interface {
	FindConversation(ctx context.Context, userA, userB uuid.UUID) (*ChatConversation, error)
	CreateConversation(ctx context.Context, conv *ChatConversation) error
	ListConversations(ctx context.Context, params ConversationListParams) ([]ChatConversation, *PaginationMeta, error)
	FindConversationByID(ctx context.Context, id uuid.UUID) (*ChatConversation, error)
	UpdateConversation(ctx context.Context, conv *ChatConversation) error
	CountActiveConversations(ctx context.Context) (int64, error)

	CreateMessage(ctx context.Context, msg *ChatMessage) error
	ListMessages(ctx context.Context, params ChatMessageListParams) ([]ChatMessage, *PaginationMeta, error)
	MarkMessagesRead(ctx context.Context, conversationID, readerID uuid.UUID) error
	CountUnreadMessages(ctx context.Context, userID uuid.UUID) (int64, error)
}

type ReplyTemplateRepository interface {
	Create(ctx context.Context, tpl *ReplyTemplate) error
	FindByID(ctx context.Context, id uuid.UUID) (*ReplyTemplate, error)
	ExistsShortcut(ctx context.Context, shortcut string, excludeID *uuid.UUID) (bool, error)
	List(ctx context.Context, params ReplyTemplateListParams) ([]ReplyTemplate, *PaginationMeta, error)
	Update(ctx context.Context, tpl *ReplyTemplate) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// ============================================================
// UseCase Interfaces
// ============================================================

type TicketUseCase interface {
	// Customer: generate presigned URL untuk upload attachment
	PresignTicketAttachment(ctx context.Context, req PresignTicketAttachmentRequest) (*PresignTicketAttachmentResponse, error)
	// Customer: buat tiket baru
	CreateTicket(ctx context.Context, customerID uuid.UUID, req CreateTicketRequest) (*TicketResponse, error)
	// CS: list semua tiket
	ListTickets(ctx context.Context, params TicketListParams) ([]TicketResponse, *PaginationMeta, error)
	// CS / Customer: detail tiket
	GetTicket(ctx context.Context, id uuid.UUID) (*TicketResponse, error)
	// CS: ambil tiket (self-assign, first-come-first-served)
	TakeTicket(ctx context.Context, csID, ticketID uuid.UUID) (*TicketResponse, error)
	// CS: update status (hanya CS pemilik tiket)
	UpdateTicketStatus(ctx context.Context, csID, ticketID uuid.UUID, req UpdateTicketStatusRequest) (*TicketResponse, error)
	// CS: tambah pesan + kirim email ke customer (hanya CS pemilik tiket)
	AddTicketMessage(ctx context.Context, senderID uuid.UUID, isCS bool, ticketID uuid.UUID, req AddTicketMessageRequest) (*TicketMessageResponse, error)
	// CS: list pesan tiket
	ListTicketMessages(ctx context.Context, ticketID uuid.UUID) ([]TicketMessageResponse, error)
}

type ChatUseCase interface {
	// Mulai / ambil conversation antara dua user
	StartOrGetConversation(ctx context.Context, initiatorID uuid.UUID, initiatorRole string, req StartConversationRequest) (*ConversationResponse, error)
	// List conversation milik user
	ListConversations(ctx context.Context, params ConversationListParams) ([]ConversationResponse, *PaginationMeta, error)
	// Detail conversation
	GetConversation(ctx context.Context, conversationID, userID uuid.UUID) (*ConversationResponse, error)
	// Kirim pesan
	SendMessage(ctx context.Context, conversationID, senderID uuid.UUID, req SendChatMessageRequest) (*ChatMessageResponse, error)
	// List pesan dalam conversation
	ListMessages(ctx context.Context, params ChatMessageListParams, userID uuid.UUID) ([]ChatMessageResponse, *PaginationMeta, error)
	// Tandai semua pesan sudah dibaca
	MarkRead(ctx context.Context, conversationID, readerID uuid.UUID) error
	// Cari user yang boleh diajak chat oleh initiatorRole
	SearchChatableUsers(ctx context.Context, initiatorID uuid.UUID, initiatorRole string, params ChatableUsersParams) ([]ChatableUserResponse, *PaginationMeta, error)
}

type ReplyTemplateUseCase interface {
	Create(ctx context.Context, csID uuid.UUID, req CreateReplyTemplateRequest) (*ReplyTemplateResponse, error)
	List(ctx context.Context, params ReplyTemplateListParams) ([]ReplyTemplateResponse, *PaginationMeta, error)
	GetByID(ctx context.Context, id uuid.UUID) (*ReplyTemplateResponse, error)
	Update(ctx context.Context, csID, id uuid.UUID, req UpdateReplyTemplateRequest) (*ReplyTemplateResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CSDashboardUseCase interface {
	GetDashboard(ctx context.Context) (*CSDashboardResponse, error)
}

type CSReportUseCase interface {
	GetSummary(ctx context.Context) (*CSReportSummaryResponse, error)
	GetTicketReport(ctx context.Context, params CSReportTicketParams) ([]CSTicketReportRow, *PaginationMeta, error)
}

type CSUserUseCase interface {
	ListUsers(ctx context.Context, params CSUserListParams) ([]UserResponse, *PaginationMeta, error)
	GetUser(ctx context.Context, userID uuid.UUID) (*UserResponse, error)
	GetUserTickets(ctx context.Context, customerID uuid.UUID, params TicketListParams) ([]TicketResponse, *PaginationMeta, error)
}
