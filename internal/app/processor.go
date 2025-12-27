// Package app содержит логику обработки файлов.
package app

import (
	"HW4/internal/dbb"
	"HW4/internal/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ProcessDir обходит директорию и обрабатывает JSON-файлы по схемам
func ProcessDir(rootPath string, config *models.Config, schemas []models.Schema, db *sql.DB) error {
	return filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			fmt.Printf("📁 Входим в папку: %s\n", path)
			return nil
		}

		if strings.HasSuffix(d.Name(), ".json") {
			return processJSONFile(path, schemas, config, db)
		}

		fmt.Printf("📎 Пропускаем: %s (не JSON)\n", path)
		return nil
	})
}

// processJSONFile извлекает метаданные из JSON по каждой схеме и сохраняет в БД
func processJSONFile(filePath string, schemas []models.Schema, config *models.Config, db *sql.DB) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("не удалось прочитать файл: %w", err)
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("некорректный JSON в %s: %w", filePath, err)
	}

	fmt.Printf("🔍 Обработка файла: %s\n", filePath)

	for i, schema := range schemas {
		extracted := extractMetadata(jsonData, schema, config.Tags)
		fmt.Printf("  Схема [%d]: %+v\n", i+1, extracted)

		err = dbb.UpsertFileMetadata(db, filePath, extracted)
		if err != nil {
			return fmt.Errorf("ошибка сохранения метаданных: %w", err)
		}
	}

	return nil
}

// extractMetadata извлекает данные по путям из Schema
func extractMetadata(data map[string]interface{}, schema models.Schema, configTags []string) map[string]interface{} {
	getString := func(path string) string {
		if val, ok := getValue(data, path); ok {
			if s, _ := val.(string); s != "" {
				return s
			}
		}
		return ""
	}

	getStringSlice := func(path string) []string {
		if val, ok := getValue(data, path); ok {
			switch v := val.(type) {
			case []interface{}:
				var res []string
				for _, item := range v {
					if s, ok := item.(string); ok {
						res = append(res, s)
					}
				}
				return res
			case string:
				return strings.Split(v, ",")
			}
		}
		return nil
	}

	name := getString(schema.Name)
	description := getString(schema.Description)
	version := getString(schema.Version)
	author := getString(schema.Author)
	tags := getStringSlice(schema.Tags)

	// Проверяем, есть ли совпадение тегов
	hasMatchingTags := false
	if len(tags) > 0 && len(configTags) > 0 {
		for _, ct := range configTags {
			for _, et := range tags {
				if ct == et {
					hasMatchingTags = true
					break
				}
			}
			if hasMatchingTags {
				break
			}
		}
	}

	return map[string]interface{}{
		"file_path":         "",
		"name":              name,
		"description":       description,
		"version":           version,
		"author":            author,
		"extracted_tags":    tags,
		"has_matching_tags": hasMatchingTags,
		"schema_used":       schema,
	}
}

// getValue — копия из config.go (можно вынести в утилиту)
func getValue(data map[string]interface{}, path string) (interface{}, bool) {
	keys := strings.Split(path, ".")
	current := data

	for _, key := range keys[:len(keys)-1] {
		val, exists := current[key]
		if !exists {
			return nil, false
		}
		if next, ok := val.(map[string]interface{}); ok {
			current = next
		} else {
			return nil, false
		}
	}

	lastKey := keys[len(keys)-1]
	val, exists := current[lastKey]
	return val, exists
}
