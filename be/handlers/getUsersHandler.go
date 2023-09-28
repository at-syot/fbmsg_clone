package handlers

import (
	"github.com/at-syot/msg_clone/libs"
	"net/http"
)

type GetUserRespUser struct {
	Id   string `json:"id"`
	Name string `json:"username"`
}
type GetUserResp struct {
	Users []GetUserRespUser `json:"users"`
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	uname := query.Get("username")

	res := GetUserResp{Users: []GetUserRespUser{}}
	for _, u := range users.MatchBy(uname) {
		res.Users = append(res.Users, GetUserRespUser{Id: u.Id.String(), Name: u.Username})
	}

	libs.WriteOKRes(w, &res)
}
