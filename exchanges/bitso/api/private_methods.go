package api

import (
	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/exchanges"
)

func (api *API) GetBalance() (*market.BalanceSnapshot, error) {
	const op = "bitso.GetBalance"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) GetAssets() (*market.AssetsSnashot, error) {
	const op = "bitso.GetAssets"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

// Orders
func (api *API) GetOrders(request *exchanges.GetOrdersRequest) ([]market.Order, error) {
	const op = "bitso.GetOrders"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) CreateOrder(orderRequest *market.OrderRequest) (*market.Order, error) {
	const op = "bitso.CreateOrder"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) CancelOrder(order *market.Order) (*market.Order, error) {
	const op = "bitso.CancelOrder"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

// Trades
func (api *API) GetTrades(request *exchanges.GetTradesRequest) ([]market.Trade, error) {
	const op = "bitso.GetPositions"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

// Positions
func (api *API) GetPositions(request *exchanges.GetPositionsRequest) ([]market.Position, error) {
	const op = "bitso.GetPositions"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) UpdatePosition(request *exchanges.UpdatePositionRequest) (*market.Position, error) {
	const op = "bitso.UpdatePosition"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) ClosePosition(request *exchanges.ClosePositionRequest) (*market.Position, error) {
	const op = "bitso.ClosePosition"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}
