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
	Ask  market.OrderBookRow
	Bid  market.OrderBookRow
}

func NewHandler() *MEXCHandler {
	return &MEXCHandler{}
}

func (h *MEXCHandler) Init(opts genericws.HandlerOptions) error {
	h.opts = opts
	return nil
}

func (h *MEXCHandler) UpdateOrderBook(data MEXCOrderBookData) {
	for _, v := range data.Asks {
		if v[1] > 0 {
			h.Ask = market.OrderBookRow{
				Price:       v[0],
				Volume:      v[1],
				AccumVolume: v[1],
			}
		}
	}

	for _, v := range data.Bids {
		if v[1] > 0 {
			h.Bid = market.OrderBookRow{
				Price:       v[0],
				Volume:      v[1],
				AccumVolume: v[1],
			}
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

func (h *MEXCHandler) Parse(in []byte) (*ws.ListenChan, error) {

	t := MEXCMsg{}
	err := json.Unmarshal(in, &t)
	if err != nil {
		return nil, err
	}

	switch t.Channel {
	case "push.ticker":
		return h.ToTickers(in)
	case "push.depth":
		return h.ToOrderBook(in)
	}

	return nil, nil
}

func (h *MEXCHandler) ToTickers(in []byte) (*ws.ListenChan, error) {
	const op = "MEXCHandler.ToTickers"
	payload := MEXCTickerPayload{}

	err := json.Unmarshal(in, &payload)
	if err != nil {
		return nil, ez.New(op, ez.EINVALID, "Failed to unmarshal payload", err)
	}

	pair, err := mexcPairToMarketPair(payload.Data.Symbol)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	marketTicker := market.Ticker{
		Time:   payload.Timestamp,
		Ask:    payload.Data.Ask1,
		Bid:    payload.Data.Bid1,
		Last:   payload.Data.LastPrice,
		Volume: 0,
		VWAP:   0,
	}
	return &ws.ListenChan{
		Pair:    pair,
		Tickers: []market.Ticker{marketTicker},
	}, nil
}

func (h *MEXCHandler) ToOrderBook(in []byte) (*ws.ListenChan, error) {
	const op = "MEXCHandler.ToOrderBook"
	payload := MEXCOrderBookPayload{}

	err := json.Unmarshal(in, &payload)
	if err != nil {
		return nil, ez.New(op, ez.EINVALID, "Failed to unmarshal payload", err)
	}

	pair, err := mexcPairToMarketPair(payload.Symbol)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	h.UpdateOrderBook(payload.Data)

	if h.Ask.Price == 0 || h.Bid.Price == 0 {
		return nil, nil
	}

	return &ws.ListenChan{
		Pair: pair,
		OrderBook: market.OrderBook{
			Time: time.Now().Unix(),
			Asks: []market.OrderBookRow{h.Ask},
			Bids: []market.OrderBookRow{h.Bid},
		},
	}, nil
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
