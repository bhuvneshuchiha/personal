package ai_interceptor

import (
	"net/http"

	"github.com/bhuvneshuhciha/project_mordoria/internal/message"
	"github.com/gin-gonic/gin"
)

// store all the messages coming from the client to send to groq
var msgSlice []message.Message

func InterceptorHandler(c *gin.Context) {
	var msg message.Message
	err := c.BindJSON(&msg)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}
	msgSlice = append(msgSlice, msg)
	c.JSON(http.StatusOK, gin.H {
		"message": "success",
	})
	// clear the slice so that the next 30 second fresh messages
	// can be stored.
	msgSlice = msgSlice[:0]
	return
}










