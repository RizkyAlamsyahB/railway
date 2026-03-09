package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

// addressModel is the GORM model for the addresses table.
type addressModel struct {
	ID              uuid.UUID `gorm:"column:id;type:uuid;primaryKey"`
	UserID          uuid.UUID `gorm:"column:user_id;type:uuid;not null"`
	Label           *string   `gorm:"column:label"`
	RecipientName   *string   `gorm:"column:recipient_name"`
	Phone           *string   `gorm:"column:phone"`
	ProvinceID      *string   `gorm:"column:province_id"`
	ProvinceName    *string   `gorm:"column:province_name"`
	CityID          *string   `gorm:"column:city_id"`
	CityName        *string   `gorm:"column:city_name"`
	DistrictID      *string   `gorm:"column:district_id"`
	DistrictName    *string   `gorm:"column:district_name"`
	SubdistrictID   *string   `gorm:"column:subdistrict_id"`
	SubdistrictName *string   `gorm:"column:subdistrict_name"`
	PostalCode      *string   `gorm:"column:postal_code"`
	AddressLine     string    `gorm:"column:address_line;not null"`
	Notes           *string   `gorm:"column:notes"`
	Latitude        *float64  `gorm:"column:latitude"`
	Longitude       *float64  `gorm:"column:longitude"`
	IsDefault       bool      `gorm:"column:is_default;not null;default:false"`
	CreatedAt       time.Time `gorm:"column:created_at;not null"`
	UpdatedAt       time.Time `gorm:"column:updated_at;not null"`
}

func (addressModel) TableName() string { return "addresses" }

// addressRepository implements domain.AddressRepository.
type addressRepository struct {
	db *gorm.DB
}

// NewAddressRepository creates a new AddressRepository backed by GORM.
func NewAddressRepository(db *gorm.DB) domain.AddressRepository {
	return &addressRepository{db: db}
}

func (r *addressRepository) Create(ctx context.Context, address *domain.Address) error {
	m := toAddressModel(address)
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *addressRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Address, error) {
	var m addressModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return toAddressDomain(m), nil
}

func (r *addressRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Address, error) {
	var rows []addressModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("is_default DESC, created_at DESC").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	addresses := make([]domain.Address, len(rows))
	for i, m := range rows {
		addresses[i] = *toAddressDomain(m)
	}
	return addresses, nil
}

func (r *addressRepository) FindDefaultByUserID(ctx context.Context, userID uuid.UUID) (*domain.Address, error) {
	var m addressModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_default = true", userID).
		First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return toAddressDomain(m), nil
}

func (r *addressRepository) Update(ctx context.Context, address *domain.Address) error {
	m := toAddressModel(address)
	return r.db.WithContext(ctx).Save(&m).Error
}

func (r *addressRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&addressModel{}).Error
}

func (r *addressRepository) ClearDefaultByUser(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&addressModel{}).
		Where("user_id = ? AND is_default = true", userID).
		Update("is_default", false).Error
}

// --- mappers ---

func toAddressModel(a *domain.Address) addressModel {
	return addressModel{
		ID:              a.ID,
		UserID:          a.UserID,
		Label:           a.Label,
		RecipientName:   a.RecipientName,
		Phone:           a.Phone,
		ProvinceID:      a.ProvinceID,
		ProvinceName:    a.ProvinceName,
		CityID:          a.CityID,
		CityName:        a.CityName,
		DistrictID:      a.DistrictID,
		DistrictName:    a.DistrictName,
		SubdistrictID:   a.SubdistrictID,
		SubdistrictName: a.SubdistrictName,
		PostalCode:      a.PostalCode,
		AddressLine:     a.AddressLine,
		Notes:           a.Notes,
		Latitude:        a.Latitude,
		Longitude:       a.Longitude,
		IsDefault:       a.IsDefault,
		CreatedAt:       a.CreatedAt,
		UpdatedAt:       a.UpdatedAt,
	}
}

func toAddressDomain(m addressModel) *domain.Address {
	return &domain.Address{
		ID:              m.ID,
		UserID:          m.UserID,
		Label:           m.Label,
		RecipientName:   m.RecipientName,
		Phone:           m.Phone,
		ProvinceID:      m.ProvinceID,
		ProvinceName:    m.ProvinceName,
		CityID:          m.CityID,
		CityName:        m.CityName,
		DistrictID:      m.DistrictID,
		DistrictName:    m.DistrictName,
		SubdistrictID:   m.SubdistrictID,
		SubdistrictName: m.SubdistrictName,
		PostalCode:      m.PostalCode,
		AddressLine:     m.AddressLine,
		Notes:           m.Notes,
		Latitude:        m.Latitude,
		Longitude:       m.Longitude,
		IsDefault:       m.IsDefault,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}
