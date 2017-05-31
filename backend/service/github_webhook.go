package service

import (
	"fmt"
	"path/filepath"

	"github.com/dpolansky/grader-ci/model"
)

type GithubWebhookService interface {
	HandleRequest(*model.GithubWebhookRequest) error
}

type githubService struct {
	config *GithubWebhookServiceConfig
}

type GithubWebhookServiceConfig struct {
	RepoService RepositoryService
	MsgService  BuildMessageService
}

func NewGithubWebhookService(cfg *GithubWebhookServiceConfig) GithubWebhookService {
	return &githubService{config: cfg}
}

func (g *githubService) HandleRequest(req *model.GithubWebhookRequest) error {
	// update the repository to reflect any changes
	if err := g.config.RepoService.UpdateRepository(&model.Repository{
		ID:        req.Repository.ID,
		Name:      req.Repository.Name,
		Owner:     req.Repository.Owner.Name,
		AvatarURL: req.Repository.Owner.AvatarURL,
	}); err != nil {
		return fmt.Errorf("failed to update repository: %v", err)
	}

	// notify the build runner to trigger a build
	return g.config.MsgService.SendBuild(&model.BuildStatus{
		Source: &model.RepositoryMetadata{
			ID:       req.Repository.ID,
			Branch:   filepath.Base(req.Ref),
			CloneURL: fmt.Sprintf("https://github.com/%s/%s", req.Repository.Owner.Name, req.Repository.Name),
		}})
}
