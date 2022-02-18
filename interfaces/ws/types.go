package ws

import (
	"github.com/vanclief/finmod/market"
)

type TickerChan struct {
	Pair  market.Pair
	Ticks []market.Ticker
}

type OrderBookChan struct {
	Pair      market.Pair
	OrderBook market.OrderBook
}
