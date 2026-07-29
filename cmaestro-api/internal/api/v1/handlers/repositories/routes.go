package repositories

import (
	"cmaestro-api/internal/api/transport/http/request"
	"cmaestro-api/internal/api/transport/http/response"
	"cmaestro-api/internal/bootstrap"
	"cmaestro-db/bucket"
	"cmaestro-db/dbutil"
	"cmaestro-db/models"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"path"
	"time"

	"github.com/google/uuid"
)

type Handler struct {
	App *bootstrap.App
}

func NewHandler(app *bootstrap.App) *Handler {
	return &Handler{
		App: app,
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := []string{
		"cactus-plane",
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(
			w,
			"internal error",
			http.StatusInternalServerError,
		)
		return
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cfg := h.App.StaticConfig.Repositories

	var artifact models.Artifact

	artifactJSON := r.FormValue(
		cfg.ArtifactMetadataFieldName,
	)

	if artifactJSON == "" {
		response.Fail(
			w,
			response.APIError{
				Status:  http.StatusBadRequest,
				Code:    "missing_artifact",
				Message: "Missing artifact metadata",
			},
		)
		return
	}

	if err := json.Unmarshal(
		[]byte(artifactJSON),
		&artifact,
	); err != nil {

		log.Printf(
			"invalid artifact JSON: %v",
			err,
		)

		response.Fail(
			w,
			response.APIError{
				Status:  http.StatusBadRequest,
				Code:    "invalid_artifact",
				Message: "Invalid artifact metadata",
			},
		)

		return
	}

	if artifact.Id == uuid.Nil {
		artifact.Id = uuid.New()
	}

	now := time.Now()

	artifact.CreatedAt = now
	artifact.UpdatedAt = now

	data, err := request.WithMultipartFile(
		r,
		cfg.FileFieldName,
		cfg.MaxUploadSize,
		func(
			file multipart.File,
			header *multipart.FileHeader,
		) (any, error) {

			defer file.Close()

			key := fmt.Sprintf(
				"uploads/repositories/%s/%s.zip",
				artifact.Id,
				uuid.New(),
			)

			object, err := h.App.Connections.ArtifactStore.UploadZip(
				r.Context(),
				file,
				header.Size,
				key,
			)

			if err != nil {
				return nil, fmt.Errorf(
					"upload ZIP to artifact store: %w",
					err,
				)
			}

			log.Printf(
				"uploaded %s (%d bytes)",
				object.Key,
				object.Size,
			)

			artifact.Hash = dbutil.NewHasherReader(file).Hash()

			artifact.Path = path.Join(
				object.Bucket,
				object.Key,
			)

			artifact.Size = object.Size
			artifact.Format = "application/zip"
			artifact.Status = "uploaded"

			return &object, nil
		},
	)

	if err != nil {
		log.Printf(
			"upload failed: %v",
			err,
		)

		response.Fail(
			w,
			cfg.Errors.ErrorNameWhenUploadFails,
		)

		return
	}

	uploadedObject, ok := data.(*bucket.Object)

	if !ok || uploadedObject == nil {

		log.Printf(
			"unexpected upload result type: %T",
			data,
		)

		response.Fail(
			w,
			cfg.Errors.ErrorWhenUploadFails,
		)

		return
	}

	if err := h.App.Services.Artifact.CreateOrUpdateArtifact(
		r.Context(),
		&artifact,
	); err != nil {

		log.Printf(
			"artifact persistence failed: %v",
			err,
		)

		response.Fail(
			w,
			response.APIError{
				Status:  http.StatusInternalServerError,
				Code:    "internal_server_error",
				Message: err.Error(),
			},
		)

		return
	}

	response.Created(
		w,
		map[string]any{
			"status":     "created",
			"id":         artifact.Id,
			"name":       artifact.Name,
			"path":       artifact.Path,
			"revision":   artifact.Revision,
			"hash":       artifact.Hash,
			"size":       artifact.Size,
			"format":     artifact.Format,
			"created_at": artifact.CreatedAt,
			"updated_at": artifact.UpdatedAt,
			"object":     uploadedObject,
		},
	)
}
