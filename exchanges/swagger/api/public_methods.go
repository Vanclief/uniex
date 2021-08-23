package api

import (
  "github.com/vanclief/ez"
  "github.com/vanclief/finmod/market"
  "strings"
  "time"
)

func (api *API) GetHistoricalData(pair *market.Pair, start, end time.Time) ([]market.Candle, error) {
  const op = "swagger.GetHistoricalData"
  pairString := strings.ToUpper(pair.Base.Symbol + pair.Quote.Symbol)

  candles, err := api.Client.GetOHLCData(pairString, "1m", start,  end)
  if err != nil {
    return nil, ez.Wrap(op, err)
  }
  return candles, nil
}