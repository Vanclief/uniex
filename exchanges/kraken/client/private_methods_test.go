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

func TestGetDepositMethods(t *testing.T) {
	methods, err := client.GetDepositMethods("XBT")
	assert.Nil(t, err)
	assert.NotNil(t, methods)

	// Use this to map the methods
	// for _, m := range methods {
	// 	fmt.Println("Method", m.Method)
	// }

	// assert.FailNow(t, "Now")
}

func TestGetDepositAddresses(t *testing.T) {
	// address, err := client.GetDepositAddresses("ETH", "Ether (Hex)")
	address, err := client.GetDepositAddresses("XBT", "Bitcoin")
	assert.Nil(t, err)
	assert.NotNil(t, address)
}
