package repositories

import (
	"cmaestro-api/internal/api/transport/http/docs"
	"cmaestro-api/internal/bootstrap"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, app *bootstrap.App, basePath string) {
	h := NewHandler(app)

	docs.Register("GET", basePath, docs.RouteMeta{
		Summary: "List repositories",
		Responses: map[string]docs.Response{
			"200": {
				Description: "Repository names",
				Content: map[string]docs.Content{
					"application/json": {Schema: docs.Schema{Type: "array", Items: &docs.Schema{Type: "string"}}},
				},
			},
		},
	})

	docs.Register("POST", basePath, docs.RouteMeta{
		Summary: "Create a repository artifact upload",
		RequestBody: &docs.RequestBody{
			Required: true,
			Content: map[string]docs.Content{
				"multipart/form-data": {
					Schema: docs.Schema{
						AnyOf: []docs.Schema{
							{
								Type:     "object",
								Required: []string{"platform.cactus.repository.source", "platform.cactus.repository.name"},
								Properties: map[string]docs.Schema{
									"platform.cactus.repository.source": {Type: "string", Format: "binary"},
									"platform.cactus.repository.name":   {Type: "string"},
								},
							},
							{
								Type:     "object",
								Required: []string{"platform.cactus.repository.source", "platform.cactus.repository.id"},
								Properties: map[string]docs.Schema{
									"platform.cactus.repository.source": {Type: "string", Format: "binary"},
									"platform.cactus.repository.id":     {Type: "string"},
								},
							},
						},
					},
				},
			},
		},
		Responses: map[string]docs.Response{
			"201": {Description: "Repository created"},
			"400": {Description: "Invalid request"},
		},
	})

	docs.Register("GET", basePath+"/{id}", docs.RouteMeta{
		Summary: "Get a repository artifact",
		Responses: map[string]docs.Response{
			"200": {Description: "Repository artifact"},
			"400": {Description: "Invalid repository ID"},
			"404": {Description: "Repository not found"},
		},
	})

	r.Get("/", h.List)
	r.Get("/{id}", h.Get)
	r.Post("/", h.Create)
}
