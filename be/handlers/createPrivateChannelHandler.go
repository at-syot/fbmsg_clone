package handlers

import (
	"net/http"

	"github.com/at-syot/msg_clone/libs"
	"github.com/at-syot/msg_clone/ws"
	"github.com/google/uuid"
)

func MockData() {
	aot := users.NewUser("aot")
	yok := users.NewUser("yok")
	ch := &ws.Channel{Id: uuid.New(), Name: "aot <> yok", Clients: make(map[*ws.Client]bool)}

	userChannels.AddChannelForUser(aot, ch)
	userChannels.AddChannelForUser(yok, ch)
	userChannels.Print(aot)
	userChannels.Print(yok)
}

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

	u := users.GetByUID(req.CreatorId)
	missingUser := false
	for _, uid := range req.UserIds {
		u := users.GetByUID(uid)
		if u == nil {
			missingUser = true
			return
		}
	}
	if u == nil || missingUser {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// create channel -> get *_Channel
	// bind *Channel to each user
	c := &ws.Channel{
		Id:      uuid.New(),
		Name:    u.Username,
		Clients: make(map[*ws.Client]bool),
	}

	// add created channel into user
	userChannels.AddChannelForUser(u, c)
	for _, uid := range req.UserIds {
		u := users.GetByUID(uid)
		userChannels.AddChannelForUser(u, c)
	}

	res := CreateChannelRes{
		ChannelId:        c.Id.String(),
		CreateChannelReq: req,
	}

	libs.WriteOKRes(w, &res)
}
