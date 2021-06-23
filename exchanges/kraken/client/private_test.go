package krakenclient

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetFundings(t *testing.T) {

	viper.AutomaticEnv()
	apiKey := viper.GetString("KRAKEN_API_KEY")
	privateKey := viper.GetString("KRAKEN_PRIVATE_KEY")

	assert.NotEmpty(t, apiKey)

	client := New(apiKey, privateKey)
	assert.NotNil(t, client)
}

func GetAccountBalance(t *testing.T) {

	viper.AutomaticEnv()
	apiKey := viper.GetString("KRAKEN_API_KEY")
	privateKey := viper.GetString("KRAKEN_PRIVATE_KEY")

	client := New(apiKey, privateKey)
	assert.NotNil(t, client)

	balance, err := client.GetAccountBalance()
	assert.Nil(t, err)
	fmt.Println("balance", balance)
	fmt.Println("err", err)
}
