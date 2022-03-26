package utils

import (
	"github.com/vanclief/finmod/market"
	"sort"
	"time"
)

func GenerateOrderBookFromMap(asks, bids map[float64]float64) market.OrderBook {
	accumVol := 0.0
	asksArray := make([][]float64, 0, len(asks))
	bidsArray := make([][]float64, 0, len(bids))
	parsedOrderBook := market.OrderBook{
		Time: time.Now().Unix(),
	}

	for price, vol := range asks {
		asksArray = append(asksArray, []float64{price, vol})
	}

	for price, vol := range bids {
		bidsArray = append(bidsArray, []float64{price, vol})
	}
	sort.Slice(asksArray, func(i, j int) bool {
		return asksArray[i][0] < asksArray[j][0]
	})
	sort.Slice(bidsArray, func(i, j int) bool {
		return bidsArray[i][0] > bidsArray[j][0]
	})

	for _, v := range asksArray {
		accumVol += v[1]
		parsedOrderBook.Asks = append(parsedOrderBook.Asks, market.OrderBookRow{
			Price:       v[0],
			Volume:      v[1],
			AccumVolume: accumVol,
		})
	}

	accumVol = 0
	for _, v := range bidsArray {
		accumVol += v[1]
		parsedOrderBook.Bids = append(parsedOrderBook.Bids, market.OrderBookRow{
			Price:       v[0],
			Volume:      v[1],
			AccumVolume: accumVol,
		})
	}

	return parsedOrderBook
}
