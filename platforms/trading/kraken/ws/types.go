package ws

import (
	"encoding/json"
	"fmt"
)

type KrakenSubscriptionRequest struct {
	Event        string             `json:"event"`
	Pair         []string           `json:"pair"`
	Subscription KrakenSubscription `json:"subscription"`
}

type KrakenSubscription struct {
	Name string `json:"name"`
}

type KrakenSubscriptionResponse struct {
	ChannelID    int                `json:"channelID"`
	ErrorMessage string             `json:"errorMessage"`
	ChannelName  string             `json:"channelName"`
	Event        string             `json:"event"`
	Pair         string             `json:"pair"`
	Status       string             `json:"status"`
	Subscription KrakenSubscription `json:"subscription"`
}

type KrakenTickerPayload struct {
	ChannelID   int                 `json:"-"`
	Ticker      KrakenTickerContent `json:"-"`
	ChannelName string              `json:"-"`
	Pair        string              `json:"-"`
}

type KrakenTickerContent struct {
	AskPrice       []float64 `json:"a"`
	BidPrice       []float64 `json:"b"`
	ClosePrice     []float64 `json:"c"`
	Volume         []float64 `json:"v"`
	VWAP           []float64 `json:"p"`
	NumberOfTrades []float64 `json:"t"`
	LowPrice       []float64 `json:"l"`
	HighPrice      []float64 `json:"h"`
	OpenPrice      []float64 `json:"o"`
}

type KrakenOrderBookContent struct {
	Asks     [][]string `json:"a"`
	Bids     [][]string `json:"b"`
	Checksum string     `json:"c"`
}

// Seq -
type Seq struct {
	Value int64 `json:"sequence"`
}

// KrakenOrderBookPayload - data structure of default Kraken WS update
type KrakenOrderBookPayload struct {
	ChannelID   int64
	Data        json.RawMessage
	ChannelName string
	Pair        string
	Sequence    Seq
}

// UnmarshalJSON - unmarshal update
func (msg *KrakenOrderBookPayload) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if len(raw) < 3 {
		return fmt.Errorf("invalid data length: %#v", raw)
	}

	if len(raw) == 5 {
		// order book can have 2 data objects
		// one for the new asks and one for the new bids
		// see https://docs.kraken.com/websockets/

		// the array is [channelid, ask, bid, channel, pair]
		ask := raw[1]
		bid := raw[2]

		// ask and bid can be merged into a single object as the keys are distinct
		if ask[len(ask)-1] != '}' || bid[0] != '{' {
			// not a bid/ask pair
			return fmt.Errorf("invalid data length/payload: %v", raw)
		}

		// merge ask + bid
		merged := make([]byte, 0, len(ask)+len(bid)-1)
		merged = append(merged, ask[0:len(ask)-1]...)
		merged = append(merged, ',')
		merged = append(merged, bid[1:]...)

		// reencode
		data, _ = json.Marshal([]json.RawMessage{
			raw[0], merged, raw[3], raw[4],
		})
	}

	body := make([]interface{}, 0)
	if len(raw) == 3 {
		body = append(body, &msg.Data, &msg.ChannelName, &msg.Sequence)
	} else {
		body = append(body, &msg.ChannelID, &msg.Data, &msg.ChannelName, &msg.Pair)
	}

	return json.Unmarshal(data, &body)
}
