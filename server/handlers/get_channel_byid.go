package handlers

import (
	"database/sql"
	"github.com/at-syot/msg_clone/db"
	"github.com/at-syot/msg_clone/libs"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type (
	getChanByIdResUser struct {
		UserId   string `json:"id"`
		Username string `json:"username"`
		Creator  bool   `json:"creator"`
	}
	getChanByIdRes struct {
		ChannelId   string               `json:"id"`
		DisplayName *string              `json:"displayname"`
		Users       []getChanByIdResUser `json:"users"`
	}
)

func GetChannelById(w http.ResponseWriter, r *http.Request) {
	chanId := chi.URLParam(r, "cid")
	log.Println(chanId)

	ctx := r.Context()
	query := `select
    c.id as channelId,
    c.displayname,
    u.id as userId,
    u.username,
    cm.creator
	from channel_members as cm
	inner join channels as c on c.id = cm.channelid
	inner join users as u on u.id = cm.userId
	where cm.channelId = $1`

	res := getChanByIdRes{Users: make([]getChanByIdResUser, 0)}
	err := db.QueryContext(ctx, query, []any{chanId}, func(rows *sql.Rows) error {
		chanUser := getChanByIdResUser{}
		if err := rows.Scan(
			&res.ChannelId,
			&res.DisplayName,
			&chanUser.UserId,
			&chanUser.Username,
			&chanUser.Creator,
		); err != nil {
			log.Printf("scan err - %s\n", err.Error())
			return err
		}

		res.Users = append(res.Users, chanUser)
		return nil
	})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	libs.WriteOKRes(w, &res)
}
