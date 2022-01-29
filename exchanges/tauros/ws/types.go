package ws

import (
	"github.com/vanclief/uniex/exchanges/ws"
)

type channelType string

var (
	ordersChannel = "orderbook"
	tickerChannel = "ticker"
)

func (c channelType) String() string {
	defaultType := ordersChannel
	if string(c) == string(ws.TickerChannel) {
		defaultType = tickerChannel
	}
	return defaultType
}