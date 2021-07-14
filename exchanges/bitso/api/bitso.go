package api

import (
	bitsoClient "github.com/vanclief/uniex/exchanges/bitso/client"
)

type API struct {
	Client *bitsoClient.Client
}

func New(apiKey string) (*API, error) {
	client := bitsoClient.New(apiKey)
	return &API{Client: client}, nil
}
