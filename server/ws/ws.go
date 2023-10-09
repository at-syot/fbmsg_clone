package ws

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/at-syot/msg_clone/db"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

// design

type (
	Message struct {
		Id          string `json:"id"`
		Event       string `json:"event"`
		SenderId    string `json:"senderId"`
		MessageType int    `json:"messageType"`
		Message     string `json:"message"`
		CreatedAt   string `json:"createdAt"`
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
		message, err := c.saveMessageDB(ctx, p)
		if err != nil {
			log.Println(err)
			return
		}

		message.MessageType = msgType
		c.Egress <- message
	}
}

func (c *Client) saveMessageDB(ctx context.Context, p []byte) (Message, error) {
	parsedMessage := Message{}
	err := json.Unmarshal(p, &parsedMessage)
	if err != nil {
		customErr := errors.New("MESSAGE RECEIVED: message format is invalid")
		return Message{}, errors.Join(err, customErr)
	}

	err = db.ExecWithTx(ctx, func(conn db.Conn) error {
		query := `
			INSERT INTO messages (
			  channelId, 
				senderId,
			  content
			) VALUES ($1, $2, $3)`
		if err := conn.Execute(
			query,
			c.Channel.Id.String(),
			parsedMessage.SenderId,
			parsedMessage.Message,
		); err != nil {
			return err
		}

		lastInsertQuery := `
			select ms.id, ms.createdAt
			from messages as ms
			where ms.channelId = $1
			order by ms.createdAt DESC
			limit 1`
		err = conn.QueryRow(
			lastInsertQuery,
			[]any{c.Channel.Id.String()},
			&parsedMessage.Id,
			&parsedMessage.CreatedAt,
		)
		if err != nil {
			return err
		}
		return nil
	})

	return parsedMessage, nil
}

func (c *Client) SendingMessage() {
	for {
		select {
		case message, ok := <-c.Egress:
			if !ok {
				return
			}

			msgBytes, err := json.Marshal(message)
			if err != nil {
				log.Printf("marchal msg err - %s\n", err.Error())
				return
			}

			log.Printf("reply message %+v\n", message)
			for client, _ := range c.Channel.Clients {
				if err := client.WSConn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
					log.Printf("writting message err %s\n", err.Error())
				}
			}
		}
	}
}
