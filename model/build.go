package model

import "time"

const (
	StatusBuildWaiting = "waiting"
	StatusBuildRunning = "running"
	StatusBuildFailed  = "failed"
	StatusBuildPassed  = "passed"
)

type BuildStatus struct {
	ID         int       `json: "id"`
	LastUpdate time.Time `json: "lastUpdate"`
	CloneURL   string    `json: "cloneURL"`
	Branch     string    `json: "branch"`
	Status     string    `json: "status"`
	Log        string    `json: "log"`
}
