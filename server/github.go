package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (s *Server) registerGithubWebhookRoutes() {
	s.Router.HandleFunc("/github", parseWebhookHTTPHandler).Methods("POST")
}

type githubWebhookRequest struct {
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

// parseWebhookHTTPHandler is an endpoint for receiving github webhook requests.
func parseWebhookHTTPHandler(rw http.ResponseWriter, req *http.Request) {
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

	b, err := json.Marshal(r)
	if err != nil {
		writeError(rw, http.StatusInternalServerError, err)
		return
	}

	writeOk(rw, b)
}
