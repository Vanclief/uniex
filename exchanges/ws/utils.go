package ws

import (
	"fmt"
	"github.com/vanclief/finmod/market"
	"strings"
)

func ToMarketPair(raw, sep string) (market.Pair, error) {
	baseQuote := strings.Split(raw, sep)
	if len(baseQuote) != 2 {
		return market.Pair{}, fmt.Errorf("fail to create market pair form '%s'", raw)
	}
	return market.Pair{
		Base:  &market.Asset{
			Symbol: strings.ToUpper(baseQuote[0]),
		},
		Quote: &market.Asset{
			Symbol: strings.ToUpper(baseQuote[1]),
		},
	}, nil
}
