package api

import (
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/vanclief/finmod/market"
)

var metaAPI *API

func init() {
	viper.AutomaticEnv()
	accountID := viper.GetString("META_API_ACCOUNT_ID")
	token := viper.GetString("META_API_TOKEN")
	metaAPI, _ = New(accountID, token)
}

func TestAPI_GetHistoricalData(t *testing.T) {
	pair := &market.Pair{
		Base: &market.Asset{
			Symbol: "#US30",
			Name:   "Dow Jones CFD",
		},
		Quote: &market.Asset{
			Symbol: "",
			Name:   "",
		},
	}

	start := time.Now().Add(-100 * time.Hour)
	end := time.Now()

	historicalData, err := metaAPI.GetHistoricalData(pair, start, end)
	assert.Nil(t, err)
	assert.NotNil(t, historicalData)
}
