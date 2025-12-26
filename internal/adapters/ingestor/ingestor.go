package adapters

import (
	"context"
	"sync"
	"time"

	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

const (
	batchSize     = 1000
	flushInterval = 5 * time.Second
	bufferSize    = 10000
)

type BatchEventIngestor struct {
	analyticsRepo ports.AnalyticsRepository
	userRepo      ports.UserRepository
	geoProvider   ports.GeoProvider

	eventChan chan *domain.ClickEvent

	mu           sync.Mutex
	eventsBatch  []*domain.ClickEvent
	userCounters map[domain.UserID]int

	done chan struct{}
	wg   sync.WaitGroup
}

func NewBatchEventIngestor(
	ar ports.AnalyticsRepository,
	ur ports.UserRepository,
	gp ports.GeoProvider,
) *BatchEventIngestor {
	s := &BatchEventIngestor{
		analyticsRepo: ar,
		userRepo:      ur,
		geoProvider:   gp,
		eventChan:     make(chan *domain.ClickEvent, bufferSize),
		eventsBatch:   make([]*domain.ClickEvent, 0, batchSize),
		userCounters:  make(map[domain.UserID]int),
		done:          make(chan struct{}),
	}

	s.wg.Add(1)
	go s.loop()

	return s
}

func (s *BatchEventIngestor) TrackClick(ctx context.Context, event *domain.ClickEvent) error {
	select {
	case s.eventChan <- event:
		return nil
	default:
		// TODO: logging: Msg("EventIngestor buffer overflow, dropping event")
		return nil
	}
}

func (s *BatchEventIngestor) Close() error {
	close(s.done)
	s.wg.Wait()
	return nil
}

func (s *BatchEventIngestor) loop() {
	defer s.wg.Done()

	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	for {
		select {
		case event := <-s.eventChan:
			s.processEvent(event)
		case <-ticker.C:
			s.flush()
		case <-s.done:
			for len(s.eventChan) > 0 {
				s.processEvent(<-s.eventChan)
			}
			s.flush()
			return
		}
	}
}

func (s *BatchEventIngestor) processEvent(event *domain.ClickEvent) {
	s.eventsBatch = append(s.eventsBatch, event)

	if event.UserID != "" {
		s.userCounters[event.UserID]++
	}

	if len(s.eventsBatch) >= batchSize {
		s.flush()
	}
}

func (s *BatchEventIngestor) flush() {
	if len(s.eventsBatch) == 0 && len(s.userCounters) == 0 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if len(s.eventsBatch) > 0 {
		if err := s.analyticsRepo.SaveBatch(ctx, s.eventsBatch); err != nil {
			// TODO: logging: Msg("Failed to flush analytics batch")
			// TODO: retry logic
		}
		s.eventsBatch = s.eventsBatch[:0]
	}

	if len(s.userCounters) > 0 {
		for userID, count := range s.userCounters {
			if err := s.userRepo.IncrementUsage(ctx, userID, count); err != nil {
				// TODO: log.Error().Str("user_id", string(userID)).Err(err).Msg("Failed to update user usage")
			}
		}
		s.userCounters = make(map[domain.UserID]int)
	}
}
