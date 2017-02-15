package worker

import "testing"

func TestParse(t *testing.T) {
	data := "language: ruby"

	cfg, err := parse([]byte(data))
	if err != nil {
		t.Fatalf("error parsing yaml: %v", err)
	}

	if cfg.Language != "ruby" {
		t.Fatalf("expected %v, got %v", "ruby", cfg.Language)
	}
}
