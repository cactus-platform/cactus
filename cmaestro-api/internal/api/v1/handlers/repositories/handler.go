package repositories

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router) {
	h := NewHandler()

	r.Get("/", h.List)
	r.Post("/", h.Create)
}
