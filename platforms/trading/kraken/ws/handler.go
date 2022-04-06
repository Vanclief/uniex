package ws

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/vanclief/uniex/utils"

	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/ws"
	"github.com/vanclief/uniex/interfaces/ws/genericws"
)

const (
	ordersChannel = "book"  // We use the book for the orders
	tickerChannel = "trade" // We use the trade for the ticker
)

type KrakenHandler struct {
	opts genericws.HandlerOptions
	asks map[string]map[float64]float64
	bids map[string]map[float64]float64
}

type TradeInfo struct {
	LastPrice  float64
	LastVolume float64
	Pair       market.Pair
}

func NewHandler() *KrakenHandler {
	return &KrakenHandler{}
}

func (h *KrakenHandler) Init(opts genericws.HandlerOptions) error {
	h.opts = opts
	h.asks = make(map[string]map[float64]float64)
	h.bids = make(map[string]map[float64]float64)
	return nil
}

func generateMapFromWsInput(input []byte) (map[string][][]float64, market.Pair) {
	var msg KrakenOrderBookPayload

	if err := json.Unmarshal(input, &msg); err != nil {
		return nil, market.Pair{}
	}

	var temp map[string]interface{}
	err := json.Unmarshal(msg.Data, &temp)
	if err != nil {
		return nil, market.Pair{}
	}

	asksBids := make(map[string][][]float64)
	for k, v := range temp {

		str, _ := json.Marshal(v)
		if !strings.Contains(string(str), "[[") {
			continue
		}

		str2 := strings.ReplaceAll(string(str), `,"r"`, "")
		str2 = strings.ReplaceAll(str2, `"`, "")
		var mapItem [][]float64

		err = json.Unmarshal([]byte(str2), &mapItem)
		if err != nil {
			return nil, market.Pair{}
		}

		asksBids[k] = mapItem
	}

	return asksBids, pairStringToMarketPair(msg.Pair)
}

func pairStringToMarketPair(in string) market.Pair {
	split := strings.Split(in, "/")
	return market.Pair{
		Base: market.Asset{
			Symbol: split[0],
		},
		Quote: market.Asset{
			Symbol: split[1],
		},
	}
}

func processTrade(in string) (*TradeInfo, error) {
	const op = "KrakenHandler.ToTradeInfo"
	leftIndex := strings.Index(in, `[[`)
	rightIndex := strings.Index(in, `]]`)
	if leftIndex == -1 || rightIndex == -1 {
		return nil, ez.New(op, ez.EINVALID, "invalid trade info", nil)
	}
	var tradeBody [][]string
	err := json.Unmarshal([]byte(in[leftIndex:rightIndex+2]), &tradeBody)
	if err != nil {
		return nil, ez.New(op, ez.EINVALID, "invalid trade info", nil)
	}
	price, _ := strconv.ParseFloat(tradeBody[0][0], 64)
	volume, _ := strconv.ParseFloat(tradeBody[0][1], 64)
	pairStr := strings.Split(in, `"trade",`)
	pairWithQuotes := strings.ReplaceAll(pairStr[1][:len(pairStr[1])-1], `"`, "")
	return &TradeInfo{
		LastPrice:  price,
		LastVolume: volume,
		Pair:       pairStringToMarketPair(pairWithQuotes),
	}, nil
}

func (h *KrakenHandler) Parse(in []byte) (ws.ListenChan, error) {

	if strings.Contains(string(in), tickerChannel) {
		return h.ToTickers(in)
	} else if strings.Contains(string(in), ordersChannel) {
		return h.ToOrderBook(in)
	}

	return ws.ListenChan{}, nil
}

func (h *KrakenHandler) ToTickers(in []byte) (ws.ListenChan, error) {
	const op = "KrakenHandler.ToTickers"

	if string(in) == `{"event":"heartbeat"}` || strings.Contains(string(in), `"status":"subscribed"`) {
		return ws.ListenChan{}, nil
	}

	tradeInfo, err := processTrade(string(in))
	if err != nil {
		return ws.ListenChan{}, ez.New(op, ez.EINVALID, "invalid trade info", nil)
	}
	return ws.ListenChan{
		IsValid: true,
		Pair:    tradeInfo.Pair,
		Tickers: []market.Ticker{
			{
				Time:   time.Now().Unix(),
				Volume: tradeInfo.LastVolume,
				Last:   tradeInfo.LastPrice,
			},
		},
	}, nil
}

