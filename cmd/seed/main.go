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

// seedUser holds the data needed to seed a single user.
type seedUser struct {
	Email    string
	Password string
	FullName string
	Phone    string
	RoleCode string
}

// devPassword is the default password used for all non-admin dev seed users.
const devPassword = "Password123!"

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
		_ = sqlDB.Close()
	}()

	log.Println("database connected successfully")
	logDatabaseContext(db)

	users := []seedUser{
		{
			Email:    cfg.Admin.Email,
			Password: cfg.Admin.Password,
			FullName: cfg.Admin.Name,
			Phone:    cfg.Admin.Phone,
			RoleCode: "admin",
		},
		{
			Email:    "umkm@dev.local",
			Password: devPassword,
			FullName: "Dev UMKM",
			Phone:    "081200000001",
			RoleCode: "umkm",
		},
		{
			Email:    "customer@dev.local",
			Password: devPassword,
			FullName: "Dev Customer",
			Phone:    "081200000002",
			RoleCode: "customer",
		},
		{
			Email:    "cs@dev.local",
			Password: devPassword,
			FullName: "Dev Customer Service",
			Phone:    "081200000003",
			RoleCode: "cs",
		},
		{
			Email:    "finance@dev.local",
			Password: devPassword,
			FullName: "Dev Finance",
			Phone:    "081200000004",
			RoleCode: "finance",
		},
	}

	for _, u := range users {
		if err := seedUserRecord(db, u); err != nil {
			log.Fatalf("failed to seed user [%s]: %v", u.Email, err)
		}
	}

	log.Println("all users seeded successfully")
	log.Printf("dev credentials (non-admin): password = %s", devPassword)

	// Find the CS user to use as creator for reply templates
	var csUser User
	if err := db.Where("lower(email) = lower(?)", "cs@dev.local").First(&csUser).Error; err != nil {
		log.Fatalf("failed to find CS user for reply template seeder: %v", err)
	}

	if err := SeedReplyTemplates(db, csUser.ID); err != nil {
		log.Fatalf("failed to seed reply templates: %v", err)
	}
}

func seedUserRecord(db *gorm.DB, u seedUser) error {
	var existing User
	result := db.Where("lower(email) = lower(?)", u.Email).First(&existing)
	if result.Error == nil {
		log.Printf("  skip   [%s] %s — already exists (id: %s)", u.RoleCode, u.Email, existing.ID)
		return nil
	}
	if result.Error != gorm.ErrRecordNotFound {
		if isUndefinedTableError(result.Error) {
			return fmt.Errorf("tables missing — run `make migrate-up` first: %w", result.Error)
		}
		return fmt.Errorf("failed to check existing user: %w", result.Error)
	}

	var role Role
	if err := db.Where("code = ?", u.RoleCode).First(&role).Error; err != nil {
		if isUndefinedTableError(err) {
			return fmt.Errorf("tables missing — run `make migrate-up` first: %w", err)
		}
		return fmt.Errorf("role %q not found: %w", u.RoleCode, err)
	}

	hash, err := auth.HashPassword(u.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	now := time.Now()
	userID := uuid.New().String()

	var phonePtr *string
	if u.Phone != "" {
		phonePtr = &u.Phone
	}

	newUser := User{
		ID:              userID,
		Email:           u.Email,
		FullName:        u.FullName,
		Phone:           phonePtr,
		PasswordHash:    hash,
		RoleID:          role.ID,
		Status:          "active",
		EmailVerifiedAt: &now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := db.Create(&newUser).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	log.Printf("  created [%s] %s (id: %s)", u.RoleCode, u.Email, userID)
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
	if err := db.Raw(`SELECT current_database() AS database, current_schema() AS schema`).Scan(&ctx).Error; err != nil {
		log.Printf("warning: failed to read database context: %v", err)
		return
	}

	log.Printf("database context: database=%s schema=%s", ctx.Database, ctx.Schema)
}
