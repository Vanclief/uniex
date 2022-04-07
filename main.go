package main

import (
	"context"
	"fmt"

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

				msg.OrderBook.Print()

				// if msg.OrderBook.Bids[0].Price > msg.OrderBook.Asks[len(msg.OrderBook.Asks)-1].Price {
				// 	fmt.Println("BEST BID", msg.OrderBook.Bids[0], "BEST ASK", msg.OrderBook.Asks[len(msg.OrderBook.Asks)-1])
				// 	panic("THIS SHOULD NOT HAPPEN")
				// }
			}
		}
	}

}
