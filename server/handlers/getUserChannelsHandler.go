package handlers

import (
	"database/sql"
	"errors"
	"github.com/at-syot/msg_clone/db"
	"github.com/at-syot/msg_clone/libs"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type (
	GetUserChannelRespChannel struct {
		Id      string `json:"id"`
		Name    string `json:"channelName"`
		Creator string `json:"creator"`
	}
	GetUserChannelResp struct {
		Channels []GetUserChannelRespChannel `json:"channels"`
	}
)

// GetUserChannelsHandler - get all available channels by userId
func GetUserChannelsHandler(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")

	ctx := r.Context()
	query := `select c.displayname, cm.channelId, cm.creator from channel_members as cm
		inner join channels as c on c.id =  cm.channelId
		where cm.userId = $1 
		group by c.displayname, cm.channelId, cm.creator`

	res := GetUserChannelResp{Channels: []GetUserChannelRespChannel{}}
	err := db.QueryContext(ctx, query, []any{uid}, func(rows *sql.Rows) error {
		resChannel := GetUserChannelRespChannel{}
		scanErr := rows.Scan(&resChannel.Name, &resChannel.Id, &resChannel.Creator)
		if scanErr != nil {
			return scanErr
		}

		res.Channels = append(res.Channels, resChannel)
		return nil
	})

	if err != nil {
		log.Printf("query context err - %s\n", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	libs.WriteOKRes(w, &res)
	return
}
