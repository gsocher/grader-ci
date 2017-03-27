package model

import "time"

const (
	StatusBuildWaiting = "waiting"
	StatusBuildRunning = "running"
	StatusBuildFailed  = "failed"
	StatusBuildError   = "error"
	StatusBuildPassed  = "passed"
)

type BuildStatus struct {
	ID           int       `json:"id"`
	RepositoryID int       `json:"repo_id"`
	LastUpdate   time.Time `json:"last_update"`
	CloneURL     string    `json:"clone_url"`
	Branch       string    `json:"branch"`
	Status       string    `json:"status"`
	Log          string    `json:"log"`
}
