package repositories

import (
	"encoding/json"
	"net/http"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
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
	users := []string{"cactus-plane"}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := map[string]any{
		"status":     "created",                  // "created" | "updated" | "failed"
		"id":         "new-uuid-here",            // upload id
		"name":       "admin",                    // repository id
		"path":       "/",                        // submission path
		"size":       12345,                      // submission size (only attached file, body params doesn't count)
		"format":     "tar.gz",                   // compression algorithm ("tar.gz" || "zip") | "tar" also works, but it's not compressed
		"revision":   0,                          // number increases at each repository submission
		"hash":       "a8728942f927248724ff4...", // submission hash
		"created_at": "1970-01-01 00:00:01",      // first submission received at
		"updated_at": "1970-01-01 00:00:01",      // latest submission received at
	}

	json.NewEncoder(w).Encode(resp)
}
