package handler

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pdf/fpdf"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// VendorOrderHandler handles vendor order management endpoints.
type VendorOrderHandler struct {
	useCase domain.VendorOrderUseCase
}

// NewVendorOrderHandler creates a new VendorOrderHandler.
func NewVendorOrderHandler(useCase domain.VendorOrderUseCase) *VendorOrderHandler {
	return &VendorOrderHandler{useCase: useCase}
}

// ListOrders handles GET /api/v1/vendors/orders.
func (h *VendorOrderHandler) ListOrders(c *gin.Context) {
	vendorID, ok := extractVendorID(c)
	if !ok {
		response.BadRequest(c, "vendor ID not found in token", nil)
		return
	}

	params := h.parseListParams(c)

	items, meta, err := h.useCase.ListOrders(c.Request.Context(), vendorID, params)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "vendor orders retrieved successfully", items, meta)
}

// ExportOrders handles GET /api/v1/vendors/orders/export.
func (h *VendorOrderHandler) ExportOrders(c *gin.Context) {
	vendorID, ok := extractVendorID(c)
	if !ok {
		response.BadRequest(c, "vendor ID not found in token", nil)
		return
	}

	params := h.parseListParams(c)

	items, err := h.useCase.ExportOrders(c.Request.Context(), vendorID, params)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	// Generate filename with timestamp
	filename := fmt.Sprintf("orders_%s.csv", time.Now().Format("20060102_150405"))

	// Set headers for CSV download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "text/csv; charset=utf-8")

	// Write CSV
	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Write header
	header := []string{
		"No. Order", "Tanggal Order", "Nama Pembeli", "Nama Produk", "Varian",
		"Qty", "Harga Satuan", "Total Harga Item", "Subtotal", "Ongkir",
		"Biaya Platform", "Total Pembayaran", "Status",
	}
	if err := writer.Write(header); err != nil {
		response.InternalServerError(c, "failed to write CSV header", nil)
		return
	}

	// Write data rows
	for _, item := range items {
		row := []string{
			item.OrderNo,
			item.OrderDate.Format("2006-01-02 15:04:05"),
			item.CustomerName,
			item.ProductName,
			item.Variant,
			strconv.Itoa(item.Qty),
			fmt.Sprintf("%.0f", item.UnitPrice),
			fmt.Sprintf("%.0f", item.LineTotal),
			fmt.Sprintf("%.0f", item.Subtotal),
			fmt.Sprintf("%.0f", item.ShippingFee),
			fmt.Sprintf("%.0f", item.PlatformFee),
			fmt.Sprintf("%.0f", item.GrandTotal),
			item.Status,
		}
		if err := writer.Write(row); err != nil {
			response.InternalServerError(c, "failed to write CSV row", nil)
			return
		}
	}
}

// parseListParams parses query parameters for list/export endpoints.
func (h *VendorOrderHandler) parseListParams(c *gin.Context) domain.VendorOrderListParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	params := domain.VendorOrderListParams{
		Page:      page,
		Limit:     limit,
		Status:    c.Query("status"),
		Search:    c.Query("search"),
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
	}

	// Parse date_from
	if dateFromStr := c.Query("date_from"); dateFromStr != "" {
		if t, err := time.Parse("2006-01-02", dateFromStr); err == nil {
			params.DateFrom = &t
		}
	}

	// Parse date_to
	if dateToStr := c.Query("date_to"); dateToStr != "" {
		if t, err := time.Parse("2006-01-02", dateToStr); err == nil {
			params.DateTo = &t
		}
	}

	return params
}

// GetOrderDetail handles GET /api/v1/vendors/orders/:orderId.
func (h *VendorOrderHandler) GetOrderDetail(c *gin.Context) {
	vendorID, ok := extractVendorID(c)
	if !ok {
		response.BadRequest(c, "vendor ID not found in token", nil)
		return
	}

	orderID, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		response.BadRequest(c, "invalid order ID", "orderId must be a valid UUID")
		return
	}

	res, err := h.useCase.GetOrderDetail(c.Request.Context(), vendorID, orderID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "order detail retrieved successfully", res)
}

