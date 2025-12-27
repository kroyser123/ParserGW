package app

import (
	"database/sql"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"github.com/kroyser123/ParserGW/internal/config"
	"github.com/kroyser123/ParserGW/internal/dbb"
	"github.com/kroyser123/ParserGW/internal/models"
)

type Stats struct {
	FilesProcessed int
	RecordsUpdated int
}

func ProcessDirectory(dbConn *sql.DB, dirPath string, jsonPaths []config.JSONPathConfig) (*Stats, error) {
	stats := &Stats{}

	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("Error accessing path %s: %v", path, err)
			return nil
		}

		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".yaml" && ext != ".yml" && ext != ".json" {
			return nil
		}

		log.Printf("Processing file: %s", path)
		stats.FilesProcessed++

		var cfg models.Config
		var parseErr error

		if ext == ".yaml" || ext == ".yml" {
			cfg, parseErr = config.ParseYAMLFile(path)
		} else {
			cfg, parseErr = config.ParseJSONFile(path, jsonPaths)
		}

		if parseErr != nil {
			log.Printf("Warning: Failed to parse %s: %v", path, parseErr)
			return nil
		}

		// Проверяем обязательное поле Name
		if cfg.Name == "" {
			log.Printf("Warning: File %s has empty 'name' field, skipping", path)
			return nil
		}

		if err := db.SaveConfig(dbConn, &cfg); err != nil {
			log.Printf("Warning: Failed to save config from %s: %v", path, err)
			return nil
		}

		stats.RecordsUpdated++
		return nil
	})

	return stats, err
}
