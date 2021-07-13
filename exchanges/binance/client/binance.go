package binanceclient

import (
	"context"
	"os"

	// 	"os"
	// 	"strings"
	// 	"time"

	goBinance "github.com/binance-exchange/go-binance"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	// 	"github.com/go-kit/kit/log"
	// 	"github.com/go-kit/kit/log/level"
	// 	"github.com/vanclief/ez"
	// 	"github.com/vanclief/go-trading-engine/config"
	// 	"github.com/vanclief/go-trading-engine/market"
	// 	"github.com/vanclief/go-trading-engine/utils"
)

// Client - Binance struct that contains the client for API calls and a context cancellable function
type Client struct {
	service   goBinance.Binance
	ctxCancel context.CancelFunc
}

func New(apiKey, secretKey string) *Client {

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = level.NewFilter(logger, level.AllowAll())
	logger = log.With(logger, "time", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	hmacSigner := &goBinance.HmacSigner{
		Key: []byte(secretKey),
	}

	ctx, cancel := context.WithCancel(context.Background())
	// use second return value for cancelling request
	binanceService := goBinance.NewAPIService(
		"https://www.binance.com",
		apiKey,
		hmacSigner,
		logger,
		ctx,
	)

	b := goBinance.NewBinance(binanceService)

	return &Client{
		service:   b,
		ctxCancel: cancel, // TODO: Did not knew what to do with cancel so passed it as parameter to struct -officialnoria
	}
}
