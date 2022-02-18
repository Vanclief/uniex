package krakenclient

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

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

// QueryOrders - Returns an existing order by ID
func (c *Client) QueryOrders(txIDs ...string) (map[string]OrderInfo, error) {
	data := url.Values{}
	const URL = "https://api.kraken.com/0/private/QueryOrders"

	switch {
	case len(txIDs) > 50:
		return nil, fmt.Errorf("Maximum count of requested orders is 50")
	case len(txIDs) == 0:
		return nil, fmt.Errorf("txIDs is required")
	default:
		data.Set("txid", strings.Join(txIDs, ","))
	}

	response := make(map[string]OrderInfo)
	err := c.httpRequest("POST", URL, data, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// AddOrder - Places a new order

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

// WithdrawFunds - Make a withdrawal request.
func (c *Client) WithdrawFunds(asset, key string, amount float64) (*WithdrawFundsResponse, error) {
	const op = "krakenclient.GetDepositAddresses"
	const URL = "https://api.kraken.com/0/private/Withdraw"

	// GetOrderBook - Gets order book for `pair` with `depth`
	data := url.Values{
		"asset":  {asset},
		"key":    {key},
		"amount": {strconv.FormatFloat(amount, 'f', -1, 64)},
	}

	response := &WithdrawFundsResponse{}
	err := c.httpRequest("POST", URL, data, &response)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return nil, nil
}

// CancelWithdraw - Cancels an asset withdrawal
