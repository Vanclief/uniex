package kucoin

import (
	"github.com/Kucoin/kucoin-go-sdk"
)

// API - Kucoin API
type API struct {
	Client *kucoin.ApiService
}

func New(apiKey, apiSecret, passphrase string) (*API, error) {

	client := kucoin.NewApiService(
		kucoin.ApiKeyOption(apiKey),
		kucoin.ApiSecretOption(apiSecret),
		kucoin.ApiPassPhraseOption(passphrase),
		kucoin.ApiKeyVersionOption(kucoin.ApiKeyVersionV2),
	)

	return &API{Client: client}, nil
}
