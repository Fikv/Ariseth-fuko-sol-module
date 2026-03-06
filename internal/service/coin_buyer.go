package service

import (
	"context"
	"errors"
	"math"

	"ariseth-fuko-sol-module/internal/client"
	"ariseth-fuko-sol-module/internal/domain"
)

const (
	wrappedSOLMint   = "So11111111111111111111111111111111111111112"
	lamportsPerSOL   = 1_000_000_000
	defaultSlippage  = 100
	defaultTxVersion = "V0"
)

type CoinBuyerService domain.CoinBuyerService

func NewCoinBuyerService(raydiumClient *client.RaydiumClient) *CoinBuyerService {
	if raydiumClient == nil {
		raydiumClient = client.NewRaydiumClient(nil)
	}

	return &CoinBuyerService{
		RaydiumClient: raydiumClient,
	}
}

func (s *CoinBuyerService) BuyByCA(ctx context.Context, req domain.BuyByCARequest) (domain.BuyByCAResponse, error) {
	if req.ContractAddress == "" {
		return domain.BuyByCAResponse{}, errors.New("contractAddress is required")
	}
	if req.WalletAddress == "" {
		return domain.BuyByCAResponse{}, errors.New("walletAddress is required")
	}
	if req.SOLAmount <= 0 {
		return domain.BuyByCAResponse{}, errors.New("solAmount must be greater than zero")
	}

	slippage := req.SlippageBps
	if slippage <= 0 {
		slippage = defaultSlippage
	}
	txVersion := req.TxVersion
	if txVersion == "" {
		txVersion = defaultTxVersion
	}

	amountLamports := uint64(math.Round(req.SOLAmount * lamportsPerSOL))

	quote, err := s.RaydiumClient.ComputeSwapBaseIn(
		ctx,
		wrappedSOLMint,
		req.ContractAddress,
		amountLamports,
		slippage,
		txVersion,
	)
	if err != nil {
		return domain.BuyByCAResponse{}, err
	}

	swapResponse := quote
	if data, ok := quote["data"].(map[string]any); ok {
		swapResponse = data
	}

	transaction, err := s.RaydiumClient.BuildSwapTransactionBaseIn(
		ctx,
		req.WalletAddress,
		txVersion,
		req.ComputeUnitPriceMicroLamports,
		swapResponse,
		true,
		false,
		req.InputTokenAccount,
		req.OutputTokenAccount,
	)
	if err != nil {
		return domain.BuyByCAResponse{}, err
	}

	return domain.BuyByCAResponse{
		InputMint:           wrappedSOLMint,
		OutputMint:          req.ContractAddress,
		InputAmountLamports: amountLamports,
		Quote:               quote,
		Transaction:         transaction,
	}, nil
}
