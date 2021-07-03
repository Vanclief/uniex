package binance

import (
	binanceclient "github.com/vanclief/uniex/exchanges/binance/client"
)

// API - Binance API
type API struct {
	Client *binanceclient.Client
}

func New(apiKey, secretKey string) (*API, error) {
	client := binanceclient.New(apiKey, secretKey)
	return &API{Client: client}, nil
}
