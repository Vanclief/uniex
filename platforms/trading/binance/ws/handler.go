package ws

import (
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/ws"
	"github.com/vanclief/uniex/interfaces/ws/genericws"
)

type binanceHandler struct{}

func NewHandler() binanceHandler {
	return binanceHandler{}
}

func (h binanceHandler) ToTickers(in []byte) (*ws.TickerChan, error) {
	return nil, nil
}

func (h binanceHandler) ToOrderBook(in []byte) (*ws.OrderBookChan, error) {
	return nil, nil
}

func (h binanceHandler) GetBaseEndpoint(pair []market.Pair) string {
	return "wss://ws.bitso.com"
}

func (h binanceHandler) GetSubscriptionsRequests(pairs []market.Pair, channelType genericws.ChannelType) ([]genericws.SubscriptionRequest, error) {
	const op = "handler.GetSubscriptionRequests"

	requests := make([]genericws.SubscriptionRequest, 0, len(pairs))

	// for _, pair := range pairs {
	// 	channel := ordersChannel
	// 	if channelType == genericws.ChannelTypeTicker {
	// 		channel = tickerChannel
	// 	}
	// 	subscriptionMessage := SubscriptionMessage{
	// 		Action: "subscribe",
	// 		Book:   strings.ToLower(pair.Symbol("_")),
	// 		Type:   channel,
	// 	}

	// 	request, err := json.Marshal(subscriptionMessage)
	// 	if err != nil {
	// 		return nil, ez.Wrap(op, err)
	// 	}

	// 	requests = append(requests, request)
	// }

	return requests, nil
}

func (h binanceHandler) VerifySubscriptionResponse(in []byte) error {
	const op = "binanceHandler.VerifySubscriptionResponse"

	// response := &SubscriptionResponse{}

	// err := json.Unmarshal(in, &response)
	// if err != nil {
	// 	return ez.Wrap(op, err)
	// }

	// if response.Response != "ok" {
	// 	return ez.New(op, ez.EINTERNAL, "Error on verify subscription response", nil)
	// }

	return nil
}
