package service

import (
	"errors"
	"testing"

	"github.com/dpolansky/grader-ci/backend/service/fakes"
)

func TestHandleRequest(t *testing.T) {
	fakeRepo := new(fakes.FakeRepositoryService)
	fakeMsg := new(fakes.FakeBuildMessageService)
	github := NewGithubWebhookService(&GithubWebhookServiceConfig{
		repoService: fakeRepo,
		msgService:  fakeMsg,
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
		repoService: fakeRepo,
		msgService:  fakeMsg,
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
		repoService: fakeRepo,
		msgService:  fakeMsg,
	})

	fakeMsg.SendBuildReturns(errors.New("send failed"))

	req := createFakeGithubWebhookRequest()
	if err := github.HandleRequest(req); err == nil {
		t.Fatalf("expected err")
	}
}

func createFakeGithubWebhookRequest() GithubWebhookRequest {
	req := GithubWebhookRequest{Ref: "refs/heads/master"}
	req.Repository.ID = 0
	req.Repository.Name = "foo"
	req.Repository.Owner.Name = "bar"
	req.Repository.Owner.AvatarURL = "avy"
	return req
}