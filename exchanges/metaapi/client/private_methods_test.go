package client

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var client *Client

func init() {
	viper.AutomaticEnv()
	accountID := viper.GetString("META_API_ACCOUNT_ID")
	token := viper.GetString("META_API_TOKEN")
	client = New(accountID, token)
}

func TestGetAccountInformation(t *testing.T) {
	response, err := client.GetAccountInformation()
	assert.Nil(t, err)
	assert.NotNil(t, response)
}

func TestGetPositions(t *testing.T) {
	response, err := client.GetPositions()
	assert.Nil(t, err)
	assert.NotNil(t, response)
}

func TestGetOrders(t *testing.T) {
	response, err := client.GetOrders()
	assert.Nil(t, err)
	assert.NotNil(t, response)
}

func TestTrade(t *testing.T) {
	request := &MetatraderTrade{}
	response, err := client.Trade(request)
	assert.Nil(t, err)
	assert.NotNil(t, response)
}
