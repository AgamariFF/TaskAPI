package main

import (
	"TaskAPI/config"
	"TaskAPI/handlers"
	"TaskAPI/logger"
	"TaskAPI/task"
	"fmt"
	"net/http"

	"TaskAPI/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	var env map[string]string
	env, err := config.LoadEnv()
	if err != nil {
		fmt.Println(err)
		env["PORT"] = "8080"
	}

	log, err := logger.NewLogger()
	if err != nil {
		fmt.Printf("Ошибка инициализации логгера: %v\n", err)
		return
	}
	defer log.Close()

	stor := task.NewStorage()

	log.Info("Запуск сервера на порте " + env["PORT"])

	r := mux.NewRouter()

	docs.SwaggerInfo.Title = "Task API"
	docs.SwaggerInfo.Description = "API для управления задачами"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	r.HandleFunc("/newtask", handlers.HandlerNewTask(stor, log)).Methods("POST")
	r.HandleFunc("/gettask", handlers.HandlerGetTask(stor, log)).Methods("GET")
	r.HandleFunc("/deletetask", handlers.HandlerDeleteTask(stor, log)).Methods("DELETE")

	if err := http.ListenAndServe(":"+env["PORT"], r); err != nil {
		log.Error(fmt.Sprintf("Ошибка при запуске сервера: %v", err))
	}
}
