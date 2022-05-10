package market

import (
	"fmt"
	"github.com/vanclief/ez"
	"math"
	"sort"
	"time"
)

// OrderBookRow - A single order from the order book
type OrderBookRow struct {
	Price       float64 `json:"price"`
	Volume      float64 `json:"volume"`
	AccumVolume float64 `json:"accum_volume"`
}

func (obr *OrderBookRow) String() string {
	return fmt.Sprintf("Price: %.3f, Volume: %.3f, AccumVolume: %.3f\n", obr.Price, obr.Volume, obr.AccumVolume)
}

type OrderBookUpdate struct {
	Price  float64 `json:"price"`
	Volume float64 `json:"volume"`
	Side   string  `json:"side"`
}

// OrderBook - A record of active buy and sell orders in a single market
type OrderBook struct {
	Time     int64          `json:"time"`
	Asks     []OrderBookRow `json:"asks"` // ordered from highest to lowest
	Bids     []OrderBookRow `json:"bids"` // ordered from lowest to highest
	MaxDepth int            `json:"max_depth"`
}

func NewOrderBook(asks, bids []OrderBookRow, maxDepth int) OrderBook {

	ob := OrderBook{
		Time:     time.Now().Unix(),
		Asks:     asks,
		Bids:     bids,
		MaxDepth: maxDepth,
	}

	ob.sort()

	return ob
}

func (ob *OrderBook) sort() {
	sort.SliceStable(ob.Asks, func(i, j int) bool {
		return ob.Asks[i].Price < ob.Asks[j].Price
	})

	sort.SliceStable(ob.Bids, func(i, j int) bool {
		return ob.Bids[i].Price > ob.Bids[j].Price
	})

	askDepth := math.Min(float64(len(ob.Asks)), float64(ob.MaxDepth))
	bidDepth := math.Min(float64(len(ob.Bids)), float64(ob.MaxDepth))

	ob.Asks = ob.Asks[:int(askDepth)]
	ob.Bids = ob.Bids[:int(bidDepth)]

	for i := range ob.Asks {
		if i == 0 {
			ob.Asks[i].AccumVolume = ob.Asks[i].Volume
		} else {
			ob.Asks[i].AccumVolume = ob.Asks[i].Volume + ob.Asks[i-1].AccumVolume
		}
	}

	for i := range ob.Bids {
		if i == 0 {
			ob.Bids[i].AccumVolume = ob.Bids[i].Volume
		} else {
			ob.Bids[i].AccumVolume = ob.Bids[i].Volume + ob.Bids[i-1].AccumVolume
		}
	}
}

func (ob *OrderBook) limit() {
	if len(ob.Asks) > ob.MaxDepth {
		ob.Asks = ob.Asks[:ob.MaxDepth]
	}

	if len(ob.Bids) > ob.MaxDepth {
		ob.Bids = ob.Bids[:ob.MaxDepth]
	}
}

func (ob *OrderBook) String() string {
	return fmt.Sprintf("Time: %v, Asks: %v, Bids: %v\n", ob.Time, ob.Asks, ob.Bids)
}

func (ob *OrderBook) Print() {
	fmt.Println("======= OrderBook =======")
	for i := range ob.Asks {
		i = len(ob.Asks) - 1 - i
		fmt.Println(ob.Asks[i])
	}
	fmt.Println("-------------------------")
	for _, bid := range ob.Bids {
		fmt.Println(bid)
	}
}

func (ob *OrderBook) removeInvalid(side string, price float64) int {
	var found bool
	rowsRemoved := 0
	if side == "ask" {
		if len(ob.Asks) == 0 || price < ob.Asks[0].Price {
			return 0
		}
		for i := range ob.Asks {
			if ob.Asks[i].Price < price {
				continue
			} else {
				startIndex := int(math.Min(float64(i+1), float64(len(ob.Asks)-1)))
				rowsRemoved = len(ob.Asks) - i
				ob.Asks = ob.Asks[startIndex:]
				found = true
				break
			}
		}
		if !found {
			ob.Asks = []OrderBookRow{}
		}
		return rowsRemoved
	} else if side == "bid" {
		if len(ob.Bids) == 0 || price > ob.Bids[0].Price {
			return 0
		}
		for i := range ob.Bids {
			if ob.Bids[i].Price > price {
				continue
			} else {
				startIndex := int(math.Min(float64(i+1), float64(len(ob.Bids)-1)))
				rowsRemoved = len(ob.Bids) - i
				ob.Bids = ob.Bids[startIndex:]
				found = true
				break
			}
		}
		if !found {
			ob.Bids = []OrderBookRow{}
		}
		return rowsRemoved
	}
	return -1
}

