package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/at-syot/msg_clone/ws"
	"github.com/google/uuid"
)

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	channelId := queries.Get("channelId")
	userId := queries.Get("userId")

	u := users.GetByUID(userId)
	if u == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	chanUUID, err := uuid.Parse(channelId)
	if err != nil {
		log.Printf("uuid.FromBytes err - %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	channel := userChannels.GetUserChannelWith_UserAndChannelId(u, chanUUID)
	if channel == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// upgrade http protocal to ws
	// get ws connection
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	fmt.Println("new connection")

	// wrap conn with _Client
	client := &ws.Client{
		Id:      uuid.New(),
		Channel: channel,
		WSConn:  conn,
		Egress:  make(chan ws.Message),
		User:    u,
	}
	channel.Clients[client] = true

	go client.ReceiveMessage()
	go client.SendingMessage()
}
