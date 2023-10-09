package db

import (
	"github.com/google/uuid"
	"time"
)

type (
	User struct {
		Id       uuid.UUID `db:"id" json:"id"`
		Username string    `db:"username" json:"username"`
	}

	Channel struct {
		Id          uuid.UUID `db:"id" json:"id"`
		DisplayName string    `db:"displayname" json:"displayname"`
	}

	Message struct {
		Id        int       `db:"id"`
		ChannelId uuid.UUID `db:"channelId"`
		Sender    uuid.UUID `db:"senderId"`
		Content   string    `db:"content"`
		CreatedAt time.Time `db:"createdAt"`
	}
)
