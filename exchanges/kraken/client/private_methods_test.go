package krakenclient

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var client *Client

func init() {
	viper.AutomaticEnv()
	apiKey := viper.GetString("KRAKEN_API_KEY")
	secretKey := viper.GetString("KRAKEN_SECRET_KEY")

	client = New(apiKey, secretKey)
}

func TestGetFundings(t *testing.T) {

}

func TestGetAccountBalance(t *testing.T) {
	balance, err := client.GetAccountBalance()
	assert.Nil(t, err)
	assert.NotNil(t, balance)
}
