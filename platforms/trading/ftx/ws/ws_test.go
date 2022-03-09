package ws

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/ws/genericws"
	"testing"
	"time"
)

func TestWebsocket(t *testing.T) {
	var opts []genericws.Option

	btc := market.Pair{
		Base:  market.Asset{Symbol: "BTC"},
		Quote: market.Asset{Symbol: "PERP"},
	}

	//usd := market.Pair{
	//	Base:  market.Asset{Symbol: "ETH"},
	//	Quote: market.Asset{Symbol: "USD"},
	//}

	opts = append(opts, genericws.WithSubscriptionTo(btc))
	//opts = append(opts, genericws.WithSubscriptionTo(usd))

	handler := NewHandler()

	opts = append(opts, genericws.WithName("FTX"))
	ws, err := genericws.NewClient(handler, opts...)

	assert.Nil(t, err)
	assert.NotNil(t, ws)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tickerChannel, err := ws.ListenTicker(ctx)
	assert.Nil(t, err)

	orderChannel, err := ws.ListenOrderBook(ctx)
	assert.Nil(t, err)
	//assert.Equal(t, "", "Test")

	for {
		select {
		case <-ctx.Done():
			return
		case order, ok := <-orderChannel:
			assert.True(t, ok)
			fmt.Println("order\n", order.Pair.String(), order.OrderBook.Asks, order.OrderBook.Bids)

		case tick, ok := <-tickerChannel:
			assert.True(t, ok)
			fmt.Println("tick", tick.Pair.String(), tick.Ticks)
		}
	}
}
