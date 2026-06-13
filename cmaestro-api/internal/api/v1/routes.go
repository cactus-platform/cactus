package v1

import (
	"cmaestro-api/internal/api/v1/handlers/repositories"
	"cmaestro-api/internal/api/v1/handlers/users"
	"cmaestro-api/internal/config"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes registers v1 routes onto the provided router.
func RegisterRoutes(r chi.Router, cfg *config.Config) {
	r.Route("/users", func(r chi.Router) {
		users.RegisterRoutes(r, cfg)
	})

	r.Route("/repositories", func(r chi.Router) {
		repositories.RegisterRoutes(r, cfg)
	})
}
