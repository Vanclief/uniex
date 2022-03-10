package ws

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/ws"
	"github.com/vanclief/uniex/interfaces/ws/genericws"
)

type MEXCHandler struct {
	Ask market.OrderBookRow
	Bid market.OrderBookRow
}

func NewHandler() MEXCHandler {
	return MEXCHandler{}
}

func (h *MEXCHandler) UpdateOrderBook(data MEXCOrderBookData) {
	for _, v := range data.Asks {
		if v[1] > 0 {
			h.Ask = market.OrderBookRow{
				Price:       v[0],
				Volume:      v[1],
				AccumVolume: v[1],
			}
			fmt.Println("Set ask", h.Ask)
		}
	}

	for _, v := range data.Bids {
		if v[1] > 0 {
			h.Bid = market.OrderBookRow{
				Price:       v[0],
				Volume:      v[1],
				AccumVolume: v[1],
			}
			fmt.Println("Set bid", h.Bid)
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

func (h *MEXCHandler) ToTickers(in []byte) (*ws.TickerChan, error) {
	const op = "MEXCHandler.ToTickers"
	payload := MEXCTickerPayload{}

	if !strings.Contains(string(in), `"channel":"push.ticker"`) {
		return nil, nil
	}

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
	return &ws.TickerChan{
		Pair:  pair,
		Ticks: []market.Ticker{marketTicker},
	}, nil
}

func (h *MEXCHandler) ToOrderBook(in []byte) (*ws.OrderBookChan, error) {
	const op = "MEXCHandler.ToOrderBook"
	payload := MEXCOrderBookPayload{}
	if !strings.Contains(string(in), `"channel":"push.depth"`) {
		return nil, nil
	}

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

	return &ws.OrderBookChan{
		Pair: pair,
		OrderBook: market.OrderBook{
			Time: time.Now().Unix(),
			Asks: []market.OrderBookRow{h.Ask},
			Bids: []market.OrderBookRow{h.Bid},
		},
	}, nil
}

func (h *MEXCHandler) GetBaseEndpoint([]market.Pair, genericws.ChannelType) string {
	return "wss://contract.mexc.com/ws"
}

func (h *MEXCHandler) GetSubscriptionsRequests(pairs []market.Pair, channelType genericws.ChannelType) ([]genericws.SubscriptionRequest, error) {
	const op = "MEXCHandler.GetSubscriptionsRequests"

	var subscriptions []genericws.SubscriptionRequest
	var method string
	if channelType == genericws.ChannelTypeTicker {
		method = "sub.ticker"
	} else if channelType == genericws.ChannelTypeOrderBook {
		method = "sub.depth"
	}

	for _, v := range pairs {

		marketSymbol := v.Symbol("_")
		subscriptionRequest := MEXCSubscriptionRequest{
			Method: method,
			Param: MEXCSymbol{
				Symbol: marketSymbol,
			},
		}
		byteSubscription, err := json.Marshal(subscriptionRequest)
		if err != nil {
			return nil, ez.New(op, ez.EINTERNAL, "error marshalling subscription request", err)
		}
		subscriptions = append(subscriptions, byteSubscription)
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
