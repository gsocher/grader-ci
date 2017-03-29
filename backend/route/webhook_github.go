package route

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/dpolansky/ci/backend/service"
	"github.com/dpolansky/ci/model"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func RegisterGithubWebhookRoutes(router *mux.Router, run service.BuildRunner, rep service.RepositoryReadWriter) {
	router.HandleFunc(pathURLGithubWebhookAPI,
		parseWebhookHTTPHandler(run, rep)).Methods("POST")
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
func parseWebhookHTTPHandler(run service.BuildRunner, rep service.RepositoryReadWriter) func(rw http.ResponseWriter, req *http.Request) {
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
