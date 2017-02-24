package worker

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dpolansky/ci/model"
)

func TestRunBuild(t *testing.T) {
	w, err := New()
	if err != nil {
		t.Fatalf("Failed to create worker, err=%v", err)
	}

	gopath := os.Getenv("GOPATH")
	repoPath := filepath.Join(gopath, "src/github.com/dpolansky/ci/worker/test/golang")

	task := &model.BuildStatus{
		CloneURL: repoPath,
		ID:       1,
	}

	err = w.RunBuild(task, os.Stdout)
	if err != nil {
		t.Fatalf("Build failed, err=%v", err)
	}
}

func TestRunBuildFail(t *testing.T) {
	w, err := New()
	if err != nil {
		t.Fatalf("Failed to create worker, err=%v", err)
	}

	gopath := os.Getenv("GOPATH")
	repoPath := filepath.Join(gopath, "src/github.com/dpolansky/ci/worker/test/golang-fail")

	task := &model.BuildStatus{
		CloneURL: repoPath,
		ID:       1,
	}

	err = w.RunBuild(task, os.Stdout)
	if err == nil {
		t.Fatalf("Expected build to fail")
	}
}

func TestGetImageForLanguage(t *testing.T) {
	var tests = []struct {
		language    string
		expectedImg string
	}{
		{"golang", "build-golang"},
		{"test", "build-test"},
	}

	for _, test := range tests {
		img, err := getImageForLanguage(test.language)
		if err != nil {
			t.Fatalf("error getting image for language %v: %v", test.language, err)
		}

		if img != test.expectedImg {
			t.Fatalf("Expected %v, got %v", test.expectedImg, img)
		}
	}
}
