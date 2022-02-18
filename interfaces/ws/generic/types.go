package generic

import (
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/ws"
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
	ToTickers(in []byte) (*ws.TickerChan, error)
	ToOrderBook(in []byte) (*ws.OrderBookChan, error)
	GetSubscriptionRequest(pair market.Pair, channelType ChannelType) ([]byte, error)
}
