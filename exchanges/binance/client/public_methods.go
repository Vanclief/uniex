package binanceclient

import (
	"math"
	"time"

	goBinance "github.com/binance-exchange/go-binance"
	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/exchanges"
	// 	"github.com/go-kit/kit/log"
	// 	"github.com/go-kit/kit/log/level"
	// 	"github.com/vanclief/ez"
)

// TimeInterval contains the UNIX timestamps to make a Binance API call with startTime and endTime, the maximum
// difference between endTime and startTime is 1000 minutes or 16 hours, 40 minutes
type TimeInterval struct {
	StartTime int64
	EndTime   int64
}

// CreateArrayOfTimestamps takes the start and end time as time.Time and calculates how many Binance API calls are
// necessary to cover the time period in minutes, every API call interval is stored in an Interval struct
func CreateArrayOfTimestamps(startTime, endTime time.Time) (timestamps []TimeInterval) {
	startUnix := startTime.Unix()
	endUnix := endTime.Unix()
	delta := int64(60 * 1000)
	loops := math.Ceil(float64(endUnix-startUnix) / float64(delta))
	if loops == 0 {
		return append(timestamps, TimeInterval{StartTime: startUnix, EndTime: endUnix})
	}

	startIndex := startUnix
	endIndex := int64(math.Min(float64(startUnix+delta), float64(endUnix)))

	for i := 0; i < int(loops); i++ {
		timestamps = append(timestamps, TimeInterval{
			StartTime: startIndex,
			EndTime:   endIndex,
		})
		startIndex += delta
		endIndex = int64(math.Min(float64(endIndex+delta), float64(endUnix)))
		//startIndex += delta + 1
		//endIndex = int64(math.Min(float64(endIndex+delta+1), float64(endUnix)))
	}
	return timestamps
}

// FetchBinanceCandles takes the start and end dates, and the crypto pair as strings, and returns the Binance Candle
// data for every minute between start and end for a crypto pair
func (b Client) FetchBinanceCandles(pair string, start, end time.Time) ([]market.Candle, error) {
	op := "binance.FetchBinanceCandles"

	var marketCandles []market.Candle

	arrayOfTimestamps := CreateArrayOfTimestamps(start, end)
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

type CandleOrderBook struct {
	candle    *market.Candle
	orderBook *market.OrderBook
}

func (b Client) FetchTicker(pair string) (*market.Ticker, error) {
	op := "binance.FetchTicker"

	thisMinute := time.Now()
	lastMinute := thisMinute.Add(time.Minute * -1)

	candles, err := b.FetchBinanceCandles(pair, lastMinute, thisMinute)
	if err != nil {
		return nil, ez.Wrap(op, err)
	} else if len(candles) == 0 {
		return nil, ez.New(op, ez.ENOTFOUND, "Candle array is empty", nil)
	}

	lastCandle := candles[len(candles)-1]

	options := &exchanges.GetOrderBookOptions{Limit: 5}

	ob, err := b.FetchOrderBook(pair, options)
	if err != nil {
		return &market.Ticker{}, ez.Wrap(op, err)
	}
	firstAsk := ob.Asks[0] // ticker ask
	firstBid := ob.Bids[0] // ticker bid

	ticker := &market.Ticker{
		Time: time.Now().Unix(),
		Candle: &market.Candle{
			Time:   lastCandle.Time,
			Open:   lastCandle.Open,
			High:   lastCandle.High,
			Low:    lastCandle.Low,
			Close:  lastCandle.Close,
			Volume: lastCandle.Volume,
		},
		Ask: &market.OrderBookRow{
			Price:       firstAsk.Price,
			Volume:      firstAsk.Quantity,
			AccumVolume: firstAsk.Quantity,
		},
		Bid: &market.OrderBookRow{
			Price:       firstBid.Price,
			Volume:      firstBid.Quantity,
			AccumVolume: firstBid.Quantity,
		},
	}
	return ticker, nil
}

func (b Client) FetchOrderBook(pair string, options *exchanges.GetOrderBookOptions) (goBinance.OrderBook, error) {
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
