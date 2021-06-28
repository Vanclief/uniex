package kraken

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/vanclief/finmod/market"
)

var krakenAPI *API

func init() {
	viper.AutomaticEnv()
	apiKey := viper.GetString("KRAKEN_API_KEY")
	secretKey := viper.GetString("KRAKEN_SECRET_KEY")

	krakenAPI, _ = New(apiKey, secretKey)
}

func TestGetTicker(t *testing.T) {

	ETH, _ := market.NewAsset("ETH", "Ethereum")
	USD, _ := market.NewAsset("USD", "US Dollar")

	ETHUSD, _ := market.NewPair(ETH, USD)

	ticker, err := krakenAPI.GetTicker(ETHUSD)
	assert.Nil(t, err)
	assert.NotNil(t, ticker)
}

func TestGetOrderBook(t *testing.T) {

	ETH, _ := market.NewAsset("ETH", "Ethereum")
	USD, _ := market.NewAsset("USD", "US Dollar")

	ETHUSD, _ := market.NewPair(ETH, USD)

	orderBook, err := krakenAPI.GetOrderBook(ETHUSD)
	assert.Nil(t, err)
	assert.NotNil(t, orderBook)
}

func TestGetBalances(t *testing.T) {
	balance, err := krakenAPI.GetBalances()
	assert.Nil(t, err)
	assert.NotNil(t, balance)
}
