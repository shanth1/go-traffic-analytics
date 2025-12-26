package v1

import (
	"net/http"

	"github.com/shanth1/gotrace/internal/core/ports"
	"github.com/shanth1/gotrace/internal/pkg/http/response"
	"github.com/shanth1/gotrace/internal/pkg/request"
)

type UserHandler struct {
	userSvc ports.UserService
}

func NewUserHandler(us ports.UserService) *UserHandler {
	return &UserHandler{
		userSvc: us,
	}
}

// GetProfileTree godoc
// @Summary      Profile Hierarchy
// @Description  Get hierarchical structure of User -> Campaigns -> Links for Tree visualization
// @Tags         Campaigns
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  domain.HierarchyNode
// @Failure      401  {object}  response.ErrorResponse "Unauthorized"
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/campaigns/tree [get]
func (h *UserHandler) GetProfileTree(w http.ResponseWriter, r *http.Request) {
	userID := request.GetUserID(r)

	tree, err := h.userSvc.GetHierarchy(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, ProfileTreeResponse{Data: tree})
}
