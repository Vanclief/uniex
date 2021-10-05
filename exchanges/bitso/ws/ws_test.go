package ws

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestClient_ListenOrders(t *testing.T) {
	host := "wss://ws.bitso.com"
	cxt, _ := context.WithDeadline(context.Background(), time.Now().Add(3 * time.Second))
	bc, err := new(host,
		WithSubscriptionTo("btc_usd", "orders"),
	)

	assert.Nil(t, err)
	msgChan, err := bc.ListenOrders(cxt)
	assert.Nil(t, err)

	for msg := range msgChan {
		//fmt.Printf("order=%+v\n", msg)
		assert.NotNil(t, msg)
	}
}

func TestClient_ListenTicker(t *testing.T) {
	host := "wss://ws.bitso.com"
	cxt, _ := context.WithDeadline(context.Background(), time.Now().Add(5 * time.Second))
	bc, err := new(host,
		WithSubscriptionTo("btc_mxn", "trades"),
	)

	assert.Nil(t, err)
	msgChan, err := bc.ListenTrade(cxt)
	assert.Nil(t, err)

	for msg := range msgChan {
		//fmt.Printf("ticker=%+v\n", msg)
		assert.NotNil(t, msg)
	}
}
