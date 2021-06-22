package kraken

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOrderBook(t *testing.T) {

	client := New("", "")
	assert.NotNil(t, client)

	ob, err := client.GetOrderBook("ETHUSD", 10)
	assert.Nil(t, err)
	assert.NotNil(t, ob)

}

func TestGetOHLCData(t *testing.T) {

	client := New("", "")
	assert.NotNil(t, client)

	ohlc, err := client.GetOHLCData("ETHUSD", 15, 0)
	assert.Nil(t, err)
	assert.NotNil(t, ohlc)

	assert.FailNow(t, "now")
}

func TestGetAssetPairs(t *testing.T) {

	client := New("", "")
	assert.NotNil(t, client)

	pairs, err := client.GetAssetPairs()
	assert.Nil(t, err)
	assert.NotNil(t, pairs)
}
