package handlers

import (
	"log"
	"net/http"

	"github.com/at-syot/msg_clone/db"
	"github.com/at-syot/msg_clone/libs"
)

func ClearDataHandler(w http.ResponseWriter, r *http.Request) {
	stmt := `
  delete from channel_members;
  delete from messages;
  delete from channels;
  delete from users;
  insert into users (username) values ('creator');`
	_, err := db.DB.Exec(stmt)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	libs.WriteOKRes(w, BaseResponse{Message: "clear data success."})
}
