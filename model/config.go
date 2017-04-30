package model

const ConfigFileName = ".ci.yml"

// Represents the yaml configuration file found in a build repository.
type Config struct {
	Language string   `yaml:"language"`
	Script   []string `yaml:"script"`
}
