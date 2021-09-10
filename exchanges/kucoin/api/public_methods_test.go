package kucoin

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/exchanges"
)

var kucoinAPI *API

func init() {
	viper.AutomaticEnv()
	apiKey := viper.GetString("KUCOIN_API_KEY")
	secretKey := viper.GetString("KUCOIN_SECRET_KEY")
	passphrase := viper.GetString("KUCOIN_PASSPHRASE")

	kucoinAPI, _ = New(apiKey, secretKey, passphrase)
}

func TestGetTicker(t *testing.T) {

	ETH, _ := market.NewAsset("ETH", "Ethereum")
	USDT, _ := market.NewAsset("USDT", "US Tether")

	ETHUSDT, _ := market.NewPair(ETH, USDT)

	ticker, err := kucoinAPI.GetTicker(ETHUSDT)
	assert.Nil(t, err)
	assert.NotNil(t, ticker)
}

func TestGetOrderBook(t *testing.T) {

	ETH, _ := market.NewAsset("ETH", "Ethereum")
	USDT, _ := market.NewAsset("USDT", "US Tether")

	ETHUSDT, _ := market.NewPair(ETH, USDT)
	options := &exchanges.GetOrderBookOptions{Limit: 100}

	orderBook, err := kucoinAPI.GetOrderBook(ETHUSDT, options)
	assert.Nil(t, err)
	assert.NotNil(t, orderBook)
}
