package kraken

import (
	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/platforms/trading/kraken/api/client"
)

// API - Kraken API
type API struct {
	Client *client.Client
}

func New(publicKey, secretKey string) (*API, error) {
	c := client.New(publicKey, secretKey)
	return &API{Client: c}, nil
}

func TranslateAsset(symbol string) (*market.Asset, error) {
	const op = "kraken.TranslateAsset"

	asset, err := market.NewCryptoAsset(symbol)
	if err == nil {
		return asset, nil
	}

	asset, err = market.NewForexAsset(symbol)
	if err == nil {
		return asset, nil
	}

	switch symbol {
	case "ZUSD":
		return market.NewForexAsset("USD")
	case "ZEUR":
		return market.NewForexAsset("EUR")

	default:
		return nil, ez.New(op, ez.ENOTFOUND, "No asset with matching symbol found", nil)
	}
}
