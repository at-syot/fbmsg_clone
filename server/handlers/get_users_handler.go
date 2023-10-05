package handlers

import (
	"database/sql"
	"github.com/at-syot/msg_clone/db"
	"github.com/at-syot/msg_clone/libs"
	"net/http"
)

type GetUserResp struct {
	Users []db.User `json:"users"`
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	uname := r.URL.Query().Get("username")

	ctx := r.Context()
	res := GetUserResp{Users: make([]db.User, 0)}
	query := `
		SELECT * FROM users 
		WHERE username LIKE '%' || $1 || '%'`
	db.QueryContext(ctx, query, []any{uname}, func(rows *sql.Rows) error {
		u := db.User{}
		if err := rows.Scan(&u.Id, &u.Username); err != nil {
			return err
		}

		res.Users = append(res.Users, u)
		return nil
	})

	libs.WriteOKRes(w, &res)
}
