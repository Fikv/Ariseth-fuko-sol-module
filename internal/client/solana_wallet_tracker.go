package client

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"ariseth-fuko-sol-module/internal/domain"
	"resty.dev/v3"
)

const (
	DefaultSolanaRPCBaseURL = "https://api.mainnet-beta.solana.com"
	lamportsPerSOL          = 1_000_000_000
	tokenProgramID          = ""
)

type SolanaWalletTrackerClient struct {
	http       *resty.Client
	rpcBaseURL string
}

func NewSolanaWalletTrackerClient(httpClient *resty.Client) *SolanaWalletTrackerClient {
	if httpClient == nil {
		httpClient = getClient()
	}

	return &SolanaWalletTrackerClient{
		http:       httpClient,
		rpcBaseURL: strings.TrimRight(DefaultSolanaRPCBaseURL, "/"),
	}
}

func (c *SolanaWalletTrackerClient) SetRPCBaseURL(baseURL string) {
	c.rpcBaseURL = strings.TrimRight(baseURL, "/")
}

func (c *SolanaWalletTrackerClient) GetWalletSnapshot(ctx context.Context, walletAddress string) (domain.WalletSnapshot, error) {
	solLamports, balanceSlot, err := c.GetSOLBalance(ctx, walletAddress)
	if err != nil {
		return domain.WalletSnapshot{}, err
	}

	tokens, tokenSlot, err := c.GetTokenBalances(ctx, walletAddress)
	if err != nil {
		return domain.WalletSnapshot{}, err
	}

	slot := balanceSlot
	if tokenSlot > slot {
		slot = tokenSlot
	}

	return domain.WalletSnapshot{
		WalletAddress: walletAddress,
		Slot:          slot,
		Balance: domain.WalletBalance{
			Lamports: solLamports,
			SOL:      float64(solLamports) / lamportsPerSOL,
		},
		Tokens:    tokens,
		UpdatedAt: time.Now().UTC(),
	}, nil
}

func (c *SolanaWalletTrackerClient) GetSOLBalance(ctx context.Context, walletAddress string) (uint64, uint64, error) {
	req := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getBalance",
		"params": []any{
			walletAddress,
			map[string]any{"commitment": "confirmed"},
		},
	}

	var out solanaRPCResponse[solanaBalanceResult]
	resp, err := c.http.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&out).
		Post(c.rpcBaseURL)
	if err != nil {
		return 0, 0, err
	}
	if resp.IsError() {
		return 0, 0, fmt.Errorf("solana request failed: %s getBalance", resp.Status())
	}
	if out.Error != nil {
		return 0, 0, fmt.Errorf("solana rpc error (%d): %s", out.Error.Code, out.Error.Message)
	}

	return out.Result.Value, out.Result.Context.Slot, nil
}

func (c *SolanaWalletTrackerClient) GetTokenBalances(ctx context.Context, walletAddress string) ([]domain.TokenBalance, uint64, error) {
	req := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getTokenAccountsByOwner",
		"params": []any{
			walletAddress,
			map[string]any{"programId": tokenProgramID},
			map[string]any{
				"encoding":   "jsonParsed",
				"commitment": "confirmed",
			},
		},
	}

	var out solanaRPCResponse[solanaTokenAccountsResult]
	resp, err := c.http.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&out).
		Post(c.rpcBaseURL)
	if err != nil {
		return nil, 0, err
	}
	if resp.IsError() {
		return nil, 0, fmt.Errorf("solana request failed: %s getTokenAccountsByOwner", resp.Status())
	}
	if out.Error != nil {
		return nil, 0, fmt.Errorf("solana rpc error (%d): %s", out.Error.Code, out.Error.Message)
	}

	tokenBalances := make([]domain.TokenBalance, 0, len(out.Result.Value))
	for _, account := range out.Result.Value {
		parsed := account.Account.Data.Parsed.Info
		uiAmount := parsed.TokenAmount.UIAmount
		if uiAmount == 0 && parsed.TokenAmount.UIAmountString != "" {
			parsedFloat, parseErr := strconv.ParseFloat(parsed.TokenAmount.UIAmountString, 64)
			if parseErr == nil {
				uiAmount = parsedFloat
			}
		}

		tokenBalances = append(tokenBalances, domain.TokenBalance{
			Mint:           parsed.Mint,
			Amount:         parsed.TokenAmount.Amount,
			Decimals:       parsed.TokenAmount.Decimals,
			UIAmount:       uiAmount,
			UIAmountString: parsed.TokenAmount.UIAmountString,
		})
	}

	return tokenBalances, out.Result.Context.Slot, nil
}

type solanaRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type solanaRPCResponse[T any] struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  T               `json:"result"`
	Error   *solanaRPCError `json:"error"`
}

type solanaContext struct {
	Slot uint64 `json:"slot"`
}

type solanaBalanceResult struct {
	Context solanaContext `json:"context"`
	Value   uint64        `json:"value"`
}

type solanaTokenAccountsResult struct {
	Context solanaContext        `json:"context"`
	Value   []solanaTokenAccount `json:"value"`
}

type solanaTokenAccount struct {
	Account solanaTokenAccountData `json:"account"`
}

type solanaTokenAccountData struct {
	Data solanaParsedData `json:"data"`
}

type solanaParsedData struct {
	Parsed solanaParsedInfo `json:"parsed"`
}

type solanaParsedInfo struct {
	Info solanaTokenInfo `json:"info"`
}

type solanaTokenInfo struct {
	Mint        string            `json:"mint"`
	TokenAmount solanaTokenAmount `json:"tokenAmount"`
}

type solanaTokenAmount struct {
	Amount         string  `json:"amount"`
	Decimals       uint8   `json:"decimals"`
	UIAmount       float64 `json:"uiAmount"`
	UIAmountString string  `json:"uiAmountString"`
}
