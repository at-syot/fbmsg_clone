package handlers

import (
	"database/sql"
	"errors"
	"github.com/at-syot/msg_clone/db"
	"github.com/at-syot/msg_clone/ws"
	"github.com/google/uuid"
	"net/http"

	"github.com/at-syot/msg_clone/libs"
)

type (
	CreateChannelReq struct {
		CreatorId string   `json:"creatorId"`
		UserIds   []string `json:"userIds"`
	}

	CreateChannelRes struct {
		ChannelId string `json:"channelId"`
		CreateChannelReq
	}
)

func CreateChannel(w http.ResponseWriter, r *http.Request) {
	req := CreateChannelReq{}
	if err := libs.ReadReqBody(r, &req); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if len(req.UserIds) == 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var newChannelId string
	ctx := r.Context()
	err := db.ExecWithTx(ctx, func(conn db.Conn) error {
		// check channel members
		var checkMemberErr error
		memberIds := []string{req.CreatorId}
		memberIds = append(memberIds, req.UserIds...)
		for _, uid := range memberIds {
			if checkMemberErr = conn.QueryRow(
				`SELECT id FROM users WHERE id = $1`,
				[]any{uid}, &[]byte{},
			); checkMemberErr != nil {
				break
			}
		}
		if checkMemberErr != nil {
			return checkMemberErr
		}

		// creating channel
		if err := conn.QueryRow(
			`INSERT INTO channels (createdBy) VALUES ($1) RETURNING id`,
			[]any{req.CreatorId}, &newChannelId,
		); err != nil {
			return err
		}

		// register members into channel
		insertCreatorQuery := `
			INSERT INTO channel_members (channelId, userId, creator)
			VALUES ($1, $2, $3)`
		if err := conn.Execute(insertCreatorQuery, newChannelId, req.CreatorId, true); err != nil {
			return err
		}

		var insertMemberErr error
		for _, uid := range req.UserIds {
			insertMemberQuery := `
				INSERT INTO channel_members (channelId, userId)
				VALUES ($1, $2)`
			if insertMemberErr = conn.Execute(insertMemberQuery, newChannelId, uid); insertMemberErr != nil {
				break
			}
		}
		if insertMemberErr != nil {
			return insertMemberErr
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	channelUUID := uuid.MustParse(newChannelId)
	channelById[channelUUID] = ws.NewChannel(channelUUID)

	res := CreateChannelRes{
		ChannelId:        newChannelId,
		CreateChannelReq: req,
	}
	libs.WriteOKRes(w, &res)
}
