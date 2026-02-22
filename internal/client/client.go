package client

import (
	"resty.dev/v3"
)

func getClient() *resty.Client {
	restyClient := resty.New()

	return restyClient
}
