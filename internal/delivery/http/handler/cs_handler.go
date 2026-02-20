package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/middleware"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// ============================================================
// Ticket Handler
// ============================================================

type TicketHandler struct {
	uc domain.TicketUseCase
}

func NewTicketHandler(uc domain.TicketUseCase) *TicketHandler {
	return &TicketHandler{uc: uc}
}

// ListTickets godoc
// GET /customer-service/tickets
func (h *TicketHandler) ListTickets(c *gin.Context) {
	params := domain.TicketListParams{
		Page:   queryInt(c, "page", 1),
		Limit:  queryInt(c, "limit", 10),
		Status: c.Query("status"),
	}
	if csIDStr := c.Query("assigned_cs_id"); csIDStr != "" {
		if id, err := uuid.Parse(csIDStr); err == nil {
			params.AssignedCSID = &id
		}
	}

	tickets, meta, err := h.uc.ListTickets(c.Request.Context(), params)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.SuccessWithMeta(c, http.StatusOK, "tickets retrieved", tickets, meta)
}

// GetTicket godoc
// GET /customer-service/tickets/:id
func (h *TicketHandler) GetTicket(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid ticket id", nil)
		return
	}
	ticket, err := h.uc.GetTicket(c.Request.Context(), id)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "ticket retrieved", ticket)
}

// UpdateTicketStatus godoc
// PATCH /customer-service/tickets/:id/status
func (h *TicketHandler) UpdateTicketStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid ticket id", nil)
		return
	}
	csID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	var req domain.UpdateTicketStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	ticket, err := h.uc.UpdateTicketStatus(c.Request.Context(), csID, id, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "ticket status updated", ticket)
}

// AssignTicket godoc
// PATCH /customer-service/tickets/:id/assign
func (h *TicketHandler) AssignTicket(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid ticket id", nil)
		return
	}
	csID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	var req domain.AssignTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	ticket, err := h.uc.AssignTicket(c.Request.Context(), csID, id, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "ticket assigned", ticket)
}

// AddTicketMessage godoc
// POST /customer-service/tickets/:id/messages
func (h *TicketHandler) AddTicketMessage(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid ticket id", nil)
		return
	}
	senderID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	var req domain.AddTicketMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	msg, err := h.uc.AddTicketMessage(c.Request.Context(), senderID, true, id, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.Created(c, "message added", msg)
}

// ListTicketMessages godoc
// GET /customer-service/tickets/:id/messages
func (h *TicketHandler) ListTicketMessages(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid ticket id", nil)
		return
	}
	msgs, err := h.uc.ListTicketMessages(c.Request.Context(), id)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "messages retrieved", msgs)
}

// CreateTicket godoc (customer endpoint)
// POST /tickets
func (h *TicketHandler) CreateTicket(c *gin.Context) {
	customerID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	var req domain.CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	ticket, err := h.uc.CreateTicket(c.Request.Context(), customerID, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.Created(c, "ticket created", ticket)
}

// ============================================================
// Chat Handler
// ============================================================

type ChatHandler struct {
	uc domain.ChatUseCase
}

func NewChatHandler(uc domain.ChatUseCase) *ChatHandler {
	return &ChatHandler{uc: uc}
}

// StartOrGetConversation godoc
// POST /customer-service/chat
func (h *ChatHandler) StartOrGetConversation(c *gin.Context) {
	initiatorID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)
	initiatorRole := c.MustGet(middleware.ContextKeyRole).(string)

	var req domain.StartConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	conv, err := h.uc.StartOrGetConversation(c.Request.Context(), initiatorID, initiatorRole, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.Created(c, "conversation ready", conv)
}

// ListConversations godoc
// GET /customer-service/chat
func (h *ChatHandler) ListConversations(c *gin.Context) {
	userID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	params := domain.ConversationListParams{
		Page:   queryInt(c, "page", 1),
		Limit:  queryInt(c, "limit", 10),
		Status: c.Query("status"),
		UserID: userID,
	}

	convs, meta, err := h.uc.ListConversations(c.Request.Context(), params)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.SuccessWithMeta(c, http.StatusOK, "conversations retrieved", convs, meta)
}

// GetConversation godoc
// GET /customer-service/chat/:conversationId
func (h *ChatHandler) GetConversation(c *gin.Context) {
	convID, err := uuid.Parse(c.Param("conversationId"))
	if err != nil {
		response.BadRequest(c, "invalid conversation id", nil)
		return
	}
	userID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	conv, err := h.uc.GetConversation(c.Request.Context(), convID, userID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "conversation retrieved", conv)
}

// SendMessage godoc
// POST /customer-service/chat/:conversationId/messages
func (h *ChatHandler) SendMessage(c *gin.Context) {
	convID, err := uuid.Parse(c.Param("conversationId"))
	if err != nil {
		response.BadRequest(c, "invalid conversation id", nil)
		return
	}
	senderID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	var req domain.SendChatMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	msg, err := h.uc.SendMessage(c.Request.Context(), convID, senderID, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.Created(c, "message sent", msg)
}

// ListMessages godoc
// GET /customer-service/chat/:conversationId/messages
func (h *ChatHandler) ListMessages(c *gin.Context) {
	convID, err := uuid.Parse(c.Param("conversationId"))
	if err != nil {
		response.BadRequest(c, "invalid conversation id", nil)
		return
	}
	userID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	params := domain.ChatMessageListParams{
		Page:           queryInt(c, "page", 1),
		Limit:          queryInt(c, "limit", 20),
		ConversationID: convID,
	}

	msgs, meta, err := h.uc.ListMessages(c.Request.Context(), params, userID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.SuccessWithMeta(c, http.StatusOK, "messages retrieved", msgs, meta)
}

// MarkRead godoc
// PATCH /customer-service/chat/:conversationId/read
func (h *ChatHandler) MarkRead(c *gin.Context) {
	convID, err := uuid.Parse(c.Param("conversationId"))
	if err != nil {
		response.BadRequest(c, "invalid conversation id", nil)
		return
	}
	userID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	if err := h.uc.MarkRead(c.Request.Context(), convID, userID); err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "messages marked as read", nil)
}

// ============================================================
// Reply Template Handler
// ============================================================

type ReplyTemplateHandler struct {
	uc domain.ReplyTemplateUseCase
}

func NewReplyTemplateHandler(uc domain.ReplyTemplateUseCase) *ReplyTemplateHandler {
	return &ReplyTemplateHandler{uc: uc}
}

// ListTemplates godoc
// GET /customer-service/reply-templates
func (h *ReplyTemplateHandler) ListTemplates(c *gin.Context) {
	params := domain.ReplyTemplateListParams{
		Page:  queryInt(c, "page", 1),
		Limit: queryInt(c, "limit", 10),
	}
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		v := isActiveStr == "true"
		params.IsActive = &v
	}

	tpls, meta, err := h.uc.List(c.Request.Context(), params)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.SuccessWithMeta(c, http.StatusOK, "templates retrieved", tpls, meta)
}

// GetTemplate godoc
// GET /customer-service/reply-templates/:id
func (h *ReplyTemplateHandler) GetTemplate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid template id", nil)
		return
	}
	tpl, err := h.uc.GetByID(c.Request.Context(), id)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "template retrieved", tpl)
}

