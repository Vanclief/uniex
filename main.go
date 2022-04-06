package main

import (
	"context"
	"fmt"
	"sort"

	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/ws/genericws"
	"github.com/vanclief/uniex/platforms/trading/kraken/ws"
)

func main() {

	var opts []genericws.Option

	// btc := market.Pair{
	// 	Base:  market.Asset{Symbol: "BTC"},
	// 	Quote: market.Asset{Symbol: "USD"},
	// }

	eth := market.Pair{
		Base:  market.Asset{Symbol: "ETH"},
		Quote: market.Asset{Symbol: "USD"},
	}

	// opts = append(opts, genericws.WithSubscriptionTo(btc))
	opts = append(opts, genericws.WithSubscriptionTo(eth))

	handler := ws.NewHandler()

	opts = append(opts, genericws.WithName("Kraken"))
	ws, _ := genericws.NewClient(handler, opts...)

	ctx := context.Background()

	wsChannel, _ := ws.Listen(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-wsChannel:
			if !ok {
				fmt.Println("Not ok")
			}

			if msg.OrderBook.Time > 0 {

				asks := msg.OrderBook.Asks
				bids := msg.OrderBook.Bids

				sort.SliceStable(asks, func(i, j int) bool {
					return asks[i].Price > asks[j].Price
				})

				fmt.Println("=======OB============")
				fmt.Println("ASKS", len(asks))
				// for i, v := range asks {
				// 	if i < len(asks)-10 {
				// 		continue
				// 	}
				// 	fmt.Println(v.Price, v.Volume)
				// }

				fmt.Println("---------------------")
				fmt.Println("BIDS", len(bids))
				// for i, v := range bids {
				// 	if i > 10 {
				// 		break
				// 	}
				// 	fmt.Println(v.Price, v.Volume)
				// }

				fmt.Println("=====================")

				// if bids[0].Price > asks[len(asks)-1].Price {
				// 	fmt.Println("BEST BID", bids[0], "BEST ASK", asks[len(asks)-1])
				// 	panic("THIS SHOULD NOT HAPPEN")
				// }
			}
		}
	}

}
