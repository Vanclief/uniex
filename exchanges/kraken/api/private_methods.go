package kraken

import (
	"reflect"
	"time"

	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
)

func (api *API) GetBalances() (*market.BalanceSnapshot, error) {
	const op = "kraken.GetBalances"

	snapshot := &market.BalanceSnapshot{
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
		snapshot.Balances = append(snapshot.Balances, *balance)
	}

	return snapshot, nil
}

func (api *API) GetPositions() ([]market.Position, error) {
	const op = "kraken.GetPositions"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) GetTrades() ([]market.Trade, error) {
	const op = "kraken.GetTrades"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

// GetOrder - Returns an existing order

// CreateOrder - Places a new order

// CancelOrder - Cancels an existing order

// GetFundingAddress - Retrieves or generates a new deposit addresses for an asset

// WithdrawAsset - Places a withdrawal request

// CancelWithdraw - Cancels an asset withdrawal
