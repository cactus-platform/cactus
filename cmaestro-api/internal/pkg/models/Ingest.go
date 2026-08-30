package models

import (
	"time"

	"github.com/google/uuid"
)

type IngestStatus string

const (
	IngestStatusPending IngestStatus = "pending"
)

type Ingest struct {
	ArtifactID uuid.UUID    `json:"artifact_id"`
	Status     IngestStatus `json:"status"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}
