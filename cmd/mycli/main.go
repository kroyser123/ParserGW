// cmd/mycli/main.go
package main

import (
	"HW4/internal/app"
	"HW4/internal/config"
	"HW4/internal/dbb"
	"flag"
	"log"
	"os"
)

var (
	configPath    string // -c: путь к app.yaml (основной конфиг)
	schemasPath   string // -s: путь к json_paths.yaml (схемы извлечения)
	directoryPath string // -d: путь к директории с JSON-файлами
)

func init() {
	flag.StringVar(&configPath, "c", "", "путь к основному конфигу (app.yaml)")
	flag.StringVar(&schemasPath, "s", "", "путь к файлу схем (json_paths.yaml)")
	flag.StringVar(&directoryPath, "d", "", "путь к директории с JSON-файлами для обработки")
}

func main() {
	flag.Parse()

	// Вывод текущей директории — для отладки путей
	wd, _ := os.Getwd()
	log.Printf("🔧 Текущая директория: %s", wd)

	// Проверка обязательных флагов
	if configPath == "" {
		log.Fatal("❌ Ошибка: флаг -c обязателен (укажите путь к app.yaml)")
	}
	if schemasPath == "" {
		log.Fatal("❌ Ошибка: флаг -s обязателен (укажите путь к json_paths.yaml)")
	}
	if directoryPath == "" {
		log.Fatal("❌ Ошибка: флаг -d обязателен (укажите путь к тестовой директории)")
	}

	// 1. Загрузить основной конфиг (app.yaml)
	cfg, err := config.ParseConfig(configPath)
	if err != nil {
		log.Fatalf("❌ Ошибка парсинга основного конфига (%s): %v", configPath, err)
	}
	log.Printf("✅ Конфиг загружен: %s (v%d), теги: %v", cfg.Name, cfg.Version, cfg.Tags)

	// 2. Загрузить схемы (json_paths.yaml)
	schemas, err := config.ParseSchemas(schemasPath)
	if err != nil {
		log.Fatalf("❌ Ошибка парсинга схем (%s): %v", schemasPath, err)
	}
	log.Printf("✅ Загружено схем: %d", len(schemas))

	// 3. Подключиться к PostgreSQL
	connStr := "user=user password=password dbname=postgres host=localhost port=5432 sslmode=disable"
	database, err := dbb.Connect(connStr)
	if err != nil {
		log.Fatalf("❌ Не удалось подключиться к БД: %v", err)
	}
	defer database.Close()

	// 4. Обработать JSON-файлы
	log.Printf("🔍 Начинаем обработку файлов в: %s", directoryPath)
	err = app.ProcessDir(directoryPath, cfg, schemas, database)
	if err != nil {
		log.Fatalf("❌ Ошибка при обработке директории: %v", err)
	}

	log.Println("✅ Готово: все файлы обработаны и сохранены в базе данных.")
}
