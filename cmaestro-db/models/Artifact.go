package models

import (
	"time"

	"github.com/google/uuid"
)

type Artifact struct {
	Id                  uuid.UUID `json:"id" gorm:"primaryKey"`
	Name                string    `json:"name" gorm:"unique"`
	Path                string    `json:"path" gorm:"unique"`
	Revision            string    `json:"revision"`
	SubmissionHistoryId uuid.UUID `json:"submission_history_id" gorm:"index"`
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
