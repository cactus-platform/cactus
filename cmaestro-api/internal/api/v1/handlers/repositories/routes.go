package repositories

import (
	"cmaestro-api/internal/api/transport/http/request"
	"cmaestro-api/internal/api/transport/http/response"
	"cmaestro-api/internal/bootstrap"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/cactus-platform/cmaestro-core/models"
	coreRepositories "github.com/cactus-platform/cmaestro-core/repositories"
	"github.com/cactus-platform/cmaestro-core/storage/dbutil"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	App *bootstrap.App
}

type createResponse struct {
	Status     models.IngestStatus `json:"status"`
	Repository *models.Repository  `json:"repository"`
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

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	artifactID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Fail(w, response.APIError{
			Status:  http.StatusBadRequest,
			Code:    "invalid_repository_id",
			Message: "Invalid repository ID",
		})
		return
	}

	artifact, err := h.App.Services.Repository.Get(r.Context(), artifactID)
	if err != nil {
		if errors.Is(err, coreRepositories.ErrRepositoryNotFound) {
			response.Fail(w, response.APIError{
				Status:  http.StatusNotFound,
				Code:    "repository_not_found",
				Message: "Repository not found",
			})
			return
		}

		log.Printf("repository lookup failed: %v", err)
		response.Fail(w, response.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "internal_server_error",
			Message: "Unable to retrieve repository",
		})
		return
	}

	response.OK(w, artifact)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cfg := h.App.StaticConfig.Repositories

	repository, err := loadRepositoryFromRequest(
		r,
		cfg.RepositoryNameFieldName,
		cfg.RepositoryIDFieldName,
	)
	if err != nil {
		if errors.Is(err, errMissingArtifactMetadata) {
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

		log.Printf(
			"invalid artifact metadata: %v",
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

	isNewRepository := repository.ID == uuid.Nil
	if isNewRepository {
		repository.ID = uuid.New()
	} else {
		existingRepository, err := h.App.Services.Repository.Get(r.Context(), repository.ID)
		if err != nil {
			if errors.Is(err, coreRepositories.ErrRepositoryNotFound) {
				response.Fail(w, response.APIError{
					Status:  http.StatusNotFound,
					Code:    "repository_not_found",
					Message: "Repository not found",
				})
				return
			}

			log.Printf("repository lookup failed: %v", err)
			response.Fail(w, response.APIError{
				Status:  http.StatusInternalServerError,
				Code:    "internal_server_error",
				Message: "Unable to retrieve repository",
			})
			return
		}

		repository = existingRepository
	}

	artifact := &models.Artifact{
		Id:   uuid.New(),
		Name: repository.Name,
	}
	repository.Artifacts = []*models.Artifact{artifact}

	now := time.Now()

	if isNewRepository {
		repository.CreatedAt = now
	}
	repository.UpdatedAt = now
	artifact.CreatedAt = now
	artifact.UpdatedAt = now

	_, err = request.WithMultipartFile(
		r,
		cfg.FileFieldName,
		cfg.MaxUploadSize,
		func(
			file multipart.File,
			header *multipart.FileHeader,
		) (any, error) {

			defer file.Close()

			hasher := dbutil.NewHasherReader(file)
			revision := uuid.NewString()

			key := fmt.Sprintf(
				"uploads/repositories/%s/%s.zip",
				repository.ID,
				revision,
			)

			object, err := h.App.Connections.ArtifactStore.UploadZip(
				r.Context(),
				hasher,
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

			artifact.Hash = hasher.Hash()
			artifact.Revision = revision

			artifact.Path = path.Join(
				object.Bucket,
				object.Key,
			)

			artifact.Size = object.Size
			artifact.Format = "application/zip"
			artifact.Status = "uploaded"
			repository.Status = artifact.Status

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

	if err := h.App.Services.Repository.CreateOrUpdate(
		r.Context(),
		repository,
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

	if err := h.App.Services.Ingest.Ingest(r.Context(), repository.ID); err != nil {
		log.Printf(
			"artifact ingest failed: %v",
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
		createResponse{
			Status:     models.IngestStatusPending,
			Repository: repository,
		},
	)
}

var errMissingArtifactMetadata = errors.New("missing artifact metadata")

func loadRepositoryFromRequest(
	r *http.Request,
	repositoryNameFieldName string,
	repositoryIDFieldName string,

) (*models.Repository, error) {
	repository := &models.Repository{
		Name: strings.TrimSpace(r.FormValue(repositoryNameFieldName)),
	}

	if repository.Name == "" {
		return nil, errMissingArtifactMetadata
	}

	if idValue := strings.TrimSpace(r.FormValue(repositoryIDFieldName)); idValue != "" {
		artifactID, err := uuid.Parse(idValue)
		if err != nil {
			return nil, fmt.Errorf("invalid repository id: %w", err)
		}

		repository.ID = artifactID
	}

	if repository.ID != uuid.Nil {
		return repository, nil
	}

	if repository.Name == "" {
		return nil, errMissingArtifactMetadata
	}

	return repository, nil
}
