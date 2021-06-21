package binance

// import (
// 	"context"
// 	"os"
// 	"strings"
// 	"time"

// 	goBinance "github.com/binance-exchange/go-binance"
// 	"github.com/go-kit/kit/log"
// 	"github.com/go-kit/kit/log/level"
// 	"github.com/vanclief/ez"
// 	"github.com/vanclief/go-trading-engine/config"
// 	"github.com/vanclief/go-trading-engine/market"
// 	"github.com/vanclief/go-trading-engine/utils"
// )

// // Interval contains the UNIX timestamps to make a Binance API call with startTime and endTime, the maximum
// // difference between endTime and startTime is 1000 minutes or 16 hours, 40 minutes
// type Interval struct {
// 	startTime int64
// 	endTime   int64
// }

// // Binance struct that contains the client for API calls and a context cancellable function
// type Binance struct {
// 	service   goBinance.Binance
// 	ctxCancel context.CancelFunc
// }

// func New(config *config.Config, env *config.Env) *Binance {
// 	var logger log.Logger
// 	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
// 	logger = level.NewFilter(logger, level.AllowAll())
// 	logger = log.With(logger, "time", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
// 	hmacSigner := &goBinance.HmacSigner{
// 		Key: []byte(env.BinanceSecret),
// 	}
// 	ctx, cancel := context.WithCancel(context.Background())
// 	// use second return value for cancelling request
// 	binanceService := goBinance.NewAPIService(
// 		config.BinanceHost,
// 		env.BinanceAPIKey,
// 		hmacSigner,
// 		logger,
// 		ctx,
// 	)
// 	b := goBinance.NewBinance(binanceService)
// 	return &Binance{
// 		service:   b,
// 		ctxCancel: cancel, // TODO: NO SABIA QUE HACER CON EL CANCEL, ASI QUE LO PASE COMO PARAMETRO A LA STRUCT
// 	}
// }

// // createArrayOfTimestamps takes the start and end time as time.Time and calculates how many Binance API calls are
// // necessary to cover the time period in minutes, every API call interval is stored in an Interval struct
// func createArrayOfTimestamps(timesArray []time.Time) []*Interval {
// 	startUnix := timesArray[0]
// 	endUnix := timesArray[1]
// 	var timestampArray []*Interval
// 	thousands := int((endUnix.Unix() - startUnix.Unix()) / 60000)
// 	cents := int((endUnix.Unix()-startUnix.Unix())/60) % 1000
// 	for i := 0; i < thousands; i++ {
// 		factor := 1000 * i
// 		temp := &Interval{
// 			startTime: startUnix.Add(time.Duration(factor) * time.Minute).Unix(),
// 			endTime:   startUnix.Add(time.Duration(999+factor) * time.Minute).Unix(),
// 		}
// 		timestampArray = append(timestampArray, temp)
// 	}
// 	timestampArray = append(timestampArray,
// 		&Interval{
// 			startTime: startUnix.Add(time.Duration(1000*thousands) * time.Minute).Unix(),
// 			endTime:   startUnix.Add(time.Duration(1000*thousands+cents-1) * time.Minute).Unix(),
// 		})
// 	return timestampArray
// }

// // parseArgs receives an array consisting of start and end dates as strings, and the crypto pair also as string
// // and parses the dates in time.RFC3339 format, and validates the crypto pair
// func parseArgs(args []string) ([]time.Time, string, error) {
// 	op := "binance.parseArgs"
// 	var parsedTimes []time.Time
// 	startTime, err := utils.RawDateToRFC3339(args[0])
// 	if err != nil {
// 		return nil, "", ez.New(op, ez.EINVALID, "cannot parse start date", nil)
// 	}
// 	endTime, err := utils.RawDateToRFC3339(args[1])
// 	if err != nil {
// 		return nil, "", ez.New(op, ez.EINVALID, "cannot parse end date", nil)
// 	}
// 	if startTime.Unix() > endTime.Unix() {
// 		startTime, endTime = endTime, startTime
// 	}
// 	parsedTimes = append(parsedTimes, startTime, endTime)
// 	return parsedTimes, strings.ToUpper(args[2]), nil
// }

// // FetchBinanceCandles takes the start and end dates, and the crypto pair as strings, and returns the Binance Candle
// // data for every minute between start and end for a crypto pair
// func (b Binance) FetchBinanceCandles(start, end, pair string) (*[]market.Candle, error) {
// 	op := "binance.FetchBinanceCandles"
// 	var marketCandles []market.Candle
// 	args := []string{start, end, pair}
// 	timesArray, uPair, err := parseArgs(args)
// 	if err != nil {
// 		return nil, ez.Wrap(op, err)
// 	}
// 	arrayOfTimestamps := createArrayOfTimestamps(timesArray)
// 	for _, v := range arrayOfTimestamps {
// 		kl, err := b.service.Klines(goBinance.KlinesRequest{
// 			Symbol:    uPair,
// 			Interval:  goBinance.Minute,
// 			StartTime: v.startTime * 1000,
// 			EndTime:   v.endTime * 1000,
// 			Limit:     1000,
// 		})
// 		if err != nil {
// 			ez.Wrap(op, err)
// 		}
// 		for _, vv := range kl {
// 			marketCandles = append(marketCandles, market.Candle{
// 				Time:   vv.OpenTime.Unix(),
// 				Open:   vv.Open,
// 				High:   vv.High,
// 				Low:    vv.Low,
// 				Close:  vv.Close,
// 				Volume: vv.Volume,
// 			})
// 		}
// 	}
// 	return &marketCandles, nil
// }
