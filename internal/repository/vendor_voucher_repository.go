package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

type vendorVoucherModel struct {
	ID             string    `gorm:"column:id;primaryKey"`
	VendorID       string    `gorm:"column:vendor_id"`
	Code           string    `gorm:"column:code"`
	Name           string    `gorm:"column:name"`
	Description    *string   `gorm:"column:description"`
	DiscountAmount float64   `gorm:"column:discount_amount"`
	QuotaTotal     int       `gorm:"column:quota_total"`
	QuotaUsed      int       `gorm:"column:quota_used"`
	StartsAt       time.Time `gorm:"column:starts_at"`
	EndsAt         time.Time `gorm:"column:ends_at"`
	IsActive       bool      `gorm:"column:is_active"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (vendorVoucherModel) TableName() string { return "vendor_vouchers" }

type vendorVoucherProductModel struct {
	VoucherID string    `gorm:"column:voucher_id"`
	ProductID string    `gorm:"column:product_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (vendorVoucherProductModel) TableName() string { return "vendor_voucher_products" }

type vendorVoucherRepository struct {
	db *gorm.DB
}

// NewVendorVoucherRepository creates a new VendorVoucherRepository backed by GORM.
func NewVendorVoucherRepository(db *gorm.DB) domain.VendorVoucherRepository {
	return &vendorVoucherRepository{db: db}
}

func (r *vendorVoucherRepository) Create(ctx context.Context, voucher *domain.VendorVoucher) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		model := toVendorVoucherModel(voucher)
		if err := tx.Create(&model).Error; err != nil {
			return err
		}

		relations := make([]vendorVoucherProductModel, 0, len(voucher.ProductIDs))
		for _, productID := range voucher.ProductIDs {
			relations = append(relations, vendorVoucherProductModel{
				VoucherID: voucher.ID.String(),
				ProductID: productID.String(),
				CreatedAt: voucher.CreatedAt,
			})
		}

		if len(relations) > 0 {
			if err := tx.Create(&relations).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *vendorVoucherRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.VendorVoucher, error) {
	var model vendorVoucherModel
	if err := r.db.WithContext(ctx).
		Where("id = ?", id.String()).
		First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	productIDs, err := r.findProductIDsByVoucherIDs(ctx, []uuid.UUID{id})
	if err != nil {
		return nil, err
	}

	voucher := toDomainVendorVoucher(model)
	voucher.ProductIDs = productIDs[id]
	return &voucher, nil
}

func (r *vendorVoucherRepository) FindByCode(ctx context.Context, vendorID uuid.UUID, code string) (*domain.VendorVoucher, error) {
	var model vendorVoucherModel
	if err := r.db.WithContext(ctx).
		Where("vendor_id = ?", vendorID.String()).
		Where("lower(code) = lower(?)", strings.TrimSpace(code)).
		First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	id, err := uuid.Parse(model.ID)
	if err != nil {
		return nil, err
	}
	productIDs, err := r.findProductIDsByVoucherIDs(ctx, []uuid.UUID{id})
	if err != nil {
		return nil, err
	}

	voucher := toDomainVendorVoucher(model)
	voucher.ProductIDs = productIDs[id]
	return &voucher, nil
}

func (r *vendorVoucherRepository) ListByVendor(ctx context.Context, vendorID uuid.UUID, params domain.VendorVoucherListParams) ([]domain.VendorVoucher, int64, error) {
	baseQuery := r.db.WithContext(ctx).
		Table("vendor_vouchers vv").
		Where("vv.vendor_id = ?", vendorID.String())

	if params.Search != "" {
		kw := "%" + strings.ToLower(strings.TrimSpace(params.Search)) + "%"
		baseQuery = baseQuery.Where("(LOWER(vv.code) LIKE ? OR LOWER(vv.name) LIKE ?)", kw, kw)
	}
	if params.IsActive != nil {
		baseQuery = baseQuery.Where("vv.is_active = ?", *params.IsActive)
	}
	if params.ProductID != nil {
		baseQuery = baseQuery.Joins("JOIN vendor_voucher_products vvp_filter ON vvp_filter.voucher_id = vv.id").
			Where("vvp_filter.product_id = ?", params.ProductID.String())
	}

	var total int64
	if err := baseQuery.Distinct("vv.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []vendorVoucherModel
	if err := baseQuery.
		Select("vv.*").
		Distinct("vv.id, vv.vendor_id, vv.code, vv.name, vv.description, vv.discount_amount, vv.quota_total, vv.quota_used, vv.starts_at, vv.ends_at, vv.is_active, vv.created_at, vv.updated_at").
		Order("vv.created_at DESC").
		Offset((params.Page - 1) * params.Limit).
		Limit(params.Limit).
		Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	voucherIDs := make([]uuid.UUID, 0, len(rows))
	for _, row := range rows {
		id, err := uuid.Parse(row.ID)
		if err != nil {
			return nil, 0, err
		}
		voucherIDs = append(voucherIDs, id)
	}

	productMap, err := r.findProductIDsByVoucherIDs(ctx, voucherIDs)
	if err != nil {
		return nil, 0, err
	}

	items := make([]domain.VendorVoucher, 0, len(rows))
	for _, row := range rows {
		item := toDomainVendorVoucher(row)
		item.ProductIDs = productMap[item.ID]
		items = append(items, item)
	}

	return items, total, nil
}

func (r *vendorVoucherRepository) Update(ctx context.Context, voucher *domain.VendorVoucher, replaceProducts bool) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vendorVoucherModel{}).
			Where("id = ?", voucher.ID.String()).
			Updates(map[string]interface{}{
				"code":            voucher.Code,
				"name":            voucher.Name,
				"description":     voucher.Description,
				"discount_amount": voucher.DiscountAmount,
				"quota_total":     voucher.QuotaTotal,
				"quota_used":      voucher.QuotaUsed,
				"starts_at":       voucher.StartsAt,
				"ends_at":         voucher.EndsAt,
				"is_active":       voucher.IsActive,
				"updated_at":      voucher.UpdatedAt,
			}).Error; err != nil {
			return err
		}

		if !replaceProducts {
			return nil
		}

		if err := tx.Where("voucher_id = ?", voucher.ID.String()).
			Delete(&vendorVoucherProductModel{}).Error; err != nil {
			return err
		}

		relations := make([]vendorVoucherProductModel, 0, len(voucher.ProductIDs))
		for _, productID := range voucher.ProductIDs {
			relations = append(relations, vendorVoucherProductModel{
				VoucherID: voucher.ID.String(),
				ProductID: productID.String(),
				CreatedAt: voucher.UpdatedAt,
			})
		}

		if len(relations) == 0 {
			return nil
		}

		return tx.Create(&relations).Error
	})
}

