package ws

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/ws/genericws"
	"testing"
)

func TestWs(t *testing.T) {

	var opts []genericws.Option

	btc := market.Pair{
		Base:  market.Asset{Symbol: "BTC"},
		Quote: market.Asset{Symbol: "USD"},
	}

	eth := market.Pair{
		Base:  market.Asset{Symbol: "ETH"},
		Quote: market.Asset{Symbol: "USD"},
	}

	opts = append(opts, genericws.WithSubscriptionTo(btc))
	opts = append(opts, genericws.WithSubscriptionTo(eth))

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
					assert.True(t, msg.OrderBook.Asks[0].Price > msg.OrderBook.Bids[0].Price)
				}
				fmt.Println("ob", msg.OrderBook.String())
			}

			//if len(msg.Tickers) > 0 && msg.Tickers[0].Time > 0 {
			//	fmt.Println("tick", msg.Tickers[0])
			//}
		}
	}
}
