package main

import (
	"fmt"
	"time"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
)

func intensiveAdd(id, num1, num2 int) (int, error) {
	fmt.Printf("Start working on job %d\n", id)
	time.Sleep(5 * time.Second)
	ret := num1 + num2
	fmt.Printf("Finished working on job %d\n", id)
	return ret, nil
}

func main() {
	server, err := machinery.NewServer(&config.Config{
		Broker:        "redis://localhost:6379/0",
		DefaultQueue:  "my_queue",
		ResultBackend: "redis://localhost:6379/1",
	})
	if err != nil {
		panic("Cannot start server")
	}

	if err := server.RegisterTask("add", intensiveAdd); err != nil {
		panic("Cannot register task")
	}

	worker := server.NewWorker("worker_1", 5)
	worker.Launch()
}