// AcceptOrder handles POST /api/v1/vendors/orders/:orderId/accept.
func (h *VendorOrderHandler) AcceptOrder(c *gin.Context) {
	vendorID, ok := extractVendorID(c)
	if !ok {
		response.BadRequest(c, "vendor ID not found in token", nil)
		return
	}

	userID, ok := extractUserID(c)
	if !ok {
		response.BadRequest(c, "user ID not found in token", nil)
		return
	}

	orderID, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		response.BadRequest(c, "invalid order ID", "orderId must be a valid UUID")
		return
	}

	res, err := h.useCase.AcceptOrder(c.Request.Context(), vendorID, orderID, userID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "order accepted successfully", res)
}

// RejectOrder handles POST /api/v1/vendors/orders/:orderId/reject.
func (h *VendorOrderHandler) RejectOrder(c *gin.Context) {
	vendorID, ok := extractVendorID(c)
	if !ok {
		response.BadRequest(c, "vendor ID not found in token", nil)
		return
	}

	userID, ok := extractUserID(c)
	if !ok {
		response.BadRequest(c, "user ID not found in token", nil)
		return
	}

	orderID, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		response.BadRequest(c, "invalid order ID", "orderId must be a valid UUID")
		return
	}

	res, err := h.useCase.RejectOrder(c.Request.Context(), vendorID, orderID, userID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "order rejected successfully", res)
}

// ShipOrder handles POST /api/v1/vendors/orders/:orderId/ship.
func (h *VendorOrderHandler) ShipOrder(c *gin.Context) {
	vendorID, ok := extractVendorID(c)
	if !ok {
		response.BadRequest(c, "vendor ID not found in token", nil)
		return
	}

	userID, ok := extractUserID(c)
	if !ok {
		response.BadRequest(c, "user ID not found in token", nil)
		return
	}

	orderID, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		response.BadRequest(c, "invalid order ID", "orderId must be a valid UUID")
		return
	}

	var req domain.ShipOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", err.Error())
		return
	}

	res, err := h.useCase.ShipOrder(c.Request.Context(), vendorID, orderID, userID, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "order shipped successfully", res)
}

// TrackWaybill handles GET /api/v1/vendors/orders/:orderId/tracking.
func (h *VendorOrderHandler) TrackWaybill(c *gin.Context) {
	vendorID, ok := extractVendorID(c)
	if !ok {
		response.BadRequest(c, "vendor ID not found in token", nil)
		return
	}

	orderID, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		response.BadRequest(c, "invalid order ID", "orderId must be a valid UUID")
		return
	}

	res, err := h.useCase.TrackWaybill(c.Request.Context(), vendorID, orderID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "tracking info retrieved successfully", res)
}

// GetOrderInvoice handles GET /api/v1/vendors/orders/:orderId/invoice.
func (h *VendorOrderHandler) GetOrderInvoice(c *gin.Context) {
	vendorID, ok := extractVendorID(c)
	if !ok {
		response.BadRequest(c, "vendor ID not found in token", nil)
		return
	}

	orderID, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		response.BadRequest(c, "invalid order ID", "orderId must be a valid UUID")
		return
	}

	res, err := h.useCase.GetOrderInvoice(c.Request.Context(), vendorID, orderID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "invoice retrieved successfully", res)
}

