package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shanth1/gotools/consts"
	"github.com/shanth1/gotools/log"
	_ "github.com/shanth1/gotrace/docs"
	"github.com/shanth1/gotrace/internal/adapters/handler/http/handlers"
	httpMw "github.com/shanth1/gotrace/internal/adapters/handler/http/middleware"
	v1 "github.com/shanth1/gotrace/internal/adapters/handler/http/v1"
	"github.com/shanth1/gotrace/internal/config"
	"github.com/shanth1/gotrace/internal/core/ports"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(
	cfg *config.Config,
	authService ports.AuthService,
	analyticsService ports.AnalyticsService,
	linkService ports.LinkService,
	redirectService ports.RedirectService,
	userService ports.UserService,
	linkRepo ports.LinkRepository,
	userRepo ports.UserRepository,
	planRepo ports.PlanRepository,
	logger log.Logger,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(httpMw.Logger(logger))
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(httpMw.Metrics)
	r.Use(middleware.Timeout(cfg.HTTP.RequestTimeout))

	redirectHandler := handlers.NewRedirectHandler(redirectService)

	// Handlers
	authHandlerV1 := v1.NewAuthHandler(authService)
	adminHandlerV1 := v1.NewAdminHandler(userService, logger)
	linkHandlerV1 := v1.NewLinkHandler(linkService)
	analyticsHandlerV1 := v1.NewAnalyticsHandler(analyticsService)

	// Middleware
	quotaMiddleware := httpMw.NewQuotaMiddleware(linkRepo, userRepo, planRepo)
	jwtAuthMiddleware := httpMw.JWTAuth(cfg)

	// --- Public Routes ---
	r.Get("/health", handlers.HealthCheck)
	r.With(quotaMiddleware.CheckClickLimit).Get("/{slug}", redirectHandler.Redirect)
	r.Group(func(sys chi.Router) {
		if cfg.Metrics.User != "" && cfg.Metrics.Password != "" {
			sys.Use(httpMw.BasicAuth(cfg.Metrics.User, cfg.Metrics.Password))
		}
		sys.Handle("/metrics", promhttp.Handler())
	})
	if cfg.Env != consts.EnvProd {
		logger.Info().Msg("Swagger UI enabled at /swagger/index.html")
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
		))
	}

	// --- API v1 Group ---
	r.Route("/api/v1", func(r chi.Router) {
		// --- Auth Routes (Public) ---
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandlerV1.Register)
			r.Post("/login", authHandlerV1.Login)
		})

		// --- Client Routes (Protected) ---
		r.Group(func(r chi.Router) {
			r.Use(jwtAuthMiddleware)

			// Campaigns
			r.Get("/campaigns", linkHandlerV1.GetCampaigns)
			r.Post("/campaigns", linkHandlerV1.CreateCampaign)
			r.Get("/campaigns/tree", linkHandlerV1.GetProfileTree)

			// Links
			r.Get("/campaigns/{id}/links", linkHandlerV1.GetLinksByCampaign)
			r.Post("/links", linkHandlerV1.CreateLink)
			r.Delete("/links/{id}", linkHandlerV1.DeleteLink)

			// Analytics (Visx Ready)
			r.Route("/analytics", func(r chi.Router) {
				r.Get("/summary", analyticsHandlerV1.GetSummary)
				r.Get("/stream", analyticsHandlerV1.GetStreamGraph)
				r.Get("/flow", analyticsHandlerV1.GetSankeyFlow)
				r.Get("/geo", analyticsHandlerV1.GetGeoMap)
				r.Get("/quality", analyticsHandlerV1.GetQualityRadar)
			})

			r.Get("/heatmap", analyticsHandlerV1.GetHeatmap) // Для Heatmap
			r.Get("/stats", analyticsHandlerV1.GetStats)     // Для BarGroup, Pies
		})

		// --- Admin Routes (Protected + Admin Role) ---
		r.Group(func(r chi.Router) {
			r.Use(jwtAuthMiddleware, httpMw.AdminOnly)

			r.Route("/admin", func(r chi.Router) {
				r.Get("/users", adminHandlerV1.GetUsers)
				r.Patch("/users/{id}/status", adminHandlerV1.UpdateUserStatus)
				r.Patch("/users/{id}/plan", adminHandlerV1.UpdateUserPlan)
				r.Get("/plans", adminHandlerV1.GetPlans)
			})
		})
	})

	return r
}
