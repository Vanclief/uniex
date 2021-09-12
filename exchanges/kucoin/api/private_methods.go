package kucoin

import (
	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/exchanges"
)

func (api *API) GetBalance() (*market.BalanceSnapshot, error) {
	const op = "kucoin.GetBalance"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) GetAssets() (*market.AssetsSnashot, error) {
	const op = "kucoin.GetAssets"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

// Orders
func (api *API) GetOrders(request *exchanges.GetOrdersRequest) ([]market.Order, error) {
	const op = "kucoin.GetOrders"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) CreateOrder(orderRequest *market.OrderRequest) (*market.Order, error) {
	const op = "kucoin.CreateOrder"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) CancelOrder(order *market.Order) (string, error) {
	const op = "kucoin.CancelOrder"
	return "", ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

// Trades
func (api *API) GetTrades(request *exchanges.GetTradesRequest) ([]market.Trade, error) {
	const op = "kucoin.GetPositions"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

// Positions
func (api *API) GetPositions(request *exchanges.GetPositionsRequest) ([]market.Position, error) {
	const op = "kucoin.GetPositions"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) UpdatePosition(position *market.Position, request *exchanges.UpdatePositionRequest) (*market.Position, error) {
	const op = "kucoin.UpdatePosition"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) ClosePosition(position *market.Position) (string, error) {
	const op = "kucoin.ClosePosition"
	return "", ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}
