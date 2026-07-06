package repositories

import (
	"cmaestro-db/repositories"
	"errors"

	"github.com/google/uuid"
)

type RepositoryCtrl interface {
	CreateOrUpdate(repo *Repository) (uuid.UUID, error)
	Get(id uuid.UUID) (*Repository, error)
	Delete(id uuid.UUID) error
}

type RepositoryCtrlImpl struct {
	repository *repositories.ArtifactRepository
}

func New(artifactRepository *repositories.ArtifactRepository) RepositoryCtrl {
	return &RepositoryCtrlImpl{
		repository: artifactRepository,
	}
}

func (r *RepositoryCtrlImpl) CreateOrUpdate(repo *Repository) (uuid.UUID, error) {
	return uuid.Nil, errors.New("not implemented")
}

func (r *RepositoryCtrlImpl) Get(id uuid.UUID) (*Repository, error) {
	return nil, errors.New("not implemented")
}

func (r *RepositoryCtrlImpl) Delete(id uuid.UUID) error {
	return errors.New("not implemented")
}
