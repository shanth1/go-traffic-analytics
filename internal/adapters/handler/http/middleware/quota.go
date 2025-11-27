package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
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

// CheckClickLimit - middleware, которое вешается на роут редиректа GET /{slug}
func (m *QuotaMiddleware) CheckClickLimit(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		slug := c.Param("slug")

		// 1. Нам нужно быстро узнать, чей это slug.
		// В идеале это должно быть в кэше (Redis).
		// Для MVP делаем быстрый lookup в LinkRepo (In-Memory это мгновенно).
		link, err := m.LinkRepo.FindBySlug(c.Request().Context(), slug)
		if err != nil {
			return next(c) // Если ссылки нет, пусть RedirectHandler вернет 404
		}

		// 2. Находим кампанию, чтобы узнать владельца (User)
		// (Тут лучше денормализация: хранить OwnerID прямо в Link, чтобы не делать лишний запрос)
		// Предположим, мы добавили UserID в Link struct
		user, err := m.UserRepo.FindByID(c.Request().Context(), link.UserID)
		if err != nil {
			return next(c)
		}

		// 3. Получаем план
		plan, err := m.PlanRepo.FindByID(c.Request().Context(), user.PlanID)
		if err != nil {
			return next(c)
		}

		// 4. ПРОВЕРКА (Gatekeeper)
		if plan.MaxClicksMonth != -1 {
			if user.ClicksCurrentMonth >= plan.MaxClicksMonth {
				// ЛИМИТ ИСЧЕРПАН!
				// Варианты действий:
				// A. Вернуть 402 Payment Required (плохо для юзера)
				// B. Редирект на заглушку "Service Suspended" (стандарт)
				return c.Redirect(http.StatusTemporaryRedirect, "https://tracebit.com/limit-reached")
			}
		}

		// Если ок - пропускаем дальше к редиректу
		return next(c)
	}
}
