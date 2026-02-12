package router

import (
	"github.com/gin-gonic/gin"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/handler"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/middleware"
)

// NewRouter sets up the Gin engine with middleware and route registration.
func NewRouter(healthHandler *handler.HealthHandler, adminUserHandler *handler.AdminUserHandler, authHandler *handler.AuthHandler, vendorHandler *handler.VendorHandler, jwtSecret string) *gin.Engine {
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
	}

	// Auth routes (public)
	authGroup := v1.Group("/auth")
	{
		authGroup.POST("/login", authHandler.Login)
	}

	// Vendor routes (public registration)
	vendorGroup := v1.Group("/vendors")
	{
		vendorGroup.POST("/register", vendorHandler.Register)
	}

	// Vendor authenticated routes (requires auth + umkm role)
	vendorAuth := v1.Group("/vendors")
	vendorAuth.Use(middleware.Auth(jwtSecret))
	vendorAuth.Use(middleware.RequireRoles("umkm"))
	{
		vendorAuth.POST("/documents/confirm", vendorHandler.ConfirmDocuments)
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
	}

	return r
}
