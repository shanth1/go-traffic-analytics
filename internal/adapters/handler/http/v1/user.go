package v1

import (
	"net/http"

	"github.com/shanth1/gotools/log"
	"github.com/shanth1/gotrace/internal/core/ports"
	"github.com/shanth1/gotrace/internal/pkg/http/request"
	"github.com/shanth1/gotrace/internal/pkg/http/response"
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
// @Summary Profile Hierarchy
// @Description Get hierarchical structure of User -> Campaigns -> Links for Tree visualization
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} ProfileTreeResponse
// @Failure 401 {object} response.ErrorWrapper "Unauthorized"
// @Failure 500 {object} response.ErrorWrapper
// @Router /api/v1/users/tree [get]
func (h *UserHandler) GetProfileTree(w http.ResponseWriter, r *http.Request) {
	userID := request.GetUserID(r)

	tree, err := h.userSvc.GetHierarchy(r.Context(), userID)
	if err != nil {
		response.Error(w, r, err)
		return
	}

	log.FromContext(r.Context()).Info().Msg("profile_tree_retrieved")

	response.OK(w, r, tree)
}
