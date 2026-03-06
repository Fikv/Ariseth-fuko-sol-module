package domain

type WalletTrackerService struct {
	Client WalletSnapshotGetter
}

type CoinBuyerService struct {
	RaydiumClient RaydiumSwapClient
}

type Wallet struct {
	Address string
	Value   float64
}
