package users

import (
	"cmaestro-api/internal/config"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, app *config.AppContext) {
	h := NewHandler(app)

	r.Get("/", h.List)
	r.Get("/{id}", h.Get)
}
