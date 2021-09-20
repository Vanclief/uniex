package exchanges

import "github.com/vanclief/finmod/market"

type GetOrderBookOptions struct {
	Limit int // Maximum number of asks/bids
}

type Status string

const (
	OpenStatus   Status = "open"
	ClosedStatus Status = "closed"
)

type GetOrdersRequest struct {
	IDs       []string `json:"ids"`
	Status    Status   `json:"status"`
	StartTime string   `json:"start_time"`
	EndTime   string   `json:"end_time"`
	Offset    int      `json:"offset"`
}

type UpdateOrderRequest struct {
	Price      float64 `json:"price"`
	StopLoss   float64 `json:"stop_loss"`
	TakeProfit float64 `json:"take_profit"`
	Volume     float64 `json:"volume"`
}

type GetTradesRequest struct {
	IDs       []string `json:"ids"`
	StartTime string   `json:"start_time"`
	EndTime   string   `json:"end_time"`
	Offset    int      `json:"offset"`
}

type GetPositionsRequest struct {
	IDs       []string    `json:"ids"`
	Status    Status      `json:"status"`
	Pair      market.Pair `json:"pair"`
	StartTime string      `json:"start_time"`
	EndTime   string      `json:"end_time"`
	Offset    int         `json:"offset"`
}

type UpdatePositionRequest struct {
	TakeProfit float64 `json:"take_profit"`
	StopLoss   float64 `json:"stop_loss"`
}
