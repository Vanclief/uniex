package polygon

import (
	"github.com/stretchr/testify/assert"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/finmod/utils"
	"os"
	//"reflect"
	"testing"
	"time"
)

const (
	polygonEnvApiKey = "POLYGON_API_KEY"
	defaultApiKey    = "UPp4fgGmfPIGsM4630r3c2vxnOhi26P8"
)

var (
	apiKey = defaultApiKey
)

func init() {
	ap := os.Getenv(polygonEnvApiKey)
	if ap != "" {
		apiKey = ap
	}
}

func createDate(t *testing.T, date string) time.Time {
	d, err := utils.RawDateToRFC3339(date)
	assert.Nil(t, err)
	return d
}

func TestPolygon_GetHistoricalDataStocks(t *testing.T) {
	p, err := New(apiKey)
	assert.Nil(t, err)

	type args struct {
		pair  *market.Pair
		start time.Time
		end   time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "should pass",
			args: args{
				pair: &market.Pair{
					Base: &market.Asset{
						Symbol: "AAPL",
						Name:   "APPLE",
					},
				},
				start: createDate(t, "2021-05-25"),
				end:   createDate(t, "2021-06-25"),
			},
			want: 17932,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.GetHistoricalData(tt.args.pair, tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHistoricalData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Len(t, got, tt.want)
		})
	}

}

func TestPolygon_GetHistoricalDataForex(t *testing.T) {
	p, err := New(apiKey,
		WithForexMarket())
	assert.Nil(t, err)

	type args struct {
		pair  *market.Pair
		start time.Time
		end   time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "should pass",
			args: args{
				pair: &market.Pair{
					Base: &market.Asset{
						Symbol: "EUR",
					},
					Quote: &market.Asset{
						Symbol: "USD",
					},
				},
				start: createDate(t, "2021-05-25"),
				end:   createDate(t, "2021-06-25"),
			},
			want: 34188,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.GetHistoricalData(tt.args.pair, tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHistoricalData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Len(t, got, tt.want)
		})
	}

}

func TestPolygon_GetTicker(t *testing.T) {
	p, err := New(apiKey)
	assert.Nil(t, err)

	// TODO: mock data
	type args struct {
		pair *market.Pair
	}
	tests := []struct {
		name    string
		args    args
		want    *market.Ticker
		wantErr bool
	}{
		{
			name: "should pass",
			args: args{
				pair: &market.Pair{
					Base: &market.Asset{
						Symbol: "AAPL",
						Name:   "APPLE",
					},
				},
			},
			want: &market.Ticker{
				Candle: &market.Candle{
					Time:   1625011101976659083,
					Open:   0,
					High:   0,
					Low:    0,
					Close:  0,
					Volume: 0,
				},
				Ask: &market.OrderBookRow{
					Price:       0,
					Volume:      0,
					AccumVolume: 0,
				},
				Bid: &market.OrderBookRow{
					Price:       0,
					Volume:      0,
					AccumVolume: 0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.GetTicker(tt.args.pair)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTicker() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, got)
			//assert.Equal(t, tt.want.Candle, got.Candle)
			//assert.Equal(t, tt.want.Ask, got.Ask)
			//assert.Equal(t, tt.want.Bid, got.Bid)
		})
	}
}
