package api

import (
	"strconv"
	"strings"
	"time"

	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/api"
)

func (api *API) GetTicker(pair *market.Pair) (*market.Ticker, error) {
	op := "bitso.GetTicker"
	requestTime := time.Now()
	symbol := strings.ToLower(pair.Base.Symbol + "_" + pair.Quote.Symbol)

	bitsoTicker, err := api.Client.GetTicker(symbol)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	mVolume, err := strconv.ParseFloat(bitsoTicker.Volume, 64)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}
	obrPriceBid, err := strconv.ParseFloat(bitsoTicker.Bid, 64)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}
	obrPriceAsk, err := strconv.ParseFloat(bitsoTicker.Ask, 64)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	ticker := &market.Ticker{
		Time:   requestTime.Unix(),
		Ask:    obrPriceAsk,
		Bid:    obrPriceBid,
		Volume: mVolume,
	}
	return ticker, nil
}

func (api *API) GetCurrentCandle(pair *market.Pair, timeframe int) (*market.Candle, error) {
	const op = "bitso.GetCurrentCandle"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) GetOrderBook(pair *market.Pair, options *api.GetOrderBookOptions) (*market.OrderBook, error) {
	op := "bitso.GetOrderBook"
	symbol := strings.ToLower(pair.Base.Symbol + "_" + pair.Quote.Symbol)
	orderBook, err := api.Client.GetOrderBook(symbol)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	parsedOrderBook := &market.OrderBook{
		Time: time.Now().Unix(),
	}
	accumVol := 0.0
	for _, v := range orderBook.Asks {
		thisPrice, err := strconv.ParseFloat(v.Price, 64)
		if err != nil {
			return nil, ez.Wrap(op, err)
		}
		thisAmount, err := strconv.ParseFloat(v.Amount, 64)
		if err != nil {
			return nil, ez.Wrap(op, err)
		}
		parsedOrderBook.Asks = append(parsedOrderBook.Asks, market.OrderBookRow{
			Price:       thisPrice,
			Volume:      thisAmount,
			AccumVolume: accumVol,
		})
		accumVol += thisAmount
	}
	accumVol = 0.0
	for _, v := range orderBook.Bids {
		thisPrice, err := strconv.ParseFloat(v.Price, 64)
		if err != nil {
			return nil, ez.Wrap(op, err)
		}
		thisAmount, err := strconv.ParseFloat(v.Amount, 64)
		if err != nil {
			return nil, ez.Wrap(op, err)
		}
		parsedOrderBook.Bids = append(parsedOrderBook.Bids, market.OrderBookRow{
			Price:       thisPrice,
			Volume:      thisAmount,
			AccumVolume: accumVol,
		})
		accumVol += thisAmount
	}
	return parsedOrderBook, nil
}

func (api *API) GetHistoricalData(pair *market.Pair, start, end time.Time, interval int) ([]market.Candle, error) {
	const op = "bitso.GetHistoricalData"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}
