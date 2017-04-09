package route

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dpolansky/ci/backend/service"
	"github.com/dpolansky/ci/model"
	"github.com/gorilla/mux"
)

func RegisterRepositoryFrontendRoutes(router *mux.Router, rep service.RepositoryReadWriter) {
	router.HandleFunc(pathURLRepositoryList,
		getRepositoryListTemplateHTTPHandler(rep)).Methods("GET")
}

func RegisterRepositoryAPIRoutes(router *mux.Router, rep service.RepositoryReadWriter) {
	router.HandleFunc(pathURLRepositoryAPI,
		getRepositoriesHTTPHandler(rep)).Methods("GET")

	router.HandleFunc(pathURLRepositoryAPI+"/{"+pathTokenOwner+"}",
		getRepositoriesByOwnerHTTPHandler(rep)).Methods("GET")
}

func getRepositoryListTemplateHTTPHandler(rep service.RepositoryReadWriter) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		reps, err := rep.GetRepositories()
		if err != nil {
			writeError(rw, http.StatusInternalServerError, fmt.Errorf("Failed to get repositories: %v", err))
			return
		}

		tempPath := filepath.Join(os.Getenv("GOPATH"), templatesDirPathFromGOPATH, "repositories.html")
		tmpl := template.Must(template.ParseFiles(tempPath))
		tmpl.Execute(rw, struct{ Repositories []*model.Repository }{reps})
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