func (ob *OrderBook) ApplyUpdate(update OrderBookUpdate) error {
	const op = "OrderBook.ApplyUpdate"

	if update.Side == "ask" {
		if update.Volume != 0 {
			ob.removeInvalid("bid", update.Price)
			found := false
			for i := range ob.Asks {
				if ob.Asks[i].Price == update.Price {
					ob.Asks[i] = OrderBookRow{Price: update.Price, Volume: update.Volume}
					found = true
					ob.limit()
					return nil
				}
			}
			if !found {
				if len(ob.Asks) == 0 {
					ob.Asks = append(ob.Asks, OrderBookRow{Price: update.Price, Volume: update.Volume})
					ob.limit()
					return nil
				}
				if len(ob.Asks) == 1 {
					if ob.Asks[0].Price > update.Price {
						ob.Asks = append(ob.Asks, ob.Asks[0])
						ob.Asks[0] = OrderBookRow{Price: update.Price, Volume: update.Volume}
						ob.limit()
						return nil
					} else {
						ob.Asks = append(ob.Asks, OrderBookRow{Price: update.Price, Volume: update.Volume})
						ob.limit()
						return nil
					}
				}
				for i := range ob.Asks {
					if ob.Asks[i].Price > update.Price {
						ob.Asks = append(ob.Asks, OrderBookRow{})
						copy(ob.Asks[i+1:], ob.Asks[i:])
						ob.Asks[i] = OrderBookRow{Price: update.Price, Volume: update.Volume}
						ob.limit()
						return nil
					}
				}
				ob.Asks = append(ob.Asks, OrderBookRow{Price: update.Price, Volume: update.Volume})
				ob.limit()
				return nil
			}
		} else {
			for i := range ob.Asks {
				if ob.Asks[i].Price == update.Price {
					ob.Asks = removeElement(ob.Asks, i)
					ob.limit()
					return nil
				}
			}
		}

	} else if update.Side == "bid" {
		if update.Volume != 0 {
			ob.removeInvalid("ask", update.Price)
			found := false
			for i := range ob.Bids {
				if ob.Bids[i].Price == update.Price {
					ob.Bids[i] = OrderBookRow{Price: update.Price, Volume: update.Volume}
					found = true
					ob.limit()
					return nil
				}
			}
			if !found {
				if len(ob.Bids) == 0 {
					ob.Bids = append(ob.Bids, OrderBookRow{Price: update.Price, Volume: update.Volume})
					ob.limit()
					return nil
				}
				if len(ob.Bids) == 1 {
					if ob.Bids[0].Price < update.Price {
						ob.Bids = append(ob.Bids, ob.Bids[0])
						ob.Bids[0] = OrderBookRow{Price: update.Price, Volume: update.Volume}
						ob.limit()
						return nil
					} else {
						ob.Bids = append(ob.Bids, OrderBookRow{Price: update.Price, Volume: update.Volume})
						ob.limit()
						return nil
					}
				}
				for i := range ob.Bids {
					if ob.Bids[i].Price < update.Price {
						ob.Bids = append(ob.Bids, OrderBookRow{})
						copy(ob.Bids[i+1:], ob.Bids[i:])
						ob.Bids[i] = OrderBookRow{Price: update.Price, Volume: update.Volume}
						ob.limit()
						return nil
					}
				}
				ob.Bids = append(ob.Bids, OrderBookRow{Price: update.Price, Volume: update.Volume})
				ob.limit()
				return nil
			}
		} else {
			for i := range ob.Bids {
				if ob.Bids[i].Price == update.Price {
					ob.Bids = removeElement(ob.Bids, i)
					ob.limit()
					return nil
				}
			}
		}
	} else {
		return ez.New(op, ez.EINVALID, "update side must be ask or bid", nil)
	}

	return nil
}

func removeElement(slice []OrderBookRow, s int) []OrderBookRow {
	return append(slice[:s], slice[s+1:]...)
}

// GetDepth - returns the accumulated volume from a determined price
// a price onwards
func (ob *OrderBook) GetDepth(price float64) (float64, error) {
	op := "OrderBook.GetDepth"

	if price <= 0 {
		return 0, ez.New(op, ez.EINVALID, "Price can't be equal or below zero", nil)
	}

	firstAsk, lastAsk := ob.Asks[0], ob.Asks[len(ob.Asks)-1]
	firstBid, lastBid := ob.Bids[0], ob.Bids[len(ob.Bids)-1]
	if firstAsk.Price <= price && price <= lastAsk.Price {
		// price is in asks interval
		for k := range ob.Asks {
			if ob.Asks[k].Price > price {
				index := math.Max(0, float64(k-1))
				return ob.Asks[int(index)].Price, nil
			}
		}
		return lastAsk.Price, nil
	} else if lastBid.Price <= price && price <= firstBid.Price {
		// price is in bids interval
		for k := range ob.Bids {
			if ob.Bids[k].Price < price {
				index := math.Max(0, float64(k-1))
				return ob.Bids[int(index)].Price, nil
			}
		}
		return lastBid.Price, nil
	}
	return 0, ez.New(op, ez.ENOTFOUND, "No depth for the selected price", nil)
}

