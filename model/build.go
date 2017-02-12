package model

import "time"

const (
	StatusBuildWaiting = "waiting"
	StatusBuildRunning = "running"
	StatusBuildFailed  = "failed"
	StatusBuildPassed  = "passed"
)

type BuildStatus struct {
	ID         string    `json: "id"`
	LastUpdate time.Time `json: "lastUpdate"`
	Language   string    `json: "language"`
	CloneURL   string    `json: "cloneURL"`
	Status     string    `json: "status"`
}
