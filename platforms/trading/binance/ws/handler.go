package ws

import (
	"encoding/json"
	"fmt"
	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/ws"
	"github.com/vanclief/uniex/interfaces/ws/genericws"
	"github.com/vanclief/uniex/utils"
	"strconv"
	"strings"
	"time"
)

type BinanceHandler struct{}

func NewHandler() BinanceHandler {
	return BinanceHandler{}
}

func tickerToMarketTicker(sAsk, sBid, sLast, sBaseAssetVolume, sWeightedAveragePrice string, eventTime int64) market.Ticker {
	ask, _ := strconv.ParseFloat(sAsk, 64)
	bid, _ := strconv.ParseFloat(sBid, 64)
	last, _ := strconv.ParseFloat(sLast, 64)
	baseAssetVolume, _ := strconv.ParseFloat(sBaseAssetVolume, 64)
	vwapNum, _ := strconv.ParseFloat(sWeightedAveragePrice, 64)
	vwap := vwapNum / baseAssetVolume

	return market.Ticker{
		Time:   eventTime,
		Ask:    ask,
		Bid:    bid,
		Last:   last,
		Volume: baseAssetVolume,
		VWAP:   vwap,
	}
}

func pairStringToPairStruct(pair string) (market.Pair, error) {
	const op = "BinanceHandler.pairStringToPairStruct"
	var exchange utils.NeededInfo
	for i := range utils.AllExchanges {
		if utils.AllExchanges[i].Symbol == pair {
			exchange = utils.AllExchanges[i]
			break
		}
	}
	if exchange.Symbol == "" {
		return market.Pair{}, ez.New(op, ez.EINVALID, fmt.Sprintf("%s is not a valid exchange", pair), nil)
	}
	return market.Pair{
		Base: &market.Asset{
			Symbol: exchange.Base,
			Name:   exchange.Base,
		},
		Quote: &market.Asset{
			Symbol: exchange.Quote,
			Name:   exchange.Quote,
		},
	}, nil
}

func (h BinanceHandler) ToTickers(in []byte) (*ws.TickerChan, error) {
	const op = "BinanceHandler.ToTickers"
	payload := StreamTickerEvent{}
	if strings.Contains(string(in), `"result"`) {
		return nil, nil
	}

	err := json.Unmarshal(in, &payload)
	if err != nil {
		return nil, ez.New(op, ez.EINVALID, "Failed to unmarshal payload", err)
	}

	marketTicker := tickerToMarketTicker(payload.Data.BestAskQuantity, payload.Data.BestBidQuantity, payload.Data.LastQuantity, payload.Data.BaseAssetVolume, payload.Data.WeightedAveragePrice, int64(payload.Data.EventTime))
	pair, err := pairStringToPairStruct(payload.Data.Symbol)
	if err != nil {
		return nil, nil
	}
	ticks := []market.Ticker{marketTicker}
	return &ws.TickerChan{
		Pair:  pair,
		Ticks: ticks,
	}, nil
}

func (h BinanceHandler) ToOrderBook(in []byte) (*ws.OrderBookChan, error) {
	const op = "BinanceHandler.ToOrderBook"
	payload := StreamOrderBookEvent{}
	if strings.Contains(string(in), `"result"`) {
		return nil, nil
	}

	err := json.Unmarshal(in, &payload)
	if err != nil {
		return nil, ez.New(op, ez.EINVALID, "Failed to unmarshal payload", err)
	}

	askPrice, _ := strconv.ParseFloat(payload.Data.BestAskPrice, 64)
	bidPrice, _ := strconv.ParseFloat(payload.Data.BestBidPrice, 64)
	askVolume, _ := strconv.ParseFloat(payload.Data.BestAskQuantity, 64)
	bidVolume, _ := strconv.ParseFloat(payload.Data.BestBidQuantity, 64)
	orderBook := market.OrderBook{
		Time: time.Now().Unix(),
		Asks: []market.OrderBookRow{{
			Price:       askPrice,
			Volume:      askVolume,
			AccumVolume: askVolume,
		}},
		Bids: []market.OrderBookRow{{
			Price:       bidPrice,
			Volume:      bidVolume,
			AccumVolume: bidVolume,
		}},
	}
	pair, err := pairStringToPairStruct(payload.Data.Symbol)
	if err != nil {
		return nil, nil
	}
	return &ws.OrderBookChan{
		Pair:      pair,
		OrderBook: orderBook,
	}, nil
}

func (h BinanceHandler) GetBaseEndpoint(pair []market.Pair, channelType genericws.ChannelType) string {
	var pairsStr string
	if channelType == "ticker" {
		for _, singlePair := range pair {
			pairsStr += strings.ToLower(singlePair.Symbol("")) + "@ticker/"
		}
	} else {
		pairsStr = "!bookTicker/"
	}

	return fmt.Sprintf("wss://fstream.binance.com:443/stream?streams=%s", pairsStr)
}

func (h BinanceHandler) GetSubscriptionsRequests(pairs []market.Pair, _ genericws.ChannelType) ([]genericws.SubscriptionRequest, error) {
	const op = "handler.GetSubscriptionRequests"

	requests := make([]genericws.SubscriptionRequest, 0, len(pairs))

	var pairsStr []string

	for _, pair := range pairs {
		pairsStr = append(pairsStr, strings.ToLower(pair.Symbol(""))+"@ticker")
	}
	subscriptionMessage := SubscriptionRequest{
		Method: "SUBSCRIBE",
		Params: pairsStr,
		ID:     1,
	}

	byteSubscriptionMessage, err := json.Marshal(subscriptionMessage)
	if err != nil {
		return nil, ez.New(op, ez.EINTERNAL, "Error parsing Subscription Message Request", err)
	}

	requests = append(requests, byteSubscriptionMessage)

	return requests, nil
}

func (h BinanceHandler) VerifySubscriptionResponse(in []byte) error {
	const op = "binanceHandler.VerifySubscriptionResponse"

	if !strings.Contains(string(in), `"result"`) {
		return nil
	}

	response := &SubscriptionResponse{}

	err := json.Unmarshal(in, &response)
	if err != nil {
		return ez.Wrap(op, err)
	}

	if response.ID != 1 {
		return ez.New(op, ez.EINTERNAL, "Error on verify subscription response", nil)
	}

	return nil
}
