package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type JsonPathsConfig struct {
	Paths map[string]string `yaml:"paths"`
}

func LoadJsonPaths(path string) (*JsonPathsConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg JsonPathsConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
