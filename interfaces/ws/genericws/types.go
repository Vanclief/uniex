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

type SubscriptionRequest []byte

type Settings struct {
	Endpoint                      string
	SubscriptionVerificationCount int
	PingTimeInterval              time.Duration
	PongWaitTime                  time.Duration
}

type WebsocketHandler interface {
	GetSettings(pair []market.Pair) (Settings, error)
	GetSubscriptionsRequests(pair []market.Pair) ([]SubscriptionRequest, error)
	VerifySubscriptionResponse(response []byte) error
	Parse(in []byte) (*ws.ListenChan, error)
}
