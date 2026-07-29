package models

import (
	"time"

	"github.com/google/uuid"
)

type Artifact struct {
	Id uuid.UUID `json:"id" gorm:"primaryKey"`

	Name string `json:"name" gorm:"unique"`
	Path string `json:"path" gorm:"unique"`

	Revision string `json:"revision"`

	Hash   string `json:"hash"`
	Size   int64  `json:"size"`
	Format string `json:"format"`

	Status string `json:"status"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
