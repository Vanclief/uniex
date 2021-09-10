package client

import (
	"fmt"
	"net/url"

	"github.com/vanclief/ez"
)

func (c *Client) GetAccountInformation() (interface{}, error) {
	const op = "MetaAPI.Client.GetAccountInformation"
	URL := fmt.Sprintf("https://mt-client-api-v1.agiliumtrade.agiliumtrade.ai/users/current/accounts/%s/accountInformation", c.AccountID)

	data := url.Values{}
	response := &MetatraderAccountInformation{}

	err := c.httpRequest("GET", URL, data, response)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return response, nil
}

func (c *Client) GetPositions() ([]MetatraderPosition, error) {
	const op = "MetaAPI.Client.GetAccountInformation"
	URL := fmt.Sprintf("https://mt-client-api-v1.agiliumtrade.agiliumtrade.ai/users/current/accounts/%s/positions", c.AccountID)

	data := url.Values{}
	response := &[]MetatraderPosition{}

	err := c.httpRequest("GET", URL, data, response)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return *response, nil
}

func (c *Client) GetOrders() ([]MetatraderOrder, error) {
	const op = "MetaAPI.Client.GetAccountInformation"
	URL := fmt.Sprintf("https://mt-client-api-v1.agiliumtrade.agiliumtrade.ai/users/current/accounts/%s/orders", c.AccountID)

	data := url.Values{}
	response := &[]MetatraderOrder{}

	err := c.httpRequest("GET", URL, data, response)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return *response, nil
}

func (c *Client) Trade(request *MetatraderTrade) (interface{}, error) {
	const op = "MetaAPI.Client.GetAccountInformation"
	URL := fmt.Sprintf("https://mt-client-api-v1.agiliumtrade.agiliumtrade.ai/users/current/accounts/%s/orders", c.AccountID)

	data := url.Values{}
	response := &MetatraderTradeResponse{}

	// TODO - Change data to be request
	err := c.httpRequest("POST", URL, data, response)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return response, nil
}
