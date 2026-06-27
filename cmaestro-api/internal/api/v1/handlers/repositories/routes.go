package repositories

import (
	"cmaestro-api/internal/api/transport/http/request"
	"cmaestro-api/internal/api/transport/http/response"
	"cmaestro-api/internal/config"
	"cmaestro-db/bucket"
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
	users := []string{"cactus-plane"}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	c := h.App.Config.Repositories

	var sourceCodeRepositoryId string
	if id := r.FormValue(c.SourceCodeIdKey); id != "" {
		sourceCodeRepositoryId = id
	} else {
		sourceCodeRepositoryId = uuid.New().String()
	}

	subCreatedAt := time.Now()

	data, err := request.WithMultipartFile(
		r,
		c.SourceCodeUploadKey,
		c.MaxUploadSize,
		func(file multipart.File, header *multipart.FileHeader) (any, error) {
			defer file.Close()

			object, err := h.App.ArtifactDB.UploadZip(
				r.Context(),
				file,
				header.Size,
				fmt.Sprintf("uploads/repositories/%s/%s.zip", sourceCodeRepositoryId, sourceCodeRepositoryId),
			)
			if err != nil {
				return nil, fmt.Errorf("upload ZIP to SeaweedFS: %w", err)
			}

			log.Printf("uploaded %s (%d bytes)", object.Key, object.Size)

			return &object, nil
		},
	)
	if err != nil {
		log.Printf("error uploading file to SeaweedFS: %v", err)
		response.Fail(w, c.Errors.ErrorNameWhenUploadFails)
		return
	}

	uploaded, ok := data.(*bucket.Object)
	if !ok || uploaded == nil {
		log.Printf("unexpected upload result type: %T", data)
		response.Fail(w, c.Errors.ErrorNameWhenUploadFails)
		return
	}

	// Use it here.
	log.Printf(
		"repository source uploaded: bucket=%s key=%s size=%d etag=%s",
		uploaded.Bucket,
		uploaded.Key,
		uploaded.Size,
		uploaded.ETag,
	)

	artifactKey := uploaded.Key
	artifactSize := uploaded.Size

	resp := map[string]any{
		"status":     "created",                               // "created" | "updated" | "failed"
		"id":         sourceCodeRepositoryId,                  // upload id
		"name":       "admin",                                 // repository id
		"path":       path.Join(uploaded.Bucket, artifactKey), // submission path
		"size":       artifactSize,                            // submission size (only attached file, body params doesn't count)
		"format":     "application/zip",                       // compression algorithm ("tar.gz" || "zip") | "tar" also works, but it's not compressed
		"revision":   0,                                       // number increases at each repository submission
		"hash":       "a8728942f927248724ff4...",              // submission hash
		"created_at": subCreatedAt,                            // first submission received at
		"updated_at": time.Now(),                              // latest submission received at
	}

	response.Created(w, resp)
	//json.NewEncoder(w).Encode(resp)
}
