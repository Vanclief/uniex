package client

type BitsoResponse struct {
	Success bool        `json:"success"`
	Payload interface{} `json:"payload"`
}

type BitsoError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type Ticker struct {
	Book      string `json:"book"`
	Volume    string `json:"volume"`
	High      string `json:"high"`
	Last      string `json:"last"`
	Low       string `json:"low"`
	Wwap      string `json:"wwap"`
	Ask       string `json:"ask"`
	Bid       string `json:"bid"`
	CreatedAt string `json:"created_at"`
	Change24  string `json:"change_24"`
}

type OrderBook struct {
	Asks      []OrderBookRow `json:"asks"`
	Bids      []OrderBookRow `json:"bids"`
	UpdatedAt string         `json:"updated_at"`
	Sequence  string         `json:"sequence"`
}

type OrderBookRow struct {
	Book   string `json:"book"`
	Price  string `json:"price"`
	Amount string `json:"amount"`
}
