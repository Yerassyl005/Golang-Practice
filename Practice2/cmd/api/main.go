package main

import (
	"log"
	"net/http"

	"practice2go/internal/handlers"
	"practice2go/middleware"
)

func main() {
	// 1. Создаём роутер
	mux := http.NewServeMux()

	// 2. Регистрируем маршрут /tasks
	mux.Handle(
		"/tasks",
		middleware.Logger(middleware.Auth(http.HandlerFunc(handlers.TasksHandler))),
	)

	// 3. Запускаем сервер
	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
