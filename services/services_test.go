package services_test

import (
	"testing"

	"TaskAPI/logger"
	"TaskAPI/services"
	"TaskAPI/task"

	"github.com/stretchr/testify/assert"
)

func createTestStorage() *task.TaskStorage {
	stor := task.NewStorage()
	return stor
}

func createTestLogger() *logger.Logger {
	log, _ := logger.NewLogger()
	return log
}

func TestAdd(t *testing.T) {
	stor := createTestStorage()
	log := createTestLogger()

	id, err := services.Add(stor, log)

	assert.NoError(t, err)
	assert.NotEmpty(t, id)

	taskFind, exist := stor.GetTask(id)
	assert.True(t, exist)
	assert.Equal(t, id, taskFind.ID)
}

func TestGet(t *testing.T) {
	stor := createTestStorage()
	log := createTestLogger()

	expectedTask := task.Task{
		ID:     "12345",
		Status: "running",
	}
	stor.AddTask(expectedTask)

	response, err := services.Get(stor, log, "12345")

	assert.NoError(t, err)
	assert.Equal(t, "12345", response.ID)
	assert.Equal(t, "running", response.Status)

	response, err = services.Get(stor, log, "unreal_id")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "task not found")
	assert.Equal(t, task.Task{}, response)
}

func TestDelete_Success_WithoutMocks(t *testing.T) {
	stor := createTestStorage()
	log := createTestLogger()

	stor.AddTask(task.Task{ID: "12345", Status: "running"})

	err := services.Delete(stor, log, "12345")

	assert.NoError(t, err)

	_, exist := stor.GetTask("12345")
	assert.False(t, exist)

	err = services.Delete(stor, log, "nonexistent_id")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Task not found")
}
