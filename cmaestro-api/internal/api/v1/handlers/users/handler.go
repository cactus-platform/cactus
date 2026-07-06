package users

import (
	"cmaestro-api/internal/bootstrap"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, app *bootstrap.App) {
	h := NewHandler(app)

	r.Get("/", h.List)
	r.Get("/{id}", h.Get)
}
