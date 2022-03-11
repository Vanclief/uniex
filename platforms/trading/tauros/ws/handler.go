package ws

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/vanclief/uniex/interfaces/ws"

	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/ws/genericws"
)

const (
	ordersChannel = "orderbook"
	tickerChannel = "ticker"
)

type TaurosHandler struct{}

func NewHandler() TaurosHandler {
	return TaurosHandler{}
}

func (h TaurosHandler) Parse(in []byte) (*ws.ListenChan, error) {
	t := Type{}
	err := json.Unmarshal(in, &t)
	if err != nil {
		return nil, err
	}

	switch t.Channel {
	case "orderbook":
		ob, pair, pErr := h.toOrderBook(in)
		if pErr != nil {
			return nil, pErr
		}
		if ob != nil {
			return &ws.ListenChan{
				Type:      ws.OrderBookType,
				Pair:      *pair,
				OrderBook: *ob,
			}, nil
		}
	case "ticker":
		ticks, pair, pErr := h.toTickers(in)
		if pErr != nil {
			return nil, pErr
		}
		if ticks != nil {
			return &ws.ListenChan{
				Type:  ws.TickerType,
				Pair:  *pair,
				Ticks: ticks,
			}, nil
		}
	}

	return nil, nil
}

func (h TaurosHandler) toTickers(in []byte) ([]market.Ticker, *market.Pair, error) {
	tick := Tick{}
	err := json.Unmarshal(in, &tick)
	if err != nil {
		return nil, nil, err
	}

	pair, err := genericws.ToMarketPair(tick.Market, "-")
	if err != nil {
		return nil, nil, err
	}

	ticks := make([]market.Ticker, 0, len(tick.Trades))
	for _, trade := range tick.Trades {
		ticks = append(ticks, transformTradeToTicker(trade))
	}

	return ticks, &pair, nil
}

func (h TaurosHandler) toOrderBook(in []byte) (*market.OrderBook, *market.Pair, error) {
	order := Order{}
	err := json.Unmarshal(in, &order)
	if err != nil {
		return nil, nil, err
	}
	orderBook := market.OrderBook{
		Asks: []market.OrderBookRow{},
		Bids: []market.OrderBookRow{},
	}

	var time int64
	for _, bid := range order.Data.Bids {
		orderRow := transformToOrderBookRow(bid)
		orderBook.Bids = append(orderBook.Bids, orderRow)
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
		orderRow := transformToOrderBookRow(ask)
		orderBook.Asks = append(orderBook.Asks, orderRow)
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

	pair, err := genericws.ToMarketPair(order.Type, "-")
	if err != nil {
		return nil, nil, err
	}
	return &orderBook, &pair, nil
}

func (h TaurosHandler) GetSettings(pair []market.Pair, channels []genericws.ChannelOpts) (genericws.Settings, error) {
	return genericws.Settings{
		Endpoint: "wss://wsv2.tauros.io",
	}, nil
}

func (h TaurosHandler) GetSubscriptionsRequests(pairs []market.Pair, channels []genericws.ChannelOpts) ([]genericws.SubscriptionRequest, error) {
	const op = "handler.GetSubscriptionRequests"
	requests := make([]genericws.SubscriptionRequest, 0, len(pairs))

	channelsMap := make(map[genericws.ChannelType]bool, len(channels))
	for _, channel := range channels {
		channelsMap[channel.Type] = true
	}

	for _, pair := range pairs {
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

func (h TaurosHandler) VerifySubscriptionResponse(in []byte) error {
	const op = "TaurosHandler.VerifySubscriptionResponse"
	response := &SubscriptionResponse{}

	err := json.Unmarshal(in, &response)
	if err != nil {
		return ez.Wrap(op, err)
	}

	if response.Response != "ok" {
		return ez.New(op, ez.EINTERNAL, "error on verify subscription response", nil)
	}

	return nil
}

func transformToOrderBookRow(ba BidAsk) market.OrderBookRow {
	orderRow := market.OrderBookRow{
		Price:       ba.Price,
		Volume:      ba.Amount,
		AccumVolume: 0,
	}
	return orderRow
}

func transformTradeToTicker(ta Trade) market.Ticker {
	ticker := market.Ticker{
		Time:   int64(ta.Timestamp),
		Last:   ta.Price,
		Volume: ta.Value,
	}
	return ticker
}

func getRequest(pair market.Pair, channel string) ([]byte, error) {
	subscriptionMessage := SubscriptionMessage{
		Action:  "subscribe",
		Market:  strings.ToUpper(pair.Symbol("-")),
		Channel: channel,
	}

	request, err := json.Marshal(subscriptionMessage)
	if err != nil {
		return nil, err
	}

	return request, nil
}
