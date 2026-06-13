package users

import (
	"cmaestro-api/internal/config"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, cfg *config.Config) {
	h := NewHandler(cfg)

	r.Get("/", h.List)
	r.Get("/{id}", h.Get)
}
