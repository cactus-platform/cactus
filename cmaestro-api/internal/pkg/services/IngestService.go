package services

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"cmaestro-api/internal/pkg/models"

	"github.com/cactus-platform/cmaestro-core/storage/keyval"
	"github.com/google/uuid"
)

type IngestService interface {
	Ingest(ctx context.Context, artifactID uuid.UUID) error
}

type IngestServiceImpl struct {
	keyVal *keyval.Client
}

func NewIngestService(keyVal *keyval.Client) IngestService {
	return &IngestServiceImpl{keyVal: keyVal}
}

func (s *IngestServiceImpl) Ingest(ctx context.Context, artifactID uuid.UUID) error {
	if s.keyVal == nil {
		return errors.New("key-val client cannot be nil")
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	now := time.Now().UTC()
	value, err := json.Marshal(models.Ingest{
		ArtifactID: artifactID,
		Status:     models.IngestStatusPending,
		CreatedAt:  now,
		UpdatedAt:  now,
	})
	if err != nil {
		return err
	}

	return s.keyVal.Set(
		"ingest:"+artifactID.String(),
		string(value),
		-1,
	)
}
