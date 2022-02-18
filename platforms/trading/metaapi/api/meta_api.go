package api

import (
	metaAPIClient "github.com/vanclief/uniex/platforms/trading/metaapi/api/client"
)

type API struct {
	Client *metaAPIClient.Client
}

func New(accountID, token string) (*API, error) {
	client := metaAPIClient.New(accountID, token)
	return &API{Client: client}, nil
}
