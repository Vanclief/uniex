package kraken

import (
	"time"

	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/api"
)

func (api *API) GetTicker(pair *market.Pair) (*market.Ticker, error) {
	const op = "kraken.GetTicker"

	requestTime := time.Now()
	symbol := pair.Base.Symbol + pair.Quote.Symbol

	tickerMap, err := api.Client.GetTicker(symbol)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	ticker := &market.Ticker{}

	// We only want the first item
	for _, value := range tickerMap {

		// Because kraken doesn't return the timestamp, we make a estimate
		// based on when did we made the request
		ticker.Time = requestTime.Unix()
		ticker.Ask = value.Ask.Price
		ticker.Bid = value.Bid.Price

		break
	}

	return ticker, nil
}

func (api *API) GetCurrentCandle(pair *market.Pair, timeframe int) (*market.Candle, error) {
	const op = "kraken.GetCurrentCandle"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) GetHistoricalData(pair *market.Pair, start, end time.Time, interval int) ([]market.Candle, error) {
	const op = "kraken.GetHistoricalData"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) GetOrderBook(pair *market.Pair, options *api.GetOrderBookOptions) (*market.OrderBook, error) {
	const op = "kraken.GetOrderBook"

	symbol := pair.Base.Symbol + pair.Quote.Symbol

	limit := 100
	if options.Limit != 0 {
		limit = options.Limit
	}

	obMap, err := api.Client.GetOrderBook(symbol, int64(limit))
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	orderBook := &market.OrderBook{Time: time.Now().Unix()}

	// We only want the first item
	for _, value := range obMap {

		// Add the Asks
		var asks []market.OrderBookRow
		askAccumVolume := float64(0)

		for _, ask := range value.Asks {

			askAccumVolume = askAccumVolume + ask.Volume

			askRow := market.OrderBookRow{
				Price:       ask.Price,
				Volume:      ask.Volume,
				AccumVolume: askAccumVolume,
			}

			asks = append(asks, askRow)
		}
		orderBook.Asks = asks

		// Add the bids
		var bids []market.OrderBookRow
		bidsAccumVolume := float64(0)

		for _, bid := range value.Bids {

			bidsAccumVolume = bidsAccumVolume + bid.Volume

			bidRow := market.OrderBookRow{
				Price:       bid.Price,
				Volume:      bid.Volume,
				AccumVolume: bidsAccumVolume,
			}

			bids = append(bids, bidRow)
		}

		orderBook.Bids = bids
	}

	return orderBook, nil
}
