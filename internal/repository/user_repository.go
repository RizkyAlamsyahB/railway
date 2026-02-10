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
	BirthDate       *time.Time `gorm:"column:birth_date"`
	Phone           *string    `gorm:"column:phone"`
	PasswordHash    string     `gorm:"column:password_hash"`
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

type userRoleModel struct {
	UserID string `gorm:"column:user_id;primaryKey"`
	RoleID int16  `gorm:"column:role_id;primaryKey"`
}

func (userRoleModel) TableName() string { return "user_roles" }

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

		model := toUserModel(user)
		if err := tx.Create(&model).Error; err != nil {
			return err
		}

		ur := userRoleModel{UserID: user.ID.String(), RoleID: role.ID}
		if err := tx.Create(&ur).Error; err != nil {
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

	roles, err := r.getUserRoles(ctx, model.ID)
	if err != nil {
		return nil, err
	}

	user := toDomainUser(&model)
	user.Roles = roles
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

	roles, err := r.getUserRoles(ctx, model.ID)
	if err != nil {
		return nil, err
	}

	user := toDomainUser(&model)
	user.Roles = roles
	return user, nil
}

func (r *userRepository) List(ctx context.Context, params domain.UserListParams) ([]domain.User, int64, error) {
	query := r.db.WithContext(ctx).Model(&userModel{})

	if params.Role != "" {
		query = query.Joins("JOIN user_roles ON user_roles.user_id = users.id").
			Joins("JOIN roles ON roles.id = user_roles.role_id").
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
		roles, _ := r.getUserRoles(ctx, m.ID)
		u.Roles = roles
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
		model := toUserModel(user)
		if err := tx.Model(&model).
			Select("full_name", "phone", "birth_date", "status", "password_hash", "updated_at").
			Updates(&model).Error; err != nil {
			return err
		}

		if err := tx.Where("user_id = ?", user.ID.String()).Delete(&userRoleModel{}).Error; err != nil {
			return err
		}

		var role roleModel
		if err := tx.Where("code = ?", roleCode).First(&role).Error; err != nil {
			return fmt.Errorf("role '%s' not found: %w", roleCode, err)
		}

		return tx.Create(&userRoleModel{UserID: user.ID.String(), RoleID: role.ID}).Error
	})
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", id.String()).Delete(&userRoleModel{}).Error; err != nil {
			return err
		}

		result := tx.Where("id = ?", id.String()).Delete(&userModel{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

func (r *userRepository) getUserRoles(ctx context.Context, userID string) ([]domain.Role, error) {
	var roles []roleModel
	err := r.db.WithContext(ctx).
		Table("roles").
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error
	if err != nil {
		return nil, err
	}
	result := make([]domain.Role, len(roles))
	for i, r := range roles {
		result[i] = domain.Role{ID: r.ID, Code: r.Code, Name: r.Name}
	}
	return result, nil
}

func toUserModel(u *domain.User) userModel {
	return userModel{
		ID:              u.ID.String(),
		Email:           u.Email,
		FullName:        u.FullName,
		BirthDate:       u.BirthDate,
		Phone:           u.Phone,
		PasswordHash:    u.PasswordHash,
		Status:          u.Status,
		EmailVerifiedAt: u.EmailVerifiedAt,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}

func toDomainUser(m *userModel) *domain.User {
	id, _ := uuid.Parse(m.ID)
	return &domain.User{
		ID:              id,
		Email:           m.Email,
		FullName:        m.FullName,
		BirthDate:       m.BirthDate,
		Phone:           m.Phone,
		PasswordHash:    m.PasswordHash,
		Status:          m.Status,
		EmailVerifiedAt: m.EmailVerifiedAt,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}
