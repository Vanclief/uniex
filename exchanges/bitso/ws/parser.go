package ws

import (
	"encoding/json"
	"sort"
	"strings"
	"time"

	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/exchanges/ws"
)

const (
	ordersChannel = "orders"
	tickerChannel = "trades"
)

type parser struct{}

func NewParser() parser {
	return parser{}
}

func (p parser) ToTickers(in []byte) (*ws.TickerChan, error) {
	if strings.Contains(string(in), "subscribe") {
		return nil, nil
	}
	tradeType := TradeType{}
	err := json.Unmarshal(in, &tradeType)
	if err != nil {
		return nil, err
	}
	if tradeType.Type == "ka" {
		return nil, nil
	}

	pair, err := ws.ToMarketPair(tradeType.Book, "_")
	if err != nil {
		return nil, err
	}

	ticks := make([]market.Ticker, 0, len(tradeType.Payload))
	for _, trade := range tradeType.Payload {
		ticks = append(ticks, toTicker(trade))
	}

	return &ws.TickerChan{
		Pair:  pair,
		Ticks: ticks,
	}, nil
}

func (p parser) ToOrderBook(in []byte) (*ws.OrderBookChan, error) {
	if strings.Contains(string(in), "subscribe") {
		return nil, nil
	}
	order := Order{}
	err := json.Unmarshal(in, &order)
	if err != nil {
		return nil, err
	}
	if order.Type == "ka" {
		return nil, nil
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

	pair, err := ws.ToMarketPair(order.Book, "_")
	if err != nil {
		return nil, err
	}
	return &ws.OrderBookChan{
		Pair:      pair,
		OrderBook: orderBook,
	}, err
}

func (p parser) GetSubscriptionRequest(pair market.Pair, channelType ws.ChannelType) ([]byte, error) {
	channel := ordersChannel
	if channelType == ws.ChannelTypeTicker {
		channel = tickerChannel
	}
	subscriptionMessage := SubscriptionMessage{
		Action: "subscribe",
		Book:   strings.ToLower(pair.Symbol("_")),
		Type:   channel,
	}

	return json.Marshal(subscriptionMessage)
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
