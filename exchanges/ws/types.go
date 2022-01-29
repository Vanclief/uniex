package ws

type ChannelType string

var (
	OrderBookChannel ChannelType = "orderbook"
	TickerChannel    ChannelType = "ticker"
)
