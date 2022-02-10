package ws

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
)

type baseClient struct {
	host              string
	subscriptionPairs []market.Pair
	parser            WebsocketParser
	name string
	// buffer for bursts or spikes in data
	buffer            int
	connectionRetries int
}

func New(host string, parser WebsocketParser, opts ...Option) (*baseClient, error) {
	c := &baseClient{
		host:              host,
		subscriptionPairs: []market.Pair{},
		parser: parser,
		buffer:            1000,
		connectionRetries: 3,
	}
	for _, opt := range opts {
		if err := opt.applyOption(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *baseClient) createClient(channelType ChannelType) (*websocket.Conn, error) {
	var err error
	var wsClient *websocket.Conn
	for i := 0; i < c.connectionRetries; i++ {
		wsClient, _, err = websocket.DefaultDialer.Dial(c.host, nil)
		if err == nil {
			for j := 0; j < c.connectionRetries; j++ {
				err = c.subscribeTo(channelType, wsClient)
				if err == nil {
					break
				}
			}
			break
		}
	}
	return wsClient, err
}

func (c *baseClient) subscribeTo(channelType ChannelType, ws *websocket.Conn) error {
	for _, mPair := range c.subscriptionPairs {
		bs, err := c.parser.GetSubscriptionRequest(mPair, channelType)
		if err != nil {
			return err
		}
		ws.WriteMessage(websocket.TextMessage, bs)
	}

	return nil
}

func (c *baseClient) ListenOrderBook(ctx context.Context) (<-chan OrderBookChan, error) {
	const op = "ws.ListenOrders"

	if len(c.subscriptionPairs) == 0 {
		return nil, ErrSubscriptionPairs
	}

	wsClient, cErr := c.createClient(ChannelTypeOrderBook)
	if cErr != nil {
		return nil, ez.Wrap(op, cErr)
	}

	chanMsgs := make(chan OrderBookChan, c.buffer)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(chanMsgs)
				return
			default:
				_, bs, bErr := wsClient.ReadMessage()
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

func (c *baseClient) ListenTicker(ctx context.Context) (<-chan TickerChan, error) {
	const op = "ws.ListenTicker"

	if len(c.subscriptionPairs) == 0 {
		return nil, ErrSubscriptionPairs
	}

	wsClient, cErr := c.createClient(ChannelTypeTicker)
	if cErr != nil {
		return nil, cErr
	}

	chanMsgs := make(chan TickerChan, c.buffer)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(chanMsgs)
				return
			default:
				_, bs, bErr := wsClient.ReadMessage()
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
