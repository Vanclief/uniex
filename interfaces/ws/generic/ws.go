package generic

import (
	"context"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/ws"
)

type baseClient struct {
	host              string
	subscriptionPairs []market.Pair
	parser            WebsocketParser
	name              string
	// buffer for bursts or spikes in data
	buffer            int
	connectionRetries int
}

func New(host string, parser WebsocketParser, opts ...Option) (*baseClient, error) {
	const op = "ws.New"

	c := &baseClient{
		host:              host,
		subscriptionPairs: []market.Pair{},
		parser:            parser,
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

func (c *baseClient) createConnection(channelType ChannelType) (*websocket.Conn, error) {
	const op = "baseClient.createConnection"
	var err error
	var wsConn *websocket.Conn

	for i := 0; i < c.connectionRetries; i++ {

		wsConn, _, err = websocket.DefaultDialer.Dial(c.host, nil)
		if err == nil {
			for j := 0; j < c.connectionRetries; j++ {
				err = c.subscribeTo(channelType, wsConn)
				if err == nil {
					break
				}
			}
			break
		}
	}

	errMsg := fmt.Sprintf("Failed to connect to websocket: %s", err.Error())
	return wsConn, ez.New(op, ez.EINTERNAL, errMsg, err)
}

func (c *baseClient) subscribeTo(channelType ChannelType, ws *websocket.Conn) error {
	const op = "baseClient.subscribeTo"

	for _, pair := range c.subscriptionPairs {
		bs, err := c.parser.GetSubscriptionRequest(pair, channelType)
		if err != nil {
			msgErr := fmt.Sprintf("Error creating subscription request: %s", err.Error())
			return ez.New(op, ez.EINTERNAL, msgErr, err)
		}
		ws.WriteMessage(websocket.TextMessage, bs)
	}

	return nil
}

// ListenOrderBook returns a channel with updates to the orderbook
func (c *baseClient) ListenOrderBook(ctx context.Context) (<-chan ws.OrderBookChan, error) {
	const op = "ws.ListenOrders"

	if len(c.subscriptionPairs) == 0 {
		return nil, ez.New(op, ez.EINVALID, ErrSubscriptionPairs, nil)
	}

	wsConn, err := c.createConnection(ChannelTypeOrderBook)
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
				//fmt.Printf("\norderbook bytes: %s\n", string(bs[:10]))
				if bErr != nil {
					log.Error().
						Str("op", op).
						Str("exchange", c.name).
						Err(bErr).
						Msg("error reading orderbook data from ws")
					continue
				}
				orderBook, pErr := c.parser.ToOrderBook(bs)
				if pErr != nil {
					log.Error().
						Str("op", op).
						Str("exchange", c.name).
						Str("bytes", string(bs)).
						Err(pErr).
						Msg("error unmarshalling orderbook data from ws")
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

// ListenOrderBook returns a channel with updates to the ticker
func (c *baseClient) ListenTicker(ctx context.Context) (<-chan ws.TickerChan, error) {
	const op = "ws.ListenTicker"

	if len(c.subscriptionPairs) == 0 {
		return nil, ez.New(op, ez.EINVALID, ErrSubscriptionPairs, nil)
	}

	wsConn, cErr := c.createConnection(ChannelTypeTicker)
	if cErr != nil {
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

				tick, pErr := c.parser.ToTickers(bs)
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
