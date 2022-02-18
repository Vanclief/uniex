package ws

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/exchanges/ws"
)

func TestWebsocket(t *testing.T) {

	host := "wss://ws.bitso.com"

	opts := []ws.Option{}

	marketPair := market.Pair{
		Base:  &market.Asset{Symbol: "BTC"},
		Quote: &market.Asset{Symbol: "MXN"},
	}

	opts = append(opts, ws.WithSubscriptionTo(marketPair))

	parser := NewParser()

	opts = append(opts, ws.WithName("Bitso"))
	ws, err := ws.New(host, parser, opts...)

	assert.Nil(t, err)
	assert.NotNil(t, ws)

	ctx := context.Background()

	tickerChannel, err := ws.ListenTicker(ctx)
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
