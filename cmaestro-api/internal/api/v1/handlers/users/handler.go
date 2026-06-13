package users

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router) {
	h := NewHandler()

	r.Get("/", h.List)
	r.Get("/{id}", h.Get)
}
