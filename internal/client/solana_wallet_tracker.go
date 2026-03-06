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

type SolanaWalletTrackerClient domain.SolanaWalletTrackerClient

func NewSolanaWalletTrackerClient(httpClient *resty.Client) *SolanaWalletTrackerClient {
	if httpClient == nil {
		httpClient = getClient()
	}

	return &SolanaWalletTrackerClient{
		HTTP:       httpClient,
		RPCBaseURL: strings.TrimRight(DefaultSolanaRPCBaseURL, "/"),
	}
}

func (c *SolanaWalletTrackerClient) SetRPCBaseURL(baseURL string) {
	c.RPCBaseURL = strings.TrimRight(baseURL, "/")
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
	resp, err := c.HTTP.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&out).
		Post(c.RPCBaseURL)
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
	resp, err := c.HTTP.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&out).
		Post(c.RPCBaseURL)
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

type solanaRPCError = domain.SolanaRPCError
type solanaRPCResponse[T any] = domain.SolanaRPCResponse[T]
type solanaContext = domain.SolanaContext
type solanaBalanceResult = domain.SolanaBalanceResult
type solanaTokenAccountsResult = domain.SolanaTokenAccountsResult
type solanaTokenAccount = domain.SolanaTokenAccount
type solanaTokenAccountData = domain.SolanaTokenAccountData
type solanaParsedData = domain.SolanaParsedData
type solanaParsedInfo = domain.SolanaParsedInfo
type solanaTokenInfo = domain.SolanaTokenInfo
type solanaTokenAmount = domain.SolanaTokenAmount
