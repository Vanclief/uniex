package api

import (
	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"strconv"
	"strings"
	"time"
)

func (api *API) GetTicker(pair *market.Pair) (*market.Ticker, error) {
	op := "bitso.GetTicker"
	requestTime := time.Now()
	symbol := strings.ToLower(pair.Base.Symbol + "_" + pair.Quote.Symbol)

	bitsoTicker, err := api.Client.GetTicker(symbol)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	mTime, err := time.Parse(time.RFC3339, bitsoTicker.CreatedAt)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}
	mHigh, err := strconv.ParseFloat(bitsoTicker.High, 64)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}
	mLow, err := strconv.ParseFloat(bitsoTicker.Low, 64)
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
		Time: requestTime.Unix(),
		Candle: &market.Candle{
			Time:   mTime.Unix(),
			Open:   0,
			High:   mHigh,
			Low:    mLow,
			Close:  0,
			Volume: mVolume,
		},
		Ask: &market.OrderBookRow{
			Price:       obrPriceAsk,
			Volume:      0,
			AccumVolume: 0,
		},
		Bid: &market.OrderBookRow{
			Price:       obrPriceBid,
			Volume:      0,
			AccumVolume: 0,
		},
	}
	return ticker, nil
}

func (api *API) GetOrderBook(pair *market.Pair) (*market.OrderBook, error) {
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
