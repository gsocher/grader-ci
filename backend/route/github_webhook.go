package route

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/dpolansky/grader-ci/backend/service"
	"github.com/dpolansky/grader-ci/model"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func RegisterGithubWebhookRoutes(router *mux.Router, githubWebhookService service.GithubWebhookService) {
	router.HandleFunc(pathURLGithubWebhookAPI,
		parseWebhookHTTPHandler(githubWebhookService)).Methods("POST")
}

// parseWebhookHTTPHandler is an endpoint for receiving github webhook requests.
func parseWebhookHTTPHandler(githubWebhookService service.GithubWebhookService) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		// respond to the request immediately since its no-reply
		writeOk(rw, nil)

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			logrus.WithError(err).Errorf("failed to read body from github webhook")
			return
		}

		var r model.GithubWebhookRequest
		if err = json.Unmarshal(body, &r); err != nil {
			logrus.WithError(err).Errorf("failed to marshal github webhook")
			return
		}

		if err := githubWebhookService.HandleRequest(&r); err != nil {
			logrus.WithError(err).Errorf("failed to handle github webhook")
		}
	}
}
