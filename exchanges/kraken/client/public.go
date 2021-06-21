package kraken

import (
	"fmt"
	"net/url"

	"github.com/vanclief/ez"
	"github.com/vanclief/uniex/exchanges/kraken/client/models"
)

// GetAssetPairs returns an array with all the tradable AssetPairs
func (c *Client) GetAssetPairs() (map[string]models.AssetPair, error) {
	const op = "Client.GetAssetPairs"
	const URL = "https://api.kraken.com/0/public/AssetPairs"

	response := make(map[string]models.AssetPair)
	err := c.httpRequest("GET", URL, nil, &response)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return response, nil
}

// GetOHLCData returns an array with
func (c *Client) GetOHLCData(pair string, interval, since int) (*models.OHLCResponse, error) {
	const op = "Client.GetOHLCData"
	const URL = "https://api.kraken.com/0/public/OHLC"

	data := url.Values{
		"pair":     {pair},
		"interval": {fmt.Sprint(interval)},
		"since":    {fmt.Sprint(since)},
	}

	response := &models.OHLCResponse{}
	err := c.httpRequest("POST", URL, data, response)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return response, nil
}
