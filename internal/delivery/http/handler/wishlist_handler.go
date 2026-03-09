package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// WishlistHandler handles customer-facing wishlist HTTP requests.
type WishlistHandler struct {
	useCase domain.WishlistUseCase
}

// NewWishlistHandler creates a new WishlistHandler.
func NewWishlistHandler(useCase domain.WishlistUseCase) *WishlistHandler {
	return &WishlistHandler{useCase: useCase}
}

// AddItem handles POST /api/v1/users/wishlist/items.
func (h *WishlistHandler) AddItem(c *gin.Context) {
	userID, ok := extractUserID(c)
	if !ok {
		response.BadRequest(c, "invalid user ID in token", nil)
		return
	}

	var req domain.AddWishlistItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	result, err := h.useCase.AddItem(c.Request.Context(), userID, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "item added to wishlist", result)
}

// ListItems handles GET /api/v1/users/wishlist.
func (h *WishlistHandler) ListItems(c *gin.Context) {
	userID, ok := extractUserID(c)
	if !ok {
		response.BadRequest(c, "invalid user ID in token", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	items, meta, err := h.useCase.ListItems(c.Request.Context(), userID, domain.WishlistListParams{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "wishlist retrieved successfully", items, meta)
}

// RemoveItem handles DELETE /api/v1/users/wishlist/products/:productId.
func (h *WishlistHandler) RemoveItem(c *gin.Context) {
	userID, ok := extractUserID(c)
	if !ok {
		response.BadRequest(c, "invalid user ID in token", nil)
		return
	}

	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		response.BadRequest(c, "invalid product ID", "productId must be a valid UUID")
		return
	}

	if err := h.useCase.RemoveItem(c.Request.Context(), userID, productID); err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "item removed from wishlist", nil)
}

// GetProductStatus handles GET /api/v1/users/wishlist/products/:productId/status.
func (h *WishlistHandler) GetProductStatus(c *gin.Context) {
	userID, ok := extractUserID(c)
	if !ok {
		response.BadRequest(c, "invalid user ID in token", nil)
		return
	}

	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		response.BadRequest(c, "invalid product ID", "productId must be a valid UUID")
		return
	}

	result, err := h.useCase.GetProductStatus(c.Request.Context(), userID, productID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "wishlist status retrieved successfully", result)
}
