package main

import (
	"fmt"
	"sync"
	"time"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
)

func sendTask(server *machinery.Server, id int) {
	task := tasks.Signature{
		Name: "add",
		Args: []tasks.Arg{
			{
				Type:  "int",
				Value: id,
			},
			{
				Type:  "int",
				Value: 4,
			},
			{
				Type:  "int",
				Value: 3,
			},
		},
	}

	job, err := server.SendTask(&task)
	if err != nil {
		fmt.Printf("Error when sending task %v\n", err)
	}
	fmt.Printf("Job %d sent\n", id)

	maxWaitTime := time.Duration(2 * time.Second)
	result, err := job.Get(maxWaitTime)
	if err != nil {
		fmt.Printf("Job %d return error: %v\n", id, err)
	}
	fmt.Printf("Job %d return sucess: %v\n", id, tasks.HumanReadableResults(result))
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

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sendTask(server, id)
		}(i)
	}

	wg.Wait()
}
