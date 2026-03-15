package service

import (
	"context"
	"time"

	"ariseth-fuko-sol-module/internal/cache"
	"ariseth-fuko-sol-module/internal/client"
	"ariseth-fuko-sol-module/internal/domain"
)

const DefaultPreMatchSnapshotCacheTTL = 30 * time.Second

type PreMatchScraperService struct {
	Client        domain.PreMatchPageGetter
	snapshotCache *cache.TTLCache[string, domain.PreMatchSnapshot]
	cacheTTL      time.Duration
}

func NewPreMatchScraperService(scraperClient *client.PreMatchScraperClient, baseURL string) *PreMatchScraperService {
	if scraperClient == nil {
		scraperClient = client.NewPreMatchScraperClient(nil, baseURL)
	}

	return &PreMatchScraperService{
		Client:        scraperClient,
		snapshotCache: cache.NewTTL[string, domain.PreMatchSnapshot](DefaultPreMatchSnapshotCacheTTL),
		cacheTTL:      DefaultPreMatchSnapshotCacheTTL,
	}
}

func (s *PreMatchScraperService) SetCacheTTL(ttl time.Duration) {
	s.cacheTTL = ttl
	s.snapshotCache.SetTTL(ttl)
}

func (s *PreMatchScraperService) ClearCache() {
	s.snapshotCache.Clear()
}

func (s *PreMatchScraperService) GetPreMatchSnapshot(ctx context.Context, sportSlug string) (domain.PreMatchSnapshot, error) {
	if snapshot, ok := s.snapshotCache.Get(sportSlug); ok {
		return snapshot, nil
	}

	html, err := s.Client.GetPreMatchHTML(ctx, sportSlug)
	if err != nil {
		return domain.PreMatchSnapshot{}, err
	}

	snapshot := domain.PreMatchSnapshot{
		SportSlug: sportSlug,
		HTML:      html,
		FetchedAt: time.Now().UTC(),
	}
	s.snapshotCache.Set(sportSlug, snapshot)

	return snapshot, nil
}
