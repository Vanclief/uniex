package ws

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

type InstanceServer struct {
	Endpoint     string `json:"endpoint"`
	Encrypt      bool   `json:"encrypt"`
	Protocol     string `json:"protocol"`
	PingInterval int    `json:"pingInterval"`
	PingTimeout  int    `json:"pingTimeout"`
}

type TokenData struct {
	Token           string           `json:"token"`
	InstanceServers []InstanceServer `json:"instanceServers"`
}

type Token struct {
	Code string
	Data TokenData
}

func GetToken() (foundToken *Token, err error) {
	const op = "kucoin.GetToken"
	endpoint := "https://api.kucoin.com/api/v1/bullet-public"
	resp, err := http.Post(endpoint, "application/json", nil)
	if err != nil {
		return nil, ez.New(op, ez.EINTERNAL, "error obtaining token", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, ez.New(op, ez.EINTERNAL, "Error reading the response", err)
	}

	err = json.Unmarshal(body, &foundToken)
	if err != nil {
		return nil, ez.New(op, ez.EINTERNAL, "Error parsing the token", err)
	}
	return foundToken, nil
}

func (p Handler) GetBaseEndpoint(pair []market.Pair) string {
	token, err := GetToken()
	if err != nil {
		return ""
	}
	connectID := uuid.New().String()
	return fmt.Sprintf("%s?token=%s&connectId=%s", token.Data.InstanceServers[0].Endpoint, token.Data.Token, connectID)
}

func (p Handler) GetSubscriptionsRequests(pair []market.Pair, channelType genericws.ChannelType) ([]genericws.SubscriptionRequest, error) {
	const op = "kucoin.GetSubscriptionsRequests"
	var topic string
	for _, v := range pair {
		topic += fmt.Sprintf("%s-%s,", v.Base.Symbol, v.Quote.Symbol)
	}
	topic = topic[:len(topic)-1]
	var subscriptionMessage SubscriptionMessageRequest
	switch channelType {
	case genericws.ChannelTypeTicker:
		subscriptionMessage = SubscriptionMessageRequest{
			ID:             1,
			Type:           "subscribe",
			Topic:          fmt.Sprintf("/market/ticker:%s", topic),
			PrivateChannel: false,
			Response:       true,
		}
	case genericws.ChannelTypeOrderBook:
		subscriptionMessage = SubscriptionMessageRequest{
			ID:             1,
			Type:           "subscribe",
			Topic:          fmt.Sprintf("/spotMarket/level2Depth5:%s", topic),
			PrivateChannel: false,
			Response:       true,
		}
	}

	byteSubscriptionMessage, err := json.Marshal(subscriptionMessage)
	if err != nil {
		return nil, ez.New(op, ez.EINTERNAL, "error marshalling subscription message", err)
	}

	return []genericws.SubscriptionRequest{byteSubscriptionMessage}, nil
}

func (p Handler) VerifySubscriptionResponse(in []byte) error {
	const op = "kucoin.VerifySubscriptionResponse"

	response := &SubscriptionMessageResponse{}

	err := json.Unmarshal(in, &response)
	if err != nil {
		return ez.Wrap(op, err)
	}

	if response.Type != "welcome" {
		return ez.New(op, ez.EINTERNAL, "Error on verify subscription response", nil)
	}
	return nil
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
		Base: &market.Asset{
			Symbol: baseQuote[0],
			Name:   baseQuote[0],
		},
		Quote: &market.Asset{
			Symbol: baseQuote[1],
			Name:   baseQuote[1],
		},
	}, nil
}

func (p Handler) ToTickers(in []byte) (*ws.TickerChan, error) {
	const op = "kucoin.ToTickers"
	if strings.Contains(string(in), "subscribe") {
		return nil, nil
	}
	if strings.Contains(string(in), `"type":"error"`) {
		return nil, nil
	}

	tradeType := TradeType{}
	err := json.Unmarshal(in, &tradeType)
	if err != nil {
		return nil, ez.New(op, ez.EINTERNAL, "error parsing trade type", err)
	}
	if tradeType.Type == "ack" {
		return nil, nil
	}

	pair, err := topicToMarketPair(tradeType.Topic)
	if err != nil {
		return nil, ez.New(op, ez.EINTERNAL, "error parsing topic", err)
	}

	ticks := []market.Ticker{toTicker(tradeType.Data)}

	return &ws.TickerChan{
		Pair:  pair,
		Ticks: ticks,
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
		Base: &market.Asset{
			Symbol: firstPair[0],
			Name:   firstPair[0],
		},
		Quote: &market.Asset{
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

func (p Handler) ToOrderBook(in []byte) (*ws.OrderBookChan, error) {
	if strings.Contains(string(in), "subscribe") {
		return nil, nil
	}
	if strings.Contains(string(in), `"type":"error"`) {
		return nil, nil
	}
	order := Order{}
	err := json.Unmarshal(in, &order)
	if err != nil {
		return nil, err
	}
	if order.Type == "ack" {
		return nil, nil
	}
	orderBook := market.OrderBook{
		Time: order.Data.Timestamp,
		Asks: []market.OrderBookRow{},
		Bids: []market.OrderBookRow{},
	}

	for _, v := range order.Data.Asks {
		obRow, err := kucoinRowToOrderBookRow(v)
		if err != nil {
			return nil, err
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
			return nil, err
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
		return nil, err
	}

	return &ws.OrderBookChan{
		Pair:      pair,
		OrderBook: orderBook,
	}, nil
}

func NewHandler() Handler {
	return Handler{}
}

func (p Handler) GetSubscriptionRequest(pair market.Pair, channelType genericws.ChannelType) ([]byte, error) {
	var subscriptionMessage SubscriptionMessageRequest
	switch channelType {
	case genericws.ChannelTypeTicker:
		subscriptionMessage = SubscriptionMessageRequest{
			ID:             1,
			Type:           "subscribe",
			Topic:          fmt.Sprintf("/market/ticker:%s-%s", pair.Base.Symbol, pair.Quote.Symbol),
			PrivateChannel: false,
			Response:       true,
		}
	case genericws.ChannelTypeOrderBook:
		subscriptionMessage = SubscriptionMessageRequest{
			ID:             1,
			Type:           "subscribe",
			Topic:          fmt.Sprintf("/spotMarket/level2Depth5:%s-%s", pair.Base.Symbol, pair.Quote.Symbol),
			PrivateChannel: false,
			Response:       true,
		}
	}

	return json.Marshal(subscriptionMessage)
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
