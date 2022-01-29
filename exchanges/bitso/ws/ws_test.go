package ws

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/vanclief/finmod/market"
	"testing"
	"time"
)

func TestClient_ListenOrders(t *testing.T) {
	host := "wss://ws.bitso.com"
	cxt, _ := context.WithDeadline(context.Background(), time.Now().Add(3 * time.Second))
	bc, err := New(host,
		WithSubscriptionTo(market.Pair{
			Base:  &market.Asset{
				Symbol: "BTC",
			},
			Quote: &market.Asset{
				Symbol: "MXN",
			},
		}),
		WithSubscriptionTo(market.Pair{
			Base:  &market.Asset{
				Symbol: "BTC",
			},
			Quote: &market.Asset{
				Symbol: "USD",
			},
		}),
	)

	assert.Nil(t, err)
	msgChan, err := bc.ListenOrderBook(cxt)
	assert.Nil(t, err)

	for msg := range msgChan {
		//fmt.Printf("order=%+v\n", msg)
		assert.NotNil(t, msg)
	}
}

func TestClient_ListenTicker(t *testing.T) {
	host := "wss://ws.bitso.com"
	cxt, _ := context.WithDeadline(context.Background(), time.Now().Add(5 * time.Second))
	bc, err := New(host,
		WithSubscriptionTo(market.Pair{
			Base:  &market.Asset{
				Symbol: "BTC",
			},
			Quote: &market.Asset{
				Symbol: "MXN",
			},
		}),
		WithSubscriptionTo(market.Pair{
			Base:  &market.Asset{
				Symbol: "BTC",
			},
			Quote: &market.Asset{
				Symbol: "USD",
			},
		}),
	)

	assert.Nil(t, err)
	msgChan, err := bc.ListenTicker(cxt)
	assert.Nil(t, err)

	for msg := range msgChan {
		//fmt.Printf("ticker=%+v\n", msg)
		assert.NotNil(t, msg)
	}
}