// DownloadOrderInvoice handles GET /api/v1/vendors/orders/:orderId/invoice/download.
func (h *VendorOrderHandler) DownloadOrderInvoice(c *gin.Context) {
	vendorID, ok := extractVendorID(c)
	if !ok {
		response.BadRequest(c, "vendor ID not found in token", nil)
		return
	}

	orderID, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		response.BadRequest(c, "invalid order ID", "orderId must be a valid UUID")
		return
	}

	invoice, err := h.useCase.GetOrderInvoice(c.Request.Context(), vendorID, orderID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	// Generate PDF
	pdfBytes, err := h.generateInvoicePDF(invoice)
	if err != nil {
		response.InternalServerError(c, "failed to generate PDF", nil)
		return
	}

	// Set headers for PDF download
	filename := fmt.Sprintf("%s.pdf", invoice.Invoice.InvoiceNumber)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Length", strconv.Itoa(len(pdfBytes)))

	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// generateInvoicePDF generates a PDF invoice from invoice data.
func (h *VendorOrderHandler) generateInvoicePDF(invoice *domain.VendorOrderInvoiceResponse) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(180, 10, "NOTA PESANAN")
	pdf.Ln(12)

	// Invoice info
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(90, 6, fmt.Sprintf("No. Invoice: %s", invoice.Invoice.InvoiceNumber))
	pdf.Cell(90, 6, fmt.Sprintf("Tanggal: %s", invoice.Invoice.IssuedAt))
	pdf.Ln(10)

	// Separator
	pdf.SetDrawColor(200, 200, 200)
	pdf.Line(15, pdf.GetY(), 195, pdf.GetY())
	pdf.Ln(5)

	// Customer info
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(180, 6, "Informasi Pembeli")
	pdf.Ln(6)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(180, 5, fmt.Sprintf("Nama: %s", invoice.Customer.Name))
	pdf.Ln(5)
	pdf.Cell(180, 5, fmt.Sprintf("Telepon: %s", invoice.Customer.Phone))
	pdf.Ln(5)
	address := fmt.Sprintf("Alamat: %s, %s, %s, %s %s",
		invoice.Customer.Address.Street,
		invoice.Customer.Address.District,
		invoice.Customer.Address.City,
		invoice.Customer.Address.Province,
		invoice.Customer.Address.PostalCode,
	)
	pdf.MultiCell(180, 5, address, "", "", false)
	pdf.Ln(5)

	// Payment info
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(180, 6, "Informasi Pembayaran")
	pdf.Ln(6)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(90, 5, fmt.Sprintf("Metode: %s", invoice.Payment.Method))
	pdf.Cell(90, 5, fmt.Sprintf("Dibayar: %s", invoice.Payment.PaidAt))
	pdf.Ln(10)

	// Items table header
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(70, 8, "Produk", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 8, "Varian", "1", 0, "L", true, 0, "")
	pdf.CellFormat(15, 8, "Qty", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "Harga", "1", 0, "R", true, 0, "")
	pdf.CellFormat(35, 8, "Subtotal", "1", 0, "R", true, 0, "")
	pdf.Ln(8)

	// Items rows
	pdf.SetFont("Arial", "", 9)
	for _, item := range invoice.Items {
		pdf.CellFormat(70, 7, truncateText(item.Name, 35), "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 7, truncateText(item.Variant, 15), "1", 0, "L", false, 0, "")
		pdf.CellFormat(15, 7, strconv.Itoa(item.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 7, formatCurrency(item.Price), "1", 0, "R", false, 0, "")
		pdf.CellFormat(35, 7, formatCurrency(item.Subtotal), "1", 0, "R", false, 0, "")
		pdf.Ln(7)
	}

	pdf.Ln(5)

	// Summary
	pdf.SetFont("Arial", "", 10)
	summaryX := 130.0
	pdf.SetX(summaryX)
	pdf.Cell(35, 6, "Subtotal:")
	pdf.Cell(30, 6, formatCurrency(invoice.Summary.Subtotal))
	pdf.Ln(6)

	pdf.SetX(summaryX)
	pdf.Cell(35, 6, "Ongkos Kirim:")
	pdf.Cell(30, 6, formatCurrency(invoice.Summary.ShippingFee))
	pdf.Ln(6)

	pdf.SetX(summaryX)
	pdf.Cell(35, 6, "Biaya Layanan:")
	pdf.Cell(30, 6, formatCurrency(invoice.Summary.ServiceFee))
	pdf.Ln(6)

	pdf.SetFont("Arial", "B", 11)
	pdf.SetX(summaryX)
	pdf.Cell(35, 7, "TOTAL:")
	pdf.Cell(30, 7, formatCurrency(invoice.Summary.Total))
	pdf.Ln(15)

	// Footer
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(128, 128, 128)
	pdf.Cell(180, 5, "Terima kasih telah berbelanja di HajiUmrohStore")

	// Output to buffer
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// truncateText truncates text to max length with ellipsis.
func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen-3] + "..."
}

// formatCurrency formats a float as Indonesian Rupiah.
func formatCurrency(amount float64) string {
	return fmt.Sprintf("Rp %.0f", amount)
}

// GetOrderShippingInfo handles GET /api/v1/vendors/orders/:orderId/shipping.
func (h *VendorOrderHandler) GetOrderShippingInfo(c *gin.Context) {
	vendorID, ok := extractVendorID(c)
	if !ok {
		response.BadRequest(c, "vendor ID not found in token", nil)
		return
	}

	orderID, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		response.BadRequest(c, "invalid order ID", "orderId must be a valid UUID")
		return
	}

	res, err := h.useCase.GetOrderShippingInfo(c.Request.Context(), vendorID, orderID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "shipping info retrieved successfully", res)
}
