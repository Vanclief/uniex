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

func TestGetCurrentCandle(t *testing.T) {

	pair := &market.Pair{
		Base: &market.Asset{
			Symbol: "#US30",
			Name:   "Dow Jones",
		},
		Quote: &market.Asset{
			Symbol: "USD",
			Name:   "US Dollar",
		},
	}

	candle, err := metaAPI.GetCurrentCandle(pair, 2)
	assert.Nil(t, err)
	assert.NotNil(t, candle)
}

func TestAPI_GetHistoricalData(t *testing.T) {
	pair := &market.Pair{
		Base: &market.Asset{
			Symbol: "#Germany30",
			Name:   "Germany DAX",
		},
		Quote: &market.Asset{
			Symbol: "USD",
			Name:   "US Dollar",
		},
	}

	start := time.Now().Add(-100 * time.Hour)
	end := time.Now()

	historicalData, err := metaAPI.GetHistoricalData(pair, start, end, 1)
	assert.Nil(t, err)
	assert.NotNil(t, historicalData)
}