func overlapCalculation(obA, obB OrderBook) (float64, error) {
	var obAAccumVolume, obBAccumVolume float64
	var obAIndex, obBIndex int
	for i := range obA.Asks {
		if obA.Asks[i].Price > obB.Bids[0].Price {
			obAIndex = i
			break
		} else {
			obAAccumVolume += obA.Asks[i].Volume
		}
	}
	for i := range obB.Bids {
		if obB.Bids[i].Price < obA.Asks[0].Price {
			obBIndex = i
			break
		} else {
			obBAccumVolume += obB.Bids[i].Volume
		}
	}

	var obAResult, obBResult float64
	var obBVolume float64

	for i := len(obB.Bids[:obBIndex]) - 1; i >= 0; i-- {
		v := obB.Bids[i]
		if obBVolume+v.Volume > obAAccumVolume {
			obBResult += (obAAccumVolume - obBVolume) * v.Price
			break
		} else {
			obBVolume += v.Volume
			obBResult += v.Price * v.Volume
		}
	}
	for _, v := range obA.Asks[:obAIndex] {
		obAResult += v.Price * v.Volume
	}

	return math.Abs(obAResult - obBResult), nil
}

func CalculateOverlap(obA, obB OrderBook) (float64, error) {
	op := "OrderBook.CalculateOverlap"
	// First determine if overlap exists
	if !(obA.Asks[0].Price < obB.Bids[0].Price || obA.Bids[0].Price > obB.Asks[0].Price) {
		return -1, ez.New(op, ez.ENOTFOUND, "No overlap", nil)
	}
	if obA.Asks[0].Price < obB.Bids[0].Price {
		return overlapCalculation(obA, obB)
	}
	if obA.Bids[0].Price > obB.Asks[0].Price {
		return overlapCalculation(obB, obA)
	}
	return 0, nil
}

// PAIR

// Pair - Quotation of two different assets or currencies, with the value of one being quoted against the other.
type Pair struct {
	Base  Asset `json:"base"`
	Quote Asset `json:"quote"`
}

// NewPair creates a new Pair from two assets
func NewPair(base, quote Asset) Pair {
	return Pair{Base: base, Quote: quote}
}

// String - Implements Stringer interface
func (p *Pair) String() string {
	return fmt.Sprintf("%s/%s\n", p.Base.Symbol, p.Quote.Symbol)
}

// Symbol - Gets the current symbol
func (p *Pair) Symbol(separator string) string {
	if p.Quote.Symbol == "" {
		return p.Base.Symbol
	} else {
		return fmt.Sprintf("%s%s%s", p.Base.Symbol, separator, p.Quote.Symbol)
	}
}

// Asset - A resource with economic value
type Asset struct {
	Name      string `json:"name"`
	Symbol    string `json:"symbol"`
	AltSymbol string `json:"alt_symbol"`
}

func (a *Asset) String() string {
	return fmt.Sprintf("%s: %s\n", a.Name, a.Symbol)
}

// NewAsset creates a new Asset from a name and a symbol
func NewAsset(symbol, name string) (asset Asset, err error) {
	const op = "market.NewAsset"

	if symbol == "" {
		return asset, ez.New(op, ez.EINVALID, "Missing asset symbol", nil)
	} else if name == "" {
		return asset, ez.New(op, ez.EINVALID, "Missing asset name", nil)
	}

	return Asset{Symbol: symbol, Name: name}, nil
}

// Ticker - The Latest price data for an asset
type Ticker struct {
	Time   int64   `json:"time"`
	Ask    float64 `json:"ask"`
	Bid    float64 `json:"bid"`
	Last   float64 `json:"last"`
	Volume float64 `json:"volume"`
	VWAP   float64 `json:"vwap"`
}

func (t *Ticker) String() string {
	return fmt.Sprintf("Time: %d, Ask: %f, Bid: %f, Last: %f, Volume: %f, VWAP: %f\n", t.Time, t.Ask, t.Bid, t.Last, t.Volume, t.VWAP)
}
