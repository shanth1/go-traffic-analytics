package v1

import (
	"net/http"
	"time"

	"github.com/shanth1/gotrace/internal/core/ports"
	"github.com/shanth1/gotrace/internal/pkg/http/response"
)

type AnalyticsHandler struct {
	service ports.AnalyticsService
}

func NewAnalyticsHandler(s ports.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: s}
}

func (h *AnalyticsHandler) parseFilter(r *http.Request) ports.AnalyticsFilter {
	query := r.URL.Query()

	// Defaults
	to := time.Now()
	from := to.AddDate(0, 0, -7) // Last 7 days

	if t := query.Get("to"); t != "" {
		if parsed, err := time.Parse(time.RFC3339, t); err == nil {
			to = parsed
		}
	}
	if f := query.Get("from"); f != "" {
		if parsed, err := time.Parse(time.RFC3339, f); err == nil {
			from = parsed
		}
	}

	return ports.AnalyticsFilter{
		CampaignID: query.Get("campaign_id"),
		LinkID:     query.Get("link_id"),
		From:       from,
		To:         to,
	}
}

// GET /analytics/summary
func (h *AnalyticsHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	filter := h.parseFilter(r)

	summary, err := h.service.GetSummary(r.Context(), filter)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{"data": summary})
}

// GET /analytics/stream?group_by=os&interval=day
func (h *AnalyticsHandler) GetStreamGraph(w http.ResponseWriter, r *http.Request) {
	filter := h.parseFilter(r)

	groupBy := r.URL.Query().Get("group_by") // os, browser, country
	if groupBy == "" {
		groupBy = "os"
	}

	data, err := h.service.GetStreamGraphData(r.Context(), filter, groupBy)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{"data": data})
}

// GET /analytics/flow (Sankey)
func (h *AnalyticsHandler) GetSankeyFlow(w http.ResponseWriter, r *http.Request) {
	filter := h.parseFilter(r)
	// stages=referer,device,country
	stages := []string{"referer", "device", "country"}

	data, err := h.service.GetSankeyData(r.Context(), filter, stages)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{"data": data})
}

// GET /analytics/geo
func (h *AnalyticsHandler) GetGeoMap(w http.ResponseWriter, r *http.Request) {
	filter := h.parseFilter(r)

	data, err := h.service.GetGeoDistribution(r.Context(), filter)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{"data": data})
}

// GET /analytics/quality
func (h *AnalyticsHandler) GetQualityRadar(w http.ResponseWriter, r *http.Request) {
	filter := h.parseFilter(r)

	data, err := h.service.GetTrafficQuality(r.Context(), filter)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{"data": data})
}
