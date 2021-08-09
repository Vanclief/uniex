package api

import (
	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/exchanges"
)

func (api *API) GetBalances() (*market.BalanceSnapshot, error) {
	const op = "bitso.GetBalances"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

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

func (api *API) GetPositions(request *exchanges.GetPositionsRequest) ([]market.Position, error) {
	const op = "bitso.GetPositions"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) GetTrades(request *exchanges.GetTradesRequest) ([]market.Trade, error) {
	const op = "bitso.GetPositions"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}
