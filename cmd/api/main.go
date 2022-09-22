package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Message struct {
	Sender   string `json:"sender" binding:"required"`
	Receiver string `json:"receiver" binding:"required"`
	Message  string `json:"message" binding:"required"`
}

func main() {
	r := gin.Default()

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, "worked")
	})
	r.POST("/message", postMessageHandler)

	r.Run("localhost:8080")
}

func postMessageHandler(c *gin.Context) {
	var message Message

	// Call BindJSON to bind the received JSON to
	if err := c.ShouldBindJSON(&message); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	fmt.Println(message)
	// push to rabbitmq queue
	c.IndentedJSON(http.StatusOK, gin.H{"message": "got message success"})
}
