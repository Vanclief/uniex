package client

import (
	"fmt"
	"net/url"

	"github.com/vanclief/ez"
)

func (c *Client) GetAccountInformation() (*MetatraderAccountInformation, error) {
	const op = "MetaAPI.Client.GetAccountInformation"
	URL := fmt.Sprintf("https://mt-client-api-v1.agiliumtrade.agiliumtrade.ai/users/current/accounts/%s/accountInformation", c.AccountID)

	data := url.Values{}
	response := &MetatraderAccountInformation{}

	err := c.httpRequest("GET", URL, data, nil, response)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return response, nil
}

func (c *Client) GetPositions() ([]MetatraderPosition, error) {
	const op = "MetaAPI.Client.GetPositions"
	URL := fmt.Sprintf("https://mt-client-api-v1.agiliumtrade.agiliumtrade.ai/users/current/accounts/%s/positions", c.AccountID)

	data := url.Values{}
	response := &[]MetatraderPosition{}

	err := c.httpRequest("GET", URL, data, nil, response)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return *response, nil
}

func (c *Client) GetDealByPositionID(id string) (*MetatraderDeal, error) {
	const op = "MetaAPI.Client.GetDealByPositionID"
	URL := fmt.Sprintf("https://mt-client-api-v1.agiliumtrade.agiliumtrade.ai/users/current/accounts/%s/history-deals/position/%s", c.AccountID, id)

	data := url.Values{}
	response := &[]MetatraderDeal{}

	err := c.httpRequest("GET", URL, data, nil, response)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	deals := *response

	return &deals[1], nil
}

func (c *Client) GetOrders() ([]MetatraderOrder, error) {
	const op = "MetaAPI.Client.GetOrders"
	URL := fmt.Sprintf("https://mt-client-api-v1.agiliumtrade.agiliumtrade.ai/users/current/accounts/%s/orders", c.AccountID)

	data := url.Values{}
	response := &[]MetatraderOrder{}

	err := c.httpRequest("GET", URL, data, nil, response)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return *response, nil
}

func (c *Client) GetOrderByID(id string) (*MetatraderOrder, error) {
	const op = "MetaAPI.Client.GetOrderByID"
	URL := fmt.Sprintf("https://mt-client-api-v1.agiliumtrade.agiliumtrade.ai/users/current/accounts/%s/history-orders/ticket/%s", c.AccountID, id)

	data := url.Values{}
	response := &[]MetatraderOrder{}

	err := c.httpRequest("GET", URL, data, nil, response)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	orders := *response

	return &orders[0], nil
}

func (c *Client) Trade(request *MetatraderTrade) (*MetatraderTradeResponse, error) {
	const op = "MetaAPI.Client.Trade"
	URL := fmt.Sprintf("https://mt-client-api-v1.agiliumtrade.agiliumtrade.ai/users/current/accounts/%s/trade", c.AccountID)

	response := &MetatraderTradeResponse{}

	err := c.httpRequest("POST", URL, nil, request, response)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return response, nil
}
