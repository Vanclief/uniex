package main

import (
	"github.com/vanclief/ez"
	kraken "github.com/vanclief/uniex/exchanges/kraken/api"
)

// Kraken - Returns a new Kraken.com exchange unified interface
func Kraken(apiKey, apiSecret string) (*Exchange, error) {
	const op = "uniex.Kraken"

	api, err := kraken.New(apiKey, apiSecret)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return &Exchange{API: api}, nil
}
