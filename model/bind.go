package model

type TestBind struct {
	SourceID   int    `json:"source_id"`
	TestID     int    `json:"test_id"`
	TestBranch string `json:"test_branch"`
}
