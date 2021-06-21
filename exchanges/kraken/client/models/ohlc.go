package models

import "encoding/json"

// OHLCResponse - response of OHLC request
type OHLCResponse struct {
	Candles map[string][]Candle `json:"-"`
	Last    int64               `json:"last"`
}

// Candle - OHLC item
type Candle struct {
	Time      int64
	Open      float64 `json:",string"`
	High      float64 `json:",string"`
	Low       float64 `json:",string"`
	Close     float64 `json:",string"`
	VolumeWAP float64 `json:",string"`
	Volume    float64 `json:",string"`
	Count     int64
}

// UnmarshalJSON -
func (item *OHLCResponse) UnmarshalJSON(buf []byte) error {
	res := make(map[string]interface{})
	if err := json.Unmarshal(buf, &res); err != nil {
		return err
	}

	last, err := getTimestamp(res["last"])
	if err != nil {
		return err
	}
	item.Last = last
	delete(res, "last")

	item.Candles = make(map[string][]Candle)
	for k, v := range res {
		items := v.([]interface{})
		item.Candles[k] = make([]Candle, len(items))
		for idx, c := range items {
			candle := c.([]interface{})

			ts, err := getTimestamp(candle[0])
			if err != nil {
				continue
			}
			open, err := getFloat64FromStr(candle[1])
			if err != nil {
				continue
			}
			high, err := getFloat64FromStr(candle[2])
			if err != nil {
				continue
			}
			low, err := getFloat64FromStr(candle[3])
			if err != nil {
				continue
			}
			close, err := getFloat64FromStr(candle[4])
			if err != nil {
				continue
			}
			vwap, err := getFloat64FromStr(candle[5])
			if err != nil {
				continue
			}
			vol, err := getFloat64FromStr(candle[6])
			if err != nil {
				continue
			}
			item.Candles[k][idx] = Candle{
				Time:      ts,
				Open:      open,
				High:      high,
				Low:       low,
				Close:     close,
				VolumeWAP: vwap,
				Volume:    vol,
				Count:     int64(candle[7].(float64)),
			}
		}
	}
	return nil
}
