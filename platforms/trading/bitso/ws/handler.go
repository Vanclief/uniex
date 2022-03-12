package ws

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/ws"
	"github.com/vanclief/uniex/interfaces/ws/genericws"
)

const (
	ordersChannel = "orders"
	tickerChannel = "trades"
)

type bitsoHandler struct {
	opts genericws.HandlerOptions
}

func NewHandler() *bitsoHandler {
	return &bitsoHandler{}
}

func (h *bitsoHandler) Init(opts genericws.HandlerOptions) error {
	h.opts = opts
	return nil
}

func (h *bitsoHandler) Parse(in []byte) (*ws.ListenChan, error) {

	t := Type{}
	err := json.Unmarshal(in, &t)
	if err != nil {
		return nil, err
	}

	if t.Type == "" || t.Book == "" {
		return nil, nil
	}

	switch t.Type {
	case "orders":
		pair, mErr := genericws.ToMarketPair(t.Book, "_")
		if mErr != nil {
			return nil, mErr
		}
		ob, pErr := h.toOrderBook(in)
		if pErr != nil {
			return nil, pErr
		}
		if ob != nil {
			return &ws.ListenChan{
				Type:      ws.OrderBookType,
				Pair:      pair,
				OrderBook: *ob,
			}, nil
		}
	case "trades":

		pair, mErr := genericws.ToMarketPair(t.Book, "_")
		if mErr != nil {
			return nil, mErr
		}
		ticks, pErr := h.toTickers(in)
		if pErr != nil {
			return nil, pErr
		}
		if ticks != nil {
			return &ws.ListenChan{
				Type:    ws.TickerType,
				Pair:    pair,
				Tickers: ticks,
			}, nil
		}
	}

	return nil, nil
}

func (h *bitsoHandler) toTickers(in []byte) ([]market.Ticker, error) {
	tradeType := TradeType{}
	err := json.Unmarshal(in, &tradeType)
	if err != nil {
		return nil, err
	}

	ticks := make([]market.Ticker, 0, len(tradeType.Payload))
	for _, trade := range tradeType.Payload {
		ticks = append(ticks, toTicker(trade))
	}

	return ticks, nil
}

func (h *bitsoHandler) toOrderBook(in []byte) (*market.OrderBook, error) {
	order := Order{}
	err := json.Unmarshal(in, &order)
	if err != nil {
		return nil, err
	}

	orderBook := market.OrderBook{
		Asks: []market.OrderBookRow{},
		Bids: []market.OrderBookRow{},
	}

	var time int64
	for _, bid := range order.Payload.Bids {
		orderBook.Bids = append(orderBook.Bids, toOrderBookRow(bid))
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
		orderBook.Asks = append(orderBook.Asks, toOrderBookRow(ask))
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
	return &orderBook, nil
}

func (h *bitsoHandler) GetSettings() (genericws.Settings, error) {
	return genericws.Settings{
		Endpoint: "wss://ws.bitso.com",
	}, nil
}

func (h *bitsoHandler) GetSubscriptionsRequests() ([]genericws.SubscriptionRequest, error) {
	const op = "handler.GetSubscriptionRequests"
	requests := make([]genericws.SubscriptionRequest, 0, len(h.opts.Pairs))

	channelsMap := make(map[genericws.ChannelType]bool, len(h.opts.Channels))
	for _, channel := range h.opts.Channels {
		channelsMap[channel.Type] = true
	}

	for _, pair := range h.opts.Pairs {
		if channelsMap[genericws.OrderBookChannel] {
			request, err := getRequest(pair, ordersChannel)
			if err != nil {
				return nil, ez.Wrap(op, err)
			}
			requests = append(requests, request)
		}

		if channelsMap[genericws.TickerChannel] {
			request, err := getRequest(pair, tickerChannel)
			if err != nil {
				return nil, ez.Wrap(op, err)
			}
			requests = append(requests, request)
		}
	}

	return requests, nil
}

func (h *bitsoHandler) VerifySubscriptionResponse(in []byte) error {
	const op = "bitsoHandler.VerifySubscriptionResponse"
	response := &SubscriptionResponse{}

	err := json.Unmarshal(in, &response)
	if err != nil {
		return ez.Wrap(op, err)
	}

	if response.Response != "ok" {
		msg := fmt.Sprintf("Error on very subscription response: %s", response.Response)
		return ez.New(op, ez.EINTERNAL, msg, nil)
	}

	return nil
}

func toOrderBookRow(ba BidAsk) market.OrderBookRow {
	orderRow := market.OrderBookRow{
		Price:       ba.Rate,
		Volume:      ba.Amount,
		AccumVolume: 0,
	}
	return orderRow
}

func toTicker(ta Trade) market.Ticker {
	ticker := market.Ticker{
		Time:   time.Now().UnixMilli(),
		Last:   ta.Rate,
		Volume: ta.Value,
	}
	return ticker
}

func getRequest(pair market.Pair, channel string) ([]byte, error) {
	subscriptionMessage := SubscriptionMessage{
		Action: "subscribe",
		Book:   strings.ToLower(pair.Symbol("_")),
		Type:   channel,
	}

	request, err := json.Marshal(subscriptionMessage)
	if err != nil {
		return nil, err
	}

	return request, nil
}
