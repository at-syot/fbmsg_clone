package ws

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type (
	Message struct {
		Id             string `json:"id"`
		Event          string `json:"event"`
		SenderId       string `json:"senderId"`
		SenderUsername string `json:"username"`
		MessageType    int    `json:"messageType"`
		Message        string `json:"message"`
		CreatedAt      string `json:"createdAt"`
	}

	Channel struct {
		sync.Mutex
		Id        uuid.UUID
		CreatorId uuid.UUID
		Name      string
		Clients   map[*Client]bool
	}

	Client struct {
		Channel  *Channel
		Id       uuid.UUID
		WSConn   *websocket.Conn
		finished chan int
		Egress   chan Message

		User *User
	}

	User struct {
		Id       uuid.UUID
		Username string
	}
)

type Users []*User
type ChannelById map[uuid.UUID]*Channel

// ---- Channel

func NewChannel(id uuid.UUID) *Channel {
	return &Channel{
		Id:      id,
		Clients: make(map[*Client]bool),
	}
}

func (c *Channel) AddClient(client *Client) {
	c.Lock()
	defer c.Unlock()
	c.Clients[client] = true
}

func (c *Channel) RemoveClient(client *Client) {
	c.Lock()
	defer c.Unlock()
	delete(c.Clients, client)
}

func (c *Channel) PrintClientsCount() {
	log.Printf("client count %d\n", len(c.Clients))
}

func (c *Channel) Print() {
	log.Printf("channelId %s\n", c.Id.String())
	for client := range c.Clients {
		log.Printf("client %+v\n", client)
	}
}

// ------ ChannelById

func (chs ChannelById) GetById(id uuid.UUID) *Channel {
	for chanId, ch := range chs {
		if chanId == id {
			return ch
		}
	}

	return nil
}

func (chs ChannelById) Print() {
	for chId, ch := range chs {
		log.Printf("k: %s, -- v : %+v\n", chId, ch)
		ch.Print()
	}
}
