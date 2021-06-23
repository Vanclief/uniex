package kraken

import (
	"fmt"
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

func (api *API) GetTicker(pair market.Pair) (*market.Ticker, error) {
	const op = "kraken.GetTicker"

	t, err := api.Client.GetTicker(pair.Base.Symbol)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	fmt.Println("t", t)

	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) GetHistoricalData(pair market.Pair, start, end time.Time) ([]market.Candle, error) {
	const op = "kraken.GetHistoricalData"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) GetOrderBook(pair market.Pair) (*market.OrderBook, error) {
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
