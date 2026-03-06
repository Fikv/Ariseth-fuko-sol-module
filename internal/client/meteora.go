package client

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"ariseth-fuko-sol-module/internal/domain"
	"resty.dev/v3"
)

const (
	DefaultMeteoraDLMMBaseURL         = "https://dlmm.datapi.meteora.ag"
	DefaultMeteoraDAMMV2BaseURL       = "https://damm-v2.datapi.meteora.ag"
	DefaultMeteoraDAMMV1BaseURL       = "https://damm-api.meteora.ag"
	DefaultMeteoraStake2EarnBaseURL   = "https://stake-for-fee-api.meteora.ag"
	DefaultMeteoraDynamicVaultBaseURL = "https://merv2-api.meteora.ag"
)

const (
	meteoraPathPools           = "/pools"
	meteoraPathPoolGroups      = "/pools/groups"
	meteoraPathPool            = "/pools/%s"
	meteoraPathPoolOHLCV       = "/pools/%s/ohlcv"
	meteoraPathPoolVolumeHist  = "/pools/%s/volume/history"
	meteoraPathProtocolMetrics = "/stats/protocol_metrics"

	meteoraPathAlphaVault        = "/alpha-vault"
	meteoraPathAlphaVaultConfigs = "/alpha-vault-configs"
	meteoraPathFarm              = "/farm"
	meteoraPathPoolConfigs       = "/pool-configs"
	meteoraPathPoolsMetrics      = "/pools-metrics"
	meteoraPathPoolsSearch       = "/pools/search"
	meteoraPathPoolsByVaultLP    = "/get_pools_by_a_vault_lp"
	meteoraPathFeeConfigByID     = "/fee-config/%s"

	meteoraPathAnalyticsAll = "/analytics/all"
	meteoraPathVaultAll     = "/vault/all"
	meteoraPathVaultFilter  = "/vault/filter"
	meteoraPathVaultByID    = "/vault/%s"

	meteoraPathVaultInfo      = "/vault_info"
	meteoraPathVaultAddresses = "/vault_addresses"
	meteoraPathVaultState     = "/vault_state/%s"
	meteoraPathAPYState       = "/apy_state/%s"
	meteoraPathAPYFilter      = "/apy_filter/%s/%d/%d"
	meteoraPathVirtualPrice   = "/virtual_price/%s/%s"
)

type MeteoraClient domain.MeteoraClient

func NewMeteoraClient(httpClient *resty.Client) *MeteoraClient {
	if httpClient == nil {
		httpClient = getClient()
	}

	return &MeteoraClient{
		HTTP:                httpClient,
		DLMMBaseURL:         strings.TrimRight(DefaultMeteoraDLMMBaseURL, "/"),
		DAMMV2BaseURL:       strings.TrimRight(DefaultMeteoraDAMMV2BaseURL, "/"),
		DAMMV1BaseURL:       strings.TrimRight(DefaultMeteoraDAMMV1BaseURL, "/"),
		Stake2EarnBaseURL:   strings.TrimRight(DefaultMeteoraStake2EarnBaseURL, "/"),
		DynamicVaultBaseURL: strings.TrimRight(DefaultMeteoraDynamicVaultBaseURL, "/"),
	}
}

type MeteoraResponse[T any] = domain.MeteoraResponse[T]

func (c *MeteoraClient) SetDLMMBaseURL(baseURL string) {
	c.DLMMBaseURL = strings.TrimRight(baseURL, "/")
}

func (c *MeteoraClient) SetDAMMV2BaseURL(baseURL string) {
	c.DAMMV2BaseURL = strings.TrimRight(baseURL, "/")
}

func (c *MeteoraClient) SetDAMMV1BaseURL(baseURL string) {
	c.DAMMV1BaseURL = strings.TrimRight(baseURL, "/")
}

func (c *MeteoraClient) SetStake2EarnBaseURL(baseURL string) {
	c.Stake2EarnBaseURL = strings.TrimRight(baseURL, "/")
}

func (c *MeteoraClient) SetDynamicVaultBaseURL(baseURL string) {
	c.DynamicVaultBaseURL = strings.TrimRight(baseURL, "/")
}

func (c *MeteoraClient) GetDLMMPools(ctx context.Context, params url.Values) (map[string]any, error) {
	return c.getMap(ctx, c.DLMMBaseURL, meteoraPathPools, params)
}

func (c *MeteoraClient) GetDLMMPoolGroups(ctx context.Context, params url.Values) (map[string]any, error) {
	return c.getMap(ctx, c.DLMMBaseURL, meteoraPathPoolGroups, params)
}

func (c *MeteoraClient) GetDLMMPoolGroup(ctx context.Context, lexicalOrderMints string, params url.Values) (map[string]any, error) {
	return c.getMap(ctx, c.DLMMBaseURL, meteoraPathPoolGroups+"/"+lexicalOrderMints, params)
}

