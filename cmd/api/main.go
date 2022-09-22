package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"time"
	"twitch_chat_analysis/util"
)

func main() {
	r := gin.Default()

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, "worked")
	})
	r.POST("/message", postMessageHandler)

	r.GET("/message/list", reportIssue)
	r.Run("localhost:8080")

}

func reportIssue(c *gin.Context) {
	sender := c.Request.URL.Query().Get("sender")
	receiver := c.Request.URL.Query().Get("receiver")
	log.Printf(sender + " " + receiver)
	client := redis.NewClient(util.RedisOptions())
	result, err := client.LRange(sender+"->"+receiver, 0, 1).Result()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "redis errr")
	}
	log.Println(result)
	c.IndentedJSON(http.StatusOK, gin.H{"result": result})

}

func postMessageHandler(c *gin.Context) {
	var message util.Message

	// Call BindJSON to bind the received JSON to
	if err := c.ShouldBindJSON(&message); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	fmt.Println(message)
	sendMessageToQueue(message)
	// push to rabbitmq queue
	c.IndentedJSON(http.StatusOK, gin.H{"message": "got message success"})
}

func sendMessageToQueue(message util.Message) {
	conn, err := amqp.Dial("amqp://user:password@localhost:7001/")
	util.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	util.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"testQueue", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	util.FailOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, _ := json.Marshal(message)
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	util.FailOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)

}