func (r *vendorVoucherRepository) Delete(ctx context.Context, id uuid.UUID, vendorID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND vendor_id = ?", id.String(), vendorID.String()).
		Delete(&vendorVoucherModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *vendorVoucherRepository) findProductIDsByVoucherIDs(ctx context.Context, voucherIDs []uuid.UUID) (map[uuid.UUID][]uuid.UUID, error) {
	result := make(map[uuid.UUID][]uuid.UUID, len(voucherIDs))
	if len(voucherIDs) == 0 {
		return result, nil
	}

	stringIDs := make([]string, 0, len(voucherIDs))
	for _, id := range voucherIDs {
		stringIDs = append(stringIDs, id.String())
		result[id] = []uuid.UUID{}
	}

	var rows []vendorVoucherProductModel
	if err := r.db.WithContext(ctx).
		Where("voucher_id IN ?", stringIDs).
		Order("created_at ASC").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		voucherID, err := uuid.Parse(row.VoucherID)
		if err != nil {
			return nil, err
		}
		productID, err := uuid.Parse(row.ProductID)
		if err != nil {
			return nil, err
		}
		result[voucherID] = append(result[voucherID], productID)
	}

	return result, nil
}

func toVendorVoucherModel(v *domain.VendorVoucher) vendorVoucherModel {
	return vendorVoucherModel{
		ID:             v.ID.String(),
		VendorID:       v.VendorID.String(),
		Code:           v.Code,
		Name:           v.Name,
		Description:    v.Description,
		DiscountAmount: v.DiscountAmount,
		QuotaTotal:     v.QuotaTotal,
		QuotaUsed:      v.QuotaUsed,
		StartsAt:       v.StartsAt,
		EndsAt:         v.EndsAt,
		IsActive:       v.IsActive,
		CreatedAt:      v.CreatedAt,
		UpdatedAt:      v.UpdatedAt,
	}
}

func toDomainVendorVoucher(m vendorVoucherModel) domain.VendorVoucher {
	id, _ := uuid.Parse(m.ID)
	vendorID, _ := uuid.Parse(m.VendorID)
	return domain.VendorVoucher{
		ID:             id,
		VendorID:       vendorID,
		Code:           m.Code,
		Name:           m.Name,
		Description:    m.Description,
		DiscountAmount: m.DiscountAmount,
		QuotaTotal:     m.QuotaTotal,
		QuotaUsed:      m.QuotaUsed,
		StartsAt:       m.StartsAt,
		EndsAt:         m.EndsAt,
		IsActive:       m.IsActive,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}
