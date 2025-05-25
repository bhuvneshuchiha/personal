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
	room := masterRoom.CreateMasterRoom()
	roomId := room.ID.String()

	go room.RunLoop()

	clientInst := &client.Client{
		ID:            uuid.New(),
		SendMessage:   make(chan *message.Message),
		MessagesCount: 0,
	}

	err := room.AddClient(clientInst, roomId)
	if err != nil {
		log.Println("Cannot add client to room", err)
		return
	}

	//Read Message
	Wg.Add(1)
	go func(clt *client.Client) {
		defer Wg.Done()
		err := websocket.WebsocketReadMessage(ws, room)
		if err != nil {
			log.Println("Error while reading the message", err)
			room.RemoveClient(clt.ID.String())
			return
		}
	}(clientInst)

	//Write Message
	Wg.Add(1)
	go func(clt *client.Client) {
		defer Wg.Done()
		for msg := range clientInst.SendMessage {
			err := websocket.WebSocketWriteMessage(ws, msg)
			if err != nil {
				log.Println("Could not write the message", err)
				room.RemoveClient(clt.ID.String())
				return
			}
		}
	}(clientInst)

	Wg.Wait()
}
