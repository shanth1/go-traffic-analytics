package v1

import (
	"net/http"

	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
	"github.com/shanth1/gotrace/internal/pkg/http/request"
	"github.com/shanth1/gotrace/internal/pkg/http/response"
)

type FeedbackHandler struct {
	svc ports.FeedbackService
}

func NewFeedbackHandler(svc ports.FeedbackService) *FeedbackHandler {
	return &FeedbackHandler{svc: svc}
}

type feedbackRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

// SendFeedback godoc
// @Summary Send feedback
// @Description Send a support request to the admin
// @Tags Feedback
// @Accept json
// @Produce json
// @Param request body feedbackRequest true "Feedback Info"
// @Success 200 {object} StatusResponse "Status: received"
// @Router /api/v1/feedback [post]
func (h *FeedbackHandler) SendFeedback(w http.ResponseWriter, r *http.Request) {
	var req feedbackRequest
	if err := request.DecodeJSON(w, r, &req); err != nil {
		response.Error(w, r, err)
		return
	}

	userID := request.GetUserID(r)

	cmd := domain.SendFeedbackCmd{
		Name:    req.Name,
		Email:   req.Email,
		Message: req.Message,
		UserID:  domain.UserID(userID),
	}

	if err := h.svc.SendFeedback(r.Context(), cmd); err != nil {
		response.Error(w, r, err)
		return
	}

	response.OK(w, r, StatusData{Status: "received"})
}
