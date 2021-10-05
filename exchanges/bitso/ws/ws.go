package ws

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	log "github.com/inconshreveable/log15"
	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"sort"
	"time"
)

type SubscribeConf struct {
	Book string
	Type subscriptionType
}

type subscriptionType string

var (
	ORDERS subscriptionType = "orders"
	TRADES subscriptionType = "trades"
)

type baseClient struct {
	host         string
	wsClient     *websocket.Conn
	subscription *SubscribeConf
	// buffer for bursts or spikes in data
	buffer            int
	connectionRetries int
}

func new(host string, opts ...Option) (*baseClient, error) {
	c := &baseClient{
		host: host,
		subscription: &SubscribeConf{
			Book: "btc_usd",
			Type: "orders",
		},
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

func (c *baseClient) connect() error {
	wsClient, _, err := websocket.DefaultDialer.Dial(c.host, nil)
	c.wsClient = wsClient
	return err
}

func (c *baseClient) subscribe() error {
	op := "bitso.subscribe"
	bytes, err := json.Marshal(&SubscriptionMessage{
		Action: "subscribe",
		Book:   c.subscription.Book,
		Type:   string(c.subscription.Type),
	})
	if err != nil {
		return ez.Wrap(op, err)
	}
	c.wsClient.WriteMessage(websocket.TextMessage, bytes)
	return nil
}

func (c *baseClient) ListenOrders(ctx context.Context) (<-chan market.OrderBook, error) {
	op := "bitso.ws.ListenOrders"
	if c.subscription.Type != ORDERS {
		return nil, ez.Wrap(op, errors.New("subscriptions type is not orders"))
	}

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
					orderRow, err := transformToOrderBookRow(&bid)
					if err != nil {
						log.Error("ws listen error", "op", op, "type", "transform", "error", err)
						continue
					}
					orderBook.Bids = append(orderBook.Bids, *orderRow)
					if time < bid.UnixTime {
						time = bid.UnixTime
					}
				}

				if len(orderBook.Bids) > 0 {
					sort.Slice(orderBook.Bids, func(i, j int) bool {
						return orderBook.Bids[i].Price > orderBook.Bids[j].Price
					})
					orderBook.Bids[0].AccumVolume = orderBook.Bids[0].Volume
					for i := 1; i < len(orderBook.Bids); i++ {
						orderBook.Bids[i].AccumVolume = orderBook.Bids[i-1].Volume + orderBook.Bids[i].Volume
					}
				}

				for _, ask := range order.Payload.Asks {
					orderRow, err := transformToOrderBookRow(&ask)
					if err != nil {
						log.Error("ws listen error", "op", op, "type", "transform", "error", err)
						continue
					}
					orderBook.Asks = append(orderBook.Asks, *orderRow)
					if time < ask.UnixTime {
						time = ask.UnixTime
					}
				}

				if len(orderBook.Asks) > 0 {
					sort.Slice(orderBook.Asks, func(i, j int) bool {
						return orderBook.Asks[i].Price < orderBook.Asks[j].Price
					})
					orderBook.Asks[0].AccumVolume = orderBook.Asks[0].Volume
					for i := 1; i < len(orderBook.Asks); i++ {
						orderBook.Asks[i].AccumVolume = orderBook.Asks[i-1].Volume + orderBook.Asks[i].Volume
					}
				}

				orderBook.Time = time

				chanMsgs <- orderBook
			}
		}
	}()

	return chanMsgs, nil
}

func (c *baseClient) ListenTrade(ctx context.Context) (<-chan market.Ticker, error) {
	op := "bitso.ws.ListenTrade"
	if c.subscription.Type != TRADES {
		return nil, ez.Wrap(op, errors.New("subscriptions type is not trades"))
	}

	chanMsgs := make(chan market.Ticker, c.buffer)

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

				tradeType := TradeType{}
				err = json.Unmarshal(bytes, &tradeType)
				if err != nil {
					log.Error("ws listen error", "op", op, "type", "reading msg", "error", err)
					continue
				}

				for _, trade := range tradeType.Payload {
					ticker := transformTradeToTicker(&trade)
					chanMsgs <- *ticker
				}

			}
		}
	}()

	return chanMsgs, nil
}

func transformToOrderBookRow(ba *BidAsk) (*market.OrderBookRow, error) {
	orderRow := &market.OrderBookRow{
		Price:       ba.Rate,
		Volume:      ba.Amount,
		AccumVolume: 0,
	}
	return orderRow, nil
}

func transformTradeToTicker(ta *Trade) *market.Ticker {
	buyOrSell := "buy"
	if ta.Type == 1 {
		buyOrSell = "sell"
	}

	ticker := &market.Ticker{
		Time:   time.Now().UnixMilli(),
		Last:   ta.Rate,
		Volume: ta.Value,
		Side:   buyOrSell,
	}

	return ticker
}
