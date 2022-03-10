package ws

type SubscriptionMessageRequest struct {
	ID             int    `json:"id"`
	Type           string `json:"type"`
	Topic          string `json:"topic"`
	PrivateChannel bool   `json:"privateChannel"`
	Response       bool   `json:"response"`
}

type SubscriptionMessageResponse struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type TradeTypeData struct {
	Sequence    string `json:"sequence"`
	Price       string `json:"price"`
	Size        string `json:"size"`
	BestAsk     string `json:"bestAsk"`
	BestAskSize string `json:"bestAskSize"`
	BestBid     string `json:"bestBid"`
	BestBidSize string `json:"bestBidSize"`
	Time        int64  `json:"time"`
}

type TradeType struct {
	Type    string        `json:"type"`
	Topic   string        `json:"topic"`
	Subject string        `json:"subject"`
	Data    TradeTypeData `json:"data"`
}

type Order struct {
	Type    string `json:"type"`
	Topic   string `json:"topic"`
	Subject string `json:"subject"`
	Data    struct {
		Asks      [][]string `json:"asks"`
		Bids      [][]string `json:"bids"`
		Timestamp int64      `json:"timestamp"`
	} `json:"data"`
}

type Type struct {
	Subject string `json:"subject"`
	Type    string `json:"type"`
}
