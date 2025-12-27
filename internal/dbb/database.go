package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/kroyser123/ParserGW/internal/models"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	connStr := "host=localhost port=5432 user=user password=password dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("database connection error: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("database ping error: %v", err)
	}

	log.Println("Connected to PostgreSQL successfully")
	return db, nil
}

func SaveConfig(db *sql.DB, config *models.Config) error {
	query := `
        INSERT INTO configs (name, description, version, author, tags)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (name) 
        DO UPDATE SET 
            description = EXCLUDED.description,
            version = EXCLUDED.version,
            author = EXCLUDED.author,
            tags = EXCLUDED.tags
    `

	_, err := db.Exec(query,
		config.Name,
		config.Description,
		config.Version,
		config.Author,
		pq.Array(config.Tags),
	)

	if err != nil {
		return fmt.Errorf("save config error: %v", err)
	}

	return nil
}
