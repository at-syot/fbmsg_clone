package ws

import (
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

type Channels []*Channel
type Users []*User
type UserChannels map[*User]*[]*Channel
type UserClient map[*User]*Client

func (users *Users) NewUser(uname string) *User {
	u := &User{Id: uuid.New(), Username: uname}
	*users = append(*users, u)
	return u
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

func (users *Users) GetByUName(uname string) *User {
	for _, u := range *users {
		if u.Username == uname {
			return u
		}
	}
	return nil
}

func (users *Users) GetByUID(id string) *User {
	for _, u := range *users {
		if u.Id.String() == id {
			return u
		}
	}
	return nil
}

func (users *Users) GetUserContacts(id uuid.UUID) []*User {
	userContacts := []*User{}
	for _, u := range *users {
		if u.Id != id {
			userContacts = append(userContacts, u)
		}
	}

	return userContacts
}

func (users *Users) Print() {
	for _, u := range *users {
		log.Printf("adress: %p\nuser %+v\n", u, u)
	}
}

// ----------------------

func (uc UserClient) RegisterNewClient(u *User) *Client {
	client := &Client{
		Id:      uuid.New(),
		WSConn:  nil,
		Egress:  make(chan Message),
		Channel: nil,
		User:    u,
	}
	uc[u] = client

	return client
}

func (uc UserClient) Print() {
	for user, client := range uc {
		log.Printf("u: {%s}, client: %+v \n", user.Id.String(), client)
	}
}

// ----------------

// ----------- Channel --------
func (cs *Channels) AddEmptyChannel(name string) *Channel {
	c := &Channel{
		Id:   uuid.New(),
		Name: name,
	}
	*cs = append(*cs, c)
	return c
}

func (cs *Channels) GetChannels() []*Channel {
	return *cs
}

// ----------- UserChannels

func (uc UserChannels) AddChannelForUser(u *User, ch *Channel) {
	uChans, ok := uc[u]
	if !ok {
		uChans = &[]*Channel{}
	}
	*uChans = append(*uChans, ch)
	uc[u] = uChans
}

func (uc UserChannels) GetUserChannels(u *User) *[]*Channel {
	return uc[u]
}

func (uc UserChannels) GetUserChannelWith_UserAndChannelId(u *User, chanId uuid.UUID) *Channel {
	chans, ok := uc[u]
	if !ok {
		return nil
	}

	for _, ch := range *chans {
		if ch.Id == chanId {
			return ch
		}
	}
	return nil
}

func (uc UserChannels) Print(u *User) {
	userChannels, ok := uc[u]
	if !ok {
		return
	}

	log.Printf("for user {%s} \n\n", u.Username)
	for _, ch := range *userChannels {
		log.Printf("address %p\nchan %+v\n", ch, ch)
	}
}

// ReceiveMessage to receive connected's messages
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

// send message
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
