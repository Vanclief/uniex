package ws

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/ws/generic"
)

func TestWebsocket(t *testing.T) {

	// host := "wss://stream.binance.com:9443/ws"

	opts := []generic.Option{}

	marketPair := market.Pair{
		Base:  &market.Asset{Symbol: "BTC"},
		Quote: &market.Asset{Symbol: "USDT"},
	}

	opts = append(opts, generic.WithSubscriptionTo(marketPair))

	handler := NewHandler()

	opts = append(opts, generic.WithName("Binance"))
	ws, err := generic.NewClient(handler, opts...)

	assert.Nil(t, err)
	assert.NotNil(t, ws)

	ctx := context.Background()

	_, err = ws.ListenTicker(ctx, "ticker")
	ez.ErrorStacktrace(err)
	assert.Nil(t, err)

	// orderChannel, err := ws.ListenOrderBook(ctx)
	// assert.Nil(t, err)

	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 		return
	// 	// case order, ok := <-orderChannel:
	// 	// 	assert.True(t, ok)
	// 	// 	fmt.Println("order", order)

	// 	case tick, ok := <-tickerChannel:
	// 		assert.True(t, ok)
	// 		fmt.Println("tick", tick)
	// 	}
	// }

}
