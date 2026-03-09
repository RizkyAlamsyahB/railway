package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/middleware"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// AddressHandler handles HTTP requests for address management.
type AddressHandler struct {
	uc domain.AddressUseCase
}

// NewAddressHandler creates a new AddressHandler.
func NewAddressHandler(uc domain.AddressUseCase) *AddressHandler {
	return &AddressHandler{uc: uc}
}

// CreateAddress godoc
// POST /api/v1/addresses
func (h *AddressHandler) CreateAddress(c *gin.Context) {
	userID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	var req domain.CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	resp, err := h.uc.Create(c.Request.Context(), userID, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "address created", resp)
}

// ListAddresses godoc
// GET /api/v1/addresses
func (h *AddressHandler) ListAddresses(c *gin.Context) {
	userID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	addresses, err := h.uc.List(c.Request.Context(), userID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "addresses loaded", addresses)
}

// GetAddress godoc
// GET /api/v1/addresses/:id
func (h *AddressHandler) GetAddress(c *gin.Context) {
	userID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	addressID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid address id", nil)
		return
	}

	resp, err := h.uc.GetByID(c.Request.Context(), userID, addressID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "address loaded", resp)
}

// UpdateAddress godoc
// PATCH /api/v1/addresses/:id
func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	userID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	addressID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid address id", nil)
		return
	}

	var req domain.UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	resp, err := h.uc.Update(c.Request.Context(), userID, addressID, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "address updated", resp)
}

// DeleteAddress godoc
// DELETE /api/v1/addresses/:id
func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	userID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	addressID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid address id", nil)
		return
	}

	err = h.uc.Delete(c.Request.Context(), userID, addressID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "address deleted", nil)
}

// SetDefaultAddress godoc
// PATCH /api/v1/addresses/:id/default
func (h *AddressHandler) SetDefaultAddress(c *gin.Context) {
	userID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	addressID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid address id", nil)
		return
	}

	resp, err := h.uc.SetDefault(c.Request.Context(), userID, addressID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "default address set", resp)
}
