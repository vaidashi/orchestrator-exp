package worker

import (
	"orchestrator-exp/task"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (a *Api) StartTaskHandler(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	te := task.TaskEvent{}
	err := d.Decode(&te)

	if err != nil {
		msg := fmt.Sprintf("Error decoding request body: %v\n", err)
		log.Print(msg)
		w.WriteHeader(400)

		e := ErrResponse {
			Message: msg,
			HTTPStatusCode: 400,
		}

		json.NewEncoder(w).Encode(e)
		return
	}

	a.Worker.AddTask(te.Task)
	log.Printf("Task %v added to queue\n", te.Task.ID)
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(te.Task)
}

func (a *Api) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(a.Worker.GetTasks())
}

func (a *Api) StopTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")

	if taskID == "" {
		log.Printf("Task ID is required\n")
		w.WriteHeader(400)
	}

	tID, _ := uuid.Parse(taskID)
	_, ok := a.Worker.Db[tID]

	if !ok {
		log.Printf("Task %v not found\n", tID)
		w.WriteHeader(404)
	}

	taskToStop := a.Worker.Db[tID]
	taskCopy := *taskToStop // not to modify the original task in the store
	taskCopy.State = task.Completed
	a.Worker.AddTask(taskCopy)

	log.Printf("Added task %v to stop container %v\n", taskToStop.ID, taskToStop.ContainerID)
	w.WriteHeader(204)
}