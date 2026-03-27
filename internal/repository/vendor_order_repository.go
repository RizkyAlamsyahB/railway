package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

// --- VendorOrderRepository ---

type vendorOrderRepository struct {
	db *gorm.DB
}

// NewVendorOrderRepository creates a new VendorOrderRepository backed by GORM.
func NewVendorOrderRepository(db *gorm.DB) domain.VendorOrderRepository {
	return &vendorOrderRepository{db: db}
}

func (r *vendorOrderRepository) ListByVendor(ctx context.Context, vendorID uuid.UUID, params domain.VendorOrderListParams) ([]domain.Order, int64, error) {
	q := r.db.WithContext(ctx).Model(&orderModel{}).Where("orders.vendor_id = ?", vendorID.String())

	// Status filter
	if params.Status != "" {
		q = q.Where("orders.order_status = ?", params.Status)
	}

	// Date filter
	if params.DateFrom != nil {
		q = q.Where("orders.placed_at >= ?", *params.DateFrom)
	}
	if params.DateTo != nil {
		// Include the entire day by adding 1 day
		endOfDay := params.DateTo.AddDate(0, 0, 1)
		q = q.Where("orders.placed_at < ?", endOfDay)
	}

	// Search filter (by order_no or customer name from users table)
	if params.Search != "" {
		searchPattern := "%" + strings.ToLower(params.Search) + "%"
		q = q.Joins("LEFT JOIN users ON users.id = orders.user_id").
			Where("LOWER(orders.order_no) LIKE ? OR LOWER(users.full_name) LIKE ?", searchPattern, searchPattern)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sorting
	orderClause := r.buildOrderClause(params.SortBy, params.SortOrder)
	q = q.Order(orderClause)

	offset := (params.Page - 1) * params.Limit
	var models []orderModel
	if err := q.
		Offset(offset).
		Limit(params.Limit).
		Find(&models).Error; err != nil {
		return nil, 0, err
	}

	orders := make([]domain.Order, len(models))
	for i, m := range models {
		orders[i] = *toDomainOrder(&m)
	}

	return orders, total, nil
}

// buildOrderClause builds the ORDER BY clause based on sort parameters.
func (r *vendorOrderRepository) buildOrderClause(sortBy, sortOrder string) string {
	// Default sort
	column := "orders.placed_at"
	direction := "DESC"

	// Map sortBy to actual column
	switch sortBy {
	case domain.VendorOrderSortByOrderDate:
		column = "orders.placed_at"
	case domain.VendorOrderSortByTotalPayment:
		column = "orders.grand_total"
	case domain.VendorOrderSortByCustomerName:
		// This requires join with users table which is done in search
		// For now, we sort by placed_at as fallback for customer_name sorting
		// A proper implementation would require subquery or join
		column = "orders.placed_at"
	}

	// Validate direction
	if strings.ToLower(sortOrder) == domain.VendorOrderSortOrderAsc {
		direction = "ASC"
	}

	return fmt.Sprintf("%s %s", column, direction)
}

func (r *vendorOrderRepository) FindByIDAndVendor(ctx context.Context, orderID, vendorID uuid.UUID) (*domain.Order, error) {
	var m orderModel
	if err := r.db.WithContext(ctx).
		Where("id = ? AND vendor_id = ?", orderID.String(), vendorID.String()).
		First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainOrder(&m), nil
}
