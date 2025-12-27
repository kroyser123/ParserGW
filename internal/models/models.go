package models

type Config struct {
	Name        string   `json:"name" yaml:"name"`
	Description string   `json:"description" yaml:"description"`
	Version     int      `json:"version" yaml:"version"`
	Author      string   `json:"author" yaml:"author"`
	Tags        []string `json:"tags" yaml:"tags"`
}
