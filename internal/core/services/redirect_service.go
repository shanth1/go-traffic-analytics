package services

import (
	"context"
	"errors"
	"log"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mssola/user_agent"
	"github.com/oschwald/geoip2-golang"
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

	geoDB *geoip2.Reader
}

func NewRedirectService(
	ctx context.Context,
	l ports.LinkRepository,
	c ports.ClickRepository,
	u ports.UserRepository,
) *RedirectService {
	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Printf("WARN: GeoIP database not found (GeoLite2-City.mmdb). Geo stats will be empty. Err: %v", err)
	}

	s := &RedirectService{
		linkRepo:  l,
		clickRepo: c,
		userRepo:  u,
		eventChan: make(chan trackingEvent, 1000),
		geoDB:     db,
	}

	go s.startBackgroundWorker(ctx)

	return s
}

func (s *RedirectService) ProcessRedirect(ctx context.Context, slug, ip, userAgentString, referer string) (string, error) {
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

func (s *RedirectService) parseUserAgent(uaString string) parsedUA {
	ua := user_agent.New(uaString)

	os := ua.OS()
	name, _ := ua.Browser()

	device := "Desktop"
	if ua.Mobile() {
		device = "Mobile"
	} else if ua.Bot() {
		device = "Bot"
	} else {
		// Простая эвристика для планшетов (iPad определяется как Mobile часто, но проверим)
		if strings.Contains(strings.ToLower(uaString), "ipad") || strings.Contains(strings.ToLower(uaString), "tablet") {
			device = "Tablet"
		}
	}

	// Fallbacks
	if os == "" {
		os = "Unknown"
	}
	if name == "" {
		name = "Unknown"
	}

	return parsedUA{
		OS:      os,
		Browser: name,
		Device:  device,
	}
}

func (s *RedirectService) resolveGeoIP(ipStr string) parsedGeo {
	// 1. Handle Localhost
	if ipStr == "127.0.0.1" || ipStr == "::1" {
		return parsedGeo{Country: "Local", City: "Host"}
	}

	// 2. Если база не загружена
	if s.geoDB == nil {
		return parsedGeo{Country: "Unknown", City: "Unknown"}
	}

	// 3. Парсинг IP
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return parsedGeo{Country: "Unknown", City: "Unknown"}
	}

	// 4. Поиск в базе
	record, err := s.geoDB.City(ip)
	if err != nil {
		return parsedGeo{Country: "Unknown", City: "Unknown"}
	}

	country := record.Country.IsoCode // "US", "RU"
	city := record.City.Names["en"]   // "New York"

	if country == "" {
		country = "Unknown"
	}
	if city == "" {
		city = "Unknown"
	}

	return parsedGeo{
		Country: country,
		City:    city,
	}
}

func (s *RedirectService) normalizeReferer(ref string) string {
	if ref == "" {
		return "Direct"
	}

	u, err := url.Parse(ref)
	if err != nil {
		// Если не удалось распарсить, возвращаем как есть (или обрезанный)
		if len(ref) > 50 {
			return ref[:50] + "..."
		}
		return ref
	}

	// Возвращаем только хост (google.com, t.co, facebook.com)
	// Это делает графики чище.
	if u.Host != "" {
		return strings.TrimPrefix(u.Host, "www.")
	}

	return "Unknown"
}
