package v1

import (
	"cmaestro-api/internal/api/v1/handlers/users"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes registers v1 routes onto the provided router.
func RegisterRoutes(r chi.Router) {
	r.Route("/users", users.RegisterRoutes)
}
