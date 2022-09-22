package util

import (
	"github.com/go-redis/redis"
	"log"
	"time"
)

type Message struct {
	Sender   string `json:"sender" binding:"required"`
	Receiver string `json:"receiver" binding:"required"`
	Message  string `json:"message" binding:"required"`
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// Default REDIS options
func RedisOptions() *redis.Options {
	return &redis.Options{
		Addr:               "localhost:6379",
		Password:           "",
		DB:                 0,
		DialTimeout:        10 * time.Second,
		ReadTimeout:        30 * time.Second,
		WriteTimeout:       30 * time.Second,
		PoolSize:           10,
		PoolTimeout:        30 * time.Second,
		IdleTimeout:        500 * time.Millisecond,
		IdleCheckFrequency: 500 * time.Millisecond,
	}
}
