package service

import (
	"github.com/dpolansky/ci/model"
	"github.com/dpolansky/ci/server/repo"
)

type RepositoryService interface {
	CreateRepository(m *model.Repository) error
	GetRepositoriesByOwner(owner string) ([]*model.Repository, error)
}

func NewRepositoryService(repo repo.RepositoryRepo) RepositoryService {
	return &repoService{
		repo: repo,
	}
}

type repoService struct {
	repo repo.RepositoryRepo
}

func (r *repoService) GetRepositoriesByOwner(owner string) ([]*model.Repository, error) {
	return r.repo.GetRepositoriesByOwner(owner)
}

func (r *repoService) CreateRepository(m *model.Repository) error {
	return r.repo.CreateRepository(m)
}
