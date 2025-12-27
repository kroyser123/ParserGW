// Package db отвечает за подключение к PostgreSQL и операции с БД.
package dbb

import (
	"database/sql"
	"log"

	"github.com/lib/pq"
)

// Connect подключается к PostgreSQL
func Connect(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	log.Println("✅ Подключено к PostgreSQL")
	return db, nil
}

// UpsertFileMetadata сохраняет извлечённые метаданные в таблицу
func UpsertFileMetadata(db *sql.DB, filePath string, meta map[string]interface{}) error {
	query := `
		INSERT INTO file_metadata (
			file_path, name, description, version, author, extracted_tags, has_matching_tags
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (file_path, name) DO UPDATE SET
			description = EXCLUDED.description,
			version = EXCLUDED.version,
			author = EXCLUDED.author,
			extracted_tags = EXCLUDED.extracted_tags,
			has_matching_tags = EXCLUDED.has_matching_tags
	`

	_, err := db.Exec(
		query,
		filePath,
		meta["name"],
		meta["description"],
		meta["version"],
		meta["author"],
		pq.Array(meta["extracted_tags"].([]string)),
		meta["has_matching_tags"].(bool),
	)
	return err
}
