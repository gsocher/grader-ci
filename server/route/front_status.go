package route

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/dpolansky/ci/model"
	"github.com/dpolansky/ci/server/service"
	"github.com/gorilla/mux"
)

const pathURLBuildListBase = "/repository"
const pathURLBuildListSuffix = "/builds"
const templatesDirPathFromGOPATH = "/src/github.com/dpolansky/ci/server/templates"
const statusListTemplatePath = "status_list.html"
const statusDetailTemplatePath = "status_detail.html"

func RegisterBuildStatusFrontendRoutes(router *mux.Router, builder service.Builder) {
	router.HandleFunc(pathURLBuildListSuffix,
		getBuildStatusListTemplateHTTPHandler(builder)).Methods("GET")

	router.HandleFunc(pathURLBuildListSuffix+"/{"+pathTokenBuildID+"}",
		getBuildStatusDetailTemplateHTTPHandler(builder)).Methods("GET")
}

func getBuildStatusListTemplateHTTPHandler(builder service.Builder) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		statuses, err := builder.GetBuilds()
		if err != nil {
			writeError(rw, http.StatusInternalServerError, err)
			return
		}

		tempPath := filepath.Join(os.Getenv("GOPATH"), templatesDirPathFromGOPATH, statusListTemplatePath)
		tmpl := template.Must(template.ParseFiles(tempPath))
		tmpl.Execute(rw, struct{ Statuses []*model.BuildStatus }{statuses})
	}
}

func getBuildStatusDetailTemplateHTTPHandler(builder service.Builder) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id, found := vars[pathTokenBuildID]
		if !found {
			writeError(rw, http.StatusBadRequest, fmt.Errorf("No build ID specified"))
			return
		}

		asInt, err := strconv.Atoi(id)
		if err != nil {
			writeError(rw, http.StatusBadRequest, fmt.Errorf("Build ID should be a number"))
			return
		}

		status, err := builder.GetBuildByID(asInt)
		if err != nil {
			writeError(rw, http.StatusInternalServerError, err)
			return
		}

		tempPath := filepath.Join(os.Getenv("GOPATH"), templatesDirPathFromGOPATH, statusDetailTemplatePath)
		tmpl := template.Must(template.ParseFiles(tempPath))
		tmpl.Execute(rw, status)
	}
}
