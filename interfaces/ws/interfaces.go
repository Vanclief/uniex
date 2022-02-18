package ws

import (
	"context"
)

// DataAPI - Unified data interface for Websockets
type DataAPI interface {
	ListenOrderBook(ctx context.Context) (<-chan OrderBookChan, error)
	ListenTicker(ctx context.Context) (<-chan TickerChan, error)
}
