package repositories

import (
	"cmaestro-api/internal/api/transport/http/response"
	"fmt"
	"net/http"
)

type Config struct {
	// Multipart field containing the uploaded artifact file
	FileFieldName string `json:"fileFieldName"`

	// Multipart field containing the Artifact JSON metadata
	ArtifactMetadataFieldName string `json:"artifactMetadataFieldName"`

	// Multipart field containing an optional existing repository ID
	RepositoryIDFieldName string `json:"repositoryIDFieldName"`

	MaxUploadSize int64 `json:"maxUploadSize"`

	Errors *Errors
}

type Errors struct {
	ErrorNameWhenUploadFails response.APIError `json:"error_name_when_upload_fails"`
	ErrorWhenUploadFails     response.APIError `json:"error_when_upload_fails"`
	ErrorWhenHashingFails    response.APIError `json:"error_when_hashing_fails"`
}

func Load() *Config {
	fileField := "platform.cactus.repository.source"
	repositoryIDField := "platform.cactus.repository.id"
	artifactMetadataField := "artifact"

	return &Config{
		FileFieldName:             fileField,
		ArtifactMetadataFieldName: artifactMetadataField,
		RepositoryIDFieldName:     repositoryIDField,
		MaxUploadSize:             10 << 20,

		Errors: &Errors{
			ErrorNameWhenUploadFails: response.APIError{
				Status: http.StatusBadRequest,
				Code:   "INVALID_SOURCE_UPLOAD",
				Message: fmt.Sprintf(
					"Invalid source code upload, field=[%s] is undefined or invalid",
					fileField,
				),
			},

			ErrorWhenUploadFails: response.APIError{
				Status:  http.StatusInternalServerError,
				Code:    "INVALID_SOURCE_UPLOAD",
				Message: "Error occurred during upload to Cactus Artifact Database",
			},

			ErrorWhenHashingFails: response.APIError{
				Status:  http.StatusInternalServerError,
				Code:    "INVALID_SOURCE_HASHING",
				Message: "Invalid source code hashing, error occurred during upload to Cactus Artifact Database",
			},
		},
	}
}
