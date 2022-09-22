package main

import (
	"encoding/json"
	"github.com/go-redis/redis"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"twitch_chat_analysis/util"
)

func main() {

	consumeMessages()

}

func consumeMessages() {
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

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	util.FailOnError(err, "Failed to register a consumer")

	client := redis.NewClient(util.RedisOptions())
	var forever chan struct{}

	go func() {
		for d := range msgs {
			var message util.Message
			json.Unmarshal(d.Body, &message)
			log.Printf("Received a message: %s", d.Body)

			err := client.LPush(message.Sender+"->"+message.Receiver, message.Message).Err()
			util.FailOnError(err, "redis lpush error")
			log.Printf("Pushed message to redis: %s", d.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}
