package genericws

import (
	"github.com/gorilla/websocket"
	"time"
)

const (
	pongWait           = 60 * time.Second
	pingTimeInterval   = 1 * time.Second
)

type wsConnHandler struct {
	host        string
	connRetries int
	conn        *websocket.Conn
}

func NewConnHandler(host string, connRetries int) *wsConnHandler {
	return &wsConnHandler{
		host:        host,
		connRetries: connRetries,
	}
}

func (w *wsConnHandler) Connect() error {
	var err error
	if !w.IsClose() {
		w.Close()
	}
	for i := 0; i < w.connRetries; i++ {
		w.conn, _, err = websocket.DefaultDialer.Dial(w.host, nil)
		if err == nil {
			break
		}
	}
	if w.conn != nil {
		w.conn.SetPongHandler(func(appData string) error {
			w.conn.SetReadDeadline(time.Now().Add(pongWait))
			time.Sleep(pingTimeInterval)
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
		if _, ok := err.(*websocket.CloseError); ok {
			w.conn.Close()
			w.conn = nil
		}
	}
	return err
}

func (w *wsConnHandler) ReadMessage() (p []byte, err error) {
	_, bs, bErr := w.conn.ReadMessage()
	if bErr != nil {
		if _, ok := bErr.(*websocket.CloseError); ok {
			w.conn.Close()
			w.conn = nil
		}
	}
	return bs, bErr
}
