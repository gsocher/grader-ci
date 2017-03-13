package service

import (
	"github.com/dpolansky/ci/model"
	"github.com/dpolansky/ci/server/repo"
)

type RepositoryService interface {
}

type repoService struct {
	repo repo.RepositoryRepo
}

func (r *repoService) GetRepositoriesByOwner(owner string) ([]*model.Repository, error) {
	return r.repo.GetRepositoriesByOwner(owner)
}

func (r *repoService) CreateRepository(cloneURL, owner string) (*model.Repository, error) {
	m := &model.Repository{
		Owner:    owner,
		CloneURL: cloneURL,
	}

	err = r.repo.CreateRepository(m)
	if err != nil {
		return nil, err
	}

	return m, err
}
