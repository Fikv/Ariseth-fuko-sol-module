package client

import "resty.dev/v3"

const (
	getTokensUrl = "https://api.dexscreener.com/token-profiles/latest/v1"
)

func GetTokens() *resty.Response {
	response, err := getClient().R().Get(getTokensUrl)

	if err != nil {
		print(err)
	}

	print(response)
	return response
	
}