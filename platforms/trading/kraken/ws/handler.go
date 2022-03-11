package ws

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/ws"
	"github.com/vanclief/uniex/interfaces/ws/genericws"
)

const (
	ordersChannel = "ticker" // We use the ticker for the orders
	tickerChannel = "trade"  // We use the trade for the ticker
)

type KrakenHandler struct {
	AskPrice  float64
	AskVolume float64
	BidPrice  float64
	BidVolume float64
}

type TradeInfo struct {
	LastPrice  float64
	LastVolume float64
	Pair       market.Pair
}

func NewHandler() KrakenHandler {
	return KrakenHandler{}
}

func getTickerArrays(in string) [][]float64 {
	startIndex := strings.Index(in, `{`)

	var tickerContent string
	for i := startIndex; i < len(in); i++ {
		if in[i] == '}' {
			tickerContent = in[startIndex : i+1]
			break
		}
	}

	splitByPoints := strings.Split(tickerContent, ":")

	var arrays [][]float64
	for _, v := range splitByPoints {
		closingBracketIndex := strings.Index(v, "]")
		if closingBracketIndex == -1 {
			continue
		}
		array := v[:closingBracketIndex+1]
		array = strings.ReplaceAll(array, `"`, "")
		var arrayFloat []float64
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

func processTicker(in string) (*KrakenTickerContent, market.Pair, error) {
	const op = "KrakenHandler.ToTickers.processTicker"
	arrays := getTickerArrays(in)
	if len(arrays) != 9 {
		return nil, market.Pair{}, ez.New(op, ez.EINVALID, "invalid ticker arrays length", nil)
	}

	pair := strings.Split(in, `"ticker",`)[1]
	marketPair := pairStringToMarketPair(strings.ReplaceAll(pair[:len(pair)-1], `"`, ""))

	return &KrakenTickerContent{
		AskPrice:       arrays[0],
		BidPrice:       arrays[1],
		ClosePrice:     arrays[2],
		Volume:         arrays[3],
		VWAP:           arrays[4],
		NumberOfTrades: arrays[5],
		LowPrice:       arrays[6],
		HighPrice:      arrays[7],
		OpenPrice:      arrays[8],
	}, marketPair, nil
}

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

func (h *KrakenHandler) Parse(in []byte) (*ws.ListenChan, error) {

	if strings.Contains(string(in), tickerChannel) {
		return h.ToTickers(in)
	} else if strings.Contains(string(in), ordersChannel) {
		return h.ToOrderBook(in)
	}

	return nil, nil
}

func (h *KrakenHandler) ToTickers(in []byte) (*ws.ListenChan, error) {
	const op = "KrakenHandler.ToTickers"

	if string(in) == `{"event":"heartbeat"}` || strings.Contains(string(in), `"status":"subscribed"`) {
		return nil, nil
	}

	tradeInfo, err := processTrade(string(in))
	if err != nil {
		return nil, ez.New(op, ez.EINVALID, "invalid trade info", nil)
	}
	return &ws.ListenChan{
		Pair: tradeInfo.Pair,
		Tickers: []market.Ticker{
			{
				Time:   time.Now().Unix(),
				Volume: tradeInfo.LastVolume,
				Last:   tradeInfo.LastPrice,
			},
		},
	}, nil
}

func (h *KrakenHandler) ToOrderBook(in []byte) (*ws.ListenChan, error) {
	const op = "KrakenHandler.ToOrderBook"

	if string(in) == `{"event":"heartbeat"}` || strings.Contains(string(in), `"status":"subscribed"`) || strings.Contains(string(in), `"as"`) {
		return nil, nil
	}

	krakenTicker, m, err := processTicker(string(in))
	if err != nil {
		return nil, ez.New(op, ez.EINVALID, "invalid ticker", err)
	}

	orderBook := market.OrderBook{
		Time: time.Now().Unix(),
		Bids: []market.OrderBookRow{
			{
				Price:       krakenTicker.BidPrice[0],
				Volume:      krakenTicker.BidPrice[2],
				AccumVolume: krakenTicker.BidPrice[2],
			},
		},
		Asks: []market.OrderBookRow{
			{
				Price:       krakenTicker.AskPrice[0],
				Volume:      krakenTicker.AskPrice[2],
				AccumVolume: krakenTicker.AskPrice[2],
			},
		},
	}

	return &ws.ListenChan{
		Pair:      m,
		OrderBook: orderBook,
	}, nil
}

func (h KrakenHandler) GetSettings(pairs []market.Pair, channels []genericws.ChannelOpts) (genericws.Settings, error) {
	return genericws.Settings{
		Endpoint: "wss://ws.kraken.com",
	}, nil
}

func (h KrakenHandler) GetSubscriptionsRequests(pairs []market.Pair, channels []genericws.ChannelOpts) ([]genericws.SubscriptionRequest, error) {
	const op = "KrakenHandler.GetSubscriptionRequests"

	var requests []genericws.SubscriptionRequest
	channelsMap := make(map[genericws.ChannelType]bool, len(channels))
	for _, channel := range channels {
		channelsMap[channel.Type] = true
	}

	pairsArray := make([]string, len(pairs))

	for i, pair := range pairs {
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

	return ez.New(op, ez.EINTERNAL, "Error subscribing to Kraken", nil)
}
