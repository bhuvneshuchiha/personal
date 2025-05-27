package ai_interceptor

import (
	"log"
	"net/http"

	"github.com/bhuvneshuhciha/project_mordoria/internal/message"
	"github.com/gin-gonic/gin"
)

// store all the messages coming from the client to send to groq
type FinalPayload struct {
	ClientId      string            `json:"client_id"`
	Payload       []message.Message `json:"payload"`
	Ai_emot_score string            `json:"ai_emot_score"`
}

func InterceptorHandler(c *gin.Context) {
	log.Println("Inside the request ----------")
	var msgBody *FinalPayload = &FinalPayload{}

	err := c.BindJSON(&msgBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		log.Println("Here is the msg body ------", msgBody)
		log.Println("Here is the error", err.Error())
		return
	}

	log.Println("this is the payload from axios:", msgBody)

	c.JSON(http.StatusOK, gin.H{
		"message":      "success",
		"all_messages": msgBody,
	})

	return
}
