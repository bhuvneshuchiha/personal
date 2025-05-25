package main

import (
	"net/http"

	"github.com/bhuvneshuhciha/project_mordoria/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H {
		"message" : "success",
		"status" : 201,
		})
	})

	r.POST("/v1/mordoria", handlers.SocketHandler)

	r.Run()
}
