package app

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/shanth1/gotools/consts"
	"github.com/shanth1/gotools/log"
	cachememory "github.com/shanth1/gotrace/internal/adapters/cache/memory"
	"github.com/shanth1/gotrace/internal/adapters/generator"
	transport "github.com/shanth1/gotrace/internal/adapters/handler/http"
	"github.com/shanth1/gotrace/internal/adapters/ingestor"
	"github.com/shanth1/gotrace/internal/adapters/repository/cached"
	"github.com/shanth1/gotrace/internal/adapters/repository/geography"
	memoryrepo "github.com/shanth1/gotrace/internal/adapters/repository/memory"
	"github.com/shanth1/gotrace/internal/config"
	"github.com/shanth1/gotrace/internal/core/services"
)

func Run(ctx, shutdownCtx context.Context, cfg *config.Config) {
	logger := log.FromContext(ctx)

	geoProvider, err := geography.NewGeoProvider(cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("new geo ip repo")
	}

	// Cache (In-Memory)
	cache := cachememory.NewCache()

	// Repositories (In-Memory)
	baseUserRepo := memoryrepo.NewUserRepo()
	baseCampRepo := memoryrepo.NewCampaignRepo()
	baseLinkRepo := memoryrepo.NewLinkRepo()
	baseClickRepo := memoryrepo.NewAnalyticRepo()
	basePlanRepo := memoryrepo.NewPlanRepo()

	// Cached Repositories
	userRepo := cached.NewUserRepo(baseUserRepo, cache, 5*time.Minute)
	campRepo := cached.NewCampaignRepo(baseCampRepo, cache, 10*time.Minute)
	linkRepo := cached.NewLinkRepo(baseLinkRepo, cache, 10*time.Minute)
	analyticRepo := cached.NewClickRepo(baseClickRepo, cache, 30*time.Second)
	planRepo := cached.NewPlanRepo(basePlanRepo, cache, 24*time.Hour)

	if cfg.Env != consts.EnvProd {
		cfg := generator.Config{
			UsersCount:      5,
			LinksPerUser:    10,
			ClicksPerLink:   20,
			DaysHistory:     30,
			PasswordDefault: "password",
		}
		seeder := generator.New(userRepo, campRepo, linkRepo, analyticRepo, planRepo)
		_ = seeder.Seed(context.Background(), cfg)
	}

	ingestor := ingestor.NewBatchEventIngestor(analyticRepo, userRepo, geoProvider)

	// Services
	authService := services.NewAuthService(userRepo, planRepo, cfg)
	analyticsService := services.NewAnalyticsService(analyticRepo)
	campaignService := services.NewCampaignService(campRepo)
	linkService := services.NewLinkService(linkRepo, userRepo, planRepo)
	redirectService := services.NewRedirectService(ctx, ingestor, linkRepo, geoProvider)
	userService := services.NewUserService(userRepo, planRepo, campRepo, linkRepo)
	billingService := services.NewBillingService(planRepo)

	httpHandler := transport.NewRouter(
		cfg,
		authService,
		analyticsService,
		campaignService,
		linkService,
		redirectService,
		userService,
		billingService,
		linkRepo,
		userRepo,
		planRepo,
		logger,
	)

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
