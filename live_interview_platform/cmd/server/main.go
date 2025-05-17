package main

import (
	"github.com/bhuvneshuchiha/live_interview_platform/internal/websocket"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H {
			"message" : "PONG",
		})
	})
	r.GET("/ws", websocket.HandleWebsocket)

	r.Run()
}
