package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Get("/users", getUsers)

	http.ListenAndServe(":8080", r)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	users := []string{"Alice", "Bob"}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
