package kucoin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vanclief/finmod/market"
)

func TestGetTicker(t *testing.T) {

	kucoin, err := New("", "", "")
	assert.Nil(t, err)

	ETH, _ := market.NewAsset("ETH", "Ethereum")
	USDT, _ := market.NewAsset("USDT", "US Tether")

	ETHUSDT, _ := market.NewPair(ETH, USDT)

	ticker, err := kucoin.GetTicker(ETHUSDT)
	assert.Nil(t, err)
	assert.NotNil(t, ticker)
}

func TestGetOrderBook(t *testing.T) {

	kucoin, err := New("", "", "")
	assert.Nil(t, err)

	ETH, _ := market.NewAsset("ETH", "Ethereum")
	USDT, _ := market.NewAsset("USDT", "US Tether")

	ETHUSDT, _ := market.NewPair(ETH, USDT)

	orderBook, err := kucoin.GetOrderBook(ETHUSDT)
	assert.Nil(t, err)
	assert.NotNil(t, orderBook)

	// fmt.Println("orderbook", orderBook)
	// assert.FailNow(t, "now")
}
