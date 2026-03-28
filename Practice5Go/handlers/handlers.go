package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"practice5/repository"
)

type Handler struct {
	Repo *repository.Repository
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {

	page,_ := strconv.Atoi(r.URL.Query().Get("page"))
	size,_ := strconv.Atoi(r.URL.Query().Get("size"))

	name := r.URL.Query().Get("name")
	email := r.URL.Query().Get("email")
	order := r.URL.Query().Get("order")

	if page == 0 {
		page = 1
	}

	if size == 0 {
		size = 5
	}

	data,_ := h.Repo.GetPaginatedUsers(page,size,name,email,order)

	json.NewEncoder(w).Encode(data)
}

func (h *Handler) GetCommonFriends(w http.ResponseWriter, r *http.Request) {

	u1,_ := strconv.Atoi(r.URL.Query().Get("user1"))
	u2,_ := strconv.Atoi(r.URL.Query().Get("user2"))

	data,_ := h.Repo.GetCommonFriends(u1,u2)

	json.NewEncoder(w).Encode(data)
}
