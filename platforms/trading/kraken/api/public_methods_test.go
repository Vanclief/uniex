package kraken

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/api"
)

func init() {
	viper.AutomaticEnv()
	apiKey := viper.GetString("KRAKEN_API_KEY")
	secretKey := viper.GetString("KRAKEN_SECRET_KEY")

	krakenAPI, _ = New(apiKey, secretKey)
}

func TestGetTicker(t *testing.T) {

	ETH, _ := market.NewAsset("ETH", "Ethereum")
	USD, _ := market.NewAsset("USD", "US Dollar")

	ETHUSD := market.NewPair(ETH, USD)

	ticker, err := krakenAPI.GetTicker(&ETHUSD)
	assert.Nil(t, err)
	assert.NotNil(t, ticker)
}

func TestGetOrderBook(t *testing.T) {

	ETH, _ := market.NewAsset("ETH", "Ethereum")
	USD, _ := market.NewAsset("USD", "US Dollar")

	ETHUSD := market.NewPair(ETH, USD)
	options := &api.GetOrderBookOptions{Limit: 100}

	orderBook, err := krakenAPI.GetOrderBook(&ETHUSD, options)
	assert.Nil(t, err)
	assert.NotNil(t, orderBook)
}
