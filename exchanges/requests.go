package exchanges

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

type GetTradesRequest struct {
	IDs       []string `json:"ids"`
	StartTime string   `json:"start_time"`
	EndTime   string   `json:"end_time"`
	Offset    int      `json:"offset"`
}

type GetPositionsRequest struct {
	IDs       []string `json:"ids"`
	Status    Status   `json:"status"`
	StartTime string   `json:"start_time"`
	EndTime   string   `json:"end_time"`
	Offset    int      `json:"offset"`
}

type UpdatePositionRequest struct {
	TakeProfit float64 `json:"take_profit"`
	StopLoss   float64 `json:"stop_loss"`
}

type ClosePositionRequest struct {
	ID string `json:"id"`
}
