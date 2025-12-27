package models

// Config - общая модель файлов (конфигов)
type Config struct {
	Name        string
	Description string
	Version     int
	Author      string
	Tags        []string
}

// Schema - модель описания JSON схемы (поле->путь)
type Schema struct {
	Name        string
	Description string
	Version     string
	Author      string
	Tags        string
}
