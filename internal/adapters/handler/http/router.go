package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	myMiddleware "github.com/shanth1/gotrace/internal/adapters/handler/http/middleware"
	v1 "github.com/shanth1/gotrace/internal/adapters/handler/http/v1"
)

func NewRouter(
	authH *v1.AuthHandler,
	linkH *v1.LinkHandler,
	redirectH *v1.RedirectHandler,
	analyticsH *v1.AnalyticsHandler,
	adminH *v1.AdminHandler,
	quotaMW *myMiddleware.QuotaMiddleware,
	authMW func(http.Handler) http.Handler,
	adminMW func(http.Handler) http.Handler,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// --- Public Routes ---
	r.With(quotaMW.CheckClickLimit).Get("/{slug}", redirectH.Redirect)

	// --- API v1 Group ---
	r.Route("/api/v1", func(r chi.Router) {
		// --- Auth Routes (Public) ---
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authH.Register)
			r.Post("/login", authH.Login)
		})

		// --- Client Routes (Protected) ---
		r.Group(func(r chi.Router) {
			r.Use(authMW)

			// Campaigns
			r.Get("/campaigns", linkH.GetCampaigns)
			r.Post("/campaigns", linkH.CreateCampaign)

			// Links
			r.Get("/campaigns/{id}/links", linkH.GetLinksByCampaign)
			r.Post("/links", linkH.CreateLink)
			r.Delete("/links/{id}", linkH.DeleteLink)

			// Analytics (Visx Ready)
			r.Route("/analytics", func(r chi.Router) {
				r.Get("/summary", analyticsH.GetSummary)
				r.Get("/stream", analyticsH.GetStreamGraph)
				r.Get("/flow", analyticsH.GetSankeyFlow)
				r.Get("/geo", analyticsH.GetGeoMap)
				r.Get("/quality", analyticsH.GetQualityRadar)
			})
		})

		// --- Admin Routes (Protected + Admin Role) ---
		r.Group(func(r chi.Router) {
			r.Use(authMW, adminMW)

			r.Route("/admin", func(r chi.Router) {
				r.Get("/users", adminH.GetUsers)
				r.Patch("/users/{id}/status", adminH.UpdateUserStatus)
				r.Patch("/users/{id}/plan", adminH.UpdateUserPlan)
				r.Get("/plans", adminH.GetPlans)
			})
		})
	})

	return r
}
