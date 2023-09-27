package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type BaseResponse struct {
	Message string `json:"message"`
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	channelId := queries.Get("channelId")
	userId := queries.Get("userId")

	u := users.GetByUID(userId)
	if u == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	chanUUID, err := uuid.Parse(channelId)
	if err != nil {
		log.Printf("uuid.FromBytes err - %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	channel := userChannels.GetUserChannelWith_UserAndChannelId(u, chanUUID)
	if channel == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// upgrade http protocal to ws
	// get ws connection
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	fmt.Println("new connection")

	// wrap conn with _Client
	client := &_Client{
		Id:      uuid.New(),
		channel: channel,
		WSConn:  conn,
		egress:  make(chan Message),
		user:    u,
	}
	channel.clients[client] = true

	go client.receiveMessage()
	go client.sendingMessage()
}

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

// createUserHandler -
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	req := CreateUserReq{}
	if err := ReadReqBody(r, &req); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	indbUser := users.GetByUName(req.Username)
	log.Println(indbUser)
	if indbUser == nil {
		user := users.NewUser(req.Username)
		userClients.RegisterNewClient(user)
		channels.addEmptyChannel(user.username)

		msg := fmt.Sprintf("username: {%s} - created", req.Username)
		resp := CreateUserResp{
			BaseResponse: &BaseResponse{Message: msg},
			UserId:       user.Id.String(),
			Username:     req.Username,
		}
		WriteOKRes(w, &resp)
	} else {
		resp := CreateUserResp{
			BaseResponse: &BaseResponse{Message: "good to go!"},
			UserId:       indbUser.Id.String(),
			Username:     req.Username,
		}
		WriteOKRes(w, &resp)
	}
}

type (
	GetUserChannelRespChannel struct {
		Id   string `json:"id"`
		Name string `json:"channelName"`
	}
	GetUserChannelResp struct {
		Channels []GetUserChannelRespChannel `json:"channels"`
	}
)

// getChannelsHandler - get all available channels by userId
func getUserChannelsHandler(w http.ResponseWriter, r *http.Request) {
	uid := mux.Vars(r)["uid"]
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
	for _, c := range *userChannels.GetUserChannels(u) {
		res.Channels = append(
			res.Channels,
			GetUserChannelRespChannel{Id: c.id.String(), Name: c.name},
		)
	}

	// build channel for this user
	WriteOKRes(w, &res)
}

// ---------

type GetUserRespUser struct {
	Id   string `json:"id"`
	Name string `json:"username"`
}
type GetUserResp struct {
	Users []GetUserRespUser `json:"users"`
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	uname := query.Get("username")

	res := GetUserResp{Users: []GetUserRespUser{}}
	for _, u := range users.MatchBy(uname) {
		res.Users = append(res.Users, GetUserRespUser{Id: u.Id.String(), Name: u.username})
	}

	WriteOKRes(w, &res)
}

// ---------- Create channels

type (
	CreateChannelReq struct {
		CreatorId  string `json:"creatorId"`
		WithUserId string `json:"withUserId"`
	}

	CreateChannelRes struct {
		ChannelId string   `json:"channelId"`
		UserIds   []string `json:"userIds"`
	}
)

func createPrivateChannel(w http.ResponseWriter, r *http.Request) {
	req := CreateChannelReq{}
	if err := ReadReqBody(r, &req); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	u := users.GetByUID(req.CreatorId)
	withUser := users.GetByUID(req.WithUserId)
	if u == nil || withUser == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// create channel -> get *_Channel
	// bind *Channel to each user
	c := &_Channel{
		id:      uuid.New(),
		name:    fmt.Sprintf("%s <> %s", u.username, withUser.username),
		clients: map[*_Client]bool{},
	}

	// add new channel to creator user
	userChannels.AddChannelForUser(u, c)
	userChannels.Print(u)

	// add new channel to withUser user
	userChannels.AddChannelForUser(withUser, c)
	userChannels.Print(withUser)

	res := CreateChannelRes{
		ChannelId: c.id.String(),
		UserIds:   []string{u.Id.String(), withUser.Id.String()},
	}
	WriteOKRes(w, &res)
}
