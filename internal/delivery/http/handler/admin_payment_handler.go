package handler

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// AdminPaymentHandler handles admin payment listing endpoints.
type AdminPaymentHandler struct {
	useCase domain.AdminPaymentUseCase
}

// NewAdminPaymentHandler creates a new AdminPaymentHandler.
func NewAdminPaymentHandler(uc domain.AdminPaymentUseCase) *AdminPaymentHandler {
	return &AdminPaymentHandler{useCase: uc}
}

func (h *AdminPaymentHandler) parseParams(c *gin.Context) domain.AdminPaymentListParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	params := domain.AdminPaymentListParams{
		Page:      page,
		Limit:     limit,
		Status:    c.Query("status"),
		Search:    c.Query("search"),
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
	}

	if start := c.Query("start_date"); start != "" {
		if t, err := time.Parse("2006-01-02", start); err == nil {
			params.StartDate = &t
		}
	}
	if end := c.Query("end_date"); end != "" {
		if t, err := time.Parse("2006-01-02", end); err == nil {
			endOfDay := t.AddDate(0, 0, 1)
			params.EndDate = &endOfDay
		}
	}

	return params
}

// List handles GET /api/v1/admin/payments
func (h *AdminPaymentHandler) List(c *gin.Context) {
	params := h.parseParams(c)

	payments, meta, err := h.useCase.List(c.Request.Context(), params)
	if err != nil {
		response.InternalServerError(c, "failed to list payments", err.Error())
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "payments retrieved successfully", payments, meta)
}

// Summary handles GET /api/v1/admin/payments/summary
func (h *AdminPaymentHandler) Summary(c *gin.Context) {
	params := h.parseParams(c)

	summary, err := h.useCase.Summary(c.Request.Context(), params)
	if err != nil {
		response.InternalServerError(c, "failed to get payment summary", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "payment summary retrieved successfully", summary)
}

// Export handles GET /api/v1/admin/payments/export
func (h *AdminPaymentHandler) Export(c *gin.Context) {
	params := h.parseParams(c)

	items, err := h.useCase.Export(c.Request.Context(), params)
	if err != nil {
		response.InternalServerError(c, "failed to export payments", err.Error())
		return
	}

	header := []string{"Tanggal", "Invoice ID", "Order No", "Customer", "Vendor", "Metode", "Nominal", "Status"}
	records := make([][]string, 0, len(items))
	for _, item := range items {
		method := ""
		if item.PaymentMethod != nil {
			method = *item.PaymentMethod
		}
		records = append(records, []string{
			item.CreatedAt.UTC().Format("2006-01-02 15:04:05"),
			item.InvoiceID.String(),
			item.OrderNo,
			item.CustomerName,
			item.VendorName,
			method,
			fmt.Sprintf("%.2f", item.Amount),
			item.Status,
		})
	}

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
	c.Header("Content-Disposition", "attachment; filename=\"payments.csv\"")
	c.Data(http.StatusOK, "text/csv", buf.Bytes())
}
