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
	linkRepo    ports.LinkRepository
	ingestor    ports.EventIngestor
	geoProvider ports.GeoProvider

	fallbackLogger log.Logger
	eventChan      chan trackingEvent
}

func NewRedirectService(
	ctx context.Context,
	i ports.EventIngestor,
	lr ports.LinkRepository,
	gp ports.GeoProvider,
) *RedirectService {
	fallbackLogger := log.FromContext(ctx)
	s := &RedirectService{
		ingestor:    i,
		linkRepo:    lr,
		geoProvider: gp,

		fallbackLogger: fallbackLogger,
		eventChan:      make(chan trackingEvent, 1000),
	}

	go s.startBackgroundWorker(ctx)

	return s
}

func (s *RedirectService) Process(ctx context.Context, slug string, meta domain.RequestMetadata) (string, error) {
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
		ip:      meta.IP,
		ua:      meta.UserAgent,
		referer: meta.Referer,
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

	location, err := s.geoProvider.Lookup(ctx, evt.ip)
	if err != nil {
		location.Country = consts.Unknown
		location.City = consts.Unknown
		logger.Warn().Err(err).Msg("get geo ip info")
	}

	event := &domain.ClickEvent{
		ID:         uuid.New().String(),
		UserID:     evt.link.UserID,
		LinkID:     evt.link.ID,
		CampaignID: evt.link.CampaignID,
		Timestamp:  evt.time,
		IP:         evt.ip,
		Country:    location.Country,
		City:       location.City,
		OS:         parsedUA.OS,
		Browser:    parsedUA.Browser,
		Device:     parsedUA.Device,
		Referer:    s.normalizeReferer(evt.referer),
	}

	if err := s.ingestor.TrackClick(ctx, event); err != nil {
		logger.Error().Err(err).Msg("track click")
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
		if len(ref) > 50 {
			return ref[:50] + "..."
		}
		return ref
	}

	if u.Host != "" {
		return strings.TrimPrefix(u.Host, "www.")
	}

	return consts.Unknown
}
