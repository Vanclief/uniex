package ws

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	log "github.com/inconshreveable/log15"
	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
)

type SubscribeConf struct {
	Book string
	Type string
}

type client struct {
	host          string
	wsClient      *websocket.Conn
	subscriptions []SubscribeConf
	// buffer for bursts or spikes in data
	buffer            int
	connectionRetries int
}

func New(host string, opts ...Option) (*client, error) {
	c := &client{
		host:              host,
		subscriptions:     []SubscribeConf{},
		buffer:            1000,
		connectionRetries: 3,
	}
	for _, opt := range opts {
		if err := opt.applyOption(c); err != nil {
			return nil, err
		}
	}

	var err error
	for i := 0; i < c.connectionRetries; i++ {
		err = c.connect()
		if err == nil {
			for j := 0; j < c.connectionRetries; j++ {
				err = c.subscribe()
				if err == nil {
					break
				}
			}
			break
		}
	}
	return c, err
}

func (c *client) connect() error {
	wsClient, _, err := websocket.DefaultDialer.Dial(c.host, nil)
	c.wsClient = wsClient
	return err
}

func (c *client) subscribe() error {
	op := "bitso.subscribe"
	for _, subConf := range c.subscriptions {
		bytes, err := json.Marshal(&SubscriptionMessage{
			Action: "subscribe",
			Book:   subConf.Book,
			Type:   "orders",
		})
		if err != nil {
			return ez.Wrap(op, err)
		}
		c.wsClient.WriteMessage(websocket.TextMessage, bytes)
	}
	return nil
}

func (c *client) Listen(ctx context.Context) <-chan market.OrderBook {
	op := "bitso.ws.Listen"
	chanMsgs := make(chan market.OrderBook, c.buffer)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(chanMsgs)
				return
			default:
				_, bytes, err := c.wsClient.ReadMessage()
				if err != nil {
					log.Error("ws listen error", "op", op, "type", "reading msg", "error", err)
					continue
				}
				order := Order{}
				err = json.Unmarshal(bytes, &order)
				if err != nil {
					continue
				}
				orderBook := market.OrderBook{
					Asks: []market.OrderBookRow{},
					Bids: []market.OrderBookRow{},
				}

				var time int64
				for _, bid := range order.Payload.Bids {
					orderRow, err := transformToOrderBookRow(bid)
					if err != nil {
						log.Error("ws listen error", "op", op, "type", "transform", "error", err)
						continue
					}
					orderBook.Bids = append(orderBook.Bids, *orderRow)
					if time < bid.UnixTime {
						time = bid.UnixTime
					}
				}

				for _, ask := range order.Payload.Asks {
					orderRow, err := transformToOrderBookRow(ask)
					if err != nil {
						log.Error("ws listen error", "op", op, "type", "transform", "error", err)
						continue
					}
					orderBook.Asks = append(orderBook.Asks, *orderRow)
					if time < ask.UnixTime {
						time = ask.UnixTime
					}
				}

				orderBook.Time = time

				chanMsgs <- orderBook
			}
		}
	}()

	return chanMsgs
}

func transformToOrderBookRow(ba BidAsk) (*market.OrderBookRow, error) {
	orderRow := &market.OrderBookRow{
		Price:  ba.Value,
		Volume: ba.Amount,
		AccumVolume: 0,
	}
	return orderRow, nil
}
