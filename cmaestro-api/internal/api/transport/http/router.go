package httptransport

import (
	"cmaestro-api/internal/config"
	"net/http"

	v1 "cmaestro-api/internal/api/v1"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter creates a chi router, applies common middleware and registers API routes.
func NewRouter(cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	// common middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// versioned API registration
	r.Route("/api", func(sr chi.Router) {
		sr.Route("/v1", func(rchi chi.Router) {
			v1.RegisterRoutes(rchi, cfg)
		})
	})

	// healthcheck
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	return r
}
