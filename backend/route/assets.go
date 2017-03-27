package route

import (
	"os"
	"path/filepath"

	"net/http"

	"github.com/gorilla/mux"
)

const assetsDirRelToGOPATH = "/src/github.com/dpolansky/ci/backend/static/assets/"

func RegisterAssetsRoute(r *mux.Router) {
	r.HandleFunc(pathURLAssets+"/{"+pathTokenFileName+"}", serveAssetsHTTPHandler()).Methods("GET")
}

func serveAssetsHTTPHandler() func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		fileName := vars[pathTokenFileName]

		gopath := os.Getenv("GOPATH")
		path := filepath.Join(gopath, assetsDirRelToGOPATH, fileName)
		http.ServeFile(rw, req, path)
	}
}
