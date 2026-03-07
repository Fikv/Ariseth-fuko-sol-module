package client

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"ariseth-fuko-sol-module/internal/domain"
	"resty.dev/v3"
)

const (
	DefaultPreMatchBaseURL   = ""
	DefaultPreMatchUserAgent = "Mozilla/5.0 (compatible; ariseth-fuko-sol-module/1.0)"
	DefaultPreMatchPath      = "/sports"
	PreMatchBaseURLEnvVar    = "PREMATCH_SOURCE_BASE_URL"
	PreMatchPathEnvVar       = "PREMATCH_SOURCE_PATH"
	PreMatchUserAgentEnvVar  = "PREMATCH_SOURCE_USER_AGENT"
)

var nonSlugCharRegex = regexp.MustCompile(`[^a-z0-9-]+`)

type PreMatchScraperClient domain.PreMatchScraperClient

func NewPreMatchScraperClient(httpClient *resty.Client, baseURL string) *PreMatchScraperClient {
	if httpClient == nil {
		httpClient = getClient()
	}
	if strings.TrimSpace(baseURL) == "" {
		baseURL = os.Getenv(PreMatchBaseURLEnvVar)
	}
	if strings.TrimSpace(baseURL) == "" {
		baseURL = DefaultPreMatchBaseURL
	}
	userAgent := strings.TrimSpace(os.Getenv(PreMatchUserAgentEnvVar))
	if userAgent == "" {
		userAgent = DefaultPreMatchUserAgent
	}

	return &PreMatchScraperClient{
		HTTP:      httpClient,
		BaseURL:   strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		UserAgent: userAgent,
	}
}

func (c *PreMatchScraperClient) SetBaseURL(baseURL string) {
	c.BaseURL = strings.TrimRight(baseURL, "/")
}

func (c *PreMatchScraperClient) SetUserAgent(userAgent string) {
	if strings.TrimSpace(userAgent) == "" {
		c.UserAgent = DefaultPreMatchUserAgent
		return
	}

	c.UserAgent = userAgent
}

func (c *PreMatchScraperClient) GetPreMatchHTML(ctx context.Context, sportSlug string) (string, error) {
	url := c.BaseURL + c.buildSportsPath(sportSlug)

	resp, err := c.HTTP.R().
		SetContext(ctx).
		SetHeader("Accept", "text/html").
		SetHeader("User-Agent", c.UserAgent).
		Get(url)
	if err != nil {
		return "", err
	}
	if resp.IsError() {
		return "", fmt.Errorf("source request failed: %s %s", resp.Status(), url)
	}

	return resp.String(), nil
}

func (c *PreMatchScraperClient) buildSportsPath(sportSlug string) string {
	basePath := normalizePreMatchPath(os.Getenv(PreMatchPathEnvVar))
	if basePath == "" {
		basePath = DefaultPreMatchPath
	}

	normalized := strings.ToLower(strings.TrimSpace(sportSlug))
	normalized = strings.ReplaceAll(normalized, " ", "-")
	normalized = nonSlugCharRegex.ReplaceAllString(normalized, "")
	normalized = strings.Trim(normalized, "-")
	if normalized == "" {
		return basePath
	}

	return basePath + "/" + normalized
}

func normalizePreMatchPath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	path = strings.Trim(path, "/")
	if path == "" {
		return ""
	}

	return "/" + path
}
