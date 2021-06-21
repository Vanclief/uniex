package main

import (
	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	kraken "github.com/vanclief/uniex/exchanges/kraken/api"
)

type Exchange struct {
	API      ExchangeAPI
	TakerFee int
}

// ExchangeAPI represents an unified exchange API
type ExchangeAPI interface {
	GetPositions() ([]market.Position, error)
	// GetPositions() ([]*market.Position, error)
	// GetTick(pair string, date time.Time) (*market.Candle, error)
	// PlaceOrder(pair string, order *market.Order) error
}

// NewKraken returns
func NewKraken() (*Exchange, error) {
	const op = "uniex.NewKraken"

	api, err := kraken.New("test", "test")
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return &Exchange{API: api}, nil
}
