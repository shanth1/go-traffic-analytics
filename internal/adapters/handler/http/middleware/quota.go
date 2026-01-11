package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/shanth1/gotools/log"
	"github.com/shanth1/gotools/logkeys"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type QuotaMiddleware struct {
	LinkRepo ports.LinkRepository
	UserRepo ports.UserRepository
	PlanRepo ports.PlanRepository
}

func NewQuotaMiddleware(l ports.LinkRepository, u ports.UserRepository, p ports.PlanRepository) *QuotaMiddleware {
	return &QuotaMiddleware{LinkRepo: l, UserRepo: u, PlanRepo: p}
}

func (m *QuotaMiddleware) CheckClickLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		slug := chi.URLParam(r, "slug")

		logger := log.FromContext(ctx).With(log.Str("slug", slug))

		link, err := m.LinkRepo.FindBySlug(ctx, slug)
		if err != nil {
			logger.Error().Err(err).Msg("find_link_failed")
			next.ServeHTTP(w, r)
			return
		}

		user, err := m.UserRepo.FindByID(ctx, link.UserID)
		if err != nil {
			logger.Error().Err(err).Any("link_id", link.ID).Msg("find_user_failed")
			next.ServeHTTP(w, r)
			return
		}

		plan, err := m.PlanRepo.FindByID(ctx, user.PlanID)
		if err != nil {
			logger.Error().Err(err).Any("link_id", link.ID).Any(logkeys.UserID, user.ID).Msg("find_plan_failed")
			next.ServeHTTP(w, r)
			return
		}

		if plan.MaxClicksMonth != -1 {
			if user.ClicksCurrentMonth >= plan.MaxClicksMonth {
				logger.Warn().Err(err).Any("link_id", link.ID).Any(logkeys.UserID, user.ID).Any("plan_id", plan.ID).Msg("click_limit_reached")
				next.ServeHTTP(w, r)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
