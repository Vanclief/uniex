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

	btc := market.Pair{
		Base:  market.Asset{Symbol: "BTC"},
		Quote: market.Asset{Symbol: "MXN"},
	}

	eth := market.Pair{
		Base:  market.Asset{Symbol: "ETH"},
		Quote: market.Asset{Symbol: "MXN"},
	}

	opts = append(opts, genericws.WithSubscriptionTo(btc))
	opts = append(opts, genericws.WithSubscriptionTo(eth))

	handler := NewHandler()

	opts = append(opts, genericws.WithName("Tauros"))
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
				fmt.Println("ob", msg.OrderBook)
			}

			if len(msg.Tickers) > 0 && msg.Tickers[0].Time > 0 {
				fmt.Println("tick", msg.Tickers[0])
			}
		}
	}
}
