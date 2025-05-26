package ai_interceptor

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type Ai_message struct {
	Client_id string
	Message   string
}

// store all the messages coming from the client to send to groq
var msgSlice []string

func InterceptorHandler(c *gin.Context) {
	var msg Ai_message
	err := c.BindJSON(&msg)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}
	msgSlice = append(msgSlice, msg.Message)
	c.JSON(http.StatusOK, gin.H {
		"message": "success",
	})
	return
}








