package ws

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/ws/generic"
)

func TestWebsocket(t *testing.T) {

	opts := []generic.Option{}

	btcMXN := market.Pair{
		Base:  &market.Asset{Symbol: "BTC"},
		Quote: &market.Asset{Symbol: "MXN"},
	}

	ethMXN := market.Pair{
		Base:  &market.Asset{Symbol: "ETH"},
		Quote: &market.Asset{Symbol: "MXN"},
	}

	opts = append(opts, generic.WithSubscriptionTo(btcMXN))
	opts = append(opts, generic.WithSubscriptionTo(ethMXN))

	handler := NewHandler()

	opts = append(opts, generic.WithName("Bitso"))
	ws, err := generic.NewClient(handler, opts...)

	assert.Nil(t, err)
	assert.NotNil(t, ws)

	ctx := context.Background()

	tickerChannel, err := ws.ListenTicker(ctx, "ticker")
	assert.Nil(t, err)

	// orderChannel, err := ws.ListenOrderBook(ctx)
	// assert.Nil(t, err)

	for {
		select {
		case <-ctx.Done():
			return
		// case order, ok := <-orderChannel:
		// 	assert.True(t, ok)
		// 	fmt.Println("order", order)

		case tick, ok := <-tickerChannel:
			assert.True(t, ok)
			fmt.Println("tick", tick)
		}
	}

}
