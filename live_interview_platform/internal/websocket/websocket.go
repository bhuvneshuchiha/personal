package websocket

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Declare an upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

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

	//Write pump
	Wg.Add(1)
	go func(rm *RoomManager) {
		ticker := time.NewTicker(pingPeriod)
		defer Wg.Done()
		defer func() {
			rm.UnregisterClient(roomId, ClientInstance)
			close(ClientInstance.Send)
			ticker.Stop()
			ClientInstance.Conn.Close()
			log.Println("Client unregistered successfully...")
		}()

		for msg := range ClientInstance.Send {
			err := ClientInstance.Conn.WriteMessage(websocket.TextMessage, []byte(msg.MessageString))
			if err != nil {
				log.Println("write error:", err)
				return
			}
		}
	}(RoomManagerInstance)

	//Read Pump
	Wg.Add(1)
	go func(rm *RoomManager) {
		defer Wg.Done()

		ClientInstance.Conn.SetReadLimit(maxMessageSize)
		ClientInstance.Conn.SetReadDeadline(time.Now().Add(pongWait))
		ClientInstance.Conn.SetPongHandler(func(string) error {
			ClientInstance.Conn.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})

		for {
			_, p, err := ClientInstance.Conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway,
					websocket.CloseAbnormalClosure) {
					log.Println("read error:", err)
				}
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
