package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// AdminContactHandler handles admin CRUD for admin contacts.
type AdminContactHandler struct {
	uc domain.AdminContactUseCase
}

// NewAdminContactHandler creates a new AdminContactHandler.
func NewAdminContactHandler(uc domain.AdminContactUseCase) *AdminContactHandler {
	return &AdminContactHandler{uc: uc}
}

// CreateAdminContact handles POST /admin/settings/contacts
func (h *AdminContactHandler) CreateAdminContact(c *gin.Context) {
	var req domain.CreateAdminContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	result, err := h.uc.CreateAdminContact(c.Request.Context(), req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.Created(c, "admin contact created", result)
}

// ListAdminContacts handles GET /admin/settings/contacts
func (h *AdminContactHandler) ListAdminContacts(c *gin.Context) {
	contacts, err := h.uc.ListAdminContacts(c.Request.Context())
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "admin contacts loaded", contacts)
}

// ListActiveAdminContacts handles GET /contacts (public)
func (h *AdminContactHandler) ListActiveAdminContacts(c *gin.Context) {
	contacts, err := h.uc.ListAdminContacts(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "failed to load admin contacts", nil)
		return
	}
	response.OK(c, "admin contacts loaded", contacts)
}

// UpdateAdminContact handles PATCH /admin/settings/contacts/:id
func (h *AdminContactHandler) UpdateAdminContact(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid contact id", nil)
		return
	}

	var req domain.UpdateAdminContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	result, err := h.uc.UpdateAdminContact(c.Request.Context(), id, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "admin contact updated", result)
}

// DeleteAdminContact handles DELETE /admin/settings/contacts/:id
func (h *AdminContactHandler) DeleteAdminContact(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid contact id", nil)
		return
	}

	if err := h.uc.DeleteAdminContact(c.Request.Context(), id); err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "admin contact deleted", nil)
}
