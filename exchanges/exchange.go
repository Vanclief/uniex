package exchanges

import (
	"time"

	"github.com/vanclief/finmod/market"
)

// Exchange - An exchange or market data API
type Exchange struct {
	Name string
	API  ExchangeAPI
}

// ExchangeAPI - Interface for an unified exchange API
type ExchangeAPI interface {
	// Public Endpoints
	GetTicker(pair *market.Pair) (*market.Ticker, error)
	GetHistoricalData(pair *market.Pair, start, end time.Time) ([]market.Candle, error)
	GetOrderBook(pair *market.Pair) (*market.OrderBook, error)
	// ListAssets() ([]market.Asset, error)

	// Private Endpoints
	GetBalances() (*market.BalanceSnapshot, error)
	// GetPositions() ([]market.Position, error)
	// GetTrades() ([]market.Trade, error)
	// GetFundings() ([]market.Funding, error)
	// GetWithdraws() ([]market.Withdraw, error)
	// GetOrders() ([]market.Order, error)
	CreateOrder(orderRequest *market.OrderRequest) (*market.Order, error)
	CancelOrder(order *market.Order) (*market.Order, error)
	// Withdraw() (market.Asset, error)
}
