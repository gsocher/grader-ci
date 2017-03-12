package route

import (
	"html/template"
	"net/http"
	"path/filepath"

	"os"

	"github.com/dpolansky/ci/model"
	"github.com/dpolansky/ci/server/service"
	"github.com/gorilla/mux"
)

const pathURLStatusListFrontend = "/status"
const templatesDirPathFromGOPATH = "/src/github.com/dpolansky/ci/server/templates/status_list.html"

func RegisterBuildStatusFrontendRoutes(router *mux.Router, builder service.Builder) {
	router.HandleFunc(pathURLStatusListFrontend, getBuildStatusListTemplateHTTPHandler(builder)).Methods("GET")
}

func getBuildStatusListTemplateHTTPHandler(builder service.Builder) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		statuses, err := builder.GetBuildStatuses()
		if err != nil {
			writeError(rw, http.StatusInternalServerError, err)
			return
		}

		templatesDir := filepath.Join(os.Getenv("GOPATH"), templatesDirPathFromGOPATH)
		tmpl := template.Must(template.ParseFiles(templatesDir))
		tmpl.Execute(rw, struct{ Statuses []*model.BuildStatus }{statuses})
	}
}
