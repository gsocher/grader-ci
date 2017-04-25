package worker

import (
	"github.com/dpolansky/grader-ci/model"

	"gopkg.in/yaml.v2"
)

func parse(b []byte) (*model.Config, error) {
	var cfg model.Config
	err := yaml.Unmarshal(b, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
