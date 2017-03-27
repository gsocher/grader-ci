package route

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dpolansky/ci/backend/service"
	"github.com/dpolansky/ci/model"
	"github.com/gorilla/mux"
)

const pathTokenOwner = "owner_name"
const pathURLRepositoryAPI = "/api/repository"

func RegisterRepositoryRoutes(router *mux.Router, rep service.RepositoryReadWriter) {
	router.HandleFunc(pathURLRepositoryAPI,
		createRepositoryHTTPHandler(rep)).Methods("POST")

	router.HandleFunc(pathURLRepositoryAPI,
		getRepositoriesHTTPHandler(rep)).Methods("GET")

	router.HandleFunc(pathURLRepositoryAPI+"/{"+pathTokenOwner+"}",
		getRepositoriesByOwnerHTTPHandler(rep)).Methods("GET")
}

func createRepositoryHTTPHandler(rep service.RepositoryReadWriter) func(rw http.ResponseWriter, req *http.Request) {
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

		err = rep.UpdateRepository(&m)
		if err != nil {
			writeError(rw, http.StatusInternalServerError, fmt.Errorf("Failed to create repository: %v", err))
			return
		}

		b, _ := json.Marshal(m)
		writeOk(rw, b)
	}
}

func getRepositoriesHTTPHandler(rep service.RepositoryReadWriter) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		reps, err := rep.GetRepositories()
		if err != nil {
			writeError(rw, http.StatusInternalServerError, fmt.Errorf("Failed to get repositories: %v", err))
			return
		}

		b, _ := json.Marshal(reps)
		writeOk(rw, b)
	}
}

func getRepositoriesByOwnerHTTPHandler(rep service.RepositoryReadWriter) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		owner, found := vars[pathTokenOwner]
		if !found {
			writeError(rw, http.StatusBadRequest, fmt.Errorf("Missing owner name in path"))
			return
		}

		reps, err := rep.GetRepositoriesByOwner(owner)
		if err != nil {
			writeError(rw, http.StatusInternalServerError, fmt.Errorf("Failed to get repositories: %v", err))
			return
		}

		b, _ := json.Marshal(reps)
		writeOk(rw, b)
	}
}
