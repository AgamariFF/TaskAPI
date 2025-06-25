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
