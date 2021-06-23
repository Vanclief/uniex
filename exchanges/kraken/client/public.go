package kraken

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/vanclief/ez"
)

// GetAssetPairs returns an array with all the tradable AssetPairs
func (c *Client) GetAssetPairs() (map[string]AssetPair, error) {
	const op = "kraken.Client.GetAssetPairs"
	const URL = "https://api.kraken.com/0/public/AssetPairs"

	response := make(map[string]AssetPair)
	err := c.httpRequest("GET", URL, nil, &response)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return response, nil
}

// GetTicker - Gets array of tickers passed through `pairs` arg.
// `pairs` - array of needed pairs. All by default if empty array passed or `pairs` is nil.
func (c *Client) GetTicker(pairs ...string) (map[string]Ticker, error) {
	const op = "kraken.Client.Ticker"
	const URL = "https://api.kraken.com/0/public/Ticker"

	var data url.Values
	if len(pairs) > 0 {
		data = url.Values{
			"pair": {strings.Join(pairs, ",")},
		}
	} else {
		return nil, ez.New(op, ez.EINVALID, "Missing pairs", nil)
	}

	response := make(map[string]Ticker)
	err := c.httpRequest("POST", URL, data, &response)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return response, nil
}

// GetOrderBook - Gets order book for `pair` with `depth`
func (c *Client) GetOrderBook(pair string, depth int64) (map[string]OrderBook, error) {
	const op = "kraken.Client.GetOrderBook"
	const URL = "https://api.kraken.com/0/public/Depth"

	// GetOrderBook - Gets order book for `pair` with `depth`
	data := url.Values{
		"pair":  {pair},
		"count": {strconv.FormatInt(depth, 10)},
	}

	response := make(map[string]OrderBook)
	err := c.httpRequest("POST", URL, data, &response)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return response, nil
}

// GetOHLCData returns an array with
func (c *Client) GetOHLCData(pair string, interval, since int) (*OHLCResponse, error) {
	const op = "kraken.Client.GetOHLCData"
	const URL = "https://api.kraken.com/0/public/OHLC"

	data := url.Values{
		"pair":     {pair},
		"interval": {fmt.Sprint(interval)},
		"since":    {fmt.Sprint(since)},
	}

	response := &OHLCResponse{}
	err := c.httpRequest("POST", URL, data, response)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return response, nil
}
