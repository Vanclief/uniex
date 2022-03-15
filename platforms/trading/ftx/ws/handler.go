package ws

import (
	"encoding/json"
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
	opts genericws.HandlerOptions
	Asks map[string]market.OrderBookRow
	Bids map[string]market.OrderBookRow
}

func NewHandler() *FTXHandler {
	return &FTXHandler{}
}

func (h *FTXHandler) Init(opts genericws.HandlerOptions) error {
	h.opts = opts
	h.Asks = make(map[string]market.OrderBookRow)
	h.Bids = make(map[string]market.OrderBookRow)
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
		return ws.ListenChan{}, nil
	}

	err := json.Unmarshal(in, &stream)
	if err != nil {
		return ws.ListenChan{}, ez.New(op, ez.EINVALID, "Failed to unmarshal payload", err)
	}

	h.ftxAskBidsToOrderBookRow(stream)

	ask, ok := h.Asks[stream.Market]
	if !ok {
		return ws.ListenChan{}, nil
	}

	bid, ok := h.Bids[stream.Market]
	if !ok {
		return ws.ListenChan{}, nil
	}

	if ask.Price == 0 || bid.Price == 0 {
		return ws.ListenChan{}, nil
	}

	return ws.ListenChan{
		IsValid: true,
		Pair:    ftxPairToMarketPair(stream.Market),
		OrderBook: market.OrderBook{
			Time: time.Now().Unix(),
			Asks: []market.OrderBookRow{ask},
			Bids: []market.OrderBookRow{bid},
		},
	}, nil
}

func (h *FTXHandler) ftxAskBidsToOrderBookRow(stream FTXOrderBookStream) {

	asks := stream.Data.Asks
	bids := stream.Data.Bids

	if len(asks) > 0 {
		for _, ask := range asks {
			volume := ask[1]
			if volume > 0 {
				ask := market.OrderBookRow{
					Price:       ask[0],
					Volume:      ask[1],
					AccumVolume: ask[1],
				}
				h.Asks[stream.Market] = ask
			}
		}
	}

	if len(bids) > 0 {
		for _, bid := range bids {
			volume := bid[1]
			if volume > 0 {
				bid := market.OrderBookRow{
					Price:       bid[0],
					Volume:      bid[1],
					AccumVolume: bid[1],
				}
				h.Bids[stream.Market] = bid
			}
		}
	}
}

func getRequest(pair market.Pair, channel string) ([]byte, error) {
	const op = "getRequest"

	market := pair.Symbol("/")
	if pair.Quote.Symbol == "PERP" {
		market = pair.Symbol("-")
	}

	subscriptionRequest := FTXSubscribeRequest{
		Operation: "subscribe",
		Channel:   string(channel),
		Market:    market,
	}

	request, err := json.Marshal(subscriptionRequest)
	if err != nil {
		return nil, ez.New(op, ez.EINTERNAL, "error marshalling subscription request", err)
	}

	return request, nil
}
