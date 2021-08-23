package client

type MetatraderCandle struct {
  Symbol string     `json:"symbol"`
  Timeframe string  `json:"timeframe"`
  Time string       `json:"time"`
  BrokerTime string `json:"brokerTime"`
  Open float64        `json:"open"`
  High float64        `json:"high"`
  Low float64         `json:"low"`
  Close float64       `json:"close"`
  TickVolume float64  `json:"tickVolume"`
  Spread float64      `json:"spread"`
  Volume float64      `json:"volume"`
}

type SwaggerResponse struct {
  Payload interface{}
}