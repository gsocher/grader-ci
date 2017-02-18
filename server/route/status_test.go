package route

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dpolansky/ci/server/service"
	"github.com/gorilla/mux"
)

func TestGetStatus(t *testing.T) {
	router := mux.NewRouter()
	amqpClient := service.NewMockClient()
	builder := service.NewBuilder(amqpClient)

	// add the route to the router
	RegisterBuildStatusRoutes(router, builder)

	ts := httptest.NewServer(router)
	defer ts.Close()

	// add a status to the build service
	status, err := builder.StartBuild("github.com/docker/docker")
	if err != nil {
		t.Fatalf("Failed to start build: %v", err)
	}

	// check if the status exists
	r, err := http.Get(ts.URL + "/status/" + status.ID)
	if err != nil {
		t.Fatalf("unexpected http client err: %v", err)
	}

	if r.StatusCode != http.StatusOK {
		t.Fatalf("expected %v got %v", http.StatusOK, r.StatusCode)
	}
}

func TestGetStatusNotFound(t *testing.T) {
	router := mux.NewRouter()
	amqpClient := service.NewMockClient()
	builder := service.NewBuilder(amqpClient)

	// add the route to the router
	RegisterBuildStatusRoutes(router, builder)

	ts := httptest.NewServer(router)
	defer ts.Close()

	r, err := http.Get(ts.URL + "/status/foo")
	if err != nil {
		t.Fatalf("unexpected http client err: %v", err)
	}

	if r.StatusCode != http.StatusNotFound {
		t.Fatalf("expected %v got %v", http.StatusNotFound, r.StatusCode)
	}
}
