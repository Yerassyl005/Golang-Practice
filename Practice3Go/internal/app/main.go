package app

import (
	"context"
	"fmt"
	"time"

	"github.com/Yerassyl005/go-practice3/internal/repository"
	"github.com/Yerassyl005/go-practice3/internal/repository/_postgres"
	"github.com/Yerassyl005/go-practice3/pkg/modules"
)

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConfig := initPostgreConfig()

	postgres := _postgres.NewPGXDialect(ctx, dbConfig)

	repos := repository.NewRepositories(postgres)

	users, err := repos.GetUsers()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Users from DB:", users)
}

func initPostgreConfig() *modules.PostgreConfig {
	return &modules.PostgreConfig{
		Host:        "localhost",
		Port:        "5432",
		Username:    "postgres",
		Password:    "erasyl25",
		DBName:      "mydb",
		SSLMode:     "disable",
		ExecTimeout: 5 * time.Second,
	}
}