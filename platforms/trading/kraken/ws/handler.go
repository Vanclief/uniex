package ws

import (
	"encoding/json"
	"fmt"
	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/ws"
	"github.com/vanclief/uniex/interfaces/ws/genericws"
	"strings"
	"time"
)

type KrakenHandler struct{}

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

func (h KrakenHandler) ToTickers(in []byte) (*ws.TickerChan, error) {
	if string(in) == `{"event":"heartbeat"}` || strings.Contains(string(in), `"status":"subscribed"`) {
		return nil, nil
	}

	const op = "KrakenHandler.ToTickers"

	arrays := getTickerArrays(string(in))
	if len(arrays) != 9 {
		return nil, ez.New(op, ez.EINVALID, "invalid ticker arrays length", nil)
	}
	krakenTickerContent := KrakenTickerContent{
		AskPrice:       arrays[0],
		BidPrice:       arrays[1],
		ClosePrice:     arrays[2],
		Volume:         arrays[3],
		VWAP:           arrays[4],
		NumberOfTrades: arrays[5],
		LowPrice:       arrays[6],
		HighPrice:      arrays[7],
		OpenPrice:      arrays[8],
	}

	ticks := []market.Ticker{{
		Time:   time.Now().Unix(),
		Ask:    krakenTickerContent.AskPrice[0],
		Bid:    krakenTickerContent.BidPrice[0],
		Last:   0,
		Volume: krakenTickerContent.Volume[0],
		VWAP:   krakenTickerContent.VWAP[0],
	},
	}

	pair := strings.Split(string(in), `"ticker",`)[1]
	marketPair := pairStringToMarketPair(strings.ReplaceAll(pair[:len(pair)-1], `"`, ""))

	return &ws.TickerChan{
		Pair:  marketPair,
		Ticks: ticks,
	}, nil
}

func (h KrakenHandler) ToOrderBook(in []byte) (*ws.OrderBookChan, error) {
	fmt.Println("ToOrderBook", string(in))
	return nil, nil
}

func (h KrakenHandler) GetBaseEndpoint(pair []market.Pair, channelType genericws.ChannelType) string {
	return "wss://ws.kraken.com"
}

func (h KrakenHandler) GetSubscriptionsRequests(pairs []market.Pair, _ genericws.ChannelType) ([]genericws.SubscriptionRequest, error) {
	const op = "KrakenHandler.GetSubscriptionRequests"
	requests := make([]genericws.SubscriptionRequest, 0, len(pairs))
	pairsArray := make([]string, len(pairs))
	for i, pair := range pairs {
		pairsArray[i] = pair.String()[:len(pair.String())-1]
	}
	subscriptionMessage := KrakenSubscriptionRequest{
		Event: "subscribe",
		Pair:  pairsArray,
		Subscription: KrakenSubscription{
			Name: "ticker",
		},
	}
	byteMessage, err := json.Marshal(subscriptionMessage)
	if err != nil {
		return nil, ez.New(op, ez.EINTERNAL, "Error parsing Subscription Message Request", err)
	}
	requests = append(requests, byteMessage)
	return requests, nil
}

func (h KrakenHandler) VerifySubscriptionResponse(in []byte) error {
	const op = "KrakenHandler.VerifySubscriptionResponse"
	if strings.Contains(string(in), `"status":"subscribed"`) {
		return nil
	}
	return ez.New(op, ez.EINTERNAL, "Error subscribing to Kraken", nil)
}
