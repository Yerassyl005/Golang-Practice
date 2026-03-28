package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Yerassyl005/go-practice3/internal/handlers"
	"github.com/Yerassyl005/go-practice3/internal/middleware"
	"github.com/Yerassyl005/go-practice3/internal/repository"
	"github.com/Yerassyl005/go-practice3/internal/repository/_postgres"
	"github.com/Yerassyl005/go-practice3/internal/usecase"
	"github.com/Yerassyl005/go-practice3/pkg/modules"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConfig := &modules.PostgreConfig{
		Host:        "localhost",
		Port:        "5432",
		Username:    "postgres",
		Password:    "erasyl25",
		DBName:      "mydb",
		SSLMode:     "disable",
		ExecTimeout: 5 * time.Second,
	}

	postgres := _postgres.NewPGXDialect(ctx, dbConfig)
	repos := repository.NewRepositories(postgres)
	usecase := usecase.NewUserUsecase(repos.UserRepository)
	userHandler := handler.NewUserHandler(usecase)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/users", userHandler.Users)
	mux.HandleFunc("/users/", userHandler.UserByID)

	finalHandler := middleware.Logging(
		middleware.APIKey(mux),
	)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", finalHandler))
}