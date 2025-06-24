package task

import (
	"context"
	"time"
)

type Task struct {
	ID         string
	Status     string
	CreatedAt  time.Time
	FinishedAt time.Time
	cancelFunc context.CancelFunc
}

type ResponseTask struct {
	ID           string        `json:"id"`
	Status       string        `json:"status"`
	CreatedAt    time.Time     `json:"created_at"`
	TimeDuration time.Duration `json:"time_duration"`
}

func TaskToResponseTask(task Task) ResponseTask {
	var duration time.Duration

	switch task.Status {
	case "completed", "canceled":
		duration = task.FinishedAt.Sub(task.CreatedAt)
	case "in_progress":
		duration = time.Since(task.CreatedAt)
	default: // "pending" и другие статусы
		duration = 0
	}

	return ResponseTask{
		ID:           task.ID,
		Status:       task.Status,
		CreatedAt:    task.CreatedAt,
		TimeDuration: time.Duration(duration.Seconds()),
	}
}
