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

type WebsocketHandler interface {
	GetSettings(pair []market.Pair, channels []ChannelOpts) (Settings, error)
	GetSubscriptionsRequests(pair []market.Pair, channels []ChannelOpts) ([]SubscriptionRequest, error)
	VerifySubscriptionResponse(response []byte) error
	Parse(in []byte) (*ws.ListenChan, error)
}
