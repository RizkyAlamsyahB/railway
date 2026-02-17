package router

import (
	"github.com/gin-gonic/gin"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/handler"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/middleware"
)

// NewRouter sets up the Gin engine with middleware and route registration.
func NewRouter(healthHandler *handler.HealthHandler, adminUserHandler *handler.AdminUserHandler, adminVendorHandler *handler.AdminVendorHandler, adminLoginHandler *handler.AdminLoginHandler, vendorHandler *handler.VendorHandler, productHandler *handler.ProductHandler, catalogHandler *handler.CatalogHandler, userHandler *handler.UserHandler, jwtSecret string) *gin.Engine {
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

		// Public catalog routes
		v1.GET("/categories", catalogHandler.ListCategories)
		v1.GET("/shipping-services", catalogHandler.ListShippingServices)
	}

	// User routes (public registration + email verification)
	userGroup := v1.Group("/users")
	{
		userGroup.POST("/register", userHandler.Register)
		userGroup.GET("/verify-email", userHandler.VerifyEmail)
		userGroup.POST("/resend-verification", userHandler.ResendVerification)
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

	return r
}
