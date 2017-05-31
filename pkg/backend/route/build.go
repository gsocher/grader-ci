package route

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"strings"

	"github.com/dpolansky/grader-ci/pkg/backend/service"
	"github.com/dpolansky/grader-ci/pkg/model"
	humanize "github.com/dustin/go-humanize"
	"github.com/gorilla/mux"
)

func RegisterBuildFrontendRoutes(router *mux.Router, build service.BuildService, rep service.RepositoryService) {
	router.HandleFunc("/{"+pathTokenRepositoryID+"}"+pathURLBuildsByRepositoryID,
		getBuildsByRepositoryIDTemplateHTTPHandler(build, rep)).Methods("GET")

	router.HandleFunc("/{"+pathTokenRepositoryID+"}"+pathURLBuildsByRepositoryID+"/{"+pathTokenBuildID+"}",
		getBuildByIDTemplateHTTPHandler(build, rep)).Methods("GET")
}

func RegisterBuildAPIRoutes(router *mux.Router, build service.BuildService) {
	router.HandleFunc(pathURLBuildAPI+"/{"+pathTokenBuildID+"}",
		getBuildStatusHTTPHandler(build)).Methods("GET")

	router.HandleFunc(pathURLRepositoryAPI+"/{"+pathTokenRepositoryID+"}/builds",
		getBuildsByRepositoryIDHTTPHandler(build)).Methods("GET")
}

func getBuildsByRepositoryIDTemplateHTTPHandler(build service.BuildService, rep service.RepositoryService) func(rw http.ResponseWriter, req *http.Request) {
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

		repo, err := rep.GetRepositoryByID(asInt)
		if err != nil {
			writeError(rw, http.StatusBadRequest, fmt.Errorf("Repository with ID %v could not be found: %v", asInt, err))
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
			TestID     string
		}

		result := []*BuildStatus{}

		// humanize times
		for _, b := range builds {
			var testID string
			if b.Tested {
				testID = strconv.Itoa(b.Test.ID)
			} else {
				testID = "none"
			}

			result = append(result, &BuildStatus{
				ID:         b.ID,
				Branch:     b.Source.Branch,
				Status:     b.Status,
				LastUpdate: humanize.Time(b.LastUpdate),
				TestID:     testID,
			})
		}

		tempPath := filepath.Join(os.Getenv("GOPATH"), templatesDirPathFromGOPATH, "builds.html")
		tmpl := template.Must(template.ParseFiles(tempPath))
		tmpl.Execute(rw, struct {
			Builds     []*BuildStatus
			Repository *model.Repository
		}{result, repo})
	}
}

func getBuildByIDTemplateHTTPHandler(build service.BuildService, rep service.RepositoryService) func(rw http.ResponseWriter, req *http.Request) {
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

		repo, err := rep.GetRepositoryByID(asInt)
		if err != nil {
			writeError(rw, http.StatusBadRequest, fmt.Errorf("Repository with ID %v could not be found: %v", asInt, err))
			return
		}

		id, found = vars[pathTokenBuildID]
		if !found {
			writeError(rw, http.StatusBadRequest, fmt.Errorf("No build ID found in path"))
			return
		}

		asInt, err = strconv.Atoi(id)
		if err != nil {
			writeError(rw, http.StatusBadRequest, fmt.Errorf("build ID is not a number: %v", id))
			return
		}

		b, err := build.GetBuildByID(asInt)
		if err != nil {
			writeError(rw, http.StatusInternalServerError, fmt.Errorf("Failed to get build for repository %v: %v", asInt, err))
			return
		}

		// custom struct to display human friendly information
		type BuildStatus struct {
			ID         int
			Branch     string
			Status     string
			LastUpdate string
			Log        []string
			Tested     bool
			TestBranch string
		}

		// humanize times
		status := &BuildStatus{
			ID:         b.ID,
			Branch:     b.Source.Branch,
			Status:     b.Status,
			LastUpdate: humanize.Time(b.LastUpdate),
			Log:        strings.Split(b.Log, "\n"),
			Tested:     b.Tested,
			TestBranch: b.Test.Branch,
		}

		var testRepo *model.Repository
		if b.Tested {
			testRepo, err = rep.GetRepositoryByID(b.Test.ID)
			if err != nil {
				writeError(rw, http.StatusBadRequest, fmt.Errorf("Repository with ID %v could not be found: %v", b.Test.ID, err))
				return
			}
		}

		tempPath := filepath.Join(os.Getenv("GOPATH"), templatesDirPathFromGOPATH, "detail.html")
		tmpl := template.Must(template.ParseFiles(tempPath))
		tmpl.Execute(rw, struct {
			Build  *BuildStatus
			Source *model.Repository
			Test   *model.Repository
		}{status, repo, testRepo})
	}
}

func getBuildStatusHTTPHandler(build service.BuildService) func(rw http.ResponseWriter, req *http.Request) {
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

func getBuildsByRepositoryIDHTTPHandler(build service.BuildService) func(rw http.ResponseWriter, req *http.Request) {
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
