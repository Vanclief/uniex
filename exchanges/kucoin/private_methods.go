package kucoin

import (
	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
)

func (api *API) GetBalances() (*market.BalanceSnapshot, error) {
	const op = "kucoin.GetBalances"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) GetOrders(orders ...*market.Order) ([]market.Order, error) {
	const op = "kucoin.GetOrders"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) CreateOrder(orderRequest *market.OrderRequest) (*market.Order, error) {
	const op = "kucoin.CreateOrder"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) CancelOrder(order *market.Order) (*market.Order, error) {
	const op = "kucoin.CancelOrder"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}
