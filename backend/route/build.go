package route

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/dpolansky/ci/backend/service"
	humanize "github.com/dustin/go-humanize"
	"github.com/gorilla/mux"
)

func RegisterBuildFrontendRoutes(router *mux.Router, build service.BuildReader) {
	router.HandleFunc("/{"+pathTokenRepositoryID+"}"+pathURLBuildsByRepositoryID,
		getBuildsByRepositoryIDTemplateHTTPHandler(build)).Methods("GET")
}

func RegisterBuildAPIRoutes(router *mux.Router, build service.BuildReader) {
	router.HandleFunc(pathURLBuildAPI+"/{"+pathTokenBuildID+"}",
		getBuildStatusHTTPHandler(build)).Methods("GET")

	router.HandleFunc(pathURLRepositoryAPI+"/{"+pathTokenRepositoryID+"}/builds",
		getBuildsByRepositoryIDHTTPHandler(build)).Methods("GET")
}

func getBuildsByRepositoryIDTemplateHTTPHandler(build service.BuildReader) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id, found := vars[pathTokenRepositoryID]
		if !found {
			writeError(rw, http.StatusBadRequest, fmt.Errorf("No repository ID found in path"))
			return
		}

		asInt, err := strconv.Atoi(id)
		if err != nil {
			writeError(rw, http.StatusBadRequest, fmt.Errorf("Repository ID is not a number: %v", id))
			return
		}

		builds, err := build.GetBuildsByRepositoryID(asInt)
		if err != nil {
			writeError(rw, http.StatusInternalServerError, fmt.Errorf("Failed to get builds for repository %v: %v", asInt, err))
			return
		}

		// custom struct to display human friendly information
		type BuildStatus struct {
			ID         int
			Branch     string
			Status     string
			LastUpdate string
		}

		result := []*BuildStatus{}

		// humanize times
		for _, b := range builds {
			result = append(result, &BuildStatus{
				ID:         b.ID,
				Branch:     b.Branch,
				Status:     b.Status,
				LastUpdate: humanize.Time(b.LastUpdate),
			})
		}

		tempPath := filepath.Join(os.Getenv("GOPATH"), templatesDirPathFromGOPATH, "builds.html")
		tmpl := template.Must(template.ParseFiles(tempPath))
		tmpl.Execute(rw, struct{ Builds []*BuildStatus }{result})
	}
}

func getBuildStatusHTTPHandler(build service.BuildReader) func(rw http.ResponseWriter, req *http.Request) {
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

		status, err := build.GetBuildByID(asInt)
		if err != nil {
			writeError(rw, http.StatusNotFound, err)
			return
		}

		b, _ := json.Marshal(status)
		writeOk(rw, b)
	}
}

func getBuildsByRepositoryIDHTTPHandler(build service.BuildReader) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		id, found := vars[pathTokenRepositoryID]
		if !found {
			writeError(rw, http.StatusBadRequest, fmt.Errorf("No repository ID found in path"))
			return
		}

		asInt, err := strconv.Atoi(id)
		if err != nil {
			writeError(rw, http.StatusBadRequest, fmt.Errorf("Repository ID is not a number: %v", id))
			return
		}

		builds, err := build.GetBuildsByRepositoryID(asInt)
		if err != nil {
			writeError(rw, http.StatusNotFound, err)
			return
		}

		b, _ := json.Marshal(builds)
		writeOk(rw, b)
	}
}
