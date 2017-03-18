package route

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dpolansky/ci/model"
	"github.com/dpolansky/ci/server/service"
	"github.com/gorilla/mux"
)

const pathTokenRepositoryID = "repository_id"
const pathTokenOwner = "owner"
const pathURLRepositoryAPI = "/api/repository"

func RegisterRepositoryRoutes(router *mux.Router, repositoryService service.RepositoryService) {
	router.HandleFunc(pathURLRepositoryAPI, createRepositoryHTTPHandler(repositoryService)).Methods("POST")
}

func createRepositoryHTTPHandler(repositoryService service.RepositoryService) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			writeError(rw, http.StatusBadRequest, err)
			return
		}

		var m model.Repository
		if err = json.Unmarshal(body, &m); err != nil {
			writeError(rw, http.StatusBadRequest, err)
			return
		}

		if m.CloneURL == "" {
			writeError(rw, http.StatusBadRequest, fmt.Errorf("Missing cloneURL"))
			return
		}

		if m.Owner == "" {
			writeError(rw, http.StatusBadRequest, fmt.Errorf("Missing owner"))
			return
		}

		err = repositoryService.CreateRepository(&m)
		if err != nil {
			writeError(rw, http.StatusInternalServerError, fmt.Errorf("Failed to create repository: %v", err))
			return
		}

		b, _ := json.Marshal(m)
		writeOk(rw, b)
	}
}
