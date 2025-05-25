package handlers

import (
	"log"
	"sync"

	"github.com/bhuvneshuhciha/project_mordoria/internal/client"
	"github.com/bhuvneshuhciha/project_mordoria/internal/masterRoom"
	"github.com/bhuvneshuhciha/project_mordoria/internal/message"
	"github.com/bhuvneshuhciha/project_mordoria/internal/websocket"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var Wg sync.WaitGroup

func SocketHandler(c *gin.Context) {
	ws := websocket.EstablishWebsocketConn(c, websocket.Upgrader)
	roomId := masterRoom.Room.ID.String()

	go masterRoom.Room.RunLoop()

	clientInst := &client.Client{
		Conn:          ws,
		ID:            uuid.New(),
		SendMessage:   make(chan *message.Message),
		MessagesCount: 0,
	}

	err := masterRoom.Room.AddClient(clientInst, roomId)
	if err != nil {
		log.Println("Cannot add client to room", err)
		return
	}

	//Read Message
	Wg.Add(1)
	go func(clt *client.Client) {
		defer Wg.Done()
		log.Println("Inside Websocket read goroutine")
		for {
			err := websocket.WebsocketReadMessage(ws, masterRoom.Room)
			if err != nil {
				log.Println("Error while reading the message", err)
				masterRoom.Room.RemoveClient(clt.ID.String())
				return
			}
		}
	}(clientInst)

	//Write Message
	Wg.Add(1)
	go func(clt *client.Client) {
		defer Wg.Done()
		for msg := range clientInst.SendMessage {
			log.Println("Inside websocket write goroutine")
			err := websocket.WebSocketWriteMessage(ws, msg)
			if err != nil {
				log.Println("Could not write the message", err)
				masterRoom.Room.RemoveClient(clt.ID.String())
				return
			}
		}
	}(clientInst)

	Wg.Wait()
}
