package repositories

import (
	"cmaestro-db/sql"
	"context"
	"errors"

	"cmaestro-db/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrArtifactIDAlreadyExists = errors.New("artifact id already exists")
	ErrArtifactNotFound        = errors.New("artifact not found")
)

type ArtifactRepository interface {
	CreateArtifact(ctx context.Context, artifact *models.Artifact) error
	GetArtifact(ctx context.Context, id uuid.UUID) (*models.Artifact, error)
	UpdateArtifact(ctx context.Context, artifact *models.Artifact) error
	ArtifactExists(ctx context.Context, id uuid.UUID) (bool, error)
}

type ArtifactRepositoryImpl struct {
	db *sql.Client
}

func NewArtifactRepository(sqlDb *sql.Client) ArtifactRepository {
	return &ArtifactRepositoryImpl{
		db: sqlDb,
	}
}

func (r *ArtifactRepositoryImpl) CreateArtifact(
	ctx context.Context,
	artifact *models.Artifact,
) error {
	if artifact == nil {
		return errors.New("artifact cannot be nil")
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64

		err := tx.Model(&models.Artifact{}).
			Where("id = ?", artifact.Id).
			Count(&count).
			Error
		if err != nil {
			return err
		}

		if count > 0 {
			return ErrArtifactIDAlreadyExists
		}

		return tx.Create(artifact).Error
	})
}

func (r *ArtifactRepositoryImpl) GetArtifact(
	ctx context.Context,
	id uuid.UUID,
) (*models.Artifact, error) {
	var artifact models.Artifact

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&artifact).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrArtifactNotFound
	}
	if err != nil {
		return nil, err
	}

	return &artifact, nil
}

func (r *ArtifactRepositoryImpl) UpdateArtifact(
	ctx context.Context,
	artifact *models.Artifact,
) error {
	if artifact == nil {
		return errors.New("artifact cannot be nil")
	}

	result := r.db.WithContext(ctx).
		Model(&models.Artifact{}).
		Where("id = ?", artifact.Id).
		Updates(artifact)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrArtifactNotFound
	}

	return nil
}

func (r *ArtifactRepositoryImpl) ArtifactExists(
	ctx context.Context,
	id uuid.UUID,
) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.Artifact{}).
		Where("id = ?", id).
		Count(&count).
		Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
