package ws

import (
	"github.com/vanclief/uniex/exchanges/ws"
)

type channelType string

var (
	ordersChannel = "orders"
	tickerChannel = "trades"
)

func (c channelType) String() string {
	defaultType := ordersChannel
	if string(c) == string(ws.TickerChannel) {
		defaultType = tickerChannel
	}
	return defaultType
}
