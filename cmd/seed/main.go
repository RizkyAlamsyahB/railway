package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/config"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/infrastructure/database"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/utils/auth"
	"gorm.io/gorm"
)

// Minimal GORM models scoped to this seed script only.
type User struct {
	ID              string     `gorm:"column:id;primaryKey"`
	Email           string     `gorm:"column:email"`
	FullName        string     `gorm:"column:full_name"`
	Phone           *string    `gorm:"column:phone"`
	PasswordHash    string     `gorm:"column:password_hash"`
	RoleID          int16      `gorm:"column:role_id"`
	Status          string     `gorm:"column:status"`
	EmailVerifiedAt *time.Time `gorm:"column:email_verified_at"`
	CreatedAt       time.Time  `gorm:"column:created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at"`
}

func (User) TableName() string { return "users" }

type Role struct {
	ID   int16  `gorm:"column:id;primaryKey"`
	Code string `gorm:"column:code"`
}

func (Role) TableName() string { return "roles" }

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if cfg.Admin.Email == "" || cfg.Admin.Password == "" || cfg.Admin.Name == "" {
		log.Fatal("ADMIN_EMAIL, ADMIN_PASSWORD, and ADMIN_NAME must be set")
	}

	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	log.Println("database connected successfully")
	logDatabaseContext(db)

	if err := seedAdmin(db, cfg.Admin.Email, cfg.Admin.Password, cfg.Admin.Name, cfg.Admin.Phone); err != nil {
		log.Fatalf("failed to seed admin: %v", err)
	}

	log.Println("admin seeding completed successfully")
}

func seedAdmin(db *gorm.DB, email, password, name, phone string) error {
	var existingUser User
	result := db.Where("lower(email) = lower(?)", email).First(&existingUser)
	if result.Error == nil {
		log.Printf("user with email %s already exists (id: %s), skipping", email, existingUser.ID)
		return nil
	}
	if result.Error != gorm.ErrRecordNotFound {
		if isUndefinedTableError(result.Error) {
			return fmt.Errorf("required tables are missing; run `make migrate-up` (or `make seed`, which now runs migrations): %w", result.Error)
		}
		return fmt.Errorf("failed to check existing user: %w", result.Error)
	}

	var adminRole Role
	if err := db.Where("code = ?", "admin").First(&adminRole).Error; err != nil {
		if isUndefinedTableError(err) {
			return fmt.Errorf("required tables are missing; run `make migrate-up` (or `make seed`, which now runs migrations): %w", err)
		}
		return fmt.Errorf("failed to find admin role (ensure migrations have been run): %w", err)
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	now := time.Now()
	userID := uuid.New().String()

	var phonePtr *string
	if phone != "" {
		phonePtr = &phone
	}

	newUser := User{
		ID:              userID,
		Email:           email,
		FullName:        name,
		Phone:           phonePtr,
		PasswordHash:    hash,
		RoleID:          adminRole.ID,
		Status:          "active",
		EmailVerifiedAt: &now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := db.Create(&newUser).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	log.Printf("admin user created successfully: id=%s, email=%s", userID, email)
	return nil
}

func isUndefinedTableError(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "42P01"
}

func logDatabaseContext(db *gorm.DB) {
	type dbContext struct {
		Database string `gorm:"column:database"`
		Schema   string `gorm:"column:schema"`
	}

	var ctx dbContext
	err := db.Raw(`SELECT current_database() AS database, current_schema() AS schema`).Scan(&ctx).Error
	if err != nil {
		log.Printf("warning: failed to read database context: %v", err)
		return
	}

	log.Printf("database context: database=%s schema=%s", ctx.Database, ctx.Schema)
}
