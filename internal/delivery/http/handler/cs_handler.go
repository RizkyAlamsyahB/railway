package handler

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/middleware"
	ws "github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/websocket"
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

// TakeTicket godoc
// PATCH /customer-service/tickets/:id/take
func (h *TicketHandler) TakeTicket(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid ticket id", nil)
		return
	}
	csID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	ticket, err := h.uc.TakeTicket(c.Request.Context(), csID, id)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "ticket taken", ticket)
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

// PresignTicketAttachment godoc
// POST /tickets/attachment/presign
func (h *TicketHandler) PresignTicketAttachment(c *gin.Context) {
	var req domain.PresignTicketAttachmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	result, err := h.uc.PresignTicketAttachment(c.Request.Context(), req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "presigned upload URL generated", result)
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
	uc  domain.ChatUseCase
	hub *ws.Hub
}

func NewChatHandler(uc domain.ChatUseCase, hub *ws.Hub) *ChatHandler {
	return &ChatHandler{uc: uc, hub: hub}
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
		Page:       queryInt(c, "page", 1),
		Limit:      queryInt(c, "limit", 10),
		UserID:     userID,
		Search:     c.Query("search"),
		RoleFilter: c.Query("role"),
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

	// Push real-time notification to the other conversation participant.
	if h.hub != nil {
		if conv, convErr := h.uc.GetConversation(c.Request.Context(), convID, senderID); convErr == nil {
			recipientID := conv.InitiatorID
			if recipientID == senderID {
				recipientID = conv.ParticipantID
			}
			if msgData, jsonErr := json.Marshal(msg); jsonErr == nil {
				h.hub.SendToUser(recipientID, ws.WSMessage{
					Type:           ws.TypeChatMessage,
					ConversationID: convID.String(),
					Data:           msgData,
				})
			}
		}
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

// SearchChatableUsers godoc
// GET /chat/users?q=&roles=cs,umkm&page=1&limit=20
func (h *ChatHandler) SearchChatableUsers(c *gin.Context) {
	initiatorID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)
	initiatorRole := c.MustGet(middleware.ContextKeyRole).(string)

	var roles []string
	if rolesStr := c.Query("roles"); rolesStr != "" {
		for _, r := range splitCSV(rolesStr) {
			if r != "" {
				roles = append(roles, r)
			}
		}
	}

	params := domain.ChatableUsersParams{
		Q:     c.Query("q"),
		Roles: roles,
		Page:  queryInt(c, "page", 1),
		Limit: queryInt(c, "limit", 20),
	}

	users, meta, err := h.uc.SearchChatableUsers(c.Request.Context(), initiatorID, initiatorRole, params)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.SuccessWithMeta(c, http.StatusOK, "users retrieved", users, meta)
}

// StartChatWithCS godoc
// POST /api/v1/chat/cs — auto-assign customer ke CS tersedia
func (h *ChatHandler) StartChatWithCS(c *gin.Context) {
	customerID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	conv, err := h.uc.StartChatWithCS(c.Request.Context(), customerID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.Created(c, "connected to customer service", conv)
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
		Page:     queryInt(c, "page", 1),
		Limit:    queryInt(c, "limit", 10),
		Category: c.Query("category"),
		Search:   c.Query("q"),
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
// GET /customer-service/dashboard?month=2&year=2026
func (h *CSDashboardHandler) GetDashboard(c *gin.Context) {
	now := time.Now()
	month := queryInt(c, "month", int(now.Month()))
	year := queryInt(c, "year", now.Year())

	if month < 1 || month > 12 {
		response.BadRequest(c, "month must be between 1 and 12", nil)
		return
	}
	if year < 2000 || year > 2100 {
		response.BadRequest(c, "year must be between 2000 and 2100", nil)
		return
	}

	params := domain.CSDashboardParams{Month: month, Year: year}
	data, err := h.uc.GetDashboard(c.Request.Context(), params)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "dashboard retrieved", data)
}

// ============================================================
// CS Report Handler  (laporan)
// ============================================================

type CSReportHandler struct {
	uc domain.CSReportUseCase
}

func NewCSReportHandler(uc domain.CSReportUseCase) *CSReportHandler {
	return &CSReportHandler{uc: uc}
}

// GetReport godoc
// GET /customer-service/reports?month=2026-01
func (h *CSReportHandler) GetReport(c *gin.Context) {
	month := c.Query("month")
	if month == "" {
		response.BadRequest(c, "query param 'month' is required (format: YYYY-MM)", nil)
		return
	}
	params := domain.CSReportParams{Month: month}
	data, err := h.uc.GetReport(c.Request.Context(), params)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "report retrieved", data)
}

// ExportReport godoc
// GET /customer-service/reports/export?month=2026-01
func (h *CSReportHandler) ExportReport(c *gin.Context) {
	month := c.Query("month")
	if month == "" {
		response.BadRequest(c, "query param 'month' is required (format: YYYY-MM)", nil)
		return
	}
	params := domain.CSReportParams{Month: month}
	rows, err := h.uc.ExportReport(c.Request.Context(), params)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	header := []string{"Ticket Number", "Subject", "Status", "Source", "Reporter", "Assigned CS", "Resolved At", "Created At"}
	records := make([][]string, len(rows))
	for i, r := range rows {
		cs := ""
		if r.AssignedCS != nil {
			cs = *r.AssignedCS
		}
		resolvedAt := ""
		if r.ResolvedAt != nil {
			resolvedAt = r.ResolvedAt.UTC().Format("2006-01-02 15:04:05")
		}
		records[i] = []string{
			r.TicketNumber,
			r.Subject,
			r.Status,
			r.Source,
			r.ReporterName,
			cs,
			resolvedAt,
			r.CreatedAt.UTC().Format("2006-01-02 15:04:05"),
		}
	}

	filename := fmt.Sprintf("laporan-tiket-%s.csv", month)
	writeCSV(c, filename, header, records)
}

// writeCSV writes a CSV response.
func writeCSV(c *gin.Context, filename string, header []string, records [][]string) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	if err := writer.Write(header); err != nil {
		response.InternalServerError(c, "failed to export CSV", err.Error())
		return
	}
	for _, record := range records {
		if err := writer.Write(record); err != nil {
			response.InternalServerError(c, "failed to export CSV", err.Error())
			return
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		response.InternalServerError(c, "failed to export CSV", err.Error())
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(http.StatusOK, "text/csv", buf.Bytes())
}

// ============================================================
// Ticket Subject Handler (list subjects for customer form)
// ============================================================

type TicketSubjectHandler struct {
	uc domain.TicketSubjectUseCase
}

func NewTicketSubjectHandler(uc domain.TicketSubjectUseCase) *TicketSubjectHandler {
	return &TicketSubjectHandler{uc: uc}
}

// ListSubjects godoc
// GET /api/v1/ticket-subjects
func (h *TicketSubjectHandler) ListSubjects(c *gin.Context) {
	subjects, err := h.uc.ListSubjects(c.Request.Context())
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "ticket subjects retrieved", subjects)
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

// splitCSV splits a comma-separated string into a trimmed slice.
func splitCSV(s string) []string {
	var out []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			part := s[start:i]
			// trim spaces
			for len(part) > 0 && part[0] == ' ' {
				part = part[1:]
			}
			for len(part) > 0 && part[len(part)-1] == ' ' {
				part = part[:len(part)-1]
			}
			if part != "" {
				out = append(out, part)
			}
			start = i + 1
		}
	}
	return out
}
