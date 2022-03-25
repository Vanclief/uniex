package ws

type SubscriptionRequest struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	ID     int      `json:"id"`
}

type SubscriptionResponse struct {
	Method interface{} `json:"method"`
	ID     int         `json:"id"`
}

type StreamTickerEvent struct {
	Stream string            `json:"stream"`
	Data   BinanceTickerData `json:"data"`
}

type StreamPartialOrderBookEvent struct {
	Stream string                      `json:"stream"`
	Data   BinancePartialOrderBookData `json:"data"`
}

type StreamUpdateOrderBookEvent struct {
	Stream string                 `json:"stream"`
	Data   BinanceUpdateOrderBook `json:"data"`
}

type BinanceTickerData struct {
	EventType            string `json:"e"`
	EventTime            int    `json:"E"`
	Symbol               string `json:"s"`
	PriceChange          string `json:"p"`
	PriceChangePercent   string `json:"P"`
	WeightedAveragePrice string `json:"w"`
	FirstTrade           string `json:"x"`
	LastPrice            string `json:"c"`
	LastQuantity         string `json:"Q"`
	BestBidPrice         string `json:"b"`
	BestBidQuantity      string `json:"B"`
	BestAskPrice         string `json:"a"`
	BestAskQuantity      string `json:"A"`
	OpenPrice            string `json:"o"`
	HighPrice            string `json:"h"`
	LowPrice             string `json:"l"`
	BaseAssetVolume      string `json:"v"`
	QuoteAssetVolume     string `json:"q"`
	OpenTime             int    `json:"O"`
	CloseTime            int    `json:"C"`
	FirstTradeID         int    `json:"F"`
	LastTradeID          int    `json:"L"`
	TotalNumberOfTrades  int    `json:"n"`
}

type BinanceUpdateOrderBook struct {
	OrderBookUpdateID int    `json:"u"`
	Symbol            string `json:"s"`
	BestBidPrice      string `json:"b"`
	BestBidQuantity   string `json:"B"`
	BestAskPrice      string `json:"a"`
	BestAskQuantity   string `json:"A"`
}

type BinancePartialOrderBookData struct {
	EventType     string     `json:"e"`
	EventTime     int64      `json:"E"`
	Symbol        string     `json:"s"`
	FirstUpdateID int64      `json:"U"`
	FinalUpdateUD int64      `json:"u"`
	Bids          [][]string `json:"b"`
	Asks          [][]string `json:"a"`
}
