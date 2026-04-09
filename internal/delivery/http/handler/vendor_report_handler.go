package handler

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/middleware"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// VendorReportHandler handles vendor revenue report endpoints.
type VendorReportHandler struct {
	useCase domain.VendorReportUseCase
}

// NewVendorReportHandler creates a new VendorReportHandler.
func NewVendorReportHandler(useCase domain.VendorReportUseCase) *VendorReportHandler {
	return &VendorReportHandler{useCase: useCase}
}

// GetReport handles GET /api/v1/vendors/reports
func (h *VendorReportHandler) GetReport(c *gin.Context) {
	vendorIDVal, _ := c.Get(middleware.ContextKeyVendorID)
	vendorID, _ := vendorIDVal.(uuid.UUID)

	params := h.parseParams(c)

	result, meta, err := h.useCase.GetReport(c.Request.Context(), vendorID, params)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "vendor report retrieved successfully", result, meta)
}

// ExportReport handles GET /api/v1/vendors/reports/export
func (h *VendorReportHandler) ExportReport(c *gin.Context) {
	vendorIDVal, _ := c.Get(middleware.ContextKeyVendorID)
	vendorID, _ := vendorIDVal.(uuid.UUID)

	params := h.parseParams(c)

	items, err := h.useCase.ExportReport(c.Request.Context(), vendorID, params)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	header := []string{"Tanggal", "Jumlah Transaksi", "Pendapatan Kotor", "Komisi Platform", "Pendapatan Bersih"}
	records := make([][]string, 0, len(items))
	for _, item := range items {
		records = append(records, []string{
			item.Date,
			strconv.FormatInt(item.TransactionCount, 10),
			fmt.Sprintf("%.2f", item.GrossRevenue),
			fmt.Sprintf("%.2f", item.PlatformCommission),
			fmt.Sprintf("%.2f", item.NetRevenue),
		})
	}

	h.writeCSV(c, "vendor-report.csv", header, records)
}

func (h *VendorReportHandler) parseParams(c *gin.Context) domain.VendorReportParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := strings.TrimSpace(c.Query("search"))
	if strings.EqualFold(search, "string") {
		search = ""
	}

	return domain.VendorReportParams{
		Month:    c.Query("month"),
		DateFrom: c.Query("date_from"),
		DateTo:   c.Query("date_to"),
		Page:     page,
		Limit:    limit,
		Search:   search,
	}
}

func (h *VendorReportHandler) writeCSV(c *gin.Context, filename string, header []string, records [][]string) {
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
