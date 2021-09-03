package oracles

import (
	"time"

	"github.com/vanclief/finmod/market"
)

// DataOracle - A Market Data Oracle
type DataOracle struct {
	Name string
	API  DataOracleAPI
}

// DataOracleAPI - Interface for a unified market data oracle API
type DataOracleAPI interface {
	GetTicker(pair *market.Pair) (*market.Ticker, error) // Will be deprecated
	GetCurrentCandle(pair *market.Pair, timeframe int) (*market.Candle, error)
	GetHistoricalData(pair *market.Pair, start, end time.Time, interval int) ([]market.Candle, error)
	// ListAssets() ([]market.Asset, error)
}
