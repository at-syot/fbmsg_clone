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
	getUserChannelRespChanUser struct {
		Id       string `json:"id"`
		Username string `json:"username"`
		Creator  bool   `json:"creator"`
	}
	getUserChannelRespChan struct {
		Id      string                       `json:"id"`
		Name    *string                      `json:"displayname"`
		Creator bool                         `json:"creator"`
		Users   []getUserChannelRespChanUser `json:"users"`
	}
	getUserChannelResp struct {
		Channels []getUserChannelRespChan `json:"channels"`
	}
)

// GetUserChannelsHandler - get all available channels by userId
func GetUserChannelsHandler(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")

	ctx := r.Context()
	query := `select 
    c.displayname, 
    cm.channelId, 
    cm.creator 
	from channel_members as cm
	inner join channels as c on c.id =  cm.channelId
	where cm.userId = $1 
	group by c.displayname, cm.channelId, cm.creator`

	// get user channels
	res := getUserChannelResp{Channels: []getUserChannelRespChan{}}
	err := db.QueryContext(ctx, query, []any{uid}, func(rows *sql.Rows) error {
		resChannel := getUserChannelRespChan{Users: make([]getUserChannelRespChanUser, 0)}
		scanErr := rows.Scan(&resChannel.Name, &resChannel.Id, &resChannel.Creator)
		if scanErr != nil {
			return scanErr
		}

		// TODO: super dirty here -> need better solution
		// get each channel members
		query := `
		select
		    u.id as userId,
		    u.username,
		    cm.creator
		from channel_members as cm
		inner join channels as c on c.id = cm.channelid
		inner join users as u on u.id = cm.userId
		where cm.channelId = $1`
		db.QueryContext(ctx, query, []any{resChannel.Id}, func(rows *sql.Rows) error {
			chanUser := getUserChannelRespChanUser{}
			if err := rows.Scan(&chanUser.Id, &chanUser.Username, &chanUser.Creator); err != nil {
				return err
			}

			resChannel.Users = append(resChannel.Users, chanUser)
			return nil
		})

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
