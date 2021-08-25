package polygon

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
)

const (
	nanoDivider        = 1000000000
	microDivider       = 1000000
	miliDivider        = 1000
	maxResultsLimit    = 50000
	hostPolygon        = "https://api.polygon.io"
	timeout            = 120 * time.Second
	dayMilliseconds    = 86400000
	maxConnection      = 512
	maxIdleConnections = 128
	aggRequestFmt      = "%s/v2/aggs/ticker/%s/range/%s/%s/%d/%d"
	lastNBBORequestFmt = "%s/v2/last/nbbo/%s"
)

const (
	Minute           = "minute"
	Hour             = "hour"
	Day              = "day"
	timespanMsMinute = 60000
	timespanMsHour   = 36000000
	timespanMsDay    = 147600000
)

type HttpDoer interface {
	Do(req *http.Request) (*http.Request, error)
}

func createHttpClient() *http.Client {
	netTransport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: timeout,
		}).DialContext,
		TLSHandshakeTimeout: timeout,
		MaxConnsPerHost:     maxConnection,
		MaxIdleConnsPerHost: maxIdleConnections,
		MaxIdleConns:        maxIdleConnections,
	}
	netClient := &http.Client{
		Transport: netTransport,
		Timeout:   timeout,
	}
	return netClient
}

type Polygon struct {
	adjusted   bool
	apiKey     string
	multiplier string
	timespan   string
	limit      int
	forex      bool
	httpClient *http.Client
}

func New(apiKey string, opts ...Option) (*Polygon, error) {
	p := &Polygon{
		apiKey:     apiKey,
		multiplier: "1",
		timespan:   "minute",
		limit:      maxResultsLimit,
		httpClient: createHttpClient(),
	}
	for _, opt := range opts {
		if err := opt.applyOption(p); err != nil {
			return nil, err
		}
	}
	return p, nil
}

func (p *Polygon) GetTicker(pair *market.Pair) (*market.Ticker, error) {
	op := "polygon.GetTicker"
	symbol := pair.Base.Symbol

	resp, err := p.getLastQuote(symbol)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	result := resp.Results

	return &market.Ticker{
		// TODO: get candle
		Candle: &market.Candle{
			Time:   result.TimeSIP / nanoDivider, // Check this is correct
			Open:   0,
			High:   0,
			Low:    0,
			Close:  0,
			Volume: 0,
		},
		Ask: &market.OrderBookRow{
			Price:       result.AskPrice,
			Volume:      float64(result.AskSize * 100),
			AccumVolume: float64(result.AskSize * 100),
		},
		Bid: &market.OrderBookRow{
			Price:       result.BidPrice,
			Volume:      float64(result.BidSize * 100),
			AccumVolume: float64(result.BidSize * 100),
		},
	}, nil
}

func (p *Polygon) GetHistoricalData(pair *market.Pair, start, end time.Time) ([]market.Candle, error) {
	op := "polygon.GetHistoricalData"
	var marketCandles []market.Candle

	unixMsStart := start.UnixNano() / 1000000
	unixMsEnd := (end.UnixNano() / 1000000) + dayMilliseconds
	symbol := pair.Base.Symbol

	if p.forex {
		symbol = fmt.Sprintf("C:%s%s", pair.Base.Symbol, pair.Quote.Symbol)
	}

	newCandles := 1
	for newCandles != 0 {
		resp, err := p.getRange(symbol, unixMsStart, unixMsEnd)
		if err != nil {
			return nil, ez.Wrap(op, err)
		}

		if resp.Count < 1 {
			return marketCandles, nil
		}

		var lastRes *RangeResult
		var cnt = 0
		for _, item := range resp.Results {
			marketCandles = append(marketCandles, market.Candle{
				Time:   item.T / miliDivider,
				Open:   item.O,
				High:   item.H,
				Low:    item.L,
				Close:  item.C,
				Volume: item.V,
			})
			lastRes = item
			cnt += 1
		}

		unixMsStart = lastRes.T + getOneTimespanUnit(p.multiplier)
		newCandles = cnt
	}

	return marketCandles, nil
}

func (p *Polygon) GetOrderBook(pair *market.Pair) (*market.OrderBook, error) {
	return nil, ez.New("", ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (p *Polygon) getRange(symbol string, unixStart, unixEnd int64) (*RangeResponse, error) {
	op := "polygon.getRange"
	url := fmt.Sprintf(aggRequestFmt, hostPolygon, symbol, p.multiplier, p.timespan, unixStart, unixEnd)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	q := req.URL.Query()
	q.Add("apiKey", p.apiKey)
	q.Add("adjusted", strconv.FormatBool(p.adjusted))
	q.Add("limit", strconv.FormatInt(int64(p.limit), 10))
	req.URL.RawQuery = q.Encode()

	httpResp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, ez.Wrap(op, fmt.Errorf("bad request statusCode=%d", httpResp.StatusCode))
	}

	resp := RangeResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return &resp, nil
}

func (p *Polygon) getLastQuote(symbol string) (*NBBOResponse, error) {
	op := "polygon.getLastQuote"
	url := fmt.Sprintf(lastNBBORequestFmt, hostPolygon, symbol)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	q := req.URL.Query()
	q.Add("apiKey", p.apiKey)
	req.URL.RawQuery = q.Encode()

	httpResp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, ez.Wrap(op, fmt.Errorf("bad request statusCode=%d", httpResp.StatusCode))
	}

	resp := NBBOResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return &resp, nil
}

func getOneTimespanUnit(ts string) int64 {
	switch ts {
	case Minute:
		return timespanMsMinute
	case Hour:
		return timespanMsHour
	case Day:
		return timespanMsDay
	}
	return timespanMsMinute
}
