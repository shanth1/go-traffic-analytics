package services

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type RedirectService struct {
	linkRepo  ports.LinkRepository
	clickRepo ports.ClickRepository
	userRepo  ports.UserRepository
}

func NewRedirectService(
	l ports.LinkRepository,
	c ports.ClickRepository,
	u ports.UserRepository,
) *RedirectService {
	return &RedirectService{
		linkRepo:  l,
		clickRepo: c,
		userRepo:  u,
	}
}

func (s *RedirectService) ProcessRedirect(
	ctx context.Context,
	slug string,
	ip string,
	userAgentString string,
	referer string,
) (string, error) {
	link, err := s.linkRepo.FindBySlug(ctx, slug)
	if err != nil {
		return "", err
	}

	if !link.IsActive {
		return "", errors.New("link is inactive") // domain.ErrLinkInactive
	}

	go s.trackAndBill(link, ip, userAgentString, referer)

	return link.TargetURL, nil
}

func (s *RedirectService) trackAndBill(
	link *domain.Link,
	ip, uaString, referer string,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	parsedUA := s.parseUserAgent(uaString)
	parsedGeo := s.resolveGeoIP(ip)

	click := &domain.ClickEvent{
		ID:        uuid.New().String(),
		LinkID:    link.ID,
		Timestamp: time.Now().UTC(),
		IP:        ip,
		Country:   parsedGeo.Country,
		City:      parsedGeo.City,
		OS:        parsedUA.OS,
		Browser:   parsedUA.Browser,
		Device:    parsedUA.Device,
		Referer:   s.normalizeReferer(referer),
	}

	if err := s.clickRepo.Save(ctx, click); err != nil {
		log.Printf("ERROR: failed to save click analytics: %v", err)
	}

	if link.CampaignID != "" {
		userID := link.UserID //

		if userID != "" {
			if err := s.userRepo.IncrementClickCount(ctx, userID); err != nil {
				log.Printf("ERROR: failed to increment user quota: %v", err)
			}
		}
	}
}

// --- Helper Structures & Methods (Private) ---

type parsedUA struct {
	OS      string
	Browser string
	Device  string
}

type parsedGeo struct {
	Country string
	City    string
}

// parseUserAgent - простая заглушка.
// В продакшене использовать библиотеку "github.com/mssola/user_agent"
func (s *RedirectService) parseUserAgent(ua string) parsedUA {
	uaLower := strings.ToLower(ua)
	res := parsedUA{
		OS:      "Unknown",
		Browser: "Unknown",
		Device:  "Desktop", // Default
	}

	// Очень примитивная эвристика для MVP
	if strings.Contains(uaLower, "iphone") || strings.Contains(uaLower, "android") {
		res.Device = "Mobile"
	}

	if strings.Contains(uaLower, "windows") {
		res.OS = "Windows"
	} else if strings.Contains(uaLower, "mac os") {
		res.OS = "macOS"
	} else if strings.Contains(uaLower, "android") {
		res.OS = "Android"
	} else if strings.Contains(uaLower, "iphone") {
		res.OS = "iOS"
	}

	if strings.Contains(uaLower, "chrome") {
		res.Browser = "Chrome"
	} else if strings.Contains(uaLower, "safari") {
		res.Browser = "Safari"
	} else if strings.Contains(uaLower, "firefox") {
		res.Browser = "Firefox"
	}

	return res
}

// resolveGeoIP - простая заглушка.
// В продакшене использовать GeoLite2 базу и либу "github.com/oschwald/geoip2-golang"
func (s *RedirectService) resolveGeoIP(ip string) parsedGeo {
	// Симуляция для MVP/Localhost
	if ip == "127.0.0.1" || ip == "::1" {
		return parsedGeo{Country: "Local", City: "Host"}
	}

	// Здесь можно вставить реальный вызов сервиса или базы
	return parsedGeo{
		Country: "US", // Default stub
		City:    "Unknown",
	}
}

func (s *RedirectService) normalizeReferer(ref string) string {
	if ref == "" {
		return "Direct"
	}
	// Можно обрезать до домена, чтобы не хранить полные пути
	// ex: https://google.com/search?q=... -> google.com
	return ref
}
