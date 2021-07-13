package api

import (
  "fmt"
  "github.com/vanclief/ez"
  "github.com/vanclief/finmod/market"
  "strings"
  "time"
)

func (api *API) GetTicker(pair *market.Pair) (*market.Ticker, error) {
  op := "bitso.GetTicker"
  requestTime := time.Now()
  symbol := strings.ToLower(pair.Base.Symbol + "_" + pair.Quote.Symbol)

  resp, err := api.Client.GetTicker(symbol)
  if err != nil {
    return nil, ez.Wrap(op, err)
  }
  fmt.Println(resp)
  
  ticker := &market.Ticker{
    Time:   requestTime.Unix(),
    Candle: &market.Candle{
      Time:   0,
      Open:   0,
      High:   0,
      Low:    0,
      Close:  0,
      Volume: 0,
    },
    Ask:    nil,
    Bid:    nil,
  }
  return ticker, nil
}