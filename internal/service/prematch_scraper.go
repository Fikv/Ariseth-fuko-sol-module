package service

import (
	"context"
	"time"

	"ariseth-fuko-sol-module/internal/client"
	"ariseth-fuko-sol-module/internal/domain"
)

type PreMatchScraperService domain.PreMatchScraperService

func NewPreMatchScraperService(scraperClient *client.PreMatchScraperClient, baseURL string) *PreMatchScraperService {
	if scraperClient == nil {
		scraperClient = client.NewPreMatchScraperClient(nil, baseURL)
	}

	return &PreMatchScraperService{
		Client: scraperClient,
	}
}

func (s *PreMatchScraperService) GetPreMatchSnapshot(ctx context.Context, sportSlug string) (domain.PreMatchSnapshot, error) {
	html, err := s.Client.GetPreMatchHTML(ctx, sportSlug)
	if err != nil {
		return domain.PreMatchSnapshot{}, err
	}

	return domain.PreMatchSnapshot{
		SportSlug: sportSlug,
		HTML:      html,
		FetchedAt: time.Now().UTC(),
	}, nil
}
