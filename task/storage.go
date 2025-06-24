package task

import (
	"sync"
)

type TaskStorage struct {
	sync.RWMutex
	tasks map[string]Task
}

func NewStorage() *TaskStorage {
	return &TaskStorage{
		tasks: make(map[string]Task),
	}
}

func (s *TaskStorage) AddTask(task Task) {
	s.Lock()
	defer s.Unlock()
	s.tasks[task.ID] = task
}

func (s *TaskStorage) GetTask(id string) (*Task, bool) {
	s.RLock()
	defer s.RUnlock()
	task, ok := s.tasks[id]
	return &task, ok
}

func (s *TaskStorage) DeleteTask(id string) bool {
	s.Lock()
	defer s.Unlock()

	if _, exists := s.tasks[id]; !exists {
		return false
	}

	delete(s.tasks, id)
	return true
}
