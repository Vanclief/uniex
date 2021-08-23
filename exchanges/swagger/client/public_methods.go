package client

import (
  "fmt"
  "github.com/vanclief/ez"
  "github.com/vanclief/finmod/market"
  binanceClient "github.com/vanclief/uniex/exchanges/binance/client"
  "net/url"
  "time"
)

func (c *Client) GetHealth() error {
  const op = "swaggerClient.GetHealth"
  const URL = "https://mt-market-data-client-api-v1.agiliumtrade.agiliumtrade.ai/health"
  data := url.Values{}
  response := &SwaggerResponse{}

  err := c.httpRequest("GET", URL, data, response)
  if err != nil {
    return ez.Wrap(op, err)
  }
  return nil
}

func (c *Client) GetOHLCData(symbol, timeframe string, startTime, endTime time.Time) ([]market.Candle, error) {
  const op = "swaggerClient.GetOHLCData"
  URL := fmt.Sprintf(`https://mt-market-data-client-api-v1.agiliumtrade.agiliumtrade.ai/users/current/accounts/%s/historical-market-data/symbols/%s/timeframes/%s/candles`, c.AccountID, symbol, timeframe)

  var marketCandles []market.Candle

  arrayOfTimestamps := binanceClient.CreateArrayOfTimestamps(startTime, endTime)
  for _, v := range arrayOfTimestamps {
    tm := time.Unix(v.EndTime, 0)
    data := url.Values{
      "startTime": {tm.String()},
      "limit": {fmt.Sprint(1000)},
    }
    fmt.Println(v)
    response := &[]MetatraderCandle{}
    err := c.httpRequest("GET", URL, data, &response)
    if err != nil {
      return nil, ez.Wrap(op, err)
    }
    for _, vv := range *response {
      parsedTime, err := time.Parse(time.RFC3339, vv.Time)
      if err != nil {
        return nil, ez.Wrap(op, err)
      }
      marketCandles = append(marketCandles, market.Candle{
        Time:   parsedTime.Unix(),
        Open:   vv.Open,
        High:   vv.High,
        Low:    vv.Low,
        Close:  vv.Close,
        Volume: vv.Volume,
      })
    }
  }
  return marketCandles, nil
}