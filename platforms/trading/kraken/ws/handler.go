package ws

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/ws"
	"github.com/vanclief/uniex/interfaces/ws/genericws"

	"github.com/buger/jsonparser"
)

const (
	ordersChannel = "book"  // We use the book for the orders
	tickerChannel = "trade" // We use the trade for the ticker
)

type KrakenHandler struct {
	opts genericws.HandlerOptions
	ob   map[string]market.OrderBook
}

type TradeInfo struct {
	LastPrice  float64
	LastVolume float64
	Pair       market.Pair
}

type orderBookRef struct {
	updates []market.OrderBookUpdate
}

func NewHandler() *KrakenHandler {
	return &KrakenHandler{}
}

func (h *KrakenHandler) Init(opts genericws.HandlerOptions) error {
	h.opts = opts
	h.ob = make(map[string]market.OrderBook)
	return nil
}

func parseUpdates(input []byte) ([]market.OrderBookUpdate, market.Pair) {
	obRef := orderBookRef{
		updates: []market.OrderBookUpdate{},
	}

	data, pair, err := getDataAndPair(input)
	if err != nil {
		return nil, market.Pair{}
	}

	for _, key := range []string{"a", "as", "b", "bs"} {
		if fErr := fillOrderBookUpdates(&obRef, data, key); fErr != nil {
			return nil, market.Pair{}
		}
	}

	return obRef.updates, pairStringToMarketPair(pair)
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

	if string(in) == `{"event":"heartbeat"}` || strings.Contains(string(in), `"status":"subscribed"`) {
		return ws.ListenChan{}, nil
	}

	updates, pair := parseUpdates(in)

	if _, ok := h.ob[pair.String()]; !ok {
		h.ob[pair.String()] = market.NewOrderBook(
			[]market.OrderBookRow{},
			[]market.OrderBookRow{},
			25,
		)
	}

	ob := h.ob[pair.String()]

	for _, update := range updates {
		ob.ApplyUpdate(update)
	}

	h.ob[pair.String()] = ob

	return ws.ListenChan{
		Pair:      pair,
		OrderBook: h.ob[pair.String()],
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

func getDataAndPair(input []byte) ([]byte, string, error) {
	var pair string
	var data []byte
	var err error

	if data, _, _, err = jsonparser.Get(input, "[1]"); err != nil {
		return data, pair, err
	}

	if pair, err = jsonparser.GetString(input, "[3]"); err != nil {
		return data, pair, err
	}

	return data, pair, nil
}

func fillOrderBookUpdates(orderBookRef *orderBookRef, data []byte, key string) error {
	var err error
	side := "ask"
	if key == "b" || key == "bs" {
		side = "bid"
	}

	if _, pErr := jsonparser.ArrayEach(data, func(dataItem []byte, dataType jsonparser.ValueType, offset int, err error) {
		priceStr, err := jsonparser.GetString(dataItem, "[0]")
		if err != nil {
			return
		}
		volumeStr, err := jsonparser.GetString(dataItem, "[1]")
		if err != nil {
			return
		}
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			return
		}
		volume, err := strconv.ParseFloat(volumeStr, 64)
		if err != nil {
			return
		}
		update := market.OrderBookUpdate{Price: price, Volume: volume, Side: side}
		orderBookRef.updates = append(orderBookRef.updates, update)
	}, key); pErr != nil {
		if !errors.Is(pErr, jsonparser.KeyPathNotFoundError) {
			return pErr
		}
	}
	return err
}
