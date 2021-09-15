package exchanges

import (
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/oracles"
)

// Exchange - An exchange or market data API
type Exchange struct {
	Name             string
	API              ExchangeAPI
	MakerFee         float64
	TakerFee         float64
	ManagedPositions bool
	HedgingMode      bool
}

// ExchangeAPI - Interface for an unified exchange API
type ExchangeAPI interface {
	// Public Endpoints
	oracles.DataOracleAPI
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
