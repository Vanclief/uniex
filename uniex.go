package main

import (
	"time"

	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	kraken "github.com/vanclief/uniex/exchanges/kraken/api"
)

// Exchange - An exchange or market data API
type Exchange struct {
	API ExchangeAPI
	// TODO?
}

// ExchangeAPI - Interface for an unified exchange API
type ExchangeAPI interface {
	// Public Endpoints
	GetTicker(asset market.Asset) ([]market.Ticker, error)
	GetHistoricalData(symbol string, start, end time.Time) ([]market.Candle, error)
	GetOrderBook(asset market.Asset) (*market.OrderBook, error)
	// ListAssets() ([]market.Asset, error)

	// Private Endpoints
	// GetPositions() ([]market.Position, error)
	// GetTrades() ([]market.Trade, error)
	// GetFundings() ([]market.Funding, error)
	// GetWithdraws() ([]market.Withdraw, error)
	// GetOrders() ([]market.Order, error)
	// CreateOrder() (market.Order, error)
	// CancelOrder() (bool, error)
	// Withdraw() (market.Asset, error)
}

// Kraken - Returns a new Kraken.com exchange unified interface
func Kraken() (*Exchange, error) {
	const op = "uniex.Kraken"

	api, err := kraken.New("test", "test")
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return &Exchange{API: api}, nil
}
