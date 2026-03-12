package domain

type WalletTrackerService struct {
	Client WalletSnapshotGetter
}

type CoinBuyerService struct {
	RaydiumClient RaydiumSwapClient
}

type PreMatchScraperService struct {
	Client PreMatchPageGetter
}

type DataAnalysisService struct{}

type Wallet struct {
	Address string
	Value   float64
}
