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

// GetOrder - Returns an existing order

// CreateOrder - Places a new order

// CancelOrder - Cancels an existing order

// GetDepositMethods
func (c *Client) GetDepositMethods(asset string) ([]DepositMethods, error) {
	const op = "krakenclient.GetDepositMethods"
	const URL = "https://api.kraken.com/0/private/DepositMethods"

	// GetOrderBook - Gets order book for `pair` with `depth`
	data := url.Values{
		"asset": {asset},
	}

	response := make([]DepositMethods, 0)
	err := c.httpRequest("POST", URL, data, &response)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return response, nil
}

// GetDepositAddresses - Retrieve (or generate a new) deposit addresses for a particular asset and method.
func (c *Client) GetDepositAddresses(asset, method string) ([]DepositAddress, error) {
	const op = "krakenclient.GetDepositAddresses"
	const URL = "https://api.kraken.com/0/private/DepositAddresses"

	// GetOrderBook - Gets order book for `pair` with `depth`
	data := url.Values{
		"asset":  {asset},
		"method": {method},
	}

	response := make([]DepositAddress, 0)
	err := c.httpRequest("POST", URL, data, &response)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return nil, nil
}

// WithdrawAsset - Places a withdrawal request

// CancelWithdraw - Cancels an asset withdrawal
