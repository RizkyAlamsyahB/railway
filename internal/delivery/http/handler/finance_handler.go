package handler

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/middleware"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// FinanceHandler handles finance endpoints.
type FinanceHandler struct {
	useCase domain.FinanceUseCase
}

type financeTransactionListData struct {
	Summary *domain.FinanceTransactionSummary `json:"summary"`
	Items   []domain.FinanceTransactionItem   `json:"items"`
}

type financePayoutListData struct {
	Summary *domain.FinancePayoutSummary `json:"summary"`
	Items   []domain.FinancePayoutItem   `json:"items"`
}

type financeRefundListData struct {
	Summary *domain.FinanceRefundSummary `json:"summary"`
	Items   []domain.FinanceRefundItem   `json:"items"`
}

// NewFinanceHandler creates a new FinanceHandler.
func NewFinanceHandler(uc domain.FinanceUseCase) *FinanceHandler {
	return &FinanceHandler{useCase: uc}
}

// Dashboard handles GET /api/v1/finance/dashboard
func (h *FinanceHandler) Dashboard(c *gin.Context) {
	month := c.Query("month")

	result, err := h.useCase.GetDashboard(c.Request.Context(), month)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "finance dashboard retrieved successfully", result)
}

// ListTransactions handles GET /api/v1/finance/transactions
func (h *FinanceHandler) ListTransactions(c *gin.Context) {
	params := h.parseListParams(c)

	items, meta, err := h.useCase.ListTransactions(c.Request.Context(), params)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	summary, err := h.useCase.GetTransactionSummary(c.Request.Context(), params.Month)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	data := financeTransactionListData{
		Summary: summary,
		Items:   items,
	}

	response.SuccessWithMeta(c, http.StatusOK, "transactions retrieved successfully", data, meta)
}

// ExportTransactions handles GET /api/v1/finance/transactions/export
func (h *FinanceHandler) ExportTransactions(c *gin.Context) {
	params := h.parseListParams(c)

	items, err := h.useCase.ExportTransactions(c.Request.Context(), params)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	header := []string{"Tanggal", "Invoice", "Customer", "vendor", "Metode", "Nominal", "Status"}
	records := make([][]string, 0, len(items)+1)
	for _, item := range items {
		method := ""
		if item.Method != nil {
			method = *item.Method
		}
		records = append(records, []string{
			formatTime(item.Date),
			item.Invoice,
			item.Customer,
			item.Vendor,
			method,
			formatAmount(item.Amount),
			item.Status,
		})
	}

	h.writeCSV(c, "transactions.csv", header, records)
}

// ListPayouts handles GET /api/v1/finance/payouts
func (h *FinanceHandler) ListPayouts(c *gin.Context) {
	params := h.parseListParams(c)

	items, meta, err := h.useCase.ListPayouts(c.Request.Context(), params)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	summary, err := h.useCase.GetPayoutSummary(c.Request.Context(), params.Month)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	data := financePayoutListData{
		Summary: summary,
		Items:   items,
	}

	response.SuccessWithMeta(c, http.StatusOK, "payouts retrieved successfully", data, meta)
}

// ExportPayouts handles GET /api/v1/finance/payouts/export
func (h *FinanceHandler) ExportPayouts(c *gin.Context) {
	params := h.parseListParams(c)

	items, err := h.useCase.ExportPayouts(c.Request.Context(), params)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	header := []string{"Vendor", "Kategori", "Order Selesai", "Nominal", "Komisi", "Net Payout", "Status"}
	records := make([][]string, 0, len(items)+1)
	for _, item := range items {
		records = append(records, []string{
			item.Vendor,
			item.VendorType,
			strconv.Itoa(item.OrderCompleted),
			formatAmount(item.Nominal),
			formatAmount(item.Commission),
			formatAmount(item.NetPayout),
			item.Status,
		})
	}

	h.writeCSV(c, "payouts.csv", header, records)
}

