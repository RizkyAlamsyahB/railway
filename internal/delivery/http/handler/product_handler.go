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

// UpdateProduct handles PUT /api/v1/vendors/products/:id
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid product ID", "id must be a valid UUID")
		return
	}

	var req domain.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	vendorIDVal, _ := c.Get(middleware.ContextKeyVendorID)
	vendorID, _ := vendorIDVal.(uuid.UUID)

	result, err := h.useCase.UpdateProduct(c.Request.Context(), vendorID, productID, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "product updated successfully", result)
}

// GetProductByID handles GET /api/v1/vendors/products/:id
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid product ID", "id must be a valid UUID")
		return
	}

	vendorIDVal, _ := c.Get(middleware.ContextKeyVendorID)
	vendorID, _ := vendorIDVal.(uuid.UUID)

	result, err := h.useCase.GetProductByID(c.Request.Context(), vendorID, productID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "product detail retrieved successfully", result)
}

// DeleteProduct handles DELETE /api/v1/vendors/products/:id
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid product ID", "id must be a valid UUID")
		return
	}

	vendorIDVal, _ := c.Get(middleware.ContextKeyVendorID)
	vendorID, _ := vendorIDVal.(uuid.UUID)

	if err := h.useCase.DeleteProduct(c.Request.Context(), vendorID, productID); err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "product deleted successfully", nil)
}

// ListProducts handles GET /api/v1/vendors/products
func (h *ProductHandler) ListProducts(c *gin.Context) {
	vendorIDVal, _ := c.Get(middleware.ContextKeyVendorID)
	vendorID, _ := vendorIDVal.(uuid.UUID)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	params := domain.VendorProductListParams{
		Page:      page,
		Limit:     limit,
		Search:    c.Query("search"),
		Status:    c.Query("status"),
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
	}

	items, meta, err := h.useCase.ListProducts(c.Request.Context(), vendorID, params)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "vendor products retrieved successfully", items, meta)
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

// ListProducts handles GET /api/v1/products
func (h *CatalogHandler) ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	params := domain.ProductListParams{
		Page:   page,
		Limit:  limit,
		Sort:   c.Query("sort"),
		Search: c.Query("search"),
	}

	items, meta, err := h.useCase.ListProducts(c.Request.Context(), params)
	if err != nil {
		response.InternalServerError(c, "failed to list products", nil)
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "products retrieved successfully", items, meta)
}

// GetProductByID handles GET /api/v1/products/:id
func (h *CatalogHandler) GetProductByID(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid product ID", "id must be a valid UUID")
		return
	}

	detail, err := h.useCase.GetProductDetail(c.Request.Context(), productID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "product detail retrieved successfully", detail)
}
