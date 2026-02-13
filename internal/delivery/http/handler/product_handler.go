package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/middleware"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase"
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
		h.handleError(c, err)
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
		h.handleError(c, err)
		return
	}

	response.OK(c, "images confirmed successfully", result)
}

func (h *ProductHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, usecase.ErrVendorNotFound):
		response.NotFound(c, "vendor not found", nil)
	case errors.Is(err, usecase.ErrVendorNotActive):
		response.Forbidden(c, "vendor is not active", nil)
	case errors.Is(err, usecase.ErrCategoryNotFound):
		response.BadRequest(c, "category not found", nil)
	case errors.Is(err, usecase.ErrShippingServiceNotFound):
		response.BadRequest(c, "one or more shipping services not found", nil)
	case errors.Is(err, usecase.ErrProductNotFound):
		response.NotFound(c, "product not found", nil)
	case errors.Is(err, usecase.ErrProductNotOwned):
		response.Forbidden(c, "product does not belong to this vendor", nil)
	case errors.Is(err, usecase.ErrImageNotFound):
		response.BadRequest(c, err.Error(), nil)
	case errors.Is(err, usecase.ErrImageNotUploaded):
		response.BadRequest(c, err.Error(), nil)
	case errors.Is(err, usecase.ErrInvalidImageContentType):
		response.BadRequest(c, err.Error(), nil)
	case errors.Is(err, usecase.ErrImageSizeOverflow):
		response.BadRequest(c, err.Error(), nil)
	case errors.Is(err, usecase.ErrDuplicatePrimaryImage):
		response.BadRequest(c, "only one image can be marked as primary", nil)
	case errors.Is(err, usecase.ErrNoPrimaryImage):
		response.BadRequest(c, "at least one image must be marked as primary", nil)
	case errors.Is(err, usecase.ErrPublishedRequiresShipping):
		response.BadRequest(c, "published product must have at least one shipping service", nil)
	case errors.Is(err, usecase.ErrTooManyImages):
		response.BadRequest(c, "maximum 10 images per product", nil)
	default:
		response.InternalServerError(c, "internal server error", nil)
	}
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
