package ws

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient_ListenOrders(t *testing.T) {
	host := "wss://wsv2.tauros.io"
	cxt, _ := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	bc, err := new(host,
		WithSubscriptionTo("btc_mxn", OrdersChannel),
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
	host := "wss://wsv2.tauros.io"
	cxt, _ := context.WithDeadline(context.Background(), time.Now().Add(100*time.Second))
	bc, err := new(host,
		WithSubscriptionTo("btc_mxn", TickerChannel),
	)

	assert.Nil(t, err)
	msgChan, err := bc.ListenTicker(cxt)
	assert.Nil(t, err)

	for msg := range msgChan {
		//fmt.Printf("ticker=%+v\n", msg)
		assert.NotNil(t, msg)
	}
}
