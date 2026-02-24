package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/middleware"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// ProductHandler handles product-related endpoints.
type ProductHandler struct {
	useCase domain.ProductUseCase
}

// NewProductHandler creates a new ProductHandler.
func NewProductHandler(uc domain.ProductUseCase) *ProductHandler {
	return &ProductHandler{useCase: uc}
}

// CreateProduct handles POST /api/v1/vendors/products
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req domain.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	vendorIDVal, _ := c.Get(middleware.ContextKeyVendorID)
	vendorID, _ := vendorIDVal.(uuid.UUID)

	result, err := h.useCase.Create(c.Request.Context(), vendorID, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.Created(c, "product created successfully", result)
}

// ConfirmImages handles POST /api/v1/vendors/products/:id/images/confirm
func (h *ProductHandler) ConfirmImages(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid product ID", "id must be a valid UUID")
		return
	}

	var req domain.ConfirmProductImagesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	vendorIDVal, _ := c.Get(middleware.ContextKeyVendorID)
	vendorID, _ := vendorIDVal.(uuid.UUID)

	result, err := h.useCase.ConfirmImages(c.Request.Context(), vendorID, productID, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "images confirmed successfully", result)
}

// CatalogHandler handles public catalog endpoints.
type CatalogHandler struct {
	useCase domain.CatalogUseCase
}

// NewCatalogHandler creates a new CatalogHandler.
func NewCatalogHandler(uc domain.CatalogUseCase) *CatalogHandler {
	return &CatalogHandler{useCase: uc}
}

// ListCategories handles GET /api/v1/categories
func (h *CatalogHandler) ListCategories(c *gin.Context) {
	categories, err := h.useCase.ListCategories(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "failed to list categories", nil)
		return
	}
	response.OK(c, "categories retrieved successfully", categories)
}

// ListShippingServices handles GET /api/v1/shipping-services
func (h *CatalogHandler) ListShippingServices(c *gin.Context) {
	services, err := h.useCase.ListShippingServices(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "failed to list shipping services", nil)
		return
	}
	response.OK(c, "shipping services retrieved successfully", services)
}
