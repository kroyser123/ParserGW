ParserGW
CLI-приложение на Go, которое парсит YAML и JSON конфиги и складывает их в PostgreSQL.

Быстрый старт
bash
# Клонируй
git clone https://github.com/kroyser123/ParserGW.git
cd ParserGW

# Подними базу
make up

# Запусти парсер
go run ./cmd/mycli -c ./json_paths.yaml -d ./tests/test_1

# Проверь что сохранилось
docker exec -it hw4-postgres-1 psql -U user -d postgres -c "SELECT * FROM configs;"
Команды
make up - запустить PostgreSQL в Docker

make down - остановить PostgreSQL

go run ./cmd/mycli -c ./json_paths.yaml -d ./папка/с/файлами - запустить парсер

Формат файлов
YAML:

yaml
name: "ServiceX"
version: 1
metadata:
  author: "John Doe"
  tags: ["example", "config"]
JSON (поддерживает разные схемы через json_paths.yaml).

Что под капотом
Ищет все YAML и JSON файлы в папке (и подпапках)

Парсит их по заданным схемам

Сохраняет в PostgreSQL (если запись с таким именем уже есть - обновляет)

Выводит статистику
