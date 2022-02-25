package ws

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/ws/genericws"
)

func TestWebsocket(t *testing.T) {

	opts := []genericws.Option{}

	btcMXN := market.Pair{
		Base:  &market.Asset{Symbol: "BTC"},
		Quote: &market.Asset{Symbol: "MXN"},
	}

	ethMXN := market.Pair{
		Base:  &market.Asset{Symbol: "ETH"},
		Quote: &market.Asset{Symbol: "MXN"},
	}

	opts = append(opts, genericws.WithSubscriptionTo(btcMXN))
	opts = append(opts, genericws.WithSubscriptionTo(ethMXN))
	opts = append(opts, genericws.SetTimeout(5))

	handler := NewHandler()

	opts = append(opts, genericws.WithName("Bitso"))
	ws, err := genericws.NewClient(handler, opts...)

	assert.Nil(t, err)
	assert.NotNil(t, ws)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tickerChannel, err := ws.ListenTicker(ctx)
	assert.Nil(t, err)

	orderChannel, err := ws.ListenOrderBook(ctx)
	assert.Nil(t, err)
	assert.Fail(t, "Test")

	for {
		select {
		case <-ctx.Done():
			return
		case order, ok := <-orderChannel:
			assert.True(t, ok)
			assert.NotNil(t, order)
			// fmt.Println("order", order.Pair.String(), order.OrderBook)

		case tick, ok := <-tickerChannel:
			assert.True(t, ok)
			assert.NotNil(t, tick)
			// fmt.Println("tick", tick.Pair.String(), tick.Ticks)
		}
	}

}
