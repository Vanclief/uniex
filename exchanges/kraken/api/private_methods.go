package kraken

import (
	"reflect"
	"time"

	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/exchanges"
)

func (api *API) GetBalance() (*market.BalanceSnapshot, error) {
	const op = "kraken.GetBalance"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) GetAssets() (*market.AssetsSnashot, error) {
	const op = "kraken.GetAssets"

	snapshot := &market.AssetsSnashot{
		Time: float64(time.Now().Unix()),
	}

	krakenBalance, err := api.Client.GetAccountBalance()
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	v := reflect.ValueOf(*krakenBalance)

	for i := 0; i < v.NumField(); i++ {

		// Translate the asset from string to actual asset
		assetName := v.Type().Field(i).Name
		asset, err := TranslateAsset(assetName)
		if err != nil {
			continue
		}

		amount := v.Field(i).Interface().(float64)
		if amount <= 0 {
			continue
		}

		balance := &market.Balance{Asset: asset, Amount: amount}
		snapshot.Assets = append(snapshot.Assets, *balance)
	}

	return snapshot, nil
}

// GetOrders - Returns existing orders with their status
func (api *API) GetOrders(request *exchanges.GetOrdersRequest) ([]market.Order, error) {
	const op = "kraken.Orders"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

// CreateOrder - Places a new order
func (api *API) CreateOrder(orderRequest *market.OrderRequest) (*market.Order, error) {
	const op = "kraken.CreateOrder"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) UpdateOrder(order *market.Order, request *exchanges.UpdateOrderRequest) (*market.Order, error) {
	const op = "kraken.UpdateOrder"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

// CancelOrder - Cancels an existing order
func (api *API) CancelOrder(order *market.Order) (string, error) {
	const op = "kraken.CancelOrder"
	return "", ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) GetTrades(request *exchanges.GetTradesRequest) ([]market.Trade, error) {
	const op = "kraken.GetPositions"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) GetPositions(request *exchanges.GetPositionsRequest) ([]market.Position, error) {
	const op = "kraken.GetPositions"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) UpdatePosition(position *market.Position, request *exchanges.UpdatePositionRequest) (*market.Position, error) {
	const op = "kraken.UpdatePosition"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) ClosePosition(position *market.Position) (string, error) {
	const op = "kraken.ClosePosition"
	return "", ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

// GetFundingAddress - Retrieves or generates a new deposit addresses for an asset

// WithdrawAsset - Places a withdrawal request

// CancelWithdraw - Cancels an asset withdrawal
