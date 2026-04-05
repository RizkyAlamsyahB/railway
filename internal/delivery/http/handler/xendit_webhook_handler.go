package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

// XenditWebhookHandler handles Xendit webhook callbacks.
type XenditWebhookHandler struct {
	vendorUseCase            domain.VendorUseCase
	financeUseCase           domain.FinanceUseCase
	webhookVerificationToken string
}

// NewXenditWebhookHandler creates a new XenditWebhookHandler.
func NewXenditWebhookHandler(vendorUseCase domain.VendorUseCase, financeUseCase domain.FinanceUseCase, webhookVerificationToken string) *XenditWebhookHandler {
	return &XenditWebhookHandler{
		vendorUseCase:            vendorUseCase,
		financeUseCase:           financeUseCase,
		webhookVerificationToken: webhookVerificationToken,
	}
}

// Payout handles POST /api/v1/webhooks/xendit/payout.
func (h *XenditWebhookHandler) Payout(c *gin.Context) {
	callbackToken := c.GetHeader("x-callback-token")
	if callbackToken != h.webhookVerificationToken {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid callback token"})
		return
	}

	var payload domain.XenditPayoutWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid payload"})
		return
	}

	if err := h.vendorUseCase.HandlePayoutWebhook(c.Request.Context(), payload); err != nil {
		log.Printf(
			"[xendit-payout-webhook] processing failed: event=%s reference_id=%s payout_id=%s status=%s err=%v",
			payload.Event,
			payload.Data.ReferenceID,
			payload.Data.ID,
			payload.Data.Status,
			err,
		)
		c.JSON(http.StatusOK, gin.H{"message": "webhook received with errors", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "webhook processed"})
}

// RefundPayout handles POST /api/v1/webhooks/xendit/refund-payout.
func (h *XenditWebhookHandler) RefundPayout(c *gin.Context) {
	callbackToken := c.GetHeader("x-callback-token")
	if callbackToken != h.webhookVerificationToken {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid callback token"})
		return
	}

	var payload domain.XenditPayoutWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid payload"})
		return
	}

	if h.financeUseCase == nil {
		c.JSON(http.StatusOK, gin.H{"message": "webhook ignored"})
		return
	}

	if err := h.financeUseCase.HandleRefundPayoutWebhook(c.Request.Context(), payload); err != nil {
		log.Printf(
			"[xendit-refund-payout-webhook] processing failed: event=%s reference_id=%s payout_id=%s status=%s err=%v",
			payload.Event,
			payload.Data.ReferenceID,
			payload.Data.ID,
			payload.Data.Status,
			err,
		)
		c.JSON(http.StatusOK, gin.H{"message": "webhook received with errors", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "webhook processed"})
}

// RefundGateway handles POST /api/v1/webhooks/xendit/refund.
// This processes refund.succeeded / refund.failed events from the Xendit Refund API.
func (h *XenditWebhookHandler) RefundGateway(c *gin.Context) {
	callbackToken := c.GetHeader("x-callback-token")
	if callbackToken != h.webhookVerificationToken {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid callback token"})
		return
	}

	var payload domain.XenditRefundWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid payload"})
		return
	}

	if h.financeUseCase == nil {
		c.JSON(http.StatusOK, gin.H{"message": "webhook ignored"})
		return
	}

	if err := h.financeUseCase.HandleRefundGatewayWebhook(c.Request.Context(), payload); err != nil {
		log.Printf(
			"[xendit-refund-gateway-webhook] processing failed: event=%s reference_id=%s refund_id=%s status=%s err=%v",
			payload.Event,
			payload.Data.ReferenceID,
			payload.Data.ID,
			payload.Data.Status,
			err,
		)
		c.JSON(http.StatusOK, gin.H{"message": "webhook received with errors", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "webhook processed"})
}
