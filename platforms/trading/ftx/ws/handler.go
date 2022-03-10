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

type FTXHandler struct {
	Ask market.OrderBookRow
	Bid market.OrderBookRow
}

func NewHandler() FTXHandler {
	return FTXHandler{}
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

func (h FTXHandler) ToTickers(in []byte) (*ws.TickerChan, error) {
	const op = "FTXHandler.ToTickers"
	payload := FTXTickerStream{}

	if !strings.Contains(string(in), `"type": "update"`) {
		return nil, nil
	}

	err := json.Unmarshal(in, &payload)
	if err != nil {
		return nil, ez.New(op, ez.EINVALID, "Failed to unmarshal payload", err)
	}

	return &ws.TickerChan{
		Pair:  ftxPairToMarketPair(payload.Market),
		Ticks: []market.Ticker{ftxDataToMarketTicker(payload.Data)},
	}, nil
}

func (h *FTXHandler) ftxAskBidsToOrderBookRow(asks, bids [][]float64) {
	if len(asks) > 0 {
		for _, ask := range asks {
			volume := ask[1]
			if volume > 0 {
				h.Ask = market.OrderBookRow{
					Price:       ask[0],
					Volume:      ask[1],
					AccumVolume: ask[1],
				}

			}
		}
	}

	if len(bids) > 0 {
		for _, bid := range bids {
			volume := bid[1]
			if volume > 0 {
				h.Bid = market.OrderBookRow{
					Price:       bid[0],
					Volume:      bid[1],
					AccumVolume: bid[1],
				}

			}
		}
	}
}

func (h FTXHandler) ToOrderBook(in []byte) (*ws.OrderBookChan, error) {
	const op = "BinanceHandler.ToOrderBook"

	payload := FTXOrderBookStream{}
	if strings.Contains(string(in), `"type": "partial"`) {
		return nil, nil
	}

	err := json.Unmarshal(in, &payload)
	if err != nil {
		return nil, ez.New(op, ez.EINVALID, "Failed to unmarshal payload", err)
	}

	h.ftxAskBidsToOrderBookRow(payload.Data.Asks, payload.Data.Bids)

	if h.Ask.Price == 0 || h.Bid.Price == 0 {
		return nil, nil
	}

	return &ws.OrderBookChan{
		Pair: ftxPairToMarketPair(payload.Market),
		OrderBook: market.OrderBook{
			Time: time.Now().Unix(),
			Asks: []market.OrderBookRow{h.Ask},
			Bids: []market.OrderBookRow{h.Bid},
		},
	}, nil
}

func (h FTXHandler) GetBaseEndpoint(pair []market.Pair, channelType genericws.ChannelType) string {
	return "wss://ftx.com/ws/"
}

func (h FTXHandler) GetSubscriptionsRequests(pairs []market.Pair, channelType genericws.ChannelType) ([]genericws.SubscriptionRequest, error) {
	const op = "FTXHandler.GetSubscriptionsRequests"

	var subscriptions []genericws.SubscriptionRequest

	for _, v := range pairs {

		market := v.Symbol("/")
		if v.Quote.Symbol == "PERP" {
			market = v.Symbol("-")
		}

		subscriptionRequest := FTXSubscribeRequest{
			Operation: "subscribe",
			Channel:   string(channelType),
			Market:    market,
		}
		byteSubscription, err := json.Marshal(subscriptionRequest)
		if err != nil {
			return nil, ez.New(op, ez.EINTERNAL, "error marshalling subscription request", err)
		}
		subscriptions = append(subscriptions, byteSubscription)
	}

	return subscriptions, nil
}

func (h FTXHandler) VerifySubscriptionResponse(in []byte) error {
	const op = "FTXHandler.VerifySubscriptionResponse"

	if strings.Contains(string(in), `"type": "subscribed"`) {
		return nil
	}
	return ez.New(op, ez.EINVALID, "invalid subscription response", nil)
}
