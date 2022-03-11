package ws

type FTXSubscribeRequest struct {
	Operation string `json:"op"`
	Channel   string `json:"channel"`
	Market    string `json:"market"`
}

type FTXSubscribeResponse struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Market  string `json:"market"`
}

type FTXTickerStream struct {
	Type    string        `json:"type"`
	Channel string        `json:"channel"`
	Market  string        `json:"market"`
	Data    FTXTickerData `json:"data"`
}

// All messages are updates (update)

type FTXTickerData struct {
	Bid     float64 `json:"bid"`
	Ask     float64 `json:"ask"`
	BidSize float64 `json:"bidSize"`
	AskSize float64 `json:"askSize"`
	Last    float64 `json:"last"`
	Time    float64 `json:"time"`
}

type FTXOrderBookStream struct {
	Type    string           `json:"type"`
	Channel string           `json:"channel"`
	Market  string           `json:"market"`
	Data    FTXOrderBookData `json:"data"`
}

type FTXOrderBookData struct {
	Time     float64     `json:"time"`
	Checksum int         `json:"checksum"`
	Bids     [][]float64 `json:"bids"`
	Asks     [][]float64 `json:"asks"`
	Action   string      `json:"action"`
}
