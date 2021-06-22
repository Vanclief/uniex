package main

import (
	"time"

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
	// Public Endpoints
	// GetTicker(symbol string) ([]market.Tick)
	GetHistorical(symbol string, start, end time.Time) ([]market.Candle, error)
	// ListAssets() ([]market.Asset, error)
	GetOrderBook(symbol string) ([]interface{}, error)

	// Private Endpoints
	GetPositions() ([]market.Position, error)
	GetTrades() ([]market.Trade, error)
	// GetFundings() ([]market.Funding, error)
	// GetWithdraws() ([]market.Withdraw, error)
	// GetOrders() ([]market.Order, error)
	// CreateOrder() (market.Order, error)
	// CancelOrder() (bool, error)
}

// NewKraken returns a new Kraken.com exchange unified interface
func NewKraken() (*Exchange, error) {
	const op = "uniex.NewKraken"

	api, err := kraken.New("test", "test")
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return &Exchange{API: api}, nil
}
