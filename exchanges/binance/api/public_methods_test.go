package binance

import (
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/exchanges"
)

var binanceAPI *API

func init() {
	viper.AutomaticEnv()
	apiKey := viper.GetString("KRAKEN_API_KEY")
	secretKey := viper.GetString("KRAKEN_SECRET_KEY")

	binanceAPI, _ = New(apiKey, secretKey)
}

func TestGetTicker(t *testing.T) {

	ETH, _ := market.NewAsset("ETH", "Ethereum")
	USD, _ := market.NewAsset("USDT", "Tether")

	ETHUSD, _ := market.NewPair(ETH, USD)

	ticker, err := binanceAPI.GetTicker(ETHUSD)
	assert.Nil(t, err)
	assert.NotNil(t, ticker)
}

func TestGetOrderBook(t *testing.T) {

	ETH, _ := market.NewAsset("ETH", "Ethereum")
	USD, _ := market.NewAsset("USDT", "Tether")

	ETHUSD, _ := market.NewPair(ETH, USD)
	options := &exchanges.GetOrderBookOptions{Limit: 100}

	orderBook, err := binanceAPI.GetOrderBook(ETHUSD, options)
	assert.Nil(t, err)
	assert.NotNil(t, orderBook)
}

func TestGetHistoricalData(t *testing.T) {

	ETH, _ := market.NewAsset("ETH", "Ethereum")
	USDT, _ := market.NewAsset("USDT", "Tether")

	start := time.Now().Add(-100 * time.Hour)
	end := time.Now()

	ETHUSDT, _ := market.NewPair(ETH, USDT)

	candles, err := binanceAPI.GetHistoricalData(ETHUSDT, start, end)
	assert.Nil(t, err)
	assert.NotNil(t, candles)
}
