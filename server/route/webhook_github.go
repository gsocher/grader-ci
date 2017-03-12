package route

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dpolansky/ci/server/service"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const pathURLGithubWebhookAPI = "/api/github"

func RegisterGithubWebhookRoutes(router *mux.Router, builder service.Builder) {
	router.HandleFunc(pathURLGithubWebhookAPI, parseWebhookHTTPHandler(builder)).Methods("POST")
}

type githubWebhookRequest struct {
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

// parseWebhookHTTPHandler is an endpoint for receiving github webhook requests.
func parseWebhookHTTPHandler(builder service.Builder) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			writeError(rw, http.StatusBadRequest, err)
			return
		}

		logrus.WithField("req", string(body)).Infof("Received request")

		var r githubWebhookRequest
		if err = json.Unmarshal(body, &r); err != nil {
			writeError(rw, http.StatusBadRequest, err)
			return
		}

		cloneURL := fmt.Sprintf("github.com/%s", r.Repository.FullName)
		status, err := builder.StartBuild(cloneURL)
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
