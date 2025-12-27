package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kroyser123/ParserGW/internal/app"
	"github.com/kroyser123/ParserGW/internal/config"
	db "github.com/kroyser123/ParserGW/internal/dbb"
)

func main() {
	configPath := flag.String("c", "./json_paths.yaml", "Path to JSON paths config")
	dirPath := flag.String("d", "", "Path to directory with config files")
	flag.Parse()

	if *dirPath == "" {
		log.Fatal("Directory path is required. Use -d flag")
	}

	
	if _, err := os.Stat(*dirPath); os.IsNotExist(err) {
		log.Fatalf("Directory %s does not exist", *dirPath)
	}

	
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	
	jsonPaths, err := config.LoadJSONPaths(*configPath)
	if err != nil {
		log.Fatalf("Failed to load JSON paths: %v", err)
	}

	
	stats, err := app.ProcessDirectory(database, *dirPath, jsonPaths)
	if err != nil {
		log.Fatalf("Failed to process directory: %v", err)
	}

	
	fmt.Printf("Processed files: %d\n", stats.FilesProcessed)
	fmt.Printf("Records updated: %d\n", stats.RecordsUpdated)
}
