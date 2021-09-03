package api

import (
	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/exchanges"
)

func (api *API) GetBalances() (*market.BalanceSnapshot, error) {
	const op = "MetaAPI.GetBalances"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

// GetOrders - Returns existing orders with their status
func (api *API) GetOrders(request *exchanges.GetOrdersRequest) ([]market.Order, error) {
	const op = "MetaAPI.Orders"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

// CreateOrder - Places a new order
func (api *API) CreateOrder(orderRequest *market.OrderRequest) (*market.Order, error) {
	const op = "MetaAPI.CreateOrder"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

// CancelOrder - Cancels an existing order
func (api *API) CancelOrder(order *market.Order) (*market.Order, error) {
	const op = "MetaAPI.CancelOrder"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) GetPositions(request *exchanges.GetPositionsRequest) ([]market.Position, error) {
	const op = "MetaAPI.GetPositions"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) GetTrades(request *exchanges.GetTradesRequest) ([]market.Trade, error) {
	const op = "MetaAPI.GetPositions"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

// GetFundingAddress - Retrieves or generates a new deposit addresses for an asset

// WithdrawAsset - Places a withdrawal request

// CancelWithdraw - Cancels an asset withdrawal
