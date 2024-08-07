package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"

	"orchestrator-exp/task"
	"orchestrator-exp/worker"
	"orchestrator-exp/manager"
)

func main() {
	whost := os.Getenv("WORKER_HOST")
    wport, _ := strconv.Atoi(os.Getenv("WORKER_PORT"))
 
    mhost := os.Getenv("MANAGER_HOST")
    mport, _ := strconv.Atoi(os.Getenv("MANAGER_PORT"))

	fmt.Println("Starting worker")

	w := worker.Worker{
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}

	wapi := worker.Api{Address: whost, Port: wport, Worker: &w}

	go w.RunTasks()
	go w.CollectStats()
	go w.UpdateTasks()
	go wapi.Start()

	fmt.Println("Starting manager")

	workers := []string{fmt.Sprintf("%s:%d", whost, wport)}
	m := manager.New(workers)
	mapi := manager.Api{Address: mhost, Port: mport, Manager: m}

	go m.ProcessTasks()
	go m.UpdateTasks()
	go m.DoHealthChecks()

	mapi.Start()
}