package main

import (
	"cmaestro-api/internal/router"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := router.New()

	// Middlewares Configuration
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/users", getUsers)
	// system
	// health
	// status

	r.ListenAndServe(":8080")
}

func getUsers(w http.ResponseWriter, _ *http.Request) {
	users := []string{"admin"}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
