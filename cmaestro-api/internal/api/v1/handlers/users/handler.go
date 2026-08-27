package users

import (
	"cmaestro-api/internal/api/transport/http/docs"
	"cmaestro-api/internal/bootstrap"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, app *bootstrap.App, basePath string) {
	h := NewHandler(app)

	docs.Register("GET", basePath, docs.RouteMeta{
		Summary: "List users",
		Responses: map[string]docs.Response{
			"200": {
				Description: "List of users",
				Content: map[string]docs.Content{
					"application/json": {Schema: docs.Schema{Type: "array", Items: &docs.Schema{Type: "string"}}},
				},
			},
		},
	})

	docs.Register("GET", basePath+"/{id}", docs.RouteMeta{
		Summary: "Get a user",
		Responses: map[string]docs.Response{
			"200": {Description: "User response"},
		},
	})

	r.Get("/", h.List)
	r.Get("/{id}", h.Get)
}
