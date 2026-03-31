package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

// --- Shipment GORM model ---

type shipmentModel struct {
	ID             string     `gorm:"column:id;primaryKey"`
	OrderID        string     `gorm:"column:order_id"`
	CourierCode    string     `gorm:"column:courier_code"`
	ServiceType    string     `gorm:"column:service_type"`
	TrackingNo     string     `gorm:"column:tracking_no"`
	ETD            string     `gorm:"column:etd"`
	ShipmentStatus string     `gorm:"column:shipment_status"`
	ShippedAt      *time.Time `gorm:"column:shipped_at"`
	DeliveredAt    *time.Time `gorm:"column:delivered_at"`
}

func (shipmentModel) TableName() string { return "shipments" }

// --- ShipmentRepository ---

type shipmentRepository struct {
	db *gorm.DB
}

// NewShipmentRepository creates a new ShipmentRepository backed by GORM.
func NewShipmentRepository(db *gorm.DB) domain.ShipmentRepository {
	return &shipmentRepository{db: db}
}

func (r *shipmentRepository) Create(ctx context.Context, s *domain.Shipment) error {
	m := shipmentModel{
		ID:             s.ID.String(),
		OrderID:        s.OrderID.String(),
		CourierCode:    s.CourierCode,
		ServiceType:    s.ServiceType,
		TrackingNo:     s.TrackingNo,
		ETD:            s.ETD,
		ShipmentStatus: s.ShipmentStatus,
		ShippedAt:      s.ShippedAt,
		DeliveredAt:    s.DeliveredAt,
	}
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *shipmentRepository) FindByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.Shipment, error) {
	var m shipmentModel
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID.String()).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainShipment(&m), nil
}

func (r *shipmentRepository) UpdateTrackingAndShip(ctx context.Context, orderID uuid.UUID, trackingNo string, shippedAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&shipmentModel{}).
		Where("order_id = ?", orderID.String()).
		Updates(map[string]interface{}{
			"tracking_no":     trackingNo,
			"shipment_status": "shipped",
			"shipped_at":      shippedAt,
		}).Error
}

func (r *shipmentRepository) UpdateDeliveredAt(ctx context.Context, orderID uuid.UUID, deliveredAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&shipmentModel{}).
		Where("order_id = ?", orderID.String()).
		Updates(map[string]interface{}{
			"shipment_status": "delivered",
			"delivered_at":    deliveredAt,
		}).Error
}

func toDomainShipment(m *shipmentModel) *domain.Shipment {
	id, _ := uuid.Parse(m.ID)
	orderID, _ := uuid.Parse(m.OrderID)
	return &domain.Shipment{
		ID:             id,
		OrderID:        orderID,
		CourierCode:    m.CourierCode,
		ServiceType:    m.ServiceType,
		TrackingNo:     m.TrackingNo,
		ETD:            m.ETD,
		ShipmentStatus: m.ShipmentStatus,
		ShippedAt:      m.ShippedAt,
		DeliveredAt:    m.DeliveredAt,
	}
}
