package main

import (
	"log"

	"github.com/labstack/echo/v4"

	"github.com/shanth1/gotrace/internal/config"
	"github.com/shanth1/gotrace/internal/utils/generator"

	// Repositories
	"github.com/shanth1/gotrace/internal/adapters/repository/memory"

	// Services
	"github.com/shanth1/gotrace/internal/core/services"

	// Handlers & Middleware
	handler "github.com/shanth1/gotrace/internal/adapters/handler/http"
	mw "github.com/shanth1/gotrace/internal/adapters/handler/http/middleware"
	v1 "github.com/shanth1/gotrace/internal/adapters/handler/http/v1"
)

func main() {
	// 1. Config
	cfg := config.LoadConfig()

	// 2. Repositories (In-Memory)
	userRepo := memory.NewUserRepo()
	campRepo := memory.NewCampaignRepo()
	linkRepo := memory.NewLinkRepo()
	clickRepo := memory.NewClickRepo()
	planRepo := memory.NewPlanRepo()

	// 3. Services
	authService := services.NewAuthService(userRepo, planRepo, cfg)
	linkService := services.NewLinkService(linkRepo, campRepo, userRepo, planRepo)
	redirectService := services.NewRedirectService(linkRepo, clickRepo, userRepo)
	analyticsService := services.NewAnalyticsService(clickRepo)
	userService := services.NewUserService(userRepo, planRepo)

	// 4. Handlers
	authHandler := v1.NewAuthHandler(authService)
	linkHandler := v1.NewLinkHandler(linkService)
	redirectHandler := v1.NewRedirectHandler(redirectService)
	analyticsHandler := v1.NewAnalyticsHandler(analyticsService)
	adminHandler := v1.NewAdminHandler(userService)

	// 5. Middleware
	quotaMW := mw.NewQuotaMiddleware(linkRepo, userRepo, planRepo)
	jwtMW := handler.InitJWTMiddleware(cfg)
	adminMW := handler.InitAdminMiddleware()

	// 6. Seeding Data (Optional, for dev convenience)
	seeder := &generator.DataSeeder{
		UserRepo:     userRepo,
		CampaignRepo: campRepo,
		LinkRepo:     linkRepo,
		ClickRepo:    clickRepo,
		PlanRepo:     planRepo,
	}
	seeder.SeedFullTopology()

	// 7. Server
	e := echo.New()

	// Init Router
	handler.NewRouter(
		e,
		authHandler,
		linkHandler,
		redirectHandler,
		analyticsHandler,
		adminHandler,
		quotaMW,
		jwtMW,
		adminMW,
	)

	log.Printf("🚀 Server starting on port %s", cfg.Port)
	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
