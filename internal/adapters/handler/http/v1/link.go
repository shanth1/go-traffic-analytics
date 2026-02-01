package v1

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/shanth1/gotools/log"
	"github.com/shanth1/gotools/ops"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
	"github.com/shanth1/gotrace/internal/pkg/http/request"
	"github.com/shanth1/gotrace/internal/pkg/http/response"
)

type LinkHandler struct {
	linkSvc ports.LinkService
	campSvc ports.CampaignService
}

func NewLinkHandler(ls ports.LinkService, cs ports.CampaignService) *LinkHandler {
	return &LinkHandler{
		linkSvc: ls,
		campSvc: cs,
	}
}

// GetLinks godoc
// @Summary Get links list
// @Description Get links with filtering, search and pagination
// @Tags Links
// @Security BearerAuth
// @Produce json
// @Param campaign_id query string false "Filter by Campaign ID"
// @Param search query string false "Search by slug or target URL"
// @Param is_active query boolean false "Filter by active status"
// @Param limit query int false "Limit (default 10)"
// @Param offset query int false "Offset (default 0)"
// @Success 200 {object} LinksListResponse
// @Failure 401 {object} response.ErrorWrapper "Unauthorized"
// @Failure 500 {object} response.ErrorWrapper
// @Router /api/v1/links [get]
func (h *LinkHandler) GetLinks(w http.ResponseWriter, r *http.Request) {
	userID := request.GetUserID(r)

	q := r.URL.Query()

	limit, err := strconv.Atoi(q.Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(q.Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	var isActivePtr *bool
	if val := q.Get("is_active"); val != "" {
		active, err := strconv.ParseBool(val)
		if err == nil {
			isActivePtr = &active
		}
	}

	filter := domain.LinkFilter{
		UserID:     userID,
		CampaignID: q.Get("campaign_id"),
		Search:     q.Get("search"),
		IsActive:   isActivePtr,
		Limit:      limit,
		Offset:     offset,
	}

	links, total, err := h.linkSvc.GetLinkList(r.Context(), filter)
	if err != nil {
		response.Error(w, r, err)
		return
	}

	if links == nil {
		links = []*domain.Link{}
	}

	log.FromContext(r.Context()).Info().
		Int("count", len(links)).
		Int("offset", offset).
		Int("limit", limit).
		Int64("total", total).
		Msg("user_links_listed")

	response.JSON(w, r, http.StatusOK, LinksListResponse{
		Data: links,
		Meta: domain.PaginationMeta{
			Total:  total,
			Limit:  limit,
			Offset: offset,
		},
	})
}

// CreateLink godoc
// @Summary Create link
// @Description Shorten a URL
// @Tags Links
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateLinkRequest true "Link Info"
// @Success 201 {object} LinkResponse
// @Failure 400 {object} response.ErrorWrapper
// @Failure 401 {object} response.ErrorWrapper "Unauthorized"
// @Failure 403 {object} response.ErrorWrapper "Limit reached or Forbidden"
// @Failure 500 {object} response.ErrorWrapper
// @Router /api/v1/links [post]
func (h *LinkHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
	const op = "v1.LinkHandler.CreateLink"

	userID := request.GetUserID(r)

	var req CreateLinkRequest
	if err := request.DecodeJSON(w, r, &req); err != nil {
		err := ops.WrapMsg(op, ops.KindInvalid, err, "invalid request body")
		response.Error(w, r, err)
		return
	}

	link, err := h.linkSvc.CreateLink(r.Context(), domain.CreateLinkCmd{
		UserID:     userID,
		CampaignID: req.CampaignID,
		TargetURL:  req.TargetURL,
		CustomSlug: "",
	})
	if err != nil {
		response.Error(w, r, err)
		return
	}

	log.FromContext(r.Context()).Info().
		Str("link_id", string(link.ID)).
		Str("slug", link.Slug).
		Msg("link_created")

	response.Created(w, r, link)
}

// DeleteLink godoc
// @Summary Delete link
// @Description Remove a link (soft delete)
// @Tags Links
// @Security BearerAuth
// @Param id path string true "Link ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} response.ErrorWrapper
// @Failure 401 {object} response.ErrorWrapper "Unauthorized"
// @Failure 404 {object} response.ErrorWrapper "Not Found"
// @Failure 500 {object} response.ErrorWrapper
// @Router /api/v1/links/{id} [delete]
func (h *LinkHandler) DeleteLink(w http.ResponseWriter, r *http.Request) {
	const op = "v1.LinkHandler.DeleteLink"

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		response.Error(w, r, ops.WrapMsg(op, ops.KindInvalid, nil, "link id is required"))
		return
	}

	id := domain.LinkID(idStr)
	userID := request.GetUserID(r)

	if err := h.linkSvc.DeleteLink(r.Context(), userID, id); err != nil {
		response.Error(w, r, err)
		return
	}

	log.FromContext(r.Context()).Info().
		Str("link_id", string(id)).
		Msg("link_deleted")

	response.NoContent(w, r)
}
