package service

import (
	"context"

	"ariseth-fuko-sol-module/internal/client"
	"ariseth-fuko-sol-module/internal/domain"
)

type WalletTrackerService domain.WalletTrackerService

func NewWalletTrackerService(walletClient *client.SolanaWalletTrackerClient) *WalletTrackerService {
	if walletClient == nil {
		walletClient = client.NewSolanaWalletTrackerClient(nil)
	}

	return &WalletTrackerService{
		Client: walletClient,
	}
}

func (s *WalletTrackerService) TrackWallet(ctx context.Context, walletAddress string) (domain.WalletSnapshot, error) {
	return s.Client.GetWalletSnapshot(ctx, walletAddress)
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
