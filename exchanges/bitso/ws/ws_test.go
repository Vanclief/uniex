package ws

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	host := "wss://ws.bitso.com"
	cxt, _ := context.WithDeadline(context.Background(), time.Now().Add(3 * time.Second))
	bc, err := New(host,
		WithSubscriptionTo("btc_usd", "orders"),
	)

	assert.Nil(t, err)
	msgChan := bc.Listen(cxt)
	for msg := range msgChan {
		//fmt.Printf("order=%+v\n", msg)
		assert.NotNil(t, msg)
	}
}
