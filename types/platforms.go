package types

import (
	"github.com/vanclief/uniex/interfaces/api"
	"github.com/vanclief/uniex/interfaces/ws"
)

type DataPlatform struct {
	Name    string
	DataAPI api.DataAPI
	DataWS  ws.DataAPI
}

type TradingPlatform struct {
	Name             string
	DataAPI          api.DataAPI
	TradingAPI       api.TradingAPI
	DataWS           ws.DataAPI
	MakerFee         float64
	TakerFee         float64
	ManagedPositions bool
	HedgingMode      bool
}
