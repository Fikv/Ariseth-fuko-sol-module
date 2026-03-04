package domain

type BuyByCARequest struct {
	ContractAddress                 string  `json:"contractAddress"`
	WalletAddress                   string  `json:"walletAddress"`
	SOLAmount                       float64 `json:"solAmount"`
	SlippageBps                     int     `json:"slippageBps"`
	TxVersion                       string  `json:"txVersion"`
	ComputeUnitPriceMicroLamports   string  `json:"computeUnitPriceMicroLamports"`
	InputTokenAccount               string  `json:"inputTokenAccount,omitempty"`
	OutputTokenAccount              string  `json:"outputTokenAccount,omitempty"`
}

type BuyByCAResponse struct {
	InputMint          string         `json:"inputMint"`
	OutputMint         string         `json:"outputMint"`
	InputAmountLamports uint64        `json:"inputAmountLamports"`
	Quote              map[string]any `json:"quote"`
	Transaction        map[string]any `json:"transaction"`
}

