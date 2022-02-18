package api

import (
	"time"

	"github.com/vanclief/finmod/market"
)

// DataAPI - Unified interface for data APIs
type DataAPI interface {
	GetTicker(pair *market.Pair) (*market.Ticker, error)
	GetCurrentCandle(pair *market.Pair, timeframe int) (*market.Candle, error)
	GetHistoricalData(pair *market.Pair, start, end time.Time, interval int) ([]market.Candle, error)
	// ListAssets() ([]market.Asset, error)
}

// TradingAPI - Unified interface for broker and exchange APIs
type TradingAPI interface {
	// Public Endpoints
	GetOrderBook(pair *market.Pair, options *GetOrderBookOptions) (*market.OrderBook, error)

	// Private Endpoints
	GetBalance() (*market.BalanceSnapshot, error)
	GetAssets() (*market.AssetsSnashot, error)

	// Orders
	GetOrders(request *GetOrdersRequest) ([]market.Order, error)
	CreateOrder(request *market.OrderRequest) (*market.Order, error)
	UpdateOrder(order *market.Order, request *UpdateOrderRequest) (*market.Order, error)
	CancelOrder(order *market.Order) (string, error)

	// Trades
	GetTrades(request *GetTradesRequest) ([]market.Trade, error)

	// Positions
	GetPositions(request *GetPositionsRequest) ([]market.Position, error)
	UpdatePosition(position *market.Position, request *UpdatePositionRequest) (*market.Position, error)
	ClosePosition(position *market.Position) (string, error)

	// Fundings & Withdraws
	// GetFundings() ([]market.Funding, error)
	// GetWithdraws() ([]market.Withdraw, error)
	// Withdraw() (market.Asset, error)
}
