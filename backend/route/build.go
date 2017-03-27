package route

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dpolansky/ci/backend/service"
	"github.com/gorilla/mux"
)

const pathTokenBuildID = "build_id"
const pathTokenRepositoryName = "repository_name"
const pathURLBuildAPI = "/api/build"

func RegisterBuildAPIRoutes(router *mux.Router, build service.BuildReader) {
	router.HandleFunc(pathURLBuildAPI+"/{"+pathTokenBuildID+"}",
		getBuildStatusHTTPHandler(build)).Methods("GET")

	router.HandleFunc(pathURLRepositoryAPI+"/{"+pathTokenOwner+"}"+"/{"+pathTokenRepositoryName+"}/builds",
		getBuildsBySourceRepositoryHTTPHandler(build)).Methods("GET")
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

func getBuildsBySourceRepositoryHTTPHandler(build service.BuildReader) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		owner, found := vars[pathTokenOwner]
		if !found {
			writeError(rw, http.StatusBadRequest, fmt.Errorf("No owner found in path"))
			return
		}

		name, found := vars[pathTokenRepositoryName]
		if !found {
			writeError(rw, http.StatusBadRequest, fmt.Errorf("No repository name found in path"))
			return
		}

		asGithubURL := fmt.Sprintf("https://github.com/%s/%s", owner, name)
		builds, err := build.GetBuildsBySourceRepositoryURL(asGithubURL)
		if err != nil {
			writeError(rw, http.StatusNotFound, err)
			return
		}

		b, _ := json.Marshal(builds)
		writeOk(rw, b)
	}
}
