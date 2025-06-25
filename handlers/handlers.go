package handlers

import (
	"TaskAPI/logger"
	"TaskAPI/services"
	"TaskAPI/task"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// ResponseTask представляет информацию о задаче
// swagger:model ResponseTask
type ResponseTask struct {
	// example: a80f84a8-4841-46a8-bb3f-10d9f775c27a
	ID string `json:"id"`

	// example: in_progress
	Status string `json:"status"`

	// example: 2025-06-25T10:00:00Z
	CreatedAt time.Time `json:"created_at"`

	// example: 1h30m
	TimeDuration string `json:"time_duration"`
}

func TaskToResponseTask(task task.Task) ResponseTask {
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
		TimeDuration: time.Duration(duration.Seconds()).String(),
	}
}

// newtask создает и запускает в работу задачу, возвращает её id
// @Summary Создание новой задачи
// @Description Создает новую задачу и возвращает её ID
// @Tags tasks
// @Accept json
// @Produce json
// @Success 202 {string} string "ID созданной задачи"
// @Router /newtask [post]
func HandlerNewTask(stor *task.TaskStorage, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Получен POST-запрос на создание новой задачи")

		if r.Method != http.MethodPost {
			log.Info("Некорректный HTTP-метод: " + r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		id, err := services.Add(stor, log)
		if err != nil {
			log.Error("Ошибка добавления задания. " + err.Error())
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
			return
		}

		log.Info(fmt.Sprintf("Задача успешно добавлена. ID: %s", id))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"id": "` + id + `"}`))
	}
}

// gettask возвращает информацию о задаче по её ID
// @Summary Получение задачи по ID
// @Description Возвращает информацию о задаче по её ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id query string true "ID задачи"
// @Success 200 {object} ResponseTask "Информация о задаче"
// @Router /gettask [get]
func HandlerGetTask(stor *task.TaskStorage, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Получен GET-запрос на получение задачи")

		if r.Method != http.MethodGet {
			log.Info("Некорректный HTTP-метод: " + r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		queryParams := r.URL.Query()
		id := queryParams.Get("id")
		log.Info(fmt.Sprintf("GET-запрос имеет ID= %s", id))

		taskFind, err := services.Get(stor, log, id)
		if err != nil {
			log.Error(fmt.Sprintf("Ошибка при получении задания: %v", err))
			if strings.Contains(err.Error(), "task not found with id") {
				http.Error(w, "Incorrect ID", http.StatusBadRequest)
			} else {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		taskResponse := TaskToResponseTask(taskFind)

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(taskResponse); err != nil {
			log.Error(fmt.Sprintf("Ошибка при записи ответа: %v", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

// deletetask удаляет задачу по её ID
// @Summary Удаление задачи по ID
// @Description Удаляет задачу по ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id query string true "ID задачи"
// @Success 200 {object} map[string]string "Задача удалена"
// @Router /deletetask [delete]
func HandlerDeleteTask(stor *task.TaskStorage, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Получен DELETE-запрос на удаление задачи")

		if r.Method != http.MethodDelete {
			log.Info("Некорректный HTTP-метод: " + r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		id := r.URL.Query().Get("id")
		if id == "" {
			log.Error("Не указан ID в query параметрах")
			http.Error(w, `{"error": "Task ID is required"}`, http.StatusBadRequest)
			return
		}

		if err := services.Delete(stor, log, id); err != nil {
			if strings.Contains(err.Error(), "Task not found") {
				log.Error("Задача не найдена с ID=" + id)
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				log.Error("Ошибка удаления задачи. " + err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(map[string]string{
			"status":  "deleted",
			"task_id": id,
		})
		if err != nil {
			log.Error(fmt.Sprintf("Ошибка при записи ответа: %v", err))
		}
	}
}
