package types

type optionApplyFunc func(platform *TradingPlatform) error

type Option interface {
	applyOption(platform *TradingPlatform) error
}

func (f optionApplyFunc) applyOption(p *TradingPlatform) error {
	return f(p)
}

// func WithWSSubscriptionTo(pair market.Pair) Option {
// 	return optionApplyFunc(func(client *TradingPlatform) error {
// 		client.subscriptionPairs = append(client.subscriptionPairs, pair)
// 		return nil
// 	})
// }

// func WithName(name string) Option {
// 	return optionApplyFunc(func(client *TradingPlatform) error {
// 		client.name = name
// 		return nil
// 	})
// }
