package websocket

import (
	"log"
	"net/http"

	"github.com/bhuvneshuhciha/project_mordoria/internal/masterRoom"
	"github.com/bhuvneshuhciha/project_mordoria/internal/message"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var Upgrader = &websocket.Upgrader {
	ReadBufferSize : 1024,
	WriteBufferSize : 1024,
	CheckOrigin : func(r *http.Request) bool {
		return true
	},
}

func EstablishWebsocketConn(c *gin.Context, Upgrader *websocket.Upgrader) *websocket.Conn{
	ws, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"message": "Failed to connected to websocket",
		})
		return nil
	}
	return ws
}


func WebsocketReadMessage(ws *websocket.Conn, room *masterRoom.MasterRoom ) error {
	_, p, err := ws.ReadMessage()
	if err != nil {
		log.Println("Error occured", err)
		return err
	}
	msgStorage := &message.Message{
		MessageString: string(p),
		ClientEmoScore: 0,
	}
	masterRoom.Room.BroadCastMessage(msgStorage)
	log.Println("Message read from websocket.go", msgStorage)

	return nil
}


func WebSocketWriteMessage(ws *websocket.Conn, msg *message.Message) error {
	err := ws.WriteMessage(websocket.TextMessage, []byte(msg.MessageString))
	log.Println("Message written", msg.MessageString)
	if err != nil {
		log.Println("Error while writing to websocket", err)
		ws.Close()
		return err
	}
	return nil
}














