package websocket

import (
	"context"
	"log"
	"net/http"
	"sync"

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

var RoomManagerInstance = &RoomManager{
	Rooms: make(map[string]*Room),
	Mu:    sync.Mutex{},
}

var Wg = &sync.WaitGroup{}

// websocket handler
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
		Role:     "client",
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//Write pump
	Wg.Add(1)
	go func(rm *RoomManager, ctx context.Context, cancel context.CancelFunc) {
		defer Wg.Done()
		defer func() {
			rm.UnregisterClient(roomId, ClientInstance)
			close(ClientInstance.Send)
			log.Println("Client unregistered successfully...")
		}()
		for {
			select {
			case <-ctx.Done():
				log.Println("context cancelled, write pump exiting early")
				return
			case msg, ok := <-ClientInstance.Send:
				if !ok {
					log.Println("client send channel closed")
					return
				}
				err := ws.WriteMessage(websocket.TextMessage, []byte(msg.MessageString))
				if err != nil {
					cancel()
					log.Println("write error:", err)
					return
				}
			}
		}
	}(RoomManagerInstance, ctx, cancel)

	//Read Pump
	Wg.Add(1)
	go func(rm *RoomManager, ctx context.Context, cancel context.CancelFunc) {
		defer Wg.Done()

		for {
			select {
			case <-ctx.Done():
				log.Println("context cancelled, exiting loop early")
				return
			default:
				_, p, err := ws.ReadMessage()
				if err != nil {
					cancel()
					log.Println("read error:", err)
					return
				}
				msg := &Message{
					MessageString: string(p),
					Sender:        "client",
					Status:        "active",
				}
				// send to room for broadcast
				rm.BroadcastToRoom(roomId, msg)
				log.Println(msg)
			}
		}
	}(RoomManagerInstance, ctx, cancel)

	Wg.Wait()
}
