package ws

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/at-syot/msg_clone/db"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"sync"
)

// design

type (
	Message struct {
		Id          string
		MessageType string
		Message     string
	}

	Channel struct {
		sync.Mutex
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

func (c *Channel) AddClient(client *Client) {
	c.Lock()
	defer c.Unlock()
	c.Clients[client] = true
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

// ------- Client

func NewClient(channel *Channel, conn *websocket.Conn) *Client {
	return &Client{
		Id:      uuid.New(),
		Channel: channel,
		WSConn:  conn,
		Egress:  make(chan Message),
	}
}

func (c *Client) ReceiveMessage() {
	for {
		msgType, p, err := c.WSConn.ReadMessage()
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		if err != nil {

			log.Printf("reading message err - %s\n", err.Error())
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// TODO: handle - socket connection unexpected close
			}

			cancel()
			return
		}

		log.Printf("receive message %s\n", string(p))
		if err = c.saveMessageDB(ctx, p); err != nil {
			log.Println(err)
			return
		}

		c.Egress <- Message{
			Id:          uuid.New().String(),
			MessageType: strconv.Itoa(msgType),
			Message:     string(p),
		}
	}
}

func (c *Client) saveMessageDB(ctx context.Context, p []byte) error {
	parsedMessage := struct {
		SenderId string `json:"senderId"`
		Event    string `json:"event"`
		Content  string `json:"message"`
	}{}
	err := json.Unmarshal(p, &parsedMessage)
	if err != nil {
		customErr := errors.New("MESSAGE RECEIVED: message format is invalid")
		return errors.Join(err, customErr)
	}

	err = db.ExecWithTx(ctx, func(conn db.Conn) error {
		query := `
			INSERT INTO messages (
			  channelId, 
				senderId,
			  content
			) VALUES ($1, $2, $3)`
		if err := conn.Execute(query, c.Channel.Id.String(), parsedMessage.SenderId, parsedMessage.Content); err != nil {
			return err
		}
		return nil
	})

	return err
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
