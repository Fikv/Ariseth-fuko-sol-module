package scanner

import (
	"context"
	"log"
	"sync"
	"time"
)

type Source interface {
	Name() string
	Run(ctx context.Context, out chan<- Event) error
}

type Service struct {
	store   *Store
	sources []Source
	logger  *log.Logger
}

func NewService(store *Store, logger *log.Logger, sources ...Source) *Service {
	return &Service{
		store:   store,
		sources: sources,
		logger:  logger,
	}
}

func (s *Service) Run(ctx context.Context) error {
	out := make(chan Event, 128)
	errCh := make(chan error, len(s.sources))

	var wg sync.WaitGroup
	for _, src := range s.sources {
		source := src
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := source.Run(ctx, out); err != nil && ctx.Err() == nil {
				errCh <- err
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case err, ok := <-errCh:
			if !ok {
				return nil
			}
			if err != nil {
				return err
			}
		case e := <-out:
			if e.SeenAt.IsZero() {
				e.SeenAt = time.Now().UTC()
			}
			if s.store.Add(e) && s.logger != nil {
				s.logger.Printf("new pool: chain=%s dex=%s pool=%s base=%s quote=%s", e.Chain, e.DEX, e.PoolAddress, e.BaseMint, e.QuoteMint)
			}
		}
	}
}
