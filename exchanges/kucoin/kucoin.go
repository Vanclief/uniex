package kucoin

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Kucoin/kucoin-go-sdk"
	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
)

// API - Kucoin API
type API struct {
	Client *kucoin.ApiService
}

func New(apiKey, apiSecret, passphrase string) (*API, error) {

	client := kucoin.NewApiService(
		kucoin.ApiKeyOption(apiKey),
		kucoin.ApiSecretOption(apiSecret),
		kucoin.ApiPassPhraseOption(passphrase),
		kucoin.ApiKeyVersionOption(kucoin.ApiKeyVersionV2),
	)

	return &API{Client: client}, nil
}

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
			TotalVolume: bestAskSize,
		},
		Bid: &market.OrderBookRow{
			Price:       bestBid,
			Volume:      bestBidSize,
			TotalVolume: bestBidSize,
		},
	}

	return ticker, nil
}

func (api *API) GetHistoricalData(pair *market.Pair, start, end time.Time) ([]market.Candle, error) {
	const op = "kucoin.GetHistoricalData"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) GetOrderBook(pair *market.Pair) (*market.OrderBook, error) {
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

	fmt.Println("kucoin", kucoinOrderBook)

	return nil, nil

	orderBook := &market.OrderBook{}

	// // We only want the first item
	// for _, value := range obMap {

	// 	// Add the Asks
	// 	asks := []market.OrderBookRow{}
	// 	askTotalVolume := float64(0)

	// 	for _, ask := range value.Asks {

	// 		askTotalVolume = askTotalVolume + ask.Volume

	// 		askRow := market.OrderBookRow{
	// 			Price:       ask.Price,
	// 			Volume:      ask.Volume,
	// 			TotalVolume: askTotalVolume,
	// 		}

	// 		asks = append(asks, askRow)
	// 	}
	// 	orderBook.Asks = asks

	// 	// Add the bids
	// 	bids := []market.OrderBookRow{}
	// 	bidsTotalVolume := float64(0)

	// 	for _, bid := range value.Asks {

	// 		bidsTotalVolume = bidsTotalVolume + bid.Volume

	// 		bidRow := market.OrderBookRow{
	// 			Price:       bid.Price,
	// 			Volume:      bid.Volume,
	// 			TotalVolume: bidsTotalVolume,
	// 		}

	// 		bids = append(bids, bidRow)
	// 	}

	// 	orderBook.Bids = bids
	// }

	return orderBook, nil
}

func (api *API) GetPositions() ([]market.Position, error) {
	const op = "kucoin.GetPositions"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) GetTrades() ([]market.Trade, error) {
	const op = "kucoin.GetTrades"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}
