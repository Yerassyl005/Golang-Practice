package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"practice5/repository"
	"practice5/handlers"

	_ "github.com/lib/pq"
)

func main() {

	conn := "user=postgres password=erasyl25 dbname=practice5 sslmode=disable"

	db,_ := sql.Open("postgres",conn)

	repo := repository.Repository{DB: db}

	h := handlers.Handler{Repo: &repo}

	http.HandleFunc("/users",h.GetUsers)
	http.HandleFunc("/common",h.GetCommonFriends)

	fmt.Println("Server started on port 8080")

	http.ListenAndServe(":8080",nil)
}
