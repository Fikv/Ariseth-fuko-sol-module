package domain

import (
	"context"
	"time"

	"resty.dev/v3"
)

type RaydiumClient struct {
	APIBaseURL  string
	SwapBaseURL string
	HTTP        *resty.Client
}

type RaydiumResponse[T any] struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Data    T      `json:"data"`
}

type MeteoraClient struct {
	HTTP *resty.Client

	DLMMBaseURL         string
	DAMMV2BaseURL       string
	DAMMV1BaseURL       string
	Stake2EarnBaseURL   string
	DynamicVaultBaseURL string
}

type MeteoraResponse[T any] struct {
	Data T `json:"data"`
}

type SolanaWalletTrackerClient struct {
	HTTP       *resty.Client
	RPCBaseURL string
}

type PreMatchScraperClient struct {
	HTTP      *resty.Client
	BaseURL   string
	UserAgent string
}

type SolanaRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SolanaRPCResponse[T any] struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  T               `json:"result"`
	Error   *SolanaRPCError `json:"error"`
}

type SolanaContext struct {
	Slot uint64 `json:"slot"`
}

type SolanaBalanceResult struct {
	Context SolanaContext `json:"context"`
	Value   uint64        `json:"value"`
}

type SolanaTokenAccountsResult struct {
	Context SolanaContext        `json:"context"`
	Value   []SolanaTokenAccount `json:"value"`
}

type SolanaTokenAccount struct {
	Account SolanaTokenAccountData `json:"account"`
}

type SolanaTokenAccountData struct {
	Data SolanaParsedData `json:"data"`
}

type SolanaParsedData struct {
	Parsed SolanaParsedInfo `json:"parsed"`
}

type SolanaParsedInfo struct {
	Info SolanaTokenInfo `json:"info"`
}

type SolanaTokenInfo struct {
	Mint        string            `json:"mint"`
	TokenAmount SolanaTokenAmount `json:"tokenAmount"`
}

type SolanaTokenAmount struct {
	Amount         string  `json:"amount"`
	Decimals       uint8   `json:"decimals"`
	UIAmount       float64 `json:"uiAmount"`
	UIAmountString string  `json:"uiAmountString"`
}

type WalletSnapshotGetter interface {
	GetWalletSnapshot(ctx context.Context, walletAddress string) (WalletSnapshot, error)
}

type RaydiumSwapClient interface {
	ComputeSwapBaseIn(
		ctx context.Context,
		inputMint string,
		outputMint string,
		amount uint64,
		slippageBps int,
		txVersion string,
	) (map[string]any, error)
	BuildSwapTransactionBaseIn(
		ctx context.Context,
		walletAddress string,
		txVersion string,
		computeUnitPriceMicroLamports string,
		swapResponse map[string]any,
		wrapSol bool,
		unwrapSol bool,
		inputTokenAccount string,
		outputTokenAccount string,
	) (map[string]any, error)
}

type PreMatchPageGetter interface {
	GetPreMatchHTML(ctx context.Context, sportSlug string) (string, error)
}

type PreMatchSnapshot struct {
	SportSlug string    `json:"sportSlug"`
	HTML      string    `json:"html"`
	FetchedAt time.Time `json:"fetchedAt"`
}
