package websocket

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

//Declare an upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//Define an endpoint handler for websocket connections
func HandleWebsocket(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	for {
		//Read the message
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(p))

		//Write the message
		er := ws.WriteMessage(messageType, p)
		if er != nil {
			log.Println(er)
			return
		}
	}
}
















