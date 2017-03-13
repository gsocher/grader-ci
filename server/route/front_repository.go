package route

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dpolansky/ci/model"
	"github.com/dpolansky/ci/server/service"
	"github.com/gorilla/mux"
)

const pathURLRepositoryFrontend = "/"
const repositoryListTemplatePath = "repository_list.html"

func RegisterRepositoryFrontendRoutes(router *mux.Router, repositoryService service.RepositoryService) {
	router.HandleFunc(pathURLRepositoryFrontend, getRepositoryListTemplateHTTPHandler(repositoryService)).Methods("GET")
}

func getRepositoryListTemplateHTTPHandler(repositoryService service.RepositoryService) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		repositories, err := repositoryService.GetRepositoriesByOwner("admin")
		if err != nil {
			writeError(rw, http.StatusInternalServerError, err)
			return
		}

		tempPath := filepath.Join(os.Getenv("GOPATH"), templatesDirPathFromGOPATH, repositoryListTemplatePath)
		tmpl := template.Must(template.ParseFiles(tempPath))
		tmpl.Execute(rw, struct{ Repositories []*model.Repository }{repositories})
	}
}
