package route

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dpolansky/ci/server/service"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const pathTokenBuildID = "build_id"
const pathURLBuildStatusAPI = "/api/status/"

func RegisterBuildStatusRoutes(router *mux.Router, builder service.Builder) {
	router.HandleFunc(pathURLBuildStatusAPI+"{"+pathTokenBuildID+"}", getBuildStatusHTTPHandler(builder)).Methods("GET")
}

func getBuildStatusHTTPHandler(builder service.Builder) func(rw http.ResponseWriter, req *http.Request) {
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
			writeError(rw, http.StatusNotFound, err)
			return
		}

		bytes, err := json.Marshal(status)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to marshal build status")
			writeError(rw, http.StatusInternalServerError, err)
			return
		}

		writeOk(rw, bytes)
	}
}
