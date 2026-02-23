package client

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"resty.dev/v3"
)

const (

	DefaultAPIBaseURL  = "https://api-v3.raydium.io"
	DefaultSwapBaseURL = "https://transaction-v1.raydium.io"
)

const (
	endpointChainTime = "/main/chain-time"
	endpointRPCs      = "/main/rpcs"
	endpointInfo      = "/main/info"

	endpointTokenList  = "/mint/list"
	endpointTokenInfo  = "/mint/ids"
	endpointTokenPrice = "/mint/price"
	endpointJupTokens  = "https://tokens.jup.ag/tokens?tags=lst,community"

	endpointPoolList      = "/pools/info/list"
	endpointPoolByID      = "/pools/info/ids"
	endpointPoolByMint    = "/pools/info/mint"
	endpointPoolByLP      = "/pools/info/lps"
	endpointPoolKeysByID  = "/pools/key/ids"
	endpointPoolLiquidity = "/pools/line/liquidity"
	endpointPoolPosition  = "/pools/line/position"

	endpointFarmInfo = "/farms/info/ids"
	endpointFarmByLP = "/farms/info/lp"
	endpointFarmKeys = "/farms/key/ids"

	endpointCLMMConfig      = "/main/clmm-config"
	endpointCLMMPoolByID    = "/clmm/pool/info/ids"
	endpointCLMMPoolKeys    = "/clmm/pool/key/ids"
	endpointCLMMPoolList    = "/clmm/pools"
	endpointCLMMPoolAPR     = "/clmm/pool/apr"
	endpointCLMMPositionAPR = "/clmm/position/apr"
	endpointCLMMPrice       = "/clmm/price"
	endpointCLMMVolume      = "/clmm/volume"
	endpointCLMMTVL         = "/clmm/tvl"
	endpointCLMMPosition    = "/clmm/position/info"
	endpointCLMMVaultInfo   = "/clmm/vault/info"

	endpointCPMMConfig    = "/main/cpmm-config"
	endpointCPMMPoolByID  = "/ammV4/pools/info/ids"
	endpointCPMMPoolList  = "/ammV4/pools"
	endpointCPMMLPTokens  = "/ammV4/lp/tokens"
	endpointCPMMVolume    = "/ammV4/volumes"
	endpointCPMMStats     = "/ammV4/stats"
	endpointCPMMTVL       = "/ammV4/tvls"
	endpointCPMMAPR       = "/ammV4/aprs"
	endpointCPMMLiquidity = "/ammV4/liquiditys"

	endpointLaunchpadList = "/main/launchpad"
)

type RaydiumClient struct {
	apiBaseURL  string
	swapBaseURL string
	http        *resty.Client
}

type RaydiumResponse[T any] struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Data    T      `json:"data"`
}

func NewRaydiumClient(httpClient *resty.Client) *RaydiumClient {
	if httpClient == nil {
		httpClient = getClient()
	}

	return &RaydiumClient{
		apiBaseURL:  strings.TrimRight(DefaultAPIBaseURL, "/"),
		swapBaseURL: strings.TrimRight(DefaultSwapBaseURL, "/"),
		http:        httpClient,
	}
}

func (c *RaydiumClient) SetAPIBaseURL(baseURL string) {
	c.apiBaseURL = strings.TrimRight(baseURL, "/")
}

func (c *RaydiumClient) SetSwapBaseURL(baseURL string) {
	c.swapBaseURL = strings.TrimRight(baseURL, "/")
}

func (c *RaydiumClient) GetChainTime(ctx context.Context) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointChainTime, nil)
}

func (c *RaydiumClient) GetRPCs(ctx context.Context) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointRPCs, nil)
}

func (c *RaydiumClient) GetInfo(ctx context.Context) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointInfo, nil)
}

func (c *RaydiumClient) GetTokenList(ctx context.Context) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointTokenList, nil)
}

func (c *RaydiumClient) GetTokenInfo(ctx context.Context, mintIDs []string) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointTokenInfo, url.Values{"mints": {strings.Join(mintIDs, ",")}})
}

func (c *RaydiumClient) GetTokenPrice(ctx context.Context, mintIDs []string) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointTokenPrice, url.Values{"mints": {strings.Join(mintIDs, ",")}})
}

func (c *RaydiumClient) GetJupiterTokenList(ctx context.Context) (RaydiumResponse[[]map[string]any], error) {
	var out []map[string]any
	_, err := c.http.R().
		SetContext(ctx).
		SetResult(&out).
		Get(endpointJupTokens)
	if err != nil {
		return RaydiumResponse[[]map[string]any]{}, err
	}
	return RaydiumResponse[[]map[string]any]{Success: true, Data: out}, nil
}

func (c *RaydiumClient) GetPoolList(ctx context.Context, params url.Values) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointPoolList, params)
}

func (c *RaydiumClient) GetPoolByIDs(ctx context.Context, ids []string) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointPoolByID, url.Values{"ids": {strings.Join(ids, ",")}})
}

func (c *RaydiumClient) GetPoolByMints(ctx context.Context, mintA, mintB string, params url.Values) (RaydiumResponse[map[string]any], error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("mint1", mintA)
	params.Set("mint2", mintB)
	return c.getMap(ctx, endpointPoolByMint, params)
}

func (c *RaydiumClient) GetPoolByLPs(ctx context.Context, lps []string) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointPoolByLP, url.Values{"lps": {strings.Join(lps, ",")}})
}

func (c *RaydiumClient) GetPoolKeysByIDs(ctx context.Context, ids []string) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointPoolKeysByID, url.Values{"ids": {strings.Join(ids, ",")}})
}

