package domain

import "time"

type HoldMode string

const (
	HoldModeTimed            HoldMode = "timed"
	HoldModeUntilSourceSells HoldMode = "source_sells"
	HoldModeManual           HoldMode = "manual"
)

type TradeAction string

const (
	TradeActionBuy  TradeAction = "buy"
	TradeActionSell TradeAction = "sell"
)

type CopyPreset struct {
	CopyAmountSOL   float64  `json:"copyAmountSOL"`
	MaxBuyPriceSOL  float64  `json:"maxBuyPriceSOL"`
	HoldMode        HoldMode `json:"holdMode"`
	HoldSeconds     int      `json:"holdSeconds"`
	TakeProfitPct   float64  `json:"takeProfitPct"`
	StopLossPct     float64  `json:"stopLossPct"`
	SlippageBps     int      `json:"slippageBps"`
	AllowAutoMirror bool     `json:"allowAutoMirror"`
}

type ExecutionSettings struct {
	DefaultBuyAmountSOL float64  `json:"defaultBuyAmountSOL"`
	MinWalletReserveSOL float64  `json:"minWalletReserveSOL"`
	Network             string   `json:"network"`
	HoldMode            HoldMode `json:"holdMode"`
	HoldSeconds         int      `json:"holdSeconds"`
	TakeProfitPct       float64  `json:"takeProfitPct"`
	StopLossPct         float64  `json:"stopLossPct"`
	SlippageBps         int      `json:"slippageBps"`
}

type WalletStats struct {
	MirroredBuys    int     `json:"mirroredBuys"`
	ClosedPositions int     `json:"closedPositions"`
	SkippedBuys     int     `json:"skippedBuys"`
	RealizedPnLSOL  float64 `json:"realizedPnLSOL"`
}

type Position struct {
	ID              string     `json:"id"`
	TokenSymbol     string     `json:"tokenSymbol"`
	TokenAddress    string     `json:"tokenAddress"`
	EntryPriceSOL   float64    `json:"entryPriceSOL"`
	CurrentPriceSOL float64    `json:"currentPriceSOL"`
	Quantity        float64    `json:"quantity"`
	InvestedSOL     float64    `json:"investedSOL"`
	HoldMode        HoldMode   `json:"holdMode"`
	OpenedAt        time.Time  `json:"openedAt"`
	TargetExitAt    *time.Time `json:"targetExitAt,omitempty"`
	ExitPending     bool       `json:"exitPending"`
	ExitReason      string     `json:"exitReason,omitempty"`
	ExitRequestedAt *time.Time `json:"exitRequestedAt,omitempty"`
	ExitAttempts    int        `json:"exitAttempts"`
	ExecutedBuy     bool       `json:"executedBuy"`
}

type LogEntry struct {
	ID        string    `json:"id"`
	Kind      string    `json:"kind"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type Notification struct {
	ID        string    `json:"id"`
	Level     string    `json:"level"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type Wallet struct {
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	Address       string      `json:"address"`
	Chain         string      `json:"chain"`
	Enabled       bool        `json:"enabled"`
	LastAction    string      `json:"lastAction"`
	UpdatedAt     time.Time   `json:"updatedAt"`
	Preset        CopyPreset  `json:"preset"`
	Stats         WalletStats `json:"stats"`
	OpenPositions []Position  `json:"openPositions"`
	Logs          []LogEntry  `json:"logs"`
}

type DashboardSummary struct {
	WalletCount    int       `json:"walletCount"`
	EnabledWallets int       `json:"enabledWallets"`
	OpenPositions  int       `json:"openPositions"`
	RealizedPnLSOL float64   `json:"realizedPnLSOL"`
	LastUpdated    time.Time `json:"lastUpdated"`
}

type DashboardState struct {
	Summary       DashboardSummary `json:"summary"`
	Wallets       []Wallet         `json:"wallets"`
	Notifications []Notification   `json:"notifications"`
}

type AddWalletRequest struct {
	Name    string     `json:"name"`
	Address string     `json:"address"`
	Enabled bool       `json:"enabled"`
	Preset  CopyPreset `json:"preset"`
}

type UpdatePresetRequest struct {
	Preset CopyPreset `json:"preset"`
}

type InjectTradeRequest struct {
	WalletID     string      `json:"walletId"`
	Action       TradeAction `json:"action"`
	TokenSymbol  string      `json:"tokenSymbol"`
	TokenAddress string      `json:"tokenAddress"`
	PriceSOL     float64     `json:"priceSOL"`
	Quantity     float64     `json:"quantity"`
	Notes        string      `json:"notes"`
}

type ObservedTrade struct {
	Action       TradeAction
	TokenSymbol  string
	TokenAddress string
	PriceSOL     float64
	Quantity     float64
	Notes        string
	ObservedAt   time.Time
}

type LeaderboardImportRequest struct {
	Limit   int        `json:"limit"`
	Enabled bool       `json:"enabled"`
	Preset  CopyPreset `json:"preset"`
}

type WalletImportEntry struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Source  string `json:"source"`
}

type LeaderboardImportResult struct {
	Imported int      `json:"imported"`
	Skipped  int      `json:"skipped"`
	Wallets  []Wallet `json:"wallets"`
}

type RealtimeTrade struct {
	TraderName    string
	TraderAddress string
	Action        TradeAction
	TokenSymbol   string
	TokenAddress  string
	PriceSOL      float64
	Quantity      float64
	Signature     string
	Fingerprint   string
}

type AutoCopyOrder struct {
	WalletID       string
	WalletName     string
	WalletAddress  string
	PositionID     string
	TokenSymbol    string
	TokenAddress   string
	AmountSOL      float64
	SlippageBps    int
	SourcePriceSOL float64
}

type AutoSellOrder struct {
	WalletID      string
	WalletName    string
	WalletAddress string
	PositionID    string
	TokenSymbol   string
	TokenAddress  string
	SlippageBps   int
	Reason        string
	RetryCount    int
	NextAttemptAt time.Time
}

type WatchedWalletRef struct {
	ID      string
	Name    string
	Address string
	Enabled bool
}
