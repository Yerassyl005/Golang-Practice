package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var db *sql.DB

func main() {

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var err error

	for {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
		}

		if err == nil {
			break
		}

		fmt.Println("Waiting for database...")
		time.Sleep(2 * time.Second)
	}

	fmt.Println("Database connected!")
	fmt.Println("Starting the Server...")

	http.HandleFunc("/items", getItems)
	http.HandleFunc("/create", createItem)
	http.HandleFunc("/delete", deleteItem)
	http.HandleFunc("/update", updateItem)

	server := &http.Server{
		Addr: ":8080",
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	fmt.Println("Shutting down gracefully...")

	db.Close()
	server.Shutdown(context.Background())
}

func getItems(w http.ResponseWriter, r *http.Request) {

	rows, _ := db.Query("SELECT id, name FROM items")

	var items []Item

	for rows.Next() {
		var item Item
		rows.Scan(&item.ID, &item.Name)
		items = append(items, item)
	}

	json.NewEncoder(w).Encode(items)
}

func createItem(w http.ResponseWriter, r *http.Request) {

	var item Item
	json.NewDecoder(r.Body).Decode(&item)

	db.Exec("INSERT INTO items(name) VALUES($1)", item.Name)

	w.Write([]byte("Item created"))
}

func deleteItem(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	db.Exec("DELETE FROM items WHERE id=$1", id)

	w.Write([]byte("Item deleted"))
}

func updateItem(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	var item Item
	json.NewDecoder(r.Body).Decode(&item)

	db.Exec("UPDATE items SET name=$1 WHERE id=$2", item.Name, id)

	w.Write([]byte("Item updated"))
}
