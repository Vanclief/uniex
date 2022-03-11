package ws

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/ws"
	"github.com/vanclief/uniex/interfaces/ws/genericws"
)

type Handler struct{}

func NewHandler() Handler {
	return Handler{}
}

func (p Handler) Parse(in []byte) (*ws.ListenChan, error) {
	t := Type{}
	err := json.Unmarshal(in, &t)
	if err != nil {
		return nil, err
	}

	if t.Type == "error" {
		return nil, fmt.Errorf("%s", string(in))
	}

	if t.Type != "message" {
		return nil, nil
	}

	switch t.Subject {
	case "level2":
		ob, pair, pErr := p.toOrderBook(in)
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
	case "trade.ticker":
		ticks, pair, pErr := p.toTickers(in)
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

func (p Handler) GetSettings(pair []market.Pair, channels []genericws.ChannelOpts) (genericws.Settings, error) {
	accessToken, err := GetToken()
	if err != nil {
		return genericws.Settings{}, err
	}
	connectID := uuid.New().String()
	endpoint := fmt.Sprintf("%s?token=%s&connectId=%s", accessToken.Data.InstanceServers[0].Endpoint, accessToken.Data.Token, connectID)
	return genericws.Settings{
		Endpoint:                      endpoint,
		SubscriptionVerificationCount: 1,
	}, nil
}

func (p Handler) GetSubscriptionsRequests(pair []market.Pair, channels []genericws.ChannelOpts) ([]genericws.SubscriptionRequest, error) {
	const op = "kucoin.GetSubscriptionsRequests"
	var topic string
	for _, v := range pair {
		topic += fmt.Sprintf("%s-%s,", v.Base.Symbol, v.Quote.Symbol)
	}
	topic = topic[:len(topic)-1]

	channelsMap := make(map[genericws.ChannelType]bool, len(channels))
	for _, channel := range channels {
		channelsMap[channel.Type] = true
	}

	subRequests := make([]genericws.SubscriptionRequest, 0, len(channels))

	if channelsMap[genericws.TickerChannel] {
		subscriptionMessage := SubscriptionMessageRequest{
			ID:             1,
			Type:           "subscribe",
			Topic:          fmt.Sprintf("/market/ticker:%s", topic),
			PrivateChannel: false,
			Response:       true,
		}
		bsSubTicker, err := json.Marshal(subscriptionMessage)
		if err != nil {
			return nil, ez.New(op, ez.EINTERNAL, "error marshalling subscription message", err)
		}
		subRequests = append(subRequests, bsSubTicker)
	}

	if channelsMap[genericws.OrderBookChannel] {
		subscriptionMessage := SubscriptionMessageRequest{
			ID:             1,
			Type:           "subscribe",
			Topic:          fmt.Sprintf("/spotMarket/level2Depth5:%s", topic),
			PrivateChannel: false,
			Response:       true,
		}
		bsSubOrder, err := json.Marshal(subscriptionMessage)
		if err != nil {
			return nil, ez.New(op, ez.EINTERNAL, "error marshalling subscription message", err)
		}
		subRequests = append(subRequests, bsSubOrder)
	}
	
	return subRequests, nil
}

func (p Handler) VerifySubscriptionResponse(in []byte) error {
	const op = "kucoin.VerifySubscriptionResponse"

	response := &SubscriptionMessageResponse{}

	err := json.Unmarshal(in, &response)
	if err != nil {
		return ez.Wrap(op, err)
	}

	if response.Type != "welcome" {
		return ez.New(op, ez.EINTERNAL, "error on verify subscription response", nil)
	}
	return nil
}

func (p Handler) toTickers(in []byte) ([]market.Ticker, *market.Pair, error) {
	const op = "kucoin.ToTickers"

	tradeType := TradeType{}
	err := json.Unmarshal(in, &tradeType)
	if err != nil {
		return nil, nil, ez.New(op, ez.EINTERNAL, "error parsing trade type", err)
	}

	pair, err := topicToMarketPair(tradeType.Topic)
	if err != nil {
		return nil, nil, ez.New(op, ez.EINTERNAL, "error parsing topic", err)
	}

	ticks := []market.Ticker{toTicker(tradeType.Data)}

	return ticks, &pair, nil
}

func (p Handler) toOrderBook(in []byte) (*market.OrderBook, *market.Pair, error) {
	order := Order{}
	err := json.Unmarshal(in, &order)
	if err != nil {
		return nil, nil, err
	}
	orderBook := market.OrderBook{
		Time: order.Data.Timestamp,
		Asks: []market.OrderBookRow{},
		Bids: []market.OrderBookRow{},
	}

	for _, v := range order.Data.Asks {
		obRow, err := kucoinRowToOrderBookRow(v)
		if err != nil {
			return nil, nil, err
		}
		orderBook.Asks = append(orderBook.Asks, obRow)
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

	for _, v := range order.Data.Bids {
		obRow, err := kucoinRowToOrderBookRow(v)
		if err != nil {
			return nil, nil, err
		}
		orderBook.Bids = append(orderBook.Bids, obRow)
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

	pair, err := getPairFromKucoinOrder(order.Topic)
	if err != nil {
		return nil, nil, err
	}

	return &orderBook, &pair, nil
}

func topicToMarketPair(topic string) (market.Pair, error) {
	const op = "kucoin.topicToMarketPair"
	pair := strings.Split(topic, "ticker:")[1]
	if pair == "" {
		return market.Pair{}, ez.New(op, ez.EINTERNAL, "error parsing topic", nil)
	}
	baseQuote := strings.Split(pair, "-")
	if len(baseQuote) != 2 {
		return market.Pair{}, ez.New(op, ez.EINTERNAL, "error parsing topic", nil)
	}
	return market.Pair{
		Base: market.Asset{
			Symbol: baseQuote[0],
			Name:   baseQuote[0],
		},
		Quote: market.Asset{
			Symbol: baseQuote[1],
			Name:   baseQuote[1],
		},
	}, nil
}

func getPairFromKucoinOrder(topic string) (market.Pair, error) {
	const op = "kucoin.getPairFromOrderBook"
	pair := strings.Split(topic, ":")
	if len(pair) != 2 {
		return market.Pair{}, ez.New(op, ez.EINTERNAL, "error parsing topic", nil)
	}
	firstPair := strings.Split(pair[1], "-")
	return market.Pair{
		Base: market.Asset{
			Symbol: firstPair[0],
			Name:   firstPair[0],
		},
		Quote: market.Asset{
			Symbol: firstPair[1],
			Name:   firstPair[1],
		},
	}, nil
}

func kucoinRowToOrderBookRow(row []string) (market.OrderBookRow, error) {
	accumVolume := 0.0
	thisAskVol, err := strconv.ParseFloat(row[1], 64)
	if err != nil {
		return market.OrderBookRow{}, err
	}
	thisAskPrice, err := strconv.ParseFloat(row[0], 64)
	if err != nil {
		return market.OrderBookRow{}, err
	}
	accumVolume += thisAskVol
	return market.OrderBookRow{
		Price:       thisAskPrice,
		Volume:      thisAskVol,
		AccumVolume: accumVolume,
	}, nil
}

func toTicker(ta TradeTypeData) market.Ticker {

	priceFloat, err := strconv.ParseFloat(ta.Price, 64)
	if err != nil {
		priceFloat = 0
	}

	sizeFloat, err := strconv.ParseFloat(ta.Size, 64)
	if err != nil {
		sizeFloat = 0
	}

	ticker := market.Ticker{
		Time:   time.Now().UnixMilli(),
		Last:   priceFloat,
		Volume: sizeFloat,
	}
	return ticker
}
