package repositories

import (
	"cmaestro-api/internal/api/transport/http/response"
	"fmt"
	"net/http"
)

type Config struct {
	SourceCodeUploadKey string `json:"sourceCodeUploadKey"`
	MaxUploadSize       int64  `json:"maxUploadSize"`
	Errors              *Errors
}

type Errors struct {
	ErrorNameWhenUploadFails response.APIError `json:"error_name_when_upload_fails"`
}

func Load() *Config {
	SCUK := "platform.cactus.repository.source"

	return &Config{
		SourceCodeUploadKey: SCUK,
		MaxUploadSize:       10 << 20,
		Errors: &Errors{
			ErrorNameWhenUploadFails: response.APIError{
				Status:  http.StatusBadRequest,
				Code:    "INVALID_SOURCE_UPLOAD",
				Message: fmt.Sprintf("Invalid source code upload, key=[%s] is undefined or invalid", SCUK),
			},
		},
	}
}