func (c *MeteoraClient) GetDLMMPool(ctx context.Context, poolAddress string, params url.Values) (map[string]any, error) {
	return c.getMap(ctx, c.DLMMBaseURL, fmt.Sprintf(meteoraPathPool, poolAddress), params)
}

func (c *MeteoraClient) GetDLMMPoolOHLCV(ctx context.Context, poolAddress string, params url.Values) (map[string]any, error) {
	return c.getMap(ctx, c.DLMMBaseURL, fmt.Sprintf(meteoraPathPoolOHLCV, poolAddress), params)
}

func (c *MeteoraClient) GetDLMMPoolVolumeHistory(ctx context.Context, poolAddress string, params url.Values) (map[string]any, error) {
	return c.getMap(ctx, c.DLMMBaseURL, fmt.Sprintf(meteoraPathPoolVolumeHist, poolAddress), params)
}

func (c *MeteoraClient) GetDLMMProtocolMetrics(ctx context.Context) (map[string]any, error) {
	return c.getMap(ctx, c.DLMMBaseURL, meteoraPathProtocolMetrics, nil)
}

func (c *MeteoraClient) GetDAMMV2Pools(ctx context.Context, params url.Values) (map[string]any, error) {
	return c.getMap(ctx, c.DAMMV2BaseURL, meteoraPathPools, params)
}

func (c *MeteoraClient) GetDAMMV2PoolGroups(ctx context.Context, params url.Values) (map[string]any, error) {
	return c.getMap(ctx, c.DAMMV2BaseURL, meteoraPathPoolGroups, params)
}

func (c *MeteoraClient) GetDAMMV2PoolGroup(ctx context.Context, lexicalOrderMints string, params url.Values) (map[string]any, error) {
	return c.getMap(ctx, c.DAMMV2BaseURL, meteoraPathPoolGroups+"/"+lexicalOrderMints, params)
}

func (c *MeteoraClient) GetDAMMV2Pool(ctx context.Context, poolAddress string, params url.Values) (map[string]any, error) {
	return c.getMap(ctx, c.DAMMV2BaseURL, fmt.Sprintf(meteoraPathPool, poolAddress), params)
}

func (c *MeteoraClient) GetDAMMV2PoolOHLCV(ctx context.Context, poolAddress string, params url.Values) (map[string]any, error) {
	return c.getMap(ctx, c.DAMMV2BaseURL, fmt.Sprintf(meteoraPathPoolOHLCV, poolAddress), params)
}

func (c *MeteoraClient) GetDAMMV2PoolVolumeHistory(ctx context.Context, poolAddress string, params url.Values) (map[string]any, error) {
	return c.getMap(ctx, c.DAMMV2BaseURL, fmt.Sprintf(meteoraPathPoolVolumeHist, poolAddress), params)
}

func (c *MeteoraClient) GetDAMMV2ProtocolMetrics(ctx context.Context) (map[string]any, error) {
	return c.getMap(ctx, c.DAMMV2BaseURL, meteoraPathProtocolMetrics, nil)
}

func (c *MeteoraClient) GetDAMMV1AlphaVaults(ctx context.Context, params url.Values) ([]map[string]any, error) {
	return c.getList(ctx, c.DAMMV1BaseURL, meteoraPathAlphaVault, params)
}

func (c *MeteoraClient) GetDAMMV1AlphaVaultConfigs(ctx context.Context, params url.Values) ([]map[string]any, error) {
	return c.getList(ctx, c.DAMMV1BaseURL, meteoraPathAlphaVaultConfigs, params)
}

func (c *MeteoraClient) GetDAMMV1Farms(ctx context.Context, params url.Values) ([]map[string]any, error) {
	return c.getList(ctx, c.DAMMV1BaseURL, meteoraPathFarm, params)
}

func (c *MeteoraClient) GetDAMMV1PoolConfigs(ctx context.Context, params url.Values) ([]map[string]any, error) {
	return c.getList(ctx, c.DAMMV1BaseURL, meteoraPathPoolConfigs, params)
}

func (c *MeteoraClient) GetDAMMV1Pools(ctx context.Context, params url.Values) ([]map[string]any, error) {
	return c.getList(ctx, c.DAMMV1BaseURL, meteoraPathPools, params)
}

func (c *MeteoraClient) GetDAMMV1PoolsMetrics(ctx context.Context, params url.Values) ([]map[string]any, error) {
	return c.getList(ctx, c.DAMMV1BaseURL, meteoraPathPoolsMetrics, params)
}

func (c *MeteoraClient) SearchDAMMV1Pools(ctx context.Context, params url.Values) ([]map[string]any, error) {
	return c.getList(ctx, c.DAMMV1BaseURL, meteoraPathPoolsSearch, params)
}

