```
POST - api/user
POST - api/create_channel 
       -- body { userIds: []string }

WS - ws/join_channel?channelId="xx"?clientId="dd"
```

---

```go
apiHandlers := []func() { }
createWsHandler := func (*Hub) {
  return func (w http.ResponseWriter, r http.Request) {
   conn, _ := upgrader.Upgrade(w, r, nil)  
	 
  }
}

hub := NewHub()

-- config routes
-- 

start service
```

---


`POST - api/user`
```mermaid
flowchart
   start --> B{Check user exist}
   B-->|exist| C[return username already exist]
   B-->|not exist| D[get UUID for new user]
   D-->E[ clientId := NewUUID
   client := NewClient with clientId
   userIdClient -newUUID- = client]
   E-->F
```

---

```mermaid
sequenceDiagram
   client->>+server: connect api/ws 
   activate server
   Note right of server: server create ws' client connection
   Note right of server: client struct { conn: *WsConn }
   Note right of server: start receive message _ *goroutine
   Note right of server: start sending message _ *goroutine
   deactivate server
   
   client->>+server: send message ws
   activate server
   Note right of server: receive message from client
   Note right of server: send new message to *client.manager.clients.egress channel
   Note right of server: sending message -> get egress channel's message
   Note right of server: sending message -> send message back to client
   server->>+client: reply message 
   deactivate server
```

```mermaid
classDiagram
    class Hub{
        map[*Channel]bool channels
        
        CreateChannel()
        DeleteChannel()
    }
    
    class Channel{
        string Id
        map[*Client]bool Clients
        
        AddClient()
        RemoveClient()
    }
    
    class Client{
        map[*Channel]bool channels
        *WsConn conn
        chan[Message] egress
        
        +RegisterConn(*WsConn)
        ReceiveMessage()
        WriteMessage()
    }

    class Message{
        string type
        string message 
    }
    
```


```mermaid
sequenceDiagram
   client->>+server: create user || login
   Note right of server: 
   
```