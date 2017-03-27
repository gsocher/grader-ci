package route

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dpolansky/ci/backend/service"
	"github.com/gorilla/mux"
)

func RegisterBuildAPIRoutes(router *mux.Router, build service.BuildReader) {
	router.HandleFunc(pathURLBuildAPI+"/{"+pathTokenBuildID+"}",
		getBuildStatusHTTPHandler(build)).Methods("GET")

	router.HandleFunc(pathURLRepositoryAPI+"/{"+pathTokenRepositoryID+"}/builds",
		getBuildsByRepositoryIDHTTPHandler(build)).Methods("GET")
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
