package worker

import (
	"bytes"
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestRunBuild(t *testing.T) {
	w, err := NewWorker()
	if err != nil {
		t.Fatalf("Failed to create worker, err=%v", err)
	}

	task := &BuildTask{
		Language: "go",
		CloneURL: "github.com/dpolansky/go-poet",
		Ctx:      context.Background(),
		ID:       uuid.New().String(),
	}

	output := bytes.NewBuffer(nil)
	exitCode, err := w.RunBuild(task, output)
	if err != nil {
		t.Fatalf("Build failed, err=%v", err)
	}

	if exitCode != 0 {
		t.Fatalf("Running build script failed, output=%v", output)
	}
}
