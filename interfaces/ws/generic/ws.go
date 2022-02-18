package generic

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

type baseClient struct {
	subscriptionPairs []market.Pair
	handler           WebsocketHandler
	name              string
	// buffer for bursts or spikes in data
	buffer            int
	connectionRetries int
}

func NewClient(handler WebsocketHandler, opts ...Option) (*baseClient, error) {
	const op = "ws.New"

	c := &baseClient{
		subscriptionPairs: []market.Pair{},
		handler:           handler,
		buffer:            1000,
		connectionRetries: 3,
	}
	for _, opt := range opts {
		if err := opt.applyOption(c); err != nil {
			return nil, ez.Wrap(op, err)
		}
	}

	return c, nil
}

func (c *baseClient) createConnection(ctx context.Context, channelType ChannelType) (*websocket.Conn, error) {
	const op = "baseClient.createConnection"
	var err error
	var wsConn *websocket.Conn

	baseEndpoint := c.handler.GetBaseEndpoint(c.subscriptionPairs)

	for i := 0; i < c.connectionRetries; i++ {
		wsConn, _, err = websocket.DefaultDialer.Dial(baseEndpoint, nil)
		if err == nil {
			break
		}
	}

	if err != nil {
		errMsg := fmt.Sprintf("Failed to dial websocket: %s", err.Error())
		return wsConn, ez.New(op, ez.EINTERNAL, errMsg, err)
	}

	for j := 0; j < c.connectionRetries; j++ {
		err = c.subscribe(ctx, channelType, wsConn)
		if err == nil {
			return wsConn, nil
		}
	}

	errMsg := fmt.Sprintf("Failed to subscribe to websocket channel: %s", err.Error())
	return wsConn, ez.New(op, ez.EINTERNAL, errMsg, err)
}

func (c *baseClient) subscribe(ctx context.Context, channelType ChannelType, wsConn *websocket.Conn) error {
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

	ctxDeadline, _ := context.WithTimeout(ctx, time.Second*5)

	counter := len(requests)

	for {
		select {
		case <-ctxDeadline.Done():
			return ez.New(op, ez.EINTERNAL, "Failed to get subscription confirmations", nil)

		default:
			_, bs, err := wsConn.ReadMessage()
			if err != nil {
				return ez.Wrap(op, err)
			}

			err = c.handler.VerifySubscriptionResponse(bs)
			if err != nil {
				return ez.Wrap(op, err)
			}

			counter--

			if counter == 0 {
				return nil
			}
		}
	}
}

// ListenOrderBook returns a channel with updates to the order-book
func (c *baseClient) ListenOrderBook(ctx context.Context) (<-chan ws.OrderBookChan, error) {
	const op = "ws.ListenOrders"

	if len(c.subscriptionPairs) == 0 {
		return nil, ez.New(op, ez.EINVALID, ErrSubscriptionPairs, nil)
	}

	wsConn, err := c.createConnection(ctx, ChannelTypeOrderBook)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	chanMsgs := make(chan ws.OrderBookChan, c.buffer)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(chanMsgs)
				return
			default:
				_, bs, bErr := wsConn.ReadMessage()
				//fmt.Printf("\norderbook bytes: %s\n", string(bs))
				if bErr != nil {
					log.Error().
						Str("op", op).
						Str("exchange", c.name).
						Err(bErr).
						Msg("error reading orderbook data from ws")
					continue
				}
				orderBook, pErr := c.handler.ToOrderBook(bs)
				if pErr != nil {
					log.Error().
						Str("op", op).
						Str("exchange", c.name).
						Str("bytes", string(bs)).
						Err(pErr).
						Msg("error unmarshalling order book data from ws")
					continue
				}

				if orderBook != nil {
					chanMsgs <- *orderBook
				}
			}
		}
	}()

	return chanMsgs, nil
}

// ListenTicker returns a channel with updates to the ticker
func (c *baseClient) ListenTicker(ctx context.Context, channelType ChannelType) (<-chan ws.TickerChan, error) {
	const op = "ws.ListenTicker"

	if len(c.subscriptionPairs) == 0 {
		return nil, ez.New(op, ez.EINVALID, ErrSubscriptionPairs, nil)
	}

	wsConn, cErr := c.createConnection(ctx, channelType)
	if cErr != nil {
		fmt.Println("error here", cErr.Error())
		return nil, ez.Wrap(op, cErr)
	}

	chanMsgs := make(chan ws.TickerChan, c.buffer)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(chanMsgs)
				return
			default:
				_, bs, bErr := wsConn.ReadMessage()
				//fmt.Printf("\nticker bytes%s\n", string(bs))
				if bErr != nil {
					log.Error().
						Str("op", op).
						Str("exchange", c.name).
						Err(bErr).
						Msg("error reading ticker data from ws")
					continue
				}

				tick, pErr := c.handler.ToTickers(bs)
				//fmt.Printf("\ntick: %+v\n", tick)
				if pErr != nil {
					log.Error().
						Str("op", op).
						Str("exchange", c.name).
						Str("bytes", string(bs)).
						Err(bErr).
						Msg("error unmarshalling ticker data from ws")
				}

				if tick != nil {
					chanMsgs <- *tick
				}
			}
		}
	}()

	return chanMsgs, nil
}
