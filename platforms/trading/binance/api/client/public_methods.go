package binanceclient

import (
	"time"

	goBinance "github.com/binance-exchange/go-binance"
	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/api"
	"github.com/vanclief/uniex/utils"
	// 	"github.com/go-kit/kit/log"
	// 	"github.com/go-kit/kit/log/level"
	// 	"github.com/vanclief/ez"
)

// FetchBinanceCandles takes the start and end dates, and the crypto pair as strings, and returns the Binance Candle
// data for every minute between start and end for a crypto pair
func (b Client) FetchBinanceCandles(pair string, start, end time.Time, interval int) ([]market.Candle, error) {
	op := "binance.FetchBinanceCandles"

	var marketCandles []market.Candle

	arrayOfTimestamps := utils.CreateArrayOfTimestamps(start, end, interval)
	for _, v := range arrayOfTimestamps {
		kl, err := b.service.Klines(goBinance.KlinesRequest{
			Symbol:    pair,
			Interval:  goBinance.Minute,
			StartTime: v.StartTime * 1000,
			EndTime:   v.EndTime * 1000,
			Limit:     1000,
		})
		if err != nil {
			return nil, ez.Wrap(op, err)
		}
		for _, vv := range kl {
			marketCandles = append(marketCandles, market.Candle{
				Time:   vv.OpenTime.Unix(),
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

func (b Client) FetchTicker(pair string) (*market.Ticker, error) {
	op := "binance.FetchTicker"

	thisMinute := time.Now()
	lastMinute := thisMinute.Add(time.Minute * -1)

	candles, err := b.FetchBinanceCandles(pair, lastMinute, thisMinute, 1)
	if err != nil {
		return nil, ez.Wrap(op, err)
	} else if len(candles) == 0 {
		return nil, ez.New(op, ez.ENOTFOUND, "Candle array is empty", nil)
	}

	options := &api.GetOrderBookOptions{Limit: 5}

	ob, err := b.FetchOrderBook(pair, options)
	if err != nil {
		return &market.Ticker{}, ez.Wrap(op, err)
	}
	firstAsk := ob.Asks[0] // ticker ask
	firstBid := ob.Bids[0] // ticker bid

	ticker := &market.Ticker{
		Time: time.Now().Unix(),
		Ask:  firstAsk.Price,
		Bid:  firstBid.Price,
	}
	return ticker, nil
}

func (b Client) FetchOrderBook(pair string, options *api.GetOrderBookOptions) (goBinance.OrderBook, error) {
	op := "binance.FetchOrderBook"

	limit := 100 // The default limit from the API

	if options.Limit != 0 {
		limit = options.Limit
	}

	obr := goBinance.OrderBookRequest{
		Symbol: pair,
		Limit:  limit,
	}
	orderBook, err := b.service.OrderBook(obr)
	if err != nil {
		return goBinance.OrderBook{}, ez.Wrap(op, err)
	}
	return *orderBook, nil
}
