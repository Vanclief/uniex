package krakenclient

import (
	"net/url"

	"github.com/vanclief/ez"
)

type PrivateRequest struct {
	Nonce string `json:"nonce"`
}

// GetAccountBalance returns the balance of the account
func (c *Client) GetAccountBalance() (*BalanceResponse, error) {
	const op = "Client.GetAccountBalance"
	const URL = "https://api.kraken.com/0/private/Balance"

	data := url.Values{}

	balance := &BalanceResponse{}
	err := c.httpRequest("POST", URL, data, balance)
	if err != nil {

		return nil, ez.Wrap(op, err)
	}

	return balance, nil
}

// CreateOrder
