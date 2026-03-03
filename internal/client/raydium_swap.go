package client

import (
	"context"
	"fmt"
	"strconv"
)

const (
	endpointSwapComputeBaseIn = "/compute/swap-base-in"
	endpointSwapTxBaseIn      = "/transaction/swap-base-in"
)

func (c *RaydiumClient) ComputeSwapBaseIn(
	ctx context.Context,
	inputMint string,
	outputMint string,
	amount uint64,
	slippageBps int,
	txVersion string,
) (map[string]any, error) {
	if slippageBps <= 0 {
		slippageBps = 100
	}
	if txVersion == "" {
		txVersion = "V0"
	}

	var out map[string]any
	if out == nil {
		out = map[string]any{}
	}

	resp, err := c.http.R().
		SetContext(ctx).
		SetQueryParam("inputMint", inputMint).
		SetQueryParam("outputMint", outputMint).
		SetQueryParam("amount", strconv.FormatUint(amount, 10)).
		SetQueryParam("slippageBps", strconv.Itoa(slippageBps)).
		SetQueryParam("txVersion", txVersion).
		SetResult(&out).
		Get(c.swapBaseURL + endpointSwapComputeBaseIn)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("raydium swap compute failed: %s", resp.Status())
	}

	return out, nil
}

func (c *RaydiumClient) BuildSwapTransactionBaseIn(
	ctx context.Context,
	walletAddress string,
	txVersion string,
	computeUnitPriceMicroLamports string,
	swapResponse map[string]any,
	wrapSol bool,
	unwrapSol bool,
	inputTokenAccount string,
	outputTokenAccount string,
) (map[string]any, error) {
	if txVersion == "" {
		txVersion = "V0"
	}

	body := map[string]any{
		"wallet":                         walletAddress,
		"txVersion":                      txVersion,
		"computeUnitPriceMicroLamports":  computeUnitPriceMicroLamports,
		"swapResponse":                   swapResponse,
		"wrapSol":                        wrapSol,
		"unwrapSol":                      unwrapSol,
	}
	if inputTokenAccount != "" {
		body["inputAccount"] = inputTokenAccount
	}
	if outputTokenAccount != "" {
		body["outputAccount"] = outputTokenAccount
	}

	var out map[string]any
	if out == nil {
		out = map[string]any{}
	}

	resp, err := c.http.R().
		SetContext(ctx).
		SetBody(body).
		SetResult(&out).
		Post(c.swapBaseURL + endpointSwapTxBaseIn)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("raydium swap transaction build failed: %s", resp.Status())
	}

	return out, nil
}

