package router

import (
	"github.com/gin-gonic/gin"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/handler"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/middleware"
	ws "github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/websocket"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

// NewRouter sets up the Gin engine with middleware and route registration.
func NewRouter(
	healthHandler *handler.HealthHandler,
	adminUserHandler *handler.AdminUserHandler,
	adminVendorHandler *handler.AdminVendorHandler,
	adminLoginHandler *handler.AdminLoginHandler,
	vendorHandler *handler.VendorHandler,
	productHandler *handler.ProductHandler,
	catalogHandler *handler.CatalogHandler,
	otpHandler *handler.OTPHandler,
	userHandler *handler.UserHandler,
	cartHandler *handler.CartHandler,
	wishlistHandler *handler.WishlistHandler,
	checkoutHandler *handler.CheckoutHandler,
	orderActionHandler *handler.OrderActionHandler,
	reviewHandler *handler.ReviewHandler,
	xenditWebhookHandler *handler.XenditWebhookHandler,
	financeHandler *handler.FinanceHandler,
	notificationHandler *handler.NotificationHandler,
	csLoginHandler *handler.CSLoginHandler,
	ticketHandler *handler.TicketHandler,
	chatHandler *handler.ChatHandler,
	replyTemplateHandler *handler.ReplyTemplateHandler,
	csDashboardHandler *handler.CSDashboardHandler,
	csUserHandler *handler.CSUserHandler,
	csReportHandler *handler.CSReportHandler,
	ticketSubjectHandler *handler.TicketSubjectHandler,
	faqHandler *handler.FAQHandler,
	bannerHandler *handler.BannerHandler,
	adminCategoryHandler *handler.AdminCategoryHandler,
	returnReasonHandler *handler.ReturnReasonHandler,
	adminContactHandler *handler.AdminContactHandler,
	adminFAQHandler *handler.AdminFAQHandler,
	vendorBannerHandler *handler.VendorBannerHandler,
	addressHandler *handler.AddressHandler,
	shippingHandler *handler.ShippingHandler,
	vendorCourierHandler *handler.VendorCourierHandler,
	vendorOrderHandler *handler.VendorOrderHandler,
	wsHandler *ws.Handler,
	jwtSecret string,
) *gin.Engine {

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// Global middleware
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", healthHandler.Check)

		// WebSocket — authenticated via ?token=<jwt> query param (any role)
		v1.GET("/ws", wsHandler.Connect)

		// Public catalog routes
		v1.GET("/products", catalogHandler.ListProducts)
		v1.GET("/products/:id", catalogHandler.GetProductByID)
		v1.GET("/products/:id/reviews", reviewHandler.ListByProduct)
		v1.GET("/products/:id/review-summary", reviewHandler.GetSummary)
		v1.GET("/categories", catalogHandler.ListCategories)
		v1.GET("/banners", bannerHandler.ListActiveBanners)
		v1.GET("/return-reasons", returnReasonHandler.ListActiveReturnReasons)
		v1.GET("/contacts", adminContactHandler.ListActiveAdminContacts)
		v1.GET("/order-statuses", orderActionHandler.ListStatuses)

		// Public vendor store banners
		v1.GET("/vendors/:id/banners", vendorBannerHandler.ListVendorBannersPublic)

	}

	// FAQ / Pusat Bantuan (requires auth + customer role)
	faqGroup := v1.Group("/faq")
	faqGroup.Use(middleware.Auth(jwtSecret))
	faqGroup.Use(middleware.RequireRoles("customer"))
	{
		faqGroup.GET("", faqHandler.ListFAQs)
		faqGroup.GET("/categories", faqHandler.ListCategories)
	}

	// User routes (public registration + email verification + login)
	otpGroup := v1.Group("/otp")
	{
		otpGroup.POST("/request", otpHandler.RequestOTP)
		otpGroup.POST("/verify", otpHandler.VerifyOTP)
	}

	// User routes (public registration + email verification + login)
	userGroup := v1.Group("/users")
	{
		userGroup.POST("/register", userHandler.Register)
		userGroup.GET("/verify-email", userHandler.VerifyEmail)
		userGroup.POST("/resend-verification", userHandler.ResendVerification)
		userGroup.POST("/login", userHandler.Login)
	}

	// User authenticated routes (requires auth + customer role)
	userAuth := v1.Group("/users")
	userAuth.Use(middleware.Auth(jwtSecret))
	userAuth.Use(middleware.RequireRoles("customer"))
	{
		userAuth.GET("/me", userHandler.GetMe)
		userAuth.GET("/cart", cartHandler.GetCart)
		userAuth.POST("/cart/items", cartHandler.AddItem)
		userAuth.PATCH("/cart/items/:itemId", cartHandler.UpdateItem)
		userAuth.DELETE("/cart/items/:itemId", cartHandler.RemoveItem)
		userAuth.DELETE("/cart", cartHandler.ClearCart)
		userAuth.GET("/wishlist", wishlistHandler.ListItems)
		userAuth.POST("/wishlist/items", wishlistHandler.AddItem)
		userAuth.DELETE("/wishlist/products/:productId", wishlistHandler.RemoveItem)
		userAuth.GET("/wishlist/products/:productId/status", wishlistHandler.GetProductStatus)
		userAuth.POST("/checkout/preview", checkoutHandler.Preview)
		userAuth.POST("/checkout", checkoutHandler.Checkout)
		userAuth.GET("/orders", orderActionHandler.ListOrders)
		userAuth.POST("/orders/:orderId/complete", orderActionHandler.Complete)
		userAuth.POST("/reviews/presign", reviewHandler.PresignImage)
		userAuth.POST("/reviews", reviewHandler.Create)
	}

	// Shipping location lookup (RajaOngkir proxy) — accessible by customer + vendor
	shippingGroup := v1.Group("/shipping")
	shippingGroup.Use(middleware.Auth(jwtSecret))
	shippingGroup.Use(middleware.RequireRoles("customer", "umkm"))
	{
		shippingGroup.GET("/provinces", shippingHandler.GetProvinces)
		shippingGroup.GET("/cities", shippingHandler.GetCities)
		shippingGroup.GET("/districts", shippingHandler.GetDistricts)
		shippingGroup.GET("/subdistricts", shippingHandler.GetSubdistricts)
		shippingGroup.GET("/couriers", vendorCourierHandler.ListCouriers)
	}

	// Address management — accessible by customer + vendor
	addressGroup := v1.Group("/addresses")
	addressGroup.Use(middleware.Auth(jwtSecret))
	addressGroup.Use(middleware.RequireRoles("customer", "umkm"))
	{
		addressGroup.POST("", addressHandler.CreateAddress)
		addressGroup.GET("", addressHandler.ListAddresses)
		addressGroup.GET("/:id", addressHandler.GetAddress)
		addressGroup.PATCH("/:id", addressHandler.UpdateAddress)
		addressGroup.DELETE("/:id", addressHandler.DeleteAddress)
		addressGroup.PATCH("/:id/default", addressHandler.SetDefaultAddress)
	}

	// Vendor routes (public registration + login)
	vendorGroup := v1.Group("/vendors")
	{
		vendorGroup.POST("/register/request-otp", vendorHandler.RequestRegistrationOTP)
		vendorGroup.POST("/register/verify-otp", vendorHandler.VerifyRegistrationOTP)
		vendorGroup.POST("/login", vendorHandler.Login)
	}

	vendorOnboarding := v1.Group("/vendors/register")
	vendorOnboarding.Use(middleware.AuthVendorOnboarding(jwtSecret))
	{
		vendorOnboarding.POST("/password", vendorHandler.SetRegistrationPassword)
		vendorOnboarding.POST("/store", vendorHandler.SaveRegistrationStore)
		vendorOnboarding.POST("/legal-document/presign", vendorHandler.PresignRegistrationLegalDocument)
		vendorOnboarding.POST("/legal-document/submit", vendorHandler.SubmitRegistrationLegal)
	}

	// Vendor authenticated routes (requires auth + umkm role)
	vendorAuth := v1.Group("/vendors")
	vendorAuth.Use(middleware.Auth(jwtSecret))
	vendorAuth.Use(middleware.RequireRoles("umkm"))
	{
		vendorAuth.GET("/me", vendorHandler.GetMe)
		vendorAuth.POST("/documents/confirm", vendorHandler.ConfirmDocuments)
		vendorAuth.POST("/products", productHandler.CreateProduct)
		vendorAuth.POST("/products/:id/images/confirm", productHandler.ConfirmImages)
		vendorAuth.GET("/balance", vendorHandler.GetBalance)
		vendorAuth.GET("/payout-channels", vendorHandler.ListPayoutChannels)
		vendorAuth.POST("/withdrawals", vendorHandler.RequestWithdrawal)

		// Chat (vendor-specific management)
		vendorAuth.GET("/chat", chatHandler.ListConversations)
		vendorAuth.GET("/chat/:conversationId", chatHandler.GetConversation)
		vendorAuth.PATCH("/chat/:conversationId/read", chatHandler.MarkRead)

		// Vendor Banner management
		vendorAuth.POST("/banners", vendorBannerHandler.CreateBanner)
		vendorAuth.GET("/banners", vendorBannerHandler.ListBanners)
		vendorAuth.GET("/banners/:id", vendorBannerHandler.GetBanner)
		vendorAuth.POST("/banners/:id/confirm", vendorBannerHandler.ConfirmBanner)
		vendorAuth.PATCH("/banners/:id", vendorBannerHandler.UpdateBanner)

		// Vendor Courier management
		vendorAuth.GET("/couriers", vendorCourierHandler.GetVendorCouriers)
		vendorAuth.PUT("/couriers", vendorCourierHandler.SetVendorCouriers)
		vendorAuth.DELETE("/couriers/:courierId", vendorCourierHandler.RemoveVendorCourier)
		vendorAuth.DELETE("/banners/:id", vendorBannerHandler.DeleteBanner)

		// Vendor Order management
		vendorAuth.GET("/orders", vendorOrderHandler.ListOrders)
		vendorAuth.GET("/orders/:orderId", vendorOrderHandler.GetOrderDetail)
		vendorAuth.POST("/orders/:orderId/accept", vendorOrderHandler.AcceptOrder)
		vendorAuth.POST("/orders/:orderId/reject", vendorOrderHandler.RejectOrder)
		vendorAuth.POST("/orders/:orderId/ship", vendorOrderHandler.ShipOrder)
		vendorAuth.GET("/orders/:orderId/tracking", vendorOrderHandler.TrackWaybill)
	}

	// Admin login route (public - no auth required)
	adminPublic := v1.Group("/admin")
	{
		adminPublic.POST("/login", adminLoginHandler.Login)
	}

	// Admin routes (requires auth + admin role)
	admin := v1.Group("/admin")
	admin.Use(middleware.Auth(jwtSecret))
	admin.Use(middleware.RequireRoles("admin"))
	{
		admin.GET("/me", adminUserHandler.GetMe)
		admin.POST("/users", adminUserHandler.Create)
		admin.GET("/users", adminUserHandler.List)
		admin.GET("/users/:id", adminUserHandler.GetByID)
		admin.PUT("/users/:id", adminUserHandler.Update)
		admin.DELETE("/users/:id", adminUserHandler.Delete)

		admin.GET("/vendors", adminVendorHandler.List)
		admin.GET("/vendors/:id", adminVendorHandler.GetByID)
		admin.PATCH("/vendors/:id/approve", adminVendorHandler.Approve)
		admin.PATCH("/vendors/:id/reject", adminVendorHandler.Reject)
		admin.PATCH("/vendors/:id/block", adminVendorHandler.Block)
		admin.PATCH("/vendors/:id/unblock", adminVendorHandler.Unblock)

		// Chat (admin-specific management)
		admin.GET("/chat", chatHandler.ListConversations)
		admin.GET("/chat/:conversationId", chatHandler.GetConversation)
		admin.PATCH("/chat/:conversationId/read", chatHandler.MarkRead)

		// Settings: Banner management
		admin.POST("/settings/banners", bannerHandler.CreateBanner)
		admin.GET("/settings/banners", bannerHandler.ListBanners)
		admin.GET("/settings/banners/:id", bannerHandler.GetBanner)
		admin.POST("/settings/banners/:id/confirm", bannerHandler.ConfirmBanner)
		admin.PATCH("/settings/banners/:id", bannerHandler.UpdateBanner)
		admin.DELETE("/settings/banners/:id", bannerHandler.DeleteBanner)

		// Settings: Master Lookup - Categories
		admin.POST("/settings/categories", adminCategoryHandler.CreateCategory)
		admin.GET("/settings/categories", adminCategoryHandler.ListCategories)
		admin.PATCH("/settings/categories/:id", adminCategoryHandler.UpdateCategory)
		admin.DELETE("/settings/categories/:id", adminCategoryHandler.DeleteCategory)

		// Settings: Master Lookup - Return Reasons
		admin.POST("/settings/return-reasons", returnReasonHandler.CreateReturnReason)
		admin.GET("/settings/return-reasons", returnReasonHandler.ListReturnReasons)
		admin.PATCH("/settings/return-reasons/:id", returnReasonHandler.UpdateReturnReason)
		admin.DELETE("/settings/return-reasons/:id", returnReasonHandler.DeleteReturnReason)

		// Settings: Master Lookup - Admin Contacts
		admin.POST("/settings/contacts", adminContactHandler.CreateAdminContact)
		admin.GET("/settings/contacts", adminContactHandler.ListAdminContacts)
		admin.PATCH("/settings/contacts/:id", adminContactHandler.UpdateAdminContact)
		admin.DELETE("/settings/contacts/:id", adminContactHandler.DeleteAdminContact)

		// Settings: Master Lookup - FAQ
		admin.POST("/settings/faq", adminFAQHandler.CreateFAQ)
		admin.GET("/settings/faq", adminFAQHandler.ListFAQs)
		admin.PATCH("/settings/faq/:id", adminFAQHandler.UpdateFAQ)
		admin.DELETE("/settings/faq/:id", adminFAQHandler.DeleteFAQ)
	}

	// Finance routes (requires auth + finance role)
	finance := v1.Group("/finance")
	finance.Use(middleware.Auth(jwtSecret))
	finance.Use(middleware.RequireRoles(domain.RoleFinance))
	{
		finance.GET("/dashboard", financeHandler.Dashboard)
		finance.GET("/transactions", financeHandler.ListTransactions)
		finance.GET("/transactions/export", financeHandler.ExportTransactions)
		finance.GET("/payouts", financeHandler.ListPayouts)
		finance.GET("/payouts/export", financeHandler.ExportPayouts)
		finance.GET("/refunds", financeHandler.ListRefunds)
		finance.GET("/refunds/export", financeHandler.ExportRefunds)
		finance.GET("/reports", financeHandler.FinancialReport)
		finance.GET("/reports/export", financeHandler.ExportFinancialReport)
		finance.PATCH("/refunds/:id/status", financeHandler.UpdateRefundStatus)
		finance.PATCH("/payouts/:id/status", financeHandler.UpdatePayoutStatus)

		// Chat (finance-specific management)
		finance.GET("/chat", chatHandler.ListConversations)
		finance.GET("/chat/:conversationId", chatHandler.GetConversation)
		finance.PATCH("/chat/:conversationId/read", chatHandler.MarkRead)
	}

	// Notifications routes (requires auth + admin/finance/cs role)
	notifications := v1.Group("/notifications")
	notifications.Use(middleware.Auth(jwtSecret))
	notifications.Use(middleware.RequireRoles(domain.RoleAdmin, domain.RoleFinance, domain.RoleCS))
	{
		notifications.GET("", notificationHandler.List)
		notifications.GET("/unread-count", notificationHandler.UnreadCount)
		notifications.PATCH("/read-all", notificationHandler.MarkAllRead)
		notifications.PATCH("/:id/read", notificationHandler.MarkRead)
	}

	// CS login route (public - no auth required)
	csPublic := v1.Group("/cs")
	{
		csPublic.POST("/login", csLoginHandler.Login)
	}

	// CS authenticated routes (requires auth + cs role)
	csAuth := v1.Group("/customer-service")
	csAuth.Use(middleware.Auth(jwtSecret))
	csAuth.Use(middleware.RequireRoles("cs"))
	{
		// Dashboard
		csAuth.GET("/dashboard", csDashboardHandler.GetDashboard)

		// Tickets
		csAuth.GET("/tickets", ticketHandler.ListTickets)
		csAuth.GET("/tickets/:id", ticketHandler.GetTicket)
		csAuth.PATCH("/tickets/:id/take", ticketHandler.TakeTicket)
		csAuth.PATCH("/tickets/:id/status", ticketHandler.UpdateTicketStatus)
		csAuth.POST("/tickets/:id/messages", ticketHandler.AddTicketMessage)
		csAuth.GET("/tickets/:id/messages", ticketHandler.ListTicketMessages)

		// Chat (CS-only management)
		csAuth.GET("/chat", chatHandler.ListConversations)
		csAuth.GET("/chat/:conversationId", chatHandler.GetConversation)
		csAuth.PATCH("/chat/:conversationId/read", chatHandler.MarkRead)

		// Reply Templates
		csAuth.GET("/reply-templates", replyTemplateHandler.ListTemplates)
		csAuth.POST("/reply-templates", replyTemplateHandler.CreateTemplate)
		csAuth.GET("/reply-templates/:id", replyTemplateHandler.GetTemplate)
		csAuth.PUT("/reply-templates/:id", replyTemplateHandler.UpdateTemplate)
		csAuth.DELETE("/reply-templates/:id", replyTemplateHandler.DeleteTemplate)

		// Data Pengguna
		csAuth.GET("/users", csUserHandler.ListUsers)
		csAuth.GET("/users/:id", csUserHandler.GetUser)
		csAuth.GET("/users/:id/tickets", csUserHandler.GetUserTickets)

		// Laporan
		csAuth.GET("/reports", csReportHandler.GetReport)
		csAuth.GET("/reports/export", csReportHandler.ExportReport)
	}

	// Customer: buat tiket baru (requires auth + customer role)
	customerTicket := v1.Group("/tickets")
	customerTicket.Use(middleware.Auth(jwtSecret))
	customerTicket.Use(middleware.RequireRoles("customer"))
	{
		customerTicket.GET("/subjects", ticketSubjectHandler.ListSubjects)
		customerTicket.POST("/attachment/presign", ticketHandler.PresignTicketAttachment)
		customerTicket.POST("", ticketHandler.CreateTicket)
	}

	// Chat: dapat diakses oleh semua role yang terautentikasi.
	// Validasi siapa boleh chat dengan siapa dilakukan di usecase (allowedChat map).
	chatGroup := v1.Group("/chat")
	chatGroup.Use(middleware.Auth(jwtSecret))
	{
		chatGroup.GET("/users", chatHandler.SearchChatableUsers)
		chatGroup.POST("", chatHandler.StartOrGetConversation)
		chatGroup.POST("/cs", chatHandler.StartChatWithCS)                       // Auto-assign customer ke CS
		chatGroup.POST("/attachment/presign", chatHandler.PresignChatAttachment) // Presign upload URL
		chatGroup.GET("", chatHandler.ListConversations)
		chatGroup.GET("/:conversationId", chatHandler.GetConversation)
		chatGroup.POST("/:conversationId/messages", chatHandler.SendMessage)
		chatGroup.GET("/:conversationId/messages", chatHandler.ListMessages)
		chatGroup.PATCH("/:conversationId/read", chatHandler.MarkRead)
	}

	// Webhook routes (public, no auth - verified via callback token)
	webhooks := v1.Group("/webhooks")
	{
		webhooks.POST("/xendit/invoice", checkoutHandler.Webhook)
		webhooks.POST("/xendit/payout", xenditWebhookHandler.Payout)
	}

	return r
}
