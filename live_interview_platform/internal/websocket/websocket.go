package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Declare an upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var RoomManagerInstance = &RoomManager{
	Rooms: make(map[string]*Room),
	Mu:    sync.Mutex{},
}

var Wg = &sync.WaitGroup{}

func HandleWebsocket(c *gin.Context) {

	roomId := c.Query("roomId")

	if roomId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "roomId was invalid",
			"status":  404,
		})
		return
	}

	// Upgrade the connection from http to websocket
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	ClientInstance := &Client{
		Conn:     ws,
		Send:     make(chan *Message),
		Room:     &Room{},
		Username: "",
		ID:       uuid.New(),
		Role:     "",
	}

	RoomManagerInstance.Mu.Lock()
	room, ok := RoomManagerInstance.Rooms[roomId]
	RoomManagerInstance.Mu.Unlock()

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server issue",
		})
		return
	}

	// Assigned a room to the client
	ClientInstance.Room = room
	// Added the client to the room manager along with the room.
	RoomManagerInstance.RegisterClient(roomId, ClientInstance)

	Wg.Add(1)
	go func(rm *RoomManager) {
		defer Wg.Done()
		defer func() {
			rm.UnregisterClient(roomId, ClientInstance)
			// close(ClientInstance.Send)
		}()

		for msg := range ClientInstance.Send {
			err := ws.WriteMessage(websocket.TextMessage, []byte(msg.MessageString))
			if err != nil {
				log.Println("write error:", err)
				return
			}
		}
	}(RoomManagerInstance)

	Wg.Add(1)
	go func(rm *RoomManager) {
		defer Wg.Done()
		for {
			_, p, err := ws.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				return
			}
			msg := &Message{
				MessageString: string(p),
				Sender:        "client",
			}
			// send to room for broadcast
			rm.BroadcastToRoom(roomId, msg)
		}
	}(RoomManagerInstance)

	Wg.Wait()
}
