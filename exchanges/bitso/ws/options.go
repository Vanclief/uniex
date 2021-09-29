package ws

import (
	"errors"
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

func WithSubscriptionTo(book string, kind subscriptionType) Option {
	return optionApplyFunc(func(client *baseClient) error {
		if kind != "orders" && kind != "diff-orders" && kind != "trades" {
			return ErrUnknownSubscriptionType
		}
		client.subscription.Book = book
		client.subscription.Type = kind
		return nil
	})
}
