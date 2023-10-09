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
	getChanMsgResponse struct {
		ChannelId string                  `json:"channelId"`
		Messages  []getChanMsgResponseRow `json:"messages"`
	}

	getChanMsgResponseRow struct {
		MessageId string `json:"id"`
		SenderId  string `json:"senderId"`
		Username  string `json:"username"`
		Message   string `json:"message"`
		CreatedAt string `json:"createdAt"`
	}
)

func GetChannelMessagesHandler(w http.ResponseWriter, r *http.Request) {
	channelId := chi.URLParam(r, "cid")

	ctx := r.Context()
	query := `select
    m.id as messageId,
    m.senderId,
    u.username,
    m.content as message,
    m.createdAt
	from channels as c
	inner join messages as m on c.id = m.channelId
	inner join users as u on u.id = m.senderId
	where c.id = $1
	order by m.createdAt ASC`

	res := getChanMsgResponse{channelId, make([]getChanMsgResponseRow, 0)}
	err := db.QueryContext(ctx, query, []any{channelId}, func(rows *sql.Rows) error {
		row := getChanMsgResponseRow{}
		if err := rows.Scan(&row.MessageId, &row.SenderId, &row.Username, &row.Message, &row.CreatedAt); err != nil {
			return err
		}
		res.Messages = append(res.Messages, row)
		return nil
	})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	libs.WriteOKRes(w, &res)
}
