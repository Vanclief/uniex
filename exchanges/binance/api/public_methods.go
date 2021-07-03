package binance

import (
	"strings"
	"time"

	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
)

func (api *API) GetTicker(pair *market.Pair) (*market.Ticker, error) {
	const op = "binance.GetTicker"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) GetHistoricalData(pair *market.Pair, start, end time.Time) ([]market.Candle, error) {
	const op = "binance.GetHistoricalData"

	pairString := pair.Base.Symbol + pair.Quote.Symbol
	pairString = strings.ToUpper(pairString)

	candles, err := api.Client.FetchBinanceCandles(pairString, start, end)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return candles, nil
}

func (api *API) GetOrderBook(pair *market.Pair) (*market.OrderBook, error) {
	const op = "binance.GetOrderBook"

	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) ListAssets() ([]market.Asset, error) {
	const op = "binance.ListAssets"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}
