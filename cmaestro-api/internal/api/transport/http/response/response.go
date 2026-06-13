package response

import (
	"encoding/json"
	"net/http"
)

type Envelope struct {
	Data  any       `json:"data,omitempty"`
	Error *APIError `json:"error,omitempty"`
}

func Write(
	w http.ResponseWriter,
	status int,
	payload Envelope,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(payload)
}
