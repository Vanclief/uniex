package ws

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