func (h *KrakenHandler) ToOrderBook(in []byte) (ws.ListenChan, error) {
	const op = "KrakenHandler.ToOrderBook"

	fmt.Println("===========================")
	fmt.Println("ob", string(in))

	if string(in) == `{"event":"heartbeat"}` || strings.Contains(string(in), `"status":"subscribed"`) {
		return ws.ListenChan{}, nil
	}

	asksBids, pair := generateMapFromWsInput(in)

	if asksBids == nil {
		return ws.ListenChan{}, ez.New(op, ez.EINVALID, "invalid order book", nil)
	}

	if _, ok := h.asks[pair.String()]; !ok {
		h.asks[pair.String()] = make(map[float64]float64)
	}

	if _, ok := h.bids[pair.String()]; !ok {
		h.bids[pair.String()] = make(map[float64]float64)
	}

	for k, v := range asksBids {
		if len(v) == 0 {
			continue
		}
		switch k {
		case "as":
			fallthrough
		case "a":
			for _, ask := range v {
				volume := ask[1]
				if volume == 0 {
					delete(h.asks[pair.String()], ask[0])
				} else {

					h.asks[pair.String()][ask[0]] = ask[1]
				}
			}
		case "bs":
			fallthrough
		case "b":
			for _, bid := range v {
				volume := bid[1]
				if volume == 0 {
					delete(h.bids[pair.String()], bid[0])
				} else {
					h.bids[pair.String()][bid[0]] = bid[1]
				}
			}
		}
	}

	parsedOrderBook := utils.GenerateOrderBookFromMap(h.asks[pair.String()], h.bids[pair.String()])

	return ws.ListenChan{
		Pair:      pair,
		OrderBook: parsedOrderBook,
		IsValid:   true,
	}, nil
}

func (h KrakenHandler) GetSettings() (genericws.Settings, error) {
	return genericws.Settings{
		Endpoint: "wss://ws.kraken.com",
	}, nil
}

func (h KrakenHandler) GetSubscriptionsRequests() ([]genericws.SubscriptionRequest, error) {
	const op = "KrakenHandler.GetSubscriptionRequests"

	var requests []genericws.SubscriptionRequest
	channelsMap := make(map[genericws.ChannelType]bool, len(h.opts.Channels))
	for _, channel := range h.opts.Channels {
		channelsMap[channel.Type] = true
	}

	pairsArray := make([]string, len(h.opts.Pairs))

	for i, pair := range h.opts.Pairs {
		pairsArray[i] = pair.String()[:len(pair.String())-1]
	}

	if channelsMap[genericws.OrderBookChannel] {
		subscriptionMessage := KrakenSubscriptionRequest{
			Event: "subscribe",
			Pair:  pairsArray,
			Subscription: KrakenSubscription{
				Name:  ordersChannel,
				Depth: 25,
			},
		}

		request, err := json.Marshal(subscriptionMessage)
		if err != nil {
			return nil, ez.New(op, ez.EINTERNAL, "Error parsing Subscription Message Request", err)
		}

		requests = append(requests, request)
	}

	if channelsMap[genericws.TickerChannel] {
		subscriptionMessage := KrakenSubscriptionRequest{
			Event: "subscribe",
			Pair:  pairsArray,
			Subscription: KrakenSubscription{
				Name: tickerChannel,
			},
		}

		request, err := json.Marshal(subscriptionMessage)
		if err != nil {
			return nil, ez.New(op, ez.EINTERNAL, "Error parsing Subscription Message Request", err)
		}

		requests = append(requests, request)
	}

	return requests, nil
}

func (h KrakenHandler) VerifySubscriptionResponse(in []byte) error {
	const op = "KrakenHandler.VerifySubscriptionResponse"

	if strings.Contains(string(in), `"status":"subscribed"`) {
		return nil
	}

	if strings.Contains(string(in), `"status":"online"`) {
		return nil
	}

	return ez.New(op, ez.EINTERNAL, "Error subscribing to Kraken", nil)
}
