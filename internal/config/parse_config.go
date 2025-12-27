// Package config отвечает за загрузку и парсинг конфигурационных файлов.
package config

import (
	"HW4/internal/models"
	"errors"
	"fmt"
	"os"

	"go.yaml.in/yaml/v3"
)

// getValue извлекает значение из вложенной map по пути вида "a.b.c"

// ParseConfig загружает и парсит основной конфиг (например, app.yaml)
func ParseConfig(configPath string) (*models.Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл конфига: %w", err)
	}

	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("ошибка парсинга YAML: %w", err)
	}

	configData, exists := raw["config"]
	if !exists {
		return nil, errors.New("отсутствует секция 'config' в YAML")
	}
	configMap, ok := configData.(map[string]interface{})
	if !ok {
		return nil, errors.New("секция 'config' должна быть объектом")
	}

	// name
	name, ok := configMap["name"].(string)
	if !ok || name == "" {
		return nil, errors.New("config.name — обязательное строковое поле")
	}

	// description
	description, ok := configMap["description"].(string)
	if !ok || description == "" {
		return nil, errors.New("config.description — обязательное строковое поле")
	}

	// version
	versionVal, ok := configMap["version"]
	if !ok {
		return nil, errors.New("config.version — обязательное поле")
	}
	var version int
	switch v := versionVal.(type) {
	case int:
		version = v
	case float64:
		version = int(v)
	default:
		return nil, errors.New("config.version должно быть числом")
	}
	if version <= 0 {
		return nil, errors.New("config.version должно быть больше 0")
	}

	// author
	author, ok := configMap["author"].(string)
	if !ok || author == "" {
		return nil, errors.New("config.author — обязательное строковое поле")
	}

	// tags
	tagsVal, ok := configMap["tags"]
	if !ok {
		return nil, errors.New("config.tags — обязательное поле")
	}
	tagsSlice, ok := tagsVal.([]interface{})
	if !ok {
		return nil, errors.New("config.tags должно быть списком")
	}
	var tags []string
	for _, t := range tagsSlice {
		if tag, ok := t.(string); ok {
			tags = append(tags, tag)
		} else {
			return nil, errors.New("все теги должны быть строками")
		}
	}
	if len(tags) == 0 {
		return nil, errors.New("config.tags не может быть пустым")
	}

	return &models.Config{
		Name:        name,
		Description: description,
		Version:     version,
		Author:      author,
		Tags:        tags,
	}, nil
}

// ParseSchemas загружает правила из json_schemas.yaml
func ParseSchemas(schemasPath string) ([]models.Schema, error) {
	data, err := os.ReadFile(schemasPath)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл схем: %w", err)
	}

	var raw map[string][]models.Schema
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("ошибка парсинга json_schemas.yaml: %w", err)
	}

	schemas, ok := raw["json_schemas"]
	if !ok {
		return nil, errors.New("в файле схем отсутствует 'json_schemas'")
	}

	if len(schemas) == 0 {
		return nil, errors.New("json_schemas не может быть пустым")
	}

	return schemas, nil
}
