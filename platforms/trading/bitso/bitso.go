package bitso

import (
	"github.com/vanclief/ez"
	"github.com/vanclief/uniex/interfaces/ws/generic"
	"github.com/vanclief/uniex/platforms/trading/bitso/api"
	"github.com/vanclief/uniex/types"
)

// New returns a new Bitso TradingPlatform.
func New(apiKey, secretKey string) (platform types.TradingPlatform, err error) {
	const op = "bitso.New"

	// Create the DataAPI
	dataAPI, err := api.New(apiKey, secretKey)
	if err != nil {
		return platform, ez.Wrap(op, err)
	}

	// Create the WebSocket
	dataWS, err := generic.New(host, parser, opts...)

	platform.Name = "Bitso"
	platform.DataAPI = dataAPI
	platform.DataWS = dataWS

	return platform, nil
}
