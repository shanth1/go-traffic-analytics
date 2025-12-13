package services

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mssola/user_agent"
	"github.com/shanth1/gotools/log"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
	"github.com/shanth1/gotrace/internal/pkg/consts"
)

type trackingEvent struct {
	ctx     context.Context
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
	geoIPRepo ports.GeoIPRepository

	fallbackLogger log.Logger
	eventChan      chan trackingEvent
}

func NewRedirectService(
	ctx context.Context,
	l ports.LinkRepository,
	c ports.ClickRepository,
	u ports.UserRepository,
	g ports.GeoIPRepository,
) *RedirectService {
	fallbackLogger := log.FromContext(ctx)
	s := &RedirectService{
		linkRepo:  l,
		clickRepo: c,
		userRepo:  u,
		geoIPRepo: g,

		fallbackLogger: fallbackLogger,
		eventChan:      make(chan trackingEvent, 1000),
	}

	go s.startBackgroundWorker(ctx)

	return s
}

func (s *RedirectService) ProcessRedirect(ctx context.Context, slug, ip, userAgentString, referer string) (string, error) {
	logger := log.FromContextOr(ctx, s.fallbackLogger)

	link, err := s.linkRepo.FindBySlug(ctx, slug)
	if err != nil {
		return "", err
	}

	if !link.IsActive {
		return "", errors.New("link is inactive")
	}

	select {
	case s.eventChan <- trackingEvent{
		ctx:     ctx,
		link:    link,
		ip:      ip,
		ua:      userAgentString,
		referer: referer,
		time:    time.Now().UTC(),
	}:
	default:
		logger.Warn().Msgf("analytics buffer full, dropping click for link %s", link.ID)
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
			s.fallbackLogger.Info().Msg("RedirectService: stopping analytics worker")
			return
		}
	}
}

func (s *RedirectService) processEvent(evt trackingEvent) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger := log.FromContextOr(ctx, s.fallbackLogger)

	parsedUA := s.parseUserAgent(evt.ua)

	country, city, err := s.geoIPRepo.GetInfo(ctx, evt.ip)
	if err != nil {
		country = consts.Unknown
		city = consts.Unknown
		logger.Warn().Err(err).Msg("get geo ip info")
	}

	click := &domain.ClickEvent{
		ID:        uuid.New().String(),
		LinkID:    evt.link.ID,
		Timestamp: evt.time,
		IP:        evt.ip,
		Country:   country,
		City:      city,
		OS:        parsedUA.OS,
		Browser:   parsedUA.Browser,
		Device:    parsedUA.Device,
		Referer:   s.normalizeReferer(evt.referer),
	}

	if err := s.clickRepo.Save(ctx, click); err != nil {
		logger.Error().Err(err).Msg("save click analytics")
	}

	if evt.link.CampaignID != "" && evt.link.UserID != "" {
		if err := s.userRepo.IncrementClickCount(ctx, evt.link.UserID); err != nil {
			logger.Error().Err(err).Msg("increment user quota")
		}
	}
}

// --- Helper Structures & Methods (Private) ---

type parsedUA struct {
	OS      string
	Browser string
	Device  string
}

func (s *RedirectService) parseUserAgent(uaString string) parsedUA {
	ua := user_agent.New(uaString)

	os := ua.OS()
	name, _ := ua.Browser()
	ua.UA()

	device := "Desktop"
	switch {
	case ua.Mobile():
		device = "Mobile"
	case ua.Bot():
		device = "Bot"
	case strings.Contains(strings.ToLower(uaString), "ipad") || strings.Contains(strings.ToLower(uaString), "tablet"):
		device = "Tablet"
	}

	// Fallbacks
	if os == "" {
		os = consts.Unknown
	}
	if name == "" {
		name = consts.Unknown
	}

	return parsedUA{
		OS:      os,
		Browser: name,
		Device:  device,
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

	return consts.Unknown
}
