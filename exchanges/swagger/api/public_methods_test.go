package api

import (
  "github.com/spf13/viper"
  "github.com/stretchr/testify/assert"
  "github.com/vanclief/finmod/market"
  "testing"
  "time"
)

var swaggerAPI *API

func init() {
 viper.AutomaticEnv()
 accountID := viper.GetString("SWAGGER_ACCOUNTID")
 token := viper.GetString("SWAGGER_TOKEN")
 swaggerAPI, _ = New(accountID, token)
}

func TestAPI_GetHistoricalData(t *testing.T) {
  pair := &market.Pair{
    Base: &market.Asset{
      Symbol: "%23US30",
      Name: "usd 30",
    },
    Quote: &market.Asset{
      Symbol: "",
      Name: "placeholder",
    },
  }

  start := time.Now().Add(-100 * time.Hour)
  end := time.Now()

  historicalData, err := swaggerAPI.GetHistoricalData(pair, start, end)
  assert.Nil(t, err)
  assert.NotNil(t, historicalData)
}
