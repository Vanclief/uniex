package ws

import (
	"github.com/vanclief/finmod/market"
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
		client.subscriptionPairs = append(client.subscriptionPairs, pair)
		return nil
	})
}

func WithName(name string) Option {
	return optionApplyFunc(func(client *baseClient) error {
		client.name = name
		return nil
	})
}