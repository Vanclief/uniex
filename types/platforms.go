package types

import (
	"github.com/vanclief/uniex/interfaces/api"
	"github.com/vanclief/uniex/interfaces/ws"
)

type DataPlatform struct {
	Name     string
	DataAPI  api.DataAPI
	PublicWS ws.PublicWS
}

type TradingPlatform struct {
	Name             string
	DataAPI          api.DataAPI
	TradingAPI       api.TradingAPI
	PublicWS         ws.PublicWS
	MakerFee         float64
	TakerFee         float64
	ManagedPositions bool
	HedgingMode      bool
}
