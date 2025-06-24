package task

import (
	"context"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

func StartTask(storage *TaskStorage) Task {
	ctx, cancel := context.WithCancel(context.Background())
	var task Task

	for {
		id := uuid.New().String()
		if _, exists := storage.GetTask(id); !exists {
			task = Task{
				ID:         id,
				Status:     "pending",
				CreatedAt:  time.Now(),
				cancelFunc: cancel,
			}
			break
		}
	}
	storage.AddTask(task)
	go executeTask(ctx, storage, task.ID)

	return task
}

func executeTask(ctx context.Context, storage *TaskStorage, id string) {
	if task, ok := storage.GetTask(id); ok {
		task.Status = "in_progress"
		storage.AddTask(*task)
	}

	select {
	case <-time.After(time.Duration(rand.Intn(120)+180) * time.Second):
		if task, ok := storage.GetTask(id); ok {
			task.Status = "completed"
			task.FinishedAt = time.Now()
			storage.AddTask(*task)
		}
	case <-ctx.Done():
		if task, ok := storage.GetTask(id); ok {
			task.Status = "canceled"
			task.FinishedAt = time.Now()
			storage.AddTask(*task)
		}
	}
}

func CancelTask(storage *TaskStorage, id string) bool {
	if task, ok := storage.GetTask(id); ok && task.Status == "in_progress" {
		task.cancelFunc()
		return true
	}
	return false
}
