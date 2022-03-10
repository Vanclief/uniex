package ws

type MEXCSubscriptionRequest struct {
	Method string     `json:"method"`
	Param  MEXCSymbol `json:"param"`
}

type MEXCSymbol struct {
	Symbol string `json:"symbol"`
}

type MEXCTickerPayload struct {
	Channel   string         `json:"channel"`
	Data      MEXCTickerData `json:"data"`
	Symbol    string         `json:"symbol"`
	Timestamp int64          `json:"ts"`
}

type MEXCTickerData struct {
	Amount24      float64 `json:"amount24"`
	Ask1          float64 `json:"ask1"`
	Bid1          float64 `json:"bid1"`
	ContractId    int     `json:"contractId"`
	FairPrice     float64 `json:"fairPrice"`
	FundingRate   float64 `json:"fundingRate"`
	High24Price   float64 `json:"high24Price"`
	HoldVol       float64 `json:"holdVol"`
	IndexPrice    float64 `json:"indexPrice"`
	LastPrice     float64 `json:"lastPrice"`
	Lower24Price  float64 `json:"lower24Price"`
	MaxBidPrice   float64 `json:"maxBidPrice"`
	MinAskPrice   float64 `json:"minAskPrice"`
	RiseFallRate  float64 `json:"riseFallRate"`
	RiseFallValue float64 `json:"riseFallValue"`
	Symbol        string  `json:"symbol"`
	Timestamp     int64   `json:"timestamp"`
	Volume24      float64 `json:"volume24"`
}

type MEXCOrderBookPayload struct {
	Channel string            `json:"channel"`
	Data    MEXCOrderBookData `json:"data"`
	Symbol  string            `json:"symbol"`
	Ts      int64             `json:"ts"`
}

type MEXCOrderBookData struct {
	Asks    [][]float64 `json:"asks"`
	Bids    [][]float64 `json:"bids"`
	Version int64       `json:"version"`
}
