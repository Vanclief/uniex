package binanceclient

import (

	// 	"os"

	"time"

	goBinance "github.com/binance-exchange/go-binance"
	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	// 	"github.com/go-kit/kit/log"
	// 	"github.com/go-kit/kit/log/level"
	// 	"github.com/vanclief/ez"
)

// TimeInterval contains the UNIX timestamps to make a Binance API call with startTime and endTime, the maximum
// difference between endTime and startTime is 1000 minutes or 16 hours, 40 minutes
type TimeInterval struct {
	startTime int64
	endTime   int64
}

// createArrayOfTimestamps takes the start and end time as time.Time and calculates how many Binance API calls are
// necessary to cover the time period in minutes, every API call interval is stored in an Interval struct
func createArrayOfTimestamps(startTime, endTime time.Time) (timestamps []TimeInterval) {

	thousands := int((endTime.Unix() - startTime.Unix()) / 60000)
	cents := int((endTime.Unix()-startTime.Unix())/60) % 1000

	for i := 0; i < thousands; i++ {
		factor := 1000 * i
		temp := TimeInterval{
			startTime: startTime.Add(time.Duration(factor) * time.Minute).Unix(),
			endTime:   startTime.Add(time.Duration(999+factor) * time.Minute).Unix(),
		}
		timestamps = append(timestamps, temp)
	}

	timestamps = append(timestamps,
		TimeInterval{
			startTime: startTime.Add(time.Duration(1000*thousands) * time.Minute).Unix(),
			endTime:   startTime.Add(time.Duration(1000*thousands+cents-1) * time.Minute).Unix(),
		})

	return timestamps
}

// FetchBinanceCandles takes the start and end dates, and the crypto pair as strings, and returns the Binance Candle
// data for every minute between start and end for a crypto pair
func (b Client) FetchBinanceCandles(pair string, start, end time.Time) ([]market.Candle, error) {
	op := "binance.FetchBinanceCandles"

	var marketCandles []market.Candle

	arrayOfTimestamps := createArrayOfTimestamps(start, end)
	for _, v := range arrayOfTimestamps {
		kl, err := b.service.Klines(goBinance.KlinesRequest{
			Symbol:    pair,
			Interval:  goBinance.Minute,
			StartTime: v.startTime * 1000,
			EndTime:   v.endTime * 1000,
			Limit:     1000,
		})
		if err != nil {
			ez.Wrap(op, err)
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
