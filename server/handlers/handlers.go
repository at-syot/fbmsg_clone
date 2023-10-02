package handlers

import (
	"github.com/at-syot/msg_clone/ws"
	"github.com/gorilla/websocket"
	"net/http"
)

type BaseResponse struct {
	Message string `json:"message"`
}

var (
	channelById = ws.ChannelById{}
	users       = ws.Users{}

	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)
