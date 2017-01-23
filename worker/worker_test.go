package worker

import (
	"context"
	"testing"

	"os"

	"github.com/google/uuid"
)

func TestRunBuild(t *testing.T) {
	w, err := NewWorker()
	if err != nil {
		t.Fatalf("Failed to create worker, err=%v", err)
	}

	task := &BuildTask{
		Language: "test",
		CloneURL: "github.com/dpolansky/go-poet",
		Ctx:      context.Background(),
		ID:       uuid.New().String(),
	}

	exitCode, err := w.RunBuild(task, os.Stdout)
	if err != nil {
		t.Fatalf("Build failed, err=%v", err)
	}

	if exitCode != 0 {
		t.Fatalf("Running build script failed")
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
