package main

import (
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/gocraft/work"
)

var redisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", ":6379")
	},
}

var enqueuer = work.NewEnqueuer("my_app", redisPool)

func main() {
	// Enqueue a job named "send_email" with the specified parameters.
	for i := 0; i < 10; i++ {
		_, err := enqueuer.Enqueue("send_email", work.Q{"id": i})
		if err != nil {
			log.Fatal(err)
		}
	}
}
