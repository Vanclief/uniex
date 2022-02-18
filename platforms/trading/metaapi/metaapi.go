package metaapi

import (
	"github.com/vanclief/ez"
	"github.com/vanclief/uniex/platforms/trading/metaapi/api"
	"github.com/vanclief/uniex/types"
)

// New returns a new MetaAPI TradingPlatform.
func New(accountID, token string) (platform types.TradingPlatform, err error) {
	const op = "metaapi.New"

	dataAPI, err := api.New(accountID, token)
	if err != nil {
		return platform, ez.Wrap(op, err)
	}

	platform.Name = "MetaAPI"
	platform.DataAPI = dataAPI

	return platform, nil
}
