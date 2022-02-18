package ws

import (
	"github.com/vanclief/finmod/market"
)

var (
	ErrUnknownSubscriptionType = "Unkown subscription type"
	ErrSubscriptionPairs       = "At least one subscription pair must be provided"
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
	Pair  market.Pair
	Ticks []market.Ticker
}

type OrderBookChan struct {
	Pair      market.Pair
	OrderBook market.OrderBook
}