func (c *RaydiumClient) GetPoolLiquidityLine(ctx context.Context, poolID string) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointPoolLiquidity, url.Values{"id": {poolID}})
}

func (c *RaydiumClient) GetPoolPositionLine(ctx context.Context, poolID string) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointPoolPosition, url.Values{"id": {poolID}})
}

func (c *RaydiumClient) GetFarmInfo(ctx context.Context, ids []string) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointFarmInfo, url.Values{"ids": {strings.Join(ids, ",")}})
}

func (c *RaydiumClient) GetFarmInfoByLP(ctx context.Context, lpMint string, page, pageSize int) (RaydiumResponse[map[string]any], error) {
	params := url.Values{"lp": {lpMint}}
	if page > 0 {
		params.Set("page", fmt.Sprintf("%d", page))
	}
	if pageSize > 0 {
		params.Set("pageSize", fmt.Sprintf("%d", pageSize))
	}
	return c.getMap(ctx, endpointFarmByLP, params)
}

func (c *RaydiumClient) GetFarmKeysByIDs(ctx context.Context, ids []string) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointFarmKeys, url.Values{"ids": {strings.Join(ids, ",")}})
}

func (c *RaydiumClient) GetCLMMConfig(ctx context.Context) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointCLMMConfig, nil)
}

func (c *RaydiumClient) GetCLMMPoolByIDs(ctx context.Context, ids []string) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointCLMMPoolByID, url.Values{"ids": {strings.Join(ids, ",")}})
}

func (c *RaydiumClient) GetCLMMPoolKeysByIDs(ctx context.Context, ids []string) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointCLMMPoolKeys, url.Values{"ids": {strings.Join(ids, ",")}})
}

func (c *RaydiumClient) GetCLMMPoolList(ctx context.Context, params url.Values) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointCLMMPoolList, params)
}

func (c *RaydiumClient) GetCLMMPoolAPR(ctx context.Context, ids []string) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointCLMMPoolAPR, url.Values{"ids": {strings.Join(ids, ",")}})
}

func (c *RaydiumClient) GetCLMMPositionAPR(ctx context.Context, params url.Values) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointCLMMPositionAPR, params)
}

func (c *RaydiumClient) GetCLMMPrice(ctx context.Context) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointCLMMPrice, nil)
}

func (c *RaydiumClient) GetCLMMVolume(ctx context.Context) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointCLMMVolume, nil)
}

func (c *RaydiumClient) GetCLMMTVL(ctx context.Context) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointCLMMTVL, nil)
}

func (c *RaydiumClient) GetCLMMPositionInfo(ctx context.Context, nftMint string) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointCLMMPosition, url.Values{"nftMint": {nftMint}})
}

func (c *RaydiumClient) GetCLMMVaultInfo(ctx context.Context, mintA, mintB string) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointCLMMVaultInfo, url.Values{"mint1": {mintA}, "mint2": {mintB}})
}

func (c *RaydiumClient) GetCPMMConfig(ctx context.Context) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointCPMMConfig, nil)
}

func (c *RaydiumClient) GetCPMMPools(ctx context.Context, params url.Values) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointCPMMPoolList, params)
}

func (c *RaydiumClient) GetCPMMPoolByIDs(ctx context.Context, ids []string) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointCPMMPoolByID, url.Values{"ids": {strings.Join(ids, ",")}})
}

func (c *RaydiumClient) GetCPMMLPTokens(ctx context.Context) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointCPMMLPTokens, nil)
}

func (c *RaydiumClient) GetCPMMVolume(ctx context.Context) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointCPMMVolume, nil)
}

func (c *RaydiumClient) GetCPMMStats(ctx context.Context) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointCPMMStats, nil)
}

func (c *RaydiumClient) GetCPMMTVL(ctx context.Context) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointCPMMTVL, nil)
}

func (c *RaydiumClient) GetCPMMAPR(ctx context.Context) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointCPMMAPR, nil)
}

func (c *RaydiumClient) GetCPMMLiquidity(ctx context.Context) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointCPMMLiquidity, nil)
}

func (c *RaydiumClient) GetLaunchpadList(ctx context.Context, params url.Values) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointLaunchpadList, params)
}

func (c *RaydiumClient) GetLaunchpadInfo(ctx context.Context, projectID string) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointLaunchpadList+"/"+projectID, nil)
}

func (c *RaydiumClient) GetLaunchpadJoined(ctx context.Context, projectID string) (RaydiumResponse[map[string]any], error) {
	return c.getMap(ctx, endpointLaunchpadList+"/"+projectID+"/mine", nil)
}

func (c *RaydiumClient) getMap(ctx context.Context, path string, params url.Values) (RaydiumResponse[map[string]any], error) {
	var out RaydiumResponse[map[string]any]
	if out.Data == nil {
		out.Data = map[string]any{}
	}

	url := c.apiBaseURL + path
	req := c.http.R().
		SetContext(ctx).
		SetResult(&out)
	if params != nil {
		req.SetQueryParamsFromValues(params)
	}

	resp, err := req.Get(url)
	if err != nil {
		return RaydiumResponse[map[string]any]{}, err
	}
	if resp.IsError() {
		return RaydiumResponse[map[string]any]{}, fmt.Errorf("raydium request failed: %s %s", resp.Status(), path)
	}

	if !out.Success {
		out.Success = true
	}
	return out, nil
}

func NewDefaultContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	return context.WithTimeout(context.Background(), timeout)
}



