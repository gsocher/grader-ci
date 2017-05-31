package service

import (
	"errors"
	"testing"

	"github.com/dpolansky/grader-ci/backend/service/fakes"
	"github.com/dpolansky/grader-ci/model"
)

func TestHandleRequest(t *testing.T) {
	fakeRepo := new(fakes.FakeRepositoryService)
	fakeMsg := new(fakes.FakeBuildMessageService)
	github := NewGithubWebhookService(&GithubWebhookServiceConfig{
		RepoService: fakeRepo,
		MsgService:  fakeMsg,
	})

	req := createFakeGithubWebhookRequest()
	if err := github.HandleRequest(req); err != nil {
		t.Fatalf("failed to handle github webhook: %v", err)
	}

	// check that dependencies were invoked properly
	if fakeRepo.UpdateRepositoryCallCount() != 1 {
		t.Fatalf("update called %v times", fakeRepo.UpdateRepositoryCallCount())
	}

	if fakeMsg.SendBuildCallCount() != 1 {
		t.Fatalf("send called %v times", fakeMsg.SendBuildCallCount())
	}
}

func TestHandleRequestErrorsIfRepoServiceErrors(t *testing.T) {
	fakeRepo := new(fakes.FakeRepositoryService)
	fakeMsg := new(fakes.FakeBuildMessageService)
	github := NewGithubWebhookService(&GithubWebhookServiceConfig{
		RepoService: fakeRepo,
		MsgService:  fakeMsg,
	})

	fakeRepo.UpdateRepositoryReturns(errors.New("update failed"))

	req := createFakeGithubWebhookRequest()
	if err := github.HandleRequest(req); err == nil {
		t.Fatalf("expected err")
	}
}

func TestHandleRequestErrorsIfMsgServiceErrors(t *testing.T) {
	fakeRepo := new(fakes.FakeRepositoryService)
	fakeMsg := new(fakes.FakeBuildMessageService)
	github := NewGithubWebhookService(&GithubWebhookServiceConfig{
		RepoService: fakeRepo,
		MsgService:  fakeMsg,
	})

	fakeMsg.SendBuildReturns(errors.New("send failed"))

	req := createFakeGithubWebhookRequest()
	if err := github.HandleRequest(req); err == nil {
		t.Fatalf("expected err")
	}
}

func createFakeGithubWebhookRequest() *model.GithubWebhookRequest {
	req := &model.GithubWebhookRequest{Ref: "refs/heads/master"}
	req.Repository.ID = 0
	req.Repository.Name = "foo"
	req.Repository.Owner.Name = "bar"
	req.Repository.Owner.AvatarURL = "avy"
	return req
}
