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

type trackingEvent struct {
	link    *domain.Link
	ip      string
	ua      string
	referer string
	time    time.Time
}

type RedirectService struct {
	linkRepo  ports.LinkRepository
	clickRepo ports.ClickRepository
	userRepo  ports.UserRepository

	eventChan chan trackingEvent
}

func NewRedirectService(
	ctx context.Context,
	l ports.LinkRepository,
	c ports.ClickRepository,
	u ports.UserRepository,
) *RedirectService {
	s := &RedirectService{
		linkRepo:  l,
		clickRepo: c,
		userRepo:  u,
		eventChan: make(chan trackingEvent, 1000),
	}

	go s.startBackgroundWorker(ctx)

	return s
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
		return "", errors.New("link is inactive")
	}

	select {
	case s.eventChan <- trackingEvent{
		link:    link,
		ip:      ip,
		ua:      userAgentString,
		referer: referer,
		time:    time.Now().UTC(),
	}:
	default:
		log.Printf("WARN: analytics buffer full, dropping click for link %s", link.ID)
	}

	return link.TargetURL, nil
}

func (s *RedirectService) startBackgroundWorker(appCtx context.Context) {
	for {
		select {
		case evt := <-s.eventChan:
			s.processEvent(evt)
		case <-appCtx.Done():
			// TODO: write to db
			log.Println("RedirectService: stopping analytics worker")
			return
		}
	}
}

func (s *RedirectService) processEvent(evt trackingEvent) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	parsedUA := s.parseUserAgent(evt.ua)
	parsedGeo := s.resolveGeoIP(evt.ip)

	click := &domain.ClickEvent{
		ID:        uuid.New().String(),
		LinkID:    evt.link.ID,
		Timestamp: evt.time,
		IP:        evt.ip,
		Country:   parsedGeo.Country,
		City:      parsedGeo.City,
		OS:        parsedUA.OS,
		Browser:   parsedUA.Browser,
		Device:    parsedUA.Device,
		Referer:   s.normalizeReferer(evt.referer),
	}

	if err := s.clickRepo.Save(ctx, click); err != nil {
		log.Printf("ERROR: failed to save click analytics: %v", err)
	}

	if evt.link.CampaignID != "" && evt.link.UserID != "" {
		if err := s.userRepo.IncrementClickCount(ctx, evt.link.UserID); err != nil {
			log.Printf("ERROR: failed to increment user quota: %v", err)
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

// parseUserAgent - mock
// TODO: "github.com/mssola/user_agent"
func (s *RedirectService) parseUserAgent(ua string) parsedUA {
	uaLower := strings.ToLower(ua)
	res := parsedUA{
		OS:      "Unknown",
		Browser: "Unknown",
		Device:  "Desktop", // Default
	}

	if strings.Contains(uaLower, "iphone") || strings.Contains(uaLower, "android") {
		res.Device = "Mobile"
	}

	switch {
	case strings.Contains(uaLower, "windows"):
		res.OS = "Windows"
	case strings.Contains(uaLower, "mac os"):
		res.OS = "macOS"
	case strings.Contains(uaLower, "android"):
		res.OS = "Android"
	case strings.Contains(uaLower, "iphone"):
		res.OS = "iOS"
	}

	switch {
	case strings.Contains(uaLower, "chrome"):
		res.Browser = "Chrome"
	case strings.Contains(uaLower, "safari"):
		res.Browser = "Safari"
	case strings.Contains(uaLower, "firefox"):
		res.Browser = "Firefox"
	}

	return res
}

// resolveGeoIP - mock
// TODO: "github.com/oschwald/geoip2-golang"
func (s *RedirectService) resolveGeoIP(ip string) parsedGeo {
	// TODO:
	if ip == "127.0.0.1" || ip == "::1" {
		return parsedGeo{Country: "Local", City: "Host"}
	}

	return parsedGeo{
		Country: "US", // Default stub
		City:    "Unknown",
	}
}

func (s *RedirectService) normalizeReferer(ref string) string {
	if ref == "" {
		return "Direct"
	}

	// TODO: https://google.com/search?q=... -> google.com
	return ref
}
