package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/vanclief/uniex/exchanges/ws"
	"sort"

	"github.com/gorilla/websocket"
	log "github.com/inconshreveable/log15"
	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
)

type SubscribeConf struct {
	Market  string
	Channel channelType
}

type baseClient struct {
	host         string
	subscription []SubscriptionMessage
	// buffer for bursts or spikes in data
	buffer            int
	connectionRetries int
}

func New(host string, opts ...Option) (*baseClient, error) {
	c := &baseClient{
		host:              host,
		subscription:      []SubscriptionMessage{},
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

func (c *baseClient) createClient(kind ws.ChannelType) (*websocket.Conn, error) {
	var err error
	var wsClient *websocket.Conn
	for i := 0; i < c.connectionRetries; i++ {
		wsClient, _, err = websocket.DefaultDialer.Dial(c.host, nil)
		if err == nil {
			for j := 0; j < c.connectionRetries; j++ {
				err = c.subscribeTo(kind, wsClient)
				if err == nil {
					break
				}
			}
			break
		}
	}
	return wsClient, err
}

func (c *baseClient) subscribeTo(kind ws.ChannelType, ws *websocket.Conn) error {
	for _, subMessage := range c.subscription {
		subMessage.Channel = channelType(kind).String()
		bytes, err := json.Marshal(&subMessage)
		if err != nil {
			return err
		}
		ws.WriteMessage(websocket.TextMessage, bytes)
	}

	return nil
}

func (c *baseClient) ListenOrders(ctx context.Context) (<-chan market.OrderBook, error) {
	op := "tauros.ListenOrders"

	wsClient, cErr := c.createClient(ws.OrderBookChannel)
	if cErr != nil {
		return  nil, ez.Wrap(op, cErr)
	}

	chanMsgs := make(chan market.OrderBook, c.buffer)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(chanMsgs)
				return
			default:
				_, bytes, err := wsClient.ReadMessage()
				if err != nil {
					log.Error("ws listen error", "op", op, "type", "reading msg", "error", err)
					continue
				}
				order := Order{}
				err = json.Unmarshal(bytes, &order)
				if err != nil || order.Channel != ordersChannel || order.Action == "subscribe" {
					continue
				}
				orderBook := market.OrderBook{
					Asks: []market.OrderBookRow{},
					Bids: []market.OrderBookRow{},
				}

				var time int64
				for _, bid := range order.Data.Bids {
					orderRow, err := transformToOrderBookRow(&bid)
					if err != nil {
						log.Error("ws listen error", "op", op, "type", "transform", "error", err)
						continue
					}
					orderBook.Bids = append(orderBook.Bids, *orderRow)
					if time < bid.UnixMs {
						time = bid.UnixMs
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

				for _, ask := range order.Data.Asks {
					orderRow, err := transformToOrderBookRow(&ask)
					if err != nil {
						log.Error("ws listen error", "op", op, "type", "transform", "error", err)
						continue
					}
					orderBook.Asks = append(orderBook.Asks, *orderRow)
					if time < ask.UnixMs {
						time = ask.UnixMs
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

func (c *baseClient) ListenTicker(ctx context.Context) (<-chan market.Ticker, error) {
	op := "tauros.ListenTicker"

	wsClient, cErr := c.createClient(ws.TickerChannel)
	fmt.Println("err", cErr)
	if cErr != nil {
		return  nil, ez.Wrap(op, cErr)
	}

	chanMsgs := make(chan market.Ticker, c.buffer)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(chanMsgs)
				return
			default:
				_, bytes, err := wsClient.ReadMessage()
				if err != nil {
					log.Error("ws listen error", "op", op, "type", "reading msg", "error", err)
					continue
				}

				tick := Tick{}
				err = json.Unmarshal(bytes, &tick)
				if err != nil || tick.Channel != tickerChannel || tick.Action == "subscribe" {
					continue
				}

				for _, trade := range tick.Trades {
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
		Price:       ba.Price,
		Volume:      ba.Amount,
		AccumVolume: 0,
	}
	return orderRow, nil
}

func transformTradeToTicker(ta *Trade) *market.Ticker {

	ticker := &market.Ticker{
		Time:   ta.Timestamp,
		Last:   ta.Price,
		Volume: ta.Value,
	}

	return ticker
}
