package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// ReviewHandler handles product review endpoints.
type ReviewHandler struct {
	useCase domain.ReviewUseCase
}

// NewReviewHandler creates a new ReviewHandler.
func NewReviewHandler(useCase domain.ReviewUseCase) *ReviewHandler {
	return &ReviewHandler{useCase: useCase}
}

// PresignImage handles POST /api/v1/users/reviews/presign.
func (h *ReviewHandler) PresignImage(c *gin.Context) {
	var req domain.PresignReviewImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	res, err := h.useCase.PresignImage(c.Request.Context(), req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "review image upload URL generated", res)
}

// Create handles POST /api/v1/users/reviews.
func (h *ReviewHandler) Create(c *gin.Context) {
	userID, ok := extractUserID(c)
	if !ok {
		response.BadRequest(c, "invalid user ID in token", nil)
		return
	}

	var req domain.CreateProductReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	res, err := h.useCase.Create(c.Request.Context(), userID, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.Created(c, "review created successfully", res)
}

// ListByProduct handles GET /api/v1/products/:id/reviews.
func (h *ReviewHandler) ListByProduct(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid product ID", "id must be a valid UUID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	items, meta, err := h.useCase.ListByProduct(c.Request.Context(), productID, domain.ProductReviewListParams{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "product reviews retrieved successfully", items, meta)
}

// GetSummary handles GET /api/v1/products/:id/review-summary.
func (h *ReviewHandler) GetSummary(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid product ID", "id must be a valid UUID")
		return
	}

	summary, err := h.useCase.GetSummary(c.Request.Context(), productID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "product review summary retrieved successfully", summary)
}
