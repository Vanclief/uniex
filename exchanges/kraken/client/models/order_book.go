package models

import (
	"encoding/json"
	"fmt"
)

// OrderBookItem - one price level in orderbook
type OrderBookItem struct {
	Price     float64
	Volume    float64
	Timestamp int64
}

// UnmarshalJSON -
func (item *OrderBookItem) UnmarshalJSON(buf []byte) error {
	var tmp []interface{}
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), 3; g != e {
		return fmt.Errorf("wrong number of fields in OrderBookItem: %d != %d", g, e)
	}

	price, err := getFloat64FromStr(tmp[0])
	if err != nil {
		return err
	}
	item.Price = price

	vol, err := getFloat64FromStr(tmp[1])
	if err != nil {
		return err
	}
	item.Volume = vol

	ts, err := getTimestamp(tmp[2])
	if err != nil {
		return err
	}
	item.Timestamp = ts

	return nil
}

// OrderBook - struct of order book levels
type OrderBook struct {
	Asks []OrderBookItem `json:"asks"`
	Bids []OrderBookItem `json:"bids"`
}
