package handlers

import (
	"github.com/at-syot/msg_clone/libs"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type (
	GetUserChannelRespChannel struct {
		Id   string `json:"id"`
		Name string `json:"channelName"`
	}
	GetUserChannelResp struct {
		Channels []GetUserChannelRespChannel `json:"channels"`
	}
)

// GetChannelsHandler - get all available channels by userId
func GetUserChannelsHandler(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")
	if users.GetByUID(uid) == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	u := users.GetByUID(uid)
	if u == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	res := GetUserChannelResp{Channels: []GetUserChannelRespChannel{}}
	_userChannels := userChannels.GetUserChannels(u)
	if _userChannels == nil {
		libs.WriteOKRes(w, &res)
		return
	}

	for _, c := range *_userChannels {
		res.Channels = append(
			res.Channels,
			GetUserChannelRespChannel{Id: c.Id.String(), Name: c.Name},
		)
	}

	// build channel for this user
	libs.WriteOKRes(w, &res)
}
