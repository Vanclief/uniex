package ws

import (
	"encoding/json"
	"fmt"

	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/ws"
	"github.com/vanclief/uniex/interfaces/ws/genericws"
)

type parser struct{}

func NewParser() parser {
	return parser{}
}

func (p parser) ToTickers(in []byte) (*ws.TickerChan, error) {

	fmt.Println("In??", in)
	return nil, nil
}

func (p parser) ToOrderBook(in []byte) (*ws.OrderBookChan, error) {
	return nil, nil
}

func (p parser) GetSubscriptionRequest(pair market.Pair, channelType genericws.ChannelType) ([]byte, error) {
	subscriptionMessage := SubscriptionMessage{
		Method: "SUBSCRIBE",
		// Book:   strings.ToLower(pair.Symbol("_")),
		ID: 1,
	}

	return json.Marshal(subscriptionMessage)
}
