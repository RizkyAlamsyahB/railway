package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
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
		HandleUsecaseError(c, err)
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
		HandleUsecaseError(c, err)
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
		HandleUsecaseError(c, err)
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
		HandleUsecaseError(c, err)
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
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "cart cleared", nil)
}
