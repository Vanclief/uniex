package ws

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
)

type SubscribeConf struct {
	Book string
	Type string
}

type client struct {
	host string
	wsClient *websocket.Conn
	subscriptions []SubscribeConf
	// buffer for bursts or spikes in data
	buffer int
}

func New(host string, opts ...Option) (*client, error) {
	c := &client{
		host:          host,
		subscriptions: []SubscribeConf{},
		buffer:        1000,
	}
	for _, opt := range opts {
		if err := opt.applyOption(c); err != nil {
			return nil, err
		}
	}

	c.connect()
	err := c.subscribe()

	return c, err
}

func (c *client) connect() error {
	// TODO: validate response to identify errors
	wsClient, _,  err :=  websocket.DefaultDialer.Dial(c.host, nil)
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

func (c *client) Listen(ctx context.Context) <- chan market.OrderBook {
	chanMsgs := make(chan market.OrderBook, c.buffer)

	go func() {
		for  {
			select {
			case <-ctx.Done():
				close(chanMsgs)
				return
			default:
				_, bytes, err := c.wsClient.ReadMessage()
				if err != nil {
					// TODO: error should be handle here or in other part
					continue
				}
				order := Order{}
				err = json.Unmarshal(bytes, &order)
				if err != nil {
					continue
				}
				orderBook := market.OrderBook{
					// TODO: fill with? default?
					Time: 0,
					Asks: []market.OrderBookRow{
					},
					Bids: []market.OrderBookRow{
					},
				}

				for _, bid := range order.Payload.Bids {
					orderRow, err := transformToOrderBookRow(bid)
					if err != nil {
						// TODO: handler errors
						continue
					}
					orderBook.Bids = append(orderBook.Bids, *orderRow)
				}

				for _, ask := range order.Payload.Asks {
					orderRow, err := transformToOrderBookRow(ask)
					if err != nil {
						// TODO: handler errors
						continue
					}
					orderBook.Asks = append(orderBook.Asks, *orderRow)
				}
				chanMsgs <- orderBook
			}
		}
	}()

	return chanMsgs
}

func transformToOrderBookRow(ba BidAsk) (*market.OrderBookRow, error){
	orderRow := &market.OrderBookRow{
		Price:       ba.Value,
		Volume:      ba.Amount,
		// TODO: fill with? default?
		AccumVolume: 0,
	}
	return orderRow, nil
}