package kraken

import (
	"github.com/vanclief/finmod/market"
	krakenclient "github.com/vanclief/uniex/exchanges/kraken/client"
)

// Kraken represents the interface
type API struct {
	Client *krakenclient.Client
}

func New(publicKey string, secretKey string) (*API, error) {
	return nil, nil
}

func (a *API) GetPositions() ([]market.Position, error) {
	return nil, nil
}
