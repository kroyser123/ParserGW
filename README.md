# Homework 3

Домашнее задание: CLI-приложение для обработки конфигурационных файлов и записи данных в PostgreSQL

## 1. Описание задания

Необходимо разработать CLI-приложение, которое:
* Принимает путь к директории в качестве аргумента командной строки: `mycli -d ./tests/src/test_1`.
* Обходит все файлы в указанной директории и её поддиректориях (рекурсивно).
* Обрабатывает **YAML-файлы** (с единой строго заданной схемой).
* Обрабатывает **JSON-файлы** (с различными возможными структурами, но с заранее известными путями к нужным данным).
* Записывает или обновляет данные в PostgreSQL (если запись уже существует, обновлять вместо вставки).
* Выводит краткую статистику о количестве обработанных файлов и обновлённых/добавленных записей.

## 2. Структура конфигурационных файлов

### YAML-файлы (единая схема)

Все YAML-файлы будут соответствовать следующей строгой схеме:

```yaml
name: "ServiceX"
description: "Some description"
version: 1
metadata:
  author: "John Doe"
  tags:
    - "example"
    - "config"

```

### JSON-файлы (разные схемы)

JSON-файлы могут иметь **разные** структуры, но известные пути к нужным полям. 

Пример 1:

```json
{
  "name": "ServiceA",
  "description": "Service A description",
  "version": 1,
  "metadata": {
    "author": "Alice",
    "tags": ["tag1", "tag2"]
  }
}
```

Пример 2:

```json
{
  "info": {
    "service_name": "ServiceB",
    "details": {
      "description": "Service B details",
      "version_number": 2
    }
  },
  "owner": {
    "creator": "Bob"
  },
  "labels": ["production", "stable"]
}
```

Приложение должно уметь распознавать JSON схему и извлекать данные по известным путям. Например:
* `name` → `"name"` (в одном формате) / `"info.service_name"` (в другом) / `"app.name"` (в третьем).
* `description` → `"description"` / `"info.details.description"` / `"app.desc"`.
* `version` → `"version"` / `"info.details.version_number"` / `"app.meta.version"`.
* `author` → `"metadata.author"` / `"owner.creator"` / `"app.meta.created_by"`.
* `tags` → `"metadata.tags"` / `"labels"` / `"app.meta.labels"`.

Чтобы это реализовать, приложение должно поддерживать конфиг-файл, в котором задаются известные пути для разных JSON-схем: `mycli -c ./json_paths.yaml -p ./tests/src/test_1`

`json_paths.yaml`
```yaml
json_schemas:
  - name: "name"
    description: "description"
    version: "version"
    author: "metadata.author"
    tags: "metadata.tags"
  - name: "info.service_name"
    description: "info.details.description"
    version: "info.details.version_number"
    author: "owner.creator"
    tags: "labels"
  - name: "app.name"
    description: "app.desc"
    version: "app.meta.version"
    author: "app.meta.created_by"
    tags: "app.meta.labels"
```

## 3. Структура БД (PostgreSQL)

Используем PostgreSQL, таблица configs:

```sql
CREATE TABLE IF NOT EXISTS configs (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    description TEXT,
    version INTEGER,
    author TEXT,
    tags TEXT[]
);
```

При вставке новых данных:
* Если запись с таким `name` уже есть, то выполняем `UPDATE`.
* Если записи нет, выполняем `INSERT`.

Для того чтобы локально развернуть БД в контейнере и подключиться необходимо:
* Установить [docker-compose](https://docs.docker.com/compose/install/)
* Установить [psql](https://www.postgresql.org/docs/current/app-psql.html) (или любой другой удобный клиент) для подключения к postgres
* Запустить команду `make up`
* Подключаемся локально к БД `psql --host=localhost --port=5432 --dbname=postgres --username=user --password`
* Таблица `configs` уже будет создана в нашей БД

## 4. Рекомендованная структура проекта
```
.
├── Makefile
├── README.md
├── cmd
│   └── mycli
│       └── main.go # точка фхода
├── docker-compose.yaml
├── go.mod
├── go.sum
├── init.sql # скрипт с созданием таблицы
├── internal
│   ├── app
│   │   └── processor.go # основная логика приложения
│   ├── config
│   │   └── parse_config.go # основная логика работы с конфигами
│   ├── db
│   │   └── database.go # основная логика приложения
│   └── models
│       └── models.go # общие модели
├── json_paths.yaml # конфиг схем
└── tests # директория с тестовыми данными
```

## 5. Запуск и проверка корректного выполнения

1. `make up` - запускаем БД
2. `psql --host=localhost --port=5432 --dbname=postgres --username=user --password` - подключаемся к БД в другом терминале (ниже бужет обозначен как `postgres=#`)
3. `go run ./cmd/mycli -c ./json_paths.yaml -d ./tests/test_1` - запускаем наш CLI на первом тесте
4. `postgres=#` > `SELECT * FROM postgres.configs ORDER BY name;` - выполним запрос на получение сохраненных конфигов
5. Сравните полученный результат в **п.4** с [RESULT.md](./tests/test_1/RESULT.md)
6. `go run ./cmd/mycli -c ./json_paths.yaml -d ./tests/test_2` - запускаем наш CLI на втором тесте
7. `postgres=#` > `SELECT * FROM postgres.configs ORDER BY name;` - опять выполним запрос на получение сохраненных конфигов
8. Сравните полученный результат в **п.7** с [RESULT.md](./tests/test_2/RESULT.md)
9. `make down`
