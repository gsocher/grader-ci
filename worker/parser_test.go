package worker

import (
	"reflect"
	"testing"
)

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

func TestParseInvalidYAML(t *testing.T) {
	data := "	language: ruby"

	_, err := parse([]byte(data))
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestParseWithBuildScript(t *testing.T) {
	data := `script:
 - echo foo
 - echo bar`

	cfg, err := parse([]byte(data))
	if err != nil {
		t.Fatalf("error parsing yaml: %v", err)
	}

	expected := []string{"echo foo", "echo bar"}

	if !reflect.DeepEqual(cfg.Script, expected) {
		t.Fatalf("expected %v, got %v", expected, cfg.Script)
	}
}
