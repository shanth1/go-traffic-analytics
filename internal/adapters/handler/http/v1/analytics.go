package v1

import (
	"net/http"
	"time"

	"github.com/shanth1/gotools/log"
	"github.com/shanth1/gotrace/internal/core/domain"
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

func (h *AnalyticsHandler) parseFilter(r *http.Request) domain.AnalyticsFilter {
	query := r.URL.Query()

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

	return domain.AnalyticsFilter{
		CampaignID: query.Get("campaign_id"),
		LinkID:     domain.LinkID(query.Get("link_id")),
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
		response.ServerError(w, r, err)
		return
	}

	log.FromContext(r.Context()).Info().Str("campaign_id", filter.CampaignID).Str("link_id", string(filter.LinkID)).Msg("analytics_summary_retrieved")

	response.Success(w, r, summary)
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
// @Param        from        query string false "Date From (RFC3339)"
// @Param        to          query string false "Date To (RFC3339)"
// @Success      200  {object}  StreamGraphResponse
// @Failure      401  {object}  response.ErrorResponse "Unauthorized"
// @Failure      403  {object}  response.ErrorResponse "Forbidden"
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/analytics/stream [get]
func (h *AnalyticsHandler) GetStreamGraph(w http.ResponseWriter, r *http.Request) {
	filter := h.parseFilter(r)

	groupBy := r.URL.Query().Get("group_by")
	if groupBy == "" {
		groupBy = consts.OS
	}

	data, err := h.service.GetStreamGraphData(r.Context(), filter, groupBy)
	if err != nil {
		response.ServerError(w, r, err)
		return
	}

	log.FromContext(r.Context()).Info().Str("group_by", groupBy).Str("campaign_id", filter.CampaignID).Str("link_id", string(filter.LinkID)).Msg("analytics_stream_graph_retrieved")

	response.Success(w, r, response.Envelope{"data": data})
}

// GetSankeyFlow godoc
// @Summary      Sankey flow
// @Description  Get flow data (Referer -> Device -> Country)
// @Tags         Analytics
// @Security     BearerAuth
// @Produce      json
// @Param        campaign_id query string false "Filter by Campaign"
// @Param        link_id     query string false "Filter by Link"
// @Param        from        query string false "Date From (RFC3339)"
// @Param        to          query string false "Date To (RFC3339)"
// @Success      200  {object}  SankeyResponse
// @Failure      401  {object}  response.ErrorResponse "Unauthorized"
// @Failure      403  {object}  response.ErrorResponse "Forbidden"
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/analytics/flow [get]
func (h *AnalyticsHandler) GetSankeyFlow(w http.ResponseWriter, r *http.Request) {
	filter := h.parseFilter(r)
	stages := []string{consts.Referer, consts.Device, consts.Country}

	data, err := h.service.GetSankeyData(r.Context(), filter, stages)
	if err != nil {
		response.ServerError(w, r, err)
		return
	}

	log.FromContext(r.Context()).Info().Str("campaign_id", filter.CampaignID).Str("link_id", string(filter.LinkID)).Msg("analytics_sankey_flow_retrieved")

	response.Success(w, r, response.Envelope{"data": data})
}

// GetGeoMap godoc
// @Summary      Geo distribution
// @Description  Get clicks count by country
// @Tags         Analytics
// @Security     BearerAuth
// @Produce      json
// @Param        campaign_id query string false "Filter by Campaign"
// @Param        link_id     query string false "Filter by Link"
// @Param        from        query string false "Date From (RFC3339)"
// @Param        to          query string false "Date To (RFC3339)"
// @Success      200  {object}  GeoResponse
// @Failure      401  {object}  response.ErrorResponse "Unauthorized"
// @Failure      403  {object}  response.ErrorResponse "Forbidden"
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/analytics/geo [get]
func (h *AnalyticsHandler) GetGeoMap(w http.ResponseWriter, r *http.Request) {
	filter := h.parseFilter(r)

	data, err := h.service.GetGeoDistribution(r.Context(), filter)
	if err != nil {
		response.ServerError(w, r, err)
		return
	}

	log.FromContext(r.Context()).Info().Str("campaign_id", filter.CampaignID).Str("link_id", string(filter.LinkID)).Msg("analytics_geo_map_retrieved")

	response.Success(w, r, response.Envelope{"data": data})
}

// GetQualityRadar godoc
// @Summary      Traffic quality
// @Description  Get quality metrics for Radar Chart
// @Tags         Analytics
// @Security     BearerAuth
// @Produce      json
// @Param        campaign_id query string false "Filter by Campaign"
// @Param        link_id     query string false "Filter by Link"
// @Param        from        query string false "Date From (RFC3339)"
// @Param        to          query string false "Date To (RFC3339)"
// @Success      200  {object}  QualityResponse
// @Failure      401  {object}  response.ErrorResponse "Unauthorized"
// @Failure      403  {object}  response.ErrorResponse "Forbidden"
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/analytics/quality [get]
func (h *AnalyticsHandler) GetQualityRadar(w http.ResponseWriter, r *http.Request) {
	filter := h.parseFilter(r)

	data, err := h.service.GetTrafficQuality(r.Context(), filter)
	if err != nil {
		response.ServerError(w, r, err)
		return
	}

	log.FromContext(r.Context()).Info().Str("campaign_id", filter.CampaignID).Str("link_id", string(filter.LinkID)).Msg("analytics_quality_radar_retrieved")

	response.Success(w, r, response.Envelope{"data": data})
}

// GetHeatmap godoc
// @Summary      Heatmap data
// @Description  Get click intensity by Day of Week and Hour
// @Tags         Analytics
// @Security     BearerAuth
// @Produce      json
// @Param        campaign_id query string false "Filter"
// @Param        link_id     query string false "Filter"
// @Param        from        query string false "Date From (RFC3339)"
// @Param        to          query string false "Date To (RFC3339)"
// @Success      200  {object}  HeatmapResponse
// @Failure      401  {object}  response.ErrorResponse "Unauthorized"
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/analytics/heatmap [get]
func (h *AnalyticsHandler) GetHeatmap(w http.ResponseWriter, r *http.Request) {
	filter := h.parseFilter(r)
	data, err := h.service.GetHeatmapData(r.Context(), filter)
	if err != nil {
		response.ServerError(w, r, err)
		return
	}

	if data == nil {
		data = []domain.HeatmapPoint{}
	}

	log.FromContext(r.Context()).Info().Str("campaign_id", filter.CampaignID).Str("link_id", string(filter.LinkID)).Msg("analytics_heatmap_retrieved")

	response.Success(w, r, response.Envelope{"data": data})
}

// GetStats godoc
// @Summary      Category Stats (Pie/Bar)
// @Description  Get top metrics for a dimension (browser, os, device)
// @Tags         Analytics
// @Security     BearerAuth
// @Produce      json
// @Param        dimension   query string false  "browser, os, device"
// @Param        campaign_id query string false "Filter by Campaign"
// @Param        link_id     query string false "Filter by Link"
// @Param        from        query string false "Date From (RFC3339)"
// @Param        to          query string false "Date To (RFC3339)"
// @Success      200  {object}  StatsResponse
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
		response.ServerError(w, r, err)
		return
	}

	log.FromContext(r.Context()).Info().Str("dimension", dim).Str("campaign_id", filter.CampaignID).Str("link_id", string(filter.LinkID)).Msg("analytics_stats_retrieved")

	response.Success(w, r, response.Envelope{"data": data})
}
