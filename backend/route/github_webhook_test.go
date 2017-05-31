package route

import (
	"bytes"
	"encoding/json"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/dpolansky/grader-ci/backend/service/fakes"
	"github.com/dpolansky/grader-ci/model"
	"github.com/gorilla/mux"
)

func TestGithubWebhook(t *testing.T) {
	router := mux.NewRouter()
	fakeGithubWebhookService := new(fakes.FakeGithubWebhookService)
	RegisterGithubWebhookRoutes(router, fakeGithubWebhookService)

	serv := httptest.NewServer(router)
	defer serv.Close()
	resp := httptest.NewRecorder()

	webhook := createFakeGithubWebhookRequest()
	b, _ := json.Marshal(webhook)
	req, _ := http.NewRequest("POST", pathURLGithubWebhookAPI, bytes.NewBuffer(b))

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status ok, got %v", resp.Code)
	}

	if fakeGithubWebhookService.HandleRequestCallCount() != 1 {
		t.Fatalf("expected service to handle request, did it %v times", fakeGithubWebhookService.HandleRequestCallCount())
	}
}

func createFakeGithubWebhookRequest() model.GithubWebhookRequest {
	req := model.GithubWebhookRequest{Ref: "refs/heads/master"}
	req.Repository.ID = 0
	req.Repository.Name = "foo"
	req.Repository.Owner.Name = "bar"
	req.Repository.Owner.AvatarURL = "avy"
	return req
}
