package ws

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// design

type (
	Message struct {
		Id          string
		MessageType string
		Message     string
	}

	Channel struct {
		Id        uuid.UUID
		CreatorId uuid.UUID
		Name      string
		Clients   map[*Client]bool
	}

	Client struct {
		Channel *Channel
		Id      uuid.UUID
		WSConn  *websocket.Conn
		Egress  chan Message

		User *User
	}

	User struct {
		Id       uuid.UUID
		Username string
	}
)

type Users []*User
type ChannelById map[uuid.UUID]*Channel

// ------------

func NewChannel(id uuid.UUID) *Channel {
	return &Channel{
		Id:      id,
		Clients: make(map[*Client]bool),
	}
}

// ---------------

func (chs ChannelById) GetById(id uuid.UUID) *Channel {
	for chanId, ch := range chs {
		if chanId == id {
			return ch
		}
	}

	return nil
}

func (chs ChannelById) Print() {
	for chanId, ch := range chs {
		fmt.Printf("channelId %s, ch %v\n", chanId, ch)
	}
}

func (users *Users) MatchBy(uname string) []*User {
	matchedUsers := []*User{}
	for _, u := range *users {
		if strings.Contains(u.Username, uname) {
			matchedUsers = append(matchedUsers, u)
		}
	}
	return matchedUsers
}

func (c *Client) ReceiveMessage() {
	for {
		msgType, p, err := c.WSConn.ReadMessage()
		if err != nil {

			log.Printf("reading message err - %s\n", err.Error())
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// TODO: handle - socket connection unexpected close
			}

			return
		}

		log.Printf("receive message %s\n", string(p))

		c.Egress <- Message{
			Id:          uuid.New().String(),
			MessageType: strconv.Itoa(msgType),
			Message:     string(p),
		}
	}
}

func (c *Client) SendingMessage() {
	for {
		select {
		case message, ok := <-c.Egress:
			if !ok {
				return
			}

			log.Printf("reply message %+v\n", message)
			for client, _ := range c.Channel.Clients {
				if err := client.WSConn.WriteMessage(websocket.TextMessage, []byte(message.Message)); err != nil {
					log.Printf("writting message err %s\n", err.Error())
				}
			}
		}
	}
}
