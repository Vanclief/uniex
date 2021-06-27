package polygon

type RangeResult struct {
	V  float64 `json:"v"`
	Vw float64 `json:"vw"`
	O  float64 `json:"o"`
	C  float64 `json:"c"`
	H  float64 `json:"h"`
	L  float64 `json:"l"`
	T  int64   `json:"t"`
	N  int     `json:"n"`
}

type RangeResponse struct {
	Ticker       string         `json:"ticker"`
	QueryCount   int            `json:"queryCount"`
	ResultsCount int            `json:"resultsCount"`
	Adjusted     bool           `json:"adjusted"`
	Results      []*RangeResult `json:"results"`
	Status       string         `json:"status"`
	RequestId    string         `json:"request_id"`
	Count        int            `json:"count"`
}

type NBBOResult struct {
	ExchangeSymbol  string  `json:"T"`
	BidPrice        float64 `json:"p"`
	BidSize         int     `json:"s"`
	BidExchangeId   int     `json:"x"`
	AskPrice        float64 `json:"P"`
	AskSize         int     `json:"S"`
	AskExchangeId   int     `json:"X"`
	SequenceNumber  int     `json:"q"`
	TimeSIP         int64   `json:"t"`
	TimeParticipant int64   `json:"y"`
	Tape            int     `json:"z"`
}

type NBBOResponse struct {
	RequestId string      `json:"request_id"`
	Status    string      `json:"status"`
	Results   *NBBOResult `json:"results"`
}
