package ws

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/ws"
	"github.com/vanclief/uniex/interfaces/ws/genericws"
)

const (
	ordersChannel = "book"  // We use the book for the orders
	tickerChannel = "trade" // We use the trade for the ticker
)

type KrakenHandler struct {
	opts genericws.HandlerOptions
	asks map[string]map[float64]float64
	bids map[string]map[float64]float64
}

type TradeInfo struct {
	LastPrice  float64
	LastVolume float64
	Pair       market.Pair
}

func NewHandler() *KrakenHandler {
	return &KrakenHandler{}
}

func (h *KrakenHandler) Init(opts genericws.HandlerOptions) error {
	h.opts = opts
	h.asks = make(map[string]map[float64]float64)
	h.bids = make(map[string]map[float64]float64)
	return nil
}

func getTickerArrays(in string) [][][]float64 {
	startIndex := strings.Index(in, `{`)
	endIndex := strings.LastIndex(in, `}`)
	tickerContent := in[startIndex : endIndex+1]

	splitByPoints := strings.Split(tickerContent, ":")

	var arrays [][][]float64
	for _, v := range splitByPoints {
		closingBracketIndex := strings.Index(v, "]]")
		if closingBracketIndex == -1 {
			continue
		}
		array := v[:closingBracketIndex+2]
		array = strings.ReplaceAll(array, `"`, "")
		var arrayFloat [][]float64
		err := json.Unmarshal([]byte(array), &arrayFloat)
		if err != nil {
			return nil
		}
		arrays = append(arrays, arrayFloat)
	}

	return arrays
}

func pairStringToMarketPair(in string) market.Pair {
	split := strings.Split(in, "/")
	return market.Pair{
		Base: market.Asset{
			Symbol: split[0],
		},
		Quote: market.Asset{
			Symbol: split[1],
		},
	}
}

//func processTicker(in string) (*KrakenTickerContent, market.Pair, error) {
//	const op = "KrakenHandler.ToTickers.processTicker"
//	arrays := getTickerArrays(in)
//	if len(arrays) != 9 {
//		return nil, market.Pair{}, ez.New(op, ez.EINVALID, "invalid ticker arrays length", nil)
//	}
//
//	pair := strings.Split(in, `"ticker",`)[1]
//	marketPair := pairStringToMarketPair(strings.ReplaceAll(pair[:len(pair)-1], `"`, ""))
//
//	return &KrakenTickerContent{
//		AskPrice:       arrays[0][0],
//		BidPrice:       arrays[0][1],
//		ClosePrice:     arrays[0][2],
//		Volume:         arrays[0][3],
//		VWAP:           arrays[0][4],
//		NumberOfTrades: arrays[0][5],
//		LowPrice:       arrays[0][6],
//		HighPrice:      arrays[0][7],
//		OpenPrice:      arrays[0][8],
//	}, marketPair, nil
//}

func processTrade(in string) (*TradeInfo, error) {
	const op = "KrakenHandler.ToTradeInfo"
	leftIndex := strings.Index(in, `[[`)
	rightIndex := strings.Index(in, `]]`)
	if leftIndex == -1 || rightIndex == -1 {
		return nil, ez.New(op, ez.EINVALID, "invalid trade info", nil)
	}
	var tradeBody [][]string
	err := json.Unmarshal([]byte(in[leftIndex:rightIndex+2]), &tradeBody)
	if err != nil {
		return nil, ez.New(op, ez.EINVALID, "invalid trade info", nil)
	}
	price, _ := strconv.ParseFloat(tradeBody[0][0], 64)
	volume, _ := strconv.ParseFloat(tradeBody[0][1], 64)
	pairStr := strings.Split(in, `"trade",`)
	pairWithQuotes := strings.ReplaceAll(pairStr[1][:len(pairStr[1])-1], `"`, "")
	return &TradeInfo{
		LastPrice:  price,
		LastVolume: volume,
		Pair:       pairStringToMarketPair(pairWithQuotes),
	}, nil
}

func (h *KrakenHandler) Parse(in []byte) (ws.ListenChan, error) {

	if strings.Contains(string(in), tickerChannel) {
		return h.ToTickers(in)
	} else if strings.Contains(string(in), ordersChannel) {
		return h.ToOrderBook(in)
	}

	return ws.ListenChan{}, nil
}

