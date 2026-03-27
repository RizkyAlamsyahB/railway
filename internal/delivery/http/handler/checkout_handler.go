package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// CheckoutHandler handles checkout and webhook HTTP requests.
type CheckoutHandler struct {
	useCase                  domain.CheckoutUseCase
	webhookVerificationToken string
	xenditBypass             bool
}

// NewCheckoutHandler creates a new CheckoutHandler.
func NewCheckoutHandler(useCase domain.CheckoutUseCase, webhookVerificationToken string, xenditBypass bool) *CheckoutHandler {
	return &CheckoutHandler{
		useCase:                  useCase,
		webhookVerificationToken: webhookVerificationToken,
		xenditBypass:             xenditBypass,
	}
}

// Checkout handles POST /api/v1/users/checkout.
func (h *CheckoutHandler) Checkout(c *gin.Context) {
	userID, ok := extractUserID(c)
	if !ok {
		response.BadRequest(c, "invalid user ID in token", nil)
		return
	}

	var req domain.CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", err.Error())
		return
	}

	result, err := h.useCase.Checkout(c.Request.Context(), userID, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.Created(c, "checkout successful", result)
}

// Preview handles POST /api/v1/users/checkout/preview.
func (h *CheckoutHandler) Preview(c *gin.Context) {
	userID, ok := extractUserID(c)
	if !ok {
		response.BadRequest(c, "invalid user ID in token", nil)
		return
	}

	var req domain.CheckoutPreviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", err.Error())
		return
	}

	result, err := h.useCase.Preview(c.Request.Context(), userID, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "checkout preview loaded", result)
}

// Webhook handles POST /api/v1/webhooks/xendit/invoice.
func (h *CheckoutHandler) Webhook(c *gin.Context) {
	// 1. Verify x-callback-token.
	callbackToken := c.GetHeader("x-callback-token")
	if callbackToken != h.webhookVerificationToken {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid callback token"})
		return
	}

	// 2. Parse payload.
	var payload domain.XenditWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid payload"})
		return
	}

	// 3. Process webhook.
	if err := h.useCase.HandleWebhook(c.Request.Context(), payload); err != nil {
		// Log the error but return 200 to prevent Xendit retries for most errors.
		// Only return non-200 for truly unexpected infrastructure failures.
		c.JSON(http.StatusOK, gin.H{"message": "webhook received with errors", "error": err.Error()})
		return
	}

	// 4. Return 200 OK per Xendit requirements.
	c.JSON(http.StatusOK, gin.H{"message": "webhook processed"})
}

// SimulatePayment handles POST /api/v1/dev/simulate-payment/:orderId.
// Only available when XENDIT_BYPASS=true.
func (h *CheckoutHandler) SimulatePayment(c *gin.Context) {
	if !h.xenditBypass {
		response.Forbidden(c, "simulate payment is only available when XENDIT_BYPASS=true", nil)
		return
	}

	orderID, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		response.BadRequest(c, "invalid order ID", nil)
		return
	}

	if err := h.useCase.SimulatePayment(c.Request.Context(), orderID); err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "payment simulated successfully", nil)
}
