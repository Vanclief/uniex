package binance

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func init() {
	viper.AutomaticEnv()
	apiKey := viper.GetString("KRAKEN_API_KEY")
	secretKey := viper.GetString("KRAKEN_SECRET_KEY")

	binanceAPI, _ = New(apiKey, secretKey)
}

func TestGetBalances(t *testing.T) {
	balances, err := binanceAPI.GetBalances()
	assert.Nil(t, err)
	assert.NotNil(t, balances)
}

func TestGetOrders(t *testing.T) {

	// orderBook, err := binanceAPI.GetOrderBook(ETHUSD)
	// assert.Nil(t, err)
	// assert.NotNil(t, orderBook)
}

func TestCreateOrder(t *testing.T) {

}

func TestCancelOrder(t *testing.T) {

}
