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

type MEXCHandler struct {
	opts genericws.HandlerOptions
	Asks map[string]market.OrderBookRow
	Bids map[string]market.OrderBookRow
}

func NewHandler() *MEXCHandler {
	return &MEXCHandler{}
}

func (h *MEXCHandler) Init(opts genericws.HandlerOptions) error {
	h.opts = opts
	h.Asks = make(map[string]market.OrderBookRow)
	h.Bids = make(map[string]market.OrderBookRow)
	return nil
}

func (h *MEXCHandler) GetSettings() (genericws.Settings, error) {
	return genericws.Settings{
		Endpoint: "wss://contract.mexc.com/ws",
	}, nil
}

func (h *MEXCHandler) GetSubscriptionsRequests() ([]genericws.SubscriptionRequest, error) {
	const op = "MEXCHandler.GetSubscriptionsRequests"

	var subscriptions []genericws.SubscriptionRequest

	channelsMap := make(map[genericws.ChannelType]bool, len(h.opts.Channels))
	for _, channel := range h.opts.Channels {
		channelsMap[channel.Type] = true
	}

	for _, pair := range h.opts.Pairs {
		if channelsMap[genericws.OrderBookChannel] {
			request, err := getRequest(pair, "sub.depth")
			if err != nil {
				return nil, ez.Wrap(op, err)
			}
			subscriptions = append(subscriptions, request)
		}

		if channelsMap[genericws.TickerChannel] {
			request, err := getRequest(pair, "sub.ticker")
			if err != nil {
				return nil, ez.Wrap(op, err)
			}
			subscriptions = append(subscriptions, request)
		}
	}

	return subscriptions, nil
}

func (h *MEXCHandler) VerifySubscriptionResponse(in []byte) error {
	const op = "MEXCHandler.VerifySubscriptionResponse"

	if strings.Contains(string(in), `"data":"success"`) {
		return nil
	}
	return ez.New(op, ez.EINVALID, "invalid subscription response", nil)
}

func (h *MEXCHandler) Parse(in []byte) (ws.ListenChan, error) {

	t := MEXCMsg{}
	err := json.Unmarshal(in, &t)
	if err != nil {
		return ws.ListenChan{}, err
	}

	switch t.Channel {
	case "push.ticker":
		return h.toTickers(in)
	case "push.depth":
		return h.toOrderBook(in)
	}

	return ws.ListenChan{}, nil
}

func (h *MEXCHandler) toTickers(in []byte) (ws.ListenChan, error) {
	const op = "MEXCHandler.toTickers"
	payload := MEXCTickerPayload{}

	err := json.Unmarshal(in, &payload)
	if err != nil {
		return ws.ListenChan{}, ez.New(op, ez.EINVALID, "Failed to unmarshal payload", err)
	}

	pair, err := mexcPairToMarketPair(payload.Data.Symbol)
	if err != nil {
		return ws.ListenChan{}, ez.Wrap(op, err)
	}

	marketTicker := market.Ticker{
		Time:   payload.Timestamp,
		Ask:    payload.Data.Ask1,
		Bid:    payload.Data.Bid1,
		Last:   payload.Data.LastPrice,
		Volume: 0,
		VWAP:   0,
	}
	return ws.ListenChan{
		IsValid: true,
		Pair:    pair,
		Tickers: []market.Ticker{marketTicker},
	}, nil
}

func (h *MEXCHandler) toOrderBook(in []byte) (ws.ListenChan, error) {
	const op = "MEXCHandler.toOrderBook"
	payload := MEXCOrderBookPayload{}

	err := json.Unmarshal(in, &payload)
	if err != nil {
		return ws.ListenChan{}, ez.New(op, ez.EINVALID, "Failed to unmarshal payload", err)
	}

	pair, err := mexcPairToMarketPair(payload.Symbol)
	if err != nil {
		return ws.ListenChan{}, ez.Wrap(op, err)
	}

	h.updateOrderBook(payload)

	ask, ok := h.Asks[payload.Symbol]
	if !ok {
		return ws.ListenChan{}, nil
	}

	bid, ok := h.Bids[payload.Symbol]
	if !ok {
		return ws.ListenChan{}, nil
	}

	if ask.Price == 0 || bid.Price == 0 {
		return ws.ListenChan{}, nil
	}

	return ws.ListenChan{
		IsValid: true,
		Pair:    pair,
		OrderBook: market.OrderBook{
			Time: time.Now().Unix(),
			Asks: []market.OrderBookRow{ask},
			Bids: []market.OrderBookRow{bid},
		},
	}, nil
}

func (h *MEXCHandler) updateOrderBook(payload MEXCOrderBookPayload) {

	for _, v := range payload.Data.Asks {
		if v[1] > 0 {
			ask := market.OrderBookRow{
				Price:       v[0],
				Volume:      v[1],
				AccumVolume: v[1],
			}
			h.Asks[payload.Symbol] = ask
		}
	}

	for _, v := range payload.Data.Bids {
		if v[1] > 0 {
			bid := market.OrderBookRow{
				Price:       v[0],
				Volume:      v[1],
				AccumVolume: v[1],
			}
			h.Bids[payload.Symbol] = bid
		}
	}
}

func mexcPairToMarketPair(in string) (market.Pair, error) {
	const op = "mexcPairToMarketPair"
	pairArray := strings.Split(in, "_")
	if len(pairArray) != 2 {
		return market.Pair{}, ez.New(op, ez.EINVALID, "invalid pair", nil)
	}
	return market.Pair{
		Base:  market.Asset{Symbol: pairArray[0]},
		Quote: market.Asset{Symbol: pairArray[1]},
	}, nil
}

func getRequest(pair market.Pair, channel string) ([]byte, error) {

	marketSymbol := pair.Symbol("_")
	subscriptionMessage := MEXCSubscriptionRequest{
		Method: channel,
		Param: MEXCSymbol{
			Symbol: marketSymbol,
		},
	}

	request, err := json.Marshal(subscriptionMessage)
	if err != nil {
		return nil, err
	}

	return request, nil
}
