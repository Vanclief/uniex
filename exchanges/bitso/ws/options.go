package ws

import (
	"errors"
	"github.com/vanclief/finmod/market"
	"strings"
)

var (
	ErrUnknownSubscriptionType = errors.New("unknown subscription type")
)

type optionApplyFunc func(client *baseClient) error

type Option interface {
	applyOption(client *baseClient) error
}

func (f optionApplyFunc) applyOption(p *baseClient) error {
	return f(p)
}

func WithSubscriptionTo(pair market.Pair) Option {
	return optionApplyFunc(func(client *baseClient) error {
		subscriptionMessage := SubscriptionMessage{
			Action: "subscribe",
			Book:   strings.ToLower(pair.Symbol("_")),
		}

		client.subscription = append(client.subscription, subscriptionMessage)
		return nil
	})
}
