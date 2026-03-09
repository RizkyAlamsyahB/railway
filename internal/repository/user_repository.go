package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

// GORM model structs (internal to repository layer).

type userModel struct {
	ID              string     `gorm:"column:id;primaryKey"`
	Email           string     `gorm:"column:email"`
	FullName        string     `gorm:"column:full_name"`
	ImageURL        *string    `gorm:"column:image_url"`
	BirthDate       *time.Time `gorm:"column:birth_date"`
	Phone           *string    `gorm:"column:phone"`
	PasswordHash    string     `gorm:"column:password_hash"`
	RoleID          int16      `gorm:"column:role_id"`
	Status          string     `gorm:"column:status"`
	EmailVerifiedAt *time.Time `gorm:"column:email_verified_at"`
	CreatedAt       time.Time  `gorm:"column:created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at"`
}

func (userModel) TableName() string { return "users" }

type roleModel struct {
	ID   int16  `gorm:"column:id;primaryKey"`
	Code string `gorm:"column:code"`
	Name string `gorm:"column:name"`
}

func (roleModel) TableName() string { return "roles" }

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository backed by GORM.
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User, roleCode string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var role roleModel
		if err := tx.Where("code = ?", roleCode).First(&role).Error; err != nil {
			return fmt.Errorf("role '%s' not found: %w", roleCode, err)
		}

		user.Role = &domain.Role{ID: role.ID, Code: role.Code, Name: role.Name}
		model := toUserModel(user)
		model.RoleID = role.ID
		if err := tx.Create(&model).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var model userModel
	if err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	role, err := r.getUserRole(ctx, model.RoleID)
	if err != nil {
		return nil, err
	}

	user := toDomainUser(&model)
	user.Role = role
	return user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var model userModel
	if err := r.db.WithContext(ctx).Where("lower(email) = lower(?)", email).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	role, err := r.getUserRole(ctx, model.RoleID)
	if err != nil {
		return nil, err
	}

	user := toDomainUser(&model)
	user.Role = role
	return user, nil
}

func (r *userRepository) FindByPhone(ctx context.Context, phone string) (*domain.User, error) {
	var model userModel
	if err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	role, err := r.getUserRole(ctx, model.RoleID)
	if err != nil {
		return nil, err
	}

	user := toDomainUser(&model)
	user.Role = role
	return user, nil
}

func (r *userRepository) ActivateUser(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&userModel{}).
		Where("id = ?", id.String()).
		Updates(map[string]interface{}{
			"status":            "active",
			"email_verified_at": now,
			"updated_at":        now,
		}).Error
}

func (r *userRepository) List(ctx context.Context, params domain.UserListParams) ([]domain.User, int64, error) {
	query := r.db.WithContext(ctx).Model(&userModel{})

	if params.Role != "" {
		query = query.Joins("JOIN roles ON roles.id = users.role_id").
			Where("roles.code = ?", params.Role)
	}

	if params.Status != "" {
		query = query.Where("users.status = ?", params.Status)
	}

	if params.Search != "" {
		search := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("(lower(users.full_name) LIKE ? OR lower(users.email) LIKE ?)", search, search)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	var models []userModel
	if err := query.Select("users.*").
		Order("users.created_at DESC").
		Offset(offset).Limit(params.Limit).
		Find(&models).Error; err != nil {
		return nil, 0, err
	}

	users := make([]domain.User, len(models))
	for i, m := range models {
		u := toDomainUser(&m)
		role, _ := r.getUserRole(ctx, m.RoleID)
		u.Role = role
		users[i] = *u
	}

	return users, total, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	model := toUserModel(user)
	return r.db.WithContext(ctx).Model(&model).
		Select("full_name", "phone", "birth_date", "status", "password_hash", "updated_at").
		Updates(&model).Error
}

func (r *userRepository) UpdateWithRole(ctx context.Context, user *domain.User, roleCode string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var role roleModel
		if err := tx.Where("code = ?", roleCode).First(&role).Error; err != nil {
			return fmt.Errorf("role '%s' not found: %w", roleCode, err)
		}

		model := toUserModel(user)
		model.RoleID = role.ID
		if err := tx.Model(&model).
			Select("full_name", "phone", "birth_date", "status", "password_hash", "role_id", "updated_at").
			Updates(&model).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("id = ?", id.String()).Delete(&userModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *userRepository) getUserRole(ctx context.Context, roleID int16) (*domain.Role, error) {
	var role roleModel
	err := r.db.WithContext(ctx).Where("id = ?", roleID).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &domain.Role{ID: role.ID, Code: role.Code, Name: role.Name}, nil
}

func toUserModel(u *domain.User) userModel {
	m := userModel{
		ID:              u.ID.String(),
		Email:           u.Email,
		FullName:        u.FullName,
		ImageURL:        u.ImageURL,
		BirthDate:       u.BirthDate,
		Phone:           u.Phone,
		PasswordHash:    u.PasswordHash,
		Status:          u.Status,
		EmailVerifiedAt: u.EmailVerifiedAt,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
	if u.Role != nil {
		m.RoleID = u.Role.ID
	}
	return m
}

func toDomainUser(m *userModel) *domain.User {
	id, _ := uuid.Parse(m.ID)
	return &domain.User{
		ID:              id,
		Email:           m.Email,
		FullName:        m.FullName,
		ImageURL:        m.ImageURL,
		BirthDate:       m.BirthDate,
		Phone:           m.Phone,
		PasswordHash:    m.PasswordHash,
		Status:          m.Status,
		EmailVerifiedAt: m.EmailVerifiedAt,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}
