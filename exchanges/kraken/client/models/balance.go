package models

// Balance is an Arrray of trade balance info
type Balance struct {
	EquivalentBalance      string `json:"eb"`
	TradeBalance           string `json:"tb"`
	MarginAmount           string `json:"m"`
	UnrealizedBalance      string `json:"n"`
	CostBasisOpenPositions string `json:"c"`
	CurrentValuation       string `json:"v"`
	Equity                 string `json:"e"`
	FreeMargin             string `json:"mf"`
	MarginLevel            string `json:"ml"`
}
