package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// -------- end design

// -----------------

var (
	channels     = Channels{}
	users        = Users{}
	userChannels = UserChannels{}

	userClients = UserClient{}

	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

//----------------- client

type ClientList map[*Client]bool
type Client struct {
	conn    *websocket.Conn
	manager *Manager
	egress  chan []byte
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		conn:    conn,
		manager: manager,
		egress:  make(chan []byte),
	}
}

func (c *Client) readMessage() {
	for {
		msgType, payload, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message %v \n", err)
			}
			break
		}

		log.Printf("messageType %v\n", msgType)
		log.Printf("payload %s \n", string(payload))
		log.Println("write to all egress")
		for client := range c.manager.clients {
			client.egress <- payload
		}
	}
}

func (c *Client) writeMessage() {
	for {
		select {
		case payload, ok := <-c.egress:
			if !ok {
				// client closed
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, payload); err != nil {
				log.Printf("answer back err %v\n", err)
			}
			log.Println("reply back success")
		}
	}
}

// --------------- manager

type Manager struct {
	clients ClientList
	sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{clients: make(ClientList)}
}

func (m *Manager) clientLength() int {
	return len(m.clients)
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.clients[client] = true
	log.Printf("clients %v", m.clientLength())
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		client.conn.Close()
		delete(m.clients, client)
	}
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from any origin
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Allow specific HTTP methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Allow specific HTTP headers
		//w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Continue processing the request
		next.ServeHTTP(w, r)
	})
}

func main() {
	aot := users.NewUser("aot")
	yok := users.NewUser("yok")
	ch := &_Channel{id: uuid.New(), name: "aot <> yok", clients: make(map[*_Client]bool)}

	userChannels.AddChannelForUser(aot, ch)
	userChannels.AddChannelForUser(yok, ch)
	userChannels.Print(aot)
	userChannels.Print(yok)

	r := mux.NewRouter()
	//r.Headers("Content-Type", "application/json")

	r.HandleFunc("/user", createUserHandler).Methods("POST")
	r.HandleFunc("/users", getUsers)
	r.HandleFunc("/channels/private", createPrivateChannel).Methods("POST")
	r.HandleFunc("/users/{uid}/channels", getUserChannelsHandler)
	r.HandleFunc("/ws", wsHandler)

	log.Println("server start pn port :3000")

	http.Handle("/", r)
	if err := http.ListenAndServe(":3000", enableCORS(r)); err != nil {
		log.Fatal(err)
	}
}
