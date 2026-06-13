package users

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// UsersList returns a simple users list.
func (h *Handler) UsersList(w http.ResponseWriter, _ *http.Request) {
	users := []string{"admin"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func RegisterRoutes(r chi.Router) {
	h := NewHandler()

	r.Get("/", h.UsersList)
}
