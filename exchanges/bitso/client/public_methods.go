package client

import (
	"github.com/vanclief/ez"
	"net/url"
)

func (c *Client) GetTicker(pair string) (*Ticker, error) {
	op := "bitsoClient.Ticker"
	URL := "https://api.bitso.com/v3/ticker/"
	if pair == "" {
		return nil, ez.New(op, ez.EINVALID, "Missing pairs", nil)
	}
	data := url.Values{
		"book": {pair},
	}

	response := &Ticker{}
	err := c.httpRequest("GET", URL, data, &response)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}
	return response, nil
}

func (c *Client) GetOrderBook(pair string) (*OrderBook, error) {
	op := "bitsoClient.GetOrderBook"
	URL := "https://api.bitso.com/v3/order_book/"

	if pair == "" {
		return nil, ez.New(op, ez.EINVALID, "Pair is not present", nil)
	}

	data := url.Values{
		"book": {pair},
	}

	response := &OrderBook{}
	err := c.httpRequest("GET", URL, data, &response)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return response, nil
}

func (c *Client) GetOHLCData() error {
	op := "bitsoClient.GetOHLCData"
	return ez.New(op, ez.ENOTIMPLEMENTED, "Method not available in Bitso API", nil)
}
