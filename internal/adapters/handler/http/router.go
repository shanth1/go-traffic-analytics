package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/shanth1/gotools/log"
	httpMw "github.com/shanth1/gotrace/internal/adapters/handler/http/middleware"
	v1 "github.com/shanth1/gotrace/internal/adapters/handler/http/v1"
	"github.com/shanth1/gotrace/internal/config"
	"github.com/shanth1/gotrace/internal/core/ports"
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
	r.Use(middleware.Timeout(cfg.HTTP.RequestTimeout))

	// Handlers
	authHandler := v1.NewAuthHandler(authService)
	linkHandler := v1.NewLinkHandler(linkService)
	redirectHandler := v1.NewRedirectHandler(redirectService)
	analyticsHandler := v1.NewAnalyticsHandler(analyticsService)
	adminHandler := v1.NewAdminHandler(userService, logger)

	// Middleware
	quotaMiddleware := httpMw.NewQuotaMiddleware(linkRepo, userRepo, planRepo)
	authMiddleware := httpMw.Auth(cfg)

	// --- Public Routes ---
	r.Get("/health", adminHandler.HealthCheck)
	r.With(quotaMiddleware.CheckClickLimit).Get("/{slug}", redirectHandler.Redirect)

	// --- API v1 Group ---
	r.Route("/api/v1", func(r chi.Router) {
		// --- Auth Routes (Public) ---
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
		})

		// --- Client Routes (Protected) ---
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)

			// Campaigns
			r.Get("/campaigns", linkHandler.GetCampaigns)
			r.Post("/campaigns", linkHandler.CreateCampaign)

			// Links
			r.Get("/campaigns/{id}/links", linkHandler.GetLinksByCampaign)
			r.Post("/links", linkHandler.CreateLink)
			r.Delete("/links/{id}", linkHandler.DeleteLink)

			// Analytics (Visx Ready)
			r.Route("/analytics", func(r chi.Router) {
				r.Get("/summary", analyticsHandler.GetSummary)
				r.Get("/stream", analyticsHandler.GetStreamGraph)
				r.Get("/flow", analyticsHandler.GetSankeyFlow)
				r.Get("/geo", analyticsHandler.GetGeoMap)
				r.Get("/quality", analyticsHandler.GetQualityRadar)
			})
		})

		// --- Admin Routes (Protected + Admin Role) ---
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware, httpMw.AdminOnly)

			r.Route("/admin", func(r chi.Router) {
				r.Get("/users", adminHandler.GetUsers)
				r.Patch("/users/{id}/status", adminHandler.UpdateUserStatus)
				r.Patch("/users/{id}/plan", adminHandler.UpdateUserPlan)
				r.Get("/plans", adminHandler.GetPlans)
			})
		})
	})

	return r
}
