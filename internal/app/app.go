package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/config"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/handler"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/router"
	ws "github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/websocket"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/infrastructure/database"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/infrastructure/email"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/infrastructure/payment"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/infrastructure/shipping"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/infrastructure/storage"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/repository"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/utils/sensitivedata"
	"gorm.io/gorm"
)

// App holds all initialized application components.
type App struct {
	Config      *config.Config
	DB          *gorm.DB
	Storage     domain.StorageProvider
	Email       domain.EmailProvider
	XenPlatform domain.XenPlatformProvider
	Router      *gin.Engine
	Scheduler   *usecase.OrderSettlementScheduler
	IMAPWorker  *usecase.IMAPWorker // nil when IMAP is disabled
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

	// Initialize Xendit XenPlatform provider
	xenPlatformProvider, err := newXenPlatformProvider(cfg.Xendit)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize xendit provider: %w", err)
	}

	log.Println("xendit xenplatform provider initialized")

	// Initialize Xendit Invoice provider
	xenditInvoiceProvider, err := newXenditInvoiceProvider(cfg.Xendit)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize xendit invoice provider: %w", err)
	}

	log.Println("xendit invoice provider initialized")

	// Initialize Xendit Payout provider
	xenditPayoutProvider, err := newXenditPayoutProvider(cfg.Xendit)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize xendit payout provider: %w", err)
	}

	log.Println("xendit payout provider initialized")

	fieldCipher, err := newSensitiveDataCipher(cfg.Sensitive)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize sensitive data protection: %w", err)
	}

	// Wire dependencies
	healthUseCase := usecase.NewHealthUseCase()
	healthHandler := handler.NewHealthHandler(healthUseCase)

	userRepo := repository.NewUserRepository(db)
	adminUserUseCase := usecase.NewAdminUserUseCase(userRepo)
	adminUserHandler := handler.NewAdminUserHandler(adminUserUseCase)

	adminAuthUseCase := usecase.NewAdminAuthUseCase(userRepo, cfg.JWT.Secret, cfg.JWT.ExpiryHours, cfg.JWT.Issuer)
	adminLoginHandler := handler.NewAdminLoginHandler(adminAuthUseCase)

	otpRepo := repository.NewOTPRepository(db)
	otpSigningSecret := cfg.OTP.Secret
	if otpSigningSecret == "" {
		otpSigningSecret = cfg.JWT.Secret
	}
	otpUseCase, err := usecase.NewOTPUseCase(otpRepo, emailProvider, cfg.OTP, otpSigningSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize otp usecase: %w", err)
	}
	otpHandler := handler.NewOTPHandler(otpUseCase)

	vendorRepo := repository.NewVendorRepository(db, fieldCipher)
	vendorOnboardingRepo := repository.NewVendorOnboardingRepository(db, fieldCipher)
	vendorUseCase := usecase.NewVendorUseCase(
		otpUseCase,
		userRepo,
		vendorRepo,
		vendorOnboardingRepo,
		storageProvider,
		xenditPayoutProvider,
		cfg.JWT.Secret,
		cfg.JWT.ExpiryHours,
		cfg.JWT.Issuer,
		usecase.VendorWithdrawalPolicy{
			FeeEstimateFixed: cfg.Withdrawal.FeeEstimateFixed,
			MinNetAmount:     cfg.Withdrawal.MinNetAmount,
		},
	)
	vendorHandler := handler.NewVendorHandler(vendorUseCase)

	adminVendorUseCase := usecase.NewAdminVendorUseCase(vendorRepo, userRepo, storageProvider, xenPlatformProvider)
	adminVendorHandler := handler.NewAdminVendorHandler(adminVendorUseCase)

	// Product & catalog feature
	categoryRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)

	productUseCase := usecase.NewProductUseCase(productRepo, vendorRepo, categoryRepo, storageProvider)
	productHandler := handler.NewProductHandler(productUseCase)

	catalogUseCase := usecase.NewCatalogUseCase(categoryRepo, productRepo, storageProvider)
	catalogHandler := handler.NewCatalogHandler(catalogUseCase)

	// User registration & email verification
	emailVerifRepo := repository.NewEmailVerificationTokenRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo, emailVerifRepo, emailProvider, cfg.App.BaseURL, cfg.JWT.Secret, cfg.JWT.ExpiryHours, cfg.JWT.Issuer)
	userHandler := handler.NewUserHandler(userUseCase, cfg.App.FrontendURL)

	// Cart feature
	cartRepo := repository.NewCartRepository(db)
	cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo, storageProvider)
	cartHandler := handler.NewCartHandler(cartUseCase)

	// Wishlist feature
	wishlistRepo := repository.NewWishlistRepository(db)
	wishlistUseCase := usecase.NewWishlistUseCase(wishlistRepo, productRepo, storageProvider)
	wishlistHandler := handler.NewWishlistHandler(wishlistUseCase)

	// Finance feature
	financeRepo := repository.NewFinanceRepository(db)
	financeUseCase := usecase.NewFinanceUseCase(financeRepo)
	financeHandler := handler.NewFinanceHandler(financeUseCase)

	// Notifications feature
	notificationRepo := repository.NewNotificationRepository(db)
	notificationUseCase := usecase.NewNotificationUseCase(notificationRepo)
	notificationHandler := handler.NewNotificationHandler(notificationUseCase)

	// CS auth feature
	csAuthUseCase := usecase.NewCSAuthUseCase(userRepo, cfg.JWT.Secret, cfg.JWT.ExpiryHours, cfg.JWT.Issuer)
	csLoginHandler := handler.NewCSLoginHandler(csAuthUseCase)

	// CS module
	ticketRepo := repository.NewTicketRepository(db)
	chatRepo := repository.NewChatRepository(db)
	replyTemplateRepo := repository.NewReplyTemplateRepository(db)
	subjectRepo := repository.NewTicketSubjectRepository(db)

	// FAQ / Pusat Bantuan
	faqRepo := repository.NewFAQRepository(db)
	faqUseCase := usecase.NewFAQUseCase(faqRepo)
	faqHandler := handler.NewFAQHandler(faqUseCase)

	// Banner management
	bannerRepo := repository.NewBannerRepository(db)
	bannerUseCase := usecase.NewBannerUseCase(bannerRepo, storageProvider)
	bannerHandler := handler.NewBannerHandler(bannerUseCase)

	// Vendor Banner management
	vendorBannerRepo := repository.NewVendorBannerRepository(db)
	vendorBannerUseCase := usecase.NewVendorBannerUseCase(vendorBannerRepo, storageProvider)
	vendorBannerHandler := handler.NewVendorBannerHandler(vendorBannerUseCase)

	// Vendor Voucher management
	vendorVoucherRepo := repository.NewVendorVoucherRepository(db)
	vendorVoucherUseCase := usecase.NewVendorVoucherUseCase(vendorVoucherRepo, productRepo)
	vendorVoucherHandler := handler.NewVendorVoucherHandler(vendorVoucherUseCase)

	// Address management (shared by customer + vendor)
	addressRepo := repository.NewAddressRepository(db)
	addressUseCase := usecase.NewAddressUseCase(addressRepo)
	addressHandler := handler.NewAddressHandler(addressUseCase)

	// Shipping location lookup (RajaOngkir)
	rajaOngkirProvider, err := newRajaOngkirProvider(cfg.RajaOngkir)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize rajaongkir provider: %w", err)
	}
	log.Println("rajaongkir provider initialized")

	shippingUseCase := usecase.NewShippingUseCase(rajaOngkirProvider)
	shippingHandler := handler.NewShippingHandler(shippingUseCase)

	// Courier management (vendor selects supported couriers)
	courierRepo := repository.NewCourierRepository(db)
	vendorCourierRepo := repository.NewVendorCourierRepository(db)
	vendorCourierUseCase := usecase.NewVendorCourierUseCase(courierRepo, vendorCourierRepo)
	vendorCourierHandler := handler.NewVendorCourierHandler(vendorCourierUseCase)

	// Master Lookup: Admin Category CRUD
	adminCategoryRepo := repository.NewAdminCategoryRepository(db)
	adminCategoryUseCase := usecase.NewAdminCategoryUseCase(adminCategoryRepo)
	adminCategoryHandler := handler.NewAdminCategoryHandler(adminCategoryUseCase)

	// Master Lookup: Return Reasons
	returnReasonRepo := repository.NewReturnReasonRepository(db)
	returnReasonUseCase := usecase.NewReturnReasonUseCase(returnReasonRepo)
	returnReasonHandler := handler.NewReturnReasonHandler(returnReasonUseCase)

	// Master Lookup: Admin Contacts
	adminContactRepo := repository.NewAdminContactRepository(db)
	adminContactUseCase := usecase.NewAdminContactUseCase(adminContactRepo)
	adminContactHandler := handler.NewAdminContactHandler(adminContactUseCase)

	// Master Lookup: Admin FAQ CRUD
	adminFAQUseCase := usecase.NewAdminFAQUseCase(faqRepo)
	adminFAQHandler := handler.NewAdminFAQHandler(adminFAQUseCase)

	orderRepo := repository.NewOrderRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

	// Admin payment listing
	adminPaymentUseCase := usecase.NewAdminPaymentUseCase(paymentRepo)
	adminPaymentHandler := handler.NewAdminPaymentHandler(adminPaymentUseCase)
	ticketUseCase := usecase.NewTicketUseCase(ticketRepo, subjectRepo, userRepo, orderRepo, paymentRepo, emailProvider, storageProvider, cfg.SMTP.FromEmail)
	csReportUseCase := usecase.NewCSReportUseCase(ticketRepo)
	ticketSubjectUseCase := usecase.NewTicketSubjectUseCase(subjectRepo)
	chatUseCase := usecase.NewChatUseCase(chatRepo, userRepo, storageProvider)
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
	csReportHandler := handler.NewCSReportHandler(csReportUseCase)
	ticketSubjectHandler := handler.NewTicketSubjectHandler(ticketSubjectUseCase)

	// Checkout & payment feature
	ledgerRepo := repository.NewLedgerRepository(db)
	shipmentRepo := repository.NewShipmentRepository(db)
	checkoutUseCase := usecase.NewCheckoutUseCase(cartRepo, productRepo, vendorRepo, orderRepo, paymentRepo, ledgerRepo, userRepo, addressRepo, vendorCourierRepo, shipmentRepo, rajaOngkirProvider, storageProvider, xenditInvoiceProvider, cfg.App.FrontendURL, cfg.Xendit.WebhookURL)
	checkoutHandler := handler.NewCheckoutHandler(checkoutUseCase, cfg.Xendit.WebhookVerificationToken, cfg.Xendit.Bypass)
	orderActionUseCase := usecase.NewOrderActionUseCase(orderRepo, shipmentRepo, paymentRepo, vendorRepo, productRepo, rajaOngkirProvider)
	orderActionHandler := handler.NewOrderActionHandler(orderActionUseCase)
	orderSettlementScheduler := usecase.NewOrderSettlementScheduler(orderRepo)

	// Vendor Order management
	vendorOrderRepo := repository.NewVendorOrderRepository(db)
	vendorOrderUseCase := usecase.NewVendorOrderUseCase(vendorOrderRepo, orderRepo, shipmentRepo, paymentRepo, userRepo, vendorRepo, rajaOngkirProvider)
	vendorOrderHandler := handler.NewVendorOrderHandler(vendorOrderUseCase)

	reviewRepo := repository.NewReviewRepository(db)
	reviewUseCase := usecase.NewReviewUseCase(reviewRepo, orderRepo, productRepo, storageProvider)
	reviewHandler := handler.NewReviewHandler(reviewUseCase)
	xenditWebhookHandler := handler.NewXenditWebhookHandler(vendorUseCase, cfg.Xendit.WebhookVerificationToken)

	// IMAP inbox poller (optional, enabled via IMAP_ENABLED=true)
	var imapWorker *usecase.IMAPWorker
	if cfg.IMAP.Enabled {
		imapReader, err := email.NewIMAPReader(cfg.IMAP)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize IMAP reader: %w", err)
		}
		pollInterval := time.Duration(cfg.IMAP.PollIntervalSec) * time.Second
		imapWorker = usecase.NewIMAPWorker(imapReader, ticketRepo, userRepo, pollInterval)
		go imapWorker.Start(context.Background())
		log.Printf("IMAP inbox poller started (interval=%ds)", cfg.IMAP.PollIntervalSec)
	} else {
		log.Println("IMAP inbox poller disabled (IMAP_ENABLED=false)")
	}

	// Setup router
	// r := router.NewRouter(healthHandler, adminUserHandler, adminVendorHandler, adminLoginHandler, vendorHandler, productHandler, catalogHandler, userHandler, cartHandler, checkoutHandler, cfg.JWT.Secret)
	r := router.NewRouter(
		healthHandler,
		adminUserHandler,
		adminVendorHandler,
		adminLoginHandler,
		vendorHandler,
		productHandler,
		catalogHandler,
		otpHandler,
		userHandler,
		cartHandler,
		wishlistHandler,
		checkoutHandler,
		orderActionHandler,
		reviewHandler,
		xenditWebhookHandler,
		financeHandler,
		notificationHandler,
		csLoginHandler,
		ticketHandler,
		chatHandler,
		replyTemplateHandler,
		csDashboardHandler,
		csUserHandler,
		csReportHandler,
		ticketSubjectHandler,
		faqHandler,
		bannerHandler,
		adminCategoryHandler,
		returnReasonHandler,
		adminContactHandler,
		adminFAQHandler,
		vendorBannerHandler,
		vendorVoucherHandler,
		addressHandler,
		shippingHandler,
		vendorCourierHandler,
		vendorOrderHandler,
		adminPaymentHandler,
		wsHandler,
		cfg.JWT.Secret,
		cfg.App.CORSAllowedOrigins,
	)

	return &App{
		Config:      cfg,
		DB:          db,
		Storage:     storageProvider,
		Email:       emailProvider,
		XenPlatform: xenPlatformProvider,
		Router:      r,
		Scheduler:   orderSettlementScheduler,
		IMAPWorker:  imapWorker,
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

// newXenPlatformProvider creates the Xendit-backed XenPlatformProvider from config.
// When cfg.Bypass is true, a no-op provider is returned for local development.
func newXenPlatformProvider(cfg config.XenditConfig) (domain.XenPlatformProvider, error) {
	if cfg.Bypass {
		return payment.NewNoopXenPlatformClient(), nil
	}
	return payment.NewXenditClient(cfg)
}

// newXenditInvoiceProvider creates the Xendit-backed XenditInvoiceProvider from config.
// When cfg.Bypass is true, a no-op provider is returned for local development.
func newXenditInvoiceProvider(cfg config.XenditConfig) (domain.XenditInvoiceProvider, error) {
	if cfg.Bypass {
		return payment.NewNoopXenditInvoiceClient(), nil
	}
	return payment.NewXenditInvoiceClient(cfg)
}

// newXenditPayoutProvider creates the Xendit-backed XenditPayoutProvider from config.
func newXenditPayoutProvider(cfg config.XenditConfig) (domain.XenditPayoutProvider, error) {
	return payment.NewXenditPayoutClient(cfg)
}

func newSensitiveDataCipher(cfg config.SensitiveDataConfig) (*sensitivedata.FieldCipher, error) {
	if cfg.EncryptionKey == "" {
		return nil, fmt.Errorf("SENSITIVE_DATA_ENCRYPTION_KEY is required")
	}
	return sensitivedata.NewFieldCipher([]byte(cfg.EncryptionKey))
}

// newRajaOngkirProvider creates the RajaOngkir API client from config.
func newRajaOngkirProvider(cfg config.RajaOngkirConfig) (domain.RajaOngkirProvider, error) {
	return shipping.NewRajaOngkirClient(cfg)
}

// Close cleans up application resources.
func (a *App) Close() error {
	if a.IMAPWorker != nil {
		a.IMAPWorker.Stop()
	}

	sqlDB, err := a.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	return sqlDB.Close()
}
