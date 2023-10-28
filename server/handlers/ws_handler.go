package handlers

import (
	"database/sql"
	"errors"
	"github.com/at-syot/msg_clone/db"
	"github.com/at-syot/msg_clone/ws"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type selectUserChannel struct {
	id          string
	displayname *string
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
		log.Println(err)
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
		log.Println(err)
		return
	}
	log.Printf("connected to channel %s\n", channelUUID)

	// wrap conn with Client
	client := ws.NewClient(channel, conn)
	channel.AddClient(client)

	go client.ReceiveMessage()
	go client.SendingMessage()
}
