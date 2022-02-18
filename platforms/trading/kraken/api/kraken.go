package kraken

import (
	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	krakenclient "github.com/vanclief/uniex/exchanges/kraken/client"
)

// API - Kraken API
type API struct {
	Client *krakenclient.Client
}

func New(publicKey, secretKey string) (*API, error) {
	client := krakenclient.New(publicKey, secretKey)
	return &API{Client: client}, nil
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
