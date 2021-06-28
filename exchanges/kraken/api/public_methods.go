package kraken

import (
	"time"

	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
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

		ticker.Candle = &market.Candle{
			Time:   requestTime.Add(-24 * time.Hour).Unix(), // The candle we get is from the past 24 hours
			Open:   value.OpeningPrice,
			High:   value.High.Price,
			Low:    value.Low.Price,
			Close:  value.Close.Price,
			Volume: value.Volume.Price,
		}

		ticker.Ask = &market.OrderBookRow{
			Price:       value.Ask.Price,
			Volume:      value.Ask.Volume,
			TotalVolume: value.Ask.WholeLotVolume,
		}

		ticker.Bid = &market.OrderBookRow{
			Price:       value.Bid.Price,
			Volume:      value.Bid.Volume,
			TotalVolume: value.Bid.WholeLotVolume,
		}
		break
	}

	return ticker, nil
}

func (api *API) GetHistoricalData(pair *market.Pair, start, end time.Time) ([]market.Candle, error) {
	const op = "kraken.GetHistoricalData"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) GetOrderBook(pair *market.Pair) (*market.OrderBook, error) {
	const op = "kraken.GetOrderBook"
	const maxDepth = 500

	symbol := pair.Base.Symbol + pair.Quote.Symbol

	obMap, err := api.Client.GetOrderBook(symbol, maxDepth)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	orderBook := &market.OrderBook{}

	// We only want the first item
	for _, value := range obMap {

		// Add the Asks
		asks := []market.OrderBookRow{}
		askTotalVolume := float64(0)

		for _, ask := range value.Asks {

			askTotalVolume = askTotalVolume + ask.Volume

			askRow := market.OrderBookRow{
				Price:       ask.Price,
				Volume:      ask.Volume,
				TotalVolume: askTotalVolume,
			}

			asks = append(asks, askRow)
		}
		orderBook.Asks = asks

		// Add the bids
		bids := []market.OrderBookRow{}
		bidsTotalVolume := float64(0)

		for _, bid := range value.Bids {

			bidsTotalVolume = bidsTotalVolume + bid.Volume

			bidRow := market.OrderBookRow{
				Price:       bid.Price,
				Volume:      bid.Volume,
				TotalVolume: bidsTotalVolume,
			}

			bids = append(bids, bidRow)
		}

		orderBook.Bids = bids
	}

	return orderBook, nil
}
