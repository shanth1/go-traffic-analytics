package v1

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/shanth1/gotrace/internal/core/ports"
	"github.com/shanth1/gotrace/internal/core/services"
)

type AnalyticsHandler struct {
	service *services.AnalyticsService
}

func NewAnalyticsHandler(s *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: s}
}

// Helper to parse query params
func (h *AnalyticsHandler) parseFilter(c echo.Context) ports.AnalyticsFilter {
	// Defaults
	to := time.Now()
	from := to.AddDate(0, 0, -7) // Last 7 days

	if t := c.QueryParam("to"); t != "" {
		if parsed, err := time.Parse(time.RFC3339, t); err == nil {
			to = parsed
		}
	}
	if f := c.QueryParam("from"); f != "" {
		if parsed, err := time.Parse(time.RFC3339, f); err == nil {
			from = parsed
		}
	}

	return ports.AnalyticsFilter{
		CampaignID: c.QueryParam("campaign_id"),
		LinkID:     c.QueryParam("link_id"),
		From:       from,
		To:         to,
	}
}

// GET /analytics/summary
func (h *AnalyticsHandler) GetSummary(c echo.Context) error {
	filter := h.parseFilter(c)
	summary, err := h.service.GetSummary(c.Request().Context(), filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"data": summary})
}

// GET /analytics/stream?group_by=os&interval=day
func (h *AnalyticsHandler) GetStreamGraph(c echo.Context) error {
	filter := h.parseFilter(c)
	groupBy := c.QueryParam("group_by") // os, browser, country
	if groupBy == "" {
		groupBy = "os"
	}

	data, err := h.service.GetStreamGraphData(c.Request().Context(), filter, groupBy)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"data": data})
}

// GET /analytics/flow (Sankey)
func (h *AnalyticsHandler) GetSankeyFlow(c echo.Context) error {
	filter := h.parseFilter(c)
	// Жестко задаем этапы потока для начала, либо берем из query params
	// stages=referer,device,country
	stages := []string{"referer", "device", "country"}

	data, err := h.service.GetSankeyData(c.Request().Context(), filter, stages)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"data": data})
}

// GET /analytics/geo
func (h *AnalyticsHandler) GetGeoMap(c echo.Context) error {
	filter := h.parseFilter(c)
	data, err := h.service.GetGeoDistribution(c.Request().Context(), filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"data": data})
}

// GET /analytics/quality
func (h *AnalyticsHandler) GetQualityRadar(c echo.Context) error {
	filter := h.parseFilter(c)
	data, err := h.service.GetTrafficQuality(c.Request().Context(), filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"data": data})
}
