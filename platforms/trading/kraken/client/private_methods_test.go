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

func TestGetAccountBalance(t *testing.T) {
	balance, err := client.GetAccountBalance()
	assert.Nil(t, err)
	assert.NotNil(t, balance)
}

func TestQueryOrders(t *testing.T) {
	orders, err := client.QueryOrders("ONVGMR-BOKKK-HES7ZW")
	assert.Nil(t, err)
	assert.NotNil(t, orders)
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

func TestWithdraw(t *testing.T) {
	// withdrawResponse, err := client.WithdrawFunds("ETH", "0x601486C5C19B035657aBe64d2f596317fa4939FB", 0.005)
	// assert.Nil(t, err)
	// assert.NotNil(t, withdrawResponse)
}
