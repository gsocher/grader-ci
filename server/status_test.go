package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetStatusNotFound(t *testing.T) {
	buildService := NewBuildService(nil)
	handler := getBuildStatusHTTPHandler(buildService)
	router := mux.NewRouter()
	router.HandleFunc("status/{"+PathTokenBuildID+"}", handler).Methods("GET")

	ts := httptest.NewServer(router)
	r, err := http.Get(ts.URL + "/status/foo")
	if err != nil {
		t.Fatalf("unexpected http client err: %v", err)
	}

	if r.StatusCode != http.StatusNotFound {
		t.Fatalf("expected %v got %v", http.StatusNotFound, r.StatusCode)
	}
}
