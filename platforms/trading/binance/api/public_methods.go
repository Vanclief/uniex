package api

import (
	"strings"
	"time"

	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/api"
)

func (api *API) GetTicker(pair *market.Pair) (*market.Ticker, error) {
	const op = "binance.GetTicker"
	pairString := strings.ToUpper(pair.Base.Symbol + pair.Quote.Symbol)
	ticker, err := api.Client.FetchTicker(pairString)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}
	return ticker, nil
}

func (api *API) GetCurrentCandle(pair *market.Pair, timeframe int) (*market.Candle, error) {
	const op = "binance.GetCurrentCandle"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) GetHistoricalData(pair *market.Pair, start, end time.Time, interval int) ([]market.Candle, error) {
	const op = "binance.GetHistoricalData"

	pairString := strings.ToUpper(pair.Base.Symbol + pair.Quote.Symbol)

	candles, err := api.Client.FetchBinanceCandles(pairString, start, end, interval)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return candles, nil
}

func (api *API) GetOrderBook(pair *market.Pair, options *api.GetOrderBookOptions) (*market.OrderBook, error) {
	const op = "binance.GetOrderBook"
	pairString := strings.ToUpper(pair.Base.Symbol + pair.Quote.Symbol)
	orderBook, err := api.Client.FetchOrderBook(pairString, options)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}
	parsedOrderBook := &market.OrderBook{
		Time: time.Now().Unix(),
	}
	accumAskVolume := 0.0
	accumBidVolume := 0.0
	for _, v := range orderBook.Asks {
		parsedOrderBook.Asks = append(parsedOrderBook.Asks, market.OrderBookRow{
			Price:       v.Price,
			Volume:      v.Quantity,
			AccumVolume: accumAskVolume,
		})
		accumAskVolume += v.Quantity
	}
	for _, v := range orderBook.Bids {
		parsedOrderBook.Bids = append(parsedOrderBook.Bids, market.OrderBookRow{
			Price:       v.Price,
			Volume:      v.Quantity,
			AccumVolume: accumBidVolume,
		})
		accumBidVolume += v.Quantity
	}

	return parsedOrderBook, nil
}

func (api *API) ListAssets() ([]market.Asset, error) {
	const op = "binance.ListAssets"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}
