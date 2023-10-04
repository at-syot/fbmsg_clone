package db

import "github.com/google/uuid"

type (
	User struct {
		Id       uuid.UUID `db:"id" json:"id"`
		Username string    `db:"username" json:"username"`
	}

	Channel struct {
		Id          uuid.UUID `db:"id" json:"id"`
		DisplayName string    `db:"displayname" json:"displayname"`
	}
)
