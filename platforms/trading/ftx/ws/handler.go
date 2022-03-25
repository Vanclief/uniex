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
	ordersChannel = "orderbook"
	tickerChannel = "ticker"
)

type FTXHandler struct {
	opts          genericws.HandlerOptions
	orderBookAsks map[string]map[float64]float64
	orderBookBids map[string]map[float64]float64
}

func NewHandler() *FTXHandler {
	return &FTXHandler{}
}

func (h *FTXHandler) Init(opts genericws.HandlerOptions) error {
	h.opts = opts
	h.orderBookAsks = make(map[string]map[float64]float64)
	h.orderBookBids = make(map[string]map[float64]float64)
	return nil
}

func (h *FTXHandler) GetSettings() (genericws.Settings, error) {
	return genericws.Settings{
		Endpoint: "wss://ftx.com/ws/",
	}, nil
}

func (h *FTXHandler) GetSubscriptionsRequests() ([]genericws.SubscriptionRequest, error) {
	const op = "FTXHandler.GetSubscriptionsRequests"

	var requests []genericws.SubscriptionRequest

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

func (h FTXHandler) VerifySubscriptionResponse(in []byte) error {
	const op = "FTXHandler.VerifySubscriptionResponse"

	if strings.Contains(string(in), `"type": "error"`) {
		return ez.New(op, ez.EINVALID, "Failed to subscribe to channel", nil)
	}

	if strings.Contains(string(in), `"type": "subscribed"`) {
		return nil
	}

	response := FTXSubscribeResponse{}
	err := json.Unmarshal(in, &response)
	if err != nil {
		return ez.Wrap(op, err)
	}

	return nil
}

func (h *FTXHandler) Parse(in []byte) (ws.ListenChan, error) {

	t := FTXTickerStream{}

	err := json.Unmarshal(in, &t)
	if err != nil {
		return ws.ListenChan{}, nil
	}

	switch t.Channel {
	case tickerChannel:
		return h.toTickers(in)

	case ordersChannel:
		return h.toOrderBook(in)
	}

	return ws.ListenChan{}, nil
}

func (h *FTXHandler) toTickers(in []byte) (ws.ListenChan, error) {
	const op = "FTXHandler.toTickers"
	payload := FTXTickerStream{}

	if !strings.Contains(string(in), `"type": "update"`) {
		return ws.ListenChan{}, nil
	}

	err := json.Unmarshal(in, &payload)
	if err != nil {
		return ws.ListenChan{}, ez.New(op, ez.EINVALID, "Failed to unmarshal payload", err)
	}

	return ws.ListenChan{
		IsValid: true,
		Pair:    ftxPairToMarketPair(payload.Market),
		Tickers: []market.Ticker{ftxDataToMarketTicker(payload.Data)},
	}, nil
}

func ftxPairToMarketPair(rawPair string) market.Pair {
	pair := strings.Split(rawPair, "/")

	if len(pair) != 2 {
		pair = strings.Split(rawPair, "-")
	}

	return market.Pair{
		Base:  market.Asset{Symbol: pair[0]},
		Quote: market.Asset{Symbol: pair[1]},
	}
}

func ftxDataToMarketTicker(data FTXTickerData) market.Ticker {
	return market.Ticker{
		Time:   int64(data.Time),
		Ask:    data.Ask,
		Bid:    data.Bid,
		Last:   data.Last,
		Volume: data.BidSize,
		VWAP:   ((data.Ask * data.AskSize) + (data.Bid * data.BidSize)) / (data.AskSize + data.BidSize),
	}
}

func (h *FTXHandler) toOrderBook(in []byte) (ws.ListenChan, error) {
	const op = "FTXHandler.toOrderBook"

	stream := FTXOrderBookStream{}

	if strings.Contains(string(in), `"type": "partial"`) {
		fmt.Println(string(in))
	}

	err := json.Unmarshal(in, &stream)
	if err != nil {
		return ws.ListenChan{}, ez.New(op, ez.EINVALID, "Failed to unmarshal payload", err)
	}

	orderBook := h.ftxAskBidsToOrderBookRow(stream)

	//asks, ok := h.Asks[stream.Market]
	//if !ok {
	//	return ws.ListenChan{}, nil
	//}
	//
	//bids, ok := h.Bids[stream.Market]
	//if !ok {
	//	return ws.ListenChan{}, nil
	//}
	//
	//if asks[0].Price == 0 || bids[0].Price == 0 {
	//	return ws.ListenChan{}, nil
	//}

	return ws.ListenChan{
		IsValid:   true,
		Pair:      ftxPairToMarketPair(stream.Market),
		OrderBook: orderBook,
	}, nil
}

func (h *FTXHandler) ftxAskBidsToOrderBookRow(stream FTXOrderBookStream) market.OrderBook {

	asks := stream.Data.Asks
	bids := stream.Data.Bids

	if _, ok := h.orderBookAsks[stream.Market]; !ok {
		h.orderBookAsks[stream.Market] = make(map[float64]float64)
	}

	if _, ok := h.orderBookBids[stream.Market]; !ok {
		h.orderBookBids[stream.Market] = make(map[float64]float64)
	}

	for _, v := range asks {
		volume := v[1]
		if volume > 0 {
			h.orderBookAsks[stream.Market][v[0]] = v[1]
		}
	}

	for _, v := range bids {
		volume := v[1]
		h.orderBookBids[stream.Market][v[0]] = volume
	}

	accumVol := 0.0

	parsedOrderBook := market.OrderBook{
		Time: time.Now().Unix(),
	}

	for k, v := range h.orderBookAsks[stream.Market] {
		accumVol += v
		parsedOrderBook.Asks = append(parsedOrderBook.Asks, market.OrderBookRow{
			Price:       k,
			Volume:      v,
			AccumVolume: accumVol,
		})
	}
	sort.Slice(parsedOrderBook.Asks, func(i, j int) bool {
		return parsedOrderBook.Asks[i].Price < parsedOrderBook.Asks[j].Price
	})

	accumVol = 0
	for k, v := range h.orderBookBids[stream.Market] {
		accumVol += v
		parsedOrderBook.Bids = append(parsedOrderBook.Bids, market.OrderBookRow{
			Price:       k,
			Volume:      v,
			AccumVolume: accumVol,
		})
	}
	sort.Slice(parsedOrderBook.Bids, func(i, j int) bool {
		return parsedOrderBook.Bids[i].Price > parsedOrderBook.Bids[j].Price
	})

	return parsedOrderBook
}

func getRequest(pair market.Pair, channel string) ([]byte, error) {
	const op = "getRequest"

	marketSymbol := pair.Symbol("/")
	if pair.Quote.Symbol == "PERP" {
		marketSymbol = pair.Symbol("-")
	}

	subscriptionRequest := FTXSubscribeRequest{
		Operation: "subscribe",
		Channel:   channel,
		Market:    marketSymbol,
	}

	request, err := json.Marshal(subscriptionRequest)
	if err != nil {
		return nil, ez.New(op, ez.EINTERNAL, "error marshalling subscription request", err)
	}

	return request, nil
}
