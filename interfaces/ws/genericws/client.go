package genericws

import (
	"context"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/interfaces/ws"
)

const (
	waitTimeForNewConn     = 1000 * time.Millisecond
	connectionRetries      = 3
	subscriptionVerifyTime = time.Second * 5
)

type baseClient struct {
	subscriptionPairs []market.Pair
	handler           WebsocketHandler
	name              string
	// buffer for bursts or spikes in data
	buffer                        int
	overrideExpectedVerifications int
	channels                      []ChannelOpts
}

func NewClient(handler WebsocketHandler, opts ...Option) (*baseClient, error) {
	const op = "ws.NewClient"

	c := &baseClient{
		subscriptionPairs: []market.Pair{},
		handler:           handler,
		buffer:            1000,
		channels:          defaultChannels,
	}
	for _, opt := range opts {
		if err := opt.applyOption(c); err != nil {
			return nil, ez.Wrap(op, err)
		}
	}

	return c, nil
}

func (c *baseClient) createConnection(ctx context.Context) (*wsConnHandler, error) {
	const op = "baseClient.createConnection"
	var err error

	if iErr := c.handler.Init(HandlerOptions{
		Pairs:    c.subscriptionPairs,
		Channels: c.channels,
	}); iErr != nil {
		return nil, ez.Wrap(op, iErr)
	}

	settings, err := c.handler.GetSettings()
	if err != nil {
		return nil, ez.New(op, ez.EINTERNAL, "", err)
	}

	wsHandler := NewConnHandler(settings.Endpoint, getOptionsFromSettings(settings))
	c.overrideExpectedVerifications = settings.SubscriptionVerificationCount

	err = wsHandler.Connect()
	if err != nil {
		errMsg := fmt.Sprintf("failed to dial websocket: %s", err.Error())
		return wsHandler, ez.New(op, ez.EINTERNAL, errMsg, err)
	}

	for j := 0; j < connectionRetries; j++ {
		err = c.subscribe(ctx, wsHandler)
		if err == nil {
			return wsHandler, nil
		}
	}

	errMsg := fmt.Sprintf("failed to subscribe to websocket channel: %s", err.Error())
	return wsHandler, ez.New(op, ez.EINTERNAL, errMsg, err)
}

func (c *baseClient) subscribe(ctx context.Context, wsConn *wsConnHandler) error {
	const op = "baseClient.subscribe"

	requests, err := c.handler.GetSubscriptionsRequests()
	if err != nil {
		msgErr := fmt.Sprintf("error creating subscription request: %s", err.Error())
		return ez.New(op, ez.EINTERNAL, msgErr, err)
	}

	if len(requests) == 0 {
		return fmt.Errorf("at least one subscription must be set")
	}

	for _, request := range requests {
		err = wsConn.WriteMessage(websocket.TextMessage, request)
		if err != nil {
			return ez.New(op, ez.EINTERNAL, "failed to write message", err)
		}
	}

	expectedRequests := len(requests)
	if c.overrideExpectedVerifications != 0 {
		expectedRequests = c.overrideExpectedVerifications
	}
	return c.verifySubscriptions(ctx, wsConn, expectedRequests)
}

func (c *baseClient) verifySubscriptions(ctx context.Context, wsConn *wsConnHandler, expectedRequest int) error {
	const op = "baseClient.verifySubscriptions"

	ctxDeadline, cancel := context.WithTimeout(ctx, subscriptionVerifyTime)
	defer cancel()

	for {
		select {
		case <-ctxDeadline.Done():
			return ez.New(op, ez.EINTERNAL, "failed to get subscription confirmations", nil)

		default:
			bs, err := wsConn.ReadMessage()
			if err != nil {
				return ez.Wrap(op, err)
			}

			err = c.handler.VerifySubscriptionResponse(bs)
			if err != nil {
				return ez.Wrap(op, err)
			}

			expectedRequest--
			if expectedRequest == 0 {
				return nil
			}
		}
	}
}

func (c *baseClient) Listen(ctx context.Context) (<-chan ws.ListenChan, error) {
	const op = "ws.Listen"

	if len(c.subscriptionPairs) == 0 {
		return nil, ez.New(op, ez.EINVALID, ErrSubscriptionPairs, nil)
	}

	wsConn := &wsConnHandler{}

	listenChan := make(chan ws.ListenChan, c.buffer)
	errCounter := 0

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(listenChan)
				return
			default:
				if wsConn == nil || wsConn.IsClose() {
					var cErr error
					wsConn, cErr = c.createConnection(ctx)
					if cErr != nil {
						log.Error().
							Str("OP", op).
							Str("Exchange", c.name).
							Err(cErr).
							Msg("error creating connection")
					}
					if wsConn == nil || wsConn.IsClose() {
						time.Sleep(waitTimeForNewConn)
						continue
					}
				}

				bs, bErr := wsConn.ReadMessage()
				if bErr != nil {
					log.Error().Str("Exchange", c.name).Err(bErr).Msg("Error reading message")
					if errCounter > 5 {
						time.Sleep(waitTimeForNewConn)
						wsConn, _ = c.createConnection(ctx)
						if wsConn.IsClose() {
							continue
						}
					} else {
						errCounter = errCounter + 1
						continue
					}
				}

				data, pErr := c.handler.Parse(bs)
				if pErr != nil {
					log.Error().
						Str("OP", op).
						Str("exchange", c.name).
						Str("bytes", string(bs)).
						Err(pErr).
						Msg("error unmarshalling order book data from ws")
					continue
				}

				if data != nil {
					listenChan <- *data
				}
			}
		}
	}()

	return listenChan, nil
}

func getOptionsFromSettings(settings Settings) WSHandlerOptions {
	opts := WSHandlerOptions{
		PingTimeInterval: pingTimeInterval,
		PongWaitTime:     pongWaitTime,
	}

	if settings.PingTimeInterval > 0 {
		opts.PingTimeInterval = settings.PingTimeInterval
	}

	if settings.PongWaitTime > 0 {
		opts.PongWaitTime = settings.PongWaitTime
	}

	return opts
}
