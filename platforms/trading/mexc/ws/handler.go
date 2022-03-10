package ws

import (
	"encoding/json"
	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/ws"
	"github.com/vanclief/uniex/interfaces/ws/genericws"
	"strings"
	"time"
)

type MEXCHandler struct{}

func NewHandler() MEXCHandler {
	return MEXCHandler{}
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

func (h MEXCHandler) ToTickers(in []byte) (*ws.TickerChan, error) {
	const op = "FTXHandler.ToTickers"
	payload := MEXCTickerPayload{}

	if !strings.Contains(string(in), `"channel":"push.ticker"`) {
		return nil, nil
	}

	err := json.Unmarshal(in, &payload)
	if err != nil {
		return nil, ez.New(op, ez.EINVALID, "Failed to unmarshal payload", err)
	}

	pair, err := mexcPairToMarketPair(payload.Data.Symbol)
	marketTicker := market.Ticker{
		Time:   payload.Timestamp,
		Ask:    payload.Data.Ask1,
		Bid:    payload.Data.Bid1,
		Last:   payload.Data.LastPrice,
		Volume: payload.Data.Volume24,
		VWAP:   0,
	}
	return &ws.TickerChan{
		Pair:  pair,
		Ticks: []market.Ticker{marketTicker},
	}, nil
}

func (h MEXCHandler) ToOrderBook(in []byte) (*ws.OrderBookChan, error) {
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

	var asks []market.OrderBookRow
	for _, v := range payload.Data.Asks {
		asks = append(asks, market.OrderBookRow{
			Price:       v[0],
			Volume:      v[1],
			AccumVolume: v[1],
		})
	}
	var bids []market.OrderBookRow
	for _, v := range payload.Data.Bids {
		bids = append(bids, market.OrderBookRow{
			Price:       v[0],
			Volume:      v[1],
			AccumVolume: v[1],
		})
	}
	return &ws.OrderBookChan{
		Pair:      pair,
		OrderBook: market.OrderBook{Time: time.Now().Unix(), Asks: asks, Bids: bids},
	}, nil
}

func (h MEXCHandler) GetBaseEndpoint([]market.Pair, genericws.ChannelType) string {
	return "wss://contract.mexc.com/ws"
}

func (h MEXCHandler) GetSubscriptionsRequests(pairs []market.Pair, channelType genericws.ChannelType) ([]genericws.SubscriptionRequest, error) {
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

func (h MEXCHandler) VerifySubscriptionResponse(in []byte) error {
	const op = "MEXCHandler.VerifySubscriptionResponse"
	if strings.Contains(string(in), `"data":"success"`) {
		return nil
	}
	return ez.New(op, ez.EINVALID, "invalid subscription response", nil)
}