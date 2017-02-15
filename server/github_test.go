package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dpolansky/ci/server/service"
)

func TestGithubWebhook(t *testing.T) {
	amqpClient := service.NewMockClient()
	builder := service.NewBuilder(amqpClient)
	s, err := New(builder)
	if err != nil {
		t.Fatalf("Failed to initialize server: %v", err)
	}

	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	// create a test webhook payload
	var payload githubWebhookRequest
	payload.Repository.FullName = "docker/docker"
	b, _ := json.Marshal(payload)
	fmt.Println(string(b))
	buf := bytes.NewBuffer(b)

	// execute post
	r, err := http.Post(ts.URL+"/github", "application/json", buf)
	if err != nil {
		t.Fatalf("unexpected http client err: %v", err)
	}

	// check error code
	if r.StatusCode != http.StatusOK {
		t.Fatalf("expected %v got %v", http.StatusOK, r.StatusCode)
	}
}
