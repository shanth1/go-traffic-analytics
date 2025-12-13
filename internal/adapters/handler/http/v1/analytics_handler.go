package v1

import (
	"net/http"
	"time"

	"github.com/shanth1/gotrace/internal/core/ports"
	"github.com/shanth1/gotrace/internal/pkg/consts"
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

// GetSummary godoc
// @Summary      Analytics summary
// @Description  Get aggregated stats
// @Tags         Analytics
// @Security     BearerAuth
// @Produce      json
// @Param        campaign_id query string false "Filter by Campaign"
// @Param        link_id     query string false "Filter by Link"
// @Param        from        query string false "Date From (RFC3339)"
// @Param        to          query string false "Date To (RFC3339)"
// @Success      200  {object}  AnalyticsSummaryResponse
// @Failure      401  {object}  response.ErrorResponse "Unauthorized"
// @Failure      403  {object}  response.ErrorResponse "Forbidden"
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/analytics/summary [get]
func (h *AnalyticsHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	filter := h.parseFilter(r)

	summary, err := h.service.GetSummary(r.Context(), filter)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, AnalyticsSummaryResponse{Data: summary})
}

// GetStreamGraph godoc
// @Summary      Stream graph data
// @Description  Get time-series data for Streamgraph/Stacked Area charts
// @Tags         Analytics
// @Security     BearerAuth
// @Produce      json
// @Param        group_by    query string false "os, browser, country"
// @Param        campaign_id query string false "Filter by Campaign"
// @Param        link_id     query string false "Filter by Link"
// @Param        from        query string false "Date From"
// @Param        to          query string false "Date To"
// @Success      200  {object}  StreamGraphResponse
// @Failure      401  {object}  response.ErrorResponse "Unauthorized"
// @Failure      403  {object}  response.ErrorResponse "Forbidden"
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/analytics/stream [get]
func (h *AnalyticsHandler) GetStreamGraph(w http.ResponseWriter, r *http.Request) {
	filter := h.parseFilter(r)

	groupBy := r.URL.Query().Get("group_by") // os, browser, country
	if groupBy == "" {
		groupBy = consts.OS
	}

	data, err := h.service.GetStreamGraphData(r.Context(), filter, groupBy)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, StreamGraphResponse{Data: data})
}

// GetSankeyFlow godoc
// @Summary      Sankey flow
// @Description  Get flow data (Referer -> Device -> Country)
// @Tags         Analytics
// @Security     BearerAuth
// @Produce      json
// @Param        campaign_id query string false "Filter by Campaign"
// @Param        link_id     query string false "Filter by Link"
// @Param        from        query string false "Date From"
// @Param        to          query string false "Date To"
// @Success      200  {object}  SankeyResponse
// @Failure      401  {object}  response.ErrorResponse "Unauthorized"
// @Failure      403  {object}  response.ErrorResponse "Forbidden"
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/analytics/flow [get]
func (h *AnalyticsHandler) GetSankeyFlow(w http.ResponseWriter, r *http.Request) {
	filter := h.parseFilter(r)
	// stages=referer,device,country
	stages := []string{consts.Referer, consts.Device, consts.Country}

	data, err := h.service.GetSankeyData(r.Context(), filter, stages)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, SankeyResponse{Data: data})
}

// GetGeoMap godoc
// @Summary      Geo distribution
// @Description  Get clicks count by country
// @Tags         Analytics
// @Security     BearerAuth
// @Produce      json
// @Param        campaign_id query string false "Filter by Campaign"
// @Param        link_id     query string false "Filter by Link"
// @Param        from        query string false "Date From"
// @Param        to          query string false "Date To"
// @Success      200  {object}  GeoResponse
// @Failure      401  {object}  response.ErrorResponse "Unauthorized"
// @Failure      403  {object}  response.ErrorResponse "Forbidden"
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/analytics/geo [get]
func (h *AnalyticsHandler) GetGeoMap(w http.ResponseWriter, r *http.Request) {
	filter := h.parseFilter(r)

	data, err := h.service.GetGeoDistribution(r.Context(), filter)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, GeoResponse{Data: data})
}

// GetQualityRadar godoc
// @Summary      Traffic quality
// @Description  Get quality metrics for Radar Chart
// @Tags         Analytics
// @Security     BearerAuth
// @Produce      json
// @Param        campaign_id query string false "Filter by Campaign"
// @Param        link_id     query string false "Filter by Link"
// @Param        from        query string false "Date From"
// @Param        to          query string false "Date To"
// @Success      200  {object}  QualityResponse
// @Failure      401  {object}  response.ErrorResponse "Unauthorized"
// @Failure      403  {object}  response.ErrorResponse "Forbidden"
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/analytics/quality [get]
func (h *AnalyticsHandler) GetQualityRadar(w http.ResponseWriter, r *http.Request) {
	filter := h.parseFilter(r)

	data, err := h.service.GetTrafficQuality(r.Context(), filter)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, QualityResponse{Data: data})
}

// GetHeatmap godoc
// @Summary      Heatmap data
// @Description  Get click intensity by Day of Week and Hour
// @Tags         Analytics
// @Security     BearerAuth
// @Produce      json
// @Param        campaign_id query string false "Filter"
// @Param        link_id     query string false "Filter"
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  response.ErrorResponse "Unauthorized"
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/analytics/heatmap [get]
func (h *AnalyticsHandler) GetHeatmap(w http.ResponseWriter, r *http.Request) {
	filter := h.parseFilter(r)
	data, err := h.service.GetHeatmapData(r.Context(), filter)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]interface{}{"data": data})
}

// GetStats godoc
// @Summary      Category Stats (Pie/Bar)
// @Description  Get top metrics for a dimension (browser, os, device)
// @Tags         Analytics
// @Security     BearerAuth
// @Produce      json
// @Param        dimension   query string true  "browser, os, device"
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  response.ErrorResponse "Unauthorized"
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/analytics/stats [get]
func (h *AnalyticsHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	filter := h.parseFilter(r)
	dim := r.URL.Query().Get("dimension")
	if dim == "" {
		dim = "browser"
	}

	data, err := h.service.GetCategoryStats(r.Context(), filter, dim)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]interface{}{"data": data})
}
