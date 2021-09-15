package api

import (
	"time"

	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/exchanges"
)

func (api *API) GetTicker(pair *market.Pair) (*market.Ticker, error) {
	const op = "MetaAPI.GetTicker"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) GetCurrentCandle(pair *market.Pair, interval int) (*market.Candle, error) {
	const op = "MetaAPI.GetCurrentCandle"

	pairString := pair.Symbol("")

	candle, err := api.Client.ReadCandle(pairString, interval)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return candle, nil
}

func (api *API) GetHistoricalData(pair *market.Pair, start, end time.Time, interval int) ([]market.Candle, error) {
	const op = "MetaAPI.GetHistoricalData"
	pairString := pair.Symbol("")

	candles, err := api.Client.GetOHLCData(pairString, start, end, interval)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return candles, nil
}

func (api *API) GetOrderBook(pair *market.Pair, options *exchanges.GetOrderBookOptions) (*market.OrderBook, error) {
	const op = "MetaAPI.GetTicker"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}
