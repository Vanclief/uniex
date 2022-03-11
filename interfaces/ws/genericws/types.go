package genericws

import (
	"time"

	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/ws"
)

var (
	ErrUnknownSubscriptionType = "unknown subscription type"
	ErrSubscriptionPairs       = "at least one subscription pair must be provided"
)

type ChannelType = string

const (
	OrderBookChannel ChannelType = "orderbook"
	TickerChannel    ChannelType = "ticker"
)

type ChannelOpts struct {
	Type ChannelType
}

var defaultChannels = []ChannelOpts{
	{
		Type: OrderBookChannel,
	},
	{
		Type: TickerChannel,
	},
}

type SubscriptionRequest []byte

type Settings struct {
	Endpoint                      string
	SubscriptionVerificationCount int
	PingTimeInterval              time.Duration
	PongWaitTime                  time.Duration
}

type HandlerOptions struct {
	Pairs    []market.Pair
	Channels []ChannelOpts
}

type WebsocketHandler interface {
	Init(opts HandlerOptions) error
	GetSettings() (Settings, error)
	GetSubscriptionsRequests() ([]SubscriptionRequest, error)
	VerifySubscriptionResponse(response []byte) error
	Parse(in []byte) (*ws.ListenChan, error)
}
