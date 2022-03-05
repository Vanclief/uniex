package genericws

import (
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/ws"
)

var (
	ErrUnknownSubscriptionType = "Unknown subscription type"
	ErrSubscriptionPairs       = "At least one subscription pair must be provided"
)

type ChannelType string

var (
	ChannelTypeOrderBook ChannelType = "orderbook"
	ChannelTypeTicker    ChannelType = "ticker"
)

type SubscriptionRequest []byte

type WebsocketHandler interface {
	GetBaseEndpoint(pair []market.Pair, channelType ChannelType) string
	GetSubscriptionsRequests(pair []market.Pair, channelType ChannelType) ([]SubscriptionRequest, error)
	VerifySubscriptionResponse(response []byte) error
	ToTickers(in []byte) (*ws.TickerChan, error)
	ToOrderBook(in []byte) (*ws.OrderBookChan, error)
}
