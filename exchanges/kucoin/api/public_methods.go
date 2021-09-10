package kucoin

import (
	"strconv"
	"time"

	"github.com/Kucoin/kucoin-go-sdk"
	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/exchanges"
)

func (api *API) GetTicker(pair *market.Pair) (*market.Ticker, error) {
	const op = "kucoin.GetTicker"

	// Convert the pair to the expected symbol by Kucoin
	symbol := pair.Base.Symbol + "-" + pair.Quote.Symbol

	// Make the request
	response, err := api.Client.TickerLevel1(symbol)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	// TODO
	// Add error catching: If I send a symbol that doesn't exist, I don't get an error

	// Parse the request
	kucoinTicker := kucoin.TickerLevel1Model{}
	err = response.ReadData(&kucoinTicker)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	bestAsk, err := strconv.ParseFloat(kucoinTicker.BestAsk, 64)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	bestAskSize, err := strconv.ParseFloat(kucoinTicker.BestAskSize, 64)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	bestBid, err := strconv.ParseFloat(kucoinTicker.BestBid, 64)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	bestBidSize, err := strconv.ParseFloat(kucoinTicker.BestBidSize, 64)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	ticker := &market.Ticker{
		Time: kucoinTicker.Time,
		Ask: &market.OrderBookRow{
			Price:       bestAsk,
			Volume:      bestAskSize,
			AccumVolume: bestAskSize,
		},
		Bid: &market.OrderBookRow{
			Price:       bestBid,
			Volume:      bestBidSize,
			AccumVolume: bestBidSize,
		},
	}

	return ticker, nil
}

func (api *API) GetCurrentCandle(pair *market.Pair, timeframe int) (*market.Candle, error) {
	const op = "kraken.GetCurrentCandle"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) GetHistoricalData(pair *market.Pair, start, end time.Time, interval int) ([]market.Candle, error) {
	const op = "kucoin.GetHistoricalData"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) GetOrderBook(pair *market.Pair, options *exchanges.GetOrderBookOptions) (*market.OrderBook, error) {
	const op = "kucoin.GetOrderBook"

	symbol := pair.Base.Symbol + "-" + pair.Quote.Symbol

	response, err := api.Client.AggregatedFullOrderBookV3(symbol)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	kucoinOrderBook := kucoin.FullOrderBookModel{}
	err = response.ReadData(&kucoinOrderBook)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	orderBook := &market.OrderBook{}

	// Add the Asks
	asks := []market.OrderBookRow{}
	askAccumVolume := float64(0)

	for _, ask := range kucoinOrderBook.Asks {

		askPrice, err := strconv.ParseFloat(ask[0], 64)
		if err != nil {
			return nil, ez.Wrap(op, err)
		}

		askVolume, err := strconv.ParseFloat(ask[1], 64)
		if err != nil {
			return nil, ez.Wrap(op, err)
		}

		askAccumVolume = askAccumVolume + askVolume

		askRow := market.OrderBookRow{
			Price:       askPrice,
			Volume:      askVolume,
			AccumVolume: askAccumVolume,
		}

		asks = append(asks, askRow)
	}
	orderBook.Asks = asks

	// Add the bids
	bids := []market.OrderBookRow{}
	bidsAccumVolume := float64(0)

	for _, bid := range kucoinOrderBook.Bids {

		bidPrice, err := strconv.ParseFloat(bid[0], 64)
		if err != nil {
			return nil, ez.Wrap(op, err)
		}

		bidVolume, err := strconv.ParseFloat(bid[1], 64)
		if err != nil {
			return nil, ez.Wrap(op, err)
		}

		bidsAccumVolume = bidsAccumVolume + bidVolume

		bidRow := market.OrderBookRow{
			Price:       bidPrice,
			Volume:      bidVolume,
			AccumVolume: bidsAccumVolume,
		}

		bids = append(bids, bidRow)
	}

	orderBook.Bids = bids

	return orderBook, nil
}

func (api *API) ListAssets() ([]market.Asset, error) {
	const op = "kucoin.ListAssets"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}
