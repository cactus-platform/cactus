package users

import (
	"cmaestro-api/internal/config"
	"encoding/json"
	"net/http"
)

type Handler struct {
	App *config.AppContext
}

func NewHandler(app *config.AppContext) *Handler {
	return &Handler{
		App: app,
	}
}

/*
func NewHandler(
	userService *service.UserService,
	logger *slog.Logger,
) *Handler {
	return &Handler{
		service: userService,
		logger:  logger,
	}
}
*/

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	users := []string{"admin"}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	// or chi.URLParam(r, "id")

	resp := map[string]string{
		"id":   id,
		"name": "admin",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