// CreateTemplate godoc
// POST /customer-service/reply-templates
func (h *ReplyTemplateHandler) CreateTemplate(c *gin.Context) {
	csID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	var req domain.CreateReplyTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	tpl, err := h.uc.Create(c.Request.Context(), csID, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.Created(c, "template created", tpl)
}

// UpdateTemplate godoc
// PUT /customer-service/reply-templates/:id
func (h *ReplyTemplateHandler) UpdateTemplate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid template id", nil)
		return
	}
	csID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	var req domain.UpdateReplyTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	tpl, err := h.uc.Update(c.Request.Context(), csID, id, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "template updated", tpl)
}

// DeleteTemplate godoc
// DELETE /customer-service/reply-templates/:id
func (h *ReplyTemplateHandler) DeleteTemplate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid template id", nil)
		return
	}
	if err := h.uc.Delete(c.Request.Context(), id); err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "template deleted", nil)
}

// ============================================================
// Dashboard Handler
// ============================================================

type CSDashboardHandler struct {
	uc domain.CSDashboardUseCase
}

func NewCSDashboardHandler(uc domain.CSDashboardUseCase) *CSDashboardHandler {
	return &CSDashboardHandler{uc: uc}
}

// GetDashboard godoc
// GET /customer-service/dashboard
func (h *CSDashboardHandler) GetDashboard(c *gin.Context) {
	data, err := h.uc.GetDashboard(c.Request.Context())
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "dashboard retrieved", data)
}

// ============================================================
// CS User Handler  (data pengguna)
// ============================================================

type CSUserHandler struct {
	uc       domain.CSUserUseCase
	ticketUC domain.TicketUseCase
}

func NewCSUserHandler(uc domain.CSUserUseCase, ticketUC domain.TicketUseCase) *CSUserHandler {
	return &CSUserHandler{uc: uc, ticketUC: ticketUC}
}

// ListUsers godoc
// GET /customer-service/users
func (h *CSUserHandler) ListUsers(c *gin.Context) {
	params := domain.CSUserListParams{
		Page:  queryInt(c, "page", 1),
		Limit: queryInt(c, "limit", 10),
		Role:  c.Query("role"),
		Query: c.Query("q"),
	}

	users, meta, err := h.uc.ListUsers(c.Request.Context(), params)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.SuccessWithMeta(c, http.StatusOK, "users retrieved", users, meta)
}

// GetUser godoc
// GET /customer-service/users/:id
func (h *CSUserHandler) GetUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid user id", nil)
		return
	}
	user, err := h.uc.GetUser(c.Request.Context(), id)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "user retrieved", user)
}

// GetUserTickets godoc
// GET /customer-service/users/:id/tickets
func (h *CSUserHandler) GetUserTickets(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid user id", nil)
		return
	}
	params := domain.TicketListParams{
		Page:  queryInt(c, "page", 1),
		Limit: queryInt(c, "limit", 10),
	}
	tickets, meta, err := h.uc.GetUserTickets(c.Request.Context(), id, params)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.SuccessWithMeta(c, http.StatusOK, "tickets retrieved", tickets, meta)
}

// ============================================================
// Helper
// ============================================================

func queryInt(c *gin.Context, key string, defaultVal int) int {
	n, err := strconv.Atoi(c.DefaultQuery(key, strconv.Itoa(defaultVal)))
	if err != nil || n <= 0 {
		return defaultVal
	}
	return n
}
