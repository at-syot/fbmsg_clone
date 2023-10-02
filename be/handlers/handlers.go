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
	channels     = ws.Channels{}
	users        = ws.Users{}
	userChannels = ws.UserChannels{}
	userClients  = ws.UserClient{}

	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)
