include .env

## help: Показать справочную информацию
.PHONY: help
help:
	@echo 'Использование:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## migrations-new name=$1: Создать новые файлы миграции для <name>
.PHONY: migrations-new
migrations-new:
	@echo 'Создание миграций для ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## migrations-up: Применить все миграции
.PHONY: migrations-up
migrations-up:
	@echo 'Применение миграций...'
	@migrate -path=./migrations -database="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost/${POSTGRES_DB}?sslmode=disable" up

## migrations-down: Откатить все миграции
.PHONY: migrations-down
migrations-down:
	@echo 'Откат миграций...'
	@migrate -path=./migrations -database="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost/${POSTGRES_DB}?sslmode=disable" down

## run: Запустить API сервер
.PHONY: run
run:
	@echo 'Запуск сервера API...'
	go run ./cmd/api -db-dsn="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost/${POSTGRES_DB}"

## control: Форматирование кода, проверка кода и обновление зависимостей
.PHONY: control
control:
	@echo 'Форматирование кода...'
	go fmt ./...
	@echo 'Проверка кода...'
	go vet ./...
	@echo 'Обновление и проверка зависимостей...'
	go mod tidy
	go mod verify
