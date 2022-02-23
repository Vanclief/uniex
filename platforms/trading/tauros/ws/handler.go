package ws

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/ws"
	"github.com/vanclief/uniex/interfaces/ws/generic"
)

const (
	ordersChannel = "orderbook"
	tickerChannel = "ticker"
)

type TaurosHandler struct{}

func NewHandler() TaurosHandler {
	return TaurosHandler{}
}

func (h TaurosHandler) ToTickers(in []byte) (*ws.TickerChan, error) {
	if strings.Contains(string(in), "subscribe") {
		return nil, nil
	}
	tick := Tick{}
	err := json.Unmarshal(in, &tick)
	if err != nil || tick.Channel != tickerChannel || tick.Action == "subscribe" {
		return nil, nil
	}

	pair, err := generic.ToMarketPair(tick.Market, "-")
	if err != nil {
		return nil, err
	}

	ticks := make([]market.Ticker, 0, len(tick.Trades))
	for _, trade := range tick.Trades {
		ticks = append(ticks, transformTradeToTicker(trade))
	}

	return &ws.TickerChan{
		Pair:  pair,
		Ticks: ticks,
	}, nil
}

func (h TaurosHandler) ToOrderBook(in []byte) (*ws.OrderBookChan, error) {
	if strings.Contains(string(in), "subscribe") {
		return nil, nil
	}
	order := Order{}
	err := json.Unmarshal(in, &order)
	if err != nil || order.Channel != ordersChannel || order.Action == "subscribe" {
		return nil, nil
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

	pair, err := generic.ToMarketPair(order.Type, "-")
	if err != nil {
		return nil, err
	}
	return &ws.OrderBookChan{
		Pair:      pair,
		OrderBook: orderBook,
	}, err
}

func (h TaurosHandler) GetBaseEndpoint(pair []market.Pair) string {
	return "wss://wsv2.tauros.io"
}

func (h TaurosHandler) GetSubscriptionsRequests(pairs []market.Pair, channelType generic.ChannelType) ([]generic.SubscriptionRequest, error) {
	const op = "handler.GetSubscriptionRequests"

	requests := make([]generic.SubscriptionRequest, 0, len(pairs))

	for _, pair := range pairs {

		channel := ordersChannel

		if channelType == generic.ChannelTypeTicker {
			channel = tickerChannel
		}

		subscriptionMessage := SubscriptionMessage{
			Action:  "subscribe",
			Market:  strings.ToUpper(pair.Symbol("-")),
			Channel: channel,
		}

		request, err := json.Marshal(subscriptionMessage)
		if err != nil {
			return nil, ez.Wrap(op, err)
		}

		requests = append(requests, request)
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

	// if response.Response != "ok" {
	// 	return ez.New(op, ez.EINTERNAL, "Error on verify subscription response", nil)
	// }

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
		Time:   ta.Timestamp,
		Last:   ta.Price,
		Volume: ta.Value,
	}
	return ticker
}
