package ws

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrUnknownSubscriptionType = errors.New("unknown subscription type")
	ErrMarketFormatError       = errors.New("market format error")
)

type optionApplyFunc func(client *baseClient) error

type Option interface {
	applyOption(client *baseClient) error
}

func (f optionApplyFunc) applyOption(p *baseClient) error {
	return f(p)
}

func WithSubscriptionTo(market string, channel channelType) Option {
	return optionApplyFunc(func(client *baseClient) error {
		if channel != "orderbook" && channel != "ticker" && channel != "trades" {
			return ErrUnknownSubscriptionType
		}
		tokens := strings.Split(market, "_")
		if len(tokens) != 2 {
			return ErrMarketFormatError
		}
		client.subscription.Channel = channel
		client.subscription.Market = fmt.Sprintf("%s-%s", strings.ToUpper(tokens[0]), strings.ToUpper(tokens[1]))
		return nil
	})
}
