package app

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/config"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/handler"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/router"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/infrastructure/database"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase"
	"gorm.io/gorm"
)

// App holds all initialized application components.
type App struct {
	Config *config.Config
	DB     *gorm.DB
	Router *gin.Engine
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

	// Wire dependencies
	healthUseCase := usecase.NewHealthUseCase()
	healthHandler := handler.NewHealthHandler(healthUseCase)

	// Setup router
	r := router.NewRouter(healthHandler, cfg.JWT.Secret)

	return &App{
		Config: cfg,
		DB:     db,
		Router: r,
	}, nil
}

// Close cleans up application resources.
func (a *App) Close() error {
	sqlDB, err := a.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	return sqlDB.Close()
}
