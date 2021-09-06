package client

import (
	"fmt"
	"net/url"
	"time"

	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/utils"
)

func (c *Client) GetHealth() error {
	const op = "MetaAPI.Client.GetHealth"
	const URL = "https://mt-market-data-client-api-v1.agiliumtrade.agiliumtrade.ai/health"
	data := url.Values{}
	response := &MetaAPIResponse{}

	err := c.httpRequest("GET", URL, data, response)
	if err != nil {
		return ez.Wrap(op, err)
	}
	return nil
}

func (c *Client) ReadCandle(symbol string, interval int) (*market.Candle, error) {
	const op = "MetaAPI.Client.ReadCandle"

	if interval <= 0 {
		return nil, ez.New(op, ez.EINVALID, "Interval must be a positive number", nil)
	}

	timeframe := fmt.Sprintf(`%dm`, interval)

	URL := fmt.Sprintf(`https://mt-client-api-v1.agiliumtrade.agiliumtrade.ai/users/current/accounts/%s/symbols/%s/current-candles/%s`, c.AccountID, url.QueryEscape(symbol), timeframe)

	response := &MetaTraderCandle{}
	err := c.httpRequest("GET", URL, nil, &response)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	parsedTime, err := time.Parse(time.RFC3339, response.Time)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	candle := &market.Candle{
		Time:   parsedTime.Unix(),
		Open:   response.Open,
		High:   response.High,
		Low:    response.Low,
		Close:  response.Close,
		Volume: response.Volume,
	}

	return candle, nil
}

func (c *Client) GetOHLCData(symbol string, startTime, endTime time.Time, interval int) ([]market.Candle, error) {
	const op = "MetaAPI.Client.GetOHLCData"
	if interval <= 0 {
		return nil, ez.New(op, ez.EINVALID, "Interval must be a positive number", nil)
	}

	timeframe := fmt.Sprintf(`%dm`, interval)

	URL := fmt.Sprintf(`https://mt-market-data-client-api-v1.agiliumtrade.agiliumtrade.ai/users/current/accounts/%s/historical-market-data/symbols/%s/timeframes/%s/candles`, c.AccountID, url.QueryEscape(symbol), timeframe)

	var marketCandles []market.Candle

	limitCandles := utils.CalculateLimit(startTime, endTime, interval, 1000)

	arrayOfTimestamps := utils.CreateArrayOfTimestamps(startTime, endTime, interval)
	for _, v := range arrayOfTimestamps {
		tm := time.Unix(v.EndTime, 0)
		data := url.Values{
			"startTime": {tm.String()},
			"limit":     {fmt.Sprint(limitCandles)},
		}

		response := &[]MetaTraderCandle{}
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
