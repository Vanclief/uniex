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

// DataOracleAPI - Interface for an unified market data oracle API
type DataOracleAPI interface {
	// Public Endpoints
	GetTicker(pair *market.Pair) (*market.Ticker, error)
	GetHistoricalData(pair *market.Pair, start, end time.Time) ([]market.Candle, error)
	// ListAssets() ([]market.Asset, error)
}
