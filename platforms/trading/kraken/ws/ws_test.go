package ws

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/ws/genericws"
)

func TestWs(t *testing.T) {

	var opts []genericws.Option

	btc := market.Pair{
		Base:  market.Asset{Symbol: "BTC"},
		Quote: market.Asset{Symbol: "USD"},
	}

	//eth := market.Pair{
	//	Base:  market.Asset{Symbol: "ETH"},
	//	Quote: market.Asset{Symbol: "USD"},
	//}

	opts = append(opts, genericws.WithSubscriptionTo(btc))
	//opts = append(opts, genericws.WithSubscriptionTo(eth))

	handler := NewHandler()

	opts = append(opts, genericws.WithName("Kraken"))
	ws, err := genericws.NewClient(handler, opts...)

	assert.Nil(t, err)
	assert.NotNil(t, ws)

	ctx := context.Background()

	wsChannel, err := ws.Listen(ctx)
	assert.Nil(t, err)

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-wsChannel:
			assert.True(t, ok)

			if msg.OrderBook.Time > 0 {
				if len(msg.OrderBook.Bids) > 0 && len(msg.OrderBook.Asks) > 0 {
					if msg.OrderBook.Asks[0].Price < msg.OrderBook.Bids[0].Price {
						msg.OrderBook.Print()
						t.FailNow()
					}
					fmt.Println(len(msg.OrderBook.Asks), len(msg.OrderBook.Bids))
					//fmt.Println("------------------------------------------------------")
					//fmt.Println(msg.OrderBook.Asks[0].Price, "\t", msg.OrderBook.Bids[0].Price)
					//fmt.Println(msg.OrderBook.Asks[len(msg.OrderBook.Asks)-1].Price, "\t", msg.OrderBook.Bids[len(msg.OrderBook.Bids)-1].Price)
				} else {
					fmt.Println("empty")
				}
				//fmt.Println("ob", msg.OrderBook.String())
			}

			//if len(msg.Tickers) > 0 && msg.Tickers[0].Time > 0 {
			//	fmt.Println("tick", msg.Tickers[0])
			//}
		}
	}
}

func BenchmarkWs(t *testing.B) {

	var opts []genericws.Option

	btc := market.Pair{
		Base:  market.Asset{Symbol: "BTC"},
		Quote: market.Asset{Symbol: "USD"},
	}

	//eth := market.Pair{
	//	Base:  market.Asset{Symbol: "ETH"},
	//	Quote: market.Asset{Symbol: "USD"},
	//}

	opts = append(opts, genericws.WithSubscriptionTo(btc))
	//opts = append(opts, genericws.WithSubscriptionTo(eth))

	handler := NewHandler()

	opts = append(opts, genericws.WithName("Kraken"))
	ws, err := genericws.NewClient(handler, opts...)

	assert.Nil(t, err)
	assert.NotNil(t, ws)

	ctx := context.Background()

	wsChannel, err := ws.Listen(ctx)
	assert.Nil(t, err)

	select {
	case <-ctx.Done():
		return
	case msg, ok := <-wsChannel:
		assert.True(t, ok)

		if msg.OrderBook.Time > 0 {
			if len(msg.OrderBook.Bids) > 0 && len(msg.OrderBook.Asks) > 0 {
				if msg.OrderBook.Asks[0].Price < msg.OrderBook.Bids[0].Price {
					fmt.Println("ob", msg.OrderBook.String())
					t.FailNow()
				}
				fmt.Println(msg.Pair.String(), len(msg.OrderBook.Asks), len(msg.OrderBook.Bids))
			}
			//fmt.Println("ob", msg.OrderBook.String())
		}

		//if len(msg.Tickers) > 0 && msg.Tickers[0].Time > 0 {
		//	fmt.Println("tick", msg.Tickers[0])
		//}
	}
}
