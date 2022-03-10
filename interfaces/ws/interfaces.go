package ws

import (
	"context"
)

// DataAPI - Unified data interface for Websockets
type DataAPI interface {
	Listen(ctx context.Context) (<-chan ListenChan, error)
}
