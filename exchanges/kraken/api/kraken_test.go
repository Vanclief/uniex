package kraken

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vanclief/finmod/market"
)

func TestGetTicker(t *testing.T) {

	kraken, err := New("", "")
	assert.Nil(t, err)

	ETH, _ := market.NewAsset("ETH", "Ethereum")
	USD, _ := market.NewAsset("USD", "US Dollar")

	ETHUSD, _ := market.NewPair(ETH, USD)

	ticker, err := kraken.GetTicker(ETHUSD)
	assert.Nil(t, err)
	assert.NotNil(t, ticker)
}
