package route

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/dpolansky/grader-ci/backend/service"
	"github.com/dpolansky/grader-ci/model"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func RegisterGithubWebhookRoutes(router *mux.Router, run service.BuildRunner, rep service.RepositoryService, bind service.TestBindService) {
	router.HandleFunc(pathURLGithubWebhookAPI,
		parseWebhookHTTPHandler(run, rep, bind)).Methods("POST")
}

type githubWebhookRequest struct {
	Ref        string `json:"ref"`
	Repository struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Owner struct {
			Name      string `json:"login"`
			AvatarURL string `json:"avatar_url"`
		} `json:"owner"`
	} `json:"repository"`
}

// parseWebhookHTTPHandler is an endpoint for receiving github webhook requests.
func parseWebhookHTTPHandler(run service.BuildRunner, rep service.RepositoryService, bind service.TestBindService) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			writeError(rw, http.StatusBadRequest, err)
			return
		}

		var r githubWebhookRequest
		if err = json.Unmarshal(body, &r); err != nil {
			writeError(rw, http.StatusBadRequest, err)
			return
		}

		repo := &model.Repository{
			ID:        r.Repository.ID,
			Name:      r.Repository.Name,
			Owner:     r.Repository.Owner.Name,
			AvatarURL: r.Repository.Owner.AvatarURL,
		}

		if err = rep.UpdateRepository(repo); err != nil {
			writeError(rw, http.StatusInternalServerError, err)
			return
		}

		m := &model.BuildStatus{
			Source: &model.RepositoryMetadata{
				ID:       repo.ID,
				Branch:   filepath.Base(r.Ref),
				CloneURL: fmt.Sprintf("https://github.com/%s/%s", r.Repository.Owner.Name, r.Repository.Name),
			},
		}

		// if this repository has a test binding, update the build to include it
		if bind, err := bind.GetTestBindBySourceRepositoryID(repo.ID); err == nil {
			testRepo, err := rep.GetRepositoryByID(bind.TestID)
			if err != nil {
				logrus.WithError(err).Errorf("Failed to get test repository using bind: %v", err)
				return
			}

			m.Tested = true
			m.Test = &model.RepositoryMetadata{
				ID:       bind.TestID,
				CloneURL: fmt.Sprintf("https://github.com/%s/%s", testRepo.Owner, testRepo.Name),
				Branch:   bind.TestBranch,
			}
		}

		status, err := run.RunBuild(m)
		if err != nil {
			logrus.WithError(err).WithField("req", string(body)).Errorf("Failed to start build")
			writeError(rw, http.StatusInternalServerError, err)
			return
		}

		bytes, err := json.Marshal(status)
		if err != nil {
			logrus.WithField("req", string(body)).WithError(err).Errorf("Failed to marshal build status")
			writeError(rw, http.StatusInternalServerError, err)
			return
		}

		writeOk(rw, bytes)
	}
}
