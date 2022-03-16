package genericws

import (
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	pongWaitTime     = 60 * time.Second
	pingTimeInterval = 1 * time.Second
	connRetries      = 3
)

type wsConnHandler struct {
	host string
	conn *websocket.Conn
	opts WSHandlerOptions
}

type WSHandlerOptions struct {
	PingTimeInterval time.Duration
	PongWaitTime     time.Duration
}

var DefaultOptions = WSHandlerOptions{
	PingTimeInterval: pingTimeInterval,
	PongWaitTime:     pongWaitTime,
}

func NewConnHandler(host string, opts WSHandlerOptions) *wsConnHandler {
	return &wsConnHandler{
		host: host,
		opts: opts,
	}
}

func (w *wsConnHandler) Connect() error {
	var err error
	if !w.IsClose() {
		w.Close()
	}
	for i := 0; i < connRetries; i++ {
		w.conn, _, err = websocket.DefaultDialer.Dial(w.host, nil)
		if err == nil {
			break
		}
	}
	if w.conn != nil {
		w.conn.SetPongHandler(func(appData string) error {
			w.conn.SetReadDeadline(time.Now().Add(w.opts.PongWaitTime))
			time.Sleep(w.opts.PingTimeInterval)
			return w.conn.WriteMessage(websocket.PingMessage, []byte{})
		})
		err = w.conn.WriteMessage(websocket.PingMessage, []byte{})
	}
	return err
}

func (w *wsConnHandler) Close() error {
	if w.conn != nil {
		return w.Close()
	}
	return nil
}

func (w *wsConnHandler) IsClose() bool {
	return w.conn == nil
}

func (w *wsConnHandler) WriteMessage(messageType int, data []byte) error {
	err := w.conn.WriteMessage(messageType, data)
	if err != nil {
		if isConnError(err) {
			w.conn.Close()
			w.conn = nil
		}
	}
	return err
}

func (w *wsConnHandler) ReadMessage() ([]byte, error) {
	_, bs, err := w.conn.ReadMessage()
	if err != nil {
		if isConnError(err) {
			w.conn.Close()
			w.conn = nil
		}
	}
	return bs, err
}

func isConnError(err error) bool {
	if _, ok := err.(*websocket.CloseError); ok ||
		strings.Contains(err.Error(), "timeout") ||
		strings.Contains(err.Error(), "connection reset by peer") ||
		strings.Contains(err.Error(), "assign requested address") ||
		websocket.IsUnexpectedCloseError(err,
			websocket.CloseAbnormalClosure,
			websocket.CloseGoingAway,
			websocket.CloseProtocolError,
			websocket.CloseInternalServerErr){
		return true
	}
	return false
}