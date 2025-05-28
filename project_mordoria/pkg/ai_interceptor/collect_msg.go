package ai_interceptor

import (
	"errors"
	"log"
	"net/http"

	finalMessage "github.com/bhuvneshuhciha/project_mordoria/pkg/final_message"
	"github.com/bhuvneshuhciha/project_mordoria/pkg/groq_controllers"
	"github.com/gin-gonic/gin"
)


func InterceptorHandler(c *gin.Context) {
	log.Println("Inside the request ----------")

	err := c.BindJSON(&finalMessage.MsgBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		log.Println("Here is the msg body ------", finalMessage.MsgBody)
		log.Println("Here is the error", err.Error())
		return
	}

	log.Println("this is the payload from axios:", finalMessage.MsgBody)

	er := groq_controllers.PrepareData(finalMessage.MsgBody, groq_controllers.DataStruct)
	if er != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"message" : er,
		})
		return
	}
	respString, e := groq_controllers.SendDataToGroq()
	if e != nil {
		log.Println("Got an error while requesting from groq, file: collect_msg.go")
		return
	}

	c.JSON(http.StatusOK, gin.H {
		"message": "success",
		"groq_response": respString,
	})

	return
}
