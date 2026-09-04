package docs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"
)

type document struct {
	OpenAPI string                           `json:"openapi"`
	Info    info                             `json:"info"`
	Servers []server                         `json:"servers"`
	Paths   map[string]map[string]*operation `json:"paths"`
}

type info struct {
	Title       string `json:"title"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

type server struct {
	URL string `json:"url"`
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
	AnyOf      []schema          `json:"anyOf,omitempty"`
}

type RouteMeta struct {
	Summary     string
	RequestBody *RequestBody
	Responses   map[string]Response
}

type RequestBody struct {
	Required bool               `json:"required"`
	Content  map[string]Content `json:"content"`
}

type Content struct {
	Schema Schema `json:"schema"`
}

type Response struct {
	Description string             `json:"description"`
	Content     map[string]Content `json:"content,omitempty"`
}

type Schema struct {
	Type       string            `json:"type,omitempty"`
	Format     string            `json:"format,omitempty"`
	Items      *Schema           `json:"items,omitempty"`
	Properties map[string]Schema `json:"properties,omitempty"`
	Required   []string          `json:"required,omitempty"`
	AnyOf      []Schema          `json:"anyOf,omitempty"`
}

var (
	pathParamPattern = regexp.MustCompile(`\{([^}]+)\}`)
	metaMu           sync.RWMutex
	routeRegistry    = map[string]map[string]RouteMeta{}
)

// Register stores OpenAPI metadata for a method and route.
func Register(method, route string, meta RouteMeta) {
	method = strings.ToLower(method)
	route = normalizeRoute(route)

	metaMu.Lock()
	defer metaMu.Unlock()

	if _, ok := routeRegistry[route]; !ok {
		routeRegistry[route] = map[string]RouteMeta{}
	}
	routeRegistry[route][method] = meta
}

func buildDocument(router chi.Routes) (document, error) {
	paths := make(map[string]map[string]*operation)

	metaMu.RLock()
	registry := make(map[string]map[string]RouteMeta, len(routeRegistry))
	for route, methods := range routeRegistry {
		registry[route] = make(map[string]RouteMeta, len(methods))
		for method, meta := range methods {
			registry[route][method] = meta
		}
	}
	metaMu.RUnlock()

	err := chi.Walk(router, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		_ = handler
		_ = middlewares

		normalizedRoute := normalizeRoute(route)
		specByMethod, ok := registry[normalizedRoute]
		if !ok {
			return nil
		}

		methodMeta, ok := specByMethod[strings.ToLower(method)]
		if !ok {
			return nil
		}

		if _, exists := paths[normalizedRoute]; !exists {
			paths[normalizedRoute] = make(map[string]*operation)
		}

		paths[normalizedRoute][strings.ToLower(method)] = &operation{
			Summary:     methodMeta.Summary,
			Parameters:  pathParameters(normalizedRoute),
			RequestBody: convertRequestBody(methodMeta.RequestBody),
			Responses:   convertResponses(methodMeta.Responses),
		}

		return nil
	})
	if err != nil {
		return document{}, err
	}

	routes := make([]string, 0, len(paths))
	for route := range paths {
		routes = append(routes, route)
	}
	sort.Strings(routes)

	resultPaths := make(map[string]map[string]*operation, len(paths))
	for _, route := range routes {
		resultPaths[route] = paths[route]
	}

	return document{
		OpenAPI: "3.0.3",
		Info: info{
			Title:       "Cmaestro API",
			Version:     "1.0.0",
			Description: "OpenAPI specification for the Cmaestro API.",
		},
		Servers: []server{{URL: "/"}},
		Paths:   resultPaths,
	}, nil
}

func normalizeRoute(route string) string {
	if route == "/" {
		return route
	}

	return strings.TrimSuffix(route, "/")
}

func pathParameters(route string) []parameter {
	matches := pathParamPattern.FindAllStringSubmatch(route, -1)
	if len(matches) == 0 {
		return nil
	}

	parameters := make([]parameter, 0, len(matches))
	seen := make(map[string]struct{}, len(matches))
	for _, match := range matches {
		name := match[1]
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		parameters = append(parameters, parameter{
			Name:     name,
			In:       "path",
			Required: true,
			Schema:   schema{Type: "string"},
		})
	}

	return parameters
}

func convertRequestBody(body *RequestBody) *requestBody {
	if body == nil {
		return nil
	}

	convertedContent := make(map[string]content, len(body.Content))
	for mimeType, item := range body.Content {
		convertedContent[mimeType] = content{Schema: convertSchema(item.Schema)}
	}

	return &requestBody{Required: body.Required, Content: convertedContent}
}

func convertResponses(responses map[string]Response) map[string]response {
	if responses == nil {
		return nil
	}

	converted := make(map[string]response, len(responses))
	for code, item := range responses {
		responseContent := make(map[string]content, len(item.Content))
		for mimeType, contentItem := range item.Content {
			responseContent[mimeType] = content{Schema: convertSchema(contentItem.Schema)}
		}
		converted[code] = response{Description: item.Description, Content: responseContent}
	}

	return converted
}

func convertSchema(s Schema) schema {
	properties := make(map[string]schema, len(s.Properties))
	for key, item := range s.Properties {
		properties[key] = convertSchema(item)
	}
	anyOf := make([]schema, len(s.AnyOf))
	for index, item := range s.AnyOf {
		anyOf[index] = convertSchema(item)
	}

	var items *schema
	if s.Items != nil {
		converted := convertSchema(*s.Items)
		items = &converted
	}

	return schema{
		Type:       s.Type,
		Format:     s.Format,
		Items:      items,
		Properties: properties,
		Required:   s.Required,
		AnyOf:      anyOf,
	}
}

// JSON serves the OpenAPI document for the API.
func JSON(router chi.Routes) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		document, err := buildDocument(router)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to build openapi document: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(document)
	}
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
