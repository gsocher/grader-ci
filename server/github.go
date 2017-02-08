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

func parseWebhookHTTPHandler(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
		return
	}

	var r githubWebhookRequest
	if err = json.Unmarshal(body, &r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
		return
	}

	b, err := json.Marshal(r)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(b)
}
