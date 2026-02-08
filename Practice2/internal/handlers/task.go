package handlers

import (
	"practice2go/models"
)

import (
	"encoding/json"
	"net/http"
	"strconv"
)

var tasks = make(map[int]models.Task)
var nextID = 1

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	case http.MethodGet:
		idStr := r.URL.Query().Get("id")

		if idStr == "" {
			list := []models.Task{}
			for _, t := range tasks {
				list = append(list, t)
			}
			json.NewEncoder(w).Encode(list)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid id"})
			return
		}

		task, ok := tasks[id]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "task not found"})
			return
		}

		json.NewEncoder(w).Encode(task)

	case http.MethodPost:
		var input struct {
			Title string `json:"title"`
		}

		json.NewDecoder(r.Body).Decode(&input)

		if input.Title == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid title"})
			return
		}

		task := models.Task{
			ID:    nextID,
			Title: input.Title,
			Done:  false,
		}

		tasks[nextID] = task
		nextID++

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(task)

	case http.MethodPatch:
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid id"})
			return
		}

		task, ok := tasks[id]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "task not found"})
			return
		}

		var input struct {
			Done bool `json:"done"`
		}

		json.NewDecoder(r.Body).Decode(&input)
		task.Done = input.Done
		tasks[id] = task

		json.NewEncoder(w).Encode(map[string]bool{"updated": true})
	}
}
