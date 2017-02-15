package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dpolansky/ci/server/service"
)

func TestGetStatus(t *testing.T) {
	amqpClient := service.NewMockClient()
	builder := service.NewBuilder(amqpClient)
	s, err := New(builder)
	if err != nil {
		t.Fatalf("Failed to initialize server: %v", err)
	}

	ts := httptest.NewServer(s.Router)
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
	amqpClient := service.NewMockClient()
	builder := service.NewBuilder(amqpClient)
	s, err := New(builder)
	if err != nil {
		t.Fatalf("Failed to initialize server: %v", err)
	}

	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	r, err := http.Get(ts.URL + "/status/foo")
	if err != nil {
		t.Fatalf("unexpected http client err: %v", err)
	}

	if r.StatusCode != http.StatusNotFound {
		t.Fatalf("expected %v got %v", http.StatusNotFound, r.StatusCode)
	}
}
