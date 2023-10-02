package handlers

import (
	"database/sql"
	"errors"
	"github.com/at-syot/msg_clone/db"
	"github.com/google/uuid"
	"net/http"

	"github.com/at-syot/msg_clone/libs"
)

type (
	CreateUserReq struct {
		Username string `json:"username"`
	}
	CreateUserResp struct {
		*BaseResponse
		UserId   string `json:"userId"`
		Username string `json:"username"`
	}
)

type (
	User struct {
		ID       uuid.UUID `db:"id"`
		Username string    `db:"username"`
	}
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	req := CreateUserReq{}
	if err := libs.ReadReqBody(r, &req); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	user := User{}
	ctx := r.Context()
	err := db.ExecWithTx(ctx, func(conn db.Conn) error {
		userQuery := `SELECT * FROM users WHERE username = $1`
		if err := conn.QueryRow(
			userQuery,
			[]any{req.Username},
			&user.ID, &user.Username,
		); err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}

			insertQuery := "INSERT INTO users (username) VALUES ($1) RETURNING *"
			if err := conn.QueryRow(
				insertQuery,
				[]any{req.Username},
				&user.ID, &user.Username,
			); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	libs.WriteOKRes(w, CreateUserResp{UserId: user.ID.String(), Username: user.Username})
}
