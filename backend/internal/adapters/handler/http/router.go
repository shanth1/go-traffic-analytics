package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	myMiddleware "github.com/shanth1/gotrace/internal/adapters/handler/http/middleware"
	v1 "github.com/shanth1/gotrace/internal/adapters/handler/http/v1"
)

func NewRouter(
	e *echo.Echo,
	authH *v1.AuthHandler,
	linkH *v1.LinkHandler,
	redirectH *v1.RedirectHandler,
	analyticsH *v1.AnalyticsHandler,
	adminH *v1.AdminHandler,
	quotaMW *myMiddleware.QuotaMiddleware,
	authMW echo.MiddlewareFunc,
	adminMW echo.MiddlewareFunc,
) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// --- Public Routes ---
	e.GET("/:slug", redirectH.Redirect, quotaMW.CheckClickLimit)

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
	client.GET("/analytics/stream", analyticsH.GetStreamGraph)
	client.GET("/analytics/flow", analyticsH.GetSankeyFlow)
	client.GET("/analytics/geo", analyticsH.GetGeoMap)
	client.GET("/analytics/quality", analyticsH.GetQualityRadar)

	// --- Admin Routes (Protected + Admin Role) ---
	admin := api.Group("/admin", authMW, adminMW)

	admin.GET("/users", adminH.GetUsers)
	admin.PATCH("/users/:id/status", adminH.UpdateUserStatus)
	admin.PATCH("/users/:id/plan", adminH.UpdateUserPlan)
	admin.GET("/plans", adminH.GetPlans)
}
