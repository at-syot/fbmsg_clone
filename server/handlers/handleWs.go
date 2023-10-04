package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/at-syot/msg_clone/db"
	"github.com/at-syot/msg_clone/ws"
	"github.com/google/uuid"
	"net/http"
)

type selectUserChannel struct {
	id          string
	displayname string
	creator     bool
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	channelId := queries.Get("channelId")
	userId := queries.Get("userId")

	ctx := r.Context()
	query := `select c.id, c.displayname, cms.creator from channels as c
    inner join channel_members cms on c.id = cms.channelId
    where c.id = $1 and cms.userid = $2`

	result := selectUserChannel{}
	if err := db.QueryRowContext(
		ctx,
		query,
		[]any{channelId, userId},
		&result.id, &result.displayname, &result.creator,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	channelUUID := uuid.MustParse(result.id)
	channel := channelById.GetById(channelUUID)

	// upgrade http to websocket protocol
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	fmt.Printf("connected to channel %s\n", channelUUID)

	// wrap conn with _Client
	client := &ws.Client{
		Id:      uuid.New(),
		Channel: channel,
		WSConn:  conn,
		Egress:  make(chan ws.Message),
	}
	channel.Clients[client] = true

	go client.ReceiveMessage()
	go client.SendingMessage()
}
