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
	ID         int                 `json:"id"`
	Source     *RepositoryMetadata `json:"source"`
	Tested     bool                `json:"tested"`
	Test       *RepositoryMetadata `json:"test"`
	LastUpdate time.Time           `json:"last_update"`
	Status     string              `json:"status"`
	Log        string              `json:"log"`
}

type RepositoryMetadata struct {
	ID       int    `json:"id"`
	CloneURL string `json:"clone_url"`
	Branch   string `json:"branch"`
}
