package app

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/shanth1/gotools/log"
	transport "github.com/shanth1/gotrace/internal/adapters/handler/http"
	"github.com/shanth1/gotrace/internal/adapters/handler/http/middleware"
	v1 "github.com/shanth1/gotrace/internal/adapters/handler/http/v1"
	"github.com/shanth1/gotrace/internal/adapters/repository/memory"
	"github.com/shanth1/gotrace/internal/config"
	"github.com/shanth1/gotrace/internal/core/services"
	"github.com/shanth1/gotrace/internal/pkg/generator"
)

func Run(ctx, shutdownCtx context.Context, cfg *config.Config) {
	logger := log.FromContext(ctx)

	// Repositories (In-Memory)
	userRepo := memory.NewUserRepo()
	campRepo := memory.NewCampaignRepo()
	linkRepo := memory.NewLinkRepo()
	clickRepo := memory.NewClickRepo()
	planRepo := memory.NewPlanRepo()

	// Services
	authService := services.NewAuthService(userRepo, planRepo, cfg)
	linkService := services.NewLinkService(linkRepo, campRepo, userRepo, planRepo)
	redirectService := services.NewRedirectService(ctx, linkRepo, clickRepo, userRepo)
	analyticsService := services.NewAnalyticsService(clickRepo)
	userService := services.NewUserService(userRepo, planRepo)

	// Handlers
	authHandler := v1.NewAuthHandler(authService)
	linkHandler := v1.NewLinkHandler(linkService)
	redirectHandler := v1.NewRedirectHandler(redirectService)
	analyticsHandler := v1.NewAnalyticsHandler(analyticsService)
	adminHandler := v1.NewAdminHandler(userService)

	// Middleware
	quotaMW := middleware.NewQuotaMiddleware(linkRepo, userRepo, planRepo)
	jwtMW := middleware.Auth(cfg)
	adminMW := middleware.AdminOnly

	httpHandler := transport.NewRouter(
		authHandler,
		linkHandler,
		redirectHandler,
		analyticsHandler,
		adminHandler,
		quotaMW,
		jwtMW,
		adminMW,
	)
	// TODO: local/develop env
	// Seeding Data (Optional, for dev convenience)
	seeder := &generator.DataSeeder{
		UserRepo:     userRepo,
		CampaignRepo: campRepo,
		LinkRepo:     linkRepo,
		ClickRepo:    clickRepo,
		PlanRepo:     planRepo,
	}
	seeder.SeedFullTopology()

	runHTTPServer(ctx, shutdownCtx, cfg, httpHandler, logger)
}

func runHTTPServer(
	ctx context.Context,
	shutdownCtx context.Context,
	cfg *config.Config,
	handler http.Handler,
	logger log.Logger,
) {
	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           handler,
		ReadTimeout:       cfg.HTTP.ReadTimeout,
		WriteTimeout:      cfg.HTTP.WriteTimeout,
		IdleTimeout:       cfg.HTTP.IdleTimeout,
		ReadHeaderTimeout: 2 * time.Second,
	}

	go func() {
		logger.Info().Msgf("starting HTTP server on %s", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal().Err(err).Msg("http server failed")
		}
	}()

	<-ctx.Done()
	logger.Info().Msg("shutting down HTTP server...")

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("http server graceful shutdown failed")
	} else {
		logger.Info().Msg("http server stopped gracefully")
	}
}
