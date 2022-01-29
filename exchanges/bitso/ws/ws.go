package ws

import (
	"context"
	"encoding/json"
	"sort"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/inconshreveable/log15"
	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/exchanges/ws"
)

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
		subMessage.Type = channelType(kind).String()
		bytes, err := json.Marshal(&subMessage)
		if err != nil {
			return err
		}
		ws.WriteMessage(websocket.TextMessage, bytes)
	}

	return nil
}

func (c *baseClient) ListenOrderBook(ctx context.Context) (<-chan market.OrderBook, error) {
	op := "ws.ListenOrders"


	wsClient, cErr := c.createClient(ws.OrderBookChannel)
	if cErr != nil {
		return nil, ez.Wrap(op, cErr)
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

func (c *baseClient) ListenTicker(ctx context.Context) (<-chan market.Ticker, error) {
	op := "ws.ListenTicker"

	wsClient, cErr := c.createClient(ws.TickerChannel)
	if cErr != nil {
		return nil, cErr
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
	//buyOrSell := "buy"
	//if ta.Type == 1 {
	//	buyOrSell = "sell"
	//}

	ticker := &market.Ticker{
		Time:   time.Now().UnixMilli(),
		Last:   ta.Rate,
		Volume: ta.Value,
		//Side:   buyOrSell,
	}

	return ticker
}
