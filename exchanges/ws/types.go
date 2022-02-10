package ws

import (
	"errors"
	"github.com/vanclief/finmod/market"
)


var (
	ErrUnknownSubscriptionType = errors.New("unknown subscription type")
	ErrSubscriptionPairs = errors.New("at least one pair should be set")
)

type ChannelType string

var (
	ChannelTypeOrderBook ChannelType = "orderbook"
	ChannelTypeTicker    ChannelType = "ticker"
)

type WebsocketParser interface {
	ToTickers(in []byte) (*TickerChan, error)
	ToOrderBook(in []byte) (*OrderBookChan, error)
	GetSubscriptionRequest(pair market.Pair, channelType ChannelType) ([]byte, error)
}

type TickerChan struct {
	Pair   market.Pair
	Ticks []market.Ticker
}

type OrderBookChan struct {
	Pair      market.Pair
	OrderBook market.OrderBook
}
