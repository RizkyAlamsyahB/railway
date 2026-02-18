package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// CartHandler handles customer-facing cart HTTP requests.
type CartHandler struct {
	useCase domain.CartUseCase
}

// NewCartHandler creates a new CartHandler.
func NewCartHandler(useCase domain.CartUseCase) *CartHandler {
	return &CartHandler{useCase: useCase}
}

// GetCart handles GET /api/v1/users/cart.
func (h *CartHandler) GetCart(c *gin.Context) {
	userID, ok := extractUserID(c)
	if !ok {
		response.BadRequest(c, "invalid user ID in token", nil)
		return
	}

	result, err := h.useCase.GetCart(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.OK(c, "cart retrieved successfully", result)
}

// AddItem handles POST /api/v1/users/cart/items.
func (h *CartHandler) AddItem(c *gin.Context) {
	userID, ok := extractUserID(c)
	if !ok {
		response.BadRequest(c, "invalid user ID in token", nil)
		return
	}

	var req domain.AddCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	result, err := h.useCase.AddItem(c.Request.Context(), userID, req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.Created(c, "item added to cart", result)
}

// UpdateItem handles PATCH /api/v1/users/cart/items/:itemId.
func (h *CartHandler) UpdateItem(c *gin.Context) {
	userID, ok := extractUserID(c)
	if !ok {
		response.BadRequest(c, "invalid user ID in token", nil)
		return
	}

	itemID, err := uuid.Parse(c.Param("itemId"))
	if err != nil {
		response.BadRequest(c, "invalid item ID", "itemId must be a valid UUID")
		return
	}

	var req domain.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	result, err := h.useCase.UpdateItem(c.Request.Context(), userID, itemID, req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.OK(c, "cart item updated", result)
}

// RemoveItem handles DELETE /api/v1/users/cart/items/:itemId.
func (h *CartHandler) RemoveItem(c *gin.Context) {
	userID, ok := extractUserID(c)
	if !ok {
		response.BadRequest(c, "invalid user ID in token", nil)
		return
	}

	itemID, err := uuid.Parse(c.Param("itemId"))
	if err != nil {
		response.BadRequest(c, "invalid item ID", "itemId must be a valid UUID")
		return
	}

	if err := h.useCase.RemoveItem(c.Request.Context(), userID, itemID); err != nil {
		h.handleError(c, err)
		return
	}

	response.OK(c, "item removed from cart", nil)
}

// ClearCart handles DELETE /api/v1/users/cart.
func (h *CartHandler) ClearCart(c *gin.Context) {
	userID, ok := extractUserID(c)
	if !ok {
		response.BadRequest(c, "invalid user ID in token", nil)
		return
	}

	if err := h.useCase.ClearCart(c.Request.Context(), userID); err != nil {
		h.handleError(c, err)
		return
	}

	response.OK(c, "cart cleared", nil)
}

// handleError maps cart usecase sentinel errors to HTTP responses.
func (h *CartHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, usecase.ErrCartItemNotFound):
		response.NotFound(c, "cart item not found", nil)
	case errors.Is(err, usecase.ErrCartItemNotOwned):
		response.Forbidden(c, "cart item does not belong to your cart", nil)
	case errors.Is(err, usecase.ErrVariantNotFound):
		response.NotFound(c, "product variant not found", nil)
	case errors.Is(err, usecase.ErrVariantNotActive):
		response.BadRequest(c, "product variant is not active", nil)
	case errors.Is(err, usecase.ErrProductNotAvailable):
		response.BadRequest(c, "product is not available", nil)
	case errors.Is(err, usecase.ErrProductNotFound):
		response.NotFound(c, "product not found", nil)
	case errors.Is(err, usecase.ErrInsufficientStock):
		response.Error(c, http.StatusConflict, "insufficient stock for requested quantity", nil)
	default:
		response.InternalServerError(c, "an unexpected error occurred", nil)
	}
}
