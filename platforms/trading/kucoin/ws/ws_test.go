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

	opts := []genericws.Option{}

	btcMXN := market.Pair{
		Base:  market.Asset{Symbol: "BTC"},
		Quote: market.Asset{Symbol: "USDT"},
	}

	ethMXN := market.Pair{
		Base:  market.Asset{Symbol: "ETH"},
		Quote: market.Asset{Symbol: "USDT"},
	}

	opts = append(opts, genericws.WithSubscriptionTo(btcMXN))
	opts = append(opts, genericws.WithSubscriptionTo(ethMXN))

	handler := NewHandler()

	opts = append(opts, genericws.WithName("Kucoin"))
	ws, err := genericws.NewClient(handler, opts...)

	assert.Nil(t, err)
	assert.NotNil(t, ws)

	ctx := context.Background()

	tickerChannel, err := ws.ListenTicker(ctx)
	assert.Nil(t, err)

	orderChannel, err := ws.ListenOrderBook(ctx)
	assert.Nil(t, err)

	for {
		select {
		case <-ctx.Done():
			return
		case order, ok := <-orderChannel:
			assert.True(t, ok)
			fmt.Println("order", order.Pair.String(), "Ask", order.OrderBook.Asks[0].Price, "Bid", order.OrderBook.Bids[0].Price)

		case tick, ok := <-tickerChannel:
			assert.True(t, ok)
			assert.NotNil(t, tick)
			fmt.Println("tick", tick.Pair.String(), tick.Ticks[0].Last)
		}
	}
}
