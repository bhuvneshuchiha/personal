package main

import (
	"net/http"

	"github.com/bhuvneshuchiha/live_interview_platform/internal/executor"
	"github.com/bhuvneshuchiha/live_interview_platform/internal/websocket"

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
	r.POST("/create-room", executor.CreateRoomHandler)

	r.Run()
}
