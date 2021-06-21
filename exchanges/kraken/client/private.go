package kraken

import (
	"net/url"

	"github.com/vanclief/ez"
	"github.com/vanclief/uniex/exchanges/kraken/client/models"
)

type PrivateRequest struct {
	Nonce string `json:"nonce"`
}

// GetAccountBalance returns the balance of the account
func (c *Client) GetAccountBalance() (*models.Balance, error) {
	const op = "Client.GetAccountBalance"
	const URL = "https://api.kraken.com/0/private/TradeBalance"

	data := url.Values{}

	balance := &models.Balance{}
	err := c.httpRequest("POST", URL, data, balance)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return balance, nil
}
