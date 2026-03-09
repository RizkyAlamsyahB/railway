package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// CompleteOrderResponse is the output DTO for customer order completion.
type CompleteOrderResponse struct {
	OrderID       uuid.UUID `json:"order_id"`
	OrderStatus   string    `json:"order_status"`
	PaymentStatus string    `json:"payment_status"`
	CompletedAt   time.Time `json:"completed_at"`
}

// CustomerOrderListParams holds query parameters for customer order listing.
type CustomerOrderListParams struct {
	Page   int
	Limit  int
	Status string
}

// CustomerOrderProductItem is product and variant info in a customer order row.
type CustomerOrderProductItem struct {
	ProductName     string `json:"product_name"`
	SelectedVariant string `json:"selected_variant"`
	Qty             int    `json:"qty"`
}

// CustomerOrderListItem is the output DTO for customer order listing.
type CustomerOrderListItem struct {
	OrderID      uuid.UUID                  `json:"order_id"`
	Items        []CustomerOrderProductItem `json:"items"`
	TotalPayment float64                    `json:"total_payment"`
	OrderDate    time.Time                  `json:"order_date"`
	Status       string                     `json:"status"`
}

// OrderActionUseCase defines customer order lifecycle actions.
type OrderActionUseCase interface {
	// CompleteByCustomer marks an eligible customer order as completed.
	CompleteByCustomer(ctx context.Context, userID, orderID uuid.UUID) (*CompleteOrderResponse, error)
	// ListByCustomer returns customer order rows with pagination metadata.
	ListByCustomer(ctx context.Context, userID uuid.UUID, params CustomerOrderListParams) ([]CustomerOrderListItem, *PaginationMeta, error)
	// ListOrderStatuses returns available order statuses.
	ListOrderStatuses(ctx context.Context) []string
}
