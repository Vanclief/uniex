package genericws

import (
	"context"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/ws"
)

const (
	waitTimeForNewConn = 100 * time.Millisecond
)


type baseClient struct {
	subscriptionPairs []market.Pair
	handler           WebsocketHandler
	name              string
	// buffer for bursts or spikes in data
	buffer            int
	connectionRetries int
	timeout           time.Duration
}

func NewClient(handler WebsocketHandler, opts ...Option) (*baseClient, error) {
	const op = "ws.NewClient"

	c := &baseClient{
		subscriptionPairs: []market.Pair{},
		handler:           handler,
		buffer:            1000,
		connectionRetries: 3,
		timeout:           time.Second * 60,
	}
	for _, opt := range opts {
		if err := opt.applyOption(c); err != nil {
			return nil, ez.Wrap(op, err)
		}
	}

	return c, nil
}

func (c *baseClient) createConnection(ctx context.Context, channelType ChannelType) (*wsConnHandler, error) {
	const op = "baseClient.createConnection"
	var err error

	baseEndpoint := c.handler.GetBaseEndpoint(c.subscriptionPairs)
	wsHandler := NewConnHandler(baseEndpoint, c.connectionRetries)

	err = wsHandler.Connect()
	if err != nil {
		errMsg := fmt.Sprintf("Failed to dial websocket: %s", err.Error())
		return wsHandler, ez.New(op, ez.EINTERNAL, errMsg, err)
	}

	for j := 0; j < c.connectionRetries; j++ {
		err = c.subscribe(ctx, channelType, wsHandler)
		if err == nil {
			return wsHandler, nil
		}
	}

	errMsg := fmt.Sprintf("Failed to subscribe to websocket channel: %s", err.Error())
	return wsHandler, ez.New(op, ez.EINTERNAL, errMsg, err)
}

func (c *baseClient) subscribe(ctx context.Context, channelType ChannelType, wsConn *wsConnHandler) error {
	const op = "baseClient.subscribe"

	requests, err := c.handler.GetSubscriptionsRequests(c.subscriptionPairs, channelType)
	if err != nil {
		msgErr := fmt.Sprintf("Error creating subscription request: %s", err.Error())
		return ez.New(op, ez.EINTERNAL, msgErr, err)
	}

	for _, request := range requests {
		err = wsConn.WriteMessage(websocket.TextMessage, request)
		if err != nil {
			return ez.New(op, ez.EINTERNAL, "Failed to write message", err)
		}
	}

	expectedRequests := len(requests)
	return c.verifySubscriptions(ctx, wsConn, expectedRequests)
}

func (c *baseClient) verifySubscriptions(ctx context.Context, wsConn *wsConnHandler, expectedRequest int) error {
	const op = "baseClient.verifySubscriptions"

	ctxDeadline, _ := context.WithTimeout(ctx, time.Second*5)
	for {
		select {
		case <-ctxDeadline.Done():
			return ez.New(op, ez.EINTERNAL, "failed to get subscription confirmations", nil)

		default:
			bs, err := wsConn.ReadMessage()
			if err != nil {
				return ez.Wrap(op, err)
			}

			err = c.handler.VerifySubscriptionResponse(bs)
			if err != nil {
				return ez.Wrap(op, err)
			}

			expectedRequest--

			if expectedRequest == 0 {
				return nil
			}
		}
	}
}

func (c *baseClient) ListenOrderBook(ctx context.Context) (<-chan ws.OrderBookChan, error) {
	const op = "ws.ListenOrders"

	if len(c.subscriptionPairs) == 0 {
		return nil, ez.New(op, ez.EINVALID, ErrSubscriptionPairs, nil)
	}

	wsConn, err := c.createConnection(ctx, ChannelTypeOrderBook)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	orderBookChan := make(chan ws.OrderBookChan, c.buffer)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(orderBookChan)
				return
			default:
				if wsConn.IsClose() {
					time.Sleep(waitTimeForNewConn)
					wsConn, _ = c.createConnection(ctx, ChannelTypeOrderBook)
					if wsConn.IsClose() {
						continue
					}
				}
				bs, bErr := wsConn.ReadMessage()
				if bErr != nil {
					log.Printf("error reading message %s", bErr)
					continue
				}

				if bErr != nil {
					log.Error().
						Str("OP", op).
						Str("Exchange", c.name).
						Err(bErr).
						Msg("Error reading orderbook data from ws")
					continue
				}

				orderBook, err := c.handler.ToOrderBook(bs)
				if err != nil {
					log.Error().
						Str("OP", op).
						Str("exchange", c.name).
						Str("bytes", string(bs)).
						Err(err).
						Msg("error unmarshalling order book data from ws")
					continue
				}

				if orderBook != nil {
					orderBookChan <- *orderBook
				}
			}
		}
	}()

	return orderBookChan, nil
}

func (c *baseClient) ListenTicker(ctx context.Context) (<-chan ws.TickerChan, error) {
	const op = "ws.ListenTicker"

	if len(c.subscriptionPairs) == 0 {
		return nil, ez.New(op, ez.EINVALID, ErrSubscriptionPairs, nil)
	}

	wsConn, cErr := c.createConnection(ctx, ChannelTypeTicker)
	if cErr != nil {
		return nil, ez.Wrap(op, cErr)
	}

	tickerChan := make(chan ws.TickerChan, c.buffer)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(tickerChan)
				return
			default:
				if wsConn.IsClose() {
					time.Sleep(waitTimeForNewConn)
					wsConn, _ = c.createConnection(ctx, ChannelTypeTicker)
					if wsConn.IsClose() {
						continue
					}
				}
				bs, bErr := wsConn.ReadMessage()
				if bErr != nil {
					log.Printf("error reading message %s", bErr)
					continue
				}

				if bErr != nil {
					log.Error().
						Str("OP", op).
						Str("Exchange", c.name).
						Err(bErr).
						Msg("Error reading ticker data from ws")
					continue
				}

				tick, err := c.handler.ToTickers(bs)
				if err != nil {
					log.Error().
						Str("OP", op).
						Str("Exchange", c.name).
						Str("Bytes", string(bs)).
						Err(err).
						Msg("Err unmarshalling ticker data from ws")
				}

				if tick != nil {
					tickerChan <- *tick
				}
			}
		}
	}()

	return tickerChan, nil
}
