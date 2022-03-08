package ws

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/ws/genericws"
)

func TestWebsocket(t *testing.T) {

	opts := []genericws.Option{}

	btcMXN := market.Pair{
		Base:  market.Asset{Symbol: "BTC"},
		Quote: market.Asset{Symbol: "MXN"},
	}

	ethMXN := market.Pair{
		Base:  market.Asset{Symbol: "USDC"},
		Quote: market.Asset{Symbol: "MXN"},
	}

	opts = append(opts, genericws.WithSubscriptionTo(btcMXN))
	opts = append(opts, genericws.WithSubscriptionTo(ethMXN))

	handler := NewHandler()

	opts = append(opts, genericws.WithName("Tauros"))
	ws, err := genericws.NewClient(handler, opts...)

	assert.Nil(t, err)
	assert.NotNil(t, ws)

	// ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	tickerChannel, err := ws.ListenTicker(ctx)
	assert.Nil(t, err)

	// orderChannel, err := ws.ListenOrderBook(ctx)
	// assert.Nil(t, err)

	for {
		select {
		case <-ctx.Done():
			assert.FailNow(t, "CTX done")
			return
		// case order, ok := <-orderChannel:
		// 	assert.True(t, ok)
		// 	fmt.Println("order", order.Pair.String(), order.OrderBook)

		case tick, ok := <-tickerChannel:
			assert.True(t, ok)
			fmt.Println("tick", tick.Pair.String(), tick.Ticks)
		}
	}
}
