package kraken

import (
	"time"

	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	krakenclient "github.com/vanclief/uniex/exchanges/kraken/client"
)

// Kraken represents the interface
type API struct {
	Client *krakenclient.Client
}

func New(publicKey string, secretKey string) (*API, error) {
	return nil, nil
}

// func (api *API)

func (api *API) GetHistorical(symbol string, start, end time.Time) ([]market.Candle, error) {
	const op = "kraken.GetHistorical"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) GetOrderBook(symbol string) ([]interface{}, error) {
	const op = "kraken.GetOrderBook"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) GetPositions() ([]market.Position, error) {
	const op = "kraken.GetPositions"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) GetTrades() ([]market.Trade, error) {
	const op = "kraken.GetTrades"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}
