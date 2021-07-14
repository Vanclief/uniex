package api

import (
	"github.com/stretchr/testify/assert"
	"github.com/vanclief/finmod/market"
	"testing"
)

func TestGetTicker(t *testing.T) {
	bitsoClient, err := New("")
	assert.Nil(t, err)
	pair := &market.Pair{
		Base: &market.Asset{
			Symbol: "btc",
			Name:   "Bitcoin",
		},
		Quote: &market.Asset{
			Symbol: "mxn",
			Name:   "Mexican Peso",
		},
	}
	ticker, err := bitsoClient.GetTicker(pair)
	assert.Nil(t, err)
	assert.NotNil(t, ticker)
}

func TestAPI_GetOrderBook(t *testing.T) {
	bitsoClient, err := New("")
	assert.Nil(t, err)
	pair := &market.Pair{
		Base: &market.Asset{
			Symbol: "btc",
			Name:   "Bitcoin",
		},
		Quote: &market.Asset{
			Symbol: "mxn",
			Name:   "Mexican Peso",
		},
	}
	ob, err := bitsoClient.GetOrderBook(pair)
	assert.Nil(t, err)
	assert.NotNil(t, ob)
}
