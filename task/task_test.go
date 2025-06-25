package task_test

import (
	"TaskAPI/task"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createTestStorage() *task.TaskStorage {
	return task.NewStorage()
}

func TestCreateTask(t *testing.T) {
	storage := createTestStorage()

	testTask := task.CreateTask(storage)

	assert.NotEmpty(t, testTask.ID)
	assert.Equal(t, "pending", testTask.Status)
	assert.NotNil(t, testTask.CreatedAt)

	storedTask, exists := storage.GetTask(testTask.ID)
	assert.True(t, exists)
	assert.Equal(t, testTask.ID, storedTask.ID)
	assert.Equal(t, "pending", storedTask.Status)
}

// func TestExecuteTask(t *testing.T) {
// 	storage := createTestStorage()

// 	testTask := task.Task{
// 		ID:        "test-task",
// 		Status:    "pending",
// 		CreatedAt: time.Now(),
// 	}
// 	storage.AddTask(testTask)

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()
// 	go task.ExecuteTask(ctx, storage, testTask.ID)

// 	time.Sleep(5 * time.Minute)

// 	storedTask, exists := storage.GetTask(testTask.ID)
// 	assert.True(t, exists)
// 	assert.Equal(t, "completed", storedTask.Status)
// 	assert.NotNil(t, storedTask.FinishedAt)
// }

func TestExecuteTask_Cancellation(t *testing.T) {
	storage := createTestStorage()

	testTask := task.Task{
		ID:        "test-task",
		Status:    "pending",
		CreatedAt: time.Now(),
	}
	storage.AddTask(testTask)

	ctx, cancel := context.WithCancel(context.Background())
	go task.ExecuteTask(ctx, storage, testTask.ID)

	time.Sleep(1 * time.Second)
	cancel()

	time.Sleep(2 * time.Second)

	// Assert
	storedTask, exists := storage.GetTask(testTask.ID)
	assert.True(t, exists)
	assert.Equal(t, "canceled", storedTask.Status)
	assert.NotNil(t, storedTask.FinishedAt)
}

func TestCancelTask_Success(t *testing.T) {
	storage := createTestStorage()

	testTask := task.CreateTask(storage)
	time.Sleep(time.Second)

	success := task.CancelTask(storage, testTask.ID)
	time.Sleep(time.Second)

	assert.True(t, success)

	storedTask, exists := storage.GetTask(testTask.ID)
	assert.True(t, exists)
	assert.Equal(t, "canceled", storedTask.Status)
	assert.NotNil(t, storedTask.FinishedAt)
}

func TestCancelTask_Fail(t *testing.T) {
	storage := createTestStorage()

	testTask := task.Task{
		ID:        "test-task",
		Status:    "pending",
		CreatedAt: time.Now(),
	}
	storage.AddTask(testTask)

	success := task.CancelTask(storage, testTask.ID)

	assert.False(t, success)

	storedTask, exists := storage.GetTask(testTask.ID)
	assert.True(t, exists)
	assert.Equal(t, "pending", storedTask.Status)
}
