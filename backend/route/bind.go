package route

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"strconv"

	"github.com/dpolansky/ci/backend/service"
	"github.com/dpolansky/ci/model"
	"github.com/gorilla/mux"
)

func RegisterBindAPIRoutes(router *mux.Router, bind service.TestBindReadWriter) {
	router.HandleFunc(pathURLTestBindAPI,
		updateTestBindHTTPHandler(bind)).Methods("POST")

	router.HandleFunc(pathURLTestBindAPI,
		getTestBindsHTTPHandler(bind)).Methods("GET")

	router.HandleFunc(pathURLTestBindAPI+"/{"+pathTokenRepositoryID+"}",
		getTestBindBySourceRepositoryIDHTTPHandler(bind)).Methods("GET")
}

func updateTestBindHTTPHandler(bind service.TestBindReadWriter) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			writeError(rw, http.StatusBadRequest, err)
			return
		}

		var b model.TestBind
		if err = json.Unmarshal(body, &b); err != nil {
			writeError(rw, http.StatusBadRequest, err)
			return
		}

		if b.SourceID == 0 {
			writeError(rw, http.StatusBadRequest, fmt.Errorf("Missing source ID"))
			return
		}

		if b.TestID == 0 {
			writeError(rw, http.StatusBadRequest, fmt.Errorf("Missing test ID"))
			return
		}

		err = bind.UpdateTestBind(&b)
		if err != nil {
			writeError(rw, http.StatusInternalServerError, err)
			return
		}

		byt, _ := json.Marshal(b)
		writeOk(rw, byt)
	}
}

func getTestBindsHTTPHandler(bind service.TestBindReadWriter) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		binds, err := bind.GetTestBinds()
		if err != nil {
			writeError(rw, http.StatusInternalServerError, fmt.Errorf("Failed to get binds: %v", err))
			return
		}

		b, _ := json.Marshal(binds)
		writeOk(rw, b)
	}
}

func getTestBindBySourceRepositoryIDHTTPHandler(bind service.TestBindReadWriter) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		sourceID, found := vars[pathTokenRepositoryID]
		if !found {
			writeError(rw, http.StatusBadRequest, fmt.Errorf("No source repository ID found in path"))
			return
		}

		asInt, err := strconv.Atoi(sourceID)
		if err != nil {
			writeError(rw, http.StatusBadRequest, fmt.Errorf("Source repository ID is not a number: %v", sourceID))
			return
		}

		bind, err := bind.GetTestBindBySourceRepositoryID(asInt)
		if err != nil {
			writeError(rw, http.StatusInternalServerError, fmt.Errorf("Failed to get bind for source repository ID=%v: %v", asInt, err))
			return
		}

		b, _ := json.Marshal(bind)
		writeOk(rw, b)
	}
}