func (c *MeteoraClient) GetDAMMV1PoolsByVaultLP(ctx context.Context, tokenAVaultLP string) (any, error) {
	body := map[string]string{"token_a_vault_lp": tokenAVaultLP}
	return c.postAny(ctx, c.DAMMV1BaseURL, meteoraPathPoolsByVaultLP, body)
}

func (c *MeteoraClient) GetDAMMV1FeeConfig(ctx context.Context, configAddress string) (map[string]any, error) {
	return c.getMap(ctx, c.DAMMV1BaseURL, fmt.Sprintf(meteoraPathFeeConfigByID, configAddress), nil)
}

func (c *MeteoraClient) GetStake2EarnAnalytics(ctx context.Context) (map[string]any, error) {
	return c.getMap(ctx, c.Stake2EarnBaseURL, meteoraPathAnalyticsAll, nil)
}

func (c *MeteoraClient) GetStake2EarnVaults(ctx context.Context, params url.Values) (map[string]any, error) {
	return c.getMap(ctx, c.Stake2EarnBaseURL, meteoraPathVaultAll, params)
}

func (c *MeteoraClient) FilterStake2EarnVaults(ctx context.Context, params url.Values) (map[string]any, error) {
	return c.getMap(ctx, c.Stake2EarnBaseURL, meteoraPathVaultFilter, params)
}

func (c *MeteoraClient) GetStake2EarnVault(ctx context.Context, vaultAddress string) (map[string]any, error) {
	return c.getMap(ctx, c.Stake2EarnBaseURL, fmt.Sprintf(meteoraPathVaultByID, vaultAddress), nil)
}

func (c *MeteoraClient) GetDynamicVaultInfo(ctx context.Context) ([]map[string]any, error) {
	return c.getList(ctx, c.DynamicVaultBaseURL, meteoraPathVaultInfo, nil)
}

func (c *MeteoraClient) GetDynamicVaultAddresses(ctx context.Context) ([]map[string]any, error) {
	return c.getList(ctx, c.DynamicVaultBaseURL, meteoraPathVaultAddresses, nil)
}

func (c *MeteoraClient) GetDynamicVaultState(ctx context.Context, tokenMint string) (map[string]any, error) {
	return c.getMap(ctx, c.DynamicVaultBaseURL, fmt.Sprintf(meteoraPathVaultState, tokenMint), nil)
}

func (c *MeteoraClient) GetDynamicVaultAPYState(ctx context.Context, tokenMint string) (map[string]any, error) {
	return c.getMap(ctx, c.DynamicVaultBaseURL, fmt.Sprintf(meteoraPathAPYState, tokenMint), nil)
}

func (c *MeteoraClient) GetDynamicVaultAPYByTimeRange(
	ctx context.Context,
	tokenMint string,
	startTimestamp int64,
	endTimestamp int64,
) (map[string]any, error) {
	return c.getMap(
		ctx,
		c.DynamicVaultBaseURL,
		fmt.Sprintf(meteoraPathAPYFilter, tokenMint, startTimestamp, endTimestamp),
		nil,
	)
}

func (c *MeteoraClient) GetDynamicVaultVirtualPrice(ctx context.Context, tokenMint, strategy string) (map[string]any, error) {
	return c.getMap(ctx, c.DynamicVaultBaseURL, fmt.Sprintf(meteoraPathVirtualPrice, tokenMint, strategy), nil)
}

func (c *MeteoraClient) getMap(
	ctx context.Context,
	baseURL string,
	path string,
	params url.Values,
) (map[string]any, error) {
	var out map[string]any
	if out == nil {
		out = map[string]any{}
	}

	resp, err := c.request(ctx, baseURL, path, params, nil, &out)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("meteora request failed: %s %s", resp.Status(), path)
	}

	return out, nil
}

func (c *MeteoraClient) getList(
	ctx context.Context,
	baseURL string,
	path string,
	params url.Values,
) ([]map[string]any, error) {
	var out []map[string]any
	resp, err := c.request(ctx, baseURL, path, params, nil, &out)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("meteora request failed: %s %s", resp.Status(), path)
	}
	return out, nil
}

func (c *MeteoraClient) postAny(
	ctx context.Context,
	baseURL string,
	path string,
	body any,
) (any, error) {
	var out any
	resp, err := c.request(ctx, baseURL, path, nil, body, &out)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("meteora request failed: %s %s", resp.Status(), path)
	}
	return out, nil
}

func (c *MeteoraClient) request(
	ctx context.Context,
	baseURL string,
	path string,
	params url.Values,
	body any,
	result any,
) (*resty.Response, error) {
	req := c.HTTP.R().
		SetContext(ctx).
		SetResult(result)
	if params != nil {
		req.SetQueryParamsFromValues(params)
	}
	if body != nil {
		req.SetBody(body)
	}

	fullURL := strings.TrimRight(baseURL, "/") + path
	if body != nil {
		return req.Post(fullURL)
	}
	return req.Get(fullURL)
}
