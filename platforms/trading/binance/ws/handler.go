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
)

type BinanceHandler struct {
	opts          genericws.HandlerOptions
	lastUpdateId  int64
	orderBookAsks map[string]map[float64]float64
	orderBookBids map[string]map[float64]float64
}

func (h *BinanceHandler) Init(opts genericws.HandlerOptions) error {
	h.opts = opts
	h.lastUpdateId = 0
	h.orderBookAsks = make(map[string]map[float64]float64)
	h.orderBookBids = make(map[string]map[float64]float64)
	return nil
}

func NewHandler() *BinanceHandler {
	return &BinanceHandler{}
}

func (h *BinanceHandler) Parse(in []byte) (ws.ListenChan, error) {

	if strings.Contains(string(in), `@bookTicker`) || strings.Contains(string(in), `@depth20@100ms`) {
		return h.ToOrderBook(in)
	} else if strings.Contains(string(in), `@ticker`) {
		return h.ToTickers(in)
	}

	return ws.ListenChan{}, nil
}

func (h *BinanceHandler) ToTickers(in []byte) (ws.ListenChan, error) {
	const op = "BinanceHandler.ToTickers"

	payload := StreamTickerEvent{}
	if strings.Contains(string(in), `"result"`) {
		return ws.ListenChan{}, nil
	}

	err := json.Unmarshal(in, &payload)
	if err != nil {
		return ws.ListenChan{}, ez.New(op, ez.EINVALID, "Failed to unmarshal payload", err)
	}

	marketTicker := tickerToMarketTicker(payload.Data)
	pair, err := pairStringToPairStruct(payload.Data.Symbol)
	if err != nil {
		return ws.ListenChan{}, nil
	}

	ticks := []market.Ticker{marketTicker}
	return ws.ListenChan{
		IsValid: true,
		Pair:    pair,
		Tickers: ticks,
	}, nil
}

func (h *BinanceHandler) ToOrderBook(in []byte) (ws.ListenChan, error) {
	const op = "BinanceHandler.ToOrderBook"

	if strings.Contains(string(in), `@depth20@100ms`) {
		payload := StreamPartialOrderBookEvent{}
		err := json.Unmarshal(in, &payload)
		if err != nil {
			return ws.ListenChan{}, ez.New(op, ez.EINVALID, "Failed to unmarshal payload", err)
		}

		if payload.Data.FinalUpdateID <= h.lastUpdateId {
			return ws.ListenChan{}, nil
		}
		h.lastUpdateId = payload.Data.FinalUpdateID

		h.orderBookAsks[payload.Data.Symbol] = make(map[float64]float64)
		h.orderBookBids[payload.Data.Symbol] = make(map[float64]float64)

		for _, ask := range payload.Data.Asks {
			priceFloat, _ := strconv.ParseFloat(ask[0], 64)
			volFloat, _ := strconv.ParseFloat(ask[1], 64)
			if volFloat == 0 {
				delete(h.orderBookAsks[payload.Data.Symbol], priceFloat)
			} else {
				h.orderBookAsks[payload.Data.Symbol][priceFloat] = volFloat
			}
		}
		for _, bid := range payload.Data.Bids {
			priceFloat, _ := strconv.ParseFloat(bid[0], 64)
			volFloat, _ := strconv.ParseFloat(bid[1], 64)
			if volFloat == 0 {
				delete(h.orderBookBids[payload.Data.Symbol], priceFloat)
			} else {
				h.orderBookBids[payload.Data.Symbol][priceFloat] = volFloat
			}
		}

		parsedOrderBook := utils.GenerateOrderBookFromMap(h.orderBookAsks[payload.Data.Symbol], h.orderBookBids[payload.Data.Symbol])

		pair, _ := pairStringToPairStruct(payload.Data.Symbol)

		return ws.ListenChan{
			IsValid:   true,
			Pair:      pair,
			OrderBook: parsedOrderBook,
		}, nil
	} else {
		payload := StreamUpdateOrderBookEvent{}
		err := json.Unmarshal(in, &payload)
		if err != nil {
			return ws.ListenChan{}, ez.New(op, ez.EINVALID, "Failed to unmarshal payload", err)
		}

		if payload.Data.OrderBookUpdateID <= h.lastUpdateId {
			return ws.ListenChan{}, nil
		}
		h.lastUpdateId = payload.Data.OrderBookUpdateID

		if _, ok := h.orderBookAsks[payload.Data.Symbol]; !ok {
			h.orderBookAsks[payload.Data.Symbol] = make(map[float64]float64)
		}
		if _, ok := h.orderBookBids[payload.Data.Symbol]; !ok {
			h.orderBookBids[payload.Data.Symbol] = make(map[float64]float64)
		}

		bestAskPrice, _ := strconv.ParseFloat(payload.Data.BestAskPrice, 64)
		bestAskQty, _ := strconv.ParseFloat(payload.Data.BestAskQuantity, 64)
		if bestAskQty == 0 {
			delete(h.orderBookAsks[payload.Data.Symbol], bestAskPrice)
		} else {
			h.orderBookAsks[payload.Data.Symbol][bestAskPrice] = bestAskQty
		}

		bestBidPrice, _ := strconv.ParseFloat(payload.Data.BestBidPrice, 64)
		bestBidQty, _ := strconv.ParseFloat(payload.Data.BestBidQuantity, 64)

		if bestBidQty == 0 {
			delete(h.orderBookBids[payload.Data.Symbol], bestBidPrice)
		} else {
			h.orderBookBids[payload.Data.Symbol][bestBidPrice] = bestBidQty
		}

		parsedOrderBook := utils.GenerateOrderBookFromMap(h.orderBookAsks[payload.Data.Symbol], h.orderBookBids[payload.Data.Symbol])

		pair, _ := pairStringToPairStruct(payload.Data.Symbol)

		return ws.ListenChan{
			IsValid:   true,
			Pair:      pair,
			OrderBook: parsedOrderBook,
		}, nil
	}
}

