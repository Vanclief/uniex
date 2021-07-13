package client

import "time"

type BitsoResponse struct {
	Success bool        `json:"success"`
	Payload interface{} `json:"payload"`
	Error   BitsoError  `json:"error"`
}

type BitsoError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type Ticker struct {
	Book      string    `json:"book"`
	Volume    string    `json:"volume"`
	High      string    `json:"high"`
	Last      string    `json:"last"`
	Low       string    `json:"low"`
	Wwap      string    `json:"wwap"`
	Ask       string    `json:"ask"`
	Bid       string    `json:"bid"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderBook struct {
	Book   string `json:"book"`
	Price  string `json:"price"`
	Amount string `json:"amount"`
}
