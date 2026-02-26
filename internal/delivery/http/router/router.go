package router

import (
	"github.com/gin-gonic/gin"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/handler"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/middleware"
	ws "github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/websocket"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

// NewRouter sets up the Gin engine with middleware and route registration.
func NewRouter(healthHandler *handler.HealthHandler, adminUserHandler *handler.AdminUserHandler, adminVendorHandler *handler.AdminVendorHandler, adminLoginHandler *handler.AdminLoginHandler, vendorHandler *handler.VendorHandler, productHandler *handler.ProductHandler, catalogHandler *handler.CatalogHandler, userHandler *handler.UserHandler, cartHandler *handler.CartHandler, financeHandler *handler.FinanceHandler, csLoginHandler *handler.CSLoginHandler, ticketHandler *handler.TicketHandler, chatHandler *handler.ChatHandler, replyTemplateHandler *handler.ReplyTemplateHandler, csDashboardHandler *handler.CSDashboardHandler, csUserHandler *handler.CSUserHandler, csReportHandler *handler.CSReportHandler, ticketSubjectHandler *handler.TicketSubjectHandler, wsHandler *ws.Handler, checkoutHandler *handler.CheckoutHandler, jwtSecret string) *gin.Engine {
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
		v1.GET("/categories", catalogHandler.ListCategories)
		v1.GET("/shipping-services", catalogHandler.ListShippingServices)
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
		userAuth.POST("/checkout", checkoutHandler.Checkout)
	}

	// Vendor routes (public registration + login)
	vendorGroup := v1.Group("/vendors")
	{
		vendorGroup.POST("/register", vendorHandler.Register)
		vendorGroup.POST("/login", vendorHandler.Login)
	}

	// Vendor authenticated routes (requires auth + umkm role)
	vendorAuth := v1.Group("/vendors")
	vendorAuth.Use(middleware.Auth(jwtSecret))
	vendorAuth.Use(middleware.RequireRoles("umkm"))
	{
		vendorAuth.POST("/documents/confirm", vendorHandler.ConfirmDocuments)
		vendorAuth.POST("/products", productHandler.CreateProduct)
		vendorAuth.POST("/products/:id/images/confirm", productHandler.ConfirmImages)
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
		finance.PATCH("/refunds/:id/status", financeHandler.UpdateRefundStatus)
		finance.PATCH("/payouts/:id/status", financeHandler.UpdatePayoutStatus)
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
	}

	return r
}
