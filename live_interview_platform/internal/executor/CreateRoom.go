package executor

import (
	"net/http"

	"github.com/bhuvneshuchiha/live_interview_platform/internal/websocket"
	"github.com/gin-gonic/gin"
)

// this will create a room for a given room manager
func CreateRoomHandler(c *gin.Context) {
	// pointer to the room manager defined in websocket.go
	rm := websocket.RoomManagerInstance
	roomId := rm.CreateRoom()

	if roomId == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create a room",
			"roomId":  roomId,
		})
		return
	}

	if roomId != "" {
		c.JSON(http.StatusCreated, gin.H{
			"message": "Room has been created successfully",
			"roomID":  roomId,
		})
		return
	}
}
