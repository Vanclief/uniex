package ws

import (
	"errors"
)

var (
	ErrUnknownSubscriptionType = errors.New("unknown subscription type")
)

type optionApplyFunc func(client *client) error

type Option interface {
	applyOption(client *client) error
}

func (f optionApplyFunc) applyOption(p *client) error {
	return f(p)
}

func WithSubscriptionTo(book string, kind string) Option {
	return optionApplyFunc(func(client *client) error {
		if kind != "orders" && kind != "diff-orders" && kind != "trades" {
			return ErrUnknownSubscriptionType
		}
		conf := SubscribeConf{
			Book: book,
			Type: kind,
		}
		client.subscriptions = append(client.subscriptions, conf)
		return nil
	})
}