package kucoin

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func init() {
	viper.AutomaticEnv()
	apiKey := viper.GetString("KUCOIN_API_KEY")
	passphrase := viper.GetString("KUCOIN_PASSPHRASE")
	secretKey := viper.GetString("KUCOIN_SECRET")

	kucoinAPI, _ = New(apiKey, passphrase, secretKey)
}

func TestGetBalance(t *testing.T) {
	balance, err := kucoinAPI.GetBalance()
	assert.Nil(t, err)
	assert.NotNil(t, balance)
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
