package main

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

	_Channel struct {
		id        uuid.UUID
		creatorId uuid.UUID
		name      string
		clients   map[*_Client]bool
	}

	_Client struct {
		channel *_Channel
		Id      uuid.UUID
		WSConn  *websocket.Conn
		egress  chan Message

		user *_User
	}

	_User struct {
		Id       uuid.UUID
		username string
	}
)

type Channels []*_Channel
type Users []*_User
type UserChannels map[*_User]*[]*_Channel
type UserClient map[*_User]*_Client

func (users *Users) NewUser(uname string) *_User {
	u := &_User{Id: uuid.New(), username: uname}
	*users = append(*users, u)
	return u
}

func (users *Users) MatchBy(uname string) []*_User {
	matchedUsers := []*_User{}
	for _, u := range *users {
		if strings.Contains(u.username, uname) {
			matchedUsers = append(matchedUsers, u)
		}
	}
	return matchedUsers
}

func (users *Users) GetByUName(uname string) *_User {
	for _, u := range *users {
		if u.username == uname {
			return u
		}
	}
	return nil
}

func (users *Users) GetByUID(id string) *_User {
	for _, u := range *users {
		if u.Id.String() == id {
			return u
		}
	}
	return nil
}

func (users *Users) GetUserContacts(id uuid.UUID) []*_User {
	userContacts := []*_User{}
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

func (uc UserClient) RegisterNewClient(u *_User) *_Client {
	client := &_Client{
		Id:      uuid.New(),
		WSConn:  nil,
		egress:  make(chan Message),
		channel: nil,
		user:    u,
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
func (cs *Channels) addEmptyChannel(name string) *_Channel {
	c := &_Channel{
		id:   uuid.New(),
		name: name,
	}
	*cs = append(*cs, c)
	return c
}

func (cs *Channels) getChannels() []*_Channel {
	return *cs
}

// ----------- UserChannels

func (uc UserChannels) AddChannelForUser(u *_User, ch *_Channel) {
	uChans, ok := userChannels[u]
	if !ok {
		uChans = &[]*_Channel{}
	}
	*uChans = append(*uChans, ch)
	userChannels[u] = uChans
}

func (uc UserChannels) GetUserChannels(u *_User) *[]*_Channel {
	return uc[u]
}

func (uc UserChannels) GetUserChannelWith_UserAndChannelId(u *_User, chanId uuid.UUID) *_Channel {
	chans, ok := uc[u]
	if !ok {
		return nil
	}

	for _, ch := range *chans {
		if ch.id == chanId {
			return ch
		}
	}
	return nil
}

func (uc UserChannels) Print(u *_User) {
	userChannels, ok := uc[u]
	if !ok {
		return
	}

	log.Printf("for user {%s} \n\n", u.username)
	for _, ch := range *userChannels {
		log.Printf("address %p\nchan %+v\n", ch, ch)
	}
}

// --------------- _Client

// receive message
func (c *_Client) receiveMessage() {
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

		c.egress <- Message{
			Id:          uuid.New().String(),
			MessageType: strconv.Itoa(msgType),
			Message:     string(p),
		}
	}
}

// send message
func (c *_Client) sendingMessage() {
	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				return
			}

			log.Printf("reply message %+v\n", message)
			for client, _ := range c.channel.clients {
				if err := client.WSConn.WriteMessage(websocket.TextMessage, []byte(message.Message)); err != nil {
					log.Printf("writting message err %s\n", err.Error())
				}
			}
		}
	}
}