func (h *BinanceHandler) GetSettings() (genericws.Settings, error) {
	var pairsStr string

	for _, singlePair := range h.opts.Pairs {
		pairsStr += strings.ToLower(singlePair.Symbol("")) + "@ticker/"
	}

	for _, singlePair := range h.opts.Pairs {
		pairsStr += strings.ToLower(singlePair.Symbol("")) + "@bookTicker/"
	}

	for _, singlePair := range h.opts.Pairs {
		pairsStr += strings.ToLower(singlePair.Symbol("")) + "@depth20@100ms/"
	}

	return genericws.Settings{
		Endpoint: fmt.Sprintf("wss://fstream.binance.com:443/stream?streams=%s", pairsStr),
	}, nil
}

func (h *BinanceHandler) GetSubscriptionsRequests() ([]genericws.SubscriptionRequest, error) {
	const op = "handler.GetSubscriptionRequests"

	requests := make([]genericws.SubscriptionRequest, 0, len(h.opts.Pairs))

	var pairsStr []string

	for _, pair := range h.opts.Pairs {
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

func (h *BinanceHandler) VerifySubscriptionResponse(in []byte) error {
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

func tickerToMarketTicker(data BinanceTickerData) market.Ticker {
	ask, _ := strconv.ParseFloat(data.BestAskPrice, 64)
	bid, _ := strconv.ParseFloat(data.BestBidPrice, 64)
	last, _ := strconv.ParseFloat(data.LastPrice, 64)
	volume, _ := strconv.ParseFloat(data.LastQuantity, 64)

	// baseAssetVolume, _ := strconv.ParseFloat(sBaseAssetVolume, 64)
	// vwapNum, _ := strconv.ParseFloat(sWeightedAveragePrice, 64)
	// vwap := vwapNum / baseAssetVolume

	return market.Ticker{
		Time:   int64(data.EventTime),
		Ask:    ask,
		Bid:    bid,
		Last:   last,
		Volume: volume,
		// VWAP:   vwap,
	}
}

func pairStringToPairStruct(pairStr string) (market.Pair, error) {
	const op = "BinanceHandler.pairStringToPairStruct"

	pair, ok := market.PairMappings[pairStr]
	if !ok {
		return market.Pair{}, ez.New(op, ez.EINVALID, fmt.Sprintf("%s is not a valid exchange", pairStr), nil)
	}

	return pair, nil
}
