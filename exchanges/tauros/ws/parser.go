package ws

import (
	"encoding/json"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/exchanges/ws"
	"sort"
	"strings"
)

const (
	ordersChannel = "orderbook"
	tickerChannel = "ticker"
)

type parser struct {}

func NewParser() parser {
	return parser{}
}

func (p parser) ToTickers(in []byte) (*ws.TickerChan, error) {
	if strings.Contains(string(in), "subscribe") {
		return nil, nil
	}
	tick := Tick{}
	err := json.Unmarshal(in, &tick)
	if err != nil || tick.Channel != tickerChannel || tick.Action == "subscribe" {
		return nil, nil
	}

	pair, err := ws.ToMarketPair(tick.Market, "-")
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

func (p parser) ToOrderBook(in []byte) (*ws.OrderBookChan, error) {
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

	pair, err :=  ws.ToMarketPair(order.Type, "-")
	if err != nil {
		return nil, err
	}
	return &ws.OrderBookChan{
		Pair: pair,
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
		Market: strings.ToUpper(pair.Symbol("-")),
		Channel: channel,
	}
	return json.Marshal(subscriptionMessage)
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
