package ws

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/ws/genericws"
)

func TestWebsocket(t *testing.T) {

	// host := "wss://stream.binance.com:9443/ws"

	opts := []genericws.Option{}

	marketPair := market.Pair{
		Base:  &market.Asset{Symbol: "BTC"},
		Quote: &market.Asset{Symbol: "USDT"},
	}

	opts = append(opts, genericws.WithSubscriptionTo(marketPair))

	handler := NewHandler()

	opts = append(opts, genericws.WithName("Binance"))
	ws, err := genericws.NewClient(handler, opts...)

	assert.Nil(t, err)
	assert.NotNil(t, ws)

	ctx := context.Background()

	_, err = ws.ListenTicker(ctx)
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
