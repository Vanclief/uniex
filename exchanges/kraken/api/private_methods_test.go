package kraken

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var krakenAPI *API

func init() {
	viper.AutomaticEnv()
	apiKey := viper.GetString("KRAKEN_API_KEY")
	secretKey := viper.GetString("KRAKEN_SECRET_KEY")

	krakenAPI, _ = New(apiKey, secretKey)
}

func TestGetBalances(t *testing.T) {
	balance, err := krakenAPI.GetBalances()
	assert.Nil(t, err)
	assert.NotNil(t, balance)
}
