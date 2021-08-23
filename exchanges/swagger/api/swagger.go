package api

import (
  swaggerClient "github.com/vanclief/uniex/exchanges/swagger/client"
)

type API struct {
  Client *swaggerClient.Client
}

func New(accountID, token string) (*API, error) {
  client := swaggerClient.New(accountID, token)
  return &API{Client: client}, nil
}