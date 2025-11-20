package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(
	e *echo.Echo,
	authH *v1.AuthHandler,
	linkH *v1.LinkHandler,
	redirectH *v1.RedirectHandler,
	analyticsH *v1.AnalyticsHandler,
	adminH *v1.AdminHandler,
	quotaMW *myMiddleware.QuotaMiddleware,
	authMW echo.MiddlewareFunc, // JWT Middleware
	adminMW echo.MiddlewareFunc, // Admin Only Middleware
) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// --- Public Routes ---
	e.GET("/:slug", redirectH.Redirect, quotaMW.CheckClickLimit) // <-- Главный роут с проверкой квот

	api := e.Group("/api/v1")

	// Auth
	api.POST("/auth/register", authH.Register)
	api.POST("/auth/login", authH.Login)

	// --- Client Routes (Protected) ---
	client := api.Group("", authMW)

	// Campaigns
	client.GET("/campaigns", linkH.GetCampaigns)
	client.POST("/campaigns", linkH.CreateCampaign)

	// Links
	client.GET("/campaigns/:id/links", linkH.GetLinksByCampaign)
	client.POST("/links", linkH.CreateLink)
	client.DELETE("/links/:id", linkH.DeleteLink)

	// Analytics (Visx Ready)
	client.GET("/analytics/summary", analyticsH.GetSummary)
	client.GET("/analytics/stream", analyticsH.GetStreamGraph)   // Area Chart / Stream
	client.GET("/analytics/flow", analyticsH.GetSankeyFlow)      // Sankey
	client.GET("/analytics/geo", analyticsH.GetGeoMap)           // Map
	client.GET("/analytics/quality", analyticsH.GetQualityRadar) // Radar

	// --- Admin Routes (Protected + Admin Role) ---
	admin := api.Group("/admin", authMW, adminMW)

	admin.GET("/users", adminH.GetUsers)
	admin.PATCH("/users/:id/status", adminH.UpdateUserStatus) // Ban/Unban
	admin.PATCH("/users/:id/plan", adminH.UpdateUserPlan)
	admin.GET("/plans", adminH.GetPlans)
	// admin.PUT("/plans", adminH.UpdatePlan) // (Optional)
}
