package ws

import (
	"github.com/vanclief/finmod/market"
)

type dataType string

const (
	TickerType    dataType = "ticker"
	OrderBookType dataType = "orderbook"
)

type TickerChan struct {
	Pair  market.Pair
	Ticks []market.Ticker
}

type OrderBookChan struct {
	Pair      market.Pair
	OrderBook market.OrderBook
}

type ListenChan struct {
	Type      dataType
	Pair      market.Pair
	OrderBook market.OrderBook
	Tickers   []market.Ticker
}
