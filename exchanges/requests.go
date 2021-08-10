package exchanges

// GetOrderBookOptions
type GetOrderBookOptions struct {
	Limit int // Maximum number of asks/bids
}

type GetOrdersRequest struct {
	IDs       []string `json:"ids"`
	Status    string   `json:"status"` // Open, Close
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
	Status    string   `json:"status"` // Open, Close
	StartTime string   `json:"start_time"`
	EndTime   string   `json:"end_time"`
	Offset    int      `json:"offset"`
}
