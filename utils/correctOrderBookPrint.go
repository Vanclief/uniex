package utils

import (
	"fmt"
	"github.com/vanclief/uniex/market"
	"math"
)

func CorrectOrderBookPrint(ob market.OrderBook) {
	fmt.Println(" ---- OrderBook ----")
	asksIndex := int(math.Min(float64(len(ob.Asks)-1), 5))
	for i := asksIndex; i >= 0; i-- {
		fmt.Printf("%.2f\t%.6f\n", ob.Asks[i].Price, ob.Asks[i].Volume)
	}
	fmt.Println("--------------------")
	bidsIndex := int(math.Min(float64(len(ob.Bids)-1), 5))
	for i := bidsIndex; i >= 0; i-- {
		fmt.Printf("%.2f\t%.6f\n", ob.Bids[i].Price, ob.Bids[i].Volume)
	}
}
