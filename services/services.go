package services

import (
	"TaskAPI/logger"
	"TaskAPI/task"
	"fmt"
)

func Add(stor *task.TaskStorage, log *logger.Logger) (string, error) {
	id := task.StartTask(stor).ID

	return id, nil
}

func Get(stor *task.TaskStorage, log *logger.Logger, id string) (task.ResponseTask, error) {
	taskFind, exist := stor.GetTask(id)
	if !exist {
		err := fmt.Errorf("task not found with id: %s", id)
		log.Error("Failed to get task. " + err.Error())
		return task.ResponseTask{}, err
	}
	taskResponse := task.TaskToResponseTask(*taskFind)

	return taskResponse, nil
}

func Delete(stor *task.TaskStorage, log *logger.Logger, id string) error {
	task.CancelTask(stor, id)

	if removed := stor.DeleteTask(id); !removed {
		return fmt.Errorf("Task not found")
	}
	return nil
}
