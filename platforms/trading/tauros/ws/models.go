package ws

type SubscriptionResponse struct{}

type SubscriptionMessage struct {
	Action  string `json:"action"`
	Market  string `json:"market"`
	Channel string `json:"channel"`
}

type Trade struct {
	Amount    float64 `json:"amount,string"`
	Value     float64 `json:"value,string"`
	Price     float64 `json:"price,string"`
	Side      string  `json:"side"`
	Timestamp float64 `json:"timestamp"`
}

type Tick struct {
	Action    string  `json:"action"`
	Channel   string  `json:"channel"`
	Market    string  `json:"market"`
	Volume    string  `json:"volume"`
	High      string  `json:"high"`
	Low       string  `json:"low"`
	Last      string  `json:"last"`
	Variation float64 `json:"variation"`
	Trades    []Trade `json:"trades"`
}

type BidAsk struct {
	Amount float64 `json:"a,string"`
	Price  float64 `json:"p,string"`
	Value  float64 `json:"v,string"`
	UnixMs int64   `json:"t"`
}

type Data struct {
	Asks []BidAsk `json:"asks"`
	Bids []BidAsk `json:"bids"`
}

type Order struct {
	Action  string `json:"action"`
	Channel string `json:"channel"`
	Type    string `json:"type"`
	Data    Data   `json:"data"`
}
