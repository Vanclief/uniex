package exchanges

import (
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/oracles"
)

// Exchange - An exchange or market data API
type Exchange struct {
	Name     string
	API      ExchangeAPI
	MakerFee float64
	TakerFee float64
}

// ExchangeAPI - Interface for an unified exchange API
type ExchangeAPI interface {
	// Public Endpoints
	oracles.DataOracleAPI
	GetOrderBook(pair *market.Pair, options *GetOrderBookOptions) (*market.OrderBook, error)

	// Private Endpoints
	GetBalances() (*market.BalanceSnapshot, error)

	// Orders
	GetOrders(request *GetOrdersRequest) ([]market.Order, error)
	CreateOrder(orderRequest *market.OrderRequest) (*market.Order, error)
	CancelOrder(order *market.Order) (*market.Order, error)

	// Trades
	GetTrades(request *GetTradesRequest) ([]market.Trade, error)

	// Positions
	GetPositions(request *GetPositionsRequest) ([]market.Position, error)
	// GetTrades() ([]market.Trade, error)
	// GetFundings() ([]market.Funding, error)
	// GetWithdraws() ([]market.Withdraw, error)
	// Withdraw() (market.Asset, error)
}