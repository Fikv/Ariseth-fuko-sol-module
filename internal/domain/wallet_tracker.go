package domain

import "time"

type WalletBalance struct {
	Lamports uint64  `json:"lamports"`
	SOL      float64 `json:"sol"`
}

type TokenBalance struct {
	Mint           string  `json:"mint"`
	Amount         string  `json:"amount"`
	Decimals       uint8   `json:"decimals"`
	UIAmount       float64 `json:"uiAmount"`
	UIAmountString string  `json:"uiAmountString"`
}

type WalletSnapshot struct {
	WalletAddress string         `json:"walletAddress"`
	Slot          uint64         `json:"slot"`
	Balance       WalletBalance  `json:"balance"`
	Tokens        []TokenBalance `json:"tokens"`
	UpdatedAt     time.Time      `json:"updatedAt"`
}
