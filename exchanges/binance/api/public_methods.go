package binance

import (
	"strings"
	"time"

	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
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

func (api *API) GetHistoricalData(pair *market.Pair, start, end time.Time) ([]market.Candle, error) {
	const op = "binance.GetHistoricalData"

	pairString := strings.ToUpper(pair.Base.Symbol + pair.Quote.Symbol)

	candles, err := api.Client.FetchBinanceCandles(pairString, start, end)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return candles, nil
}

func (api *API) GetOrderBook(pair *market.Pair) (*market.OrderBook, error) {
	const op = "binance.GetOrderBook"
	pairString := strings.ToUpper(pair.Base.Symbol + pair.Quote.Symbol)
	orderBook, err := api.Client.FetchOrderBook(pairString)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}
	parsedOrderBook := &market.OrderBook{
		Time: float64(time.Now().Unix()),
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
