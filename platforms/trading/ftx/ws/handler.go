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

type FTXHandler struct{}

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

func ftxAskBidsToOrderBookRow(asks, bids [][]float64) (parsedAsks market.OrderBookRow, parsedBids market.OrderBookRow) {
	if len(asks) == 0 {
		parsedAsks = market.OrderBookRow{}
	} else {
		ask := asks[0]
		parsedAsks = market.OrderBookRow{
			Price:       ask[0],
			Volume:      ask[1],
			AccumVolume: ask[1],
		}
	}
	if len(bids) == 0 {
		parsedBids = market.OrderBookRow{}
	} else {
		bid := bids[0]
		parsedBids = market.OrderBookRow{
			Price:       bid[0],
			Volume:      bid[1],
			AccumVolume: bid[1],
		}
	}
	return parsedAsks, parsedBids
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

	asks, bids := ftxAskBidsToOrderBookRow(payload.Data.Asks, payload.Data.Bids)
	return &ws.OrderBookChan{
		Pair:      ftxPairToMarketPair(payload.Market),
		OrderBook: market.OrderBook{Time: time.Now().Unix(), Asks: []market.OrderBookRow{asks}, Bids: []market.OrderBookRow{bids}},
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
