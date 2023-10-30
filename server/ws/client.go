package ws

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/at-syot/msg_clone/db"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func NewClient(channel *Channel, conn *websocket.Conn) *Client {
	return &Client{
		Id:       uuid.New(),
		Channel:  channel,
		WSConn:   conn,
		finished: make(chan int),
		Egress:   make(chan Message),
	}
}

func (c *Client) ReceiveMessage() {
	for {
		msgType, p, err := c.WSConn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// TODO: handle - socket connection close error
				log.Println(err)
				c.selfClose()
			}

			break
		}

		log.Printf("receive message %s\n", string(p))
		// TODO: will comeback to this context
		message, err := c.saveMessageDB(context.Background(), p)
		if err != nil {
			log.Println(err)
			c.selfClose()
			break
		}

		message.MessageType = msgType
		c.Egress <- message
	}
}

func (c *Client) SendingMessage() {
	for {
		select {
		case message, ok := <-c.Egress:
			if !ok {
				c.selfClose()
				return
			}

			msgBytes, err := json.Marshal(message)
			if err != nil {
				log.Printf("marchal msg err - %s\n", err.Error())
				return
			}

			log.Printf("reply message %+v\n", message)
			for client := range c.Channel.Clients {
				if err := client.WSConn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
					log.Printf("writting message err %s\n", err.Error())
					c.selfClose()
					return
				}
			}
		case <-c.finished:
			return
		}
		/*
			    for goroutine leak testing

			    default:
						time.Sleep(time.Second)
						log.Printf("client id %s is working\n", c.Id.String())
					}
		*/
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
      select ms.id, ms.createdAt, u.username
      from messages as ms
      inner join users as u on u.id = ms.senderId
      where ms.channelId = $1
      order by ms.createdAt DESC
      limit 1`
		err = conn.QueryRow(
			lastInsertQuery,
			[]any{c.Channel.Id.String()},
			&parsedMessage.Id,
			&parsedMessage.CreatedAt,
			&parsedMessage.SenderUsername,
		)
		if err != nil {
			return err
		}
		return nil
	})

	return parsedMessage, nil
}

func (c *Client) selfClose() {
	c.WSConn.Close()
	c.Channel.RemoveClient(c)
	c.Channel.PrintClientsCount()
	c.finished <- 0
}
