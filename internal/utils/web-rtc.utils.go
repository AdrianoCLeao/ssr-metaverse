package utils

import (
	"net/http"

	"github.com/gorilla/websocket"
	"ssr-metaverse/internal/core/error"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, *error.APIError) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, &error.APIError{
			Code:    500,
			Message: "WebSocket upgrade failed",
		}
	}

	return conn, nil
}