func (h *KrakenHandler) ToTickers(in []byte) (ws.ListenChan, error) {
	const op = "KrakenHandler.ToTickers"

	if string(in) == `{"event":"heartbeat"}` || strings.Contains(string(in), `"status":"subscribed"`) {
		return ws.ListenChan{}, nil
	}

	tradeInfo, err := processTrade(string(in))
	if err != nil {
		return ws.ListenChan{}, ez.New(op, ez.EINVALID, "invalid trade info", nil)
	}
	return ws.ListenChan{
		IsValid: true,
		Pair:    tradeInfo.Pair,
		Tickers: []market.Ticker{
			{
				Time:   time.Now().Unix(),
				Volume: tradeInfo.LastVolume,
				Last:   tradeInfo.LastPrice,
			},
		},
	}, nil
}

func (h *KrakenHandler) ToOrderBook(in []byte) (ws.ListenChan, error) {
	const op = "KrakenHandler.ToOrderBook"

	if string(in) == `{"event":"heartbeat"}` || strings.Contains(string(in), `"status":"subscribed"`) {
		return ws.ListenChan{}, nil
	}

	arrays := getTickerArrays(string(in))
	if len(arrays) == 0 {
		return ws.ListenChan{}, nil
	}

	temp := strings.Split(string(in), `,`)
	pairString := temp[len(temp)-1][:len(temp[len(temp)-1])-1]
	marketPair := pairStringToMarketPair(strings.ReplaceAll(pairString, `"`, ""))

	if _, ok := h.asks[marketPair.String()]; !ok {
		h.asks[marketPair.String()] = make(map[float64]float64)
	}

	if _, ok := h.bids[marketPair.String()]; !ok {
		h.bids[marketPair.String()] = make(map[float64]float64)
	}

	asks := arrays[0]
	var bids [][]float64
	if len(arrays) == 2 {
		bids = arrays[1]
	}

	for _, v := range asks {
		h.asks[marketPair.String()][v[0]] = v[1]
	}

	for _, v := range bids {
		h.bids[marketPair.String()][v[0]] = v[1]
	}

	accumVol := 0.0

	parsedOrderBook := market.OrderBook{
		Time: time.Now().Unix(),
	}
	for k, v := range h.asks[marketPair.String()] {
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
	for k, v := range h.bids[marketPair.String()] {
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

	return ws.ListenChan{
		Pair:      marketPair,
		OrderBook: parsedOrderBook,
		IsValid:   true,
	}, nil
}

func (h KrakenHandler) GetSettings() (genericws.Settings, error) {
	return genericws.Settings{
		Endpoint: "wss://ws.kraken.com",
	}, nil
}

func (h KrakenHandler) GetSubscriptionsRequests() ([]genericws.SubscriptionRequest, error) {
	const op = "KrakenHandler.GetSubscriptionRequests"

	var requests []genericws.SubscriptionRequest
	channelsMap := make(map[genericws.ChannelType]bool, len(h.opts.Channels))
	for _, channel := range h.opts.Channels {
		channelsMap[channel.Type] = true
	}

	pairsArray := make([]string, len(h.opts.Pairs))

	for i, pair := range h.opts.Pairs {
		pairsArray[i] = pair.String()[:len(pair.String())-1]
	}

	if channelsMap[genericws.OrderBookChannel] {
		subscriptionMessage := KrakenSubscriptionRequest{
			Event: "subscribe",
			Pair:  pairsArray,
			Subscription: KrakenSubscription{
				Name: ordersChannel,
			},
		}

		request, err := json.Marshal(subscriptionMessage)
		if err != nil {
			return nil, ez.New(op, ez.EINTERNAL, "Error parsing Subscription Message Request", err)
		}

		requests = append(requests, request)
	}

	if channelsMap[genericws.TickerChannel] {
		subscriptionMessage := KrakenSubscriptionRequest{
			Event: "subscribe",
			Pair:  pairsArray,
			Subscription: KrakenSubscription{
				Name: tickerChannel,
			},
		}

		request, err := json.Marshal(subscriptionMessage)
		if err != nil {
			return nil, ez.New(op, ez.EINTERNAL, "Error parsing Subscription Message Request", err)
		}

		requests = append(requests, request)
	}

	return requests, nil
}

func (h KrakenHandler) VerifySubscriptionResponse(in []byte) error {
	const op = "KrakenHandler.VerifySubscriptionResponse"

	if strings.Contains(string(in), `"status":"subscribed"`) {
		return nil
	}

	if strings.Contains(string(in), `"status":"online"`) {
		return nil
	}

	return ez.New(op, ez.EINTERNAL, "Error subscribing to Kraken", nil)
}
