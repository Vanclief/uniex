package client

type MetaAPIResponse struct {
	Payload interface{}
}

type MetaAPIError struct {
	ID      int    `json:"id"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

type MetaTraderCandle struct {
	Symbol     string  `json:"symbol"`
	Timeframe  string  `json:"timeframe"`
	Time       string  `json:"time"`
	BrokerTime string  `json:"brokerTime"`
	Open       float64 `json:"open"`
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	Close      float64 `json:"close"`
	TickVolume float64 `json:"tickVolume"`
	Spread     float64 `json:"spread"`
	Volume     float64 `json:"volume"`
}

type MetatraderAccountInformation struct {
	Platform     string  `json:"platform"`
	Broker       string  `json:"broker"`
	Currency     string  `json:"currency"`
	Server       string  `json:"server"`
	Balance      float64 `json:"balance"`
	Equity       float64 `json:"equity"`
	Margin       float64 `json:"margin"`
	FreeMargin   float64 `json:"freeMargin"`
	Leverage     float64 `json:"leverage"`
	MarginLevel  int64   `json:"marginLevel"`
	TradeAllowed bool    `json:"tradeAllowed"`
	InvestorMode bool    `json:"investorMode"`
	MarginMode   string  `json:"marginMode"`
	Name         string  `json:"name"`
	Login        int64   `json:"login"`
	Credit       float64 `json:"credit"`
}

type MetatraderPosition struct {
	ID                          string  `json:"id"`
	Type                        string  `json:"type"`
	Symbol                      string  `json:"symbol"`
	Magic                       int     `json:"magic"`
	Time                        string  `json:"time"`
	BrokerTime                  string  `json:"brokerTime"`
	UpdateTime                  string  `json:"updateTime"`
	OpenPrice                   float64 `json:"openPrice"`
	CurrentPrice                float64 `json:"currentPrice"`
	CurrentTickValue            float64 `json:"currentTickValue"`
	StopLoss                    float64 `json:"stopLoss"`
	TakeProfit                  float64 `json:"takeProfit"`
	Volume                      float64 `json:"volume"`
	Swap                        float64 `json:"swap"`
	Profit                      float64 `json:"profit"`
	Comment                     string  `json:"comment"`
	ClientID                    string  `json:"clientId"`
	Commission                  float64 `json:"commission"`
	Reason                      string  `json:"reason"`
	UnrealizedProfit            float64 `json:"unrealizedProfit"`
	RealizedProfit              float64 `json:"realizedProfit"`
	AccountCurrencyExchangeRate float64 `json:"accountCurrencyExchangeRate"`
	BrokerComment               string  `json:"brokerComment"`
}

type MetatraderOrder struct {
	ID                          string  `json:"id"`
	Type                        string  `json:"type"`
	State                       string  `json:"state"`
	Magic                       int     `json:"magic"`
	Time                        string  `json:"time"`
	BrokerTime                  string  `json:"brokerTime"`
	DoneTime                    string  `json:"doneTime"`
	BrokerDoneTime              string  `json:"brokerDoneTime"`
	Symbol                      string  `json:"symbol"`
	OpenPrice                   float64 `json:"openPrice"`
	StopLimitPrice              float64 `json:"stopLimitPrice"`
	CurrentPrice                float64 `json:"currentPrice"`
	StopLoss                    float64 `json:"stopLoss"`
	TakeProfit                  float64 `json:"takeProfit"`
	Volume                      float64 `json:"volume"`
	CurrentVolume               float64 `json:"currentVolume"`
	PositionID                  string  `json:"positionId"`
	Comment                     string  `json:"comment"`
	BrokerComment               string  `json:"brokerComment"`
	ClientID                    string  `json:"clientId"`
	Platform                    string  `json:"platform"`
	Reason                      string  `json:"reason"`
	FillingMode                 string  `json:"fillingMode"`
	ExpirationType              string  `json:"expirationType"`
	ExpirationTime              string  `json:"expirationTime"`
	AccountCurrencyExchangeRate float64 `json:"accountCurrencyExchangeRate"`
	CloseByPosition             float64 `json:"closeByPosition"`
}

type MetatraderTrade struct {
	ActionType        string   `json:"actionType"`
	Symbol            string   `json:"symbol"`
	Volume            float64  `json:"volume"`
	OpenPrice         float64  `json:"openPrice"`
	StopLoss          float64  `json:"stopLoss"`
	TakeProfit        float64  `json:"takeProfit"`
	StopLossUnits     string   `json:"stopLossUnits"`
	TakeProfitUnits   string   `json:"takeProfitUnits"`
	OrderID           string   `json:"orderId"`
	PositionID        string   `json:"positionId"`
	Comment           string   `json:"comment"`
	ClientID          string   `json:"clientId"`
	Magic             int      `json:"magic"`
	Slippage          float64  `json:"slippage"`
	FillingModes      []string `json:"fillingModes"`
	CloseByPositionId string   `json:"closeByPositionId"`
	StopLimitPrice    float64  `json:"stopLimitPrice"`
}

type MetatraderTradeResponse struct {
	NumericCode int    `json:"numericCode"`
	StringCode  string `json:"stringCode"`
	Message     string `json:"message"`
	OrderID     string `json:"orderId"`
	PositionID  string `json:"positionId"`
}
