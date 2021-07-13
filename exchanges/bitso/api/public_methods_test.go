package api

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/vanclief/finmod/market"
	"testing"
)

var bitsoAPI *API

func init() {
	viper.AutomaticEnv()
	apiKey := viper.GetString("BITSO_API_KEY")
	apiSecret  := viper.GetString("BITSO_SECRET")

	bitsoAPI, _ = New(apiKey, apiSecret)
}

func TestGetTicker(t *testing.T) {
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
	ticker, err := bitsoAPI.GetTicker(pair)
	assert.Nil(t, err)
	assert.NotNil(t, ticker)
}

func TestAPI_GetOrderBook(t *testing.T) {
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
	ob, err := bitsoAPI.GetOrderBook(pair)
	assert.Nil(t, err)
	assert.NotNil(t, ob)
}
