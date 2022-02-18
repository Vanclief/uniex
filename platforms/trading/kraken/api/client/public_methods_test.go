package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTicker(t *testing.T) {

	krakenClient := New("", "")
	assert.NotNil(t, krakenClient)

	tick, err := krakenClient.GetTicker("ETHUSD")
	assert.Nil(t, err)
	assert.NotNil(t, tick)
}

func TestGetOrderBook(t *testing.T) {

	krakenClient := New("", "")
	assert.NotNil(t, krakenClient)

	ob, err := krakenClient.GetOrderBook("ETHUSD", 10)
	assert.Nil(t, err)
	assert.NotNil(t, ob)

}

func TestGetOHLCData(t *testing.T) {

	krakenClient := New("", "")
	assert.NotNil(t, krakenClient)

	ohlc, err := krakenClient.GetOHLCData("ETHUSD", 15, 0)
	assert.Nil(t, err)
	assert.NotNil(t, ohlc)
}

func TestGetAssetPairs(t *testing.T) {

	krakenClient := New("", "")
	assert.NotNil(t, krakenClient)

	pairs, err := krakenClient.GetAssetPairs()
	assert.Nil(t, err)
	assert.NotNil(t, pairs)
}
