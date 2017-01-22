package main

import (
	"context"
	"log"
	"os"

	"github.com/dpolansky/ci/worker"
	"github.com/google/uuid"
)

func main() {
	w, err := worker.NewWorker()
	if err != nil {
		log.Fatalf("Failed to create worker, err=%v", err)
	}

	task := &worker.BuildTask{
		Language: "go",
		CloneURL: "github.com/dpolansky/go-poet",
		Ctx:      context.Background(),
		ID:       uuid.New().String(),
	}

	w.RunBuild(task, os.Stdout)
}
