package v1

import (
	"net/http"

	"github.com/shanth1/gotools/log"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
	"github.com/shanth1/gotrace/internal/pkg/http/response"
)

type BillingHandler struct {
	service ports.BillingService
}

func NewBillingHandler(s ports.BillingService) *BillingHandler {
	return &BillingHandler{service: s}
}

// GetPlans godoc
// @Summary List plans
// @Description Get all available plans
// @Tags Billing
// @Produce json
// @Success 200 {object} PlansListResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/billing/plans [get]
func (h *BillingHandler) GetPlans(w http.ResponseWriter, r *http.Request) {
	plans, err := h.service.GetAllPlans(r.Context())
	if err != nil {
		response.Error(w, r, err)
		return
	}

	if plans == nil {
		plans = []*domain.Plan{}
	}

	log.FromContext(r.Context()).Info().Msg("billing_plans_retrieved")

	response.OK(w, r, plans)
}
