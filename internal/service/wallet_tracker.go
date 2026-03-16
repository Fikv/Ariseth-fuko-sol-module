package service

import (
	"context"
	"time"

	"ariseth-fuko-sol-module/internal/cache"
	"ariseth-fuko-sol-module/internal/client"
	"ariseth-fuko-sol-module/internal/domain"
)

const DefaultWalletSnapshotCacheTTL = 5 * time.Second

type WalletTrackerService struct {
	Client        domain.WalletSnapshotGetter
	snapshotCache *cache.TTLCache[string, domain.WalletSnapshot]
	cacheTTL      time.Duration
}

func NewWalletTrackerService(walletClient *client.SolanaWalletTrackerClient) *WalletTrackerService {
	if walletClient == nil {
		walletClient = client.NewSolanaWalletTrackerClient(nil)
	}

	return &WalletTrackerService{
		Client:        walletClient,
		snapshotCache: cache.NewTTL[string, domain.WalletSnapshot](DefaultWalletSnapshotCacheTTL),
		cacheTTL:      DefaultWalletSnapshotCacheTTL,
	}
}

func (s *WalletTrackerService) SetCacheTTL(ttl time.Duration) {
	s.cacheTTL = ttl
	s.snapshotCache.SetTTL(ttl)
}

func (s *WalletTrackerService) ClearCache() {
	s.snapshotCache.Clear()
}

func (s *WalletTrackerService) TrackWallet(ctx context.Context, walletAddress string) (domain.WalletSnapshot, error) {
	if snapshot, ok := s.snapshotCache.Get(walletAddress); ok {
		return cloneWalletSnapshot(snapshot), nil
	}

	snapshot, err := s.Client.GetWalletSnapshot(ctx, walletAddress)
	if err != nil {
		return domain.WalletSnapshot{}, err
	}

	s.snapshotCache.Set(walletAddress, cloneWalletSnapshot(snapshot))
	return snapshot, nil
}

func (s *WalletTrackerService) TrackWallets(ctx context.Context, walletAddresses []string) ([]domain.WalletSnapshot, error) {
	snapshots := make([]domain.WalletSnapshot, 0, len(walletAddresses))

	for _, address := range walletAddresses {
		snapshot, err := s.TrackWallet(ctx, address)
		if err != nil {
			return nil, err
		}

		snapshots = append(snapshots, snapshot)
	}

	return snapshots, nil
}

func cloneWalletSnapshot(snapshot domain.WalletSnapshot) domain.WalletSnapshot {
	cloned := snapshot
	if snapshot.Tokens != nil {
		cloned.Tokens = append([]domain.TokenBalance(nil), snapshot.Tokens...)
	}
	return cloned
}
