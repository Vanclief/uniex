package ws

import (
	"context"
)

// PublicWS - Unified data interface for Websockets
type PublicWS interface {
	Listen(ctx context.Context) (<-chan ListenChan, error)
}
