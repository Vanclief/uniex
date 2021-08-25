package api

import (
	"strings"
	"time"

	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
)

func (api *API) GetTicker(pair *market.Pair) (*market.Ticker, error) {
	const op = "MetaAPI.GetHistoricalData"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) GetHistoricalData(pair *market.Pair, start, end time.Time) ([]market.Candle, error) {
	const op = "MetaAPI.GetHistoricalData"
	pairString := strings.ToUpper(pair.Base.Symbol + pair.Quote.Symbol)

	candles, err := api.Client.GetOHLCData(pairString, "1m", start, end)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}
	return candles, nil
}
