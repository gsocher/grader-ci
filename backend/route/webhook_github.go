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

const pathURLGithubWebhookAPI = "/api/webhook/github"

func RegisterGithubWebhookRoutes(router *mux.Router, run service.BuildRunner) {
	router.HandleFunc(pathURLGithubWebhookAPI,
		parseWebhookHTTPHandler(run)).Methods("POST")
}

type githubWebhookRequest struct {
	Ref        string `json:"ref"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

// parseWebhookHTTPHandler is an endpoint for receiving github webhook requests.
func parseWebhookHTTPHandler(run service.BuildRunner) func(rw http.ResponseWriter, req *http.Request) {
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

		m := &model.BuildStatus{
			CloneURL: fmt.Sprintf("https://github.com/%s", r.Repository.FullName),
			Branch:   filepath.Base(r.Ref),
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