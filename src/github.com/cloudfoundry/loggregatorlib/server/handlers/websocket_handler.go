package handlers

import (
	"github.com/cloudfoundry/gosteno"
	"github.com/cloudfoundry/loggregatorlib/server"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type websocketHandler struct {
	messages  <-chan []byte
	keepAlive time.Duration
	logger    *gosteno.Logger
}

func NewWebsocketHandler(m <-chan []byte, keepAlive time.Duration, logger *gosteno.Logger) *websocketHandler {
	return &websocketHandler{messages: m, keepAlive: keepAlive, logger: logger}
}

func (h *websocketHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.logger.Debugf("websocket handler: ServeHTTP entered with request %v", r)
	defer h.logger.Debugf("websocket handler: ServeHTTP exited")

	ws, err := websocket.Upgrade(rw, r, nil, 0, 0)
	if err != nil {
		http.Error(rw, "Not a websocket handshake", http.StatusBadRequest)
		h.logger.Debugf("websocket handler: Not a websocket handshake: %s", err.Error())
		return
	}
	defer ws.Close()
	defer ws.WriteControl(websocket.CloseMessage, []byte{}, time.Time{})
	keepAliveExpired := make(chan struct{})

	// TODO: remove this loop (but keep ws.ReadMessage()) once we retire support in the cli for old style keep alives
	go func() {
		for {
			_, _, err := ws.ReadMessage()
			if err != nil {
				return
			}
		}
	}()

	go func() {
		server.NewKeepAlive(ws, h.keepAlive).Run()
		close(keepAliveExpired)
		h.logger.Debugf("websocket handler: Connection from %s timed out", r.RemoteAddr)
	}()

	for {
		select {
		case <-keepAliveExpired:
			return
		case message, ok := <-h.messages:
			if !ok {
				return
			}
			err = ws.WriteMessage(websocket.BinaryMessage, message)
			if err != nil {
				h.logger.Debugf("websocket handler: Error writing to websocket: %s", err.Error())
				return
			}
		}
	}
}
