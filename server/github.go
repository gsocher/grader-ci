package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dpolansky/ci/server/amqp"
	"github.com/dpolansky/ci/worker"
	"github.com/sirupsen/logrus"
)

func (s *Server) registerGithubWebhookRoutes() {
	s.Router.HandleFunc("/github", parseWebhookHTTPHandler(s.amqpClient)).Methods("POST")
}

type githubWebhookRequest struct {
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

// parseWebhookHTTPHandler is an endpoint for receiving github webhook requests.
func parseWebhookHTTPHandler(amqpWriter amqp.Writer) func(rw http.ResponseWriter, req *http.Request) {
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

		// create a new build task
		task := &worker.BuildTask{
			CloneURL: fmt.Sprintf("github.com/%v", r.Repository.FullName),
			Language: "golang",
		}
		fmt.Printf("%+v\n", task)

		b, err := json.Marshal(task)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to marshal build task")
			writeError(rw, http.StatusInternalServerError, fmt.Errorf("Failed to start build"))
		}

		err = amqpWriter.SendToQueue("jobs", b)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to send build to queue")
			writeError(rw, http.StatusInternalServerError, fmt.Errorf("Failed to start build"))
		}

		writeOk(rw, []byte{})
	}
}
