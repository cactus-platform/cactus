package repositories

import (
	"cmaestro-api/internal/api/transport/http/response"
	"fmt"
	"net/http"
)

type Config struct {
	SourceCodeUploadKey string `json:"sourceCodeUploadKey"`
	SourceCodeIdKey     string `json:"sourceCodeIdKey"`
	MaxUploadSize       int64  `json:"maxUploadSize"`
	Errors              *Errors
}

type Errors struct {
	ErrorNameWhenUploadFails response.APIError `json:"error_name_when_upload_fails"`
	ErrorWhenUploadFails     response.APIError `json:"error_when_upload_fails"`
	ErrorWhenHashingFails    response.APIError `json:"error_when_hashing_fails"`
}

func Load() *Config {
	SCUK := "platform.cactus.repository.source"
	SCIK := "platform.cactus.repository.id"

	return &Config{
		SourceCodeUploadKey: SCUK,
		SourceCodeIdKey:     SCIK,
		MaxUploadSize:       10 << 20,
		Errors: &Errors{
			ErrorNameWhenUploadFails: response.APIError{
				Status:  http.StatusBadRequest,
				Code:    "INVALID_SOURCE_UPLOAD",
				Message: fmt.Sprintf("Invalid source code upload, key=[%s] is undefined or invalid", SCUK),
			},
			ErrorWhenUploadFails: response.APIError{
				Status:  http.StatusInternalServerError,
				Code:    "INVALID_SOURCE_UPLOAD",
				Message: "Error occurred during upload to Cactus Artifact Database",
			},
			ErrorWhenHashingFails: response.APIError{
				Status:  http.StatusInternalServerError,
				Code:    "INVALID_SOURCE_HASHING",
				Message: "Invalid source code hashing, Error occurred during upload to Cactus Artifact Database",
			},
		},
	}
}
