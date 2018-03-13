package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

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

type worker struct {
	DB string
}

func (w *worker) SendEmail(job *work.Job) error {
	id := job.ArgInt64("id")
	if err := job.ArgError(); err != nil {
		return err
	}

	fmt.Printf("Send email %d using db %s\n", id, w.DB)
	time.Sleep(5 * time.Second)
	fmt.Printf("Finished send email %d\n", id)
	return nil
}

func main() {
	wk := worker{}
	wk.DB = "mysql"

	pool := work.NewWorkerPool(wk, 5, "my_app", redisPool)
	pool.JobWithOptions("send_email", work.JobOptions{
		MaxFails: 3,
	}, wk.SendEmail)
	pool.Start()
	fmt.Println("Worker pool started")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	pool.Stop()
}
