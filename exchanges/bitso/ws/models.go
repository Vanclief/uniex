package ws

type SubscriptionMessage struct {
	Action string `json:"action"`
	Book   string `json:"book"`
	Type   string `json:"type"`
}

type Type struct {
	Type string `json:"type"`
}

type Trade struct {
	Type    string `json:"type"`
	Book    string `json:"book"`
	Payload []struct {
		I int    `json:"i"`
		A string `json:"a"`
		R string `json:"r"`
		V string `json:"v"`
	} `json:"payload"`
}

type DiffOrder struct {
	Type     string `json:"type"`
	Book     string `json:"book"`
	Sequence int    `json:"sequence"`
	Payload  []struct {
		D int64  `json:"d"`
		R string `json:"r"`
		T int    `json:"t"`
		A string `json:"a"`
		V string `json:"v"`
		O string `json:"o"`
	} `json:"payload"`
}

type BidAsk struct {
	Rate      float64 `json:"r"`
	Amount    float64 `json:"a"`
	Value     float64 `json:"v"`
	SellOrBuy int     `json:"t"`
	UnixTime  int64   `json:"d"`
}

type Order struct {
	Type    string `json:"type"`
	Book    string `json:"book"`
	Payload struct {
		Bids []BidAsk `json:"bids"`
		Asks []BidAsk `json:"asks"`
	} `json:"payload"`
}
