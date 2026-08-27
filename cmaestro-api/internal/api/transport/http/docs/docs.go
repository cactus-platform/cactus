package docs

import (
	"encoding/json"
	"net/http"
)

type document struct {
	OpenAPI string          `json:"openapi"`
	Info    info            `json:"info"`
	Servers []server        `json:"servers"`
	Paths   map[string]path `json:"paths"`
}

type info struct {
	Title       string `json:"title"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

type server struct {
	URL string `json:"url"`
}

type path struct {
	Get  *operation `json:"get,omitempty"`
	Post *operation `json:"post,omitempty"`
}

type operation struct {
	Summary     string              `json:"summary,omitempty"`
	Parameters  []parameter         `json:"parameters,omitempty"`
	RequestBody *requestBody        `json:"requestBody,omitempty"`
	Responses   map[string]response `json:"responses"`
}

type parameter struct {
	Name     string `json:"name"`
	In       string `json:"in"`
	Required bool   `json:"required"`
	Schema   schema `json:"schema"`
}

type requestBody struct {
	Required bool               `json:"required"`
	Content  map[string]content `json:"content"`
}

type content struct {
	Schema schema `json:"schema"`
}

type response struct {
	Description string             `json:"description"`
	Content     map[string]content `json:"content,omitempty"`
}

type schema struct {
	Type       string            `json:"type,omitempty"`
	Format     string            `json:"format,omitempty"`
	Items      *schema           `json:"items,omitempty"`
	Properties map[string]schema `json:"properties,omitempty"`
	Required   []string          `json:"required,omitempty"`
}

func buildDocument() document {
	return document{
		OpenAPI: "3.0.3",
		Info: info{
			Title:       "Cmaestro API",
			Version:     "1.0.0",
			Description: "OpenAPI specification for the Cmaestro API.",
		},
		Servers: []server{{URL: "/"}},
		Paths: map[string]path{
			"/healthz": {
				Get: &operation{
					Summary: "Health check",
					Responses: map[string]response{
						"200": {Description: "Service is healthy"},
					},
				},
			},
			"/api/v1/users": {
				Get: &operation{
					Summary: "List users",
					Responses: map[string]response{
						"200": {
							Description: "List of users",
							Content: map[string]content{
								"application/json": {Schema: schema{Type: "array", Items: &schema{Type: "string"}}},
							},
						},
					},
				},
			},
			"/api/v1/users/{id}": {
				Get: &operation{
					Summary: "Get a user",
					Parameters: []parameter{{
						Name:     "id",
						In:       "path",
						Required: true,
						Schema:   schema{Type: "string"},
					}},
					Responses: map[string]response{
						"200": {Description: "User response"},
					},
				},
			},
			"/api/v1/repositories": {
				Get: &operation{
					Summary: "List repositories",
					Responses: map[string]response{
						"200": {
							Description: "Repository names",
							Content: map[string]content{
								"application/json": {Schema: schema{Type: "array", Items: &schema{Type: "string"}}},
							},
						},
					},
				},
				Post: &operation{
					Summary: "Create a repository artifact upload",
					RequestBody: &requestBody{
						Required: true,
						Content: map[string]content{
							"multipart/form-data": {
								Schema: schema{
									Type:     "object",
									Required: []string{"platform.cactus.repository.source", "platform.cactus.repository.name"},
									Properties: map[string]schema{
										"platform.cactus.repository.source": {Type: "string", Format: "binary"},
										"platform.cactus.repository.name":   {Type: "string"},
										"platform.cactus.repository.id":     {Type: "string"},
									},
								},
							},
						},
					},
					Responses: map[string]response{
						"201": {Description: "Artifact created"},
						"400": {Description: "Invalid request"},
					},
				},
			},
		},
	}
}

// JSON serves the OpenAPI document for the API.
func JSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(buildDocument())
}

// UI serves a minimal Swagger UI page backed by the local OpenAPI document.
func UI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Cmaestro API Docs</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
    <style>
      body { margin: 0; background: #0f172a; }
      #swagger-ui { background: white; }
    </style>
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
      window.onload = () => {
        window.ui = SwaggerUIBundle({
          url: '/openapi.json',
          dom_id: '#swagger-ui',
          deepLinking: true,
          presets: [SwaggerUIBundle.presets.apis],
          layout: 'BaseLayout'
        });
      };
    </script>
  </body>
</html>`))
}
