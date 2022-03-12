package genericws

import (
	"fmt"
	"strings"

	"github.com/vanclief/finmod/market"
)

func ToMarketPair(str, sep string) (market.Pair, error) {
	baseQuote := strings.Split(str, sep)

	if len(baseQuote) != 2 {
		return market.Pair{}, fmt.Errorf("Failed to create market pair from '%s'", str)
	}

	return market.Pair{
		Base: market.Asset{
			Symbol: strings.ToUpper(baseQuote[0]),
		},
		Quote: market.Asset{
			Symbol: strings.ToUpper(baseQuote[1]),
		},
	}, nil
}
