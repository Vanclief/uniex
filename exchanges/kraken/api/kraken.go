package kraken

import (
	"time"

	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	krakenclient "github.com/vanclief/uniex/exchanges/kraken/client"
)

// Kraken represents the interface
type API struct {
	Client *krakenclient.Client
}

func New(publicKey string, secretKey string) (*API, error) {
	client := krakenclient.New(publicKey, secretKey)
	return &API{Client: client}, nil
}

func (api *API) GetTicker(pair *market.Pair) (*market.Ticker, error) {
	const op = "kraken.GetTicker"

	requestTime := time.Now()
	symbol := pair.Base.Symbol + pair.Quote.Symbol

	tickerMap, err := api.Client.GetTicker(symbol)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	ticker := &market.Ticker{}

	// We only want the first item
	for _, value := range tickerMap {

		// Because kraken doesn't return the timestamp, we make a estimate
		// based on when did we made the request
		ticker.Time = requestTime.Unix()

		ticker.Candle = &market.Candle{
			Time:   requestTime.Add(-24 * time.Hour).Unix(), // The candle we get is from the past 24 hours
			Open:   value.OpeningPrice,
			High:   value.High.Price,
			Low:    value.Low.Price,
			Close:  value.Close.Price,
			Volume: value.Volume.Price,
		}

		ticker.Ask = &market.OrderBookRow{
			Price:       value.Ask.Price,
			Volume:      value.Ask.Volume,
			TotalVolume: value.Ask.WholeLotVolume,
		}

		ticker.Bid = &market.OrderBookRow{
			Price:       value.Bid.Price,
			Volume:      value.Bid.Volume,
			TotalVolume: value.Bid.WholeLotVolume,
		}
		break
	}

	return ticker, nil
}

func (api *API) GetHistoricalData(pair *market.Pair, start, end time.Time) ([]market.Candle, error) {
	const op = "kraken.GetHistoricalData"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) GetOrderBook(pair *market.Pair) (*market.OrderBook, error) {
	const op = "kraken.GetOrderBook"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) GetPositions() ([]market.Position, error) {
	const op = "kraken.GetPositions"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) GetTrades() ([]market.Trade, error) {
	const op = "kraken.GetTrades"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}
