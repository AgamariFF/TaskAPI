package handlers

import (
	"TaskAPI/logger"
	"TaskAPI/services"
	"TaskAPI/task"
	"encoding/json"
	"fmt"
	"net/http"
)

func HandlerNewTask(stor *task.TaskStorage, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		id, err := services.Add(stor, log)
		if err != nil {
			log.Error("Ошибка добавления задания. " + err.Error())
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(id))
	}
}

func HandlerGetTask(stor *task.TaskStorage, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		queryParams := r.URL.Query()

		id := queryParams.Get("id")

		taskResponse, err := services.Get(stor, log, id)
		if err != nil {
			log.Error(fmt.Sprintf("Ошибка при получении задания: %v", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(taskResponse); err != nil {
			log.Error(fmt.Sprintf("Ошибка при записи ответа: %v", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

func HandlerDeleteTask(stor *task.TaskStorage, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, `{"error": "Task ID is required"}`, http.StatusBadRequest)
			return
		}

		if err := services.Delete(stor, log, id); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "deleted",
			"task_id": id,
		})
	}
}
