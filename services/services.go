package services

import (
	"TaskAPI/logger"
	"TaskAPI/task"
	"fmt"
)

func Add(stor *task.TaskStorage, log *logger.Logger) (string, error) {
	id := task.CreateTask(stor).ID

	return id, nil
}

func Get(stor *task.TaskStorage, log *logger.Logger, id string) (task.Task, error) {
	taskFind, exist := stor.GetTask(id)
	if !exist {
		err := fmt.Errorf("task not found with id: %s", id)
		log.Error("Ошибка получения задания. " + err.Error())
		return task.Task{}, err
	}

	return *taskFind, nil
}

func Delete(stor *task.TaskStorage, log *logger.Logger, id string) error {
	task.CancelTask(stor, id)

	if removed := stor.DeleteTask(id); !removed {
		return fmt.Errorf("Task not found")
	}
	return nil
}
