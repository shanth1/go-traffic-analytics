package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5"
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

// TODO: added logger
func (m *QuotaMiddleware) CheckClickLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")

		ctx := r.Context()

		link, err := m.LinkRepo.FindBySlug(ctx, slug)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		user, err := m.UserRepo.FindByID(ctx, link.UserID)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		plan, err := m.PlanRepo.FindByID(ctx, user.PlanID)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		if plan.MaxClicksMonth != -1 {
			if user.ClicksCurrentMonth >= plan.MaxClicksMonth {
				// TODO: alert
				next.ServeHTTP(w, r)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
