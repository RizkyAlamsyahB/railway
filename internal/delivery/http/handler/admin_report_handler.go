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
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// AdminReportHandler handles admin report endpoints.
type AdminReportHandler struct {
	useCase domain.AdminReportUseCase
}

// NewAdminReportHandler creates a new AdminReportHandler.
func NewAdminReportHandler(useCase domain.AdminReportUseCase) *AdminReportHandler {
	return &AdminReportHandler{useCase: useCase}
}

// GetReport handles GET /api/v1/admin/reports
func (h *AdminReportHandler) GetReport(c *gin.Context) {
	params, ok := parseAdminReportParams(c, false)
	if !ok {
		return
	}

	result, meta, err := h.useCase.GetReport(c.Request.Context(), params)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "admin report retrieved successfully", result, meta)
}

// ExportReport handles GET /api/v1/admin/reports/export
func (h *AdminReportHandler) ExportReport(c *gin.Context) {
	params, ok := parseAdminReportParams(c, true)
	if !ok {
		return
	}

	items, err := h.useCase.ExportReport(c.Request.Context(), params)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	header := []string{"Periode", "Total Transaksi", "DP", "Pelunasan", "Refund"}
	records := make([][]string, 0, len(items))
	for _, item := range items {
		records = append(records, []string{
			item.Period,
			formatAdminReportAmount(item.TotalTransaction),
			formatAdminReportAmount(item.DownPayment),
			formatAdminReportAmount(item.Settlement),
			formatAdminReportAmount(item.Refund),
		})
	}

	writeAdminReportCSV(c, fmt.Sprintf("admin-report-%d.csv", params.Year), header, records)
}

func parseAdminReportParams(c *gin.Context, allowDateRange bool) (domain.AdminReportParams, bool) {
	now := time.Now()
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	year := queryInt(c, "year", now.Year())
	if year < 2000 || year > 2100 {
		response.BadRequest(c, "year must be between 2000 and 2100", nil)
		return domain.AdminReportParams{}, false
	}

	search := strings.TrimSpace(c.Query("search"))
	if strings.EqualFold(search, "string") {
		search = ""
	}

	fromDate := strings.TrimSpace(c.Query("from_date"))
	toDate := strings.TrimSpace(c.Query("to_date"))
	if !allowDateRange {
		fromDate = ""
		toDate = ""
	}

	if fromDate != "" {
		if _, err := time.Parse("2006-01-02", fromDate); err != nil {
			response.BadRequest(c, "from_date must use YYYY-MM-DD format", nil)
			return domain.AdminReportParams{}, false
		}
	}
	if toDate != "" {
		if _, err := time.Parse("2006-01-02", toDate); err != nil {
			response.BadRequest(c, "to_date must use YYYY-MM-DD format", nil)
			return domain.AdminReportParams{}, false
		}
	}
	if fromDate != "" && toDate != "" {
		from, _ := time.Parse("2006-01-02", fromDate)
		to, _ := time.Parse("2006-01-02", toDate)
		if from.After(to) {
			response.BadRequest(c, "from_date must be before or equal to to_date", nil)
			return domain.AdminReportParams{}, false
		}
	}

	return domain.AdminReportParams{
		Year:      year,
		Page:      page,
		Limit:     limit,
		Search:    search,
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
		FromDate:  fromDate,
		ToDate:    toDate,
	}, true
}

func writeAdminReportCSV(c *gin.Context, filename string, header []string, records [][]string) {
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

func formatAdminReportAmount(amount float64) string {
	return strconv.FormatFloat(amount, 'f', 0, 64)
}
