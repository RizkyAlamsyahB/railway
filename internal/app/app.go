package app

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/config"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/handler"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/router"
	ws "github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/websocket"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/infrastructure/database"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/infrastructure/email"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/infrastructure/storage"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/repository"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase"
	"gorm.io/gorm"
)

// App holds all initialized application components.
type App struct {
	Config  *config.Config
	DB      *gorm.DB
	Storage domain.StorageProvider
	Email   domain.EmailProvider
	Router  *gin.Engine
}

// Initialize loads config, connects to the database, wires all dependencies,
// and returns a fully configured App ready to serve requests.
func Initialize() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	log.Printf("starting %s in %s mode", cfg.App.Name, cfg.App.Env)

	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("database connected successfully")

	// Initialize storage provider
	storageProvider, err := newStorageProvider(cfg.Storage)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}

	log.Printf("storage provider initialized (provider=%s)", cfg.Storage.Provider)

	// Initialize email provider
	emailProvider, err := newEmailProvider(cfg.SMTP)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize email provider: %w", err)
	}

	log.Println("email provider initialized (SMTP)")

	// Wire dependencies
	healthUseCase := usecase.NewHealthUseCase()
	healthHandler := handler.NewHealthHandler(healthUseCase)

	userRepo := repository.NewUserRepository(db)
	adminUserUseCase := usecase.NewAdminUserUseCase(userRepo)
	adminUserHandler := handler.NewAdminUserHandler(adminUserUseCase)

	adminAuthUseCase := usecase.NewAdminAuthUseCase(userRepo, cfg.JWT.Secret, cfg.JWT.ExpiryHours, cfg.JWT.Issuer)
	adminLoginHandler := handler.NewAdminLoginHandler(adminAuthUseCase)

	vendorRepo := repository.NewVendorRepository(db)
	vendorUseCase := usecase.NewVendorUseCase(userRepo, vendorRepo, storageProvider, cfg.JWT.Secret, cfg.JWT.ExpiryHours, cfg.JWT.Issuer)
	vendorHandler := handler.NewVendorHandler(vendorUseCase)

	adminVendorUseCase := usecase.NewAdminVendorUseCase(vendorRepo, userRepo, storageProvider)
	adminVendorHandler := handler.NewAdminVendorHandler(adminVendorUseCase)

	// Product & catalog feature
	categoryRepo := repository.NewCategoryRepository(db)
	shippingServiceRepo := repository.NewShippingServiceRepository(db)
	productRepo := repository.NewProductRepository(db)

	productUseCase := usecase.NewProductUseCase(productRepo, vendorRepo, categoryRepo, shippingServiceRepo, storageProvider)
	productHandler := handler.NewProductHandler(productUseCase)

	catalogUseCase := usecase.NewCatalogUseCase(categoryRepo, shippingServiceRepo)
	catalogHandler := handler.NewCatalogHandler(catalogUseCase)

	// User registration & email verification
	emailVerifRepo := repository.NewEmailVerificationTokenRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo, emailVerifRepo, emailProvider, cfg.App.BaseURL, cfg.JWT.Secret, cfg.JWT.ExpiryHours, cfg.JWT.Issuer)
	userHandler := handler.NewUserHandler(userUseCase, cfg.App.FrontendURL)

	// Cart feature
	cartRepo := repository.NewCartRepository(db)
	cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo, storageProvider)
	cartHandler := handler.NewCartHandler(cartUseCase)

	// CS auth feature
	csAuthUseCase := usecase.NewCSAuthUseCase(userRepo, cfg.JWT.Secret, cfg.JWT.ExpiryHours, cfg.JWT.Issuer)
	csLoginHandler := handler.NewCSLoginHandler(csAuthUseCase)

	// CS module
	ticketRepo := repository.NewTicketRepository(db)
	chatRepo := repository.NewChatRepository(db)
	replyTemplateRepo := repository.NewReplyTemplateRepository(db)

	ticketUseCase := usecase.NewTicketUseCase(ticketRepo)
	chatUseCase := usecase.NewChatUseCase(chatRepo, userRepo)
	replyTemplateUseCase := usecase.NewReplyTemplateUseCase(replyTemplateRepo)
	csDashboardUseCase := usecase.NewCSDashboardUseCase(ticketRepo, chatRepo)
	csUserUseCase := usecase.NewCSUserUseCase(userRepo, ticketRepo)

	// WebSocket Hub (runs in background goroutine)
	hub := ws.NewHub()
	wsHandler := ws.NewHandler(hub, cfg.JWT.Secret)

	ticketHandler := handler.NewTicketHandler(ticketUseCase)
	chatHandler := handler.NewChatHandler(chatUseCase, hub)
	replyTemplateHandler := handler.NewReplyTemplateHandler(replyTemplateUseCase)
	csDashboardHandler := handler.NewCSDashboardHandler(csDashboardUseCase)
	csUserHandler := handler.NewCSUserHandler(csUserUseCase, ticketUseCase)

	// Setup router
	r := router.NewRouter(healthHandler, adminUserHandler, adminVendorHandler, adminLoginHandler, vendorHandler, productHandler, catalogHandler, userHandler, cartHandler, csLoginHandler, ticketHandler, chatHandler, replyTemplateHandler, csDashboardHandler, csUserHandler, wsHandler, cfg.JWT.Secret)

	return &App{
		Config:  cfg,
		DB:      db,
		Storage: storageProvider,
		Email:   emailProvider,
		Router:  r,
	}, nil
}

// newStorageProvider creates the appropriate StorageProvider based on config.
func newStorageProvider(cfg config.StorageConfig) (domain.StorageProvider, error) {
	switch cfg.Provider {
	case "s3":
		return storage.NewS3Storage(cfg)
	default:
		return nil, fmt.Errorf("unsupported storage provider: %q", cfg.Provider)
	}
}

// newEmailProvider creates the SMTP-backed EmailProvider from config.
func newEmailProvider(cfg config.SMTPConfig) (domain.EmailProvider, error) {
	return email.NewSMTPSender(cfg)
}

// Close cleans up application resources.
func (a *App) Close() error {
	sqlDB, err := a.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	return sqlDB.Close()
}