// ListRefunds handles GET /api/v1/finance/refunds
func (h *FinanceHandler) ListRefunds(c *gin.Context) {
	params := h.parseListParams(c)

	items, meta, err := h.useCase.ListRefunds(c.Request.Context(), params)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	summary, err := h.useCase.GetRefundSummary(c.Request.Context(), params.Month)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	data := financeRefundListData{
		Summary: summary,
		Items:   items,
	}

	response.SuccessWithMeta(c, http.StatusOK, "refunds retrieved successfully", data, meta)
}

// ExportRefunds handles GET /api/v1/finance/refunds/export
func (h *FinanceHandler) ExportRefunds(c *gin.Context) {
	params := h.parseListParams(c)

	items, err := h.useCase.ExportRefunds(c.Request.Context(), params)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	header := []string{"ID", "Customer", "Vendor", "Order ID", "Alasan", "Nominal", "Status"}
	records := make([][]string, 0, len(items)+1)
	for _, item := range items {
		reason := ""
		if item.Reason != nil {
			reason = *item.Reason
		}
		records = append(records, []string{
			item.ID.String(),
			item.Customer,
			item.Vendor,
			item.OrderID.String(),
			reason,
			formatAmount(item.Amount),
			item.Status,
		})
	}

	h.writeCSV(c, "refunds.csv", header, records)
}

// FinancialReport handles GET /api/v1/finance/reports
func (h *FinanceHandler) FinancialReport(c *gin.Context) {
	month := c.Query("month")

	result, err := h.useCase.GetFinancialReport(c.Request.Context(), month)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "financial report retrieved successfully", result)
}

// ExportFinancialReport handles GET /api/v1/finance/reports/export
func (h *FinanceHandler) ExportFinancialReport(c *gin.Context) {
	month := c.Query("month")

	report, err := h.useCase.GetFinancialReport(c.Request.Context(), month)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	header := []string{"Tanggal", "Nominal"}
	records := make([][]string, 0, len(report.DailyIncome)+1)
	for _, item := range report.DailyIncome {
		records = append(records, []string{
			item.Date.UTC().Format("2006-01-02"),
			formatAmount(item.Amount),
		})
	}

	h.writeCSV(c, "financial_report.csv", header, records)
}

// UpdateRefundStatus handles PATCH /api/v1/finance/refunds/:id/status
func (h *FinanceHandler) UpdateRefundStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid refund ID", "id must be a valid UUID")
		return
	}

	var req domain.UpdateRefundStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	actorID, _ := c.Get(middleware.ContextKeyUserID)
	actorUUID, _ := actorID.(uuid.UUID)

	result, err := h.useCase.UpdateRefundStatus(c.Request.Context(), id, req.Status, actorUUID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "refund status updated successfully", result)
}

// UpdatePayoutStatus handles PATCH /api/v1/finance/payouts/:id/status
func (h *FinanceHandler) UpdatePayoutStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid payout ID", "id must be a valid UUID")
		return
	}

	var req domain.UpdatePayoutStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	actorID, _ := c.Get(middleware.ContextKeyUserID)
	actorUUID, _ := actorID.(uuid.UUID)

	result, err := h.useCase.UpdatePayoutStatus(c.Request.Context(), id, req.Status, actorUUID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "payout status updated successfully", result)
}

func (h *FinanceHandler) parseListParams(c *gin.Context) domain.FinanceListParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := strings.TrimSpace(c.Query("search"))
	if strings.EqualFold(search, "string") {
		// Ignore OpenAPI placeholder value so it doesn't unintentionally filter results.
		search = ""
	}

	return domain.FinanceListParams{
		Month:    c.Query("month"),
		DateFrom: c.Query("date_from"),
		DateTo:   c.Query("date_to"),
		Page:     page,
		Limit:    limit,
		Status:   c.Query("status"),
		Search:   search,
	}
}

func (h *FinanceHandler) writeCSV(c *gin.Context, filename string, header []string, records [][]string) {
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
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(http.StatusOK, "text/csv", buf.Bytes())
}

func formatTime(t time.Time) string {
	return t.UTC().Format("2006-01-02 15:04:05")
}

func formatAmount(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}
